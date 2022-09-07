package main

import (
	"fmt"
	"github.com/das08/pdf2text"
	"github.com/goccy/go-json"
	"golang.org/x/text/width"
	"io/ioutil"
	"strconv"
	"strings"
)

const (
	YEAR = 2021
)

type RakutanInfo struct {
	Faculty       string `json:"faculty"`
	LectureName   string `json:"lecture_name"`
	RegisterTotal int    `json:"register_total"`
	PassedTotal   int    `json:"passed_total"`
}

func (r *RakutanInfo) Print() {
	fmt.Printf("FN: %s, LN: %s, RT: %d, PT: %d \n", r.Faculty, r.LectureName, r.RegisterTotal, r.PassedTotal)
}

func main() {
	content, err := readPdf2(fmt.Sprintf("pdf/%d.pdf", YEAR))
	if err != nil {
		panic(err)
	}
	for _, r := range content {
		r.Print()
	}
	file, _ := json.MarshalIndent(content, "", " ")
	_ = ioutil.WriteFile(fmt.Sprintf("export/%d.json", YEAR), file, 0644)
	return
}

func isFacultyName(text pdf.Text) bool {
	return text.X >= 19.500 && text.X <= 20.0
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

func readPdf2(path string) ([]RakutanInfo, error) {
	f, r, err := pdf.Open(path)
	defer f.Close()
	if err != nil {
		return []RakutanInfo{}, err
	}
	totalPage := 150

	var rakutanInfos []RakutanInfo

	for pageIndex := totalPage; pageIndex <= totalPage; pageIndex++ {
		p := r.Page(pageIndex)
		if p.V.IsNull() {
			continue
		}

		rakutanInfo := RakutanInfo{}
		var _facultyName, _lectureName, _regStr, _passStr string
		var ok bool

		texts := p.Content().Text
		for _, text := range texts {
			_facultyName, _ = getText(isFacultyName, text, _facultyName, &rakutanInfo.Faculty)
			_lectureName, _ = getText(isLectureName, text, _lectureName, &rakutanInfo.LectureName)
			_regStr, _ = getText(isRegisterTotal, text, _regStr, &rakutanInfo.RegisterTotal)
			_passStr, ok = getText(isPassedTotal, text, _passStr, &rakutanInfo.PassedTotal)

			// If the text is the last one of the sentence, append rakutanInfo to rakutanInfos
			// and reset rakutanInfo
			if _passStr == "" && ok {
				rakutanInfos = append(rakutanInfos, rakutanInfo)
				rakutanInfo = RakutanInfo{}
			}
		}
	}

	return rakutanInfos, nil
}
