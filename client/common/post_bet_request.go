package common

import (
	"os"
	"strconv"
)

type PostBetRequest struct {
	FirstName string
	LastName  string
	Document  string
	Birthdate string
	Number    int64
}

func BuildPostBetRequest() PostBetRequest {
	numeroStr := os.Getenv("NUMERO")
	numeroInt, err := strconv.ParseInt(numeroStr, 10, 64)
	if err != nil {
		log.Fatalf("Error NUMERO to int64: %v", err)
	}

	return PostBetRequest{
		FirstName: os.Getenv("NOMBRE"),
		LastName:  os.Getenv("APELLIDO"),
		Document:  os.Getenv("DOCUMENTO"),
		Birthdate: os.Getenv("NACIMIENTO"),
		Number:    numeroInt,
	}
}
