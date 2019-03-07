create table if not exists users (
    id int not null auto_increment primary key,
    email varchar(255) not null,
    passHash binary(60) not null,
    userName varchar(255) not null,
    firstName varchar(64) not null,
    lastName varchar(128) not null,
    photoUrl varchar(128) not null,
    unique key(email),
    unique key(userName)
);

create table if not exists signinuser (
    id2 int not null auto_increment primary key,
    userID int not null,
    signinDatetime datetime not null,
    ipAddress varchar(50) not null
);