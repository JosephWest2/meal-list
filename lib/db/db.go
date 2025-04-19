package db

import (
	"errors"
	"josephwest2/meal-list/assert"
	"log"
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type DB struct {
	gormdb *gorm.DB
}

func InitDB(connectionString string) DB {
	db, err := gorm.Open(postgres.Open(connectionString), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to db: ", err)
	}
	migrate(db)
	return DB{gormdb: db}
}

type RecipeQueryParams struct {
	Category string
	Name     string
}

func (db *DB) GetRecipeByID(id uint) (*Recipe, error) {
	var recipe Recipe
	err := db.gormdb.Model(&Recipe{}).Where(&Recipe{ID: id}).Preload(clause.Associations).Preload("Ingredients.Unit").Preload("Ingredients.Ingredient").First(&recipe).Error
	if err != nil {
		return nil, err
	}
	return &recipe, nil
}

func (db *DB) GetAllCategories() ([]RecipeCategory, error) {
	var categories []RecipeCategory
	err := db.gormdb.Find(&categories).Error
	if err != nil {
		return nil, err
	}
	return categories, nil
}

func (db *DB) GetAllIngredients() ([]Ingredient, error) {
	var ingredients []Ingredient
	err := db.gormdb.Find(&ingredients).Error
	if err != nil {
		return nil, err
	}
	return ingredients, nil
}

func (db *DB) GetAllUnits() ([]Unit, error) {
	var units []Unit
	err := db.gormdb.Find(&units).Error
	if err != nil {
		return nil, err
	}
	return units, nil
}

func (db *DB) GetRecipes(recipeQueryParams *RecipeQueryParams) []Recipe {
	var recipes []Recipe
	var query = db.gormdb.Model(&Recipe{}).Limit(10)
	if recipeQueryParams.Category != "" {
		var categoryId uint
		db.gormdb.Model(&RecipeCategory{}).Where("name LIKE ?", recipeQueryParams.Category).Find(&categoryId)
		query = query.Where("recipe_category_id = ?", recipeQueryParams.Category)
	}
	if recipeQueryParams.Name != "" {
		query = query.Where("name LIKE ?", "%"+recipeQueryParams.Name+"%")
	}
	query.Joins("RecipeCategory").Find(&recipes)
	return recipes
}

func (db *DB) GetSessionByID(sessionID string) (*Session, error) {
	var session Session
	err := db.gormdb.Model(&Session{}).Where("ID = ?", sessionID).First(&session).Error
	if err != nil {
		return nil, err
	}
	return &session, nil
}

func (db *DB) ClearSession(sessionID string) error {
	return db.gormdb.Model(&Session{}).Where("ID = ?", sessionID).Delete(&Session{}).Error
}

func (db *DB) CreateSession(sessionID string, userID uint) error {
	return db.gormdb.Create(&Session{ID: sessionID, UserID: userID, CreatedAt: time.Now()}).Error
}

func (db *DB) GetUserByUsername(username string) (*User, error) {
	var user User
	err := db.gormdb.Model(&User{}).Where("username = ?", username).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (db *DB) GetUserByID(id uint) (*User, error) {
	var user User
	err := db.gormdb.Model(&User{}).Where("id = ?", id).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (db *DB) GetIngredientByID(id uint) (*Ingredient, error) {
	var ingredient Ingredient
	err := db.gormdb.Model(&Ingredient{}).Where("id = ?", id).First(&ingredient).Error
	if err != nil {
		return nil, err
	}
	return &ingredient, nil
}

func (db *DB) GetUnitByID(id uint) (*Unit, error) {
	var unit Unit
	err := db.gormdb.Model(&Unit{}).Where("id = ?", id).First(&unit).Error
	if err != nil {
		return nil, err
	}
	return &unit, nil
}

func (db *DB) CreateUser(username string, passwordUnhashed string, role Role) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(passwordUnhashed), bcrypt.DefaultCost)
	if err != nil {
		return errors.New("failed to create user")
	}
	err = db.gormdb.Create(&User{Username: username, PasswordHash: string(hash), Role: role}).Error
	if err != nil {
		return err
	}
	return nil
}

func (db *DB) createRecipeIngredients(recipeIngredients []RecipeIngredient) error {
	err := db.gormdb.Create(recipeIngredients).Error
	if err != nil {
		return err
	}
	return nil
}

func (db *DB) CreateRecipe(recipe Recipe) error {
	err := db.gormdb.Create(&recipe).Error
	if err != nil {
		return err
	}
	return nil
}

func (db *DB) UpdateUserRole(username string, role Role) error {
	return db.gormdb.Model(&User{}).Where("username = ?", username).Update("role", role).Error
}

type ListItemDetails struct {
	Name               string
	Quantity           float64
	Unit               string
	UnitCategory       string
	ConversionFactor   float64
	IngredientCategory string
}

func (db *DB) AddToList(listID uint, listItemDetails ListItemDetails) error {
	var listIngredient ListIngredient
	var ingredient Ingredient
	err := db.gormdb.Model(&Ingredient{}).Where("name = ?", listItemDetails.Name).First(&ingredient).Error
	if err != nil {
		return err
	}
	listIngredient.IngredientID = ingredient.ID
	listIngredient.Quantity = listItemDetails.Quantity
	var unit *Unit
	db.gormdb.Model(&Unit{}).Where("name = ?", listItemDetails.Unit).First(unit)
	listIngredient.Unit = unit
	return db.gormdb.Model(&List{ID: listID}).Association("Ingredients").Append(&listIngredient)
}

func (db *DB) GetRecipeIngredients(recipeID uint) ([]RecipeIngredient, error) {
	var ingredients []RecipeIngredient
	err := db.gormdb.Model(&RecipeIngredient{}).Where("recipe_id = ?", recipeID).Preload("Unit").Preload("Ingredient").Find(&ingredients).Error
	if err != nil {
		return nil, err
	}
	return ingredients, nil
}

func (db *DB) AddRecipeToListOrIncrement(listID uint, recipeID uint) error {
	var currentListIngredients []ListIngredient
	db.gormdb.Model(&ListIngredient{}).Where("list_id = ? AND recipe_id = ?", listID, recipeID).Find(&currentListIngredients)
	currentListIngredientsMap := make(map[uint]ListIngredient)
	for _, listIngredient := range currentListIngredients {
		currentListIngredientsMap[listIngredient.IngredientID] = listIngredient
	}
	recipeIngredients, err := db.GetRecipeIngredients(recipeID)
	if err != nil {
		return err
	}
	for _, recipeIngredient := range recipeIngredients {
		_, contains := currentListIngredientsMap[recipeIngredient.IngredientID]
		if contains {
			temp := currentListIngredientsMap[recipeIngredient.IngredientID]
			temp.Quantity += recipeIngredient.Quantity
			err = db.gormdb.Model(&ListIngredient{}).Where("ID = ?", temp.ID).Updates(temp).Error
			assert.Assert(err == nil, "Failed to update list item")
		} else {
			listIngredient := ListIngredient{
				ListID:       listID,
				IngredientID: recipeIngredient.Ingredient.ID,
				Quantity:     recipeIngredient.Quantity,
				Unit:         recipeIngredient.Unit,
				RecipeID:     &recipeID,
			}
			err = db.gormdb.Create(&listIngredient).Error
			assert.Assert(err == nil, "Failed to create list item")
		}
	}
	return nil
}

func (db *DB) RemoveFromList(listID uint, listIngredientID uint) error {
	return db.gormdb.Unscoped().Model(&List{ID: listID}).Association("Ingredients").Unscoped().Delete(&ListIngredient{ID: listIngredientID})
}

func (db *DB) GetOrCreateListByUserID(userID uint) List {
	var list List
	err := db.gormdb.Model(&List{UserID: userID}).Preload("Ingredients").Preload("Ingredients.Ingredient").Preload("Ingredients.Unit").First(&list).Error
	if err != nil {
		list = List{UserID: userID}
		db.gormdb.Create(&list)
	}
	return list
}

func (db *DB) GetListByID(listID uint) (*List, error) {
	var list List
	err := db.gormdb.Model(&List{ID: listID}).Preload("Ingredients").First(&list).Error
	if err != nil {
		return nil, err
	}
	return &list, nil
}
