package postgres_test

import (
	"context"
	"testing"

	"github.com/shanvl/garbage/internal/eventsvc/postgres"
)

func TestValidateSchema(t *testing.T) {
	t.Run("ok case", func(t *testing.T) {
		if err := postgres.ValidateSchema(context.Background(), db); err != nil {
			t.Errorf("ValidateSchema() error = %v", err)
		}
	})
}
