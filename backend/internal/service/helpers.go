package service

import (
	"database/sql"
	"strconv"

	"github.com/jbechler2/grant-tool/backend/internal/repository"
)

func toNullStringFromPtr(s *string) sql.NullString {
	if s == nil {
		return sql.NullString{Valid: false}
	}

	return sql.NullString{
		String: *s,
		Valid:  *s != "",
	}
}

func toNullString(s string) sql.NullString {
	return sql.NullString{
		String: s,
		Valid:  s != "",
	}
}

func float64ToNullString(f *float64) sql.NullString {
	if f == nil {
		return sql.NullString{Valid: false}
	}
	return sql.NullString{
		String: strconv.FormatFloat(*f, 'f', 2, 64),
		Valid:  true,
	}
}

func nullStringToFloat64(s sql.NullString) *float64 {
	if !s.Valid {
		return nil
	}
	f, err := strconv.ParseFloat(s.String, 64)
	if err != nil {
		return nil
	}
	return &f
}

func mergeString(input *string, existing string) string {
	if input == nil {
		return existing
	}

	if *input == "" {
		return existing
	}

	return *input
}

func mergeFloat64(input *float64, existing sql.NullString) sql.NullString {
	if input == nil {
		return existing
	}
	return float64ToNullString(input)
}

func mergeVisibility(input *string, existing repository.GrantVisibility) repository.GrantVisibility {
	if input == nil {
		return existing
	}
	return repository.GrantVisibility(*input)
}
