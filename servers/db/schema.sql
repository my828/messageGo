create table if not exists users (
    id int not null auto_increment primary key,
    email varchar(255) not null,
    pass_hash varchar(64) not null,
    user_name varchar(255) not null,
    first_name varchar(64) not null,
    last_name varchar(128) not null,
    photo_url varchar(128) not null,
    unique key(email, user_name)
)
