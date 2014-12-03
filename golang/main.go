package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"
)

// Book e book
type Book struct {
	Name     string
	Acronym  string
	Chapters int
}

var folder = "/Users/thiago/go/src/github.com/thslopes/Moody"
var prefix = "<p class=\"calibre2\"><b class=\"calibre1\">"
var patternChapter = "<p class=\"calibre2\"><b class=\"calibre1\" id=\"%s%d\">%s"
var patternBook = "<p class=\"calibre2\"><b class=\"calibre1\" id=\"%s\">%s"
var patternIndex = "<a href=\"#%s\">%s</a>"
var patternHeaderItem = "<p class=\"calibre2\"><a href=\"#%s%s\"><b class=\"calibre1\">%s</b></a></p>"
var patternChapterIndex = "<a href=\"#%s%d\">%d</a>"
var books = initBooks()
var book Book
var nextBook Book
var chapter = 1
var bookIdx = -1

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func write(f *os.File, text string) {
	_, err := io.WriteString(f, text)
	check(err)
}

func getFiles() []os.FileInfo {
	files, err := ioutil.ReadDir(folder)
	check(err)
	return files
}

func getFileContent(fileName string) []string {
	dat, err := ioutil.ReadFile(fmt.Sprintf("%s/%s", folder, fileName))
	check(err)
	lines := strings.Split(string(dat), "\n")
	return lines
}

func (book *Book) printHeader() string {
	header := fmt.Sprintf(patternHeaderItem, book.Acronym, "intro", "INTRODUÇÃO")
	header += fmt.Sprintf(patternHeaderItem, book.Acronym, "outline", "ESBOÇO")
	header += fmt.Sprint("<p class=\"calibre2\"><b class=\"calibre1\">")
	for i := 1; i <= book.Chapters; i++ {
		header += fmt.Sprintf(patternChapterIndex, book.Acronym, i, i)
	}
	header += fmt.Sprint("</b></p>")

	//fmt.Sprint(header)

	return header
}

func printChapter(line string) string {
	trimLine := strings.TrimLeft(line, " ")
	io := strings.Index(trimLine, fmt.Sprintf("%s%s %d", prefix, book.Name, chapter))
	io2 := strings.Index(trimLine, fmt.Sprintf("%s%s cap. %d", prefix, book.Name, chapter))
	if io == 0 || io2 == 0 {
		line = fmt.Sprintf(patternChapter, book.Acronym, chapter, strings.Replace(line, prefix, "", 1))
		chapter++
		//		fmt.Sprint(line)
	}
	return line
}

func printIndex(printHeaders bool) string {
	index := fmt.Sprint("<b class=\"calibre1\">COMENTÁRIO  BÍBLICO  MOODY </b>")
	index += fmt.Sprint("<p class=\"calibre2\"><b class=\"calibre1\">Moody Bible Institute of Chicago </b></p>")
	index += fmt.Sprint("<p class=\"calibre2\"><b class=\"calibre1\">Clique num livro bíblico para o comentário</b></p>")
	index += fmt.Sprint("<p class=\"calibre2\"><b class=\"calibre1\">ANTIGO TESTAMENTO</b></p>")
	index += fmt.Sprint("<p class=\"calibre2\"><b class=\"calibre1\" >")
	for i := 0; i < 39; i++ {
		book := books[i]
		index += fmt.Sprintf(patternIndex, book.Acronym, book.Name)
		if printHeaders {
			index += book.printHeader()
		}
	}
	index += fmt.Sprint("</b></p>")
	index += fmt.Sprint("<p class=\"calibre2\"><b class=\"calibre1\">NOVO TESTAMENTO</b></p>")
	index += fmt.Sprint("<p class=\"calibre2\"><b class=\"calibre1\">")
	for i := 39; i < 66; i++ {
		book := books[i]
		index += fmt.Sprintf(patternIndex, book.Acronym, book.Name)
		if printHeaders {
			index += book.printHeader()
		}
	}
	index += fmt.Sprint("</b></p>")

	//fmt.Print(index)
	return index
}

func printIntroduction(line string, bookHeader bool) (string, bool) {
	trimLine := strings.TrimLeft(line, " ")
	io := strings.Index(trimLine, fmt.Sprintf("%sINTRODUÇÃO </b></p>", prefix))
	if io == 0 {
		line = fmt.Sprintf(patternBook, book.Acronym+"intro", strings.Replace(line, prefix, "", 1))
		bookHeader = false
		//fmt.Sprint(line)
	}
	return line, bookHeader
}

func printOutline(line string) string {
	trimLine := strings.TrimLeft(line, " ")
	io := strings.Index(trimLine, fmt.Sprintf("%sESBOÇO </b></p>", prefix))
	if io == 0 {
		line = fmt.Sprintf(patternBook, book.Acronym+"outline", strings.Replace(line, prefix, "", 1))
		//fmt.Sprint(line)
	}
	return line
}

