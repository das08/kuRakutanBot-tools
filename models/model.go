package models

import (
	"bytes"
	"fmt"
	"github.com/goccy/go-json"
	"strconv"
	"strings"
)

type RakutanPDF struct {
	FacultyName   string `json:"faculty_name"`
	LectureName   string `json:"lecture_name"`
	RegisterTotal int    `json:"register_total"`
	PassedTotal   int    `json:"passed_total"`
}

func (r *RakutanPDF) Print() {
	fmt.Printf("FN: %s, LN: %s, RT: %d, PT: %d \n", r.FacultyName, r.LectureName, r.RegisterTotal, r.PassedTotal)
}

type RakutanEntry struct {
	ID            int       `json:"id"`
	FacultyName   string    `json:"faculty_name"`
	LectureName   string    `json:"lecture_name"`
	RegisterTotal []NullInt `json:"register_total"`
	PassedTotal   []NullInt `json:"passed_total"`
	KakomonURL    string    `json:"kakomon_url"`
}

func Join(a []NullInt) string {
	var b []string
	for _, v := range a {
		if v.Valid {
			b = append(b, strconv.Itoa(v.Int))
		} else {
			b = append(b, "NULL")
		}
	}
	return "{" + strings.Join(b, ",") + "}"
}

func (r *RakutanEntry) ToRakutanCSV() RakutanCSV {
	rakutanCSV := RakutanCSV{
		ID:                 r.ID,
		FacultyName:        r.FacultyName,
		LectureName:        r.LectureName,
		RegisterTotalArray: Join(r.RegisterTotal),
		PassedTotalArray:   Join(r.PassedTotal),
		KakomonURL:         r.KakomonURL,
	}
	return rakutanCSV
}

type RakutanCSV struct {
	ID                 int    `csv:"id"`
	FacultyName        string `csv:"faculty_name"`
	LectureName        string `csv:"lecture_name"`
	RegisterTotalArray string `csv:"register"`
	PassedTotalArray   string `csv:"passed"`
	KakomonURL         string `csv:"-"`
}

type UserDataCSV struct {
	UID          string `csv:"uid"`
	IsVerified   string `csv:"is_verified"`
	RegisteredAt string `csv:"registered_at"`
	VerifiedAt   string `csv:"verified_at"`
}

type UserFavoritesCSV struct {
	UID       string `csv:"uid"`
	ID        int    `csv:"id"`
	CreatedAt string `csv:"created_at"`
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
