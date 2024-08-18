create table admins(
    id serial primary key,
    name varchar,
    password varchar
);

create table admin_sessions(
    id serial primary key,
    admin_id int references admins(id),
    token varchar
              );