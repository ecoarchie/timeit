package pgxmapper

import (
	"time"

	"github.com/jackc/pgx/v5/pgtype"
)

func PgxIntervalToDuration(i pgtype.Interval) time.Duration {
	return time.Duration(i.Microseconds * 1000)
}

func DurationToPgxInterval(d time.Duration) pgtype.Interval {
	return pgtype.Interval{
		Microseconds: d.Microseconds(),
		Valid:        true,
	}
}

func TimeToPgxTimestamp(t time.Time) pgtype.Timestamp {
	return pgtype.Timestamp{
		Time:  t,
		Valid: true,
	}
}

func PgxTimestampToTime(ts pgtype.Timestamp) time.Time {
	return ts.Time
}

func TimeToPgxDate(t time.Time) pgtype.Date {
	return pgtype.Date{
		Time:  t,
		Valid: true,
	}
}

func StringToPgxText(s string) pgtype.Text {
	return pgtype.Text{
		String: s,
		Valid:  true,
	}
}
