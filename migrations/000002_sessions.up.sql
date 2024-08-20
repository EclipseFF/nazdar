create table users(
                      id serial primary key,
                      phone_number varchar unique,
                      name varchar,
                      createdAt timestamp
);

create table sessions(
                         user_id int references users(id),
                         token varchar not null,
                         primary key(user_id, token)
);