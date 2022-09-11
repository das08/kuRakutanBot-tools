package main

import (
	"github.com/das08/kuRakutanBot-tools/models"
	"github.com/gocarina/gocsv"
	"github.com/goccy/go-json"
	"os"
	"time"
)

type MongoUserData struct {
	Id struct {
		Oid string `json:"$oid"`
	} `json:"_id"`
	Uid   string `json:"uid"`
	Count struct {
		Message int `json:"message"`
		Rakutan int `json:"rakutan"`
		Onitan  int `json:"onitan"`
	} `json:"count"`
	RegisterTime int  `json:"register_time"`
	Verified     bool `json:"verified"`
	VerifiedTime int  `json:"verified_time"`
}

func saveToCSV(rakutanCSVs []models.UserDataCSV) {
	file, _ := os.OpenFile("../export/users.csv", os.O_RDWR|os.O_CREATE, os.ModePerm)
	defer file.Close()
	gocsv.MarshalFile(rakutanCSVs, file)
}

func main() {
	// load json
	file, err := os.Open("../mongodump/user_data.json")
	if err != nil {
		panic(err)
	}
	defer file.Close()
	// decode json
	decoder := json.NewDecoder(file)
	var data []MongoUserData
	err = decoder.Decode(&data)
	if err != nil {
		panic(err)
	}

	var userDataCSV []models.UserDataCSV
	for _, v := range data {
		registerDate := time.Unix(int64(v.RegisterTime), 0).Format("2006-01-02 15:04:05.000000+00")
		verifiedDate := time.Unix(int64(v.VerifiedTime), 0).Format("2006-01-02 15:04:05.000000+00")
		isVerified := "FALSE"

		if v.RegisterTime == 0 {
			registerDate = ""
		}
		if v.VerifiedTime == 0 {
			verifiedDate = ""
		}
		if v.Verified {
			isVerified = "TRUE"
		}

		userDataCSV = append(userDataCSV, models.UserDataCSV{
			UID:          v.Uid,
			IsVerified:   isVerified,
			RegisteredAt: registerDate,
			VerifiedAt:   verifiedDate,
		})
	}
	saveToCSV(userDataCSV)
}
