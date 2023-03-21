package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

const VERSION = "1.0.0"

type config struct {
	port string
	env  string
	db   struct {
		uri      string
		database string
	}

	smtp struct {
		prov     string
		sport    int
		user     string
		password string
	}

	personal struct {
		provider  string
		pport     int
		puser     string
		ppassword string
	}
}

type application struct {
	config   config
	infoLog  *log.Logger
	errorLog *log.Logger
}

func (app *application) StartServer() error {
	srv := &http.Server{
		Addr:              fmt.Sprintf(":%v", app.config.port),
		Handler:           app.Router(),
		IdleTimeout:       30 * time.Second,
		ReadTimeout:       30 * time.Second,
		ReadHeaderTimeout: 30 * time.Second,
		WriteTimeout:      30 * time.Second,
	}

	app.infoLog.Printf("API starting at PORT %s in %s mode", app.config.port, app.config.env)
	return srv.ListenAndServe()
}

func main() {

	var c config

	port := os.Getenv("PORT")

	if port == "" {
		port = "8005"
	}

	flag.StringVar(&c.env, "env", "DEVELOPMENT", "{DEVELOPMENT,MAINTAINANCE,PRODUCTION}")
	flag.StringVar(&c.port, "API port", port, "PORT OF API")
	flag.StringVar(&c.db.uri, "URI", "mongodb+srv://manuel:adminadmin@ismandb.mcjzh2z.mongodb.net/test", "URI OF THE DATABASE")
	flag.StringVar(&c.db.database, "DB DATABASE", "main", "DATABASE NAME")
	flag.StringVar(&c.smtp.prov, "SMTP PROVIDER", "ns148.hostgator.mx", "SMTP PROVIDER")
	flag.IntVar(&c.smtp.sport, "SMTP PORT", 465, "SMTP PORT")
	flag.StringVar(&c.smtp.user, "SMTP USER", "contacto@ismanpublicidad.com", "SMTP USER")
	flag.StringVar(&c.smtp.password, "SMTP PASSWORD", "ismanadminadmin", "SMTP PASSWORD")
	flag.StringVar(&c.personal.provider, "PERSONAL PROVIDER", "smtp.gmail.com", "SMTP PROVIDER")
	flag.IntVar(&c.personal.pport, "PERSONAL PORT", 587, "SMTP PORT")
	flag.StringVar(&c.personal.puser, "PERSONAL USER", "ismanpublicidad@gmail.com", "SMTP USER")
	flag.StringVar(&c.personal.ppassword, "PERSONAL PASSWORD", "ismanpublicidad12345", "SMTP PASSWORD")
	flag.Parse()

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime|log.Lshortfile)
	errorLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime|log.Lshortfile)

	app := &application{
		config:   c,
		infoLog:  infoLog,
		errorLog: errorLog,
	}

	err := app.StartServer()
	if err != nil {
		app.errorLog.Fatal(err.Error())
		return
	}

}
