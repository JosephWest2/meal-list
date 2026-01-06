create or replace function before_update_updated_at() returns trigger as
$body$
begin
    if row(new.*::text) is distinct from row(old.*::text) then
        new.updated_at = now();
    end if;

    return new;
end;
$body$
language plpgsql;

do $body$
declare t text;
begin
    for t in
        select table_name
        from information_schema.columns
        where (
            column_name = 'updated_at'
            and (
                select 1
                from information_schema.triggers
                where trigger_name = 'before_update_updated_at_' || table_name
            ) is null
        )
    loop
        execute format('
            create trigger before_update_updated_at_%s
            before update on %i
            for each row execute procedure before_update_updated_at();
        ', t, t);
    end loop;
end;
$body$
language plpgsql;

create or replace function create_recipe_category(name text) returns integer as $$
declare
    category_id integer;
begin
    insert into recipe_categories (recipe_category_name) values (name) returning recipe_category_id into category_id;
    return category_id;
end
$$ language plpgsql;

create or replace function create_ingredient_category(name text) returns integer as $$
declare
    category_id integer;
begin
    insert into ingredient_categories (ingredient_category_name) values (name) returning ingredient_category_id into category_id;
    return category_id;
end
$$ language plpgsql;

create or replace function create_unit(name text, category text, conversion_factor real) returns integer as $$
declare
    new_unit_id integer;
begin
    insert into units (unit_name, unit_category, conversion_factor) values (name, category::unit_category, conversion_factor) returning unit_id into new_unit_id;
    return new_unit_id;
end
$$ language plpgsql;

create or replace function create_ingredient(name text, category_id integer, unit_id integer) returns integer as $$
declare
    new_ingredient_id integer;
begin
    insert into ingredients (ingredient_name, unit_id) values (name, unit_id) returning ingredient_id into new_ingredient_id;
    insert into ingredients_to_categories (ingredient_id, category_id) values (new_ingredient_id, category_id);
    return new_ingredient_id;
end
$$ language plpgsql;

create or replace function create_recipe(name text, directions text, source_url text, image_filename text) returns integer as $$
declare
    new_recipe_id integer;
begin
    insert into recipes (recipe_name, directions, source_url, image_filename) values (name, directions, source_url, image_filename) returning recipe_id into new_recipe_id;
    return new_recipe_id;
end
$$ language plpgsql;

create or replace function seed_db() returns void as $body$
declare
    breakfast_id integer;
    meal_id integer;
    snack_id integer;
    dinner_id integer;

    produce_id integer;
    meat_id integer;
    dairy_id integer;
    frozen_id integer;
    bakery_id integer;
    canned_id integer;
    deli_id integer;
    pasta_id integer;
    spices_id integer;

    kilogram_id integer;
    liter_id integer;
    gram_id integer;
    pound_id integer;
    ounce_id integer;
    cup_id integer;
    teaspoon_id integer;
    tablespoon_id integer;
    piece_id integer;

    egg_id integer;
    milk_id integer;
    ricotta_cheese_id integer;
    garlic_id integer;
    beef_id integer;
    sausage_id integer;
    parsley_id integer;
    tomato_paste_id integer;
    lasagna_noodles_id integer;
    onion_id integer;
    parmesan_cheese_id integer;
    italian_seasoning_id integer;
    pasta_sauce_id integer;
    salt_id integer;
    sugar_id integer;
    flour_id integer;
    baking_powder_id integer;

    good_old_fashioned_pancakes_id integer;
    easy_homemade_lasagna_id integer;

begin
    breakfast_id := create_recipe_category('breakfast');
    meal_id := create_recipe_category('meal');
    snack_id := create_recipe_category('snack');
    dinner_id := create_recipe_category('dinner');

    produce_id := create_ingredient_category('produce');
    meat_id := create_ingredient_category('meat');
    dairy_id := create_ingredient_category('dairy');
    frozen_id := create_ingredient_category('frozen');
    bakery_id := create_ingredient_category('bakery');
    canned_id := create_ingredient_category('canned');
    deli_id := create_ingredient_category('deli');
    pasta_id := create_ingredient_category('pasta, rice, cereal');
    spices_id := create_ingredient_category('spices, seasoning');

    kilogram_id := create_unit('kilogram', 'mass', 0.001);
    liter_id := create_unit('liter', 'volume', 0.001);
    gram_id := create_unit('gram', 'mass', 1);
    pound_id := create_unit('pound', 'mass', 0.002205);
    ounce_id := create_unit('ounce', 'mass', 0.035274);
    cup_id := create_unit('cup', 'volume', 0.004226);
    teaspoon_id := create_unit('teaspoon', 'volume', 0.20288);
    tablespoon_id := create_unit('tablespoon', 'volume', 0.06762);
    piece_id := create_unit('piece', 'count', 1);

    egg_id := create_ingredient('egg', produce_id, piece_id);
    milk_id := create_ingredient('milk', dairy_id, cup_id);
    ricotta_cheese_id := create_ingredient('ricotta cheese', dairy_id, cup_id);
    garlic_id := create_ingredient('garlic', produce_id, piece_id);
    beef_id := create_ingredient('beef', meat_id, pound_id);
    sausage_id := create_ingredient('sausage', meat_id, pound_id);
    parsley_id := create_ingredient('parsley', produce_id, cup_id);
    tomato_paste_id := create_ingredient('tomato paste', canned_id, cup_id);
    lasagna_noodles_id := create_ingredient('lasagna noodles', pasta_id, piece_id);
    onion_id := create_ingredient('onion', produce_id, piece_id);
    parmesan_cheese_id := create_ingredient('parmesan cheese', dairy_id, cup_id);
    italian_seasoning_id := create_ingredient('italian seasoning', spices_id, tablespoon_id);
    pasta_sauce_id := create_ingredient('pasta sauce', canned_id, cup_id);
    salt_id := create_ingredient('salt', spices_id, teaspoon_id);
    sugar_id := create_ingredient('sugar', spices_id, cup_id);
    flour_id := create_ingredient('flour', bakery_id, cup_id);
    baking_powder_id := create_ingredient('baking powder', bakery_id, teaspoon_id);

    good_old_fashioned_pancakes_id := create_recipe(
        $$Easy Homemade Lasagna,
        1. Boil pasta: Cook lasagna noodles in salted water according to the recipe.
        2. Prepare meat sauce: Cook sausage, beef, onion, and garlic. Drain, then add pasta sauce and simmer to thicken.
        3. Combine cheese mixture: Mix cheese ingredients together in a bowl.
        4. Layer & bake: Layer the meat sauce, cheese mixture, and lasagna noodles, then bake until the top is golden brown.$$,
        'https://www.spendwithpennies.com/easy-homemade-lasagna/',
        'lasagna.jpg'
    );
    easy_homemade_lasagna_id := create_recipe(
        $$Good Old Fashioned Pancakes,
        1. Sift flour, baking powder, sugar, and salt into a large bowl.
        2. Make a well in the center and add milk, melted butter, and egg; mix until smooth.
        3. Heat a lightly oiled griddle or pan over medium-high heat.
        4. Pour or scoop about 1/4 cup of batter onto the griddle for each pancake.
        5. Cook until bubbles form and edges are dry (about 2–3 minutes), then flip.
        6. Cook the other side until browned.
        7. Repeat with remaining batter.$$,
        'https://www.allrecipes.com/recipe/21014/good-old-fashioned-pancakes/',
        'pancakes.jpg'
    );

    insert into recipes_to_recipe_categories (recipe_id, recipe_category_id) values
        (good_old_fashioned_pancakes_id, breakfast_id),
        (easy_homemade_lasagna_id, meal_id);

    insert into recipes_to_ingredients (recipe_id, ingredient_id, recipe_ingredient_quantity) values
        (good_old_fashioned_pancakes_id, egg_id, 1),
        (good_old_fashioned_pancakes_id, milk_id, 1),
        (good_old_fashioned_pancakes_id, ricotta_cheese_id, 1),
        (good_old_fashioned_pancakes_id, garlic_id, 1),
        (good_old_fashioned_pancakes_id, beef_id, 1),
        (good_old_fashioned_pancakes_id, lasagna_noodles_id, 10),
        (good_old_fashioned_pancakes_id, tomato_paste_id, 1),
        (good_old_fashioned_pancakes_id, italian_seasoning_id, 1),
        (easy_homemade_lasagna_id, flour_id, 1),
        (easy_homemade_lasagna_id, baking_powder_id, 1),
        (easy_homemade_lasagna_id, sugar_id, 1),
        (easy_homemade_lasagna_id, salt_id, 1),
        (easy_homemade_lasagna_id, milk_id, 1);

end
$body$ language plpgsql;
