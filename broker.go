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

func (broker *DataBroker) getDayByName(weekday string) (day *Day) {
	day = new(Day)
	broker.dbhandler.QueryModel(day, "name = ?", weekday)
	return
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

	for i := range avail {
		avail[i] = 10
	}

	for _, ticket := range tickets {
		avail[ticket.TierID-1]--
	}

	return
}

func (broker *DataBroker) GetDayIdByName(dayName string) int {
	day := broker.getDayByName(dayName)
	return int(day.DayID)
}

func (broker *DataBroker) GetShowsByPlaytime(weekday string, hour int) (shows []Show) {
	var time Time

	day := broker.getDayByName(weekday)
	broker.dbhandler.QueryModel(&time, "hour = ?", hour)
	broker.dbhandler.QueryModel(&shows, "day_id = ? AND time_id = ?", day.DayID, time.TimeID)

	return
}

func (broker *DataBroker) GetTicketDatesBySerials(serials []string) (dates []time.Time) {
	for _, serial := range serials {
		var t Ticket
		broker.dbhandler.QueryModel(&t, "serial = ?", serial)
		if t.Date.Year() != 1 { //if the year == 1 then the date object was not filled, no ticket found
			dates = append(dates, t.Date)
		}
	}

	return
}

func (broker *DataBroker) GetTicketCountByDay(date *time.Time) int {
	loc, _ := time.LoadLocation("America/Chicago")
	show1 := time.Date(date.Year(), date.Month(), date.Day(), 14, 0, 0, 0, loc)
	show2 := time.Date(date.Year(), date.Month(), date.Day(), 20, 0, 0, 0, loc)

	t1, err := broker.dbhandler.QueryModelAndCount(&[]Ticket{}, "date = ?", show1)
	if err != nil {
		panic(err)
	}

	t2, err := broker.dbhandler.QueryModelAndCount(&[]Ticket{}, "date = ?", show2)
	if err != nil {
		panic(err)
	}

	return t1 + t2
}
