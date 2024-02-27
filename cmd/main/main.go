package main

import (
	//импортируем пакет с сервером
	"complex_expression_calculator/http/server"
	"net/http"
)

// функция запуска калькулятора
func main() {
	//подключаем css шаблон
	fs := http.FileServer(http.Dir("http/templates"))
	http.Handle("/http/templates/", http.StripPrefix("/http/templates/", fs))
	//запускаем сервер! :D
	server.HandleRequest()
}
