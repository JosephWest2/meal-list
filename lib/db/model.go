package db

import (
	"gorm.io/gorm"
)

func migrate(db *gorm.DB) {
	db.AutoMigrate(
		&IngredientCategory{},
		&Ingredient{},
		&Recipe{},
		&IngredientQuantity{},
		&Unit{},
	)
}


// Fields are converted to snake case in postgres

type Recipe struct {
	gorm.Model
	Name                 string               `gorm:"unique"`
	IngredientQuantities []IngredientQuantity `gorm:"many2many:recipe_ingredients;"`
	Directions           string
	RecipeCategoryID     uint
	RecipeCategory       RecipeCategory
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
	Weight UnitCategory = iota
	Volume
	Count
	Other
)

type Unit struct {
	ID               uint
	Name             string `gorm:"unique"`
	UnitCategory     uint   `gorm:"unique"`
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
