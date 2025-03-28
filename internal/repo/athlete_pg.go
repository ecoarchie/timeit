package repo

import (
	"context"
	"fmt"

	"github.com/ecoarchie/timeit/internal/database"
	"github.com/ecoarchie/timeit/internal/entity"
	"github.com/ecoarchie/timeit/pkg/pgxmapper"
	"github.com/ecoarchie/timeit/pkg/postgres"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

type AthleteQuery interface {
	GetAthleteByID(ctx context.Context, id uuid.UUID) (database.GetAthleteByIDRow, error)
	CreateOrUpdateAthlete(ctx context.Context, arg database.CreateOrUpdateAthleteParams) (database.Athlete, error)
	AddChipBib(ctx context.Context, arg database.AddChipBibParams) (database.ChipBib, error)
	AddEventAthlete(ctx context.Context, arg database.AddEventAthleteParams) (database.EventAthlete, error)
	DeleteAthleteByID(ctx context.Context, athleteID uuid.UUID) error
	DeleteAthletesWithRaceID(ctx context.Context, raceID uuid.UUID) error
	DeleteAthletesWithEventID(ctx context.Context, eventID uuid.UUID) error
	DeleteChipBib(ctx context.Context, arg database.DeleteChipBibParams) error
	DeleteChipBibWithEventID(ctx context.Context, arg database.DeleteChipBibWithEventIDParams) error
	DeleteChipBibWithRaceID(ctx context.Context, raceID uuid.UUID) error
	DeleteAthleteSplit(ctx context.Context, arg database.DeleteAthleteSplitParams) error
	GetEventAthlete(ctx context.Context, athleteID uuid.UUID) (database.EventAthlete, error)
	GetCategoryForAthlete(ctx context.Context, arg database.GetCategoryForAthleteParams) (database.Category, error)
	GetEventAthleteRecordsC(ctx context.Context, arg database.GetEventAthleteRecordsCParams) ([]database.GetEventAthleteRecordsCRow, error)
	GetSplitsForEvent(ctx context.Context, eventID uuid.UUID) ([]database.Split, error)
	CreateAthleteSplits(ctx context.Context, arg database.CreateAthleteSplitsParams) error
	GetEventIDsWithWavesStarted(ctx context.Context, raceID uuid.UUID) ([]uuid.UUID, error)
	CreateAthleteBulk(ctx context.Context, arg []database.CreateAthleteBulkParams) (int64, error)
	AddChipBibBulk(ctx context.Context, arg []database.AddChipBibBulkParams) (int64, error)
	AddEventAthleteBulk(ctx context.Context, arg []database.AddEventAthleteBulkParams) (int64, error)
	GetSplitsForRace(ctx context.Context, raceID uuid.UUID) ([]database.Split, error)
	GetManualAthleteSplits(ctx context.Context, arg database.GetManualAthleteSplitsParams) ([]database.GetManualAthleteSplitsRow, error)
	SetStatus(ctx context.Context, arg database.SetStatusParams) error
	WithTx(tx pgx.Tx) *database.Queries
}

type AthleteRepoPG struct {
	q  AthleteQuery
	pg *postgres.Postgres
}

func NewAthleteRepoPG(q AthleteQuery, pg *postgres.Postgres) *AthleteRepoPG {
	return &AthleteRepoPG{
		q:  q,
		pg: pg,
	}
}

func (ar *AthleteRepoPG) WithTx(tx pgx.Tx) *AthleteRepoPG {
	return &AthleteRepoPG{
		q:  ar.q.WithTx(tx),
		pg: ar.pg,
	}
}

func (ar *AthleteRepoPG) SaveAthleteBulk(ctx context.Context, raceID uuid.UUID, athletes []*entity.Athlete) (int64, error) {
	// err = ar.DeleteAthletesForRace(ctx, raceID)
	// if err != nil {
	// 	return 0, fmt.Errorf("save athlete bulk: delete athletes for race with ID: %s", err.Error())
	// }
	createPms := make([]database.CreateAthleteBulkParams, 0, len(athletes))
	chipBibPms := make([]database.AddChipBibBulkParams, 0, len(athletes)) // FIXME if athlete has more than 1 chip, this must be rewritten
	eventAthletePms := make([]database.AddEventAthleteBulkParams, 0, len(athletes))
	for _, a := range athletes {
		ap := database.CreateAthleteBulkParams{
			ID:              a.ID,
			RaceID:          a.RaceID,
			FirstName:       pgxmapper.StringToPgxText(a.FirstName),
			LastName:        pgxmapper.StringToPgxText(a.LastName),
			Gender:          database.CategoryGender(a.Gender),
			DateOfBirth:     pgxmapper.TimeToPgxDate(a.DateOfBirth),
			Phone:           pgxmapper.StringToPgxText(a.Phone),
			AthleteComments: pgxmapper.StringToPgxText(a.Comments),
		}
		createPms = append(createPms, ap)

		cb := database.AddChipBibBulkParams{
			RaceID:  a.RaceID,
			EventID: a.EventID,
			Chip:    int32(a.Chip),
			Bib:     int32(a.Bib),
		}
		chipBibPms = append(chipBibPms, cb)

		ea := database.AddEventAthleteBulkParams{
			RaceID:     a.RaceID,
			EventID:    a.EventID,
			AthleteID:  a.ID,
			WaveID:     a.WaveID,
			CategoryID: a.CategoryID,
			Bib:        int32(a.Bib),
		}
		eventAthletePms = append(eventAthletePms, ea)
	}

	tx, err := ar.pg.Pool.Begin(ctx)
	if err != nil {
		return 0, err
	}
	defer tx.Rollback(ctx)
	qtx := ar.WithTx(tx)

	err = qtx.q.DeleteChipBibWithRaceID(ctx, raceID)
	if err != nil {
		return 0, fmt.Errorf("error deleting chipbib for race = %s", raceID)
	}

	err = qtx.q.DeleteAthletesWithRaceID(ctx, raceID)
	if err != nil {
		return 0, fmt.Errorf("error deleting athletes for race = %s", raceID)
	}

	createdCount, err := qtx.q.CreateAthleteBulk(ctx, createPms)
	if err != nil {
		return 0, fmt.Errorf("save athlete bulk: error creating athletes")
	}

	_, err = qtx.q.AddChipBibBulk(ctx, chipBibPms)
	if err != nil {
		return 0, fmt.Errorf("save athlete bulk: error creating chip bib")
	}

	_, err = qtx.q.AddEventAthleteBulk(ctx, eventAthletePms)
	if err != nil {
		return 0, fmt.Errorf("save athlete bulk: error creating event-athlete record")
	}

	err = tx.Commit(ctx)
	if err != nil {
		return 0, fmt.Errorf("save athlete bulk: transaction commit error")
	}

	return createdCount, nil
}

func (ar *AthleteRepoPG) SaveAthlete(ctx context.Context, p *entity.Athlete) error {
	tx, err := ar.pg.Pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)
	qtx := ar.WithTx(tx)
	aParams := database.CreateOrUpdateAthleteParams{
		ID:              p.ID,
		RaceID:          p.RaceID,
		FirstName:       pgxmapper.StringToPgxText(p.FirstName),
		LastName:        pgxmapper.StringToPgxText(p.LastName),
		Gender:          database.CategoryGender(p.Gender),
		DateOfBirth:     pgxmapper.TimeToPgxDate(p.DateOfBirth),
		Phone:           pgxmapper.StringToPgxText(p.Phone),
		AthleteComments: pgxmapper.StringToPgxText(p.Comments),
	}

	_, err = qtx.q.CreateOrUpdateAthlete(ctx, aParams)
	if err != nil {
		return err
	}
	cParams := database.AddChipBibParams{
		RaceID:  p.RaceID,
		EventID: p.EventID,
		Chip:    int32(p.Chip),
		Bib:     int32(p.Bib),
	}
	_, err = qtx.q.AddChipBib(ctx, cParams)
	if err != nil {
		return err
	}

	eaParams := database.AddEventAthleteParams{
		RaceID:     p.RaceID,
		EventID:    p.EventID,
		AthleteID:  p.ID,
		WaveID:     p.WaveID,
		CategoryID: p.CategoryID,
		Bib:        int32(p.Bib),
	}

	_, err = qtx.q.AddEventAthlete(ctx, eaParams)
	if err != nil {
		return err
	}
	return tx.Commit(ctx)
}

