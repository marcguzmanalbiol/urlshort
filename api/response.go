package api

type Message struct {
	Message string `json:"message"`
}

type SuccessfullyCreated struct {
	Message  string `json:"message"`
	ShortURL string `json:"ShortURL"`
}

type Detail struct {
	OriginalURL    string `json:"OriginalURL"`
	ShortURL       string `json:"ShortURL"`
	ExpirationDate string `json:"ExpirationDate"`
	UsedCount      int    `json:"UsedCount"`
}
