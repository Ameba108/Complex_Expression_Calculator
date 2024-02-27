package postgreSql

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "1234"
	dbname   = "Expressions"
)

// структура, благодаря которой мы сможем вывести результаты расчетов из базы данных
type Answer struct {
	Id               int
	Expression_input string
	Number           float64
}

func GetDbConnection() *sql.DB {
	//открываем базу данных
	psqlconn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
	db, err := sql.Open("postgres", psqlconn)
	if err != nil {
		//что-то пошло не так - выводим в консоль сообщение об ошибке
		fmt.Println(err)
	}
	return db
}

// функция, которая будет отправлять вводимое выражение и ответ вычисления этого выражения в базу данных
func InsertAnswer(db *sql.DB, expression string, result float64) error {
	//отправляем все в базу данных
	_, err := db.Exec("INSERT INTO \"Input\" (\"Expression\", \"Answer\") VALUES ($1, $2)", expression, result)
	return err
}

// выводим на страничке "История" все сохраненные данные из БД
func GetAnswers(db *sql.DB) ([]Answer, error) {
	//выводим все из БД
	rows, err := db.Query("SELECT \"Id\", \"Expression\", \"Answer\" FROM \"Input\"")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var answers []Answer
	for rows.Next() {
		var answer Answer
		err := rows.Scan(&answer.Id, &answer.Expression_input, &answer.Number)
		if err != nil {
			return nil, err
		}
		answers = append(answers, answer)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return answers, nil
}
