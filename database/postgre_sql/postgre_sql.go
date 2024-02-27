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

type Answer struct {
	Id               int
	Expression_input string
	Number           float64
}

func GetDbConnection() *sql.DB {
	psqlconn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
	db, err := sql.Open("postgres", psqlconn)
	if err != nil {
		fmt.Println(err)
	}
	return db
}

func InsertAnswer(db *sql.DB, expression string, result float64) error {
	_, err := db.Exec("INSERT INTO \"Input\" (\"Expression\", \"Answer\") VALUES ($1, $2)", expression, result)
	return err
}

func GetAnswers(db *sql.DB) ([]Answer, error) {
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
