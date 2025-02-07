package repo

import (
	"context"

	"github.com/ecoarchie/timeit/internal/database"
	"github.com/ecoarchie/timeit/internal/entity"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type RaceQuery interface {
	CreateRace(ctx context.Context, arg database.CreateRaceParams) (database.Race, error)
}

type RaceRepoPG struct {
	q RaceQuery
}

func NewRaceRepoPG(q RaceQuery) *RaceRepoPG {
	return &RaceRepoPG{
		q: q,
	}
}

func (rr RaceRepoPG) Create(ctx context.Context, r entity.Race) (uuid.UUID, error) {
	var date pgtype.Date
	if err := date.Scan(r.RaceDate); err != nil {
		return uuid.Nil, err
	}
	params := database.CreateRaceParams{
		ID:       r.ID,
		Name:     r.Name,
		RaceDate: date,
		Timezone: r.Timezone,
	}
	race, err := rr.q.CreateRace(ctx, params)
	if err != nil {
		return uuid.Nil, err
	}
	return race.ID, nil
}
