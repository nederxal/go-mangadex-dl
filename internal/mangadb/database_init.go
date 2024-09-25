package mangadb

import (
	"database/sql"
	"encoding/csv"
	"os"

	log "github.com/sirupsen/logrus"
)

const REQ string = `
CREATE TABLE 'mangas' (
'id'	INTEGER UNIQUE,
'name'	TEXT,
'name_UUID'	TEXT UNIQUE,
'next_chapter'	INTEGER,
'langue'	TEXT,
PRIMARY KEY('id' AUTOINCREMENT))
`

// If the database doesn't exists at the defined path create it and fill it
func CreateDatabase(pathDB, UUIDList string) {
	tmp, err := os.Create(pathDB)
	if err != nil {
		log.Error(err)
	}
	tmp.Close()

	db, err := sql.Open("sqlite3", pathDB)
	if err != nil {
		log.Error(err)
	}

	db.Exec(REQ)
	defer db.Close()

	stat, err := os.Stat(UUIDList)
	if err != nil {
		log.Warn("File not found won't add new mangas")
	} else {
		if stat.Size() > 0 {
			// Parcourir le CSV des UUID de mangas et du premier chapitre à télécharger -> ajouter dans la base -> continuer le programme
			file, _ := os.Open(UUIDList)
			defer file.Close()

			reader := csv.NewReader(file)
			csvDoubleTab, _ := reader.ReadAll()
			AddMangas(db, csvDoubleTab)

		} else {
			log.Warn("File is empty.")
		}
	}

}