func (ar *AthleteRepoPG) GetEventIDsWithWavesStarted(ctx context.Context, raceID uuid.UUID) ([]uuid.UUID, error) {
	return ar.q.GetEventIDsWithWavesStarted(ctx, raceID)
}

func (ar *AthleteRepoPG) GetCategoryFor(ctx context.Context, p *entity.Athlete) (uuid.NullUUID, bool, error) {
	params := database.GetCategoryForAthleteParams{
		EventID:  p.EventID,
		Gender:   database.CategoryGender(p.Gender),
		DateFrom: pgxmapper.TimeToPgxTimestamp(p.DateOfBirth),
	}

	c, err := ar.q.GetCategoryForAthlete(ctx, params)
	if err != nil {
		if c.ID == uuid.Nil {
			return uuid.NullUUID{}, false, nil
		}
		return uuid.NullUUID{}, false, err
	}
	return uuid.NullUUID{
		UUID:  c.ID,
		Valid: true,
	}, true, nil
}

func (ar *AthleteRepoPG) GetAthleteWithChip(chip int) (*entity.Athlete, error) {
	panic("not implemented")
}

func (ar *AthleteRepoPG) GetAthleteByID(ctx context.Context, athleteID uuid.UUID) (*entity.Athlete, error) {
	a, err := ar.q.GetAthleteByID(ctx, athleteID)
	if err != nil {
		return nil, err
	}

	athlete := &entity.Athlete{
		ID:          a.ID,
		RaceID:      a.RaceID,
		EventID:     a.EventID,
		WaveID:      a.WaveID,
		Bib:         int(a.Bib),
		Chip:        int(a.Chip),
		FirstName:   a.FirstName.String,
		LastName:    a.LastName.String,
		Gender:      entity.CategoryGender(a.Gender),
		DateOfBirth: a.DateOfBirth.Time,
		CategoryID:  a.CategoryID,
		Phone:       a.Phone.String,
		Comments:    a.AthleteComments.String,
	}
	return athlete, nil
}

