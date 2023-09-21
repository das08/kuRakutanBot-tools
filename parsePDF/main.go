package main

import (
	"fmt"
	"github.com/das08/kuRakutanBot-tools/models"
	"github.com/das08/pdf2text"
	"github.com/goccy/go-json"
	"golang.org/x/text/width"
	"io/ioutil"
	"strconv"
	"strings"
	"time"
)

var YEAR = []int{2017, 2018, 2019, 2020, 2021}

func main() {
	start := time.Now()
	for _, year := range YEAR {
		content, err := readPdf2(fmt.Sprintf("pdf/%d.pdf", year))
		if err != nil {
			panic(err)
		}
		//for _, r := range content {
		//	r.Print()
		//}
		fmt.Println("processed year: ", year)
		file, _ := json.MarshalIndent(content, "", " ")
		_ = ioutil.WriteFile(fmt.Sprintf("export/%d.json", year), file, 0644)
	}

	end := time.Now()
	fmt.Printf("Process Ended in: %fs\n", (end.Sub(start)).Seconds())
	return
}

func isFacultyName(text pdf.Text) bool {
	return text.X >= 19.0 && text.X <= 35.0
}

func isLectureName(text pdf.Text) bool {
	return text.X >= 170.0 && text.X <= 370.0
}

func isRegisterTotal(text pdf.Text) bool {
	return text.X >= 390.0 && text.X <= 415.0
}

func isPassedTotal(text pdf.Text) bool {
	return text.X >= 470.0 && text.X <= 490.0
}

// getText returns the appended text and a boolean value indicating whether the text is the last one of the sentence.
// If the text is the last one of the sentence, the text is assigned to `dest` and the return text is set to "".
func getText[T *string | *int](validator func(pdf.Text) bool, text pdf.Text, init string, dest T) (string, bool) {
	if validator(text) {
		init = init + text.S
	} else if init != "" {
		// If the text is the last one of the sentence, format text
		formatted := formatter(init)

		switch p := any(dest).(type) {
		case *string:
			*p = formatted
		case *int:
			*p, _ = strconv.Atoi(formatted)
		}
		return "", true
	}
	return init, false
}

// formatter formats the text
func formatter(text string) string {
	text = strings.ReplaceAll(text, "�", "")
	text = strings.TrimSpace(text)
	text = width.Fold.String(text)
	return text
}

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

type ValidateError struct {
	Err string
}

func (e *ValidateError) Error() string {
	return e.Err
}

func validator(rakutanPDF models.RakutanPDF) (bool, error) {
	var validFacultyName = []string{"文学部", "教育学部", "法学部", "経済学部", "理学部", "医学部",
		"医学部(人間健康科学科)", "薬学部", "工学部", "農学部", "総合人間学部", "文学研究科", "教育学研究科",
		"法学研究科", "経済学研究科", "理学研究科", "医学研究科", "医学研究科(人間健康科学系専攻", "薬学研究科",
		"工学研究科", "農学研究科", "人間・環境学研究科", "エネルギー科学研究科", "アジア・アフリカ地域研究研究科",
		"情報学研究科", "生命科学研究科", "地球環境学舎", "公共政策教育部", "経営管理教育部", "法学研究科(法科大学院)",
		"総合生存学館", "国際高等教育院",
	}
	if rakutanPDF.FacultyName == "" || rakutanPDF.LectureName == "" {
		return false, &ValidateError{Err: "FacultyName or LectureName is empty"}
	}
	if len(rakutanPDF.FacultyName) > 100 {
		return false, &ValidateError{Err: "FacultyName is too long"}
	}
	if rakutanPDF.PassedTotal > rakutanPDF.RegisterTotal {
		return false, &ValidateError{Err: "PassedTotal is bigger than RegisterTotal"}
	}
	if rakutanPDF.RegisterTotal == 0 {
		return false, &ValidateError{Err: "RegisterTotal is zero"}
	}

	if ok := contains(validFacultyName, rakutanPDF.FacultyName); !ok {
		//fmt.Println("Invalid FacultyName:", len(rakutanPDF.FacultyName))
		fmt.Printf("Invalid FacultyName: %s %d\n", rakutanPDF.FacultyName, len(rakutanPDF.FacultyName))
		return false, &ValidateError{Err: fmt.Sprintf("Invalid FacultyName: %s", rakutanPDF.FacultyName)}
	}
	return true, nil
}

func readPdf2(path string) ([]models.RakutanPDF, error) {
	f, r, err := pdf.Open(path)
	defer f.Close()
	if err != nil {
		return []models.RakutanPDF{}, err
	}
	totalPage := r.NumPage()

	var rakutanInfos []models.RakutanPDF

	for pageIndex := 1; pageIndex <= totalPage; pageIndex++ {
		p := r.Page(pageIndex)
		if p.V.IsNull() {
			continue
		}

		rakutanPDF := models.RakutanPDF{}
		var _facultyName, _lectureName, _regStr, _passStr string
		var isEnd bool

		texts := p.Content().Text
		for _, text := range texts {
			_facultyName, _ = getText(isFacultyName, text, _facultyName, &rakutanPDF.FacultyName)
			_lectureName, _ = getText(isLectureName, text, _lectureName, &rakutanPDF.LectureName)
			_regStr, _ = getText(isRegisterTotal, text, _regStr, &rakutanPDF.RegisterTotal)
			_passStr, isEnd = getText(isPassedTotal, text, _passStr, &rakutanPDF.PassedTotal)

			if isEnd {
				// Append to list if rakutanInfo is all set and validated
				if ok, err := validator(rakutanPDF); ok {
					rakutanInfos = append(rakutanInfos, rakutanPDF)
					rakutanPDF = models.RakutanPDF{}
				} else if err != nil {
					fmt.Printf("Error: %s, rakutanInfo: %v\n", err.Error(), rakutanPDF)
				}
			}
		}
	}

	fmt.Println("Total Entries: ", len(rakutanInfos))
	return rakutanInfos, nil
}

//Total RakutanInfo:  8824
//processed year:  2018
//Total RakutanInfo:  7248
//processed year:  2019
//Total RakutanInfo:  8733
//processed year:  2020
//Total RakutanInfo:  8704
//processed year:  2021
//Process Ended in: 19.644771s
