package main

import (
	"fmt"
	"strings"
	"time"
)

type RefundHandler struct {
	broker *DataBroker
}

func (refunder *RefundHandler) refund() {
	fmt.Println("Please provide the serial number(s), separating them with a space if you are " +
		"trying to refund more than one ticket at a time")

	var serials []string
	if userInput := readcmd("serials"); userInput == "" || userInput == " " {
		return
	} else {
		serials = strings.Split(userInput, " ")
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
