package main

import (
	"errors"
	"fmt"
	"math/rand"
	"reflect"
	"time"

	"github.com/jinzhu/gorm"

	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

type Data struct {
	db *gorm.DB
	ran *rand.Rand
}

func (handle *Data) Init() {
	if handle.db != nil {
		panic("already connected")
	}

	var err error
	handle.db, err = gorm.Open("sqlite3", "/tmp/test.db")

	if err != nil {
		panic("failed connection to db")
	}

	s1 := rand.NewSource(time.Now().UnixNano())
	handle.ran = rand.New(s1)
}

func (handle *Data) Close() {
	_ = handle.db.Close()
}

func (handle *Data) FillModelById(resources interface{}, id int) error {
	t := reflect.Indirect(reflect.ValueOf(resources))

	switch t.Kind() {
	case reflect.Struct:
		if err := handle.db.First(resources, id).Error; err != nil {
			return err
		}
	default:
		return errors.New("Resources must be a struct!")
	}

	return nil
}

func (handle *Data) FillModels(resources interface{}) error {
	t := reflect.Indirect(reflect.ValueOf(resources))

	switch t.Kind() {
	case reflect.Slice:
		if err := handle.db.Find(resources).Error; err != nil {
			return err
		}
	default:
		return errors.New("Resources must be a slice!")
	}

	return nil
}

func (handle *Data) LoadAssociations(resources interface{}, assocations ...string) error {
	return nil
}

func (handle *Data) GenPlan(plan *map[string]string) {
	var movies []Movie
	var movCount int

	handle.db.Find(&movies).Count(&movCount)

	// var showCnt int
	shows := make([][]Show, movCount)

	for i, movie := range movies {
		handle.db.Where("movie_id = ?", movie.MovieID).Find(&shows[i])

		for _, show := range shows[i] {
			var day Day
			var time Time
			handle.db.Model(&show).Related(&day)
			handle.db.Model(&show).Related(&time)
			(*plan)[fmt.Sprint(show.Screen)+day.Name+time.Desc] = movie.Title
		}
	}
}

func (handle *Data) GetAval(date *time.Time, show *Show) (avail []int) {
	var all int
	var tickets []Ticket

	handle.db.AutoMigrate(&Ticket{})

	handle.db.Model(&Tier{}).Count(&all)
	handle.db.Where("date = ? AND show_id = ?", date, show.ShowID).Find(&tickets)

	avail = make([]int, all)

	for i:= range avail {
		avail[i] = 10
	}

	for _, ticket := range tickets {
		avail[ticket.TierID-1]--
	}

	return
}

func (handle *Data) CreaTic(date *time.Time, show *Show, amount int, tier int) (serials []string){

	tickets := make([]Ticket, amount)

	for  i := 0; i < amount; i++ {
		var ticket Ticket
		ticket.Serial = handle.serial()
		ticket.TierID = uint(tier)
		ticket.Date = *date
		tickets[i] = ticket
		serials = append(serials, ticket.Serial)
	}

	handle.db.Model(show).Association("Tickets").Append(tickets)
	return
}

func (handle *Data) DayId(day string) int {

	var d Day
	handle.db.Where("name = ?", day).Find(&d)
	return int(d.DayID)
}

func (handle *Data) GetTic(serials []string) (dates []time.Time) {

	for _, serial := range serials {
		var t Ticket
		handle.db.Where("serial = ?", serial).Find(&t)
		dates = append(dates, t.Date)
	}

	return
}

func (handle *Data) GetTicD(date *time.Time) int {

	var t1, t2 int
	loc, _ := time.LoadLocation("America/Chicago")

	show1 := time.Date(date.Year(), date.Month(), date.Day(), 14, 0, 0, 0, loc)
	show2 := time.Date(date.Year(), date.Month(), date.Day(), 20, 0, 0, 0, loc)

	handle.db.Where("date  = ?", show1).Find(&[]Ticket{}).Count(&t1)
	handle.db.Where("date  = ?", show2).Find(&[]Ticket{}).Count(&t2)

	return t1 + t2
}

func (handle *Data) GetTiD2(date *time.Time, showID uint) (tickets int, vacant int) {

	handle.db.Model(&Ticket{}).Where("date = ? AND show_id = ?", date, showID).Count(&tickets)
	vacant = 40 - tickets

	return
}

func (handle *Data) TimDatS (weekday string, hour int) (shows []Show) {

	var d Day
	var t Time

	handle.db.Where("name = ?", weekday).Find(&d)
	handle.db.Where("hour = ?", hour).Find(&t)

	dayID, hourID := d.DayID, t.TimeID

	handle.db.Where("day_id = ? AND time_id = ?", dayID, hourID).Find(&shows)

	return
}

func (handle *Data) DelTic(serials []string) {

	for _, serial := range serials {
		handle.db.Where("serial = ?", serial).Delete(&Ticket{})
	}
}

func (handle *Data) Movies() (movies []Movie) {
	handle.db.Find(&movies)
	return
}

func (handle *Data) Shows(movie *Movie) (shows []Show) {

	handle.db.Where("movie_id = ?", movie.MovieID).Find(&shows)

	for i := range shows {
		var d Day
		var t Time

		handle.db.Model(&shows[i]).Related(&d)
		handle.db.Model(&shows[i]).Related(&t)

		shows[i].Day = d
		shows[i].Time = t
	}

	return
}

func (handle *Data) ShowMov(shows *[]Show) (movies []string) {

	movies = make([]string, len(*shows))

	for i, show := range *shows {
		var m Movie
		handle.db.Where("movie_id = ?", show.MovieID).Find(&m)
		movies[i] = m.Title
	}

	return
}

func (handle *Data) Plays(shows *[]Show) (day2h map[string]int) {

	var d Day
	var t Time

	day2h = make(map[string]int)

	for _, show := range *shows {
		handle.db.Model(show).Related(&d)
		handle.db.Model(show).Related(&t)
		day2h[d.Name] = t.Hour
	}

	return
}


func (handle *Data) Pricing(tier int) float64 {
	var t Tier
	handle.db.Where("tier_id = ?", tier).Find(&t)
	return t.Price
}