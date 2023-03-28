package athome

import (
	"encoding/json"
	"io"
	"net/http"
	"net/url"

	log "github.com/sirupsen/logrus"
)

const URLATHOME = "https://api.mangadex.org/at-home/server/"

// Get MangaDex@Home server URL
func GetAtHome(chapterUUID string) AtHome {
	url, err := url.JoinPath(URLATHOME, chapterUUID)
	if err != nil {
		log.Panic(err)
	}

	resp, err := http.Get(url)
	if err != nil {
		log.Panic(err)
	}

	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)

	var ah *AtHome
	errj := json.Unmarshal(body, &ah)
	if errj != nil {
		log.Panic(errj)
	}
	return *ah
}
