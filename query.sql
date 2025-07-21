-- name: DeleteIngredient :exec
delete from ingredients where id = $1;

-- name: CreateIngredient :one
insert into ingredients (name, ingredient_category_id) values ($1, $2) returning *;

-- name: UpdateIngredient :one
update ingredients set name = $2, ingredient_category_id = $3 where id = $1 returning *;

-- name: GetRecipe :one
select * from recipes where id = $1;

-- name: GetRecipeAndAssociatedData :one
select * from recipes
join recipes_to_recipe_categories on recipes.id = recipes_to_recipe_categories.recipe_id
join recipe_categories on recipes_to_recipe_categories.recipe_category_id = recipe_categories.id
join recipes_to_ingredients on recipes.id = recipes_to_ingredients.recipe_id
join ingredients on recipes_to_ingredients.ingredient_id = ingredients.id
join units on ingredients.unit_id = units.id
where recipes.id = $1;

-- name: GetAllRecipeCategories :many
select * from recipe_categories;

-- name: GetRecipesWithOffset :many
select * from recipes offset $1 limit $2;

-- name: GetAllIngredientCategories :many
select * from ingredient_categories;

-- name: GetIngredientsWithOffset :many
select * from ingredients offset $1 limit $2;

-- name: GetAllUnits :many
select * from units;

-- name: GetRecipesByNameAndCategoryID :many
select * from recipes_to_recipe_categories
join recipes on recipes_to_recipe_categories.recipe_id = recipes.id 
where recipes.name = $1 and recipes_to_recipe_categories.recipe_category_id = $2;

-- name: DeleteRecipeCategory :exec
delete from recipe_categories where id = $1;

-- name: CreateRecipeCategory :one
insert into recipe_categories (name) values ($1) returning *;

-- name: UpdateRecipeCategory :one
update recipe_categories set name = $2 where id = $1 returning *;

-- name: CreateUserSession :one
insert into user_sessions (id, user_id, last_access) values ($1, $2, now()) returning *;

-- name: GetUserSessionBySessionID :one
select * from user_sessions where id = $1;

-- name: DeleteUserSessionBySessionID :exec
delete from user_sessions where id = $1;

-- name: DeleteUserSessionByUserID :exec
delete from user_sessions where user_id = $1;

-- name: GetUserByUsername :one
select * from users where username = $1;

-- name: GetUserByID :one
select * from users where id = $1;

-- name: CreateUser :one
insert into users (username, password_hash, role) values ($1, $2, $3) returning *;

-- name: DeleteUser :exec
delete from users where id = $1;