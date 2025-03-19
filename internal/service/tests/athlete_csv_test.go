package service

import (
	"slices"
	"testing"
)

func TestValidHeaders(t *testing.T) {
	parser := NewAthleteCSVParser("data.csv", ";")

	validHeadersToTest := []string{"event", "wave", "bib", "tag", "name", "surname", "gender", "date of birth", "phone", "comments"}
	got := parser.ValidateHeaders(validHeadersToTest)
	if got != nil {
		t.Error("error headers slice for valid headers must be empty")
	}
}

func TestInvalidHeaders(t *testing.T) {
	parser := NewAthleteCSVParser("data.csv", ";")
	t.Run("not nil error headers", func(t *testing.T) {
		invalidHeadersToTest := []string{"invalid header", "event", "wave", "bib", "tag", "name", "surname", "gender", "date of birth", "phone", "comments"}
		got := parser.ValidateHeaders(invalidHeadersToTest)
		if got == nil {
			t.Error("len of erros headers slice must be greater than 0")
		}
	})

	t.Run("invalid headers", func(t *testing.T) {
		invalidHeadersToTest := []string{"invalid header", "event", "test", "wave", "bib", "tag", "name", "surname", "gender", "date of birth", "phone", "comments"}
		got := parser.ValidateHeaders(invalidHeadersToTest)
		want := []string{"invalid header", "test"}
		if !slices.Equal(got, want) {
			t.Errorf("got %v, want %v", got, want)
		}
	})

	t.Run("error headers len", func(t *testing.T) {
		validHeadersToTest := []string{"invalid header", "event", "test", "wave", "test2", "bib", "tag", "name", "surname", "gender", "date of birth", "phone", "comments"}
		got := parser.ValidateHeaders(validHeadersToTest)
		if len(got) != 3 {
			t.Errorf("Len of got - %d,  want - %d", len(got), 3)
		}
	})
}
