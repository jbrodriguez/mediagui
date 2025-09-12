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
CREATE VIRTUAL TABLE movietitle USING fts5(content="movie", title, original_title, file_title)
/* movietitle(title,original_title,file_title) */;
CREATE TABLE IF NOT EXISTS 'movietitle_data'(id INTEGER PRIMARY KEY, block BLOB);
CREATE TABLE IF NOT EXISTS 'movietitle_idx'(segid, term, pgno, PRIMARY KEY(segid, term)) WITHOUT ROWID;
CREATE TABLE IF NOT EXISTS 'movietitle_docsize'(id INTEGER PRIMARY KEY, sz BLOB);
CREATE TABLE IF NOT EXISTS 'movietitle_config'(k PRIMARY KEY, v) WITHOUT ROWID;
CREATE TRIGGER movie_bu BEFORE UPDATE ON movie BEGIN
        DELETE FROM movietitle WHERE rowid=old.rowid;
END;
CREATE TRIGGER movie_bd BEFORE DELETE ON movie BEGIN
        DELETE FROM movietitle WHERE rowid=old.rowid;
END;
CREATE TRIGGER movie_au AFTER UPDATE ON movie BEGIN
        INSERT INTO movietitle(rowid, title, original_title, file_title) VALUES (new.rowid, new.title, new.original_title, new.file_title);
END;
CREATE TRIGGER movie_ai AFTER INSERT ON movie BEGIN
        INSERT INTO movietitle(rowid, title, original_title, file_title) VALUES (new.rowid, new.title, new.original_title, new.file_title);
END;
CREATE VIRTUAL TABLE moviegenre USING fts5(content="movie", genres)
/* moviegenre(genres) */;
CREATE TABLE IF NOT EXISTS 'moviegenre_data'(id INTEGER PRIMARY KEY, block BLOB);
CREATE TABLE IF NOT EXISTS 'moviegenre_idx'(segid, term, pgno, PRIMARY KEY(segid, term)) WITHOUT ROWID;
CREATE TABLE IF NOT EXISTS 'moviegenre_docsize'(id INTEGER PRIMARY KEY, sz BLOB);
CREATE TABLE IF NOT EXISTS 'moviegenre_config'(k PRIMARY KEY, v) WITHOUT ROWID;
CREATE TRIGGER genre_bu BEFORE UPDATE ON movie BEGIN
  DELETE FROM moviegenre WHERE rowid=old.rowid;
END;
CREATE TRIGGER genre_bd BEFORE DELETE ON movie BEGIN
  DELETE FROM moviegenre WHERE rowid=old.rowid;
END;
CREATE TRIGGER genre_au AFTER UPDATE ON movie BEGIN
  INSERT INTO moviegenre(rowid, genres) VALUES (new.rowid, new.genres);
END;
CREATE TRIGGER genre_ai AFTER INSERT ON movie BEGIN
  INSERT INTO moviegenre(rowid, genres) VALUES (new.rowid, new.genres);
END;
CREATE VIRTUAL TABLE moviecountry USING fts5(content="movie", countries)
/* moviecountry(countries) */;
CREATE TABLE IF NOT EXISTS 'moviecountry_data'(id INTEGER PRIMARY KEY, block BLOB);
CREATE TABLE IF NOT EXISTS 'moviecountry_idx'(segid, term, pgno, PRIMARY KEY(segid, term)) WITHOUT ROWID;
CREATE TABLE IF NOT EXISTS 'moviecountry_docsize'(id INTEGER PRIMARY KEY, sz BLOB);
CREATE TABLE IF NOT EXISTS 'moviecountry_config'(k PRIMARY KEY, v) WITHOUT ROWID;
CREATE TRIGGER country_bu BEFORE UPDATE ON movie BEGIN
  DELETE FROM moviecountry WHERE rowid=old.rowid;
END;
CREATE TRIGGER country_bd BEFORE DELETE ON movie BEGIN
  DELETE FROM moviecountry WHERE rowid=old.rowid;
END;
CREATE TRIGGER country_au AFTER UPDATE ON movie BEGIN
  INSERT INTO moviecountry(rowid, countries) VALUES (new.rowid, new.countries);
