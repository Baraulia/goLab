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
    adress varchar(255),
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
    condition int not null
);

CREATE TYPE status AS ENUM ('open', 'closed');

CREATE TABLE issue_act
(
    id serial not null unique primary key,
    user_id int references users(id) not null,
    listbook_id int references list_books(id) not null,
    rental_time int not null,
    return_date timestamp with time zone not null,
    pre_cost decimal not null,
    cost decimal not null,
    status status not null
);

CREATE TABLE return_act
(
    id serial not null unique primary key,
    issue_act_id int references issue_act(id) not null,
    return_date timestamp with time zone not null,
    foto varchar(255) array[5],
    fine decimal,
    condition_decrese int,
    rating int
);
insert into genre (genre_name) values ('Novel2'), ('Fantasy'), ('Detective'), ('Advanture'), ('Erotic'), ('Triller'), ('Philosophical'), ('Satire'), ('Comedy'), ('Crime'), ('Horror'), ('Business');

insert into authors (author_name, author_foto) values ('Редьярд Киплинг', 'Path1'),
                                                            ('Марк Твен', 'Path2'),
                                                            ('Джордж Оруэлл', 'Path3'),
                                                            ('Максим Горький', 'Path4'),
                                                            ('Александр Куприн', 'Path5'),
                                                            ('Иван Бунин', 'Path6'),
                                                            ('Томас Манн', 'Path7'),
                                                            ('Джек Лондон', 'Path8'),
                                                            ('Франц Кафка', 'Path9'),
                                                            ('Борис Пастернак', 'Path10'),
                                                            ('Агата Кристи', 'Path11'),
                                                            ('Михаил Булгаков', 'Path12'),
                                                            ('Эрнест Хемингузй', 'Path13'),
                                                            ('Антуан де Сент-Экзюпери', 'Path14'),
                                                            ('Бернард Шоу', 'Path15'),
                                                            ('Артур Конан Дойл', 'Path16'),
                                                            ('Эмиль Золя', 'Path17');

insert into users (surname, user_name, patronymic, pasp_number, email, adress, birth_date) values ('Барауля', 'Сергей', 'Михайлович', '123456', 'baraulia@yandex.ru', 'Minsk', '1965-07-20'),
                                                                                                  ('Иванов', 'Андрей', 'Александрович', '123457', 'baraulia1@yandex.ru', 'Pemza', '1990-08-15'),
                                                                                                  ('Петров', 'Илья', 'Дмитриевич', '123476', 'baraulia2@yandex.ru', 'Tagil', '1915-02-15'),
                                                                                                  ('Сидоров', 'Борис', 'Ильич', '127456', 'baraulia3@yandex.ru', 'Gomel', '2006-06-26'),
                                                                                                  ('Ульянов', 'Андрей', 'Михайлович', '173456', 'baraulia4@yandex.ru', 'Vitebsk', '1995-09-20'),
                                                                                                  ('Чаушеску', 'Сергей', 'Сидорович', '723456', 'baraulia5@yandex.ru', 'Donetsk', '1975-10-11'),
                                                                                                  ('Рэмбо', 'Борис', 'Андреевич', '823456', 'baraulia6@yandex.ru', 'Moskow', '1983-03-18'),
                                                                                                  ('Сталин', 'Алексей', 'Сергеевич', '183456', 'baraulia7@yandex.ru', 'Kolodoschi', '1999-01-09'),
                                                                                                  ('Чиполино', 'Дмитрий', 'Алексеевич', '128456', 'baraulia8@yandex.ru', 'Grodno', '2000-05-05'),
                                                                                                  ('Мариарти', 'Андрей', 'Дмитриевич', '123856', 'baraulia9@yandex.ru', 'Volgograd', '2010-12-07'),
                                                                                                  ('Сергеев', 'Сергей', 'Ильич', '123486', 'baraulia10@yandex.ru', 'Vladivostok', '2002-11-02'),
                                                                                                  ('Никитин', 'Илья', 'Сергеевич', '123458', 'baraulia11@yandex.ru', 'Vileyka', '1998-10-22'),
                                                                                                  ('Тиньков', 'Сидор', 'Александрович', '123459', 'baraulia12@yandex.ru', 'Tombov', '1983-09-12');

