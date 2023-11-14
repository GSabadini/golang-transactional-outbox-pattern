package repository

import (
	"database/sql"
	"time"
)

func NewNullInt64(i int64) sql.NullInt64 {
	if i == 0 {
		return sql.NullInt64{}
	}

	return sql.NullInt64{
		Int64: i,
		Valid: true,
	}
}

func NewNullFloat64(i float64) sql.NullFloat64 {
	if i == 0 {
		return sql.NullFloat64{}
	}

	return sql.NullFloat64{
		Float64: i,
		Valid:   true,
	}
}

func NewNullTime(t time.Time) sql.NullTime {
	if t.IsZero() {
		return sql.NullTime{}
	}

	return sql.NullTime{
		Time:  t,
		Valid: true,
	}
}

func NewNullString(s string) sql.NullString {
	if len(s) == 0 {
		return sql.NullString{}
	}

	return sql.NullString{
		String: s,
		Valid:  true,
	}
}

func NewNullBool(b bool) sql.NullBool {
	return sql.NullBool{
		Bool:  b,
		Valid: true,
	}
}
