package main

import (
	"bufio"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

func (handle *DBHandler) serial() (code string) {

	hasher := sha256.New()
	hasher.Write([]byte(strconv.Itoa(handle.ran.Intn(500))+time.Now().String()))
	code = hex.EncodeToString(hasher.Sum(nil))[:8]
	return
}

func readcmd(prompt string) (text string) {

	reader := bufio.NewReader(os.Stdin)
	fmt.Print(prompt+"> ")
	text, _ = reader.ReadString('\n')
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

func genDate(data *DBHandler, shows *[]Show) (showt map[time.Time]Show) {

	today := time.Now()
	todayId := data.DayId(today.Weekday().String())

	showt = make(map[time.Time]Show)

	for _, show := range *shows {
		diff := (int(show.Day.DayID) - todayId) % 7

		if diff == 0 {

			r := today.AddDate(0, 0, 7)

			if show.Time.Hour < today.Hour() {
				showt[time.Date(today.Year(), today.Month(), today.Day(), show.Time.Hour, 0, 0, 0, today.Location())] = show
			}

			showt[time.Date(today.Year(), today.Month(), r.Day(), show.Time.Hour, 0, 0, 0, today.Location())] = show
		} else {
			if diff < 0 {
				diff += 7
			}

			r := today.AddDate(0, 0, diff)
			showt[time.Date(r.Year(), r.Month(), r.Day(), show.Time.Hour, 0, 0, 0, today.Location())] = show
		}
	}

	return
}

func pricing(data *DBHandler, tier int, day *Day, time *Time) (price float64) {

	var dis float64
	price = data.Pricing(tier)

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
