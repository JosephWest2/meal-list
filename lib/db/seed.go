package db

import (
	"fmt"
)

func (db *DB) Seed() {

	// Recipe Categories
	breakfastCategory := RecipeCategory{Name: "Breakfast"}
	db.gormdb.Create(&breakfastCategory)

	mealCategory := RecipeCategory{Name: "Meal"}
	db.gormdb.Create(&mealCategory)

	snackCategory := RecipeCategory{Name: "Snack"}
	db.gormdb.Create(&snackCategory)

	dinnerCategory := RecipeCategory{Name: "Dinner"}
	db.gormdb.Create(&dinnerCategory)

	// Ingredient Categories
	produceCategory := IngredientCategory{Name: "Produce"}
	db.gormdb.Create(&produceCategory)

	meatCategory := IngredientCategory{Name: "Meat"}
	db.gormdb.Create(&meatCategory)

	dairyCategory := IngredientCategory{Name: "Dairy"}
	db.gormdb.Create(&dairyCategory)

	frozenCategory := IngredientCategory{Name: "Frozen"}
	db.gormdb.Create(&frozenCategory)

	bakeryCategory := IngredientCategory{Name: "Bakery"}
	db.gormdb.Create(&bakeryCategory)

	cannedCategory := IngredientCategory{Name: "Canned"}
	db.gormdb.Create(&cannedCategory)

	deliCategory := IngredientCategory{Name: "Deli"}
	db.gormdb.Create(&deliCategory)

	pastaRiceCerealCategory := IngredientCategory{Name: "Pasta, Rice, Cereal"}
	db.gormdb.Create(&pastaRiceCerealCategory)

	spicesSeasoningCategory := IngredientCategory{Name: "Spices, Seasoning"}
	db.gormdb.Create(&spicesSeasoningCategory)

	// Ingredients
	egg := Ingredient{Name: "Egg", Category: produceCategory}
	db.gormdb.Create(&egg)

	milk := Ingredient{Name: "Milk", Category: dairyCategory}
	db.gormdb.Create(&milk)

	ricottaCheese := Ingredient{Name: "Ricotta Cheese", Category: dairyCategory}
	db.gormdb.Create(&ricottaCheese)

	garlic := Ingredient{Name: "Garlic", Category: produceCategory}
	db.gormdb.Create(&garlic)

	beef := Ingredient{Name: "Beef", Category: meatCategory}
	db.gormdb.Create(&beef)

	sausage := Ingredient{Name: "Sausage", Category: meatCategory}
	db.gormdb.Create(&sausage)

	parsley := Ingredient{Name: "Parsley", Category: produceCategory}
	db.gormdb.Create(&parsley)

	tomatoPaste := Ingredient{Name: "Tomato Paste", Category: cannedCategory}
	db.gormdb.Create(&tomatoPaste)

	lasagnaNoodles := Ingredient{Name: "Lasagna Noodles", Category: pastaRiceCerealCategory}
	db.gormdb.Create(&lasagnaNoodles)

	onion := Ingredient{Name: "Onion", Category: produceCategory}
	db.gormdb.Create(&onion)

	parmesanCheese := Ingredient{Name: "Parmesan Cheese", Category: dairyCategory}
	db.gormdb.Create(&parmesanCheese)

	italianSeasoning := Ingredient{Name: "Italian Seasoning", Category: spicesSeasoningCategory}
	db.gormdb.Create(&italianSeasoning)

	pastaSauce := Ingredient{Name: "Pasta Sauce", Category: cannedCategory}
	db.gormdb.Create(&pastaSauce)

	salt := Ingredient{Name: "Salt", Category: spicesSeasoningCategory}
	db.gormdb.Create(&salt)

	sugar := Ingredient{Name: "Sugar", Category: spicesSeasoningCategory}
	db.gormdb.Create(&sugar)

	flour := Ingredient{Name: "Flour", Category: bakeryCategory}
	db.gormdb.Create(&flour)

	bakingPowder := Ingredient{Name: "Baking Powder", Category: bakeryCategory}
	db.gormdb.Create(&bakingPowder)

	// Units
	kilogram := Unit{Name: "kilogram", UnitCategory: Mass, ConversionFactor: 1}
	db.gormdb.Create(&kilogram)

	liter := Unit{Name: "liter", UnitCategory: Volume, ConversionFactor: 1}
	db.gormdb.Create(&liter)

	gram := Unit{Name: "gram", UnitCategory: Mass, ConversionFactor: 0.001}
	db.gormdb.Create(&gram)

	pound := Unit{Name: "pound", UnitCategory: Mass, ConversionFactor: 0.453592}
	db.gormdb.Create(&pound)

	ounce := Unit{Name: "ounce", UnitCategory: Mass, ConversionFactor: 0.0283495}
	db.gormdb.Create(&ounce)

	cup := Unit{Name: "cup", UnitCategory: Volume, ConversionFactor: 0.236588}
	db.gormdb.Create(&cup)

	teaspoon := Unit{Name: "teaspoon", UnitCategory: Volume, ConversionFactor: 0.00492892}
	db.gormdb.Create(&teaspoon)

	tablespoon := Unit{Name: "tablespoon", UnitCategory: Volume, ConversionFactor: 0.0147868}
	db.gormdb.Create(&tablespoon)

	piece := Unit{Name: "piece", UnitCategory: Count, ConversionFactor: 1}
	db.gormdb.Create(&piece)

	// Recipes

	recipes := []Recipe{
		{
			Name: "Easy Homemade Lasagna",
			Directions: `
				1. Boil pasta: Cook lasagna noodles in salted water according to the recipe.
				2. Prepare meat sauce: Cook sausage, beef, onion, and garlic. Drain, then add pasta sauce and simmer to thicken.
				3. Combine cheese mixture: Mix cheese ingredients together in a bowl.
				4. Layer & bake: Layer the meat sauce, cheese mixture, and lasagna noodles, then bake until the top is golden brown.	
			`,
			RecipeSourceURL: "https://www.spendwithpennies.com/easy-homemade-lasagna/",
			ImageName:        "lasagna.jpg",
			RecipeCategory:  mealCategory,
			IngredientQuantities: []IngredientQuantity{
				{
					Ingredient: egg,
					Quantity:   1,
					Unit:       &piece,
				},
				{
					Ingredient: milk,
					Quantity:   1,
					Unit:       &cup,
				},
				{
					Ingredient: ricottaCheese,
					Quantity:   1,
					Unit:       &cup,
				},
				{
					Ingredient: garlic,
					Quantity:   1,
					Unit:       &piece,
				},
				{
					Ingredient: beef,
					Quantity:   1,
					Unit:       &pound,
				},
				{
					Ingredient: lasagnaNoodles,
					Quantity:   10,
					Unit:       &piece,
				},
				{
					Ingredient: tomatoPaste,
					Quantity:   1,
					Unit:       &cup,
				},
				{
					Ingredient: italianSeasoning,
					Quantity:   1,
					Unit:       &cup,
				},
			},
		},
		{
			Name: "Good Old Fashioned Pancakes",
			Directions: `
				1. Sift flour, baking powder, sugar, and salt into a large bowl.  
				2. Make a well in the center and add milk, melted butter, and egg; mix until smooth.  
				3. Heat a lightly oiled griddle or pan over medium-high heat.  
				4. Pour or scoop about 1/4 cup of batter onto the griddle for each pancake.  
				5. Cook until bubbles form and edges are dry (about 2–3 minutes), then flip.  
				6. Cook the other side until browned.  
				7. Repeat with remaining batter.	
			`,
			RecipeSourceURL: "https://www.allrecipes.com/recipe/21014/good-old-fashioned-pancakes/",
			ImageName:        "pancakes.jpg",
			RecipeCategory:  breakfastCategory,
			IngredientQuantities: []IngredientQuantity{
				{
					Ingredient: flour,
					Quantity:   1,
					Unit:       &cup,
				},
				{
					Ingredient: bakingPowder,
					Quantity:   1,
					Unit:       &teaspoon,
				},
				{
					Ingredient: sugar,
					Quantity:   1,
					Unit:       &cup,
				},
				{
					Ingredient: salt,
					Quantity:   1,
					Unit:       &teaspoon,
				},
				{
					Ingredient: milk,
					Quantity:   1,
					Unit:       &cup,
				},
			},
		},
	}
	db.gormdb.Create(&recipes)
	fmt.Println(recipes)

    db.CreateUser("admin", "admin", 1)
}
