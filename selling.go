package main

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

type SellHandler struct {
	broker *DataBroker
}

func (seller *SellHandler) StartSelling() {
	var maxLen int

	movies := seller.broker.GetAllMovies()
	for i, mov := range movies {
		length := len("(" + strconv.Itoa(i+1) + ") " + mov.Title)
		if i == 0 || length > maxLen {
			maxLen = length
		}
	}

	for i, mov := range movies {
		if i%3 == 0 {
			fmt.Println()
		}

		item := "(" + strconv.Itoa(i+1) + ") " + mov.Title
		fmt.Print(item + strings.Repeat(" ", 1+maxLen-len(item)))
	}

	fmt.Println()
	choice, err := choose("Which movie")

	if err != nil {
		return
	}

	movie := movies[choice-1]
	shows := seller.broker.GetShowsByMovie(&movie)

	fmt.Println(movie.Title)

	// if there is a show today check if it has already been shown
	// if no, list twice
	// otherwise once
	//today := time.Now()
	//now := time.Now().Format("3 pm")

	//fmt.Println(today.String() + " " + now)
	dates := genDate(seller.broker, &shows)
	chDat := make([]time.Time, len(dates))
	i := 0

	for date, show := range dates {

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
	avail := seller.broker.GetAvailableTickets(&date, &show)

	fmt.Println(movie.Title + " on " + date.Format("Jan 2, Mon 3 pm"))

	for i, amount := range avail {
		price := pricing(seller.broker, i+1, &show.Day, &show.Time)
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
		serials := seller.broker.CreateTickets(&date, &show, amount, tierCho)
		output := fmt.Sprintf("The serials are %v. The customer will need them in case they want a refund!", serials)
		fmt.Println(output)
	}
}
