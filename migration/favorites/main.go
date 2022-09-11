package main

import (
	"fmt"
	"github.com/das08/kuRakutanBot-tools/models"
	"github.com/gocarina/gocsv"
	"github.com/goccy/go-json"
	"golang.org/x/text/width"
	"os"
	"strings"
	"time"
)

type MongoFavorites struct {
	Id struct {
		Oid string `json:"$oid"`
	} `json:"_id"`
	Uid         string `json:"uid"`
	Id1         int    `json:"id"`
	LectureName string `json:"lecture_name"`
}

type MongoRakutanInfo struct {
	Id struct {
		Oid string `json:"$oid"`
	} `json:"_id"`
	Id1         int    `json:"id"`
	FacultyName string `json:"faculty_name"`
	LectureName string `json:"lecture_name"`
	Detail      []struct {
		Year     int  `json:"year"`
		Accepted *int `json:"accepted"`
		Total    *int `json:"total"`
	} `json:"detail"`
	Omikuji string `json:"omikuji"`
	Url     string `json:"url"`
}

type UserFavorites struct {
	UID         string `json:"uid"`
	OldID       int    `json:"old_id"`
	NewID       int    `json:"new_id"`
	FacultyName string `json:"faculty_name"`
	LectureName string `json:"lecture_name"`
}

func loadMongoFavorites() []MongoFavorites {
	// load json
	file, err := os.Open("../mongodump/favorites.json")
	if err != nil {
		panic(err)
	}
	defer file.Close()
	// decode json
	decoder := json.NewDecoder(file)
	var data []MongoFavorites
	err = decoder.Decode(&data)
	if err != nil {
		panic(err)
	}
	return data
}

func loadMongoRakutanInfo() []MongoRakutanInfo {
	// load json
	file, err := os.Open("../mongodump/rakutan2021.json")
	if err != nil {
		panic(err)
	}
	defer file.Close()
	// decode json
	decoder := json.NewDecoder(file)
	var data []MongoRakutanInfo
	err = decoder.Decode(&data)
	if err != nil {
		panic(err)
	}
	return data
}

func loadRakutanInfoCSV() []models.RakutanCSV {
	var newRakutanInfo []models.RakutanCSV
	file, err := os.OpenFile("../../mergeJSON/export/2021.csv", os.O_RDWR|os.O_CREATE, os.ModePerm)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	if err := gocsv.UnmarshalFile(file, &newRakutanInfo); err != nil { // Load clients from file
		panic(err)
	}
	//for _, v := range newRakutanInfo {
	//	fmt.Printf("id: %d, faculty_name: %s, lecture_name: %s\n", v.ID, v.FacultyName, v.LectureName)
	//}
	return newRakutanInfo
}

func createMongoRakutanMap(rakutanInfo []MongoRakutanInfo) map[int]MongoRakutanInfo {
	rakutanMap := make(map[int]MongoRakutanInfo)
	for _, v := range rakutanInfo {
		rakutanMap[v.Id1] = v
	}
	return rakutanMap
}

func createNewRakutanMap(rakutanInfo []models.RakutanCSV) map[string]models.RakutanCSV {
	rakutanMap := make(map[string]models.RakutanCSV)
	for _, v := range rakutanInfo {
		rakutanMap[fmt.Sprintf("%s:%s", v.FacultyName, v.LectureName)] = v
	}
	return rakutanMap
}

func formatter(text string) string {
	text = strings.ReplaceAll(text, "�", "")
	text = strings.TrimSpace(text)
	text = width.Fold.String(text)
	return text
}

func saveToCSV(rakutanCSVs []models.UserFavoritesCSV) {
	file, _ := os.OpenFile("../export/favorites.csv", os.O_RDWR|os.O_CREATE, os.ModePerm)
	defer file.Close()
	gocsv.MarshalFile(rakutanCSVs, file)
}

func main() {
	mongoFavorites := loadMongoFavorites()
	mongoRakutanInfo := loadMongoRakutanInfo()
	newRakutanInfo := loadRakutanInfoCSV()
	rakutanMap := createMongoRakutanMap(mongoRakutanInfo)
	newRakutanMap := createNewRakutanMap(newRakutanInfo)
	var userFavorites []UserFavorites
	for _, v := range mongoFavorites {
		rakutanInfo, ok := rakutanMap[v.Id1]
		if ok {
			userFavorites = append(userFavorites, UserFavorites{
				UID:         v.Uid,
				OldID:       v.Id1,
				FacultyName: rakutanInfo.FacultyName,
				LectureName: rakutanInfo.LectureName,
			})
		} else {
			fmt.Printf("id: %d not found", v.Id1)
			panic("id not found")
		}
	}
	// print output
	//for _, v := range userFavorites {
	//	fmt.Printf("uid: %s, old_id: %d, faculty_name: %s, lecture_name: %s\n", v.UID, v.OldID, v.FacultyName, v.LectureName)
	//}

	// create new id
	var favoritesCSV []models.UserFavoritesCSV
	for _, v := range userFavorites {
		v.FacultyName = formatter(v.FacultyName)
		v.LectureName = formatter(v.LectureName)
		newRakutanInfo, ok := newRakutanMap[fmt.Sprintf("%s:%s", v.FacultyName, v.LectureName)]

		createdAt := time.Now().Format("2006-01-02 15:04:05.000000+00")

		if ok {
			favoritesCSV = append(favoritesCSV, models.UserFavoritesCSV{
				UID:       v.UID,
				ID:        newRakutanInfo.ID,
				CreatedAt: createdAt,
			})
		} else {
			fmt.Errorf("id: %d, fac: %s, lec: %s\n", v.OldID, v.FacultyName, v.LectureName)
		}
	}
	saveToCSV(favoritesCSV)
}
