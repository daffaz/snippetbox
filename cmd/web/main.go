package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	"github.com/daffaz/snippetbox/pkg/models/mysql"
	_ "github.com/go-sql-driver/mysql"
)

type application struct {
	infoLog  *log.Logger
	errorLog *log.Logger
	snippets *mysql.SnippetModel
}

func main() {
	infoLog := log.New(os.Stdout, "[INFO]\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "[ERROR]\t", log.Ldate|log.Ltime|log.Lshortfile)

	dsn := "CONFIDENTAL:CONFIDENTAL@/CONFIDENTAL?parseTime=true"

	db, err := openDB(dsn)
	if err != nil {
		errorLog.Fatal(err)
	}

	defer db.Close()

	app := &application{
		infoLog:  infoLog,
		errorLog: errorLog,
		snippets: &mysql.SnippetModel{DB: db},
	}

	srv := http.Server{
		Addr:     ":4000",
		Handler:  app.routes(),
		ErrorLog: errorLog,
	}
	infoLog.Printf("Running in port %s\n", srv.Addr)
	errorLog.Fatal(srv.ListenAndServe())
}

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}
