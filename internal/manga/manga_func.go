package manga

import (
	"database/sql"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"strconv"

	log "github.com/sirupsen/logrus"
)

// Only needed data to gather from mangadex
type mangaDexInfo struct {
	Data struct {
		Attributes struct {
			Title struct {
				En string `json:"en"` // title will always exist in English, otherwise check in altTitle but may won't exist ...
			}
			Status      string `json:"status"`
			LastChapter string `json:"lastChapter"` // aaaaand it's not an int ...
		}
	}
}

const GETMANGA = "https://api.mangadex.org/manga/"

// Just to get manga name and insert it into database
func GetMangaNameFromUUID(uuid string) string {
	mangaUrl, err := url.JoinPath(GETMANGA, uuid)
	if err != nil {
		log.Error("url foireuse")
	}

	resp, err := http.Get(mangaUrl)
	if err != nil {
		log.Error("Error get URL")
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	var mangaDexInfo *mangaDexInfo

	err = json.Unmarshal(body, &mangaDexInfo)
	if err != nil {
		log.Panic(err)
	}

	return mangaDexInfo.Data.Attributes.Title.En
}

// To run at the end and clean database from ended mangas
func GetMangaStatus(db *sql.DB, name, mangaUUID string, id, chapter int) bool {
	mangaUrl, err := url.JoinPath(GETMANGA, mangaUUID)
	if err != nil {
		log.Error("url foireuse")
	}

	resp, err := http.Get(mangaUrl)
	if err != nil {
		log.Error("Error get URL")
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	var mangaDexInfo *mangaDexInfo

	err = json.Unmarshal(body, &mangaDexInfo)
	if err != nil {
		log.Panic(err)
	}

	if mangaDexInfo.Data.Attributes.Status == "completed" && mangaDexInfo.Data.Attributes.LastChapter == strconv.Itoa(chapter) {
		return true
	}

	// if manga.Data.Attributes.Status == "completed" && manga.Data.Attributes.LastChapter != strconv.Itoa(chapter) {
	// 	fmt.Println(name + " pas terminé ... à voir si il manque des chapitres dans la langue")
	// 	return true
	// }

	return false
}
