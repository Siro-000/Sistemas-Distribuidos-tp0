package common

import (
	"fmt"
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

func BuildPostBetRequest() (PostBetRequest, error) {
	numeroStr := os.Getenv("NUMERO")
	numeroInt, err := strconv.ParseInt(numeroStr, 10, 64)
	if err != nil {
		return PostBetRequest{}, fmt.Errorf("error NUMERO to int64: %v", err)
	}

	return PostBetRequest{
		FirstName: os.Getenv("NOMBRE"),
		LastName:  os.Getenv("APELLIDO"),
		Document:  os.Getenv("DOCUMENTO"),
		Birthdate: os.Getenv("NACIMIENTO"),
		Number:    numeroInt,
	}, nil
}
