package calculator

import (
	"fmt"
	"strconv"
	"strings"
	"sync"
	"unicode"
)

type Stack struct {
	mu           sync.Mutex
	i            float64
	data         [100]float64
	Input_string string
	Result       string
}

func NewStack(n string) *Stack {
	return &Stack{Input_string: n}
}

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

func (s *Stack) InfixToPostfix() {
	s.mu.Lock()
	defer s.mu.Unlock()
	var number strings.Builder
	var output strings.Builder
	var stack []string
	for i, char := range s.Input_string {
		switch {
		case unicode.IsDigit(char) || char == '.':
			number.WriteRune(char)
			// Если следующий символ - это оператор или скобка, добавляем текущее число в выходную строку
			if i+1 < len(s.Input_string) && (unicode.IsDigit(rune(s.Input_string[i+1])) || s.Input_string[i+1] == '.' || s.Input_string[i+1] == '(' || s.Input_string[i+1] == ')' || s.Input_string[i+1] == '+' || s.Input_string[i+1] == '-' || s.Input_string[i+1] == '*' || s.Input_string[i+1] == '/') {
				continue
			}
			output.WriteString(number.String())
			output.WriteRune(' ')
			number.Reset()
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
			if number.Len() > 0 {
				output.WriteString(number.String())
				output.WriteRune(' ')
				number.Reset()
			}
			if char != ' ' {
				for len(stack) > 0 && precedence(stack[len(stack)-1]) >= precedence(string(char)) {
					output.WriteString(stack[len(stack)-1])
					output.WriteRune(' ')
					stack = stack[:len(stack)-1]
				}
				stack = append(stack, string(char))
			}
		}
	}

	// Если в конце строки осталось число, добавляем его в выходную строку
	if number.Len() > 0 {
		output.WriteString(number.String())
		output.WriteRune(' ')
	}

	for len(stack) > 0 {
		output.WriteString(stack[len(stack)-1])
		output.WriteRune(' ')
		stack = stack[:len(stack)-1]
	}

	s.Result = output.String()
}

func (s *Stack) Push(i float64) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.i+1 >= float64(len(s.data)) {
		fmt.Println("Стек заполнен!")
	}

	s.data[int(s.i)] = i
	s.i++
}

func (s *Stack) Pop() float64 {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.i-1 < 0 {
		fmt.Println("Нет данных для удаления!")
	}

	s.i--
	return s.data[int(s.i)]
}

func (s *Stack) Size() int {
	return int(s.i)
}

func Calculate(input string) float64 {
	stack := new(Stack)

	for i := 0; i < len(input); i++ {
		c := rune(input[i])

		if unicode.IsDigit(c) {
			stack.Push(parseNumber(input, (&i)))
		}
		switch c {
		case '+', '-', '*', '/':
			if stack.Size() < 2 {
				fmt.Println("Не хватает значений для операции в стеке " + string(c))
			}
			// Теперь можно безопасно извлечь два числа
			b := stack.Pop()
			a := stack.Pop()
			switch c {
			case '+':
				stack.Push(a + b)
			case '-':
				stack.Push(a - b)
			case '*':
				stack.Push(a * b)
			case '/':
				if b == 0 {
					fmt.Println("Деление на ноль!")
				}
				stack.Push(a / b)
			}
		}
	}

	if stack.Size() != 1 {
		fmt.Println("Стек недостаточно заполнен!")
	}

	return stack.Pop()
}

func parseNumber(input string, i *int) float64 {
	find_number := *i

	// Найти конец числа, учитывая точку
	for *i < len(input) && (unicode.IsDigit(rune(input[*i])) || rune(input[*i]) == '.') {
		*i++
	}
	num, _ := strconv.ParseFloat(input[find_number:*i], 32)
	return num
}
