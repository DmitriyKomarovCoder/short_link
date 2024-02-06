# Тестовое задание для стажера-разработчика 
## Задача



Реализовать сервис, предоставляющий API по созданию сокращённых ссылок.



### Ссылка должна быть:

—Уникальной; на один оригинальный URL должна ссылаться только одна сокращенная ссылка;

—Длиной 10 символов;

—Из символов латинского алфавита в нижнем и верхнем регистре, цифр и символа _ (подчеркивание).



### Сервис должен быть написан на Go и принимать следующие запросы по http:

1. Метод Post, который будет сохранять оригинальный URL в базе и возвращать сокращённый.

2. Метод Get, который будет принимать сокращённый URL и возвращать оригинальный.

#### Условие со звёздочкой (будет большим плюсом):

Сделать работу сервиса через GRPC, то есть составить proto и реализовать сервис с двумя соответствующими эндпойнтами





### Решение должно соответствовать условиям:

— Сервис распространён в виде Docker-образа;

— В качестве хранилища ожидаем in-memory решение и PostgreSQL. Какое хранилище использовать, указывается параметром при запуске сервиса;

— Реализованный функционал покрыт Unit-тестами.


