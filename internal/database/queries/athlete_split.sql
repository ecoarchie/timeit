-- name: CreateAthleteSplits :exec
INSERT INTO athlete_split
(race_id, event_id, split_id, athlete_id, tod, gun_time, net_time)
VALUES($1, $2, $3, $4, $5, $6, $7)
ON CONFLICT (race_id, event_id, split_id, athlete_id) DO UPDATE
SET tod=EXCLUDED.tod, gun_time=EXCLUDED.gun_time, net_time=EXCLUDED.net_time;

-- name: GetManualAthleteSplits :many
SELECT ast.race_id, ast.event_id, ast.split_id, ast.athlete_id, ast.tod, ast.gun_time, ast.net_time, ea.category_id, a.gender
FROM athlete_split ast
join event_athlete ea on ea.athlete_id = ast.athlete_id and ea.race_id = ast.race_id and ea.event_id = ast.event_id
join athletes a on ea.athlete_id = a.id
WHERE ast.race_id = $1 AND ast.event_id = $2 AND is_manual IS TRUE;

-- name: DeleteAthleteSplit :exec
DELETE FROM athlete_split
WHERE race_id = $1 AND athlete_ID = $2;