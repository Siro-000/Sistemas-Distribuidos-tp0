package common

import (
	"os"
)

type PostBetRequest struct {
	FirstName string `json:"nombre"`
	LastName  string `json:"apellido"`
	Document  string `json:"documento"`
	Birthdate string `json:"nacimiento"`
	Number    string `json:"numero"`
}

func BuildPostBetRequest() PostBetRequest {
	return PostBetRequest{
		FirstName: os.Getenv("NOMBRE"),
		LastName:  os.Getenv("APELLIDO"),
		Document:  os.Getenv("dni"),
		Birthdate: os.Getenv("NACIMIENTO"),
		Number:    os.Getenv("NUMERO"),
	}
}
