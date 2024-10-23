package manga

import (
	"database/sql"
	"encoding/csv"
	"encoding/json"
	"fmt"
	athome "go-mangadex-dl/internal/atHome"
	"go-mangadex-dl/internal/chapter"
	mdb "go-mangadex-dl/internal/mangadb"
	"io"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"time"

	log "github.com/sirupsen/logrus"
)

// Seulement ce qui est nécessaire depuis Mangadex
type mangaDexInfo struct {
	Data struct {
		Attributes struct {
			Title struct {
				En string `json:"en"` // On récupère le titre toujours en anglais
			}
			Status      string `json:"status"`
			LastChapter string `json:"lastChapter"` // aaaaand it's not an int ...
		}
	}
}

type myMangas struct {
	Id          int
	Name        string
	UUID        string
	NextChapter int
	Langue      string
}

const GETMANGA = "https://api.mangadex.org/manga/"

func AddMangas(db *sql.DB, uuidlist string) {
	stat, err := os.Stat(uuidlist)
	if err != nil {
		log.Warn("File not found won't add new mangas")
	} else {
		if stat.Size() > 0 {
			// Parcourir le CSV des UUID de mangas et du premier chapitre à télécharger -> ajouter dans la base -> continuer le programme
			file, _ := os.Open(uuidlist)
			defer file.Close()

			reader := csv.NewReader(file)
			csvDoubleTab, _ := reader.ReadAll()

			for _, line := range csvDoubleTab {
				name := getMangaNameFromUUID(line[0])
				langue := line[1]
				next_chapter, _ := strconv.Atoi(line[2])
				fmt.Println(name, langue)
				time.Sleep(200 * time.Millisecond)
				log.Infof("Ajout de %s à partir du chapitre %d", name, next_chapter)
				mdb.InsertDB(db, name, line[0], langue, next_chapter)
			}
		} else {
			log.Warn("File is empty.")
		}
	}
}

// Liste les mangas à télécharger et ensuite va chercher les chapitres 1 par 1
func ListMangas(db *sql.DB) {

	var liste []myMangas
	rows := mdb.GetAllMangas(db)

	for rows.Next() {
		var manga myMangas
		err := rows.Scan(&manga.Id, &manga.Name, &manga.UUID, &manga.NextChapter, &manga.Langue)
		if err != nil {
			log.Fatal(err)
		}

		liste = append(liste, manga)
	}
	// Obliger de fermer le rows ici sinon on ne peut pas mettre la db à jour au fil de l'eau
	rows.Close()

	for _, l := range liste {
		for {
			log.Infof("Manga : %s chapitre %d", l.Name, l.NextChapter)
			c := chapter.GetChapter(l.UUID, l.NextChapter, l.Langue)

			if c.Total == 0 || len(c.ChapterData) == 0 {
				log.Warnf("%s chapitre %d vide / inexistant", l.Name, l.NextChapter)
				mdb.UpdateDB(db, l.Id, l.NextChapter)
				break
			}

			ah := athome.GetAtHome(c.ChapterData[0].Id)
			//TODO: Gérer si tout le chapitre n'a pas été DL
			if chapter.Download(ah, l.Name, strconv.Itoa(l.NextChapter)) {
				mdb.UpdateDB(db, l.Id, l.NextChapter+1)
			}
			l.NextChapter += 1
		}
	}
}

// Just to get manga name and insert it into database
func getMangaNameFromUUID(uuid string) string {
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

	if mangaDexInfo.Data.Attributes.Status == "completed" && mangaDexInfo.Data.Attributes.LastChapter != strconv.Itoa(chapter) {
		log.Info(name + " pas terminé ... à voir si il manque des chapitres dans la langue")
	}

	return false
}
