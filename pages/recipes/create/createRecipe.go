package create

import (
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

func Post(context *app.AppContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
        assert.Assert(auth.IsAuthorized(context.DB, r, auth.AdminRole), "Admin role required")
        r.ParseMultipartForm(5 << 20)

        recipe := db.Recipe{}

        // Name
        name := r.FormValue("name")
        recipe.Name = name

        directions := r.FormValue("directions")
        recipe.Directions = directions

        category := r.FormValue("category")
        categoryInt, err := strconv.ParseUint(category, 10, 64)
        assert.Assert(err == nil, "failed to convert category to uint")
        recipe.RecipeCategoryID = uint(categoryInt)

        sourceUrl := r.FormValue("recipe-source-url")
        recipe.RecipeSourceURL = sourceUrl

        ingredients := r.MultipartForm.Value["ingredient[]"]
        for _, ingredient := range ingredients {
            println(ingredient)
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
        if err != nil {
            panic(err)
        }
        defer dst.Close()
        io.Copy(dst, imageFile)
	}
}
