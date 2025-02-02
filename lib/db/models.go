package db

import (
	"time"

	"gorm.io/gorm"
)

func migrate(db *gorm.DB) {
	db.AutoMigrate(
		&IngredientCategory{},
		&Ingredient{},
		&Recipe{},
		&RecipeCategory{},
		&IngredientQuantity{},
		&Unit{},
		&User{},
		&Session{},
		&List{},
		&ListItem{},
	)
}

// Fields are converted to snake case in postgres

type Recipe struct {
	ID                   uint
	CreatedAt            time.Time
	UpdatedAt            time.Time
	Name                 string               `gorm:"unique"`
	IngredientQuantities []IngredientQuantity `gorm:"many2many:recipe_ingredients;"`
	Directions           string
	RecipeCategoryID     uint
	RecipeCategory       RecipeCategory
	ImageRef             string
	RecipeSourceURL      string
}

type RecipeCategory struct {
	ID   uint
	Name string `gorm:"unique"`
}

type IngredientQuantity struct {
	ID           uint
	IngredientID uint
	Ingredient   Ingredient
	Quantity     float64
	UnitID       uint
	Unit         Unit
}

type UnitCategory uint

const (
	// Kg Unit
	Mass UnitCategory = iota
	// L Unit
	Volume
	// unitless
	Count
	// unitless
	Other
)

type Unit struct {
	ID               uint
	Name             string `gorm:"unique"`
	UnitCategory     UnitCategory
	ConversionFactor float64
}

type Ingredient struct {
	ID         uint
	Name       string `gorm:"unique"`
	CategoryID uint
	Category   IngredientCategory
}

type IngredientCategory struct {
	ID   uint
	Name string `gorm:"unique"`
}

type Role uint

const (
	StandardUser Role = iota
	Admin
)

type User struct {
	ID           uint
	CreatedAt    time.Time
	UpdatedAt    time.Time
	Username     string `gorm:"unique"`
	Role         Role
	PasswordHash string
}

type Session struct {
	ID        string `gorm:"primarykey"`
	UserID    uint
	User      User
	CreatedAt time.Time
}

type List struct {
	ID        uint
	ListItems []ListItem
	UserID    uint
	User      User
	CreatedAt time.Time
}

type ListItem struct {
	ID     uint
	ListID uint

	Name               string
	Quantity           float64
	Unit               string
	UnitCategory       string
	ConversionFactor   float64
	IngredientCategory string
}
