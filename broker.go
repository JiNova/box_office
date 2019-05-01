package main

import (
	"strconv"
	"time"
)

type DataBroker struct {
	dbhandler DBHandler
}

func (broker *DataBroker) Init() {
	broker.dbhandler.Init()
}

func (broker *DataBroker) Close() {
	broker.dbhandler.Close()
}

func (broker *DataBroker) GetAvailableTickets(date *time.Time, show *Show) (avail []int) {
	var all int
	var tickets []Ticket

	all, err := broker.dbhandler.CountEntries(&Tier{})
	if err != nil {
		panic(err)
	}

	query := Query{"date = ? AND show_id = ?", []string{date.String(), strconv.Itoa(int(show.ShowID))}}
	if err := broker.dbhandler.QueryModel(&tickets, &query); err != nil {
		panic(err)
	}
	avail = make([]int, all)

	for i:= range avail {
		avail[i] = 10
	}

	for _, ticket := range tickets {
		avail[ticket.TierID-1]--
	}

	return
}
