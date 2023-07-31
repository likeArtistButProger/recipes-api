package recipe

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"test-api/controllers/recipe"
)

var Handler = http.NewServeMux()

func init() {
	Handler.HandleFunc("/recipe/create", handleCreateRecipe)
	Handler.HandleFunc("/recipe/get-list", handleGetRecipes)
	Handler.HandleFunc("/recipe/get-one-detailed", handleGetRecipeDetailed)
	Handler.HandleFunc("/recipe/update", handleUpdateRecipe)
	Handler.HandleFunc("/recipe/delete", handleRemoveRecipe)
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
		log.Println(err)
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
	var reqBody recipe.EditRecipeRequest

	dec := json.NewDecoder(r.Body)

	err := dec.Decode(&reqBody)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, `{ "error": "invalid req body, make sure its structured as a recipe" }`)

		return
	}

	err = recipe.EditRecipeGeneral(reqBody)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, `{ "error": "couldn't change general info of recipe" }`)

		return
	}

	err = recipe.EditRecipeSteps(reqBody)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, `{ "error": "couldn't change steps of recipe" }`)

		return
	}

	err = recipe.EditRecipeIngredients(reqBody)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, `{ "error": "couldn't change ingredients of recipe" }`)

		return
	}
}

func handleRemoveRecipe(w http.ResponseWriter, r *http.Request) {
	var reqBody struct {
		Id int `json:"id"`
	}

	dec := json.NewDecoder(r.Body)
	err := dec.Decode(&reqBody)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, `{ "error": "bad recipe id provided" }`)

		return
	}

	err = recipe.RemoveRecipe(reqBody.Id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, `{ "error": "internal server error" }`)

		return
	}
}
