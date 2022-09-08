package models

import (
	"bytes"
	"fmt"
	"github.com/goccy/go-json"
)

type RakutanPDF struct {
	Faculty       string `json:"faculty"`
	LectureName   string `json:"lecture_name"`
	RegisterTotal int    `json:"register_total"`
	PassedTotal   int    `json:"passed_total"`
}

func (r *RakutanPDF) Print() {
	fmt.Printf("FN: %s, LN: %s, RT: %d, PT: %d \n", r.Faculty, r.LectureName, r.RegisterTotal, r.PassedTotal)
}

type RakutanEntry struct {
	ID            int       `json:"id"`
	FacultyName   string    `json:"faculty_name"`
	LectureName   string    `json:"lecture_name"`
	RegisterTotal []NullInt `json:"register_total"`
	PassedTotal   []NullInt `json:"passed_total"`
	KakomonURL    string    `json:"kakomon_url"`
}

var nullLiteral = []byte("null")

type NullInt struct {
	Int   int
	Valid bool // Valid is true if Segments is not NULL
}

func (s *NullInt) UnmarshalJSON(b []byte) error {
	if bytes.Equal(b, nullLiteral) {
		return nil
	}

	err := json.Unmarshal(b, &s.Int)
	if err == nil {
		s.Valid = true
		return nil
	}

	return err
}

func (s NullInt) MarshalJSON() ([]byte, error) {
	if s.Valid {
		return json.Marshal(s.Int)
	} else {
		return nullLiteral, nil
	}
}
