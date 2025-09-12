-- Create the main movie table
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

-- Create basic indexes
CREATE INDEX movie_title_idx ON movie (title);
CREATE INDEX movie_location_idx ON movie (location);
CREATE INDEX movie_filetype_idx ON movie (filetype);

-- Create FTS5 virtual table for movie titles
CREATE VIRTUAL TABLE movietitle USING fts5(content="movie", title, original_title, file_title);

-- Create triggers for movietitle FTS table
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

-- Create FTS5 virtual table for movie genres
CREATE VIRTUAL TABLE moviegenre USING fts5(content="movie", genres);

-- Create triggers for moviegenre FTS table
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

-- Create FTS5 virtual table for movie countries
CREATE VIRTUAL TABLE moviecountry USING fts5(content="movie", countries);

-- Create triggers for moviecountry FTS table
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

-- Create FTS5 virtual table for movie directors
CREATE VIRTUAL TABLE moviedirector USING fts5(content="movie", director);

-- Create triggers for moviedirector FTS table
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

-- Create FTS5 virtual table for movie actors
CREATE VIRTUAL TABLE movieactor USING fts5(content="movie", actors);

-- Create triggers for movieactor FTS table
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

-- Create FTS5 virtual table for movie locations
CREATE VIRTUAL TABLE movielocation USING fts5(content="movie", location);

-- Create triggers for movielocation FTS table
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
