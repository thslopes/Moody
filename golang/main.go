package main

import (
    "fmt"
    "io/ioutil"
    "os"
    "strings"
)

func check(e error) {
    if e != nil {
        panic(e)
    }
}

var folder = "/home/thiago/github.com/thslopes/Moody"
var prefix = "<p class=\"calibre2\"><b class=\"calibre1\">"
var formatoCapitulo = "<p class=\"calibre2\"><b class=\"calibre1\" id=\"%s%d\">%s\n"
var formatoLivro = "<p class=\"calibre2\"><b class=\"calibre1\" id=\"%s\">%s\n"
var livros, siglas = Init()
var livro string
var sigla string
var proxLivro string
var capitulo = 1
var livroIdx = -1

func GetFiles() []os.FileInfo {
  files, err := ioutil.ReadDir(folder)
  check(err)
  return files
}

func GetFileContent(fileName string) []string {
  dat, err := ioutil.ReadFile(fmt.Sprintf("%s/%s",folder, fileName))
  check(err)
  lines := strings.Split(string(dat), "\n")
  return lines
}

func FormatCapitulo(line string) string {
  trimLine := strings.TrimLeft(line, " ")
  io := strings.Index(trimLine, fmt.Sprintf("%s%s %d", prefix, livro, capitulo))
  io2 := strings.Index(trimLine, fmt.Sprintf("%s%s cap. %d", prefix, livro, capitulo))
  if io == 0 || io2 == 0 {
    line = fmt.Sprintf(formatoCapitulo, sigla, capitulo, strings.Replace(line, prefix, "", 1))
    capitulo++
  }
  return line
}

func FormatIntroducao(line string) string {
  trimLine := strings.TrimLeft(line, " ")
  io := strings.Index(trimLine, fmt.Sprintf("%sINTRODUÇÃO </b></p>", prefix))
  if io == 0 {
    line = fmt.Sprintf(formatoLivro, sigla + "intro", strings.Replace(line, prefix, "", 1))
  }
  return line
}

func FormatEsboco(line string) string {
  trimLine := strings.TrimLeft(line, " ")
  io := strings.Index(trimLine, fmt.Sprintf("%sESBOÇO </b></p>", prefix))
  if io == 0 {
    line = fmt.Sprintf(formatoLivro, sigla + "esboco", strings.Replace(line, prefix, "", 1))
  }
  return line
}

func FormatLivro(line string) string {
  trimLine := strings.TrimLeft(line, " ")
  if livroIdx < 65 {
    proxLivro = strings.ToUpper(livros[livroIdx + 1])
    if proxLivro == "CANTARES" {
      proxLivro = "CANTARES DE SALOMÃO"
    }
  }
  io := strings.Index(trimLine, fmt.Sprintf("%s%s </b></p>", prefix, proxLivro))
  io2 := strings.Index(trimLine, fmt.Sprintf("%s  </b> <b class=\"calibre1\">%s </b></p>", prefix, proxLivro))
  if io == 0 || io2 == 0 {
    livroIdx += 1
    livro = livros[livroIdx]
    sigla = siglas[livroIdx]
    if livro == "Salmos" {
      livro = "Salmo"
    }
    capitulo = 1
    line = fmt.Sprintf(formatoLivro, sigla, strings.Replace(line, prefix, "", 1))
  }
  return line
}

func IsBody(isBody bool, isHeader bool, line string) (bool, bool) {
  if !isHeader {
    isBody = true
  }
  trimLine := strings.TrimLeft(line, " ")
  io := strings.Index(trimLine, "<body")
  if io == 0 {
    isHeader = false
  }

  return isBody, isHeader
}

func main() {
    livro = livros[0]
    proxLivro = livros[0]
    for _, file := range GetFiles() {
      if strings.Index(file.Name(), "index") == 0 {
        isBody := false
        isHeader := true
        for _, line := range GetFileContent(file.Name()) {
          if isBody, isHeader = IsBody(isBody, isHeader, line); isBody {
            line = FormatCapitulo(line)
            line = FormatLivro(line)
            line = FormatEsboco(line)
            line = FormatIntroducao(line)
          }
        }
      }
    }
}


func Init() ([66]string, [66]string) {
  livros := [66]string{"Gênesis","Êxodo","Levítico","Números","Deuteronômio","Josué","Juízes","Rute","1 Samuel","2 Samuel","1 Reis","2 Reis","1 Crônicas","2 Crônicas","Esdras","Neemias","Ester","Jó","Salmos","Provérbios","Eclesiastes","Cantares","Isaías","Jeremias","Lamentações","Ezequiel","Daniel","Oséias","Joel","Amós","Obadias","Jonas","Miquéias","Naum","Habacuque","Sofonias","Ageu","Zacarias","Malaquias","Mateus","Marcos","Lucas","João","Atos","Romanos","1 Coríntios","2 Coríntios","Gálatas","Efésios","Filipenses","Colossenses","1 Tessalonicenses","2 Tessalonicenses","1 Timóteo","2 Timóteo","Tito","Filemom","Hebreus","Tiago","1 Pedro","2 Pedro","1 João","2 João","3 João","Judas","Apocalipse"}
  sigla := [66]string{"gn","ex","lv","nm","dt","js","jz","rt","1sm","2sm","1re","2re","1cr","2cr","ed","ne","et","jo","sl","pv","ec","ct","is","jr","lm","ez","dn","os","jl","am","ob","jn","mq","na","hc","sf","ag","zc","mq","mt","mc","lc","ja","at","rm","1co","2co","gl","ef","fp","cl","1ts","2ts","1tm","2tm","tt","fl","hb","tg","1pe","2pe","1jo","2jo","3jo","jd","ap"}
  return livros, sigla
}
