-- name: CreateAthleteSplits :exec
INSERT INTO athlete_split
(race_id, event_id, split_id, athlete_id, tod, gun_time, net_time)
VALUES($1, $2, $3, $4, $5, $6, $7)
ON CONFLICT (race_id, event_id, split_id, athlete_id) DO UPDATE
SET tod=EXCLUDED.tod, gun_time=EXCLUDED.gun_time, net_time=EXCLUDED.net_time;