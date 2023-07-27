package recipe

import (
	"database/sql"
	"errors"
	"test-api/db"
)

type RecipeStep struct {
	Id     int    `json:"id"`
	Number int    `json:"number"`
	Text   string `json:"text"`
}

type Recipe struct {
	Id          int          `json:"id"`
	Title       string       `json:"title"`
	Description string       `json:"description"`
	Ingredients []string     `json:"ingredients"`
	Steps       []RecipeStep `json:"steps"`
}

type RecipeGeneral struct {
	Id          int    `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
}

type RecipeAddRequest struct {
	Title       string       `json:"title"`
	Description string       `json:"description"`
	Ingredients []string     `json:"ingredients"`
	Steps       []RecipeStep `json:"steps"`
}

func CreateRecipe(recipe RecipeAddRequest) (int, error) {
	var (
		err      error
		recipeId int
	)

	// NOTE(Nikita): Add main recipe info
	sqlStatement := `
		INSERT INTO recipes(title, description, created_at, updated_at)
		VALUES ($1, $2, now(), now())
		RETURNING id
	`

	err = db.Conn.QueryRow(sqlStatement, recipe.Title, recipe.Description).Scan(&recipeId)
	if err != nil {
		return 0, err
	}

	// NOTE(Nikita): Add ingredients, it would be better to have storage of ingredients instead of storing them
	// 				 with recipe id to reduce duplicates, but its a story for another time
	sqlStatement = `
		INSERT INTO recipes_ingredients()
		VALUE ($1, $2)
		ON CONFLICT DO NOTHING
	`

	stmt, err := db.Conn.Prepare(sqlStatement)
	if err != nil {
		return 0, err
	}

	for _, ingredient := range recipe.Ingredients {
		_, err := stmt.Exec(recipeId, ingredient)
		if err != nil {
			return 0, err
		}
	}

	sqlStatement = `
		INSERT INTO recipe_steps(number, step_text, created_at, updated_at)
		VALUES ($1, $2, now(), now())
	`

	stmt, err = db.Conn.Prepare(sqlStatement)
	if err != nil {
		return 0, err
	}

	for _, step := range recipe.Steps {
		_, err = stmt.Exec(step.Number, step.Text)
		if err != nil {
			return 0, err
		}
	}

	return recipeId, nil
}

func GetRecipes() ([]RecipeGeneral, error) {
	result := make([]RecipeGeneral, 0)

	sqlStatement := `
		SELECT id, title, description
		FROM recipes(id, title, description)
	`

	rows, err := db.Conn.Query(sqlStatement)
	if err != nil {
		if err == sql.ErrNoRows {
			return []RecipeGeneral{}, nil
		}

		return nil, err
	}

	for rows.Next() {
		recipe := RecipeGeneral{}

		err = rows.Scan(&recipe.Id, &recipe.Title, &recipe.Description)
		if err != nil {
			return nil, err
		}

		result = append(result, recipe)
	}

	return result, nil
}

func GetRecipe(recipeId int) (*Recipe, error) {
	var (
		result Recipe
		err    error
	)

	sqlStatement := `
		SELECT 
			r.id
			r.title
			r.description
			rs.id,
			rs.number,
			rs.step_text,
			ri.name
		FROM
			recipes r
		LEFT JOIN
			recipe_steps rs ON r.id = rs.recipe_id
		LEFT JOIN
			recipes_ingredients ri ON r.id = ri.recipe_id
		WHERE
			r.id = $1
		ORDER BY
			rs.number
	`

	rows, err := db.Conn.Query(sqlStatement, recipeId)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("Recipe with specified id doesn't exist")
		}

		return nil, err
	}

	for rows.Next() {
		// _, err := rows.Scan()
	}

	return &result, nil
}

type EditRecipeRequest struct {
	Id          int    `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
}

func EditRecipeGeneral(editRecipeReq EditRecipeRequest) error {
	sqlStatement := `
		UPDATE recipes (title, description)
		SET title = $1, description = $2
		WHERE id = $3
	`

	_, err := db.Conn.Exec(sqlStatement, editRecipeReq.Title, editRecipeReq.Description, editRecipeReq.Id)
	if err != nil {
		return err
	}

	return nil
}

func RemoveRecipe(recipeId int) error {
	var err error

	sqlStatement := `
		DELETE FROM recipes * WHERE id = $1
	`

	_, err = db.Conn.Exec(sqlStatement, recipeId)
	if err != nil {
		return err
	}

	sqlStatement = `
		DELETE FROM recipe_steps * WHERE recipe_id = $1
	`

	_, err = db.Conn.Exec(sqlStatement, recipeId)
	if err != nil {
		return err
	}

	sqlStatement = `
		DELETE FROM recipes_ingredients * WHERE recipe_id = $1
	`

	_, err = db.Conn.Exec(sqlStatement, recipeId)
	if err != nil {
		return err
	}

	return nil
}
