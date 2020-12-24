DROP TABLE IF EXISTS movie;
DROP TABLE IF EXISTS movietitle;
DROP TABLE IF EXISTS moviegenre;
DROP TABLE IF EXISTS moviecountry;
DROP TABLE IF EXISTS moviedirector;
DROP TABLE IF EXISTS movieactor;

DROP INDEX IF EXISTS movie_filetype_idx;
DROP INDEX IF EXISTS movie_location_idx;
DROP INDEX IF EXISTS movie_title_idx;

DROP TRIGGER IF EXISTS movie_ai;
DROP TRIGGER IF EXISTS movie_au;
DROP TRIGGER IF EXISTS movie_bd;
DROP TRIGGER IF EXISTS movie_bu;

DROP TRIGGER IF EXISTS genre_ai;
DROP TRIGGER IF EXISTS genre_au;
DROP TRIGGER IF EXISTS genre_bd;
DROP TRIGGER IF EXISTS genre_bu;

DROP TRIGGER IF EXISTS country_ai;
DROP TRIGGER IF EXISTS country_au;
DROP TRIGGER IF EXISTS country_bd;
DROP TRIGGER IF EXISTS country_bu;

DROP TRIGGER IF EXISTS director_ai;
DROP TRIGGER IF EXISTS director_au;
DROP TRIGGER IF EXISTS director_bd;
DROP TRIGGER IF EXISTS director_bu;

DROP TRIGGER IF EXISTS actor_ai;
DROP TRIGGER IF EXISTS actor_au;
DROP TRIGGER IF EXISTS actor_bd;
DROP TRIGGER IF EXISTS actor_bu;

CREATE TABLE movie
(
  title text,
  original_title text,
  file_title text,
  year integer,
  runtime integer,
  tmdb_id integer,
  imdb_id text,
  overview text,
  tagline text,
  resolution text,
  filetype text,
  location text,
  cover text,
  backdrop text,
  genres text,
  vote_average integer,
  vote_count integer,
  countries text,
  added text,
  modified text,
  last_watched text,
  all_watched text,
  count_watched integer,
  score integer,
  director text,
  writer text,
  actors text,
  awards text,
  imdb_rating integer,
  imdb_votes integer,
  show_if_duplicate integer,
  stub integer
);
CREATE INDEX movie_title_idx ON movie (title);
CREATE INDEX movie_location_idx ON movie (location);
CREATE INDEX movie_filetype_idx ON movie (filetype);

/* titles */
CREATE VIRTUAL TABLE movietitle USING fts4(content="movie", title, original_title, file_title);
CREATE TRIGGER movie_bu BEFORE UPDATE ON movie BEGIN
	DELETE FROM movietitle WHERE docid=old.rowid;
END;

CREATE TRIGGER movie_bd BEFORE DELETE ON movie BEGIN
	DELETE FROM movietitle WHERE docid=old.rowid;
END;

CREATE TRIGGER movie_au AFTER UPDATE ON movie BEGIN
	INSERT INTO movietitle(docid, title, original_title, file_title) VALUES (new.rowid, new.title, new.original_title, new.file_title);
END;

CREATE TRIGGER movie_ai AFTER INSERT ON movie BEGIN
	INSERT INTO movietitle(docid, title, original_title, file_title) VALUES (new.rowid, new.title, new.original_title, new.file_title);
END;

/* genres */
CREATE VIRTUAL TABLE moviegenre USING fts4(content="movie", genres);
CREATE TRIGGER genre_bu BEFORE UPDATE ON movie BEGIN
  DELETE FROM moviegenre WHERE docid=old.rowid;
END;

CREATE TRIGGER genre_bd BEFORE DELETE ON movie BEGIN
  DELETE FROM moviegenre WHERE docid=old.rowid;
END;

CREATE TRIGGER genre_au AFTER UPDATE ON movie BEGIN
  INSERT INTO moviegenre(docid, genres) VALUES (new.rowid, new.genres);
END;

CREATE TRIGGER genre_ai AFTER INSERT ON movie BEGIN
  INSERT INTO moviegenre(docid, genres) VALUES (new.rowid, new.genres);
END;

/* country */
CREATE VIRTUAL TABLE moviecountry USING fts4(content="movie", countries);
CREATE TRIGGER country_bu BEFORE UPDATE ON movie BEGIN
  DELETE FROM moviecountry WHERE docid=old.rowid;
END;

CREATE TRIGGER country_bd BEFORE DELETE ON movie BEGIN
  DELETE FROM moviecountry WHERE docid=old.rowid;
END;

CREATE TRIGGER country_au AFTER UPDATE ON movie BEGIN
  INSERT INTO moviecountry(docid, countries) VALUES (new.rowid, new.countries);
END;

CREATE TRIGGER country_ai AFTER INSERT ON movie BEGIN
  INSERT INTO moviecountry(docid, countries) VALUES (new.rowid, new.countries);
END;

/* director */
CREATE VIRTUAL TABLE moviedirector USING fts4(content="movie", director);
CREATE TRIGGER director_bu BEFORE UPDATE ON movie BEGIN
  DELETE FROM moviedirector WHERE docid=old.rowid;
END;

CREATE TRIGGER director_bd BEFORE DELETE ON movie BEGIN
  DELETE FROM moviedirector WHERE docid=old.rowid;
END;

CREATE TRIGGER director_au AFTER UPDATE ON movie BEGIN
  INSERT INTO moviedirector(docid, director) VALUES (new.rowid, new.director);
END;

CREATE TRIGGER director_ai AFTER INSERT ON movie BEGIN
  INSERT INTO moviedirector(docid, director) VALUES (new.rowid, new.director);
END;

/* actor */
CREATE VIRTUAL TABLE movieactor USING fts4(content="movie", actors);
CREATE TRIGGER actor_bu BEFORE UPDATE ON movie BEGIN
  DELETE FROM movieactor WHERE docid=old.rowid;
END;

CREATE TRIGGER actor_bd BEFORE DELETE ON movie BEGIN
  DELETE FROM movieactor WHERE docid=old.rowid;
END;

CREATE TRIGGER actor_au AFTER UPDATE ON movie BEGIN
  INSERT INTO movieactor(docid, actors) VALUES (new.rowid, new.actors);
END;

CREATE TRIGGER actor_ai AFTER INSERT ON movie BEGIN
  INSERT INTO movieactor(docid, actors) VALUES (new.rowid, new.actors);
END;

/* location */
CREATE VIRTUAL TABLE movielocation USING fts4(content="movie", location);
CREATE TRIGGER location_bu BEFORE UPDATE ON movie BEGIN
  DELETE FROM movielocation WHERE docid=old.rowid;
END;

CREATE TRIGGER location_bd BEFORE DELETE ON movie BEGIN
  DELETE FROM movielocation WHERE docid=old.rowid;
END;

CREATE TRIGGER location_au AFTER UPDATE ON movie BEGIN
  INSERT INTO movielocation(docid, location) VALUES (new.rowid, new.location);
END;

CREATE TRIGGER location_ai AFTER INSERT ON movie BEGIN
  INSERT INTO movielocation(docid, location) VALUES (new.rowid, new.location);
END;
