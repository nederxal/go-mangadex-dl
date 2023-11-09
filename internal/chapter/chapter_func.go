package chapter

import (
	"archive/zip"
	"bufio"
	"encoding/json"
	"fmt"
	athome "go-mangadex-dl/internal/atHome"
	"io"
	"net/http"
	"net/url"
	"os"
	"path"
	"strconv"
	"time"

	log "github.com/sirupsen/logrus"
)

const BASEURLCHAPTER string = "https://api.mangadex.org/chapter/"

// On va créer directement le cbz dans là où ça doit être télécharger (oui smart phrase on est bien)
func Download(ah athome.AtHome, mangaName, mangaNextChapter string) {
	// On s'assure que le répertoire du manga existe (au pire ça le créé)
	destFolder := path.Join(os.Getenv("HOME"), "MangadexDownloads", mangaName)
	err := os.MkdirAll(destFolder, os.ModePerm)
	if err != nil {
		log.Panic(err)
	}

	// Creation de l'archive cbz
	cbzPath := path.Join(destFolder, mangaNextChapter+".cbz")
	cbzChapter, err := os.Create(cbzPath)
	if err != nil {
		panic(err)
	}
	defer cbzChapter.Close()

	cbzWriter := zip.NewWriter(cbzChapter)

	// On va récupérer les pages une par une et les ajouter à l'archive
	for _, page := range ah.Chapter.Data {

		pageUrl, err := url.JoinPath(ah.BaseUrl, "data", ah.Chapter.Hash, page)
		fmt.Println(pageUrl)
		if err != nil {
			log.Error("url foireuse")
		}

		resp, err := http.Get(pageUrl)
		if err != nil {
			log.Error("Error get URL")
		}
		defer resp.Body.Close()

		r := bufio.NewReader(resp.Body)

		w, err := cbzWriter.Create(page)
		if err != nil {
			panic(err)
		}

		if _, err := io.Copy(w, r); err != nil {
			panic(err)
		}

		// Temps de pause pour pas se faire striker parle limiteur de mangadex
		time.Sleep(200 * time.Millisecond)
	}

	cbzWriter.Close()

}

// récupère la structure du chapitre
func GetChapter(mangaUUID string, chapter int, lang string) ChapterStruct {
	url, err := url.Parse(BASEURLCHAPTER)
	if err != nil {
		log.Panic(err)
	}

	q := url.Query()
	q.Add("limit", "1")                     // --> Quantités de chapitres ressortis --> si plusieurs chapitres de la même langue on prend le premier proposé
	q.Add("manga", mangaUUID)               // --> UUID du manga
	q.Add("chapter", strconv.Itoa(chapter)) // --> Numero du chapitre recherché à convertir en string parce que url.Query
	q.Add("translatedLanguage[]", lang)     // --> langue ... fr possible mais surtout en anglais
	// Les 3 options suivantes sont obligatoires (au moins une)
	q.Add("contentRating[]", "safe")
	q.Add("contentRating[]", "suggestive")
	q.Add("contentRating[]", "erotica")

	url.RawQuery = q.Encode()

	r, err := http.Get(url.String())
	if err != nil {
		log.Panic(err)
	}
	defer r.Body.Close()

	body, _ := io.ReadAll(r.Body)

	var chap *ChapterStruct
	err = json.Unmarshal(body, &chap)
	if err != nil {
		log.Panic(err)
	}

	return *chap
}
