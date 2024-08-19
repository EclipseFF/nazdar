create table if not exists items (
                                     id serial primary key,
                                     name varchar not null,
                                    price integer not null,
                                     description varchar,
                                     images varchar[]
);