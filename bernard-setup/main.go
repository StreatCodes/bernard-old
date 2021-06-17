package main

import (
	"fmt"
	"log"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	"github.com/streatcodes/bernard/bernard-server/db"
)

func main() {
	var email, password string
	fmt.Print("Enter an email for the admin user: ")
	fmt.Scanf("%s", &email)

	fmt.Print("Enter the password for the new user: ")
	fmt.Scanf("%s", &password)

	//Create database
	fmt.Println("Writing database")
	database, err := sqlx.Open("sqlite3", "bernard.db")
	if err != nil {
		log.Fatalf("Error opening DB: %s\n", err)
	}
	database.MustExec(`PRAGMA foreign_keys;`)
	database.MustExec(db.Schema)

	//Insert admin user
	_, err = db.CreateUser(database, email, password, db.RoleManageUsers|db.RoleMangeHosts)
	if err != nil {
		log.Fatalf("Error creating user: %s\n", err)
	}

	fmt.Println("User and database created")
}
