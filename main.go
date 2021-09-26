// main.go

package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

const VERSION = "0.1.7"

func main() {
	pwd := os.Getenv("DB_PASSWORD")
	host := os.Getenv("DB_HOST")
	db, err := sql.Open("mysql", fmt.Sprintf("root:%s@tcp(%s)/mysql", pwd, host))
	if err != nil {
		panic(err)
	}
	defer db.Close()
	// make sure connection is available
	err = db.Ping()
	if err != nil {
		panic(err)
	}

	stmt, err := db.Prepare("INSERT userinfo SET username=?, description=?")
	if err != nil {
		panic(err)
	}
	_, err = stmt.Exec("KubeVela", "It's a test user")
	if err != nil {
		panic(err)
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprintf(w, "Version: %s\n", VERSION)
	})
	http.HandleFunc("/db", func(w http.ResponseWriter, r *http.Request) {
		rows, err := db.Query("select * from userinfo;")
		if err != nil {
			_, _ = fmt.Fprintf(w, "Error: %v\n", err)
		}
		for rows.Next() {
			var username string
			var description string
			err = rows.Scan(&username, &description)
			if err != nil {
				_, _ = fmt.Fprintf(w, "Scan Error: %v\n", err)
			}
			_, _ = fmt.Fprintf(w, "User: %s Description: %s\n", username, description)
		}
	})

	if err := http.ListenAndServe(":8088", nil); err != nil {
		println(err.Error())
	}
}
