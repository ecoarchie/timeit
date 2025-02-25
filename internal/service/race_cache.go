package service

import (
	"sync"

	"github.com/ecoarchie/timeit/internal/entity"
	"github.com/google/uuid"
)

type (
	RaceID       = uuid.UUID
	EventID      = uuid.UUID
	SplitID      = uuid.UUID
	WaveID       = uuid.UUID
	TimeReaderID = uuid.UUID
)

type RaceCache struct {
	Races       map[RaceID]*entity.Race
	Events      map[EventID]*entity.Event
	Splits      map[SplitID]*entity.Split
	EventSplits map[EventID][]*entity.Split
	Waves       map[WaveID]*entity.Wave
	TimeReaders map[TimeReaderID]*entity.TimeReader
	mu          sync.RWMutex
}

func NewRaceCache() *RaceCache {
	return &RaceCache{
		Races:       make(map[RaceID]*entity.Race),
		Events:      make(map[EventID]*entity.Event),
		Splits:      make(map[SplitID]*entity.Split),
		EventSplits: make(map[EventID][]*entity.Split),
		Waves:       make(map[WaveID]*entity.Wave),
		TimeReaders: make(map[TimeReaderID]*entity.TimeReader),
	}
}

func (rc *RaceCache) StoreRaceConfig(cfg *entity.RaceConfig) {
	rc.clearRaceCache(cfg.ID)
	rc.Races[cfg.ID] = &entity.Race{}
	rc.Races[cfg.ID].Name = cfg.Name
	rc.Races[cfg.ID].Timezone = cfg.Timezone

	for _, l := range cfg.TimeReaders {
		rc.StoreTimeReader(l)
	}
	for _, e := range cfg.Events {
		rc.StoreEvent(e.Event)
		for _, tp := range e.Splits {
			rc.StoreSplit(tp)
		}
		for _, w := range e.Waves {
			rc.StoreWave(w)
		}

		rc.StoreEventSplits(e.ID, e.Splits)
	}
}

func (rc *RaceCache) clearRaceCache(raceID uuid.UUID) {
	rc.removeEventsForRace(raceID)
	rc.removeSplitsForRace(raceID)
	rc.removeWavesForRace(raceID)
	rc.removeTimeReadersForRace(raceID)
}

func (rc *RaceCache) GetSplitsForEventTimeReader(eventID uuid.UUID, TimeReaderID uuid.UUID) []*entity.Split {
	rc.mu.RLock()
	tps := []*entity.Split{}
	for _, tp := range rc.EventSplits[eventID] {
		if tp.TimeReaderID == TimeReaderID {
			tps = append(tps, tp)
		}
	}
	rc.mu.RUnlock()
	return tps
}

// Events
func (rc *RaceCache) StoreEvent(e *entity.Event) {
	rc.mu.Lock()
	rc.Events[e.ID] = e
	rc.mu.Unlock()
}

func (rc *RaceCache) GetEvent(id uuid.UUID) (*entity.Event, bool) {
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
			rc.removeEventSplitsForEvent(e.ID)
		}
	}
}

// Timing Points
func (rc *RaceCache) StoreSplit(tp *entity.Split) {
	rc.mu.Lock()
	rc.Splits[tp.ID] = tp
	rc.mu.Unlock()
}

func (rc *RaceCache) GetSplit(id uuid.UUID) (*entity.Split, bool) {
	rc.mu.RLock()
	tp, found := rc.Splits[id]
	rc.mu.RUnlock()
	return tp, found
}

func (rc *RaceCache) RemoveSplit(id uuid.UUID) {
	rc.mu.Lock()
	delete(rc.Splits, id)
	rc.mu.Unlock()
}

func (rc *RaceCache) removeSplitsForRace(raceID uuid.UUID) {
	for _, tp := range rc.Splits {
		if tp.RaceID == raceID {
			rc.RemoveSplit(tp.ID)
		}
	}
}

// Waves
func (rc *RaceCache) StoreWave(w *entity.Wave) {
	rc.mu.Lock()
	rc.Waves[w.ID] = w
	rc.mu.Unlock()
}

func (rc *RaceCache) GetWave(id uuid.UUID) (*entity.Wave, bool) {
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

// TimeReaders
func (rc *RaceCache) StoreTimeReader(l *entity.TimeReader) {
	rc.mu.Lock()
	rc.TimeReaders[l.ID] = l
	rc.mu.Unlock()
}

func (rc *RaceCache) GetTimeReader(id uuid.UUID) (*entity.TimeReader, bool) {
	rc.mu.RLock()
	l, found := rc.TimeReaders[id]
	rc.mu.RUnlock()
	return l, found
}

func (rc *RaceCache) RemoveTimeReader(id uuid.UUID) {
	rc.mu.Lock()
	delete(rc.TimeReaders, id)
	rc.mu.Unlock()
}

func (rc *RaceCache) removeTimeReadersForRace(raceID uuid.UUID) {
	for _, l := range rc.TimeReaders {
		if l.RaceID == raceID {
			rc.RemoveTimeReader(l.ID)
		}
	}
}

// Event splits
func (rc *RaceCache) StoreEventSplits(eventID uuid.UUID, tps []*entity.Split) {
	rc.mu.Lock()
	rc.EventSplits[eventID] = tps
	rc.mu.Unlock()
}

func (rc *RaceCache) removeEventSplitsForEvent(eventID uuid.UUID) {
	rc.mu.Lock()
	delete(rc.EventSplits, eventID)
	rc.mu.Unlock()
}
