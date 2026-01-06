create type role as enum (
    'standard',
    'admin'
);

create table if not exists users (
    user_id serial primary key,
    username varchar(255) unique not null,
    password_hash varchar(255) not null,
    role role not null default 'standard',

    created_at timestamp not null default now(),
    updated_at timestamp not null default now()
);

create table if not exists user_sessions (
    session_id varchar(255) primary key,
    user_id integer not null references users(user_id) on delete cascade,
    last_access timestamp not null default now()
);

create table if not exists recipe_categories (
    recipe_category_id serial primary key,
    recipe_category_name varchar(255) unique not null
);

create table if not exists recipes (
    recipe_id serial primary key,
    recipe_name varchar(255) unique not null,
    directions text not null,
    source_url varchar(255),
    image_filename varchar(255) not null,

    created_at timestamp not null default now(),
    updated_at timestamp not null default now()
);

create table if not exists recipes_to_recipe_categories (
    recipe_to_recipe_category_id serial primary key,
    recipe_id integer not null references recipes(recipe_id) on delete cascade,
    recipe_category_id integer not null references recipe_categories(recipe_category_id) on delete cascade
);

create type unit_category as enum (
    'mass',
    'volume',
    'count'
);

create table if not exists units (
    unit_id serial primary key,
    unit_name varchar(255) unique not null,
    conversion_factor real not null,
    unit_category unit_category
);

create table if not exists ingredient_categories (
    ingredient_category_id serial primary key,
    ingredient_category_name varchar(255) unique not null
);

create table if not exists ingredients (
    ingredient_id serial primary key,
    ingredient_name varchar(255) unique not null,
    unit_id integer references units(unit_id) on delete set null
);

create table if not exists ingredients_to_categories (
    ingredient_to_category_id serial primary key,
    ingredient_id integer not null references ingredients(ingredient_id) on delete cascade,
    category_id integer not null references ingredient_categories(ingredient_category_id) on delete cascade
);

create table if not exists recipes_to_ingredients (
    recipe_id integer not null references recipes(recipe_id) on delete cascade,
    ingredient_id integer not null references ingredients(ingredient_id) on delete restrict,
    recipe_ingredient_quantity real not null,
    primary key (recipe_id, ingredient_id)
);

create table if not exists lists (
    list_id serial primary key,
    list_name varchar(255),
    user_id integer not null references users(user_id) on delete cascade,

    created_at timestamp not null default now(),
    updated_at timestamp not null default now()
);

create table if not exists lists_to_recipes (
    list_id integer not null references lists(list_id) on delete cascade,
    recipe_id integer not null references recipes(recipe_id) on delete cascade,
    list_recipe_quantity real not null,
    primary key (list_id, recipe_id)
);

create table if not exists custom_list_items (
    custom_list_item_id serial primary key,
    custom_list_item_name varchar(255) not null,
    custom_list_item_quantity real,
    custom_unit varchar(255),
    unit_id integer references units(unit_id) on delete set null,
    check ((custom_unit is null) or (unit_id is null))
);

create table if not exists list_items (
    list_item_id serial primary key,
    list_id integer not null references lists(list_id) on delete cascade,
    list_item_quantity real not null,
    list_item_associated_recipe_id integer references recipes(recipe_id) on delete set null,
    ingredient_id integer references ingredients(ingredient_id) on delete cascade,
    custom_list_item_id integer references custom_list_items(custom_list_item_id) on delete cascade,
    check
        ((ingredient_id is not null and custom_list_item_id is null)
        or (ingredient_id is null and custom_list_item_id is not null))
);

