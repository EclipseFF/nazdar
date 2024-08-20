create table user_items(
    user_id int references users(id),
    item_id int references items(id),
    count int,
    primary key(user_id, item_id)
);

create table user_orders(
    id serial primary key,
    user_id int references users(id),
    item_ids int[],
    count int[],
    total_price int
);