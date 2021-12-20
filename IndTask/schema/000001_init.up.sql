CREATE TABLE books
(
    id serial not null unique primary key,
    book_name varchar(255) not null unique,
    cost decimal not null,
    cover varchar(255) not null,
    published int not null,
    pages integer not null,
    amount integer not null
    authors_id []integer not null,
    genres_id []integer not null

);

CREATE TABLE authors
(
    id serial not null unique primary key,
    author_name varchar(255) not null unique,
    author_foto varchar(255)
);

CREATE TABLE users
(
    id serial not null unique primary key,
    surname  varchar(255) not null unique,
    user_name varchar(255) not null,
    patronymic varchar(255),
    pasp_number varchar(50) unique,
    email varchar(255) not null unique,
    adress varchar(255),
    birth_date timestamp with time zone not null
);

CREATE TABLE book_author
(
    id serial not null unique primary key,
    book_id int references books(id) on delete set null,
    author_id int references authors(id) on delete set null
);


CREATE TABLE genre
(
    id serial not null unique primary key,
    genre_name varchar(255) not null unique
);

CREATE TABLE book_genre
(
    id serial not null unique primary key,
    book_id int references books(id) on delete set null,
    genre_id int references genre(id) on delete set null
);

CREATE TABLE list_books
(
    id serial not null unique primary key,
    book_id int references books(id) not null,
    issued bool default false not null,
    rent_number int not null,
    rent_cost int not null,
    reg_date timestamp with time zone not null,
    condition int not null
);

CREATE TABLE issue_act
(
    id serial not null unique primary key,
    user_id int references users(id) not null,
    book_id int references books(id) not null,
    rental_time interval not null,
    return_date timestamp with time zone not null,
    pre_cost decimal not null,
    status bool not null
);

CREATE TABLE return_act
(
    id serial not null unique primary key,
    user_id int references users(id) not null,
    book_id int references books(id) not null,
    cost decimal not null,
    return_date timestamp with time zone not null,
    foto varchar[],
    fine decimal,
    condition_decrese int,
    rating int
);


