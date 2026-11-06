package domain

import (
	"strings"
)

type Deck struct {
	ID          string
	Name        string
	Format      string
	Description string
}

func NewDeck(name, format, description string) (Deck, error) {
	name = strings.TrimSpace(name)
	if name == "" {
		return Deck{}, ErrDeckNameRequired
	}

	return Deck{
		Name:        name,
		Format:      strings.TrimSpace(format),
		Description: strings.TrimSpace(description),
	}, nil
}

func UpdateDeck(oldDeck, newDeck Deck) (Deck, error) {
	if newDeck.Name == "" {
		newDeck.Name = oldDeck.Name
	}
	if newDeck.Format == "" {
		newDeck.Format = oldDeck.Format
	}
	if newDeck.Description == "" {
		newDeck.Description = oldDeck.Description
	}
	return newDeck, nil
}
