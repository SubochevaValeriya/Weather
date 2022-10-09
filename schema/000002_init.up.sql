CREATE TABLE subscription
(
    id serial not null primary key unique,
    city varchar(255) not null unique,
    subscription_date date not null
);

CREATE TABLE weather

(
    id serial not null unique primary key,
    city_id int not null references subscription(id),
    weather_date date not null,
    weather float
);

CREATE TABLE weather_archive

(
    id serial not null unique primary key,
    city_id int not null references subscription(id),
    weather_date date not null,
    weather float
);