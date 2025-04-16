package create

import (
	"encoding/json"
	"io"
	"josephwest2/meal-list/assert"
	"josephwest2/meal-list/lib/app"
	"josephwest2/meal-list/lib/auth"
	"josephwest2/meal-list/lib/db"
	"josephwest2/meal-list/pages"
	"net/http"
	"os"
	"strconv"
)

func Get(context *app.AppContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		assert.Assert(auth.IsAuthorized(context.DB, r, auth.AdminRole), "Admin role required")
		categories, err := context.DB.GetAllCategories()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		units, err := context.DB.GetAllUnits()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		ingredients, err := context.DB.GetAllIngredients()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		pages.RenderPage("Recipes", createRecipe(categories, units, ingredients), nil, w, r)
	}
}

type IngredientInput struct {
	ID       uint    `json:"id"`
	Quantity float64 `json:"quantity"`
	UnitID   int     `json:"unitid"`
}

func Post(context *app.AppContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		assert.Assert(auth.IsAuthorized(context.DB, r, auth.AdminRole), "Admin role required")
		r.ParseMultipartForm(5 << 20)

		name := r.FormValue("name")

		directions := r.FormValue("directions")

		category := r.FormValue("category")
		categoryInt, err := strconv.ParseUint(category, 10, 64)
		assert.Assert(err == nil, "failed to convert category to uint")

		sourceUrl := r.FormValue("recipe-source-url")

		ingredients := r.MultipartForm.Value["ingredient[]"]
		ingredientQuantities := make([]db.IngredientQuantity, 0, len(ingredients))
		for _, ingredientString := range ingredients {
			println(ingredientString)
			var ingredientParsed IngredientInput
			err := json.Unmarshal([]byte(ingredientString), &ingredientParsed)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				println("failed to parse ingredient json: " + err.Error())
				return
			}
			dbIngredient, err := context.DB.GetIngredientByID(ingredientParsed.ID)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				println("failed to get ingredient from db: " + err.Error())
				return
			}
			var dbUnit *db.Unit
			if ingredientParsed.UnitID >= 0 {
				dbUnit, err = context.DB.GetUnitByID(uint(ingredientParsed.UnitID))
				if err != nil {
					w.WriteHeader(http.StatusBadRequest)
					println("failed to get unit from db: " + err.Error())
					return
				}
			}
			ingredientQuantities = append(ingredientQuantities, db.IngredientQuantity{
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

		creationPath := "./static/recipe-images/" + handler.Filename
		_, err = os.Stat(creationPath)
		if !os.IsNotExist(err) {
			w.WriteHeader(http.StatusConflict)
			w.Write([]byte("Image: " + handler.Filename + " already exists"))
			return
		}

		dst, err := os.Create("./static/recipe-images/" + handler.Filename)
		assert.Assert(err != nil)
		defer dst.Close()
		io.Copy(dst, imageFile)

		// recipe insertion into db
		recipe := db.Recipe{
			Name:                 name,
			Directions:           directions,
			RecipeSourceURL:      sourceUrl,
			ImageName:            handler.Filename,
			RecipeCategoryID:     uint(categoryInt),
			IngredientQuantities: ingredientQuantities,
		}
		err = context.DB.CreateRecipe(recipe)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
}
