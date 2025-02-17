package entity

import (
	"time"
)

type RaceConfig struct {
	*Race
	PhysicalLocations []*PhysicalLocation `json:"locations"`
	Events            []*EventConfig      `json:"events"`
}

type EventConfig struct {
	*Event
	TimingPoints []*TimingPoint `json:"timing_points"`
	Waves        []*Wave        `json:"waves"`
	Categories   []*Category    `json:"categories"`
}

type RaceFormData struct {
	Id       string    `json:"race_id"`
	Name     string    `json:"name"`
	RaceDate time.Time `json:"race_date"`
	Timezone string    `json:"timezone"`
}
