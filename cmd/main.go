package main

import (
	"database/sql"
	m "go-mangadex-dl/internal/manga"
	"os"
	"path"

	log "github.com/sirupsen/logrus"
)

func main() {
	f, err := os.OpenFile(path.Join(os.Getenv("HOME"), "MangadexDownloads", "log.txt"), os.O_WRONLY|os.O_CREATE, 0755)
	if err != nil {
		log.Panic("Can't create log file !")
	}

	log.SetOutput(f)

	pathDB := path.Join(os.Getenv("HOME"), "MangadexDownloads", "db.sqlite")

	if _, err := os.Stat(pathDB); err == nil {
		db, err := sql.Open("sqlite3", pathDB)

		if err != nil {
			log.Fatal(err)
		}

		defer db.Close()
		m.ListMangas(db)

	} else {
		log.Fatal("Database doesn't exist, bye")
	}
}
