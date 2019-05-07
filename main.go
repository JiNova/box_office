package main

import (
	"fmt"
	"time"

	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

func report(broker *DataBroker) {

	fmt.Println("Create a report for \n (1) A specific showtime \n (2) All shows on a specific date")
	choice, err := choose("report")

	if err != nil {
		return
	}

	if choice == 1 {

		fmt.Println("Which date? (mm/dd/yyyy)")
		input := readcmd("date")

		date, err := time.Parse("01/02/2006", input)

		if err != nil {
			fmt.Println("I did not understand that, sorry :(")
			return
		}

		fmt.Println("Which time? (h am/pm)")
		input = readcmd("time")

		hour, err := time.Parse("3 pm", input)

		if err != nil {
			fmt.Println("I did not understand that, sorry :(")
			return
		}

		loc, _ := time.LoadLocation("America/Chicago")
		date = time.Date(date.Year(), date.Month(), date.Day(), hour.Hour(), 0, 0, 0, loc)

		shows := broker.GetShowsByPlaytime(date.Weekday().String(), date.Hour())
		movTit := broker.GetMovieTitlesByShows(shows)

		for i, show := range shows {
			output := fmt.Sprintf("(%v) %v on Screen %v", i+1, movTit[i], show.Screen)
			fmt.Println(output)
		}

		choice, err = choose("which show")

		if err != nil {
			return
		}

		tickets, vacant := broker.GetSoldVacantTicketsByShow(&date, shows[choice-1].ShowID)

		output := fmt.Sprintf("%v on %v sold %v tickets, %v seats empty", movTit[choice-1],
			date.Format("Jan 2, 2006"), tickets, vacant)
		fmt.Println(output)

	} else if choice == 2 {

		fmt.Println("Which date? (mm/dd/yyyy)")
		input := readcmd("date")

		date, err := time.Parse("01/02/2006", input)

		if err != nil {
			fmt.Println("I did not understand that, sorry :(")
			return
		}

		output := fmt.Sprintf("On %v we sold %v tickets", date.Format("Jan 2, Mon"), broker.GetTicketCountByDay(&date))
		fmt.Println(output)

	} else {
		fmt.Println("I did not understand that, sorry :(")
	}
}

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
			refund(&broker)
		case "report":
			report(&broker)
		case "quit":
			quit = true
		default:
			fmt.Println("I did not understand that, sorry :(")
		}
	}
}
