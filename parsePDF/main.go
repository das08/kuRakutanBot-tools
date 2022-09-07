package main

import (
	"fmt"
	"github.com/das08/pdf2text"
	"strings"
)

type RakutanInfo struct {
	Faculty       string `json:"faculty"`
	LectureName   string `json:"lecture_name"`
	RegisterTotal string `json:"register_total"`
	PassedTotal   string `json:"passed_total"`
}

func (r *RakutanInfo) Print() {
	fmt.Printf("FN: %s, LN: %s, RT: %s, PT: %s \n", r.Faculty, r.LectureName, r.RegisterTotal, r.PassedTotal)
}

func main() {
	content, err := readPdf2("pdf/2021.pdf")
	if err != nil {
		panic(err)
	}
	fmt.Println(content)
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
func getText(validator func(pdf.Text) bool, text pdf.Text, init string, dest *string) (string, bool) {
	if validator(text) {
		init = init + text.S
	} else if init != "" {
		// If the text is the last one of the sentence, trim text
		init = strings.TrimSpace(init)
		init = strings.ReplaceAll(init, "�", "")
		*dest = init
		return "", true
	}
	return init, false
}

func readPdf2(path string) (string, error) {
	f, r, err := pdf.Open(path)
	defer f.Close()
	if err != nil {
		return "", err
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

	for _, r := range rakutanInfos {
		r.Print()
	}
	return "", nil
}
