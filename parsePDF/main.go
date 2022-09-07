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

func isSameSentence(text pdf.Text, lastTextStyle pdf.Text) bool {
	return (text.Font == lastTextStyle.Font) && (text.FontSize == lastTextStyle.FontSize) && (text.X == lastTextStyle.X)
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

func getText(validator func(pdf.Text) bool, text pdf.Text, init string, dest *string) (string, bool) {
	if validator(text) {
		init = init + text.S
	} else if init != "" {
		init = strings.TrimSpace(init)
		init = strings.ReplaceAll(init, "�", "")
		*dest = init
		return "", true
	}
	return init, false
}

func readPdf2(path string) (string, error) {
	f, r, err := pdf.Open(path)
	// remember close file
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
		var lastTextStyle pdf.Text
		rakutanInfo := RakutanInfo{}
		//var facultyName string
		var _facultyName, _lectureName, _regStr, _passStr string
		var ok bool

		texts := p.Content().Text
		for _, text := range texts {
			_facultyName, _ = getText(isFacultyName, text, _facultyName, &rakutanInfo.Faculty)
			_lectureName, _ = getText(isLectureName, text, _lectureName, &rakutanInfo.LectureName)
			_regStr, _ = getText(isRegisterTotal, text, _regStr, &rakutanInfo.RegisterTotal)
			_passStr, ok = getText(isPassedTotal, text, _passStr, &rakutanInfo.PassedTotal)

			//getText(isLectureName, text, _lectureName, "lecture_name")
			//getText(isRegisterTotal, text, _regStr, "register_total")
			//getText(isPassedTotal, text, _passStr, "passed_total")
			//if isFacultyName(text) {
			//	_facultyName = _facultyName + text.S
			//	continue
			//} else if _facultyName != "" {
			//	fmt.Printf("Faculty: %s \n", _facultyName)
			//	_facultyName = ""
			//}

			if _passStr == "" && ok {
				//fmt.Printf("Rakutan: %v \n", rakutanInfo)
				rakutanInfos = append(rakutanInfos, rakutanInfo)
				rakutanInfo = RakutanInfo{}
			}

			//if isLectureName(text) {
			//	_lectureName = _lectureName + text.S
			//	continue
			//} else if _lectureName != "" {
			//	fmt.Printf("Lecture: %s \n", _lectureName)
			//	_lectureName = ""
			//}
			//
			//if isRegisterTotal(text) {
			//	_regStr = _regStr + text.S
			//	continue
			//} else if _regStr != "" {
			//	fmt.Printf("Register: %s \n", _regStr)
			//	_regStr = ""
			//}
			//
			//if isPassedTotal(text) {
			//	_passStr = _passStr + text.S
			//	continue
			//} else if _passStr != "" {
			//	fmt.Printf("Passed: %s \n", _passStr)
			//	_passStr = ""
			//}

			if isSameSentence(text, lastTextStyle) {
				lastTextStyle.S = lastTextStyle.S + text.S
			} else {
				fmt.Printf("x: %f, y: %f, content: %s \n", lastTextStyle.X, lastTextStyle.Y, lastTextStyle.S)
				lastTextStyle = text
			}
		}
		//fmt.Printf("x: %f, y: %f, content: %s \n", lastTextStyle.X, lastTextStyle.Y, lastTextStyle.S)
		//fmt.Printf("Faculty: %s \n", faculty)
	}
	for _, r := range rakutanInfos {
		r.Print()
	}
	return "", nil
}
