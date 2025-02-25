-- +goose Up
-- +goose StatementBegin
-- Table: races
CREATE TABLE races (
  id UUID PRIMARY KEY,
  name TEXT NOT NULL UNIQUE,
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
  PRIMARY KEY (id),
  CHECK (from_age <= to_age),
  UNIQUE (id, race_id, event_id)
);


-- Table: reader_records
CREATE TABLE reader_records (
  id SERIAL PRIMARY KEY,
  race_id UUID NOT NULL,
  chip INTEGER NOT NULL,
  tod TIMESTAMPTZ NOT NULL,
  reader_name TEXT NOT NULL,
  can_use BOOLEAN NOT NULL DEFAULT TRUE
);

-- Performance indexes for reader_records
CREATE INDEX idx_box_records_chip_reader_name ON reader_records (chip, reader_name);
CREATE INDEX idx_box_records_race_chip ON reader_records (race_id, chip);

-- Table: time_readers
CREATE TABLE time_readers (
  id UUID NOT NULL,
  race_id UUID NOT NULL REFERENCES races(id) ON DELETE CASCADE,
  reader_name TEXT NOT NULL,
  PRIMARY KEY (id),
  UNIQUE (race_id, reader_name)
);

-- Enum: tp_type
DROP TYPE IF EXISTS tp_type;
CREATE TYPE tp_type AS ENUM ('start', 'standard', 'finish');


-- Table: splits
CREATE TABLE splits (
  id UUID NOT NULL,
  race_id UUID NOT NULL,
  event_id UUID NOT NULL,
  name TEXT NOT NULL,
  type tp_type NOT NULL,
  distance_from_start INTEGER NOT NULL,
  time_reader_id UUID NOT NULL,
  min_time_sec INTEGER DEFAULT 0,
  max_time_sec INTEGER DEFAULT 0,
  min_lap_time_sec INTEGER DEFAULT 0,
  PRIMARY KEY (race_id, event_id, id),
  FOREIGN KEY (race_id, event_id) REFERENCES events (race_id, id) ON DELETE CASCADE,
  FOREIGN KEY (time_reader_id) REFERENCES time_readers(id) ON DELETE CASCADE
);

-- Table: athletes
CREATE TABLE athletes (
  id UUID PRIMARY KEY,
  race_id UUID NOT NULL,
  first_name TEXT DEFAULT 'athlete',
  last_name TEXT DEFAULT 'unknown',
  gender category_gender NOT NULL DEFAULT 'unknown',
  date_of_birth DATE,
  phone TEXT,
  comments TEXT,
  created_at TIMESTAMPTZ DEFAULT NOW(),
  updated_at TIMESTAMPTZ DEFAULT NOW(),
  FOREIGN KEY (race_id) REFERENCES races (id) ON DELETE CASCADE
);

CREATE INDEX idx_athletes_first_last ON athletes (first_name, last_name);

-- Table: event_athlete
CREATE TABLE event_athlete (
  race_id UUID NOT NULL,
  event_id UUID NOT NULL,
  athlete_id UUID NOT NULL REFERENCES athletes(id),
  wave_id UUID NOT NULL,
  category_id UUID,
  bib INTEGER,
  PRIMARY KEY (race_id, event_id, athlete_id),
  FOREIGN KEY (race_id) REFERENCES races (id) ON DELETE CASCADE
);

-- Table: chip_bib
CREATE TABLE chip_bib (
  race_id UUID NOT NULL REFERENCES races(id) ON DELETE CASCADE,
  event_id UUID NOT NULL,
  chip INTEGER NOT NULL,
  bib INTEGER NOT NULL,
  PRIMARY KEY (chip, bib, race_id, event_id)
);

-- Table: athlete_split
CREATE TABLE athlete_split (
  race_id UUID NOT NULL,
  event_id UUID NOT NULL,
  split_id UUID NOT NULL,
  athlete_id UUID NOT NULL,
  tod TIMESTAMP NOT NULL,
  gun_time BIGINT NOT NULL,
  net_time BIGINT NOT NULL,
  PRIMARY KEY (race_id, event_id, split_id, athlete_id),
  FOREIGN KEY (race_id) REFERENCES races(id) ON DELETE CASCADE 
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS event_athlete;
DROP TABLE IF EXISTS splits;
DROP TABLE IF EXISTS time_readers;
DROP TABLE IF EXISTS athlete_split;
DROP TABLE IF EXISTS athletes;
DROP TABLE IF EXISTS waves;
DROP TABLE IF EXISTS categories;
DROP TABLE IF EXISTS events;
DROP TABLE IF EXISTS reader_records;
DROP TABLE IF EXISTS chip_bib;
DROP TABLE IF EXISTS races;

DROP TYPE IF EXISTS category_gender;
DROP TYPE IF EXISTS tp_type;
-- +goose StatementEnd
