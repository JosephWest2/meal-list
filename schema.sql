create type role as enum (
    'standard',
    'admin'
);

create table if not exists users (
    id serial primary key,
    username varchar(255) unique not null,
    password_hash varchar(255) not null,
    role role,

    created_at timestamp not null default now(),
    updated_at timestamp not null default now()
);

create table if not exists user_sessions (
    id varchar(255) primary key,
    user_id integer not null references users(id) on delete cascade,
    last_access timestamp not null default now()
);

create table if not exists recipe_categories (
    id serial primary key,
    name varchar(255) unique not null
);

create table if not exists recipes (
    id serial primary key,
    name varchar(255) unique not null,
    directions text not null,
    source_url varchar(255),
    image_filename varchar(255) not null,

    created_at timestamp not null default now(),
    updated_at timestamp not null default now()
);

create table if not exists recipes_to_recipe_categories (
    id serial primary key,
    recipe_id integer not null references recipes(id) on delete cascade,
    recipe_category_id integer not null references recipe_categories(id) on delete cascade
);

create type unit_category as enum (
    'mass',
    'volume',
    'count'
);

create table if not exists units (
    id serial primary key,
    name varchar(255) unique not null,
    category unit_category
);

create table if not exists ingredient_categories (
    id serial primary key,
    name varchar(255) unique not null
);

create table if not exists ingredients (
    id serial primary key,
    name varchar(255) unique not null,
    ingredient_category_id integer references ingredient_categories(id) on delete set null,
    unit_id integer references units(id) on delete set null
);

create table if not exists recipes_to_ingredients (
    recipe_id integer not null references recipes(id) on delete cascade,
    ingredient_id integer not null references ingredients(id) on delete restrict,
    quantity real not null,
    primary key (recipe_id, ingredient_id)
);

create table if not exists lists (
    id serial primary key,
    name varchar(255),
    user_id integer references users(id) on delete cascade,

    created_at timestamp not null default now(),
    updated_at timestamp not null default now()
);

create table if not exists lists_to_recipes (
    list_id integer not null references lists(id) on delete cascade,
    recipe_id integer not null references recipes(id) on delete cascade,
    quantity real not null,
    primary key (list_id, recipe_id)
);

create table if not exists custom_list_items (
    id serial primary key,
    list_id integer not null references lists(id) on delete cascade,
    name varchar(255),
    quantity real,
    custom_unit varchar(255),
    unit_id integer references units(id) on delete set null,
    ingredient_id integer references ingredients(id) on delete cascade,
    check 
        ((name is not null and ingredient_id is null)
        or (name is null and ingredient_id is not null))
);

create table if not exists checked_list_items (
    id serial primary key,
    list_id integer not null references lists(id) on delete cascade,
    ingredient_id integer references ingredients(id) on delete cascade,
    custom_list_item_id integer references custom_list_items(id) on delete cascade,
    check 
        ((ingredient_id is not null and custom_list_item_id is null)
        or (ingredient_id is null and custom_list_item_id is not null))
);

CREATE OR REPLACE FUNCTION before_update_updated_at() RETURNS trigger AS
$BODY$
BEGIN
    IF row(NEW.*::text) IS DISTINCT FROM row(OLD.*::text) THEN
        NEW.updated_at = now();
    END IF;

    RETURN NEW;
END;
$BODY$
LANGUAGE plpgsql;

DO $BODY$
DECLARE t text;
BEGIN
    FOR t IN
        SELECT table_name
        FROM information_schema.columns
        WHERE (
            column_name = 'updated_at'
            AND (
                SELECT 1
                FROM information_schema.triggers
                WHERE trigger_name = 'before_update_updated_at_' || table_name
            ) IS NULL
        )
    LOOP
        EXECUTE format('
            CREATE TRIGGER before_update_updated_at_%s
            BEFORE UPDATE ON %I
            FOR EACH ROW EXECUTE PROCEDURE before_update_updated_at();
        ', t, t);
    END loop;
END;
$BODY$
LANGUAGE plpgsql;