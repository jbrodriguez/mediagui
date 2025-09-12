-- Drop location triggers and table
DROP TRIGGER IF EXISTS location_ai;
DROP TRIGGER IF EXISTS location_au;
DROP TRIGGER IF EXISTS location_bd;
DROP TRIGGER IF EXISTS location_bu;
DROP TABLE IF EXISTS movielocation;

-- Drop actor triggers and table
DROP TRIGGER IF EXISTS actor_ai;
DROP TRIGGER IF EXISTS actor_au;
DROP TRIGGER IF EXISTS actor_bd;
DROP TRIGGER IF EXISTS actor_bu;
DROP TABLE IF EXISTS movieactor;

-- Drop director triggers and table
DROP TRIGGER IF EXISTS director_ai;
DROP TRIGGER IF EXISTS director_au;
DROP TRIGGER IF EXISTS director_bd;
DROP TRIGGER IF EXISTS director_bu;
DROP TABLE IF EXISTS moviedirector;

-- Drop country triggers and table
DROP TRIGGER IF EXISTS country_ai;
DROP TRIGGER IF EXISTS country_au;
DROP TRIGGER IF EXISTS country_bd;
DROP TRIGGER IF EXISTS country_bu;
DROP TABLE IF EXISTS moviecountry;

-- Drop genre triggers and table
DROP TRIGGER IF EXISTS genre_ai;
DROP TRIGGER IF EXISTS genre_au;
DROP TRIGGER IF EXISTS genre_bd;
DROP TRIGGER IF EXISTS genre_bu;
DROP TABLE IF EXISTS moviegenre;

-- Drop movietitle triggers and table
DROP TRIGGER IF EXISTS movie_ai;
DROP TRIGGER IF EXISTS movie_au;
DROP TRIGGER IF EXISTS movie_bd;
DROP TRIGGER IF EXISTS movie_bu;
DROP TABLE IF EXISTS movietitle;

-- Drop basic indexes
DROP INDEX IF EXISTS movie_filetype_idx;
DROP INDEX IF EXISTS movie_location_idx;
DROP INDEX IF EXISTS movie_title_idx;

-- Drop the main table
DROP TABLE IF EXISTS movie;
