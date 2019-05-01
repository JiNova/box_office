package main

import (
	"github.com/jinzhu/gorm"
	"time"
)

type Movie struct {
	MovieID uint `gorm:"primary_key"`
	Title   string
	Length  uint
	Shows   []Show `gorm:"foreignkey:MovieID"`
}

type Show struct {
	ShowID  uint `gorm:"primary_key"`
	Screen  uint
	Tickets []Ticket `gorm:"foreignkey:ShowID"`
	Day     Day
	DayID   uint
	Time    Time
	TimeID  uint
	MovieID uint
}

type Ticket struct {
	gorm.Model
	Serial string `gorm:"unique;not null"`
	TierID uint
	Tier   Tier
	Date time.Time
	ShowID uint
}

type Tier struct {
	TierID uint `gorm:"primary_key"`
	Price  float64
}

type Day struct {
	DayID uint `gorm:"primary_key"`
	Name  string
}

type Time struct {
	TimeID uint `gorm:"primary_key"`
	Desc   string
	Hour   int
}
