package main

import (
	"crypto/sha256"
	"encoding/hex"
	"math/rand"
	"strconv"
	"time"
)

type DataBroker struct {
	dbhandler DBHandler
	ran       *rand.Rand
}

func (broker *DataBroker) Init() {
	broker.dbhandler.Init()
	s1 := rand.NewSource(time.Now().UnixNano())
	broker.ran = rand.New(s1)
}

func (broker *DataBroker) Close() {
	broker.dbhandler.Close()
}

func (broker *DataBroker) generateTicketSerial() (code string) {

	hasher := sha256.New()
	hasher.Write([]byte(strconv.Itoa(broker.ran.Intn(500)) + time.Now().String()))
	code = hex.EncodeToString(hasher.Sum(nil))[:8]
	return
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
	location, _ := time.LoadLocation("America/Chicago")
	matinee := time.Date(date.Year(), date.Month(), date.Day(), 14, 0, 0, 0, location)
	nightShow := time.Date(date.Year(), date.Month(), date.Day(), 20, 0, 0, 0, location)

	matineeTickets, err := broker.dbhandler.QueryModelAndCount(&[]Ticket{}, "date = ?", matinee)
	if err != nil {
		panic(err)
	}

	nightTickets, err := broker.dbhandler.QueryModelAndCount(&[]Ticket{}, "date = ?", nightShow)
	if err != nil {
		panic(err)
	}

	return matineeTickets + nightTickets
}

func (broker *DataBroker) GetSoldVacantTicketsByShow(date *time.Time, showID uint) (sold int, vacant int) {
	sold, err := broker.dbhandler.QueryModelAndCount(&[]Ticket{}, "date = ? AND show_id = ?", date, showID)
	if err != nil {
		panic(err)
	}
	vacant = 40 - sold
	return
}

func (broker *DataBroker) GetAllMovies() (movies []Movie) {
	if err := broker.dbhandler.FillModels(&movies); err != nil {
		panic(err)
	}
	return
}

func (broker *DataBroker) GetShowsByMovie(movie *Movie) (shows []Show) {
	movie.Shows = []Show{}
	if err := broker.dbhandler.LoadAssociations(movie, "Shows"); err != nil {
		panic(err)
	}

	for i := range movie.Shows {
		show := &(movie.Shows[i])
		show.Day = Day{}
		show.Time = Time{}
		if err := broker.dbhandler.LoadRelated(show, &(show.Day)); err != nil {
			panic(err)
		}

		if err := broker.dbhandler.LoadRelated(show, &(show.Time)); err != nil {
			panic(err)
		}
	}

	return movie.Shows
}

func (broker *DataBroker) CreateTickets(date *time.Time, show *Show, amount int, tier int) (serials []string) {
	tickets := make([]Ticket, amount)

	for i := range tickets {
		ticket := &(tickets[i])
		ticket.Serial = broker.generateTicketSerial()
		ticket.TierID = uint(tier)
		ticket.Date = *date
		serials = append(serials, ticket.Serial)
	}

	broker.dbhandler.CreateAssociations(show, "Tickets", tickets)
	return
}

func (broker *DataBroker) DeleteTicketsBySerial(serials []string) {
	for _, serial := range serials {
		if err := broker.dbhandler.QueryModelAndDeleteData(&Ticket{}, "serial = ?", serial); err != nil {
			panic(err)
		}
	}
}

func (broker *DataBroker) GetMovieTitlesByShows(shows []Show) (movieTitles []string) {
	movieTitles = make([]string, len(shows))

	for i, show := range shows {
		var m Movie
		broker.dbhandler.QueryModel(&m, "movie_id = ?", show.MovieID)
		movieTitles[i] = m.Title
	}

	return
}

func (broker *DataBroker) GetMovieById(movieID int) (movie *Movie) {
	movie = new(Movie)
	broker.dbhandler.FillModelById(movie, movieID)
	return
}

func (broker *DataBroker) GetShowById(showID int) (show *Show) {
	show = new(Show)
	broker.dbhandler.FillModelById(show, showID)
	return
}

func (broker *DataBroker) IsValidSerial(serial string) bool {
	if err := broker.dbhandler.QueryModel(&Ticket{}, "serial = ?", serial); err != nil {
		if err.Error() == "record not found" {
			return false
		} else {
			panic(err)
		}
	} else {
		return true
	}
}

func (broker *DataBroker) GetTierPrice(tierNum int) float64 {
	var tier Tier
	if err := broker.dbhandler.QueryModel(&tier, "tier_id = ?", tierNum); err != nil {
		panic(err)
	}
	return tier.Price
}
