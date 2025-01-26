package db

import (
	"fmt"
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var db *gorm.DB

func Setup() {
	postgresConnectionString := os.Getenv("POSTGRES_CONNECTION_STRING")
	fmt.Println("connectionString:", postgresConnectionString)
	var err error
	db, err = gorm.Open(postgres.Open(postgresConnectionString), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to db", err)
	}
	migrate(db)
}

type RecipeQueryParams struct {
	Category string
	Name     string
}

func GetRecipes(recipeQueryParams *RecipeQueryParams) []Recipe {
	var recipes []Recipe
	var query = db.Model(&Recipe{}).Limit(10)
	if recipeQueryParams.Category != "" {
		var categoryId uint
		db.Model(&RecipeCategory{}).Where("name LIKE ?", recipeQueryParams.Category).Find(&categoryId)
		query = query.Where("recipe_category_id = ?", recipeQueryParams.Category)
	}
	if recipeQueryParams.Name != "" {
		query = query.Where("name LIKE ?", "%"+recipeQueryParams.Name+"%")
	}
	query.Find(&recipes)
	return recipes
}

func SeedRecipes() {
	category := RecipeCategory{Name: "Breakfast"}
	db.Create(&category)

	recipes := []Recipe{
		{Name: "Pancakes", RecipeCategoryID: category.ID},
		{Name: "Oatmeal", RecipeCategoryID: category.ID},
	}
	db.Create(&recipes)
}