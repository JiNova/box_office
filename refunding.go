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

	dates := refunder.broker.GetTicketDatesBySerials(serials)
	now := time.Now()

	for i, date := range dates {
		if date.Before(now) {
			fmt.Println("Not refunding " + serials[i] + ", show already took place!")
			serials[i] = serials[len(serials)-1]
			serials = serials[:len(serials)-1]
		}
	}

	refunder.broker.DeleteTicketsBySerial(serials)
	fmt.Println("All your eligible tickets have been refunded!")
}

func (refunder *RefundHandler) GetSerialsFromUser() (serials []string, err error) {
	fmt.Println("Please provide the serial number(s), separating them with a space if you are " +
		"trying to refund more than one ticket at a time")

	if userInput := readcmd("serials"); userInput == "" || userInput == " " {
		return nil, errors.New("Invalid input!")
	} else {
		serials = strings.Split(userInput, " ")
	}

	return
}
