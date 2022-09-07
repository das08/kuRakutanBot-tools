package models

import "fmt"

type RakutanInfo struct {
	Faculty       string `json:"faculty"`
	LectureName   string `json:"lecture_name"`
	RegisterTotal int    `json:"register_total"`
	PassedTotal   int    `json:"passed_total"`
}

func (r *RakutanInfo) Print() {
	fmt.Printf("FN: %s, LN: %s, RT: %d, PT: %d \n", r.Faculty, r.LectureName, r.RegisterTotal, r.PassedTotal)
}
