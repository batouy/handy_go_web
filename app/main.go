package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	"blogdemo.batou.cn/common/models"
	_ "github.com/go-sql-driver/mysql"
)

type application struct {
	infoLog  *log.Logger
	errorLog *log.Logger
	blogs    *models.BlogModel
}

func main() {
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERR\t", log.Ldate|log.Ltime|log.Lshortfile)

	// 数据库连接，parseTime指令的作用是将SQL的TIME和DATE字段转为Go的time.Time对象
	dsn := "root:123456@tcp(127.0.0.1:3305)/blog?parseTime=true"

	db, err := dbConnect(dsn)
	if err != nil {
		errorLog.Fatal(err.Error())
	}

	defer db.Close()

	app := &application{
		infoLog:  infoLog,
		errorLog: errorLog,
		blogs:    &models.BlogModel{DB: db},
	}

	srv := &http.Server{
		Addr:     ":4000",
		Handler:  app.routes(),
		ErrorLog: errorLog,
	}

	infoLog.Print("启动web服务，端口4000")
	err = srv.ListenAndServe()
	errorLog.Fatal(err)
}

func dbConnect(dsn string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn)

	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}
