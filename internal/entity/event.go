package entity

import (
	"cmp"
	"fmt"
	"slices"
	"time"

	"github.com/ecoarchie/timeit/internal/controller/httpv1/dto"
	"github.com/ecoarchie/timeit/pkg/validator"
	"github.com/google/uuid"
)

type Event struct {
	ID               uuid.UUID   `json:"event_id"`
	RaceID           uuid.UUID   `json:"race_id"`
	Name             string      `json:"event_name"`
	DistanceInMeters int         `json:"distance_in_meters"`
	EventDate        time.Time   `json:"event_date"`
	Splits           []*Split    `json:"splits"`
	Waves            []*Wave     `json:"waves"`
	Categories       []*Category `json:"categories"`
}

func NewEvent(e *dto.EventDTO, ss []*dto.SplitDTO, trs []*dto.TimeReaderDTO, ww []*dto.WaveDTO, cc []*dto.CategoryDTO, v *validator.Validator) *Event {
	v.Check(e.DistanceInMeters > 0, "event distance", "must be greater than 0")
	eventDate, _ := time.Parse(time.RFC3339, e.EventDate)

	// Splits
	v.Check(len(ss) != 0, "splits", "event must have at least one split")
	if !v.Valid() {
		return nil
	}
	splits := make([]*Split, 0, len(ss))
	for _, s := range ss {
		spl := NewSplit(s, trs, v)
		if !v.Valid() {
			return nil
		}
		splits = append(splits, spl)
	}
	var splitsNames []string
	splitTypeQty := make(map[SplitType]int)
	for _, split := range splits {
		splitsNames = append(splitsNames, split.Name)
		splitTypeQty[split.Type]++
	}
	v.Check(validator.Unique(splitsNames), "splits", "must have unique names for event")
	v.Check(splitTypeQty[SplitTypeStart] < 2, "split with type start", "must be 0 or 1")
	v.Check(splitTypeQty[SplitTypeFinish] == 1, "split with type finish", "must be only 1")

	// Waves
	v.Check(len(ww) > 0, "waves", "must be at least one for event")
	if !v.Valid() {
		return nil
	}
	waves := make([]*Wave, 0, len(ww))
	var wavesNames []string
	for _, w := range ww {
		wave := NewWave(w, v)
		if !v.Valid() {
			return nil
		}
		wavesNames = append(wavesNames, w.Name)
		waves = append(waves, wave)
	}
	v.Check(validator.Unique(wavesNames), "waves", "must have unique names for event")

	// Categories. Event may have no categories at all
	categories := make([]*Category, 0, len(cc))
	if len(cc) > 0 {
		var categoryNames []string
		for _, c := range cc {
			cat := NewCategory(c, eventDate, v)
			if !v.Valid() {
				return nil
			}
			categoryNames = append(categoryNames, c.Name)
			categories = append(categories, cat)
		}
		v.Check(validator.Unique(categoryNames), "categories", "must have unique names for event")
		CheckCategoriesBoundary(categories, v)
	}

	event := &Event{
		ID:               e.ID,
		RaceID:           e.RaceID,
		Name:             e.Name,
		DistanceInMeters: e.DistanceInMeters,
		EventDate:        eventDate,
		Splits:           splits,
		Waves:            waves,
		Categories:       categories,
	}
	event.AssignLapToSplits()
	fmt.Println(event)
	return event
}

type (
	SplitID  = uuid.UUID
	ReaderID = uuid.UUID
)

func (e *Event) AssignLapToSplits() {
	slices.SortFunc(e.Splits, func(a, b *Split) int {
		return cmp.Compare(a.DistanceFromStart, b.DistanceFromStart)
	})
	m := make(map[ReaderID]SplitID, len(e.Splits))
	for _, s := range e.Splits {
		if _, ok := m[s.TimeReaderID]; ok {
			s.PreviousLapSplitID = uuid.NullUUID{
				UUID:  m[s.TimeReaderID],
				Valid: true,
			}
			m[s.TimeReaderID] = s.ID
		}
	}
}

func (e Event) String() string {
	return fmt.Sprintf(
		"Event {\n"+
			"  ID: %s\n"+
			"  RaceID: %s\n"+
			"  Name: %q\n"+
			"  Distance: %d meters\n"+
			"  Event Date: %s\n"+
			"  Splits: %s\n"+
			"  Waves: %s\n"+
			"  Categories: %s\n"+
			"}",
		e.ID,
		e.RaceID,
		e.Name,
		e.DistanceInMeters,
		e.EventDate.Format(time.DateOnly),
		formatSlice(e.Splits),
		formatSlice(e.Waves),
		formatSlice(e.Categories),
	)
}

func formatSlice[T fmt.Stringer](items []T) string {
	if len(items) == 0 {
		return "[]"
	}

	result := "[\n"
	for _, item := range items {
		result += fmt.Sprintf("    %s,\n", item.String())
	}
	result += "  ]"
	return result
}