END;
CREATE TRIGGER country_ai AFTER INSERT ON movie BEGIN
  INSERT INTO moviecountry(rowid, countries) VALUES (new.rowid, new.countries);
END;
CREATE VIRTUAL TABLE moviedirector USING fts5(content="movie", director)
/* moviedirector(director) */;
CREATE TABLE IF NOT EXISTS 'moviedirector_data'(id INTEGER PRIMARY KEY, block BLOB);
CREATE TABLE IF NOT EXISTS 'moviedirector_idx'(segid, term, pgno, PRIMARY KEY(segid, term)) WITHOUT ROWID;
CREATE TABLE IF NOT EXISTS 'moviedirector_docsize'(id INTEGER PRIMARY KEY, sz BLOB);
CREATE TABLE IF NOT EXISTS 'moviedirector_config'(k PRIMARY KEY, v) WITHOUT ROWID;
CREATE TRIGGER director_bu BEFORE UPDATE ON movie BEGIN
  DELETE FROM moviedirector WHERE rowid=old.rowid;
END;
CREATE TRIGGER director_bd BEFORE DELETE ON movie BEGIN
  DELETE FROM moviedirector WHERE rowid=old.rowid;
END;
CREATE TRIGGER director_au AFTER UPDATE ON movie BEGIN
  INSERT INTO moviedirector(rowid, director) VALUES (new.rowid, new.director);
END;
CREATE TRIGGER director_ai AFTER INSERT ON movie BEGIN
  INSERT INTO moviedirector(rowid, director) VALUES (new.rowid, new.director);
END;
CREATE VIRTUAL TABLE movieactor USING fts5(content="movie", actors)
/* movieactor(actors) */;
CREATE TABLE IF NOT EXISTS 'movieactor_data'(id INTEGER PRIMARY KEY, block BLOB);
CREATE TABLE IF NOT EXISTS 'movieactor_idx'(segid, term, pgno, PRIMARY KEY(segid, term)) WITHOUT ROWID;
CREATE TABLE IF NOT EXISTS 'movieactor_docsize'(id INTEGER PRIMARY KEY, sz BLOB);
CREATE TABLE IF NOT EXISTS 'movieactor_config'(k PRIMARY KEY, v) WITHOUT ROWID;
CREATE TRIGGER actor_bu BEFORE UPDATE ON movie BEGIN
  DELETE FROM movieactor WHERE rowid=old.rowid;
END;
CREATE TRIGGER actor_bd BEFORE DELETE ON movie BEGIN
  DELETE FROM movieactor WHERE rowid=old.rowid;
END;
CREATE TRIGGER actor_au AFTER UPDATE ON movie BEGIN
  INSERT INTO movieactor(rowid, actors) VALUES (new.rowid, new.actors);
END;
CREATE TRIGGER actor_ai AFTER INSERT ON movie BEGIN
  INSERT INTO movieactor(rowid, actors) VALUES (new.rowid, new.actors);
END;
CREATE VIRTUAL TABLE movielocation USING fts5(content="movie", location)
/* movielocation(location) */;
CREATE TABLE IF NOT EXISTS 'movielocation_data'(id INTEGER PRIMARY KEY, block BLOB);
CREATE TABLE IF NOT EXISTS 'movielocation_idx'(segid, term, pgno, PRIMARY KEY(segid, term)) WITHOUT ROWID;
CREATE TABLE IF NOT EXISTS 'movielocation_docsize'(id INTEGER PRIMARY KEY, sz BLOB);
CREATE TABLE IF NOT EXISTS 'movielocation_config'(k PRIMARY KEY, v) WITHOUT ROWID;
CREATE TRIGGER location_bu BEFORE UPDATE ON movie BEGIN
  DELETE FROM movielocation WHERE rowid=old.rowid;
END;
CREATE TRIGGER location_bd BEFORE DELETE ON movie BEGIN
  DELETE FROM movielocation WHERE rowid=old.rowid;
END;
CREATE TRIGGER location_au AFTER UPDATE ON movie BEGIN
  INSERT INTO movielocation(rowid, location) VALUES (new.rowid, new.location);
END;
CREATE TRIGGER location_ai AFTER INSERT ON movie BEGIN
  INSERT INTO movielocation(rowid, location) VALUES (new.rowid, new.location);
END;
