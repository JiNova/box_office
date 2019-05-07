package main

import (
	"fmt"
	"strings"
	"time"
)

func refund(broker *DataBroker) {

	fmt.Println("Please provide the serial number(s), separating them with a space if there you are " +
		"trying to refund more than one ticket at a time")
	serials := strings.Split(readcmd("serials"), " ")

	dates := broker.GetTicketDatesBySerials(serials)
	now := time.Now()

	for i, date := range dates {
		if date.Before(now) {
			fmt.Println("Not refunding " + serials[i] + ", show already took place!")
			serials[i] = serials[len(serials)-1]
			serials = serials[:len(serials)-1]
		}
	}

	broker.DeleteTicketsBySerial(serials)
	fmt.Println("All your eligible tickets have been refunded!")
}
