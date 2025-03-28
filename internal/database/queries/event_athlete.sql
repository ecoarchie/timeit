-- name: AddEventAthlete :one
INSERT INTO event_athlete
(race_id, event_id, athlete_id, wave_id, category_id, bib)
VALUES ($1, $2, $3, $4, $5, $6)
ON CONFLICT (race_id, event_id, athlete_id)
DO UPDATE
SET event_id=EXCLUDED.event_id, wave_id=EXCLUDED.wave_id, category_id=EXCLUDED.category_id, bib=EXCLUDED.bib 
RETURNING *;

-- name: AddEventAthleteBulk :copyfrom
INSERT INTO event_athlete
(race_id, event_id, athlete_id, wave_id, category_id, bib)
VALUES ($1, $2, $3, $4, $5, $6);

-- name: GetEventAthlete :one
SELECT race_id, event_id, athlete_id, wave_id, category_id, bib, status_id
FROM event_athlete
WHERE athlete_id=$1;

-- name: SetStatus :exec
UPDATE event_athlete
SET status_id = $1
WHERE athlete_id = $2 AND race_id = $3 AND event_id = $4;

-- name: GetEventAthleteRecordsC :many
with distinct_rr_tod as (
    select distinct tr.id, rr.tod, rr.chip, rr.race_id
    from reader_records rr
    join time_readers tr on
        tr.reader_name = rr.reader_name
        and tr.race_id = rr.race_id
    where rr.can_use is true
)
select 
    ea.athlete_id,
    ea.category_id,
    ea.bib,
    cb.chip,
    a.gender,
    s.status_full,
    w.start_time as wave_start,
    (
        select array_agg(row(d.id, d.tod)::rr_tod order by d.tod)::rr_tod[]
        from distinct_rr_tod d
        where d.race_id = ea.race_id
          and d.chip = cb.chip
    ) as rr_tod
from
    event_athlete ea
join statuses s on ea.status_id = s.status_id
join waves w on
    w.race_id = ea.race_id
    and w.event_id = ea.event_id
    and w.id = ea.wave_id
join chip_bib cb on
    cb.race_id = ea.race_id
    and cb.event_id = ea.event_id
    and cb.bib = ea.bib
join athletes a on
    a.id = ea.athlete_id
    and a.race_id = ea.race_id
	where ea.race_id = $1 
		and ea.event_id = $2 
		and w.is_launched is true;
