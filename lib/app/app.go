package app

import (
	"context"
	"josephwest2/meal-list/lib/sqlc"
)

type App struct {
	Queries      *sqlc.Queries
	QueryContext context.Context
}
