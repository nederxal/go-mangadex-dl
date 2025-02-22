package mangadb

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
	log "github.com/sirupsen/logrus"
)

func GetAllMangas(db *sql.DB) *sql.Rows {
	rows, err := db.Query("select * from mangas;")
	if err != nil {
		log.Warn(err)
	}
	return rows
}

func UpdateDB(db *sql.DB, id, nextChapter int) {
	_, err := db.Exec(`UPDATE mangas SET next_chapter = ? WHERE id = ?`, nextChapter, id)
	if err != nil {
		log.Error(err)
	}
}

func InsertDB(db *sql.DB, name, name_UUID, langue string, next_chapter int) {
	insertManga := "INSERT INTO mangas(name, name_UUID, next_chapter, langue) VALUES(?, ? ,?, ?)"
	stmt, err := db.Prepare(insertManga)
	if err != nil {
		log.Error(nil)
	}
	_, err = stmt.Exec(name, name_UUID, next_chapter, langue)
	if err != nil {
		log.Error(err)
	}

}

func RemoveFromDB(db *sql.DB, uuid string) {
	_, err := db.Exec(`delete from mangas where name_UUID = ?`, uuid)
	if err != nil {
		log.Warn(err)

	}
}
