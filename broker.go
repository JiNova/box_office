package main

import (
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

	if err := broker.dbhandler.QueryModel(&tickets, "date = ? AND show_id = ?", date, show.ShowID); err != nil {
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

func (broker *DataBroker) GetDayIdByName(dayName string) int {
	var day Day
	broker.dbhandler.QueryModel(&day, "name = ?", dayName)
	return int(day.DayID)
}

func (broker *DataBroker) GetShowsByPlaytime(weekday string, hour int) (shows []Show) {
	return
}