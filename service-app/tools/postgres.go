package main

import (
	"context"
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"log"
	"time"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "diwakar"
	password = "root"
	dbname   = "postgres"
)

var db *sql.DB

func main() {
	var err error
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	db, err = sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	if err := db.Ping(); err != nil {
		panic(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	//Insert(ctx)
	//Insert2(ctx)
	//delete(ctx)
	//querySingleRecords(ctx)
	QueryMultipleRecords(ctx)
}

func Insert(ctx context.Context) {
	sqlStatement := `INSERT INTO users (age, email, first_name,last_name)
					VALUES ($1, $2, $3, $4)`
	//db.Exec()
	res, err := db.ExecContext(ctx, sqlStatement, 33, "abc@email.com", "dev", "kumar") // for zero to one values getting back from the query

	if err != nil {
		log.Fatalln(err)
	}
	log.Println(res.LastInsertId()) // this is not supported

}
func Insert2(ctx context.Context) {

	sqlStatement := `INSERT INTO users (age, email, first_name,last_name)
					VALUES ($1, $2, $3, $4)
					RETURNING id,email
					`
	var (
		id    int
		email string
	)
	err := db.QueryRowContext(ctx, sqlStatement, 34, "abc11@email.com", "dev11", "kumar1").Scan(&id, &email)

	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println(id, email)

}
func delete(ctx context.Context) {

	sqlStatement := `Delete FROM users
                     where id =$1;
`

	res, err := db.ExecContext(ctx, sqlStatement, 2)
	if err != nil {
		log.Println(err)
		return
	}

	fmt.Println(res.RowsAffected())

}

func querySingleRecords(ctx context.Context) {

	sqlStatement := `Select id, email FROM users where id = $1;`

	var (
		id    int
		email string
	)

	err := db.QueryRowContext(ctx, sqlStatement, 1).Scan(&id, &email)

	if err != nil {
		log.Println(err)

	}

	switch err {
	case sql.ErrNoRows:
		log.Println("no rows returned")
	case nil:
		fmt.Println(id, email)
	default:
		log.Println(err)

	}

}

func QueryMultipleRecords(ctx context.Context) {

	rows, err := db.QueryContext(ctx, "Select id, email FROM users LIMIT $1", 4)
	if err != nil {
		log.Fatalln(err)
	}

	defer rows.Close()

	for rows.Next() {
		var (
			id    int
			email string
		)

		err = rows.Scan(&id, &email)
		if err != nil {
			log.Println(err)
		}

		fmt.Println(id, email)

	}

	if err := rows.Err(); err != nil {
		log.Fatalln(err)
	}

}
