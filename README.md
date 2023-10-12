# Практичні 2-го етапу курсу з Go

|  # | Короткий опис                                                  | Умова | Рішення |
|---:|:---------------------------------------------------------------|-------|---------|
|  1 | Вивести усіх користувачів таблицею у консоль через `Printf`... | [#1](https://github.com/grescher/goc-s2-psets/issues/1) | [p1.0](https://github.com/grescher/goc-s2-psets/releases/tag/p1.0) |
|  2 | До таблиці додати перелік книжок, порахувати середній вік...   | [#2](https://github.com/grescher/goc-s2-psets/issues/2) | [p2.0](https://github.com/grescher/goc-s2-psets/releases/tag/p2.0) |
|  3 | Відсортувати користувачів за сумою середн. віку читача кожної книжки ... | [#5](https://github.com/grescher/goc-s2-psets/issues/5) | [p3.0](https://github.com/grescher/goc-s2-psets/releases/tag/p3.0) |
|  4 | Типізовані константи: додати поле "тип користувача"...         | [#6](https://github.com/grescher/goc-s2-psets/issues/6) | [p4.0](https://github.com/grescher/goc-s2-psets/releases/tag/p4.0) |
|  5 | Декодування даних з бінарного буферу                           | [#8](https://github.com/grescher/goc-s2-psets/issues/8) | [p5.0](https://github.com/grescher/goc-s2-psets/releases/tag/p5.0) |
|  6 | Спільний індекс — стан активності користувачів                 | [#10](https://github.com/grescher/goc-s2-psets/issues/10) | [p6.0](https://github.com/grescher/goc-s2-psets/releases/tag/p6.0) |
|  7 | Збереження у файл. Елементарний інтерпретатор на `Scanf`       | [#12](https://github.com/grescher/goc-s2-psets/issues/12) | [p7.0](https://github.com/grescher/goc-s2-psets/releases/tag/p7.0) |
|  8 | Розпізнавання виразів з цифрами. Буферизований ввод.           | [#13](https://github.com/grescher/goc-s2-psets/issues/13) | |
|  9 | Пересилання даних між клієнтом і сервером                      | [#14](https://github.com/grescher/goc-s2-psets/issues/14) | |
| 10 | http-сервер, який отримує дані з БД (API)                      | [#15](https://github.com/grescher/goc-s2-psets/issues/15) | |
| 12 | Запустити наш застосунок в Докері (Dockerfile)                 | [#17](https://github.com/grescher/goc-s2-psets/issues/17) | |

## Запуск

```term
go run . [FILE]

# or

./app [FILE]
```

Під час виклику запуску, як зі скомпільованого файла, так і з `go run`, перший аргумент — назва файлу з даними. Якщо конкретний файл не вказано, то застосунок за замовчанням створить порожній файл, який буде знаходитись за відносною адресою `./datafiles/test.database`.
