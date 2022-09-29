package chapter

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
)

func Download() {
	// URL Example : https://uploads.mangadex.org/data/<HASH>/<IMG>
	resp, err := http.Get("https://uploads.mangadex.org/data/c6fc48b6fb79d5604c0847d21ed53227/1-002cd644db59fce2592f6c187796a17cfdbabcf8829a4ac6b0a8e2d70c393210.png")
	if err != nil {
		fmt.Println("Error get URL")
	}
	defer resp.Body.Close()

	r := bufio.NewReader(resp.Body)

	output, _ := os.Create("toto.jpg")
	defer output.Close()

	w := bufio.NewWriter(output)

	r.WriteTo(w)
}

func GetChapter(mangaUUID string, chapter string, lang string) ChapterStruct {
	url, err := url.Parse("https://api.mangadex.org/chapter/")
	if err != nil {
		panic(err)
	}

	q := url.Query()
	q.Add("limit", "1")
	q.Add("manga", mangaUUID)           // --> UUID du manga
	q.Add("chapter", chapter)           // --> Numero du chapitre recherché
	q.Add("translatedLanguage[]", lang) // --> langue ... fr possible mais surtout en anglais
	// Les 3 options suivantes sont obligatoires (au moins une)
	q.Add("contentRating[]", "safe")
	q.Add("contentRating[]", "suggestive")
	q.Add("contentRating[]", "erotica")

	url.RawQuery = q.Encode()

	r, err := http.Get(url.String())
	if err != nil {
		panic(err)
	}

	defer r.Body.Close()
	body, _ := io.ReadAll(r.Body)

	var chap *ChapterStruct
	errj := json.Unmarshal(body, &chap)
	if errj != nil {
		panic(errj)
	}

	return *chap
}
