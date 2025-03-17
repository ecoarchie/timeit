-- name: AddEventAthlete :one
INSERT INTO event_athlete
(race_id, event_id, athlete_id, wave_id, category_id, bib)
VALUES ($1, $2, $3, $4, $5, $6)
ON CONFLICT (race_id, event_id, athlete_id)
DO UPDATE
SET event_id=EXCLUDED.event_id, wave_id=EXCLUDED.wave_id, category_id=EXCLUDED.category_id, bib=EXCLUDED.bib 
RETURNING *;

-- name: GetEventAthlete :one
SELECT race_id, event_id, athlete_id, wave_id, category_id, bib
FROM event_athlete
WHERE athlete_id=$1;

-- name: GetEventAthleteRecords :many
select 
	ea.athlete_id,
	ea.bib,
	cb.chip,
	ea.category_id,
	a.gender,
	w.start_time as wave_start,
	(select array_agg(rr.tod order by rr.tod)::timestamp[]
	from reader_records rr
	where rr.race_id = ea.race_id and rr.chip = cb.chip and rr.can_use is true) as records,
	(select array_agg(tr.id order by rr.tod)::uuid[]
	from reader_records rr
	join time_readers tr on tr.reader_name = rr.reader_name and tr.race_id = rr.race_id
	where rr.race_id = ea.race_id and rr.chip = cb.chip and rr.can_use is true) as reader_ids
	from event_athlete ea
	join waves w on w.race_id = ea.race_id and w.event_id = ea.event_id  and w.id = ea.wave_id
	join chip_bib cb on cb.race_id = ea.race_id and cb.event_id = ea.event_id and cb.bib = ea.bib
	join athletes a on a.id = ea.athlete_id and a.race_id = ea.race_id
	where ea.race_id = $1 
		and ea.event_id = $2 
		and w.is_launched is true;

-- name: GetEventAthleteRecordsJSON :many
select 
	ea.race_id,
	ea.event_id,
	ea.wave_id,
	ea.athlete_id,
	ea.bib,
	cb.chip,
	w.start_time as wave_start,
	(select json_agg(json_build_object('tod', rr.tod, 'reader', tr.id)) from reader_records rr 
	join time_readers tr on tr.reader_name = rr.reader_name and tr.race_id = rr.race_id
	where rr.race_id = ea.race_id and rr.chip = cb.chip and rr.can_use is true) as j_recs
from event_athlete ea
join waves w on w.race_id = ea.race_id and w.event_id = ea.event_id  and w.id = ea.wave_id
join chip_bib cb on cb.race_id = ea.race_id and cb.event_id = ea.event_id and cb.bib = ea.bib
where ea.race_id = $1
	and ea.event_id = $2
	and w.is_launched is true;

-- name: GetEventAthleteRecordsC :many
select 
	ea.athlete_id,
	ea.bib,
	cb.chip,
	ea.category_id,
	a.gender,
	w.start_time as wave_start,
	(select array_agg(row(tr.id, rr.tod)::rr_tod order by rr.tod)::rr_tod[]
	from
		reader_records rr
	join time_readers tr on
		tr.reader_name = rr.reader_name
		and tr.race_id = rr.race_id
	where
		rr.race_id = ea.race_id
		and rr.chip = cb.chip
		and rr.can_use is true) as rr_tod
	from event_athlete ea
	join waves w on w.race_id = ea.race_id and w.event_id = ea.event_id  and w.id = ea.wave_id
	join chip_bib cb on cb.race_id = ea.race_id and cb.event_id = ea.event_id and cb.bib = ea.bib
	join athletes a on a.id = ea.athlete_id and a.race_id = ea.race_id
	where ea.race_id = $1 
		and ea.event_id = $2 
		and w.is_launched is true;
