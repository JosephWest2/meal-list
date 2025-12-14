-- name: DeleteIngredient :exec
delete from ingredients where id = $1;

-- name: CreateIngredient :one
insert into ingredients (name, ingredient_category_id) values ($1, $2) returning *;

-- name: UpdateIngredient :one
update ingredients set name = $2, ingredient_category_id = $3 where id = $1 returning *;

-- name: GetRecipeByID :one
select * from recipes where id = $1;

-- name: GetRecipeAndAssociatedData :one
select * from recipes
join recipes_to_recipe_categories on recipes.id = recipes_to_recipe_categories.recipe_id
join recipe_categories on recipes_to_recipe_categories.recipe_category_id = recipe_categories.id
join recipes_to_ingredients on recipes.id = recipes_to_ingredients.recipe_id
join ingredients on recipes_to_ingredients.ingredient_id = ingredients.id
join units on ingredients.unit_id = units.id
where recipes.id = $1;

-- name: GetRecipeAssociatedCategories :many
select recipe_categories.* from recipe_categories
join recipes_to_recipe_categories on recipe_categories.id = recipes_to_recipe_categories.recipe_category_id
where recipes_to_recipe_categories.recipe_id = $1;

-- name: GetRecipeAssociatedIngredientsAndUnits :many
select * from ingredients
join recipes_to_ingredients on ingredients.id = recipes_to_ingredients.ingredient_id
join units on ingredients.unit_id = units.id
where recipes_to_ingredients.recipe_id = $1;

-- name: GetAllRecipeCategories :many
select * from recipe_categories;

-- name: GetRecipesWithOffset :many
select * from recipes offset $1 limit $2;

-- name: GetAllIngredients :many
select * from ingredients;

-- name: GetAllIngredientCategories :many
select * from ingredient_categories;

-- name: GetIngredientsWithOffset :many
select * from ingredients offset $1 limit $2;

-- name: GetIngredientByID :one
select * from ingredients where id = $1;

-- name: GetAllUnits :many
select * from units;

-- name: GetUnitByID :one
select * from units where id = $1;

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

-- name: GetListsByUserID :many
select * from lists where user_id = $1;

-- name: CreateList :one
insert into lists (name, user_id) values ($1, $2) returning *;

-- name: GetListByName :one
select * from lists where name = $1;

-- name: GetListByID :one
select * from lists where id = $1;

-- name: GetListRecipesAndFullRecipeIngredients :one
select * from lists
join lists_to_recipes on lists.id = lists_to_recipes.list_id
join recipes on lists_to_recipes.recipe_id = recipes.id
join recipes_to_ingredients on recipes.id = recipes_to_ingredients.recipe_id
join ingredients on recipes_to_ingredients.ingredient_id = ingredients.id
join units on ingredients.unit_id = units.id
where lists.id = $1;

-- name: GetFreeListItemsAndFullIngredients :one
select list_id from lists
join free_list_items on lists.id = free_list_items.list_id
join ingredients on free_list_items.ingredient_id = ingredients.id
join units on ingredients.unit_id = units.id
join custom_list_items on free_list_items.custom_list_item_id = custom_list_items.id
join units on custom_list_items.unit_id = units.id
where lists.id = $1;

-- name: GetListAssociatedRecipeByID :one
select recipes.* from lists_to_recipes
join recipes on lists_to_recipes.recipe_id = recipes.id
where lists_to_recipes.list_id = @list_id::int and recipes.id = @recipe_id::int;

-- name: GetAllListAssociatedRecipes :many
select recipes.*, lists_to_recipes.quantity from lists_to_recipes
join recipes on lists_to_recipes.recipe_id = recipes.id
where lists_to_recipes.list_id = @list_id::int;

-- name: GetListLeftJoinListToRecipe :one
select lists.*, lists_to_recipes.quantity from lists
left join lists_to_recipes on lists.id = lists_to_recipes.list_id
where lists.id = @list_id::int and lists_to_recipes.recipe_id = @recipe_id::int;

-- name: AddToListRecipeQuantity :exec
update lists_to_recipes set quantity = quantity + $3 where list_id = $1 and recipe_id = $2;

-- name: UpdateListToRecipeQuantity :exec
update lists_to_recipes set quantity = $3 where list_id = $1 and recipe_id = $2;

-- name: CreateListToRecipe :one
insert into lists_to_recipes (list_id, recipe_id, quantity) values ($1, $2, $3) returning *;

-- name: GetListAssociatedRecipesAndQuantities :many
select * from lists_to_recipes
join recipes on lists_to_recipes.recipe_id = recipes.id
where lists_to_recipes.list_id = $1;

-- name: CreateRecipe :one
insert into recipes (name, directions, source_url, image_filename) values ($1, $2, $3, $4) returning *;
