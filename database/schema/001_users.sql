-- +goose Up
create table users(id integer primary key autoincrement, username text, password text)
-- +goose Down
drop table users;
