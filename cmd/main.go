package main

import (
	athome "go-mangadex-dl/internal/atHome"
	chapter "go-mangadex-dl/internal/chapter"
)

func main() {

	// mangaName := "Hajime No Ippo"
	mangaUUID := "f7888782-0727-49b0-95ec-a3530c70f83b"
	nextChapter := "1"
	lang := "en"

	c := chapter.GetChapter(mangaUUID, nextChapter, lang)

	// on garde que ce chapitre doit être retéléchargé
	if c.ChapterData[0].ChapterAttributes.Pages == 0 {
		panic("chapitre vide")
	} else {
		ah := athome.GetAtHome(c.ChapterData[0].Id)
		// refaire le paramétrage de où DL le manga --> dans le /home/$USER/<manganame>/<chapternumber>/
		// Gérer si tout le chapitre n'a pas été DL
		chapter.Download(ah)
	}

}
