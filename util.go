package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

func readcmd(prompt string) (text string) {

	reader := bufio.NewReader(os.Stdin)
	fmt.Print(prompt + "> ")
	text, err := reader.ReadString('\n')
	if err != nil && err.Error() != "EOF" {
		panic(err)
	}

	text = strings.TrimSuffix(text, "\n")
	return
}

func choose(prompt string) (choice int, err error) {

	choice, err = strconv.Atoi(readcmd(prompt))

	if err != nil {
		fmt.Println("Something went wrong")
	}

	return
}

func genPlaytimeFromShow(broker *DataBroker, show *Show, today *time.Time) *time.Time {
	todayDayValue := broker.GetDayIdByName(today.Weekday().String())
	diff := (int(show.Day.DayID) - todayDayValue) % 7
	if diff < 0 {
		diff += 7
	}

	location, _ := time.LoadLocation("America/Chicago")
	playtime := time.Date(today.Year(), today.Month(), today.Day()+diff, show.Time.Hour, 0, 0, 0, location)
	return &playtime
}

func pricing(broker *DataBroker, tier int, day *Day, time *Time) (price float64) {

	var dis float64
	price = broker.GetTierPrice(tier)

	// 10 percent discount Mon-Thu
	if 1 <= day.DayID && day.DayID <= 4 {
		dis += price * 0.1
	}

	// 5 percent discount for matinee
	if time.TimeID == 1 {
		dis += price * 0.05
	}

	price -= dis
	return
}
