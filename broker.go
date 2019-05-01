package main

import "time"

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

	broker.dbhandler.db.Where("date = ? AND show_id = ?", date, show.ShowID).Find(&tickets)
	//broker.dbhandler.QueryModel(&tickets, Query{"date = ? AND show_id = ?", })
	avail = make([]int, all)

	for i:= range avail {
		avail[i] = 10
	}

	for _, ticket := range tickets {
		avail[ticket.TierID-1]--
	}

	return
}
