
CREATE SCHEMA IF NOT EXISTS public;

ALTER SCHEMA public OWNER TO postgres;

SET default_tablespace = '';
SET default_table_access_method = heap;

SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SELECT pg_catalog.set_config('search_path', '', false);
SET check_function_bodies = false;
SET xmloption = content;
SET client_min_messages = warning;
SET row_security = off;

CREATE TABLE public.recipes (
    id serial4 NOT NULL UNIQUE,
    title VARCHAR(255) NOT NULL,
    description TEXT NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL
);

ALTER TABLE public.recipes OWNER TO postgres;

CREATE TABLE public.recipe_steps (
    id serial4 NOT NULL UNIQUE,
    recipe_id int4 NOT NULL,
    number int4 NOT NULL,
    step_text TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL
);

ALTER TABLE public.recipe_steps OWNER TO postgres;

CREATE TABLE public.recipes_ingredients (
    recipe_id int4 NOT NULL,
    name VARCHAR(64) NOT NULL,
    CONSTRAINT ingredient_unique UNIQUE(recipe_id, name)
);

ALTER TABLE public.recipes_ingredients OWNER TO postgres;