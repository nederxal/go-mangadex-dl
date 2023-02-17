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

func QueryDatabase() {

	db, err := sql.Open("sqlite3", "../db.sqlite")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	sqlQuery := "select * from mangas;"

	rows, err := db.Query(sqlQuery)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		var manga Mangas
		err = rows.Scan(&manga.id, &manga.name, &manga.name_UUID, &manga.next_chapter, &manga.langue)
		if err != nil {
			log.Fatal(err)
		}

		// on garde que ce chapitre doit être retéléchargé
		for {
			c := chapter.GetChapter(manga.name_UUID, manga.next_chapter, manga.langue)
			if len(c.ChapterData) == 0 {
				fmt.Printf("%s chapitre %d vide / inexistant", manga.name, manga.next_chapter)
				//TODO: mettre à jour la DB avec le prochain chapitre à télécharger
				break
			} else {
				ah := athome.GetAtHome(c.ChapterData[0].Id)
				//TODO: Gérer si tout le chapitre n'a pas été DL
				chapter.Download(ah, manga.name, strconv.Itoa(manga.next_chapter))
			}
			manga.next_chapter += 1
		}
	}
}
