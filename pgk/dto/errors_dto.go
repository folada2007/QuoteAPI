package dto

type Errors struct {
	Type    string `json:"type"`
	Message string `json:"message"`
}

type ResponseErrors struct {
	Errors []Errors `json:"errors"`
}
