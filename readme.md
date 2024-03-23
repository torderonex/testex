# Записка к тестовому заданию

В ходе выполнения задания был полностью реализован требуемый функционал. Сервис предоставляет возможность

- Создавать новые команды
- Получать список всех команд
- Получать информацию об одной команде
- Исполнять команду
- Останавливать выполнение команды
- Следить за логами команды
- Параллельно запускать неограниченное кол-во команд
  Сервис покрыт тестами, также имеет настроенный пайплайн в Gitlab CI и упакован в докер вместе с базой данных. В дополнение добавлена поддержка Windows. Систему можно выставить в конфигурационном файле (UNIX по умолчанию).

Прим. (config.yaml) :

```yaml
os: "win"
```

## Инструкция по запуску

1. Установите Docker, Docker-compose
2. Соберите приложение при помощи:

```bash
make build && make run
```

или

```bash
docker-compose build testex
```

## API

### Add Command

- **URL**: `/commands/add`
- **Method**: `POST`
- **Description**: Добавляет новую команду. Вовращает id добавленной команды.
- **Request Body**: `{ "alias": "string", "script": "string" }`
- **Response**: `{ "id": 1 }`

### Execute Command

- **URL**: `/commands/execute`
- **Method**: `POST`
- **Description**: Выполняет команду по псевдониму. Возвращает id выполняющейся команды.
- **Request Body**:
  `{   "alias": "string" }`
- **Response**:
  `{   "id": "int" }`

### Get Command

- **URL**: `/commands/{alias}`
- **Method**: `GET`
- **Description**: Возвращает информацию о команде.
- **URL Parameters**:
  - `alias`: Псеводним для команды
- **Response**:
  `{"id": "int", "alias": "string", "script": "string" }`

### Get All Commands

- **URL**: `/commands`
- **Method**: `GET`
- **Description**: Возвращает информацию обо всех командах.
- **Response**: Массив всех команд:
  `[{"id": "int","alias": "string", "script": "string" }, ... ]`

### Stop Command

- **URL**: `/commands/stop`
- **Method**: `POST`
- **Description**: Останавливает выполнение команды.
- **Request Body**:
  `{   "id": "int" }`

### Get Logs

- **URL**: `/commands/logs/{id}`
- **Method**: `GET`
- **Description**: Возвращает логи исполняемой команды.
- **URL Parameters**:
  - `id`: ID исполняемой команды.
- **Response**: Array of log objects:
  `[ {"id": "int", "message": "string", "executed_command_id" : "int" }, ... ]`

# Используемые технологии

- Docker, Docker Compose
- Windows, Linux Ubuntu
- Gitlab-CI
- Make
  Из сторонних пакетов Go использовались:
- viper (для файлов конфигурации)
- slqx, pq (для работы с Postgres)
- testify (для удобного тестирования)
