package recipe

import (
	"database/sql"
	"errors"
	"test-api/db"

	"github.com/lib/pq"
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

func (step *RecipeStep) less(other *RecipeStep) bool {
	return step.Number < other.Number
}

func (step *RecipeStep) greater(other *RecipeStep) bool {
	return step.Number > other.Number
}

func quickSortSteps(steps []RecipeStep) []RecipeStep {
	var result []RecipeStep

	if len(steps) <= 1 {
		return steps
	}

	pivot := &steps[len(steps)/2]
	left := make([]RecipeStep, 0)
	right := make([]RecipeStep, 0)
	equal := make([]RecipeStep, 0)

	for _, step := range steps {
		if step.less(pivot) {
			left = append(left, step)
		} else if step.greater(pivot) {
			right = append(right, step)
		} else {
			equal = append(equal, step)
		}
	}

	left = quickSortSteps(left)
	right = quickSortSteps(right)

	result = append(append(left, equal...), right...)
	return result
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
		INSERT INTO recipes_ingredients(recipe_id, name)
		VALUES ($1, $2)
		ON CONFLICT DO NOTHING
	`

	stmt, err := db.Conn.Prepare(sqlStatement)
	if err != nil {
		return 0, err
	}
	defer stmt.Close()

	for _, ingredient := range recipe.Ingredients {
		_, err := stmt.Exec(recipeId, ingredient)
		if err != nil {
			return 0, err
		}
	}

	sqlStatement = `
		INSERT INTO recipe_steps(recipe_id, number, step_text, created_at, updated_at)
		VALUES ($1, $2, $3, now(), now())
	`

	stmt, err = db.Conn.Prepare(sqlStatement)
	if err != nil {
		return 0, err
	}
	defer stmt.Close()

	for _, step := range recipe.Steps {
		_, err = stmt.Exec(recipeId, step.Number, step.Text)
		if err != nil {
			return 0, err
		}
	}

	return recipeId, nil
}

func GetRecipes() ([]RecipeGeneral, error) {
	result := make([]RecipeGeneral, 0)

	// CREATE TABLE public.recipes (
	// 	id serial4 NOT NULL UNIQUE,
	// 	title VARCHAR(255) NOT NULL,
	// 	description TEXT NULL,
	// 	created_at TIMESTAMP NOT NULL,
	// 	updated_at TIMESTAMP NOT NULL
	// );

	sqlStatement := `
		SELECT id, title, description
		FROM recipes
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
	defer rows.Close()

	return result, nil
}

func GetRecipe(recipeId int) (*Recipe, error) {
	var (
		result = Recipe{}
		err    error
	)

	// TODO(Nikita): Try to parse it later via 1 request to db, but not that important for a moment.
	// sqlStatement := `
	// 	SELECT
	// 		r.id,
	// 		r.title,
	// 		r.description,
	// 		rs.id,
	// 		rs.number,
	// 		rs.step_text,
	// 		ri.name
	// 	FROM
	// 		recipes r
	// 	LEFT JOIN
	// 		recipe_steps rs ON r.id = rs.recipe_id
	// 	LEFT JOIN
	// 		recipes_ingredients ri ON r.id = ri.recipe_id
	// 	WHERE
	// 		r.id = $1
	// 	ORDER BY
	// 		rs.number
	// `

	sqlStatement := ` SELECT id, title, description FROM recipes WHERE id = $1`

	rows, err := db.Conn.Query(sqlStatement, recipeId)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("Recipe with specified id doesn't exist")
		}

		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(&result.Id, &result.Title, &result.Description)
		if err != nil {
			return nil, err
		}
	}

	sqlStatement = ` SELECT id, number, step_text FROM recipe_steps WHERE recipe_id = $1 `
	rows, err = db.Conn.Query(sqlStatement, recipeId)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("Recipe with specified id doesn't exist")
		}

		return nil, err
	}

	steps := make([]RecipeStep, 0)
	for rows.Next() {
		step := RecipeStep{}

		err := rows.Scan(&step.Id, &step.Number, &step.Text)
		if err != nil {
			return nil, err
		}

		steps = append(steps, step)
	}

	result.Steps = steps

	sqlStatement = ` SELECT name FROM recipes_ingredients WHERE recipe_id = $1 `
	rows, err = db.Conn.Query(sqlStatement, recipeId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	ingredients := make([]string, 0)
	for rows.Next() {
		var ingredient string

		err := rows.Scan(&ingredient)
		if err != nil {
			return nil, err
		}

		ingredients = append(ingredients, ingredient)
	}
	result.Ingredients = ingredients

	return &result, nil
}

type EditRecipeRequest struct {
	Id          int          `json:"id"`
	Title       string       `json:"title"`
	Description string       `json:"description"`
	Steps       []RecipeStep `json:"steps"`
	Ingredients []string     `json:"ingredients"`
}

func EditRecipeGeneral(editRecipeReq EditRecipeRequest) error {
	sqlStatement := `
		UPDATE recipes
		SET title = $1, description = $2, updated_at = now()
		WHERE id = $3
	`

	_, err := db.Conn.Exec(sqlStatement, editRecipeReq.Title, editRecipeReq.Description, editRecipeReq.Id)
	if err != nil {
		return err
	}

	return nil
}

func EditRecipeSteps(editRecipeReq EditRecipeRequest) error {
	sortedSteps := quickSortSteps(editRecipeReq.Steps)

	for i := 1; i <= len(sortedSteps); i++ {
		sortedSteps[i-1].Number = i
	}

	sqlStatement := `
		DELETE FROM recipe_steps *
		WHERE recipe_id = $1
	`

	_, err := db.Conn.Exec(sqlStatement, editRecipeReq.Id)
	if err != nil {
		return err
	}

	sqlStatement = `
		INSERT INTO recipe_steps(recipe_id, number, step_text, created_at, updated_at)
		VALUES ($1, $2, $3, now(), now())
	`

	stmt, err := db.Conn.Prepare(sqlStatement)
	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, step := range sortedSteps {
		_, err := stmt.Exec(editRecipeReq.Id, step.Number, step.Text)
		if err != nil {
			return err
		}
	}

	return nil
}

func EditRecipeIngredients(editRecipeReq EditRecipeRequest) error {
	sqlStatement := `
		DELETE FROM recipes_ingredients *
		WHERE recipe_id = $1
	`

	_, err := db.Conn.Exec(sqlStatement, editRecipeReq.Id)
	if err != nil {
		return err
	}

	sqlStatement = `
		INSERT INTO recipes_ingredients (recipe_id, name)
		SELECT $1, name FROM unnest($2::VARCHAR[]) name
		ON CONFLICT (recipe_id, name)
		DO NOTHING
	`

	_, err = db.Conn.Exec(sqlStatement, editRecipeReq.Id, pq.Array(editRecipeReq.Ingredients))
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
