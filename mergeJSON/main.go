package main

import (
	"fmt"
	"github.com/das08/kuRakutanBot-tools/models"
	"github.com/gocarina/gocsv"
	"github.com/goccy/go-json"
	"io/ioutil"
	"os"
	"sort"
)

// YEAR MUST BE IN DESCENDING ORDER
var YEAR = []int{2022, 2021, 2020, 2019, 2018}

var id = InitialId

const (
	BaseYear    = 2022
	InitialId   = 10001
	IdIncrement = 10000
)

func readJSON() map[int][]models.RakutanPDF {
	rakutanPDFs := map[int][]models.RakutanPDF{}
	for _, year := range YEAR {
		data, err := ioutil.ReadFile(fmt.Sprintf("../parsePDF/export/%d.json", year))
		if err != nil {
			fmt.Errorf("failed to read file: %v", err)
		}

		var rakutanPDF []models.RakutanPDF
		err = json.Unmarshal(data, &rakutanPDF)
		if err != nil {
			fmt.Errorf("failed to unmarshal: %v", err)
		}
		rakutanPDFs[year] = rakutanPDF
	}
	return rakutanPDFs
}

func initialize() *models.RakutanEntry {
	return &models.RakutanEntry{
		ID:            id,
		FacultyName:   "",
		LectureName:   "",
		RegisterTotal: make([]models.NullInt, len(YEAR)),
		PassedTotal:   make([]models.NullInt, len(YEAR)),
		KakomonURL:    "",
	}
}

func updateEntryTotal(entry *models.RakutanEntry, year int, r models.RakutanPDF) {
	if entry.RegisterTotal[BaseYear-year].Valid {
		entry.RegisterTotal[BaseYear-year].Int += r.RegisterTotal
	} else {
		entry.RegisterTotal[BaseYear-year] = models.NullInt{Int: r.RegisterTotal, Valid: true}
	}
	if entry.PassedTotal[BaseYear-year].Valid {
		entry.PassedTotal[BaseYear-year].Int += r.PassedTotal
	} else {
		entry.PassedTotal[BaseYear-year] = models.NullInt{Int: r.PassedTotal, Valid: true}
	}
}

func saveToCSV(rakutanCSVs []models.RakutanCSV) {
	file, _ := os.OpenFile(fmt.Sprintf("export/%d.csv", BaseYear), os.O_RDWR|os.O_CREATE, os.ModePerm)
	defer file.Close()
	gocsv.MarshalFile(rakutanCSVs, file)
}

func main() {
	var rakutanCSVs []models.RakutanCSV
	rakutanPDFs := readJSON()

	// Merge rakutanPDFs into rakutanEntries
	rakutanEntryMap := make(map[string]*models.RakutanEntry)
	for _, year := range YEAR {
		id = InitialId + IdIncrement*(BaseYear-year)

		for _, r := range rakutanPDFs[year] {
			// Key is a combination of faculty name and lecture name since there are multiple entries
			// for the same lecture with different faculty name
			key := fmt.Sprintf("%s:%s", r.FacultyName, r.LectureName)

			var entry *models.RakutanEntry
			// If the key is already in the map, accumulate the register and passed total
			if old, ok := rakutanEntryMap[key]; ok {
				entry = old
				updateEntryTotal(entry, year, r)
			} else {
				entry = initialize()
				entry.FacultyName = r.FacultyName
				entry.LectureName = r.LectureName
				entry.RegisterTotal[BaseYear-year] = models.NullInt{Int: r.RegisterTotal, Valid: true}
				entry.PassedTotal[BaseYear-year] = models.NullInt{Int: r.PassedTotal, Valid: true}
				id += 1
			}
			rakutanEntryMap[key] = entry
		}
	}

	fmt.Println(len(rakutanEntryMap))

	for _, entry := range rakutanEntryMap {
		rakutanCSVs = append(rakutanCSVs, entry.ToRakutanCSV())
	}
	// sort rakutanEntries by ID
	sort.Slice(rakutanCSVs, func(i, j int) bool {
		return rakutanCSVs[i].ID < rakutanCSVs[j].ID
	})

	saveToCSV(rakutanCSVs)
	//fmt.Println(*rakutanEntryMap["国際高等教育院:線形代数学A"])
	//for _, entry := range rakutanEntryMap {
	//	if entry.ID > 50000 {
	//		fmt.Println(entry)
	//	}
	//}
}
