package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"time"
)

var (
	InfoLogger  *log.Logger
	UrlEntities []ShortUrlEntity
	Counter     = 0
	letters     = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890")
)

type ShortUrlEntity struct {
	id          int
	originalUrl string
	shortUrl    string
}

type RequestBodyGenerateUrl struct {
	Url string
}

func handleGenerateUrl(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		return
	}
	requestBodyBytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		InfoLogger.Println("deu ruim lendo o corpo")
	}
	// 										*RequestBodyGenerateUrl will throw an error
	var requestBodyStruct RequestBodyGenerateUrl
	err = json.Unmarshal(requestBodyBytes, &requestBodyStruct)
	if err != nil {
		InfoLogger.Println("Erro ao fazer unmarshall do json")
		io.WriteString(w, `
		{
			"mensagem": "requisicao invalida"
		}`)
		return
	}
	path := randomSequence(7)
	var shortUrlEntity ShortUrlEntity
	Counter += 1
	shortUrlEntity.id = Counter
	shortUrlEntity.shortUrl = path
	shortUrlEntity.originalUrl = requestBodyStruct.Url
	UrlEntities = append(UrlEntities, shortUrlEntity)
	w.Header().Set("content-type", "application/json")
	io.WriteString(w, fmt.Sprintf(`
	{
		"url": "%s"
	}
	`, "http://localhost:8003/"+path))
	return
}

func handleCheckUrl(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		return
	}
	path := strings.Split(r.URL.Path[1:], "/")[0]
	validPath := false
	var rightUrlEntity ShortUrlEntity
	for i := 0; i < len(UrlEntities); i++ {
		if strings.Compare(UrlEntities[i].shortUrl, path) == 0 {
			validPath = true
			rightUrlEntity = UrlEntities[i]
			break
		}
	}
	InfoLogger.Println(validPath)
	if validPath {
		http.Redirect(w, r, rightUrlEntity.originalUrl, 301)
	}
	return
}

func randomSequence(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func main() {
	InfoLogger = log.New(
		os.Stdout,
		"REQUEST INFO: ",
		log.Ldate|log.Ltime|log.Lmsgprefix,
	)

	multiplexer := http.NewServeMux()

	multiplexer.HandleFunc("/url", handleGenerateUrl)
	multiplexer.HandleFunc("/", handleCheckUrl)

	server := &http.Server{
		Addr:           ":8003",
		Handler:        multiplexer,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20, // its equivalent to (1 * 2) * 20 bytes = 1mb.
	}
	server.ListenAndServe()
}
