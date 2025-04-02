package recipes

import (
	"josephwest2/meal-list/lib/app"
	"josephwest2/meal-list/lib/auth"
	"josephwest2/meal-list/lib/db"
	"josephwest2/meal-list/pages"
	"net/http"
	"os"
)

func Get(context *app.AppContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

        dbRecipes := context.DB.GetRecipes(&db.RecipeQueryParams{})
        isAdmin := auth.IsAuthorized(context.DB, r, auth.AdminRole);
        var categories []db.RecipeCategory
        var imageNames []string
        if isAdmin {
            categories, _ = context.DB.GetAllCategories()
            f, err := os.Open("./static/recipe-images/")
            if err != nil {
                panic("failed to read recipe-images directory")
            }
            imageNames, _ = f.Readdirnames(0)
        }
        pages.RenderPage("Recipes", recipes(dbRecipes, isAdmin, categories, imageNames), nil, w, r)
	}
}

func Post(context *app.AppContext) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        isAdmin := auth.IsAuthorized(context.DB, r, auth.AdminRole);
        if !isAdmin {
            w.WriteHeader(http.StatusUnauthorized)
            return
        }
        w.WriteHeader(http.StatusOK)
    }
}
