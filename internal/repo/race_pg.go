package repo

import (
	"context"

	"github.com/ecoarchie/timeit/internal/database"
	"github.com/ecoarchie/timeit/internal/entity"
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

func (rr RaceRepoPG) SaveRaceConfig(ctx context.Context, r entity.RaceConfig) error {
	return nil
}
