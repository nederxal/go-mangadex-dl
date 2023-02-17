package mangadb

import (
	"database/sql"
	"fmt"
	athome "go-mangadex-dl/internal/atHome"
	"go-mangadex-dl/internal/chapter"
	"log"
	"strconv"

	_ "github.com/mattn/go-sqlite3"
)

const SQLQUERY string = "select * from mangas;"

func QueryDatabase() {
	db, err := sql.Open("sqlite3", "../db.sqlite")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	rows, err := db.Query(SQLQUERY)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		var manga Mangas
		err = rows.Scan(&manga.id, &manga.name, &manga.nameUUID, &manga.nextChapter, &manga.langue)
		if err != nil {
			log.Fatal(err)
		}

		for {
			c := chapter.GetChapter(manga.nameUUID, manga.nextChapter, manga.langue)

			if len(c.ChapterData) == 0 {
				fmt.Printf("%s chapitre %d vide / inexistant", manga.name, manga.nextChapter)
				defer updateDB(db, manga.id, manga.nextChapter)
				break
			} else {
				ah := athome.GetAtHome(c.ChapterData[0].Id)
				//TODO: Gérer si tout le chapitre n'a pas été DL
				chapter.Download(ah, manga.name, strconv.Itoa(manga.nextChapter))
			}
			manga.nextChapter += 1
		}
	}
}

func updateDB(db *sql.DB, id, nextChapter int) {
	_, err := db.Exec(`update mangas set next_chapter = ? where id = ?`, nextChapter, id)
	if err != nil {
		panic(err)
	}
}
