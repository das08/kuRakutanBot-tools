package main

import (
	"fmt"
	"github.com/das08/kuRakutanBot-migration/models"
	"github.com/goccy/go-json"
	"io/ioutil"
)

var YEAR = []int{2021}

const (
	BASE_YEAR    = 2021
	INITIAL_ID   = 100001
	ID_INCREMENT = 100000
)

func readJSON() map[int][]models.RakutanPDF {
	rakutanPDFs := map[int][]models.RakutanPDF{}
	for _, year := range YEAR {
		data, err := ioutil.ReadFile(fmt.Sprintf("../parsePDF/export/%d.json", year))
		//fmt.Printf(string(data))
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

func main() {
	//var rakutanEntries []models.RakutanEntry
	rakutanPDFs := readJSON()
	rakutanPDFMap := make(map[string]models.RakutanPDF)
	for year, rakutanPDF := range rakutanPDFs {
		for _, r := range rakutanPDF {
			key := fmt.Sprintf("%d:%s:%s", year, r.Faculty, r.LectureName)

			// If the key is already in the map, accumulate the register and passed total
			if old, ok := rakutanPDFMap[key]; ok {
				r = models.RakutanPDF{
					Faculty:       r.Faculty,
					LectureName:   r.LectureName,
					RegisterTotal: old.RegisterTotal + r.RegisterTotal,
					PassedTotal:   old.PassedTotal + r.PassedTotal,
				}
			}
			rakutanPDFMap[key] = r
		}
	}
	fmt.Println(len(rakutanPDFMap))
	fmt.Println(rakutanPDFMap["2021:情報学研究科:研究論文"])
}
