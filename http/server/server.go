package server

import (
	postgreSql "complex_expression_calculator/database/postgre_sql"
	"complex_expression_calculator/pkg/calculator"
	"fmt"
	"html/template"

	"net/http"
)

type Answer struct {
	Id               int
	Expression_input string
	Number           float64
}

var answer = Answer{}

func index_page(w http.ResponseWriter, r *http.Request) {
	db := postgreSql.GetDbConnection()
	defer db.Close()

	if r.Method == "POST" {
		expression := r.FormValue("data")
		stack := calculator.NewStack(expression)
		answer.Expression_input = expression
		stack.InfixToPostfix()
		result := calculator.Calculate(stack.Result)
		answer.Number = result
		err := postgreSql.InsertAnswer(db, expression, result)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(result)
	} else {
		fmt.Println("Что-то пошло не так!")
	}

	tmplate, err := template.ParseFiles("http/templates/index_page.html")
	if err != nil {
		fmt.Println("Возникли проблемы с шаблонами")
	}
	tmplate.Execute(w, answer)
}

func history_page(w http.ResponseWriter, r *http.Request) {
	db := postgreSql.GetDbConnection()
	defer db.Close()

	answers, err := postgreSql.GetAnswers(db)
	if err != nil {
		fmt.Println(err)
	}

	tmplate, err := template.ParseFiles("http/templates/history_page.html")
	if err != nil {
		fmt.Println("Возникли проблемы с шаблонами")
		return
	}
	tmplate.Execute(w, answers)
}

func HandleRequest() {
	http.HandleFunc("/calc", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "http/templates/index_page.html")
	})
	http.HandleFunc("/", index_page)

	http.HandleFunc("/history", history_page)
	http.ListenAndServe(":8080", nil)
}
