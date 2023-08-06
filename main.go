package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/mattn/go-sqlite3"
	"github.com/mtstnt/ioc-go-experiment/controllers"
	"github.com/mtstnt/ioc-go-experiment/ioc"
	"github.com/mtstnt/ioc-go-experiment/services"
)

func getDB() *sql.DB {
	v, err := sql.Open("sqlite3", "./db")
	if err != nil {
		panic(err)
	}
	return v
}

func main() {
	c := ioc.NewContainer()
	c.Register(services.UserServiceFactory)
	c.Register(getDB)

	v, err := ioc.Inject[controllers.User](c, &controllers.User{})
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println(v.GetA())
	fmt.Println(v.MyDB.Query("CREATE TABLE t (id INTEGER)"))
}
