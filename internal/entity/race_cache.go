package entity

import (
	"sync"

	"github.com/google/uuid"
)

type (
	EventID       = uuid.UUID
	TimingPointID = uuid.UUID
	WaveID        = uuid.UUID
	LocationID    = uuid.UUID
)

type RaceCache struct {
	Events            map[EventID]*Event
	TimingPoints      map[TimingPointID]*TimingPoint
	EventTimingPoints map[EventID][]*TimingPoint
	Waves             map[WaveID]*Wave
	Locations         map[LocationID]*PhysicalLocation
	mu                sync.RWMutex
}

func NewRaceCache() *RaceCache {
	return &RaceCache{
		Events:            make(map[EventID]*Event),
		TimingPoints:      make(map[TimingPointID]*TimingPoint),
		EventTimingPoints: make(map[EventID][]*TimingPoint),
		Waves:             make(map[WaveID]*Wave),
		Locations:         make(map[LocationID]*PhysicalLocation),
	}
}

func (rc *RaceCache) StoreRaceConfig(cfg *RaceConfig) {
	rc.clearRaceConfig(cfg.ID)
	for _, l := range cfg.PhysicalLocations {
		rc.StoreLocation(l)
	}
	for _, e := range cfg.Events {
		rc.StoreEvent(e.Event)
		for _, tp := range e.TimingPoints {
			rc.StoreTimingPoint(tp)
		}
		for _, w := range e.Waves {
			rc.StoreWave(w)
		}

		rc.StoreEventTimingPoints(e.ID, e.TimingPoints)
	}
}

func (rc *RaceCache) clearRaceConfig(raceID uuid.UUID) {
	rc.removeEventsForRace(raceID)
	rc.removeTimingPointsForRace(raceID)
	rc.removeWavesForRace(raceID)
	rc.removeLocationsForRace(raceID)
}

func (rc *RaceCache) GetTimingPointsForEventLocation(eventID uuid.UUID, boxName string) []*TimingPoint {
	rc.mu.RLock()
	tps := []*TimingPoint{}
	for _, tp := range rc.EventTimingPoints[eventID] {
		if tp.BoxName == boxName {
			tps = append(tps, tp)
		}
	}
	rc.mu.RUnlock()
	return tps
}

// Events
func (rc *RaceCache) StoreEvent(e *Event) {
	rc.mu.Lock()
	rc.Events[e.ID] = e
	rc.mu.Unlock()
}

func (rc *RaceCache) GetEvent(id uuid.UUID) (*Event, bool) {
	rc.mu.RLock()
	e, found := rc.Events[id]
	rc.mu.RUnlock()
	return e, found
}

func (rc *RaceCache) RemoveEvent(id uuid.UUID) {
	rc.mu.Lock()
	delete(rc.Events, id)
	rc.mu.Unlock()
}

func (rc *RaceCache) removeEventsForRace(raceID uuid.UUID) {
	for _, e := range rc.Events {
		if e.RaceID == raceID {
			rc.RemoveEvent(e.ID)
			rc.removeEventTimingPointsForEvent(e.ID)
		}
	}
}

// Timing Points
func (rc *RaceCache) StoreTimingPoint(tp *TimingPoint) {
	rc.mu.Lock()
	rc.TimingPoints[tp.ID] = tp
	rc.mu.Unlock()
}

func (rc *RaceCache) GetTimingPoint(id uuid.UUID) (*TimingPoint, bool) {
	rc.mu.RLock()
	tp, found := rc.TimingPoints[id]
	rc.mu.RUnlock()
	return tp, found
}

func (rc *RaceCache) RemoveTimingPoint(id uuid.UUID) {
	rc.mu.Lock()
	delete(rc.TimingPoints, id)
	rc.mu.Unlock()
}

func (rc *RaceCache) removeTimingPointsForRace(raceID uuid.UUID) {
	for _, tp := range rc.TimingPoints {
		if tp.RaceID == raceID {
			rc.RemoveTimingPoint(tp.ID)
		}
	}
}

// Waves
func (rc *RaceCache) StoreWave(w *Wave) {
	rc.mu.Lock()
	rc.Waves[w.ID] = w
	rc.mu.Unlock()
}

func (rc *RaceCache) GetWave(id uuid.UUID) (*Wave, bool) {
	rc.mu.RLock()
	w, found := rc.Waves[id]
	rc.mu.RUnlock()
	return w, found
}

func (rc *RaceCache) RemoveWave(id uuid.UUID) {
	rc.mu.Lock()
	delete(rc.Waves, id)
	rc.mu.Unlock()
}

func (rc *RaceCache) removeWavesForRace(raceID uuid.UUID) {
	for _, w := range rc.Waves {
		if w.RaceID == raceID {
			rc.RemoveWave(w.ID)
		}
	}
}

// Physical Locations
func (rc *RaceCache) StoreLocation(l *PhysicalLocation) {
	rc.mu.Lock()
	rc.Locations[l.ID] = l
	rc.mu.Unlock()
}

func (rc *RaceCache) GetLocation(id uuid.UUID) (*PhysicalLocation, bool) {
	rc.mu.RLock()
	l, found := rc.Locations[id]
	rc.mu.RUnlock()
	return l, found
}

func (rc *RaceCache) RemoveLocation(id uuid.UUID) {
	rc.mu.Lock()
	delete(rc.Locations, id)
	rc.mu.Unlock()
}

func (rc *RaceCache) removeLocationsForRace(raceID uuid.UUID) {
	for _, l := range rc.Locations {
		if l.RaceID == raceID {
			rc.RemoveLocation(l.ID)
		}
	}
}

// Event timing points
func (rc *RaceCache) StoreEventTimingPoints(eventID uuid.UUID, tps []*TimingPoint) {
	rc.mu.Lock()
	rc.EventTimingPoints[eventID] = tps
	rc.mu.Unlock()
}

func (rc *RaceCache) removeEventTimingPointsForEvent(eventID uuid.UUID) {
	rc.mu.Lock()
	delete(rc.EventTimingPoints, eventID)
	rc.mu.Unlock()
}