insert into books  (book_name, cost, cover, published, pages, amount) values ('Том сойер', 50, 'Pathbook1', 2021, 456, 2),
                                                                             ('Белый клык', 65, 'Pathbook2', 2020, 560, 1),
                                                                             ('Нельзя молчать!', 32, 'Pathbook3', 2018, 400, 2),
                                                                             ('Процесс', 50, 'Pathbook4', 2021, 288, 3),
                                                                             ('Волшебная гора', 25, 'Pathbook5', 2019, 928, 1),
                                                                             ('Деньги', 50, 'Pathbook6', 2012, 512, 2),
                                                                             ('Конан-варвар', 38, 'Pathbook7', 2021, 680, 3),
                                                                             ('Цезарь и Клеопатра', 22, 'Pathbook8', 2021, 412, 1),
                                                                             ('Маленький принц', 15, 'Pathbook9', 2014, 220, 2),
                                                                             ('Мастер и Маргарита', 62, 'Pathbook10', 2021, 576, 2),
                                                                             ('Доктор Живаго', 13, 'Pathbook11', 2021, 608, 3),
                                                                             ('По ком звонит колокол', 50, 'Pathbook1', 2019, 640, 2);

insert into list_books  (book_id, issued, rent_number, rent_cost, reg_date, condition) values (1, false, 0, 1.15, '2022-01-01',100),
                                                                                              (1, false, 0, 1.15, '2022-01-01',100),
                                                                                              (2, false, 0, 1.50, '2022-01-01',100),
                                                                                              (3, false, 0, 0.74, '2022-01-01',100),
                                                                                              (3, false, 0, 0.74, '2022-01-01',100),
                                                                                              (4, false, 0, 1.15, '2022-01-01',100),
                                                                                              (4, false, 0, 1.15, '2022-01-01',100),
                                                                                              (4, false, 0, 1.15, '2022-01-01',100),
                                                                                              (5, false, 0, 0.58, '2022-01-01',100),
                                                                                              (6, false, 0, 1.15, '2022-01-01',100),
                                                                                              (6, false, 0, 1.15, '2022-01-01',100),
                                                                                              (7, false, 0, 0.87, '2022-01-01',100),
                                                                                              (7, false, 0, 0.87, '2022-01-01',100),
                                                                                              (7, false, 0, 0.87, '2022-01-01',100),
                                                                                              (8, false, 0, 0.51, '2022-01-01',100),
                                                                                              (9, false, 0, 0.35, '2022-01-01',100),
                                                                                              (9, false, 0, 0.35, '2022-01-01',100),
                                                                                              (10, false, 0, 1.43, '2022-01-01',100),
                                                                                              (10, false, 0, 1.43, '2022-01-01',100),
                                                                                              (11, false, 0, 0.30, '2022-01-01',100),
                                                                                              (11, false, 0, 0.30, '2022-01-01',100),
                                                                                              (11, false, 0, 0.30, '2022-01-01',100),
                                                                                              (12, false, 0, 1.15, '2022-01-01',100),
                                                                                              (12, false, 0, 1.15, '2022-01-01',100);


insert into book_genre  (book_id, genre_id) values (1, 4), (2, 4), (3, 8), (4, 3), (5, 6), (6, 1), (7, 2), (8, 8), (9, 2), (10, 11), (11, 10), (12, 3);

insert into book_author  (book_id, author_id) values (1, 2), (2, 8), (3, 4), (4, 9), (5, 7), (6, 17), (7, 16), (8, 15), (9, 14), (10, 12), (11, 10), (12, 13);