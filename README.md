# Yalantis test task

## Task conditions

Создать REST сервис, в котором пользователь может:

- зарегистрировать учетную запись через login, password. авторизоваться;
- заполнить, изменить свой профиль;
- просмотреть профиля других зарегистрированных пользователей;

Бонус:

- логирование ошибок;
- пользователь может установить себе аватар. если изображение больше 160x160px сервис должен её сжать;
- тесты или код который покрывается тестами с минимальной модификацией (interfaces).

На усмотрение кандидата: DB(RDBMS), ORM, images storage, router, libs

## Notes and sources

The application is modified and dockerized task from medium.com Go sources of REST api draft example

## Build and start

``` bash
docker-compose up -d
```

## Usage examples

### register new user

``` bash
curl -X POST http://localhost:8080/api/user/new -H 'Content-Type: application/json' -d '{"email": "test.user@gmail.com", "password": "mypassword"}'
```

### login with test user

``` bash
curl -X POST http://localhost:8080/api/user/login -H 'Content-Type: application/json' -d '{"email": "test.user@gmail.com", "password": "mypassword"}'
```

### use token to get user by id

``` bash
curl http://localhost:8080/api/user/1 -H 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJVc2VySWQiOjJ9.DzOJ7GHkPwiDE3T78dFMriY96VwzytQSBV7-c64dxx8'
```

### use token to get registered users list

``` bash
curl http://localhost:8080/api/users -H 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJVc2VySWQiOjJ9.DzOJ7GHkPwiDE3T78dFMriY96VwzytQSBV7-c64dxx8'
```

### update user's info

Allowed fields to update:

- email
- password
- about

#### update examples:

``` bash
curl -X PUT http://localhost:8080/api/user/1/update -H 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJVc2VySWQiOjJ9.DzOJ7GHkPwiDE3T78dFMriY96VwzytQSBV7-c64dxx8' -d '{"about": "Hello!"}'
```

few params:

``` bash
curl -X PUT http://localhost:8080/api/user/1/update -H 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJVc2VySWQiOjJ9.DzOJ7GHkPwiDE3T78dFMriY96VwzytQSBV7-c64dxx8' -d '{"about": "Hello, World!", "email": "test.guy@gmail.com"}'
```
