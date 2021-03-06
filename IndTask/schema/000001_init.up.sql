CREATE TABLE books
(
    id serial not null unique primary key,
    book_name varchar(255) not null,
    cost decimal not null,
    cover varchar(255) not null,
    published int not null,
    pages integer not null,
    amount integer not null
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
    surname  varchar(255) not null default null,
    user_name varchar(255) not null default null,
    patronymic varchar(255),
    pasp_number varchar(50) unique,
    email varchar(255) not null unique default null,
    address varchar(255),
    birth_date date not null
);

CREATE TABLE book_author
(
    book_id int references books(id) on delete cascade,
    author_id int references authors(id) on delete cascade,
    PRIMARY KEY(book_id, author_id)
);


CREATE TABLE genre
(
    id serial not null unique primary key,
    genre_name varchar(255) not null unique
);

CREATE TABLE book_genre
(
    book_id int references books(id) on delete cascade,
    genre_id int references genre(id) on delete cascade,
    PRIMARY KEY(book_id, genre_id)

);

CREATE TABLE list_books
(
    id serial not null unique primary key,
    book_id int references books(id) not null,
    issued bool not null,
    rent_number int not null,
    rent_cost decimal not null,
    reg_date timestamp with time zone not null,
    condition int not null,
    scrapped bool
);

CREATE TYPE stat AS ENUM ('open', 'closed');

CREATE TABLE act
(
    id serial not null unique primary key,
    user_id int references users(id) not null,
    listbook_id int references list_books(id) not null,
    rental_time int not null,
    return_date timestamp with time zone not null,
    pre_cost decimal not null,
    cost decimal,
    status stat not null,
    actual_return_date timestamp with time zone,
    foto varchar(255) array[5],
    fine decimal,
    condition_decrese int,
    rating int
);

insert into genre (genre_name) values ('Novel'), ('Fantasy'), ('Detective'), ('Adventure'), ('Erotic'), ('Triller'), ('Philosophical'), ('Satire'), ('Comedy'), ('Crime'), ('Horror'), ('Business');

insert into authors (author_name, author_foto) values ('?????????????? ??????????????', 'images/authors/Rudyard_kipling.jpg'),
                                                            ('???????? ????????', 'images/authors/Mark_Tven.jpeg'),
                                                            ('???????????? ????????????', 'images/authors/George_Orwel.jpg'),
                                                            ('???????????? ??????????????', 'images/authors/Maksim_Gorky.jpeg'),
                                                            ('?????????????????? ????????????', 'images/authors/Aleksandr_Kuprin.jpg'),
                                                            ('???????? ??????????', 'images/authors/Ivan_Bunin.jpg'),
                                                            ('?????????? ????????', 'images/authors/Tomas_Mann.jpg'),
                                                            ('???????? ????????????', 'images/authors/Jack_London.jpg'),
                                                            ('?????????? ??????????', 'images/authors/Franc_Kafka.jpg'),
                                                            ('?????????? ??????????????????', 'images/authors/Boris_Pasternak.jpg'),
                                                            ('?????????? ????????????', 'images/authors/Agata_Kristi.jpg'),
                                                            ('???????????? ????????????????', 'images/authors/Mihail_Bulgakov.jpg'),
                                                            ('???????????? ??????????????????', 'images/authors/Ernest_Heminguey.jpg'),
                                                            ('???????????? ???? ????????-????????????????', 'images/authors/Antuan_de_Sent-Ekzupery.png'),
                                                            ('?????????????? ??????', 'images/authors/Bernard_Show.jpg'),
                                                            ('???????????? ????????????', 'images/authors/Robert_Govard.jpeg'),
                                                            ('?????????? ????????', 'images/authors/Emil_Zolia.jpg');

insert into users (surname, user_name, patronymic, pasp_number, email, address, birth_date) values ('??????????????', '????????????', '????????????????????', '123456', 'baraulia@yandex.ru', 'Minsk', '1965-07-20'),
                                                                                                  ('????????????', '????????????', '??????????????????????????', '123457', 'baraulia1@yandex.ru', 'Pemza', '1990-08-15'),
                                                                                                  ('????????????', '????????', '????????????????????', '123476', 'baraulia2@yandex.ru', 'Tagil', '1915-02-15'),
                                                                                                  ('??????????????', '??????????', '??????????', '127456', 'baraulia3@yandex.ru', 'Gomel', '2006-06-26'),
                                                                                                  ('??????????????', '????????????', '????????????????????', '173456', 'baraulia4@yandex.ru', 'Vitebsk', '1995-09-20'),
                                                                                                  ('????????????????', '????????????', '??????????????????', '723456', 'baraulia5@yandex.ru', 'Donetsk', '1975-10-11'),
                                                                                                  ('??????????', '??????????', '??????????????????', '823456', 'baraulia6@yandex.ru', 'Moskow', '1983-03-18'),
                                                                                                  ('????????????', '??????????????', '??????????????????', '183456', 'baraulia7@yandex.ru', 'Kolodoschi', '1999-01-09'),
                                                                                                  ('????????????????', '??????????????', '????????????????????', '128456', 'baraulia8@yandex.ru', 'Grodno', '2000-05-05'),
                                                                                                  ('????????????????', '????????????', '????????????????????', '123856', 'baraulia9@yandex.ru', 'Volgograd', '2010-12-07'),
                                                                                                  ('??????????????', '????????????', '??????????', '123486', 'baraulia10@yandex.ru', 'Vladivostok', '2002-11-02'),
                                                                                                  ('??????????????', '????????', '??????????????????', '123458', 'baraulia11@yandex.ru', 'Vileyka', '1998-10-22'),
                                                                                                  ('??????????????', '??????????', '??????????????????????????', '123459', 'baraulia12@yandex.ru', 'Tombov', '1983-09-12');

