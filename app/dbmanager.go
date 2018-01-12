package app

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
)

const (
	dbNname            = "trading"
	candelsTableStruct = `(
		id INT(10) PRIMARY KEY NOT NULL AUTO_INCREMENT,
		mts BIGINT(20),
		open FLOAT,
		close FLOAT,
		hight FLOAT,
		low FLOAT,
		volume FLOAT,
		ts TIMESTAMP DEFAULT 'CURRENT_TIMESTAMP' NOT NULL
		);`
)



func getMySQL() string {
	b, err := ioutil.ReadFile("dev/mtest_users.sql") // just pass the file name
	if err != nil {
		fmt.Println("ERROR Getting SQL file: " + err.Error())
	}
	return string(b) // convert content to a 'string'
}

func createTables() {
	_, err := DB.Exec(`CREATE TABLE IF NOT EXISTS trades` + candelsTableStruct)
	check(err)
}

func shareData(w http.ResponseWriter, r *http.Request) {
	qb := fmt.Sprintf("SELECT mts, open, close, hight, low, volume FROM %v LIMIT 100", r.URL.Path[len("/data/"):])
	rows, err := DB.Query(qb)
	//rows, err := DB.Query("SELECT mts, open, close, high, low, volume FROM trades")
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
