package main

import (
	"fmt"
	"github.com/das08/kuRakutanBot-tools/models"
	"github.com/gocarina/gocsv"
	"github.com/goccy/go-json"
	"io/ioutil"
	"strconv"
)

var YEAR = []int{2022, 2021, 2020, 2019, 2018}

const (
	BaseYear = 2022
)

func readJSON() map[int][]models.RakutanPDF {
	rakutanPDFs := map[int][]models.RakutanPDF{}
	for _, year := range YEAR {
		data, err := ioutil.ReadFile(fmt.Sprintf("../../parsePDF/export/%d.json", year))
		if err != nil {
			panic(fmt.Errorf("failed to read file: %v", err))
		}

		var rakutanPDF []models.RakutanPDF
		err = json.Unmarshal(data, &rakutanPDF)
		if err != nil {
			panic(fmt.Errorf("failed to unmarshal: %v", err))
		}
		rakutanPDFs[year] = rakutanPDF
	}

	return rakutanPDFs
}

func validate(rakutanPDFs map[int][]models.RakutanPDF, rakutanCSVs []models.RakutanCSV) int {
	// accumulate same faculty and lecture name per year
	accumulatedPDFs := map[int]map[string]models.RakutanPDF{}
	for _, year := range YEAR {
		accumulatedPDFs[year] = map[string]models.RakutanPDF{}
		for _, r := range rakutanPDFs[year] {
			key := fmt.Sprintf("%s:%s", r.FacultyName, r.LectureName)
			if old, ok := accumulatedPDFs[year][key]; ok {
				old.RegisterTotal += r.RegisterTotal
				old.PassedTotal += r.PassedTotal
				accumulatedPDFs[year][key] = old
			} else {
				accumulatedPDFs[year][key] = r
			}
		}
	}

	errorCount := 0
	for _, year := range YEAR {
		index := BaseYear - year
		for _, pdf := range accumulatedPDFs[year] {
			for _, csv := range rakutanCSVs {
				if pdf.FacultyName == csv.FacultyName && pdf.LectureName == csv.LectureName {
					if strconv.Itoa(pdf.RegisterTotal) != csv.RegisterTotalArray.ToArray()[index] {
						fmt.Printf("RegisterTotal is not matched: [%d] %s, %s, (%d, %s)\n", year, pdf.FacultyName, pdf.LectureName, pdf.RegisterTotal, csv.RegisterTotalArray.ToArray()[index])
						errorCount++
					}
					if strconv.Itoa(pdf.PassedTotal) != csv.PassedTotalArray.ToArray()[index] {
						fmt.Printf("PassedTotal is not matched:[%d] %s, %s, (%d, %s)\n", year, pdf.FacultyName, pdf.LectureName, pdf.PassedTotal, csv.PassedTotalArray.ToArray()[index])
						errorCount++
					}
				}
			}
		}
	}

	return errorCount
}

func readCSV() []models.RakutanCSV {
	var rakutanCSVs []models.RakutanCSV
	data, err := ioutil.ReadFile(fmt.Sprintf("../export/%d.csv", BaseYear))
	if err != nil {
		panic(err)
	}
	err = gocsv.UnmarshalBytes(data, &rakutanCSVs)
	if err != nil {
		panic(err)
	}
	return rakutanCSVs
}

func main() {
	rakutanPDFs := readJSON()
	rakutanCSVs := readCSV()
	errCount := validate(rakutanPDFs, rakutanCSVs)
	fmt.Printf("error count: %d\n", errCount)
}
