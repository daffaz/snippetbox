package main

import (
	"crypto/tls"
	"database/sql"
	"html/template"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/daffaz/snippetbox/pkg/models/mysql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/golangcollege/sessions"
)

type application struct {
	infoLog       *log.Logger
	errorLog      *log.Logger
	session       *sessions.Session
	snippets      *mysql.SnippetModel
	templateCache map[string]*template.Template
}

func main() {
	infoLog := log.New(os.Stdout, "[INFO]\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "[ERROR]\t", log.Ldate|log.Ltime|log.Lshortfile)

	dsn := "CONFIDENTAL:CONFIDENTAL@/CONFIDENTAL?parseTime=true"
	secret := "s6Ndh+pPbnzHbS*+9Pk8qGWhTzbpa@ge"

	db, err := openDB(dsn)
	if err != nil {
		errorLog.Fatal(err)
	}

	defer db.Close()
	templateCache, err := newTemplateCache("./ui/html/")
	if err != nil {
		errorLog.Fatal(err)
	}

	session := sessions.New([]byte(secret))
	session.Lifetime = 12 * time.Hour
	session.Secure = true

	app := &application{
		infoLog:       infoLog,
		errorLog:      errorLog,
		session:       session,
		snippets:      &mysql.SnippetModel{DB: db},
		templateCache: templateCache,
	}

	tlsConfig := &tls.Config{
		PreferServerCipherSuites: true,
		CurvePreferences:         []tls.CurveID{tls.X25519, tls.CurveP256},
	}

	srv := http.Server{
		Addr:      ":4000",
		Handler:   app.routes(),
		ErrorLog:  errorLog,
		TLSConfig: tlsConfig,
	}

	infoLog.Printf("Running in port %s\n", srv.Addr)
	errorLog.Fatal(srv.ListenAndServeTLS("./tls/cert.pem", "./tls/key.pem"))
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
