package mangadb

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
	log "github.com/sirupsen/logrus"
)

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

// func AddMangas(db *sql.DB, ze double tab) int {
// 	for line in csvDoubleTab
//		mangaName := getMangaName(UUID)
//		db.Exec/Query ? avec UUID/Name/Chapitre/Lange
// }
