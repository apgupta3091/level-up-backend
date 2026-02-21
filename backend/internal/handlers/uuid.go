package handlers

import (
	"fmt"

	"github.com/google/uuid"
)

func parseUUID(s string) (uuid.UUID, error) {
	id, err := uuid.Parse(s)
	if err != nil {
		return uuid.UUID{}, fmt.Errorf("invalid uuid: %w", err)
	}
	return id, nil
}
