package cfg

const (
	//delta = 12255925sec
	AppStartTime = 1511036725 //[Seconds] == 2017-11-18T20:25:25+00:00
	FirstPole    = 1498780800 // [Seconds] == 2017-06-30T00:00:00+00:00

	Alias = "BTCUSD"
)

var Connections = []string{
	//"test",
	"root:111@tcp(127.0.0.1:3306)/",
	"root:12345678@tcp(127.0.0.1:3306)/",
	"root@tcp(127.0.0.1:3306)/",
	"root:111@localhost/",
}
