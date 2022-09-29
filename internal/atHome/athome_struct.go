package athome

/* Reponse type :

{
  "result": "ok",
  "baseUrl": "https://uploads.mangadex.org",
  "chapter": {
    "hash": "c6fc48b6fb79d5604c0847d21ed53227",
    "data": [
      "1-002cd644db59fce2592f6c187796a17cfdbabcf8829a4ac6b0a8e2d70c393210.png",
      "2-1d04e9252cd45897b3b5de9fc591db1b278b3808167214fd0f90579e872898a5.jpg"
	  ...
    ],
    "dataSaver": [
      "1-4592b5350e20cd9c94050d64996e481bd7d742e25c8acc296fc90b75f954ac7e.jpg",
      "2-261147535269e69c4c3f58e92ffc66a96b038327270d88cf2aafa89df89b677f.jpg"
	  ...
    ]
  }
}
*/

type AtHome struct {
	Result  string `json:"result"`
	BaseUrl string `json:"baseUrl"`
	Chapter struct {
		Hash      string   `json:"hash"`
		Data      []string `json:"data"`
		DataSaver []string `json:"dataSaver"`
	} `json:"chapter"`
}
