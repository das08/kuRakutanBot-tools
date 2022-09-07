package main

import (
	"fmt"
	"github.com/das08/kuRakutanBot-migration/models"
	"github.com/goccy/go-json"
)

func main() {
	a := models.RakutanEntry{
		FacultyName: "test",
	}
	b, _ := json.Marshal(a)
	fmt.Println(string(b))
}
