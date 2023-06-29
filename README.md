# Lecture 1. Data Formatting

Вивести усіх користувачів через Printf у консоль таблицею таким чином, щоби
усі колонки мали однаковий розмір (за найдовшою клітинкою).

Значення у полях знаходяться у різних діапазонах, на ваш розсуд - придумати
як читабельніше вивести їх.

## Example

```md
|                             Name                             |    Age     | Active |       Mass       |
|--------------------------------------------------------------|------------|--------|------------------|
| John Doe                                                     | 30         | true   | 80.0             |
| Jake Doe                                                     | 20         | false  | 60.0             |
|  Jane Doe                                                    | 150        | true   | 0.75             |
| \t                                                           | -10        | true   | 8000.0           |
| Vm0weE5GVXhUblJWV0dSUFZtMW9WVll3WkRSV1ZteDBaRVYwVmsxWGVGWlZi | 0          | true   | 0.0              |
| VEZIWVd4S2MxTnNiR0ZXVm5Cb1ZsVmFWMVpWTVVWaGVqQTk=\nVm0weE5GVX |            |        |                  |
| hUblJWV0dSUFZtMW9WVll3WkRSV1ZteDBaRVYwVmsxWGVGWlZiVEZIWVd4S2 |            |        |                  |
| MxTnNiR0ZXVm5Cb1ZsVmFWMVpWTVVWaGVqQTk=                       |            |        |                  |
| \x00\x10\x20\x30\x40\x50\x60\x70                             | 0          | true   | 0.0              |
| Billy Bones                                                  | -130000    | false  | 3141567.98765457 |
| Billy Bones Jr.                                              | 1234567890 | true   | 3141567.98765457 |
```
