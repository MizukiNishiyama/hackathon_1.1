package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sourse/controller"
	"sourse/dao"
	"sourse/usecase"
	"syscall"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/oklog/ulid/v2"
)

var db *sql.DB

func main() {
	db = initDB()

	userDao := &dao.UserDao{DB: db}
	searchUserController := &controller.SearchUserController{SearchUserUseCase: &usecase.SearchUserUseCase{UserDao: userDao}}
	registerUserController := &controller.RegisterUserController{RegisterUserUseCase: &usecase.RegisterUserUseCase{UserDao: userDao}}
	http.HandleFunc("/user", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			searchUserController.Handle(w, r)
		case http.MethodPost:
			registerUserController.Handle(w, r)
		default:
			log.Printf("BadRequest(status code = 400)")
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	})

	closeDBWithSysCall()

	log.Println("Listening...")
	if err := http.ListenAndServe(":8000", nil); err != nil {
		log.Fatal(err)
	}
}

func initDB() *sql.DB {
	// DB接続のための準備
	//mysqlUser := os.Getenv("MYSQL_USER")
	//mysqlPwd := os.Getenv("MYSQL_PWD")
	//mysqlHost := os.Getenv("MYSQL_HOST")
	//mysqlDatabase := os.Getenv("MYSQL_DATABASE")

	mysqlUser := "root"
	mysqlPwd := "ramen102"
	mysqlHost := "34.172.193.162:3306"
	mysqlDatabase := "hackathon"

	connStr := fmt.Sprintf("%s:%s@tcp(%s)/%s", mysqlUser, mysqlPwd, mysqlHost, mysqlDatabase)
	fmt.Println(connStr)
	_db, err := sql.Open("mysql", connStr)

	if err != nil {
		log.Fatalf("fail: sql.Open, %v\n", err)
	}
	if err := _db.Ping(); err != nil {
		log.Fatalf("fail: _db.Ping, %v\n", err)
	}
	return _db
}

func closeDBWithSysCall() {
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGTERM, syscall.SIGINT)
	go func() {
		s := <-sig
		log.Printf("received syscall, %v", s)

		if err := db.Close(); err != nil {
			log.Fatal(err)
		}
		log.Printf("success: db.Close()")
		os.Exit(0)
	}()
}
