package main

import (
	"errors"
	"fmt"
	"strings"
	"time"
)

type RefundHandler struct {
	broker *DataBroker
}

func (refunder *RefundHandler) ProceedRefund() {

	serials, err := refunder.GetSerialsFromUser()
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	serials, err = refunder.FilterSerialsForValid(serials)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	refunder.broker.DeleteTicketsBySerial(serials)
	fmt.Println("All your eligible tickets have been refunded!")
}

func (refunder *RefundHandler) GetSerialsFromUser() (serials []string, err error) {
	fmt.Println("Please provide the serial number(s), separating them with a space if you are " +
		"trying to refund more than one ticket at a time")

	var userInput string
	if userInput = readcmd("serials"); userInput == "" || userInput == " " {
		return nil, errors.New("invalid input")
	}

	serials = strings.Split(userInput, " ")
	return
}

func (refunder *RefundHandler) FilterSerialsForValid(serials []string) ([]string, error) {
	var filteredSerials []string

	for i, _ := range serials {
		if refunder.broker.IsValidSerial(serials[i]) {
			filteredSerials = append(filteredSerials, serials[i])
		}
	}

	dates := refunder.broker.GetTicketDatesBySerials(filteredSerials)
	now := time.Now()

	for i, date := range dates {
		if date.Before(now) {
			fmt.Println("Not refunding " + serials[i] + ", show already took place!")
			filteredSerials[i] = serials[len(serials)-1]
			filteredSerials = serials[:len(serials)-1]
		}
	}

	return filteredSerials, nil
}
