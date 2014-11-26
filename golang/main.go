package main

import (
	"fmt"
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

var folder = "/home/thiago/go/src/github.com/thslopes/Moody"
var prefix = "<p class=\"calibre2\"><b class=\"calibre1\">"
var patternChapter = "<p class=\"calibre2\"><b class=\"calibre1\" id=\"%s%d\">%s\n"
var patternBook = "<p class=\"calibre2\"><b class=\"calibre1\" id=\"%s\">%s\n"
var patternIndex = "<a href=\"#%s\">%s</a>\n"
var patternHeaderItem = "<p class=\"calibre2\"><a href=\"%s%s\"><b class=\"calibre1\">%s</b></a></p>\n"
var patternChapterIndex = "<a href=\"#%s%d\">%d</a>\n"
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

func (book *Book) printHeader() {
	fmt.Printf(patternHeaderItem, book.Acronym, "intro", "INTRODUÇÃO")
	fmt.Printf(patternHeaderItem, book.Acronym, "outline", "ESBOÇO")
	fmt.Print("<p class=\"calibre2\"><b class=\"calibre1\">\n")
	for i := 1; i <= book.Chapters; i++ {
		fmt.Printf(patternChapterIndex, book.Acronym, i, i)
	}
	fmt.Print("</b></p>\n")
}

func printChapter(line string) string {
	trimLine := strings.TrimLeft(line, " ")
	io := strings.Index(trimLine, fmt.Sprintf("%s%s %d", prefix, book.Name, chapter))
	io2 := strings.Index(trimLine, fmt.Sprintf("%s%s cap. %d", prefix, book.Name, chapter))
	if io == 0 || io2 == 0 {
		line = fmt.Sprintf(patternChapter, book.Acronym, chapter, strings.Replace(line, prefix, "", 1))
		chapter++
		fmt.Printf(line)
	}
	return line
}

func printIndex() {
	fmt.Print("<b class=\"calibre1\">COMENTÁRIO  BÍBLICO  MOODY </b>\n")
	fmt.Print("<p class=\"calibre2\"><b class=\"calibre1\">Moody Bible Institute of Chicago </b></p>\n")
	fmt.Print("<p class=\"calibre2\"><b class=\"calibre1\">Clique num livro bíblico para o comentário</b></p>\n")
	fmt.Print("<p class=\"calibre2\"><b class=\"calibre1\">ANTIGO TESTAMENTO</b></p>\n")
	fmt.Print("<p class=\"calibre2\"><b class=\"calibre1\" >\n")
	for i := 0; i < 39; i++ {
		fmt.Printf(patternIndex, books[i].Acronym, books[i].Name)
	}
	fmt.Print("</b></p>\n")
	fmt.Print("<p class=\"calibre2\"><b class=\"calibre1\">NOVO TESTAMENTO</b></p>\n")
	fmt.Print("<p class=\"calibre2\"><b class=\"calibre1\">\n")
	for i := 39; i < 66; i++ {
		fmt.Printf(patternIndex, books[i].Acronym, books[i].Name)
	}
	fmt.Print("</b></p>\n")
}

func printIntroduction(line string) string {
	trimLine := strings.TrimLeft(line, " ")
	io := strings.Index(trimLine, fmt.Sprintf("%sINTRODUÇÃO </b></p>", prefix))
	if io == 0 {
		line = fmt.Sprintf(patternBook, book.Acronym+"intro", strings.Replace(line, prefix, "", 1))
		fmt.Printf(line)
	}
	return line
}

func printOutline(line string) string {
	trimLine := strings.TrimLeft(line, " ")
	io := strings.Index(trimLine, fmt.Sprintf("%sESBOÇO </b></p>", prefix))
	if io == 0 {
		line = fmt.Sprintf(patternBook, book.Acronym+"outline", strings.Replace(line, prefix, "", 1))
		fmt.Printf(line)
	}
	return line
}

func printBook(line string) string {
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
		//		fmt.Printf("{\"%s\",\"%s\",%d},", book.Name, book.Acronym, chapter-1)
		bookIdx++
		book = books[bookIdx]
		if book.Name == "Salmos" {
			book.Name = "Salmo"
		}
		chapter = 1
		line = fmt.Sprintf(patternBook, book.Acronym, strings.Replace(line, prefix, "", 1))
		fmt.Printf(line)
		book.printHeader()
	}
	return line
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

func main() {
	printIndex()
	book = books[0]
	nextBook = books[0]
	for _, file := range getFiles() {
		if strings.Index(file.Name(), "index") == 0 {
			inBody := false
			inHeader := true
			for _, line := range getFileContent(file.Name()) {
				if inBody, inHeader = isBody(inBody, inHeader, line); inBody {
					line = printChapter(line)
					line = printBook(line)
					line = printOutline(line)
					line = printIntroduction(line)
				}
			}
		}
	}
	//	fmt.Printf("{\"%s\",\"%s\",%d}\n", book.Name, book.Acronym, chapter-1)
}

func initBooks() [66]Book {
	books := [66]Book{{"Gênesis", "gn", 50}, {"Êxodo", "ex", 40}, {"Levítico", "lv", 27}, {"Números", "nm", 36}, {"Deuteronômio", "dt", 34}, {"Josué", "js", 24}, {"Juízes", "jz", 21}, {"Rute", "rt", 4}, {"1 Samuel", "1sm", 31}, {"2 Samuel", "2sm", 24}, {"1 Reis", "1re", 22}, {"2 Reis", "2re", 25}, {"1 Crônicas", "1cr", 29}, {"2 Crônicas", "2cr", 36}, {"Esdras", "ed", 10}, {"Neemias", "ne", 13}, {"Ester", "et", 10}, {"Jó", "jo", 42}, {"Salmos", "sl", 150}, {"Provérbios", "pv", 31}, {"Eclesiastes", "ec", 12}, {"Cantares", "ct", 8}, {"Isaías", "is", 66}, {"Jeremias", "jr", 52}, {"Lamentações", "lm", 5}, {"Ezequiel", "ez", 48}, {"Daniel", "dn", 12}, {"Oséias", "os", 14}, {"Joel", "jl", 3}, {"Amós", "am", 9}, {"Obadias", "ob", 1}, {"Jonas", "jn", 4}, {"Miquéias", "mq", 7}, {"Naum", "na", 3}, {"Habacuque", "hc", 3}, {"Sofonias", "sf", 3}, {"Ageu", "ag", 2}, {"Zacarias", "zc", 14}, {"Malaquias", "mq", 4}, {"Mateus", "mt", 28}, {"Marcos", "mc", 16}, {"Lucas", "lc", 24}, {"João", "ja", 21}, {"Atos", "at", 28}, {"Romanos", "rm", 16}, {"1 Coríntios", "1co", 16}, {"2 Coríntios", "2co", 13}, {"Gálatas", "gl", 6}, {"Efésios", "ef", 6}, {"Filipenses", "fp", 4}, {"Colossenses", "cl", 4}, {"1 Tessalonicenses", "1ts", 5}, {"2 Tessalonicenses", "2ts", 3}, {"1 Timóteo", "1tm", 6}, {"2 Timóteo", "2tm", 4}, {"Tito", "tt", 3}, {"Filemom", "fl", 1}, {"Hebreus", "hb", 13}, {"Tiago", "tg", 5}, {"1 Pedro", "1pe", 5}, {"2 Pedro", "2pe", 3}, {"1 João", "1jo", 5}, {"2 João", "2jo", 1}, {"3 João", "3jo", 1}, {"Judas", "jd", 1}, {"Apocalipse", "ap", 22}}
	return books
}
