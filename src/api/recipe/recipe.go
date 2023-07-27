package recipe

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"test-api/controllers/recipe"
)

var Handler = http.NewServeMux()

func init() {
	Handler.HandleFunc("/create", handleCreateRecipe)
	Handler.HandleFunc("/get-list", handleGetRecipes)
	Handler.HandleFunc("/get-one-detailed", handleGetRecipeDetailed)
	Handler.HandleFunc("/update", handleUpdateRecipe)
	Handler.HandleFunc("/delete", handleRemoveRecipe)
}

func handleCreateRecipe(w http.ResponseWriter, r *http.Request) {
	var (
		reqBody recipe.RecipeAddRequest
		err     error
	)

	dec := json.NewDecoder(r.Body)

	err = dec.Decode(&reqBody)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, `{ "error": "bad json provided" }`)

		return
	}

	if reqBody.Title == "" {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, `{ "error": "recipe title is empty" }`)

		return
	}

	if len(reqBody.Ingredients) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, `{ "error": "recipe has no ingredients" }`)

		return
	}

	if len(reqBody.Steps) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, `{ "error": "recipe has no steps" }`)

		return
	}

	recipeId, err := recipe.CreateRecipe(reqBody)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, `{ "error": "internal server error" }`)

		return
	}

	w.WriteHeader(http.StatusCreated)
	fmt.Fprint(w, recipeId)
}

func handleGetRecipes(w http.ResponseWriter, r *http.Request) {
	var resp struct {
		Recipes []recipe.RecipeGeneral `json:"recipes"`
	}

	recipes, err := recipe.GetRecipes()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, `{ "error": "internal server error" }`)

		return
	}

	resp.Recipes = recipes

	respBytes, err := json.Marshal(resp)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, `{ "error": "internal server error" }`)

		return
	}

	fmt.Fprint(w, string(respBytes))
}

func handleGetRecipeDetailed(w http.ResponseWriter, r *http.Request) {
	var (
		recipeIdStr = r.URL.Query().Get("id")
		resp        struct {
			Recipe recipe.Recipe `json:"recipe"`
		}
	)

	recipeId, err := strconv.Atoi(recipeIdStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, `{ "error": "invalid recipe id" }`)

		return
	}

	recipeDetailed, err := recipe.GetRecipe(recipeId)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, `{ "error": "internal server error" }`)

		return
	}

	resp.Recipe = *recipeDetailed

	respBytes, err := json.Marshal(resp)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, `{ "error": "internal server error" }`)

		return
	}

	fmt.Fprint(w, string(respBytes))
}

func handleUpdateRecipe(w http.ResponseWriter, r *http.Request) {

}

func handleRemoveRecipe(w http.ResponseWriter, r *http.Request) {

}
