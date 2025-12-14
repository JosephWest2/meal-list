package createRecipe

import (
	"encoding/json"
	"io"
	"josephwest2/meal-list/assert"
	"josephwest2/meal-list/lib/app"
	"josephwest2/meal-list/lib/auth"
	"josephwest2/meal-list/lib/db"
	"josephwest2/meal-list/lib/sqlc"
	"josephwest2/meal-list/pages"
	"net/http"
	"os"
	"strconv"
	"strings"
)

func Get(context *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		assert.Assert(auth.IsAuthorized(context, r, sqlc.RoleAdmin), "Admin role required")
		categories, err := context.Queries.GetAllRecipeCategories(context.QueryContext)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		units, err := context.Queries.GetAllUnits(context.QueryContext)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		ingredients, err := context.Queries.GetAllIngredients(context.QueryContext)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		pages.RenderPage(context, "Recipes", createRecipe(categories, units, ingredients), nil, w, r)
	}
}

type IngredientInput struct {
	ID       int32    `json:"id"`
	Quantity float64 `json:"quantity"`
	UnitID   int32     `json:"unitid"`
}

const imageDirectory = "./static/recipeImages/"

func Post(context *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		assert.Assert(auth.IsAuthorized(context, r, sqlc.RoleAdmin), "Admin role required")
		r.ParseMultipartForm(5 << 20)

		name := r.FormValue("name")

		directions := r.FormValue("directions")

		categoryIDString := r.FormValue("category-id")
		categoryID, err := strconv.ParseUint(categoryIDString, 10, 64)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		sourceUrl := r.FormValue("recipe-source-url")

		ingredients := r.MultipartForm.Value["ingredient[]"]
		recipeIngredients := make([]db.RecipeIngredient, 0, len(ingredients))
		for _, ingredientString := range ingredients {
			var ingredientParsed IngredientInput
			err := json.Unmarshal([]byte(ingredientString), &ingredientParsed)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				println("failed to parse ingredient json: " + err.Error())
				return
			}
			dbIngredient, err := context.Queries.GetIngredientByID(context.QueryContext, ingredientParsed.ID)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				println("failed to get ingredient from db: " + err.Error())
				return
			}
			var unit sqlc.Unit
			if ingredientParsed.UnitID >= 0 {
				unit, err = context.Queries.GetUnitByID(context.QueryContext, ingredientParsed.UnitID)
			}
			recipeIngredients = append(recipeIngredients, db.RecipeIngredient{
				Ingredient: *dbIngredient,
				Quantity:   ingredientParsed.Quantity,
				Unit:       dbUnit,
			})
		}

		// Image
		imageFile, handler, err := r.FormFile("image")
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		defer imageFile.Close()

		// replace whitespace with underscore, append file extension
		newFileName := strings.ReplaceAll(name, " ", "_") + "." + strings.Split(handler.Filename, ".")[1]
		creationPath := imageDirectory + newFileName
		_, err = os.Stat(creationPath)
		if os.IsExist(err) {
			w.WriteHeader(http.StatusConflict)
			w.Write([]byte("Image: " + newFileName + " already exists"))
			return
		}
		err = os.MkdirAll(imageDirectory, os.ModePerm)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		dst, err := os.Create(imageDirectory + newFileName)
		assert.Assert(err == nil)
		defer dst.Close()
		io.Copy(dst, imageFile)

		_, err := context.Queries.CreateRecipe(context.QueryContext, sqlc.CreateRecipeParams{
			Name:             recipe.Name,
			Directions:       recipe.Directions,
			SourceUrl:        recipe.RecipeSourceURL,
			ImageFilename:    recipe.RecipeImage,
		})
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
}
