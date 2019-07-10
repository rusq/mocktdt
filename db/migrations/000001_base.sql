-- +goose Up
create table if not exists personas (
    id integer primary key
    ,name varchar(50)
    ,dob timestamp
);

-- +goose Down
drop table personas;
