package app

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	"forum/internal/controller"
	"forum/internal/repository"
	"forum/internal/service"

	_ "github.com/mattn/go-sqlite3"
)

const PORT = ":8080"

func Run() {
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	DB, _ := sql.Open("sqlite3", "forum.db")
	defer DB.Close()
	repo := repository.NewRepository(DB)
	service := service.NewService(repo)
	handler := &controller.Handler{ErrorLog: errorLog, Service: service}
	err := repository.CreateDB(DB)
	if err != nil {
		errorLog.Println(err)
		return
	}
	srv := &http.Server{
		Addr:     PORT,
		ErrorLog: errorLog,
		Handler:  controller.Routes(handler),
	}
	infoLog.Printf("listening on http://localhost" + PORT)
	err = srv.ListenAndServe()
	errorLog.Fatal(err)
}
