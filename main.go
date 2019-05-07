package main

import (
	"fmt"

	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

func main() {
	var broker DataBroker
	broker.Init()
	defer broker.Close()

	quit := false
	var input string

	for !quit {
		input = readcmd("main menu")

		switch input {
		case "sell":
			seller := SellHandler{&broker}
			seller.ProceedSelling()
		case "refund":
			refunder := RefundHandler{&broker}
			refunder.ProceedRefund()
		case "report":
			reporter := ReportHandler{&broker}
			reporter.ProceedReporting()
		case "quit":
			quit = true
		default:
			fmt.Println("I did not understand that, sorry :(")
		}
	}
}
