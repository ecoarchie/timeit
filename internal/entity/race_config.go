package entity

import (
	"fmt"
	"time"
)

type RaceConfig struct {
	*Race
	TimeReaders []*TimeReader  `json:"time_readers"`
	Events      []*EventConfig `json:"events"`
}

type EventConfig struct {
	*Event
	Splits     []*Split    `json:"splits"`
	Waves      []*Wave     `json:"waves"`
	Categories []*Category `json:"categories"`
}

type RaceFormData struct {
	Id       string    `json:"race_id"`
	Name     string    `json:"name"`
	RaceDate time.Time `json:"race_date"`
	Timezone string    `json:"timezone"`
}

func (rc RaceConfig) String() string {
	return fmt.Sprintf(`{
	RaceID: %s,
	Name: %s,
	Timezone: %s,
	Events: 
		%+v
}`, rc.Name, rc.Name, rc.Timezone, rc.Events)
}
