package main

import (
	"fmt"
	athome "go-mangadex-dl/internal/atHome"
	chapter "go-mangadex-dl/internal/chapter"
)

func main() {

	// mangaName := "Hajime No Ippo"
	mangaUUID := "f7888782-0727-49b0-95ec-a3530c70f83b"
	nextChapter := "1"
	lang := "en"

	r := chapter.GetChapter(mangaUUID, nextChapter, lang)
	fmt.Println(r)

	ah := athome.GetAtHome("6289e6ee-32b4-4020-b8f1-5de988008732")
	fmt.Println(ah)
}
