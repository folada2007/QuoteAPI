package dto

type ResponseQuotes struct {
	ID     int64  `json:"id"`
	Quote  string `json:"quote"`
	Author string `json:"author"`
}

type RequestQuote struct {
	Quote  string `json:"quote"`
	Author string `json:"author"`
}