// func (ar *AthleteRepoPG) TruncateAndSaveBulkAthletes(ctx context.Context, raceID uuid.UUID, [])

func (ar *AthleteRepoPG) DeleteAthletesForRace(ctx context.Context, raceID uuid.UUID) error {
	tx, err := ar.pg.Pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)
	qtx := ar.WithTx(tx)

	err = qtx.q.DeleteChipBibWithRaceID(ctx, raceID)
	if err != nil {
		return fmt.Errorf("error deleting chipbib for race = %s", raceID)
	}

	err = qtx.q.DeleteAthletesWithRaceID(ctx, raceID)
	if err != nil {
		return fmt.Errorf("error deleting athletes for race = %s", raceID)
	}

	return tx.Commit(ctx)
}

func (ar *AthleteRepoPG) DeleteAthletesForRaceWithEventID(ctx context.Context, raceID, eventID uuid.UUID) error {
	tx, err := ar.pg.Pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)
	qtx := ar.WithTx(tx)

	dParams := database.DeleteChipBibWithEventIDParams{
		RaceID:  raceID,
		EventID: eventID,
	}
	err = qtx.q.DeleteChipBibWithEventID(ctx, dParams)
	if err != nil {
		return fmt.Errorf("error deleting chipbib for eventID = %s", eventID)
	}

	err = qtx.q.DeleteAthletesWithEventID(ctx, eventID)
	if err != nil {
		return fmt.Errorf("error deleting athletes for race = %s", raceID)
	}

	return tx.Commit(ctx)
}

func (ar *AthleteRepoPG) DeleteAthlete(ctx context.Context, a *entity.Athlete) error {
	tx, err := ar.pg.Pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)
	qtx := ar.WithTx(tx)
	cbParams := database.DeleteChipBibParams{
		RaceID: a.RaceID,
		Chip:   int32(a.Chip),
		Bib:    int32(a.Bib),
	}
	err = qtx.q.DeleteChipBib(ctx, cbParams)
	if err != nil {
		return err
	}

	err = qtx.q.DeleteAthleteSplit(ctx, database.DeleteAthleteSplitParams{
		RaceID:    a.RaceID,
		AthleteID: a.ID,
	})
	if err != nil {
		return err
	}

	err = qtx.q.DeleteAthleteByID(ctx, a.ID)
	if err != nil {
		return err
	}
	return tx.Commit(ctx)
}

func (ar *AthleteRepoPG) SaveAthleteSplits(ctx context.Context, as []database.CreateAthleteSplitsParams) error {
	for _, sp := range as {
		err := ar.q.CreateAthleteSplits(ctx, sp)
		if err != nil {
			fmt.Printf("error creating athleteID split: %s\n", sp.AthleteID)
			return err
		}
	}
	return nil
}

const tempTableCreate = `
	CREATE TEMPORARY TABLE athlete_split_tmp (
		LIKE athlete_split INCLUDING ALL, visited BOOLEAN
	) ON COMMIT DROP;
`

