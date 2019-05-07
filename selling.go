package main

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

type SellHandler struct {
	broker *DataBroker
}

func (seller *SellHandler) StartSelling() {

	movie, err := seller.PresentMovies()
	shows := seller.broker.GetShowsByMovie(movie)

	fmt.Println(movie.Title)

	// if there is a show today check if it has already been shown
	// if no, list twice
	// otherwise once
	show, date, err := seller.ChooseShow(movie, shows)
	avail := seller.broker.GetAvailableTickets(date, show)

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
		serials := seller.broker.CreateTickets(date, show, amount, tierCho)
		output := fmt.Sprintf("The serials are %v. The customer will need them in case they want a refund!", serials)
		fmt.Println(output)
	}
}

func (seller *SellHandler) PresentMovies() (movie *Movie, err error) {
	movies := seller.broker.GetAllMovies()

	var maxItemLength int
	for i, mov := range movies {
		length := len("(" + strconv.Itoa(i+1) + ") " + mov.Title)
		if i == 0 || length > maxItemLength {
			maxItemLength = length
		}
	}

	for i, mov := range movies {
		if i%3 == 0 {
			fmt.Println()
		}

		item := "(" + strconv.Itoa(i+1) + ") " + mov.Title
		fmt.Print(item + strings.Repeat(" ", 1+maxItemLength-len(item)))
	}

	fmt.Println()
	choice, err := choose("Which movie")

	if err != nil {
		return
	}

	movie = &(movies[choice-1])
	return
}

func (seller *SellHandler) ChooseShow(movie *Movie, shows []Show) (show *Show, playtime *time.Time, err error) {

	for i, show := range shows {
		prompt := fmt.Sprintf("(%v) %v, at %v", i+1, show.Day.Name, show.Time.Desc)
		fmt.Println(prompt)
	}

	choice, err := choose("Which show")
	if err != nil {
		return
	} else if choice < 0 || choice > len(shows) {
		return nil, nil, errors.New("Invalid choice!")
	}

	show = &(shows[choice-1])
	now := time.Now()
	playtime = genPlaytimeFromShow(seller.broker, show, &now)

	if show.Day.Name == now.Weekday().String() && now.Hour() < show.Time.Hour {
		fmt.Println("(1) Today or (2) in 7 days?")
		if choice, err = choose("When"); err != nil {
			return nil, nil, errors.New("Invalid choice!")
		} else if choice < 1 || choice > 2 {
			return
		} else if choice == 2 {
			playtime.AddDate(0, 0, 7)
		}
	}

	return
}
