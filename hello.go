package main

import (
	"database/sql"
	"fmt"
    
    _ "github.com/go-sql-driver/mysql"
)

type user struct {
	Name string
	Email string
	Pass string
}

func main() {
	dbDriver := "mysql"
    dbUser := "root"
    dbPass := "root"
    dbName := "users"
    db, err := sql.Open(dbDriver, dbUser+":"+dbPass+"@tcp(127.0.0.1:3306)/"+dbName)
    if err != nil {
        panic(err.Error())
    }else{
		fmt.Println("Database accessed")
		u := user{}
		selDb,err := db.Query("SELECT * FROM info")
		if err != nil {
			panic(err.Error())
		}
		for selDb.Next(){
			var name, email, pass string
			err1 := selDb.Scan(&name,&email,&pass)
			if err1 != nil {
				panic(err1.Error())
			}
			u.Name = name
			u.Email = email
			u.Pass = pass
			fmt.Println(u.Name+" "+u.Email+" "+u.Pass)
		}
	}
	
}