package calculator

import (
	"fmt"
	"strconv"
	"strings"
	"sync"
	"unicode"
)

// реализация стека
type Stack struct {
	mu           sync.Mutex   //мьютекс, чтобы не было проблем работы со стеком в случае конкурентности
	i            float64      //индекс стека
	data         [100]float64 //массив с фиксированным размером
	Input_string string       //вводимая строка (например "2+2*2" и тд)
	Result       string       //результат (здесь в роли результата является строка, преобразованная в постфикс)
	//то есть резальтатом "2 + 2" будет "2 2 +"
}

// реализуем стек
func NewStack(n string) *Stack {
	return &Stack{Input_string: n}
}

// расставляем приоритеты знаков, чтобы в случае умножения/деления эти операции выполнялись первыми
func precedence(op string) float64 {
	switch op {
	case "+", "-":
		return 1
	case "*", "/":
		return 2
	default:
		return 0
	}
}

// функция, которая преобразует строку в постфикс
func (s *Stack) InfixToPostfix() {
	s.mu.Lock()
	defer s.mu.Unlock()
	var number strings.Builder
	var output strings.Builder
	var stack []string
	//функция пробегается по строке, определяя знаки операции и числа
	for i, char := range s.Input_string {
		switch {
		//в случае, если число с точкой (десятичная дробь)
		case unicode.IsDigit(char) || char == '.':
			number.WriteRune(char)
			// Если следующий символ - это оператор или скобка, добавляем текущее число в выходную строку
			if i+1 < len(s.Input_string) && (unicode.IsDigit(rune(s.Input_string[i+1])) || s.Input_string[i+1] == '.' || s.Input_string[i+1] == '(' || s.Input_string[i+1] == ')' || s.Input_string[i+1] == '+' || s.Input_string[i+1] == '-' || s.Input_string[i+1] == '*' || s.Input_string[i+1] == '/') {
				continue
			}
			output.WriteString(number.String())
			output.WriteRune(' ')
			number.Reset()
			//находим открывающуюся скобку и добавляем ее в стек
		case char == '(':
			stack = append(stack, string(char))
		case char == ')':
			for len(stack) > 0 && stack[len(stack)-1] != "(" {
				output.WriteString(stack[len(stack)-1])
				output.WriteRune(' ')
				stack = stack[:len(stack)-1]
			}
			stack = stack[:len(stack)-1] // Удаляем открывающую скобку
		default: // Оператор или пробел
			//проверяем на наличие символов оператора или пробела (должно быть больше нуля)
			if number.Len() > 0 {
				//если находим - отправляем в output
				output.WriteString(number.String())
				output.WriteRune(' ')
				number.Reset()
			}
			//если символ char не является пробелом, добавляем символы (операции) в зависимости от приоритета в output
			if char != ' ' {
				for len(stack) > 0 && precedence(stack[len(stack)-1]) >= precedence(string(char)) {
					output.WriteString(stack[len(stack)-1])
					output.WriteRune(' ')
					stack = stack[:len(stack)-1]
				}
				//добавляем char в стек
				stack = append(stack, string(char))
			}
		}
	}

	// Если в конце строки осталось число, добавляем его в выходную строку
	if number.Len() > 0 {
		output.WriteString(number.String())
		output.WriteRune(' ')
	}
	//получаем результат
	for len(stack) > 0 {
		output.WriteString(stack[len(stack)-1])
		output.WriteRune(' ')
		stack = stack[:len(stack)-1]
	}

	s.Result = output.String()
}

// функция, добавляющая данные в стек
func (s *Stack) Push(i float64) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.i+1 >= float64(len(s.data)) {
		fmt.Println("Стек заполнен!")
	}

	s.data[int(s.i)] = i
	s.i++
}

// функция, удаляющая данные из стека
func (s *Stack) Pop() float64 {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.i-1 < 0 {
		fmt.Println("Нет данных для удаления!")
	}

	s.i--
	return s.data[int(s.i)]
}

// получения размера стека
func (s *Stack) Size() int {
	return int(s.i)
}

// функция калькулятора
// не просто так была создана функция, переводящая входную строку в постфикс
// калькулятор реализован с помощью обратной польской аннотации
func Calculate(input string) float64 {
	//создаем стек
	stack := new(Stack)

	for i := 0; i < len(input); i++ {
		c := rune(input[i])
		//если c является числом - добавляем в стек
		if unicode.IsDigit(c) {
			stack.Push(parseNumber(input, (&i)))
		}
		switch c {
		//если введено неверное выражение (например: 2+ ), то выводим сообщение об ошибке в консоль
		case '+', '-', '*', '/':
			if stack.Size() < 2 {
				fmt.Println("Не хватает значений для операции в стеке " + string(c))
			}
			//Извлекаем два числа из стека
			b := stack.Pop()
			a := stack.Pop()
			//считаем, в зависимости от операции и отправляем полученный результат в стек
			switch c {
			case '+':
				stack.Push(a + b)
			case '-':
				stack.Push(a - b)
			case '*':
				stack.Push(a * b)
			case '/':
				if b == 0 {
					//при делении на ноль выводим сообщение об ошибке в консоль
					fmt.Println("Деление на ноль!")
				}
				stack.Push(a / b)
			}
		}
	}
	//если стек пуст, выводим сообщение об ошибке в консоль
	if stack.Size() != 1 {
		fmt.Println("Стек недостаточно заполнен!")
	}
	//удаляем из стека элементы, чтобы не было ошибок
	//например, чтобы не повторялось одно действие несколько раз подряд
	return stack.Pop()
}

// функция получения числа
func parseNumber(input string, i *int) float64 {
	find_number := *i

	//Найти конец числа, учитывая точку
	for *i < len(input) && (unicode.IsDigit(rune(input[*i])) || rune(input[*i]) == '.') {
		*i++
	}
	num, _ := strconv.ParseFloat(input[find_number:*i], 32)
	return num
}
