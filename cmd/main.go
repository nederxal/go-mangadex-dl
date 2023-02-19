package main

import (
	"fmt"
	mdb "go-mangadex-dl/internal/mangadb"
	"os"
	"path"
)

func main() {
	pathDB := path.Join(os.Getenv("HOME"), "MangadexDownloads", "db.sqlite")
	if _, err := os.Stat(pathDB); err == nil {
		mdb.QueryDatabase(pathDB)
	} else {
		fmt.Println("The database doesn't exist !")
		os.Exit(1)
	}

}
