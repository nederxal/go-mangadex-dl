package main

import (
	"database/sql"
	"fmt"
	athome "go-mangadex-dl/internal/atHome"
	chapter "go-mangadex-dl/internal/chapter"
	"log"
	"strconv"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	type mangas struct {
		id           int
		name         string
		name_UUID    string
		next_chapter int
		langue       string
	}

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
		var manga mangas
		// mangaName := "Hajime No Ippo"
		// mangaUUID := "f7888782-0727-49b0-95ec-a3530c70f83b"
		// nextChapter := "1"
		// lang := "en"
		err = rows.Scan(&manga.id, &manga.name, &manga.name_UUID, &manga.next_chapter, &manga.langue)
		if err != nil {
			log.Fatal(err)
		}
		// fmt.Println(manga.name, manga.name_UUID, manga.next_chapter, manga.langue)

		// on garde que ce chapitre doit être retéléchargé
		for {
			c := chapter.GetChapter(manga.name_UUID, manga.next_chapter, manga.langue)
			if len(c.ChapterData) == 0 {
				fmt.Printf("%s chapitre %d vide / inexistant", manga.name, manga.next_chapter)
				break
			} else {
				ah := athome.GetAtHome(c.ChapterData[0].Id)
				// refaire le paramétrage de où DL le manga --> dans le /home/$USER/<manganame>/<chapternumber>/
				// Gérer si tout le chapitre n'a pas été DL
				chapter.Download(ah, manga.name, strconv.Itoa(manga.next_chapter))
			}
			manga.next_chapter += 1
		}
		// mettre à jour la DB avec le prochain chapitre à télécharger
	}

}
