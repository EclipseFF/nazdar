create table category
(
    id   serial primary key,
    name varchar unique
);

create table item_category
(
    item_id int references items(id),
    category_id int references category(id),
    primary key(item_id, category_id)
);
