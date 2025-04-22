package main

import (
	"fmt"
	"log"
	"net/http"

	"practice/library-system-v3/handler"

	"github.com/gorilla/schema"
	// "github.com/gorilla/securecookie"
	"github.com/gorilla/sessions"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func main() {
	var createTable = `
			CREATE TABLE IF NOT EXISTS books (
				id serial,
				book_name text,
				author text,
				category text,
				book_description text,
				book_cover text,
				is_available boolean,

				primary key(id)
			);`

	var createUserTable = `
		CREATE TABLE IF NOT EXISTS users (
			user_id serial,
			first_name text,
			last_name text,
			email text,
			password text,
			is_verified boolean,

			primary key(user_id)
		);`

	db, err := sqlx.Connect("postgres", "user=postgres password=Anubis0912 dbname=library_v3 sslmode=disable")
	if err != nil {
		log.Fatalln(err)
	}
	db.MustExec(createTable)
	db.MustExec(createUserTable)

	decoder := schema.NewDecoder()
	decoder.IgnoreUnknownKeys(true)
	store := sessions.NewCookieStore([]byte("op0Rg9Wct9j2Po3EEmiZbjW6rVwPypsC"))
	r := handler.GetHandler(db, decoder, store)

	fmt.Println("Server Starting...")
	err = http.ListenAndServe(":8080", r)
	if err != nil {
		log.Fatal("Server Not Found", err)
	}
}
