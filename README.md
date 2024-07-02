# Тестовое задание Efficent Mobile 2024

## Сервис тайм-трекинга Go,PostgreSQL

##### Автор: [Виноградов Данил](https://t.me/japsty)
##### Ссылка на задание [тут](https://github.com/Japsty/EM_test_task_2/blob/main/task.md)
##### Ссылка на [Postman](https://www.postman.com/japsty/workspace/em-test-task/collection/29141861-d601e92d-12c2-45a7-b128-c77592806717?action=share&creator=29141861) коллекцию
##### Ссылка на yaml файл [Swagger](https://github.com/Japsty/EM_test_task_2/blob/main/docs/swagger.yaml)

### Запуск

Запуск из корневой директории проекта
```
docker compose --env-file .env up
```

Для мока API использовал Prism
```
prism mock mockAPI.yaml -h 0.0.0.0  
```
Хендлеры практически полностью покрыты тестами.
Репозитории покрыты тестами только на успешное выполнение.

Билд в два этапа, чтобы не тянуть весь проект в образ.
