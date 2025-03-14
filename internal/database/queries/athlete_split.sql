-- name: CreateAthleteSplits :copyfrom
INSERT INTO athlete_split
(race_id, event_id, split_id, athlete_id, tod, gun_time, net_time)
VALUES($1, $2, $3, $4, $5, $6, $7);