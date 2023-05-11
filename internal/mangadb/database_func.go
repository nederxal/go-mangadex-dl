package mangadb

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
	log "github.com/sirupsen/logrus"
)

const QUERYMANGAS string = "select * from mangas;"

func QueryDB(db *sql.DB) *sql.Rows {
	rows, err := db.Query(QUERYMANGAS)
	if err != nil {
		log.Warn(err)
	}
	return rows
}

func UpdateDB(db *sql.DB, id, nextChapter int) {
	_, err := db.Exec(`update mangas set next_chapter = ? where id = ?`, nextChapter, id)
	if err != nil {
		panic(err)
	}
}

func RemoveFromDB(db *sql.DB, id int) {
	_, err := db.Exec(`delete from mangas where id = ?`, id)
	if err != nil {
		fmt.Println(err)
		log.Fatal(err)

	}
}
