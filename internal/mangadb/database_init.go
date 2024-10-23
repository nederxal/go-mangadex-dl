package mangadb

import (
	"database/sql"
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
func CreateDatabase(pathDB string) {
	tmp, err := os.Create(pathDB)
	if err != nil {
		log.Error(err)
	}
	tmp.Close()

	db, err := sql.Open("sqlite3", pathDB)
	if err != nil {
		log.Error(err)
	}
	defer db.Close()

	db.Exec(REQ)
}
