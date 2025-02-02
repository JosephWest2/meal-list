package db

import (
	"errors"
	"fmt"
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
	err := db.gormdb.Model(&Recipe{}).Where(&Recipe{ID: id}).Preload(clause.Associations).Preload("IngredientQuantities.Unit").Preload("IngredientQuantities.Ingredient").First(&recipe).Error
	if err != nil {
		return nil, err
	}
	return &recipe, nil
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
	listItem := ListItem{
		ListID:             listID,
		Name:               listItemDetails.Name,
		Quantity:           listItemDetails.Quantity,
		Unit:               listItemDetails.Unit,
		UnitCategory:       listItemDetails.UnitCategory,
		ConversionFactor:   listItemDetails.ConversionFactor,
		IngredientCategory: listItemDetails.IngredientCategory,
	}
	return db.gormdb.Model(&List{ID: listID}).Association("ListItems").Append(&listItem)
}

func (db *DB) RemoveFromList(listID uint, listItemID uint) error {
	return db.gormdb.Unscoped().Model(&List{ID: listID}).Association("ListItems").Unscoped().Delete(&ListItem{ID: listItemID})
}

func (db *DB) GetListByUserID(userID uint) List {
	var list List
	err := db.gormdb.Model(&List{UserID: userID}).Preload("ListItems").First(&list).Error
	if err != nil {
		list = List{UserID: userID}
		db.gormdb.Create(&list)
	}
	fmt.Println(list)
	return list
}

func (db *DB) GetListByID(listID uint) (*List, error) {
	var list List
	err := db.gormdb.Model(&List{ID: listID}).Preload("ListItems").First(&list).Error
	if err != nil {
		return nil, err
	}
	return &list, nil
}
