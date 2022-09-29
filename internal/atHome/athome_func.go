package athome

import (
	"encoding/json"
	"io"
	"net/http"
	"net/url"
)

func GetAtHome(chapterUUID string) AtHome {
	url, err := url.JoinPath("https://api.mangadex.org/at-home/server/", chapterUUID)
	if err != nil {
		panic(err)
	}

	resp, err := http.Get(url)
	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)

	var ah *AtHome
	errj := json.Unmarshal(body, &ah)
	if errj != nil {
		panic(errj)
	}
	return *ah
}
