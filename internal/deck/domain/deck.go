package domain

import "strings"

type Deck struct {
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
