package main

import (
	"crypto/tls"
	"database/sql"
	"flag"
	"log"
	"net/http"
	"os"
	"text/template"
	"time"

	"blogdemo.batou.cn/common/models"
	"github.com/alexedwards/scs/mysqlstore"
	"github.com/alexedwards/scs/v2"
	"github.com/go-playground/form/v4"
	_ "github.com/go-sql-driver/mysql"
)

type application struct {
	debug          bool
	infoLog        *log.Logger
	errorLog       *log.Logger
	blogs          *models.BlogModel
	users          *models.UserModel
	templateCache  map[string]*template.Template
	formDecoder    *form.Decoder
	sessionManager *scs.SessionManager
}

type contextKey string

const isAuthenticated = contextKey("isAuthenticated")

func main() {
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERR\t", log.Ldate|log.Ltime|log.Lshortfile)

	debug := flag.Bool("debug", false, "是否启用Debug模式")
	flag.Parse()

	// 数据库连接，parseTime指令的作用是将SQL的TIME和DATE字段转为Go的time.Time对象
	dsn := "root:123456@tcp(127.0.0.1:3305)/blog?parseTime=true&loc=Local"

	db, err := dbConnect(dsn)
	if err != nil {
		errorLog.Fatal(err.Error())
	}

	defer db.Close()

	templateCache, err := newTemplateCache()
	if err != nil {
		errorLog.Fatal(err)
	}

	sessionManager := scs.New()
	sessionManager.Store = mysqlstore.New(db)
	sessionManager.Lifetime = 12 * time.Hour

	app := &application{
		debug:          *debug,
		infoLog:        infoLog,
		errorLog:       errorLog,
		blogs:          &models.BlogModel{DB: db},
		users:          &models.UserModel{DB: db},
		templateCache:  templateCache,
		formDecoder:    form.NewDecoder(),
		sessionManager: sessionManager,
	}

	tlsConfig := &tls.Config{
		CurvePreferences: []tls.CurveID{tls.X25519, tls.CurveP256},
	}

	srv := &http.Server{
		Addr:      ":4000",
		Handler:   app.routes(),
		ErrorLog:  errorLog,
		TLSConfig: tlsConfig,
		// 配置服务超时信息，提高可用性
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	infoLog.Print("启动web服务，端口4000")

	// 本地测试用的密钥文件可以通过 Go 安装目录下的 generate_cert.go 文件创建
	err = srv.ListenAndServeTLS("./tls/cert.pem", "./tls/key.pem")
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
