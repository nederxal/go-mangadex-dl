package chapter

import (
	"bufio"
	"encoding/json"
	"fmt"
	athome "go-mangadex-dl/internal/atHome"
	"io"
	"net/http"
	"net/url"
	"os"
	"time"
)

func Download(ah athome.AtHome) {
	// URL Example : https://uploads.mangadex.org/data/<HASH>/<IMG>
	for _, page := range ah.Chapter.Data {
		pageUrl, e := url.JoinPath(ah.BaseUrl, "data", ah.Chapter.Hash, page)
		if e != nil {
			panic("url foireuse")
		}

		fmt.Println(pageUrl)
		resp, err := http.Get(pageUrl)
		if err != nil {
			fmt.Println("Error get URL")
		}
		defer resp.Body.Close()

		r := bufio.NewReader(resp.Body)

		output, _ := os.Create(page)
		defer output.Close()

		w := bufio.NewWriter(output)

		r.WriteTo(w)
		time.Sleep(200 * time.Millisecond)

	}

}

// récupère la structure du chapitre
func GetChapter(mangaUUID string, chapter string, lang string) ChapterStruct {
	url, err := url.Parse("https://api.mangadex.org/chapter/")
	if err != nil {
		panic(err)
	}

	q := url.Query()
	q.Add("limit", "1")                 // --> Quantités de chapitres ressortis --> si plusieurs chapitres de la même langue on prend le premier proposé
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