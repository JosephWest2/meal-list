package app

import (
	"josephwest2/meal-list/lib/db"
)

type AppContext struct {
	DB *db.DB
}
