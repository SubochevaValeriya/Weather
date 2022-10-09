CREATE TABLE subscription
(
    id serial not null unique,
    city varchar(255) not null unique primary key,
    subscription_date date not null
);

CREATE TABLE weather

(
    id serial not null unique primary key,
    city_id int not null references subscription(id),
    city varchar(255) not null,
    weather_date date not null,
    weather int
);