insert into books  (book_name, cost, cover, published, pages, amount) values ('?????? ??????????', 50, 'images/book_covers/Tom_Soyer_2021.jpg', 2021, 456, 2),
                                                                             ('?????????? ????????', 65, 'images/book_covers/Beliy_klik.2020.jpg', 2020, 560, 1),
                                                                             ('???????????? ??????????????!', 32, 'images/book_covers/Nelza_molchat_2018.jpg', 2018, 400, 2),
                                                                             ('??????????????', 50, 'images/book_covers/Process_2021.jpg', 2021, 288, 3),
                                                                             ('?????????????????? ????????', 25, 'images/book_covers/Volshebnaya_gora_2019.jpg', 2019, 928, 1),
                                                                             ('????????????', 50, 'images/book_covers/Dengi_2012.jpg', 2012, 512, 2),
                                                                             ('??????????-????????????', 38, 'images/book_covers/Conan_varvar_2021.jpg', 2021, 680, 3),
                                                                             ('???????????? ?? ??????????????????', 22, 'images/book_covers/Cezar_i_Kleopatra_2021.jpeg', 2021, 412, 1),
                                                                             ('?????????????????? ??????????', 15, 'images/book_covers/malenkiy_princ_2014.jpg', 2014, 220, 2),
                                                                             ('???????????? ?? ??????????????????', 62, 'images/book_covers/Master_i_Margarita_20121.jpg', 2021, 576, 2),
                                                                             ('???????????? ????????????', 13, 'images/book_covers/Doctor_Zivago_2021.jpg', 2021, 608, 3),
                                                                             ('???? ?????? ???????????? ??????????????', 50, 'images/book_covers/Po_kom_zvonit_kolokol_2019.jpg', 2019, 640, 2);

insert into list_books  (book_id, issued, rent_number, rent_cost, reg_date, condition, scrapped) values (1, false, 0, 0.65, '2022-01-01',100, false),
                                                                                                        (1, false, 0, 0.65, '2022-01-01',75, false),
                                                                                                        (2, false, 0, 0.845, '2022-01-01',100, false),
                                                                                                        (3, false, 0, 0.416, '2022-01-01',100, false),
                                                                                                        (3, false, 0, 0.416, '2022-01-01',100, false),
                                                                                                        (4, false, 0, 0.65, '2022-01-01',100, false),
                                                                                                        (4, false, 0, 0.65, '2022-01-01',100, false),
                                                                                                        (4, false, 0, 0.65, '2022-01-01',100, false),
                                                                                                        (5, false, 0, 0.325, '2022-01-01',100, false),
                                                                                                        (6, false, 0, 0.65, '2022-01-01',70, false),
                                                                                                        (6, false, 0, 0.65, '2022-01-01',100, false),
                                                                                                        (7, false, 0, 0.494, '2022-01-01',100, false),
                                                                                                        (7, false, 0, 0.494, '2022-01-01',100, false),
                                                                                                        (7, false, 0, 0.494, '2022-01-01',100, false),
                                                                                                        (8, false, 0, 0.286, '2022-01-01',100, false),
                                                                                                        (9, false, 0, 0.195, '2022-01-01',65, false),
                                                                                                        (9, false, 0, 0.195, '2022-01-01',100, false),
                                                                                                        (10, false, 0, 0.806, '2022-01-01',100, false),
                                                                                                        (10, false, 0, 0.806, '2022-01-01',100, false),
                                                                                                        (11, false, 0, 0.169, '2022-01-01',100, false),
                                                                                                        (11, false, 0, 0.169, '2022-01-01',100, false),
                                                                                                        (11, false, 0, 0.169, '2022-01-01',100, false),
                                                                                                        (12, false, 0, 0.65, '2022-01-01',100, false),
                                                                                                        (12, false, 0, 0.65, '2022-01-01',100, false);


insert into book_genre  (book_id, genre_id) values (1, 4), (2, 4), (3, 8), (4, 3), (5, 6), (6, 1), (7, 2), (8, 8), (9, 2), (10, 11), (11, 10), (12, 3);

insert into book_author  (book_id, author_id) values (1, 2), (2, 8), (3, 4), (4, 9), (5, 7), (6, 17), (7, 16), (8, 15), (9, 14), (10, 12), (11, 10), (12, 13);