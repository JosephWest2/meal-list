package db

import (
	"time"

	"gorm.io/gorm"
)

func migrate(db *gorm.DB) {
	db.AutoMigrate(
		&IngredientCategory{},
		&Ingredient{},
		&RecipeIngredient{},
		&Recipe{},
		&RecipeCategory{},
		&Unit{},
		&User{},
		&Session{},
		&List{},
		&ListIngredient{},
		&CustomListItem{},
	)
}

// Fields are converted to snake case in postgres

type Recipe struct {
	ID               uint
	CreatedAt        time.Time
	UpdatedAt        time.Time
	Name             string `gorm:"unique"`
	Ingredients      []RecipeIngredient
	Directions       string
	RecipeCategoryID uint
	RecipeCategory   RecipeCategory
	RecipeImage      string
	RecipeSourceURL  string
}

type RecipeCategory struct {
	ID   uint
	Name string `gorm:"unique"`
}

type RecipeIngredient struct {
	ID           uint
	RecipeID     uint
	Recipe       Recipe
	IngredientID uint
	Ingredient   Ingredient
	Quantity     float64
	UnitID       *uint
	Unit         *Unit
}

type UnitCategory string

const (
	// Kg Unit
	Mass UnitCategory = "Mass"
	// L Unit
	Volume UnitCategory = "Volume"
	// unitless
	Count UnitCategory = "Count"
	// unitless
	Other UnitCategory = "Other"
)

type Unit struct {
	ID               uint
	Name             string `gorm:"unique"`
	UnitCategory     UnitCategory
	ConversionFactor float64
}

type Ingredient struct {
	ID         uint               `json:"id"`
	Name       string             `gorm:"unique" json:"name"`
	CategoryID uint               `json:"categoryID"`
	Category   IngredientCategory `json:"category"`
}

type IngredientCategory struct {
	ID   uint
	Name string `gorm:"unique"`
}

type Role uint

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
	ID          uint
	Ingredients []ListIngredient
	CustomListItems []CustomListItem
	UserID      uint
	User        User
	CreatedAt   time.Time
}

type ListIngredient struct {
	ID     uint
	ListID uint
	List   List

	IngredientID uint
	Ingredient   Ingredient

	UnitID *uint
	Unit   *Unit

	Quantity float64

	RecipeID *uint
	Recipe   *Recipe
}

type CustomListItem struct {
	ID uint
	ListID uint
	List List

	Name string
	Amount string
}