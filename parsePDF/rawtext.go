package main

import (
	"fmt"
	"github.com/das08/pdf2text"
)

func main() {
	content, err := readPdfRaw("pdf/2021.pdf")
	if err != nil {
		panic(err)
	}
	fmt.Println(content)
	return
}

func isSameSentenceRaw(text pdf.Text, lastTextStyle pdf.Text) bool {
	return (text.Font == lastTextStyle.Font) && (text.FontSize == lastTextStyle.FontSize) && (text.X == lastTextStyle.X)
}

func readPdfRaw(path string) (string, error) {
	f, r, err := pdf.Open(path)
	// remember close file
	defer f.Close()
	if err != nil {
		return "", err
	}
	totalPage := 94

	for pageIndex := totalPage; pageIndex <= totalPage; pageIndex++ {
		p := r.Page(pageIndex)
		if p.V.IsNull() {
			continue
		}
		var lastTextStyle pdf.Text

		texts := p.Content().Text
		for _, text := range texts {
			if isSameSentenceRaw(text, lastTextStyle) {
				lastTextStyle.S = lastTextStyle.S + text.S
			} else {
				fmt.Printf("x: %f, y: %f, content: %s \n", lastTextStyle.X, lastTextStyle.Y, lastTextStyle.S)
				lastTextStyle = text
			}
		}
		fmt.Printf("x: %f, y: %f, content: %s \n", lastTextStyle.X, lastTextStyle.Y, lastTextStyle.S)
	}
	return "", nil
}