const mergeWithRanks = `
	merge into athlete_split asl
	using (
		select
			ats.race_id,
			ats.event_id,
			ats.split_id,
			ats.athlete_id,
			ats.tod,
			ats.gun_time,
			ats.net_time,
			ats.visited,
			CASE
				WHEN ss.status_id IN (2, 3) and a.gender <> 'unknown' and ats.visited is true THEN
				RANK() OVER (PARTITION BY ats.race_id, ats.event_id, ats.split_id, a.gender, ats.visited ORDER BY ats.gun_time ASC)
			END AS gun_rank_gender,
			CASE 
						WHEN ea.category_id IS NOT NULL AND ss.status_id IN (2, 3) and ats.visited is true THEN
						RANK() OVER (PARTITION BY ats.race_id, ats.event_id, ats.split_id, ea.category_id, ats.visited ORDER BY ats.gun_time ASC) 
			END AS gun_rank_category,
			CASE
				WHEN ss.status_id IN (2, 3) and ats.visited is true THEN
				RANK() OVER (PARTITION BY ats.race_id, ats.event_id, ats.split_id, ats.visited ORDER BY ats.gun_time ASC)
			END AS gun_rank_overall,
			CASE
				WHEN ss.status_id IN (2, 3) and a.gender <> 'unknown' and ats.visited is true THEN
				RANK() OVER (PARTITION BY ats.race_id, ats.event_id, ats.split_id, a.gender, ats.visited ORDER BY ats.net_time ASC)
			END AS net_rank_gender,
			CASE 
					WHEN ea.category_id IS NOT NULL AND ss.status_id IN (2, 3) and ats.visited is true THEN
					RANK() OVER (PARTITION BY ats.race_id, ats.event_id, ats.split_id, ea.category_id, ats.visited ORDER BY ats.net_time ASC) 
			END AS net_rank_category,
			CASE
				WHEN ss.status_id IN (2, 3) and ats.visited is true THEN
				RANK() OVER (PARTITION BY ats.race_id, ats.event_id, ats.split_id, ats.visited ORDER BY ats.net_time ASC)
			END AS net_rank_overall
		from athlete_split_tmp ats
		join athletes a on a.id = ats.athlete_id
		join event_athlete ea on ea.athlete_id  = ats.athlete_id and ea.race_id = ats.race_id and ea.event_id = ats.event_id
		join statuses ss on ea.status_id = ss.status_id
		join splits s on s.id = ats.split_id and s.race_id = ats.race_id and s.event_id = ats.event_id
		order by
			ats.gun_time
	) ats
	on asl.race_id = ats.race_id and asl.event_id = ats.event_id and asl.split_id = ats.split_id and asl.athlete_id = ats.athlete_id
	when matched and ats.visited is FALSE then
		DELETE
	when matched then update set
		tod = ats.tod,
		gun_time = ats.gun_time,
		net_time = ats.net_time,
		gun_rank_gender = ats.gun_rank_gender,
		gun_rank_category = ats.gun_rank_category,
		gun_rank_overall = ats.gun_rank_overall,
		net_rank_gender = ats.net_rank_gender,
		net_rank_category = ats.net_rank_category,
		net_rank_overall = ats.net_rank_overall
	when not matched and ats.visited is FALSE then DO NOTHING 
	when not matched then insert 
		(race_id, event_id, split_id, athlete_id, tod, gun_time, net_time, gun_rank_gender, gun_rank_category, gun_rank_overall, net_rank_gender, net_rank_category, net_rank_overall)
		values (ats.race_id, ats.event_id, ats.split_id, ats.athlete_id, ats.tod, ats.gun_time, ats.net_time, ats.gun_rank_gender, ats.gun_rank_category, ats.gun_rank_overall, ats.net_rank_gender, ats.net_rank_category, ats.net_rank_overall)
`

func (ar *AthleteRepoPG) SaveBulkAthleteSplits(ctx context.Context, raceID uuid.UUID, as []*entity.AthleteSplit) error {
	tx, err := ar.pg.Pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	_, err = tx.Exec(ctx, tempTableCreate)
	if err != nil {
		fmt.Println("Error executing creation temp table: ", err)
		return err
	}

	var linkedParams [][]interface{}
	for _, p := range as {
		if p != nil {
			linkedParams = append(linkedParams, []interface{}{p.RaceID, p.EventID, p.SplitID, p.AthleteID, p.TOD, p.GunTime, p.NetTime, p.IsVisited()})
		}
	}
	_, err = tx.CopyFrom(ctx, []string{"athlete_split_tmp"}, []string{"race_id", "event_id", "split_id", "athlete_id", "tod", "gun_time", "net_time", "visited"}, pgx.CopyFromRows(linkedParams))
	if err != nil {
		fmt.Println("Error executing copyfrom athlete splits: ", err)
		return err
	}

	_, err = tx.Exec(ctx, mergeWithRanks)
	if err != nil {
		fmt.Println("Error executing rank query: ", err)
		return err
	}
	err = tx.Commit(ctx)
	if err != nil {
		fmt.Println("Error commiting transaction: ", err)
		return err
	}
	return nil
}

