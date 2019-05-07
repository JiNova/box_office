package main

import (
	"fmt"
	"time"
)

type ReportHandler struct {
	broker *DataBroker
}

func (reporter *ReportHandler) ProceedReporting() {

	fmt.Println("Create a report for \n (1) A specific showtime \n (2) All shows on a specific date")
	choice, err := choose("report")
	if err != nil {
		return
	}

	if choice == 1 {
		day, err := reporter.GetDateFromUser()
		if err != nil {
			return
		}

		hour := reporter.GetTimeFromUser()
		loc, _ := time.LoadLocation("America/Chicago")
		reportDate := time.Date(day.Year(), day.Month(), day.Day(), hour.Hour(), 0, 0, 0, loc)

		show, movieTitle, err := reporter.GetSpecificShowFromUser(&reportDate)
		if err != nil {
			fmt.Println("Something went wrong, please try again!", err)
			return
		}

		tickets, vacant := reporter.broker.GetSoldVacantTicketsByShow(&reportDate, show.ShowID)
		output := fmt.Sprintf("%v on %v sold %v tickets, %v seats empty", movieTitle,
			reportDate.Format("Jan 2, 2006"), tickets, vacant)
		fmt.Println(output)

	} else if choice == 2 {

		date, err := reporter.GetDateFromUser()
		if err != nil {
			return
		}

		output := fmt.Sprintf("On %v we sold %v tickets", date.Format("Jan 2, Mon"), reporter.broker.GetTicketCountByDay(date))
		fmt.Println(output)

	} else {
		fmt.Println("I did not understand that, sorry :(")
	}
}

func (reporter *ReportHandler) GetDateFromUser() (*time.Time, error) {
	fmt.Println("Which date? (mm/dd/yyyy)")
	input := readcmd("date")

	date, err := time.Parse("01/02/2006", input)

	if err != nil {
		fmt.Println("I did not understand that, sorry :(")
		return nil, err
	}

	return &date, nil
}

func (reporter *ReportHandler) GetTimeFromUser() *time.Time {
	fmt.Println("Which time? (h am/pm)")
	input := readcmd("time")

	hour, err := time.Parse("3 pm", input)

	if err != nil {
		fmt.Println("I did not understand that, sorry :(")
		return nil
	}

	return &hour
}

func (reporter *ReportHandler) GetSpecificShowFromUser(reportDate *time.Time) (show *Show, movieTitle string, err error) {
	shows := reporter.broker.GetShowsByPlaytime(reportDate.Weekday().String(), reportDate.Hour())
	movieTitles := reporter.broker.GetMovieTitlesByShows(shows)

	for i, show := range shows {
		output := fmt.Sprintf("(%v) %v on Screen %v", i+1, movieTitles[i], show.Screen)
		fmt.Println(output)
	}

	choice, err := choose("which show")

	if err != nil {
		return
	}

	show = &(shows[choice-1])
	movieTitle = movieTitles[choice-1]
	return
}
