package app

import (
	"context"
	"josephwest2/meal-list/lib/sqlc"
)

type AppContext struct {
	Queries      *sqlc.Queries
	QueryContext context.Context
}
