package main

import (
	"database/sql"
	"fmt"
	m "go-mangadex-dl/internal/manga"
	mdb "go-mangadex-dl/internal/mangadb"
	"os"
	"path"

	log "github.com/sirupsen/logrus"
)

func main() {
	// Voir pour la création du dossier MangadexDownloads dans le $HOME
	logFile := path.Join(os.Getenv("HOME"), "MangadexDownloads", "log.txt")
	pathDB := path.Join(os.Getenv("HOME"), "MangadexDownloads", "db.sqlite")
	UUIDList := path.Join(os.Getenv("HOME"), "MangadexDownloads", "uuid.csv")

	f, err := os.OpenFile(logFile, os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		fmt.Println(err)
		panic("Can't create log file !")
	}

	log.SetOutput(f)

	_, err = os.Stat(pathDB)
	if err != nil {
		log.Warn("Database doesn't exist ... creating it ...")
	}
	// Creation si n'existe pas et ajoute des mangas
	mdb.CreateDatabase(pathDB, UUIDList)

	db, _ := sql.Open("sqlite3", pathDB)

	defer db.Close()
	os.Exit(1)
	m.ListMangas(db)
}
