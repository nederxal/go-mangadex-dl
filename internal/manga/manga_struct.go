package manga

// Only needed data
type Manga struct {
	Data struct {
		Attributes struct {
			Status      string `json:"status"`
			LastChapter string `json:"lastChapter"` // aaaaand it's not an int ...
		}
	}
}
