package app

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"
	"database/sql"
	"datarequester/cfg"
)


var operationsNeeded = (cfg.AppStartTime - cfg.FirstPole) / 60 / 500 //needed requests count

type Site struct {
	Url string `json:"url"`
}

type Exchange [][]float64

type ExchangeStructred struct {
	MTS    uint64  `db:"mts"`  //millisecond time stamp
	Open   float64 `db:"open"` //First execution during the time frame
	Close  float64 //Last execution during the time frame
	High   float64 //Highest execution during the time frame
	Low    float64 //Lowest execution during the timeframe
	Volume float64 //Quantity of symbol traded within the timeframe
}

var DB = DBConnect()

func DBConnect() *sql.DB {
	var (
		//dbConnect string
		db  *sql.DB
		err error
	)

	for _, el := range cfg.Connections {
		db, err = sql.Open("mysql", el)
		//TODO: remove it, just check DB exists
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


func Main() {
	ready := make(chan bool)
	config, err := ioutil.ReadFile("config/config.json")
	check(err)
	fmt.Println(string(config))
	var buf Site
	err = json.Unmarshal([]byte(config), &buf)
	check(err)

	//SQL
	createTables()
	go server(ready)

	fmt.Println("PARSING")
	var rowsCount int64 = 0

	for i := 0; i < 0; i++ {
		//url := buf.Url + "5m:tBTCUSD/hist?start=1507106517000&end=1509307200000&limit=500&_=1509738036450"
		endTime := cfg.AppStartTime - i*30000 //500*60s
		startTime := endTime - 30000
		//check data in DB
		qb := `SELECT count(1) FROM trades WHERE trades.mts BETWEEN ` + strconv.Itoa(startTime) + "000 AND " + strconv.Itoa(endTime) + "000"
		rows, err := DB.Query(qb)
		check(err)
		var count int
		rows.Next()
		err = rows.Scan(&count)
		err = rows.Close()
		check(err)
		if count >= 500 {continue}

		url := buf.Url + fmt.Sprintf("1m:t%v/hist?start=%v000&end=%v000&limit=500", cfg.Alias, startTime, endTime)
		//
		fmt.Println(url)
		resp, err := http.Get(url)

		body, err := ioutil.ReadAll(resp.Body)
		resp.Body.Close()
		if resp.StatusCode != 200 {
			fmt.Println(resp.StatusCode)
			fmt.Println(string(body))
			continue
		}
		check(err)
		var exchangeBuf Exchange
		err = json.Unmarshal([]byte(body), &exchangeBuf)
		check(err)
		fmt.Printf("Body length: %v\n", len(exchangeBuf))

		for j := range exchangeBuf {
			insertTrade := fmt.Sprintf("INSERT IGNORE INTO trades VALUES(DEFAULT, %v, %v, %v, %v, %v, %v, DEFAULT)",
				exchangeBuf[j][0],
				exchangeBuf[j][1],
				exchangeBuf[j][2],
				exchangeBuf[j][3],
				exchangeBuf[j][4],
				exchangeBuf[j][5],
			)
			result, err := DB.Exec(insertTrade)
			check(err)
			isAffected, err := result.RowsAffected()
			check(err)
			rowsCount += isAffected
		}
		time.Sleep((1200 + time.Duration(generateDelay())) * time.Millisecond)
	}
	<-ready
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

var jsn = `{
  "url":"https://api.bitfinex.com/v2/candles/trade:"
}`

var server = func(ready chan bool){
	http.HandleFunc("/data/", shareData)
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/", fs)
	http.ListenAndServe(":8090", nil)
	ready <- true
}