func printBook(line string, bookHeader bool) (string, bool, bool) {
	title := false
	trimLine := strings.TrimLeft(line, " ")
	if bookIdx < 65 {
		nextBook = books[bookIdx+1]
		nextBook.Name = strings.ToUpper(nextBook.Name)
		if nextBook.Name == "CANTARES" {
			nextBook.Name = "CANTARES DE SALOMÃO"
		}
	}
	io := strings.Index(trimLine, fmt.Sprintf("%s%s </b></p>", prefix, nextBook.Name))
	io2 := strings.Index(trimLine, fmt.Sprintf("%s  </b> <b class=\"calibre1\">%s </b></p>", prefix, nextBook.Name))
	if io == 0 || io2 == 0 {
		//		fmt.Sprint("{\"%s\",\"%s\",%d},", book.Name, book.Acronym, chapter-1)
		bookIdx++
		book = books[bookIdx]
		if book.Name == "Salmos" {
			book.Name = "Salmo"
		}
		chapter = 1
		line = fmt.Sprintf(patternBook, book.Acronym, strings.Replace(line, prefix, "", 1))
		//fmt.Sprint(line)
		line += book.printHeader()
		bookHeader = true
		title = true
	}
	return line, bookHeader, title
}

func isBody(inBody bool, isHeader bool, line string) (bool, bool) {
	if !isHeader {
		inBody = true
	}
	trimLine := strings.TrimLeft(line, " ")
	io := strings.Index(trimLine, "<body")
	if io == 0 {
		isHeader = false
	}

	return inBody, isHeader
}

func printFullIndex() {
	f, err := os.OpenFile(folder+"/new/header.html", os.O_CREATE|os.O_TRUNC|os.O_RDWR, 0666)
	defer f.Close()
	check(err)
	write(f, printIndex(false) + "\n")
	write(f, printIndex(true))

}

func main() {
	printFullIndex()
	book = books[0]
	nextBook = books[0]
	bookHeader := false
	title := false
	for _, file := range getFiles() {
		if strings.Index(file.Name(), "index") == 0 {
			f, err := os.OpenFile(fmt.Sprintf("%s/new/%s", folder, file.Name()), os.O_CREATE|os.O_TRUNC|os.O_RDWR, 0666)
			defer f.Close()
			check(err)
			for _, line := range getFileContent(file.Name()) {
				line = printChapter(line)
				line, bookHeader, title = printBook(line, bookHeader)
				line = printOutline(line)
				line, bookHeader = printIntroduction(line, bookHeader)
				if !bookHeader || title {
					write(f, line + "\n")
				}
				title = false
			}
		}
	}
	//	fmt.Sprint("{\"%s\",\"%s\",%d}", book.Name, book.Acronym, chapter-1)
	//write(f, "</body></html>")
}

func initBooks() [66]Book {
	books := [66]Book{{"Gênesis", "gn", 50}, {"Êxodo", "ex", 40}, {"Levítico", "lv", 27}, {"Números", "nm", 36}, {"Deuteronômio", "dt", 34}, {"Josué", "js", 24}, {"Juízes", "jz", 21}, {"Rute", "rt", 4}, {"1 Samuel", "1sm", 31}, {"2 Samuel", "2sm", 24}, {"1 Reis", "1re", 22}, {"2 Reis", "2re", 25}, {"1 Crônicas", "1cr", 29}, {"2 Crônicas", "2cr", 36}, {"Esdras", "ed", 10}, {"Neemias", "ne", 13}, {"Ester", "et", 10}, {"Jó", "jo", 42}, {"Salmos", "sl", 150}, {"Provérbios", "pv", 31}, {"Eclesiastes", "ec", 12}, {"Cantares", "ct", 8}, {"Isaías", "is", 66}, {"Jeremias", "jr", 52}, {"Lamentações", "lm", 5}, {"Ezequiel", "ez", 48}, {"Daniel", "dn", 12}, {"Oséias", "os", 14}, {"Joel", "jl", 3}, {"Amós", "am", 9}, {"Obadias", "ob", 1}, {"Jonas", "jn", 4}, {"Miquéias", "mq", 7}, {"Naum", "na", 3}, {"Habacuque", "hc", 3}, {"Sofonias", "sf", 3}, {"Ageu", "ag", 2}, {"Zacarias", "zc", 14}, {"Malaquias", "mq", 4}, {"Mateus", "mt", 28}, {"Marcos", "mc", 16}, {"Lucas", "lc", 24}, {"João", "ja", 21}, {"Atos", "at", 28}, {"Romanos", "rm", 16}, {"1 Coríntios", "1co", 16}, {"2 Coríntios", "2co", 13}, {"Gálatas", "gl", 6}, {"Efésios", "ef", 6}, {"Filipenses", "fp", 4}, {"Colossenses", "cl", 4}, {"1 Tessalonicenses", "1ts", 5}, {"2 Tessalonicenses", "2ts", 3}, {"1 Timóteo", "1tm", 6}, {"2 Timóteo", "2tm", 4}, {"Tito", "tt", 3}, {"Filemom", "fl", 1}, {"Hebreus", "hb", 13}, {"Tiago", "tg", 5}, {"1 Pedro", "1pe", 5}, {"2 Pedro", "2pe", 3}, {"1 João", "1jo", 5}, {"2 João", "2jo", 1}, {"3 João", "3jo", 1}, {"Judas", "jd", 1}, {"Apocalipse", "ap", 22}}
	return books
}
