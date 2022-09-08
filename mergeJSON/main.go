package main

import (
	"fmt"
	"github.com/das08/kuRakutanBot-tools/models"
	"github.com/goccy/go-json"
	"io/ioutil"
)

// YEAR MUST BE IN DESCENDING ORDER
var YEAR = []int{2021, 2020, 2019, 2018}

var id = InitialId

const (
	BaseYear    = 2021
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
		Omikuji:       models.Normal,
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

func getLatestTotalCount(entry *models.RakutanEntry) (int, int) {
	var latestRegisterCount, latestPassedCount, i, j int
	var register, passed models.NullInt
	for i, register = range entry.RegisterTotal {
		if register.Valid {
			latestRegisterCount = register.Int
			break
		}
	}
	for j, passed = range entry.PassedTotal {
		if passed.Valid {
			latestPassedCount = passed.Int
			break
		}
	}
	if i != j {
		panic("i and j should be the same")
	}
	return latestRegisterCount, latestPassedCount
}

func main() {
	//var rakutanEntries []models.RakutanEntry
	rakutanPDFs := readJSON()

	// Merge rakutanPDFs into rakutanEntries
	rakutanEntryMap := make(map[string]*models.RakutanEntry)
	for _, year := range YEAR {
		id = InitialId + IdIncrement*(BaseYear-year)

		for _, r := range rakutanPDFs[year] {
			// Key is a combination of faculty name and lecture name since there are multiple entries
			// for the same lecture with different faculty name
			key := fmt.Sprintf("%s:%s", r.Faculty, r.LectureName)

			var entry *models.RakutanEntry
			// If the key is already in the map, accumulate the register and passed total
			if old, ok := rakutanEntryMap[key]; ok {
				entry = old
				updateEntryTotal(entry, year, r)
			} else {
				entry = initialize()
				entry.FacultyName = r.Faculty
				entry.LectureName = r.LectureName
				entry.RegisterTotal[BaseYear-year] = models.NullInt{Int: r.RegisterTotal, Valid: true}
				entry.PassedTotal[BaseYear-year] = models.NullInt{Int: r.PassedTotal, Valid: true}
				id += 1
			}
			rakutanEntryMap[key] = entry
		}
	}

	// Judge omikuji type
	for key, entry := range rakutanEntryMap {
		latestRegisterCount, latestPassedCount := getLatestTotalCount(entry)
		//fmt.Printf("%s: %d, %d\n", entry.LectureName, latestRegisterCount, latestPassedCount)
		if float64(latestRegisterCount)*0.76 <= float64(latestPassedCount) &&
			latestRegisterCount > 12 && entry.FacultyName == "国際高等教育院" {
			fmt.Printf("%s: %d, %d\n", entry.LectureName, latestRegisterCount, latestPassedCount)
			rakutanEntryMap[key].Omikuji = models.Rakutan
		}
		if float64(latestRegisterCount)*0.31 >= float64(latestPassedCount) && latestRegisterCount > 4 {
			rakutanEntryMap[key].Omikuji = models.Onitan
		}
	}

	fmt.Println(len(rakutanEntryMap))
	fmt.Println(*rakutanEntryMap["国際高等教育院:線形代数学A"])
	//for _, entry := range rakutanEntryMap {
	//	if entry.ID > 40000 {
	//		fmt.Println(entry)
	//	}
	//}

}
