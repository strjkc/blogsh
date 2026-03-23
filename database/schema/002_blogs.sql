-- +goose Up
create table blogs(id integer primary key autoincrement, content text, user_id integer, foreign key(user_id) references users(id) on delete cascade)
-- +goose Down
drop table blogs
