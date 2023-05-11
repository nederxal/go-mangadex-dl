package manga

import (
	"database/sql"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"strconv"

	athome "go-mangadex-dl/internal/atHome"
	"go-mangadex-dl/internal/chapter"
	mdb "go-mangadex-dl/internal/mangadb"

	log "github.com/sirupsen/logrus"
)

const GETMANGA = "https://api.mangadex.org/manga/"

func ListMangas(db *sql.DB) {
	rows := mdb.QueryDB(db)
	defer rows.Close()

	// parcours les mangas à télécharger
	for rows.Next() {
		var manga mdb.Mangas
		err := rows.Scan(&manga.Id, &manga.Name, &manga.UUID, &manga.NextChapter, &manga.Langue)
		if err != nil {
			log.Fatal(err)
		}

		// Tant qu'on trouve un chapitre pas vide on télécharge
		for {
			c := chapter.GetChapter(manga.UUID, manga.NextChapter, manga.Langue)

			if len(c.ChapterData) == 0 {
				log.Warnf("%s chapitre %d vide / inexistant", manga.Name, manga.NextChapter)
				defer mdb.UpdateDB(db, manga.Id, manga.NextChapter)
				break
			}

			ah := athome.GetAtHome(c.ChapterData[0].Id)
			//TODO: Gérer si tout le chapitre n'a pas été DL
			chapter.Download(ah, manga.Name, strconv.Itoa(manga.NextChapter))

			if getMangaStatus(db, manga.Name, manga.UUID, manga.Id, manga.NextChapter) {
				defer mdb.RemoveFromDB(db, manga.Id)
			}

			manga.NextChapter += 1
		}
	}
}

func getMangaStatus(db *sql.DB, name, mangaUUID string, id, chapter int) bool {
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

	var manga *Manga

	err = json.Unmarshal(body, &manga)
	if err != nil {
		log.Panic(err)
	}

	if manga.Data.Attributes.Status == "completed" && manga.Data.Attributes.LastChapter == strconv.Itoa(chapter) {
		return true
	}

	// if manga.Data.Attributes.Status == "completed" && manga.Data.Attributes.LastChapter != strconv.Itoa(chapter) {
	// 	fmt.Println(name + " pas terminé ... à voir si il manque des chapitres dans la langue")
	// 	return true
	// }

	return false
}
