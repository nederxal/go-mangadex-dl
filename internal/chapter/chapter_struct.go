package chapter

/* Response type, vu qu'on ne prend les chapitres que 1 par 1 le tableau de DATA est seul.

{
  "result": "ok",
  "response": "collection",
  "data": [
    {
      "id": "89f27f84-951b-4c88-9ebb-798ece0fbcd2",
      "type": "chapter",
      "attributes": {
        "volume": "1",
        "chapter": "6",
        "title": "Shadow-Boxing",
        "translatedLanguage": "en",
        "externalUrl": null,
        "publishAt": "2018-02-05T18:38:15+00:00",
        "readableAt": "2018-02-05T18:38:15+00:00",
        "createdAt": "2018-02-05T18:38:15+00:00",
        "updatedAt": "2018-02-05T18:38:15+00:00",
        "pages": 20,
        "version": 1
      },
      "relationships": [ // cette merde est variable
        {
          "id": "b708d4dc-dbaf-4267-8482-7dc844490d50",
          "type": "scanlation_group"
        },
        {
          "id": "f7888782-0727-49b0-95ec-a3530c70f83b",
          "type": "manga"
        },
        {
          "id": "201fa20a-08a9-415d-ad63-77bbf40c8a0e",
          "type": "user"
        }
      ]
    }
  ],
  "limit": 1,
  "offset": 0,
  "total": 6
}

*/

type Attributes struct {
	Volume             string `json:"volume"`
	Chapter            string `json:"chapter"`
	Title              string `json:"title"`
	TranslatedLanguage string `json:"translatedLanguage"`
	ExternalUrl        string `json:"externalUrl"`
	PublishAt          string `json:"publishAt"`
	ReadableAt         string `json:"readableAt"`
	CreatedAt          string `json:"createdAt"`
	UpdatedAt          string `json:"updatedAt"`
	Pages              int    `json:"pages"`
	Version            int    `json:"version"`
}

type Relationships struct {
	Id   string `json:"id"`
	Type string `json:"type"`
}

type Data struct {
	Id                   string          `json:"id"`
	Type                 string          `json:"type"`
	ChapterAttributes    Attributes      `json:"attributes"`
	ChapterRelationships []Relationships `json:"relationships"`
}

type ChapterStruct struct {
	Result      string `json:"result"`
	Response    string `json:"response"`
	ChapterData []Data `json:"data"`
	Limit       int    `json:"limit"`
	Offset      int    `json:"offset"`
	Total       int    `json:"total"`
}
