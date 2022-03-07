package models

import (
	"errors"
	"fmt"
	"strings"
)

type Album struct {
	ID     int     `json:"id"`
	Title  string  `json:"title"`
	Artist string  `json:"artist"`
	Price  float64 `json:"price"`
}

func (a Album) AlbumIsValid() error {
	var validationErrors []string
	const MISSING_REQUIRED_FIELD_MSG = "missing required field %s"

	if a.Artist == "" {
		validationErrors = append(validationErrors, fmt.Sprintf(MISSING_REQUIRED_FIELD_MSG, "artist"))
	}
	if len(validationErrors) > 0 {
		for _, e := range validationErrors {
			fmt.Println(e)
		}

		errorMessage := fmt.Sprintf("The following errors were encountered: %s", strings.Join(validationErrors, "; "))
		return errors.New(errorMessage)
	}

	return nil
}