func (ar *AthleteRepoPG) GetManualAthleteSplits(ctx context.Context, raceID, eventID uuid.UUID) (map[uuid.UUID][]*entity.AthleteSplit, error) {
	params := database.GetManualAthleteSplitsParams{
		RaceID:  raceID,
		EventID: eventID,
	}
	ms, err := ar.q.GetManualAthleteSplits(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("get manual athlete splits")
	}
	res := make(map[uuid.UUID][]*entity.AthleteSplit)
	for _, m := range ms {
		res[m.AthleteID] = append(res[m.AthleteID], &entity.AthleteSplit{
			RaceID:     m.RaceID,
			EventID:    m.EventID,
			AthleteID:  m.AthleteID,
			SplitID:    m.SplitID,
			TOD:        pgxmapper.PgxTimestampToTime(m.Tod),
			GunTime:    pgxmapper.PgxIntervalToDuration(m.GunTime),
			NetTime:    pgxmapper.PgxIntervalToDuration(m.NetTime),
			Gender:     entity.CategoryGender(m.Gender),
			CategoryID: m.CategoryID,
		})
	}
	return res, nil
}

func (ar *AthleteRepoPG) GetAthleteSplitResults(ctx context.Context, raceID uuid.UUID) error {
	splits, err := ar.q.GetSplitsForRace(ctx, raceID)
	if err != nil {
		return err
	}

	m := map[uuid.UUID][]database.Split{}
	for _, s := range splits {
		m[s.EventID] = append(m[s.EventID], s)
	}

	// res := map[uuid.UUID][]entity.AthleteSplitResults{}

	// FIXME complete the func
	return nil
}

func (ar *AthleteRepoPG) GetRecordsAndSplitsForEventAthlete(ctx context.Context, raceID, eventID uuid.UUID) ([]database.GetEventAthleteRecordsCRow, []*entity.Split, error) {
	ss, err := ar.q.GetSplitsForEvent(ctx, eventID)
	if err != nil {
		return nil, nil, err
	}

	eaParams := database.GetEventAthleteRecordsCParams{
		RaceID:  raceID,
		EventID: eventID,
	}
	records, err := ar.q.GetEventAthleteRecordsC(ctx, eaParams)
	if err != nil {
		return nil, nil, err
	}

	splits := []*entity.Split{}
	for _, s := range ss {
		split := &entity.Split{
			ID:                 s.ID,
			RaceID:             s.RaceID,
			EventID:            s.EventID,
			Name:               s.SplitName,
			Type:               entity.SplitType(s.SplitType),
			DistanceFromStart:  int(s.DistanceFromStart),
			TimeReaderID:       s.TimeReaderID,
			MinTime:            pgxmapper.PgxIntervalToDuration(s.MinTime),
			MaxTime:            pgxmapper.PgxIntervalToDuration(s.MaxTime),
			MinLapTime:         pgxmapper.PgxIntervalToDuration(s.MinLapTime),
			PreviousLapSplitID: s.PreviousLapSplitID,
		}
		splits = append(splits, split)
	}

	return records, splits, nil
}

func (ar *AthleteRepoPG) UpdateStatus(ctx context.Context, status entity.Status, raceID, eventID, athleteID uuid.UUID) error {
	var statusID pgtype.Int4
	switch status {
	case entity.NYS:
		statusID = pgtype.Int4{
			Int32: 1,
			Valid: true,
		}
	case entity.RUN:
		statusID = pgtype.Int4{
			Int32: 2,
			Valid: true,
		}
	case entity.FIN:
		statusID = pgtype.Int4{
			Int32: 3,
			Valid: true,
		}
	case entity.DSQ:
		statusID = pgtype.Int4{
			Int32: 4,
			Valid: true,
		}
	case entity.QRT:
		statusID = pgtype.Int4{
			Int32: 5,
			Valid: true,
		}
	case entity.DNS:
		statusID = pgtype.Int4{
			Int32: 6,
			Valid: true,
		}
	case entity.DNF:
		statusID = pgtype.Int4{
			Int32: 7,
			Valid: true,
		}
	default:
		statusID = pgtype.Int4{
			Int32: 1,
			Valid: true,
		}
	}

	sParam := database.SetStatusParams{
		StatusID:  statusID,
		AthleteID: athleteID,
		RaceID:    raceID,
		EventID:   eventID,
	}
	err := ar.q.SetStatus(ctx, sParam)
	if err != nil {
		return err
	}
	return nil
}
