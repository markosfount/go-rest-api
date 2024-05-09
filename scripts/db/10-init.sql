CREATE USER "user" WITH PASSWORD 'password';
ALTER DATABASE movies OWNER TO "user";
CREATE TABLE movies (
                        id SERIAL,
                        movieID varchar(50) NOT NULL UNIQUE,
                        movieName varchar(50) NOT NULL,
                        overview text,
                        runtime smallint,
                        PRIMARY KEY (id)
);
GRANT ALL ON movies TO "user";
GRANT ALL ON SEQUENCE movies_id_seq TO "user";
-- CREATE DATABASE users;
-- ALTER DATABASE users OWNER TO "user";
CREATE TABLE users (
                        id SERIAL,
                        username varchar(50) NOT NULL UNIQUE,
                        password varchar(50) NOT NULL,
                        PRIMARY KEY (id)
);
GRANT ALL ON users TO "user";
GRANT ALL ON SEQUENCE users_id_seq TO "user";
