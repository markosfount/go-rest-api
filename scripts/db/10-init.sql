CREATE USER "user" WITH PASSWORD 'password';
-- CREATE DATABASE "movies" OWNER "user";
CREATE TABLE movies (
                        id SERIAL,
                        movieID varchar(50) NOT NULL UNIQUE,
                        movieName varchar(50) NOT NULL,
                        PRIMARY KEY (id)
);
