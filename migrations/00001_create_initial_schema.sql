-- +goose Up
-- +goose StatementBegin
-- Table: races
CREATE TABLE races (
  id UUID PRIMARY KEY,
  name TEXT NOT NULL UNIQUE,
  race_date DATE NOT NULL,
  timezone TEXT NOT NULL
);

-- Table: events
CREATE TABLE events (
  id UUID PRIMARY KEY,
  race_id UUID NOT NULL REFERENCES races(id) ON DELETE CASCADE,
  name TEXT NOT NULL,
  distance_in_meters INTEGER NOT NULL,
  event_date TIMESTAMPTZ NOT NULL,
  UNIQUE (race_id, id)
);

-- Table: waves
CREATE TABLE waves (
  id UUID NOT NULL,
  race_id UUID NOT NULL REFERENCES races(id) ON DELETE CASCADE,
  event_id UUID NOT NULL REFERENCES events(id) ON DELETE CASCADE,
  name VARCHAR NOT NULL,
  start_time TIMESTAMPTZ NOT NULL,
  is_launched BOOLEAN NOT NULL DEFAULT FALSE,
  PRIMARY KEY (race_id, event_id, id)
);

-- Enum: category_gender
DROP TYPE IF EXISTS category_gender;
CREATE TYPE category_gender AS ENUM ('male', 'female', 'mixed', 'unknown');

-- Table: categories
CREATE TABLE categories (
  id UUID NOT NULL,
  race_id UUID NOT NULL REFERENCES races(id) ON DELETE CASCADE,
  event_id UUID NOT NULL REFERENCES events(id) ON DELETE CASCADE,
  name TEXT NOT NULL,
  gender category_gender NOT NULL,
  from_age INTEGER NOT NULL,
  from_race_date BOOLEAN NOT NULL DEFAULT FALSE,
  to_age INTEGER NOT NULL,
  to_race_date BOOLEAN NOT NULL DEFAULT FALSE,
  PRIMARY KEY (race_id, event_id, id),
  CHECK (from_age <= to_age)
);


-- Table: box_records
CREATE TABLE box_records (
  id SERIAL PRIMARY KEY,
  race_id UUID NOT NULL REFERENCES races(id),
  chip INTEGER NOT NULL,
  tod TIMESTAMPTZ NOT NULL,
  box_name TEXT NOT NULL,
  can_use BOOLEAN NOT NULL DEFAULT TRUE
);

-- Performance indexes for box_records
CREATE INDEX idx_box_records_chip_box_name ON box_records (chip, box_name);
CREATE INDEX idx_box_records_race_chip ON box_records (race_id, chip);

-- Enum: tp_type
DROP TYPE IF EXISTS tp_type;
CREATE TYPE tp_type AS ENUM ('start', 'standard', 'finish');


-- Table: timing_points
CREATE TABLE timing_points (
  id UUID NOT NULL,
  race_id UUID NOT NULL,
  event_id UUID NOT NULL,
  name TEXT NOT NULL,
  type tp_type NOT NULL,
  distance_from_start INTEGER NOT NULL,
  box_name TEXT NOT NULL,
  min_time_sec INTEGER DEFAULT 0,
  max_time_sec INTEGER DEFAULT 0,
  min_lap_time_sec INTEGER DEFAULT 0,
  PRIMARY KEY (race_id, event_id, id),
  FOREIGN KEY (race_id, event_id) REFERENCES events (race_id, id) ON DELETE CASCADE
);

-- Table: participants
CREATE TABLE participants (
  id UUID PRIMARY KEY,
  first_name TEXT DEFAULT 'athlete',
  last_name TEXT DEFAULT 'unknown',
  gender category_gender NOT NULL DEFAULT 'unknown',
  date_of_birth DATE,
  phone TEXT,
  created_at TIMESTAMPTZ DEFAULT NOW(),
  updated_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_participants_first_last ON participants (first_name, last_name);

-- Table: event_participant
CREATE TABLE event_participant (
  race_id UUID NOT NULL,
  event_id UUID NOT NULL,
  participant_id UUID NOT NULL REFERENCES participants(id),
  wave_id UUID NOT NULL,
  category_id UUID NOT NULL,
  bib INTEGER,
  PRIMARY KEY (race_id, event_id, participant_id)
);

-- Table: chip_bib
CREATE TABLE chip_bib (
  race_id UUID NOT NULL,
  event_id UUID NOT NULL,
  chip INTEGER NOT NULL,
  bib INTEGER NOT NULL,
  PRIMARY KEY (chip, bib, race_id, event_id)
);

-- Table: physical_locations
CREATE TABLE physical_locations (
  race_id UUID NOT NULL REFERENCES races(id) ON DELETE CASCADE,
  box_name TEXT NOT NULL,
  PRIMARY KEY (race_id, box_name)
);

-- Table: event_location
CREATE TABLE event_location (
  event_id UUID NOT NULL REFERENCES events(id),
  race_id UUID NOT NULL,
  box_name TEXT NOT NULL,
  PRIMARY KEY (event_id, race_id, box_name),
  FOREIGN KEY (race_id, box_name) REFERENCES physical_locations(race_id, box_name) 
);

-- Table: participant_timing_point
CREATE TABLE participant_timing_point (
  race_id UUID NOT NULL,
  event_id UUID NOT NULL,
  timing_point_id UUID NOT NULL,
  participant_id UUID NOT NULL,
  tod TIMESTAMP NOT NULL,
  gun_time BIGINT NOT NULL,
  net_time BIGINT NOT NULL,
  PRIMARY KEY (race_id, event_id, timing_point_id, participant_id)
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS box_records;
DROP TABLE IF EXISTS chip_bib;
DROP TABLE IF EXISTS event_location;
DROP TABLE IF EXISTS event_participant;
DROP TABLE IF EXISTS physical_locations;
DROP TABLE IF EXISTS participant_timing_point;
DROP TABLE IF EXISTS participants;
DROP TABLE IF EXISTS waves;
DROP TABLE IF EXISTS categories;
DROP TABLE IF EXISTS timing_points;
DROP TABLE IF EXISTS events;
DROP TABLE IF EXISTS races;

DROP TYPE IF EXISTS category_gender;
DROP TYPE IF EXISTS tp_type;
-- +goose StatementEnd
