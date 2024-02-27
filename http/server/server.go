package server

import (
	//импортируем пакеты с postgreSQL и функцией калькулятора
	postgreSql "complex_expression_calculator/database/postgre_sql"
	"complex_expression_calculator/pkg/calculator"
	"fmt"
	"html/template"

	"net/http"
)

// структура Answer, в которой мы будем хранить само вводимое выражение (expression_input) и ответ выражения (number)
type Answer struct {
	Id               int
	Expression_input string
	Number           float64
}

// реализуем структуру Answer
var answer = Answer{}

// главная страница
func index_page(w http.ResponseWriter, r *http.Request) {
	//подключаем базу данных, чтобы передавать в нее данные
	db := postgreSql.GetDbConnection()
	defer db.Close()

	//проверяем метод запроса. Если запрос POST - наш калькулятор начинает работать и считать полученное выражение
	if r.Method == "POST" {
		//получаем вводимое выражение
		expression := r.FormValue("data")
		//создаем стек
		stack := calculator.NewStack(expression)
		//сохраняем в структуре Answer вводимую строку
		answer.Expression_input = expression
		//переводим строку в вид постфикс, чтобы калькулятор смог посчитать
		stack.InfixToPostfix()
		//считаем с поомощью функции Calculate
		result := calculator.Calculate(stack.Result)
		//сохраняем ответ в структуре Answer
		answer.Number = result
		//если возникла ошибка с переносом данных в БД - выводим сообщение ошибки в консоль
		err := postgreSql.InsertAnswer(db, expression, result)
		if err != nil {
			fmt.Println(err)
		}
		//выводим результат в консоль (я сделала этот принт по своим нуждам, чтобы проверять работу сервера и калькулятора)
		//в принципе, этот принт не нужен, ведь результат и так выводится на страничке сервера
		fmt.Println(result)
	} else {
		//в любом другом случае выводим в консоль сообщение о том, что что-то случилось
		fmt.Println("Что-то пошло не так!")
	}
	//подключение шаблонов
	tmplate, err := template.ParseFiles("http/templates/index_page.html")
	if err != nil {
		//если возникли проблемы с шаблонами, выводим сообщение об ошибке
		fmt.Println("Возникли проблемы с шаблонами")
	}
	tmplate.Execute(w, answer)
}

// страница истории запросов
func history_page(w http.ResponseWriter, r *http.Request) {
	//подключаем базу данных
	db := postgreSql.GetDbConnection()
	defer db.Close()
	//выводим всю историю запросов
	answers, err := postgreSql.GetAnswers(db)
	if err != nil {
		fmt.Println(err)
	}
	//подключение шаблонов
	tmplate, err := template.ParseFiles("http/templates/history_page.html")
	if err != nil {
		fmt.Println("Возникли проблемы с шаблонами")
		return
	}
	tmplate.Execute(w, answers)
}

// определяем пути
func HandleRequest() {
	http.HandleFunc("/calc", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "http/templates/index_page.html")
	})
	http.HandleFunc("/", index_page)

	http.HandleFunc("/history", history_page)
	http.ListenAndServe(":8080", nil)
}
