package main

import (
	"context"
	"database/sql"
	"errors"
	"flag"
	"fmt"

	"go-getting-started/internal/handler"
	"go-getting-started/internal/service"

	"log"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"

	r "go-getting-started/internal/rss"

	// m "go-getting-started/internal/common"
	m "go-getting-started/internal/common"
)

// const (
// 	userDB     = "root3"
// 	passwordDB = "Aa123456$"

// 	redisAddress = "localhost:6379"
// 	redisPassword = ""
// 	//port       = 8080
// )

var (
	port      = env("PORT", "8080")
	originStr = env("ORIGIN", "http://localhost:"+port)
	//todo remove
	databaseURL      = env("DATABASE_URL", m.DatabaseURL)
	databaseURL_read = env("DATABASE_URL", m.DatabaseURL_read)

	redisAddress  = env("REDIS_ADDRESS", m.RedisAddress)
	redisPassword = env("REDIS_PASSWORD", m.RedisPassword)

	// docker version
	// databaseURL      = env("JAWSDB_URL", "root:toor@tcp(alpha:3306)/socialnet?parseTime=true")
	// databaseURL_read = env("JAWSDB_URL", "root:toor@tcp(slave:3306)/socialnet?parseTime=true")

	tokenKey = env("TOKEN_KEY", "supersecretkeyyoushouldnotcommit")

	// smtpHost     = env("SMTP_HOST", "smtp.mailtrap.io")
	// smtpPort     = env("SMTP_PORT", "25")
	// smtpUsername = mustEnv("SMTP_USERNAME")
	// smtpPassword = mustEnv("SMTP_PASSWORD")
)

func main() {
	godotenv.Load()
	if err := run(); err != nil {
		log.Fatalln(err)
	}
}

func run() error {

	//if databaseURL != "root3:Aa123456$@/socialnet" {
	// if databaseURL != "root3:Aa123456$@tcp(alpha:1441)/socialnet" {
	// 	databaseURL += `?useUnicode=true&characterEncoding=utf-8&reconnect=true`
	// }

	allowedOrigin := os.Getenv("ALLOWED_ORIGIN")
	var useNats bool
	flag.StringVar(&allowedOrigin, "allowed-origin", allowedOrigin, "Allowed origin to do requests to this server. If empty, anyone will have access")
	flag.BoolVar(&useNats, "nats", false, "Whether use nats")
	flag.Parse()

	origin, err := url.Parse(originStr)
	if err != nil || !origin.IsAbs() {
		return errors.New("invalid origin url")
	}

	time.Sleep(14 * time.Second)
	//////////////////////////////////////////

	rss := r.New()
	rss.Init(redisAddress, redisPassword)

	/////////////////////////////////////
	//db part
	log.Println(databaseURL)
	log.Println("-------------------  before database opening")
	db, err := sql.Open("mysql", databaseURL)
	log.Println("--------------------after database opening")

	if err != nil {
		log.Println(" ---------------------  ERROR database opening")
		panic(err)
	}

	db_read, err := sql.Open("mysql", databaseURL_read)

	if err != nil {
		log.Println(" ---------------------  ERROR database  opening for only reading")
		panic(err)
	}

	err = db.Ping()
	if err != nil {
		log.Println(" ---------------------  ERROR database PING")
		panic(err)
	}

	err = db_read.Ping()
	if err != nil {
		log.Println(" ---------------------  ERROR database PING for only reading")
		panic(err)
	}

	// See "Important settings" section.
	db.SetConnMaxLifetime(time.Minute * 3)
	db.SetMaxOpenConns(1000)
	db.SetMaxIdleConns(100)

	defer db.Close()

	// See "Important settings" section.
	db_read.SetConnMaxLifetime(time.Minute * 3)
	db_read.SetMaxOpenConns(1000)
	db_read.SetMaxIdleConns(100)

	defer db_read.Close()

	// if err = db.Ping(); err != nil {
	// 	return fmt.Errorf("could not ping to db: %v", err)
	// }
	log.Println("database opened successfully")
	service := service.New(
		db,
		db_read,
		rss,
		//transport,
		//sender,
		*origin,
		tokenKey,
	)
	server := http.Server{
		Addr:              ":" + port,
		Handler:           handler.New(service, origin.Hostname() == "localhost"),
		ReadHeaderTimeout: time.Second * 5,
		ReadTimeout:       time.Second * 15,
	}

	errs := make(chan error, 2)

	go func() {
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, os.Interrupt, os.Kill)

		<-quit

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
		defer cancel()
		if err := server.Shutdown(ctx); err != nil {
			errs <- fmt.Errorf("could not shutdown server: %v", err)
			return
		}

		errs <- ctx.Err()
	}()

	go func() {
		log.Printf("accepting connections on port %s\n", port)
		log.Printf("starting server at %s\n", origin)
		if err = server.ListenAndServe(); err != http.ErrServerClosed {
			errs <- fmt.Errorf("could not listen and serve: %v", err)
			return
		}
		errs <- nil
	}()

	return <-errs
}

func env(key, fallbackValue string) string {
	s, ok := os.LookupEnv(key)
	if !ok {
		return fallbackValue
	}
	return s
}
