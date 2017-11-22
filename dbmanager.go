package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
)

const dbNname = "trading"

var connections = []string{
	"test",
	"root:111@tcp(127.0.0.1:3306)/",
	"root:12345678@tcp(127.0.0.1:3306)/",
	"root@tcp(127.0.0.1:3306)/",
	"root:111@localhost/",
}

func getMySQL() string {
	b, err := ioutil.ReadFile("dev/mtest_users.sql") // just pass the file name
	if err != nil {
		fmt.Println("ERROR Getting SQL file: " + err.Error())
	}
	return string(b) // convert content to a 'string'
}

func DBConnect() *sql.DB {
	var (
		//dbConnect string
		db  *sql.DB
		err error
	)

	for _, el := range connections {
		db, err = sql.Open("mysql", el)
		result, _ := db.Exec("CREATE DATABASE IF NOT EXISTS " + dbNname)
		db, err = sql.Open("mysql", el+dbNname)
		if result != nil {
			fmt.Println("CONNECTED TO DB-SERVER: " + el)
			break
		}
	}
	if err != nil {
		fmt.Println("DB Initialization finished with errors:")
		panic(err)
	} else {
		fmt.Printf("DB Initialization finished successfully [%v]\n", dbNname)
	}
	// defer db.Close()
	return db
}

func createTables () {
	_, err := DB.Exec(`CREATE TABLE IF NOT EXISTS trades	(
		    id INT(10) PRIMARY KEY NOT NULL AUTO_INCREMENT,
		    mts BIGINT(20),
		    open FLOAT,
		    close FLOAT,
		    hight FLOAT,
		    low FLOAT,
		    volume FLOAT,
		    ts TIMESTAMP DEFAULT 'CURRENT_TIMESTAMP' NOT NULL
		);`)
	check(err)
}

func shareData(w http.ResponseWriter, r *http.Request) {
	qb := fmt.Sprintf("SELECT mts, open, close, hight, low, volume FROM %v LIMIT 100", r.URL.Path[len("/data/"):])
	rows, err := DB.Query(qb)
	//rows, err := DB.Query("SELECT mts, open, close, hight, low, volume FROM trades")
	check(err)
	var result [][]float64

	var buf struct {
		mts    float64
		open   float64
		close  float64
		high   float64
		low    float64
		volume float64
	}
	for rows.Next() {
		err = rows.Scan(
			&buf.mts,
			&buf.open,
			&buf.close,
			&buf.high,
			&buf.low,
			&buf.volume,
		)
		check(err)
		result = append(result, []float64{buf.mts, buf.open, buf.close, buf.high, buf.low, buf.volume})
	}
	out, err := json.Marshal(result)
	check(err)
	err = rows.Close()
	check(err)
	fmt.Fprintf(w, "%v", string(out))
}

//func page (w http.ResponseWriter, r *http.Request)
