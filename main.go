package main

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

func sell (data *DBHandler) {

	var maxLen int

	movies := data.Movies()
	for i,mov := range movies {
		length := len("("+strconv.Itoa(i+1)+") "+mov.Title)
		if i == 0 || length > maxLen {
			maxLen = length
		}
	}

	for i,mov := range movies {
		if i%3 == 0 {
			fmt.Println()
		}

		item := "("+strconv.Itoa(i+1)+") "+mov.Title
		fmt.Print(item+strings.Repeat(" ", 1+maxLen-len(item)))
	}

	fmt.Println()
	choice, err := choose("Which movie")

	if err != nil {
		return
	}

	movie := movies[choice-1]
	shows := data.Shows(&movie)

	fmt.Println(movie.Title)


	// if there is a show today check if it has already been shown
	// if no, list twice
	// otherwise once
	//today := time.Now()
	//now := time.Now().Format("3 pm")

	//fmt.Println(today.String() + " " + now)
	dates := genDate(data, &shows)
	chDat := make([]time.Time, len(dates))
	i := 0

	for date,show := range dates {

		chDat[i] = date
		prompt := fmt.Sprintf("(%v) %v %v on Screen %v", i+1, date.Format("Jan 2, Mon"), show.Time.Desc, show.Screen)
		fmt.Println(prompt)
		i++
	}

	choice, err = choose("Which show")

	if err != nil {
		return
	}

	date := chDat[choice-1]
	show := dates[date]
	avail := data.GetAval(&date, &show)

	fmt.Println(movie.Title + " on " + date.Format("Jan 2, Mon 3 pm"))

	for i,amount := range avail{
		price := pricing(data, i+1, &show.Day, &show.Time)
		prompt := fmt.Sprintf("Tier %v: %v left, $%.2f each", i+1, amount, price)
		fmt.Println(prompt)
	}

	tierCho, err := choose("Which tier")

	if err != nil {
		return
	}

	amount, err := choose("How many")

	if err != nil {
		return
	}

	if avail[tierCho-1]-amount < 0 {
		fmt.Println("Not enough seats left!")
		//FIXME: loop this
		return
	} else {
		serials := data.CreaTic(&date, &show, amount, tierCho)
		output := fmt.Sprintf("The serials are %v. The customer will need them in case they want a refund!", serials)
		fmt.Println(output)
	}
}

func refund(data *DBHandler) {

	fmt.Println("Please provide the serial number(s), separating them with a space if there you are " +
					"trying to refund more than one ticket at a time")
	serials := strings.Split(readcmd("serials"), " ")

	dates := data.GetTic(serials)
	now := time.Now()

	for i,date := range dates {
		if date.Before(now) {
			fmt.Println("Not refunding " + serials[i] + ", show already took place!")
			serials[i] = serials[len(serials)-1]
			serials = serials[:len(serials)-1]
		}
	}

	data.DelTic(serials)

	fmt.Println("All your eligible tickets have been refunded!")
}

func report(data *DBHandler) {

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

		shows := data.TimDatS(date.Weekday().String(),  date.Hour())
		movTit := data.ShowMov(&shows)

		for i, show := range shows {
			output := fmt.Sprintf("(%v) %v on Screen %v", i+1, movTit[i], show.Screen)
			fmt.Println(output)
		}

		choice, err = choose("which show")

		if err != nil {
			return
		}

		tickets, vacant := data.GetTiD2(&date, shows[choice-1].ShowID)

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

		output := fmt.Sprintf("On %v we sold %v tickets", date.Format("Jan 2, Mon") , data.GetTicD(&date))
		fmt.Println(output)

	} else {
		fmt.Println("I did not understand that, sorry :(")
	}
}

func main() {

	var data DBHandler
	data.Init()

	defer data.Close()

	quit := false
	var input string

	for !quit {
		input = readcmd("main menu")

		switch input {
		case "sell":
			sell(&data)
		case "refund":
			refund(&data)
		case "report":
			report(&data)
		case "quit":
			quit = true
		default:
			fmt.Println("I did not understand that, sorry :(")
		}
	}
}
