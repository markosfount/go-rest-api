CREATE USER "user" WITH PASSWORD 'password';
ALTER DATABASE movies OWNER TO "user";
CREATE TABLE movies (
                        id SERIAL,
                        movieID varchar(50) NOT NULL UNIQUE,
                        movieName varchar(50) NOT NULL,
                        PRIMARY KEY (id)
);
GRANT ALL ON movies TO "user";
GRANT ALL ON SEQUENCE movies_id_seq TO "user";
