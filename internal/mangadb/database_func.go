package mangadb

import (
	"database/sql"
	athome "go-mangadex-dl/internal/atHome"
	"go-mangadex-dl/internal/chapter"
	m "go-mangadex-dl/internal/manga"
	"strconv"
	"time"

	_ "github.com/mattn/go-sqlite3"
	log "github.com/sirupsen/logrus"
)

type myMangas struct {
	Id          int
	Name        string
	UUID        string
	NextChapter int
	Langue      string
}

func ListMangas(db *sql.DB) {
	rows := QueryDB(db)
	defer rows.Close()

	// parcours les mangas à télécharger
	for rows.Next() {
		var manga myMangas
		err := rows.Scan(&manga.Id, &manga.Name, &manga.UUID, &manga.NextChapter, &manga.Langue)
		if err != nil {
			log.Fatal(err)
		}

		// Tant qu'on trouve un chapitre pas vide on télécharge
		for {
			c := chapter.GetChapter(manga.UUID, manga.NextChapter, manga.Langue)

			if len(c.ChapterData) == 0 {
				log.Warnf("%s chapitre %d vide / inexistant", manga.Name, manga.NextChapter)
				defer UpdateDB(db, manga.Id, manga.NextChapter)
				break
			}

			ah := athome.GetAtHome(c.ChapterData[0].Id)
			//TODO: Gérer si tout le chapitre n'a pas été DL
			chapter.Download(ah, manga.Name, strconv.Itoa(manga.NextChapter))

			if m.GetMangaStatus(db, manga.Name, manga.UUID, manga.Id, manga.NextChapter) {
				defer RemoveFromDB(db, manga.Id)
			}

			manga.NextChapter += 1
		}
	}
}

func QueryDB(db *sql.DB) *sql.Rows {
	rows, err := db.Query("select * from mangas;")
	if err != nil {
		log.Warn(err)
	}
	return rows
}

func UpdateDB(db *sql.DB, id, nextChapter int) {
	_, err := db.Exec(`update mangas set next_chapter = ? where id = ?`, nextChapter, id)
	if err != nil {
		log.Error(err)
	}
}

func RemoveFromDB(db *sql.DB, id int) {
	_, err := db.Exec(`delete from mangas where id = ?`, id)
	if err != nil {
		log.Warn(err)

	}
}

func AddMangas(db *sql.DB, csvDoubleTab [][]string) {
	for _, line := range csvDoubleTab {
		time.Sleep(200 * time.Millisecond) //avoid ban
		mangaName := m.GetMangaNameFromUUID(line[0])
		lang := line[1]
		chapter := line[2]

		insertManga := "INSERT INTO mangas(name, name_UUID, next_chapter, langue) VALUES(?, ? ,?, ?)"
		stmt, err := db.Prepare(insertManga)
		if err != nil {
			log.Error(nil)
		}

		_, err = stmt.Exec(mangaName, line[0], chapter, lang)
		if err != nil {
			log.Error(err)
		}
	}

	// 	db.Exec/Query ? avec UUID/Name/Chapitre/Langue
}
