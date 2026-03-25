-- +goose Up
create table blogs(id integer primary key autoincrement,
    title text not null unique, 
    content text not null,
    category text not null,
    tags text not null,
    createdAt text not null,
    updatedAt text not null
);
-- +goose Down
drop table blogs;
