package postgres_test

import (
	"context"
	"testing"

	"github.com/shanvl/garbage/internal/authsvc/postgres"
)

func TestValidateSchema(t *testing.T) {
	t.Parallel()
	t.Run("ok case", func(t *testing.T) {
		if err := postgres.ValidateSchema(context.Background(), db); err != nil {
			t.Errorf("ValidateSchema() error = %v", err)
		}
	})
}
