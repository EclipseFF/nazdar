create table user_items(
    user_id int references users(id),
    item_id int references items(id),
    count int,
    primary key(user_id, item_id)
);