# Log Parser 3000

Микросервис на Go для парсинга лог-файлов InfiniBand, агрегации топологии и хранения данных в PostgreSQL.

---

## Требования
- Go 1.25
- Docker
- Docker Compose
- Архив с логами (`.zip` или `.tar.gz`) — положить в папку `data/`
- Архив должен содержать файл `ibdiagnet2.db_csv`
---

## Быстрый старт

```bash
# 1. Клонировать репозиторий
git clone https://github.com/Davidianol/Log-Parser-3000.git
cd log_parser3000

# 2. Положить архив с логами
mkdir -p data
cp /path/to/log.zip data/

# 3. Поднять сервисы
docker compose up -d

# 4. Проверить что сервис запустился
curl -s http://localhost:8080/api/log/1 | jq
```

В проекте используются миграции: после первого запуска появятся схемы и таблицы в постгрес, после перезапусков повторных миграций не будет.

---

## Переменные окружения

| Переменная     | По умолчанию                                                     | Описание                             |
|----------------|------------------------------------------------------------------|--------------------------------------|
| `DATABASE_URL` | `postgres://postgres:postgres@db:5432/logparser?sslmode=disable` | Строка подключения к PostgreSQL      |
| `PORT`         | `8080`                                                           | Порт HTTP-сервера                    |
| `LOG_LEVEL`    | `info`                                                           | Уровень логов: `debug/info/warn/error` |

Переопределить можно через `.env`.

Пример:
```env
PORT=9090
LOG_LEVEL=debug
```

---

## API

### POST /api/parse/

Парсит архив из папки `data/` и сохраняет результат в БД. Поддерживаются форматы `.zip` и `.tar.gz`. Архив должен содержать файл `ibdiagnet2.db_csv`.

**Тело запроса:**
```json
{ "path": "log.zip" }
```

**Ответ 200:**
```json
{ "log_id": 1 }
```

**Ответ при ошибочном файле (422):**
```json
{ "log_id": 2, "error": "START_PORTS section not found or empty" }
```

```bash
curl -s -X POST http://localhost:8080/api/parse/ \
  -H "Content-Type: application/json" \
  -d '{"path": "log.zip"}' | jq
```

> Повторная отправка одного и того же файла создаёт новую запись лога с новым `log_id`.

---

### GET /api/log/{log_id}

Мета-информация о логе: статус, количество узлов, дата загрузки.

**Ответ 200:**
```json
{
  "id": 1,
  "filename": "log.zip",
  "status": "done",
  "node_count": 5,
  "uploaded_at": "2026-05-14T22:00:00Z",
  "error_message": null
}
```

```bash
curl -s http://localhost:8080/api/log/1 | jq
```

Возможные значения `status`: `done`, `error`.

---

### GET /api/topology/{log_id}

Список узлов и групп топологии для данного лога.

**Ответ 200:**
```json
{
  "Nodes": [
    {
      "id": 1,
      "log_id": 1,
      "node_guid": "0xb8599f0300ebd7a0",
      "node_desc": "switch-leaf-01",
      "node_type": 2,
      "num_ports": 40,
      "system_image_guid": "0xb8599f0300ebd7a0",
      "port_guid": "0xb8599f0300ebd7a0"
    }
  ],
  "Groups": [
    {
      "Key": "host_active",
      "NodeGUIDs": ["0xabc123", "0xdef456"]
    },
    {
      "Key": "switch_active",
      "NodeGUIDs": ["0xb8599f0300ebd7a0"]
    }
  ]
}
```

```bash
curl -s http://localhost:8080/api/topology/1 | jq
```

---

### GET /api/node/{node_id}

Детальная информация об узле по его ID.

**Ответ 200:**
```json
{
  "id": 1,
  "log_id": 1,
  "node_guid": "0xb8599f0300ebd7a0",
  "node_desc": "switch-leaf-01",
  "node_type": 2,
  "num_ports": 40,
  "system_image_guid": "0xb8599f0300ebd7a0",
  "port_guid": "0xb8599f0300ebd7a0"
}
```

```bash
curl -s http://localhost:8080/api/node/1 | jq
```

---

### GET /api/port/{node_id}

Все порты узла по его ID.

**Ответ 200:**
```json
[
  {
    "id": 1,
    "node_id": 1,
    "node_guid": "0xb8599f0300ebd7a0",
    "port_num": 1,
    "port_state": 4,
    "port_phy_state": 5,
    "lid": 1,
    "link_speed_actv": 2048,
    "link_width_actv": 2
  }
]
```

```bash
curl -s http://localhost:8080/api/port/1 | jq
```

---

## Топология и возможные связи (F-3)

Чтобы построить полный граф с ребрами, нужны поля данные о соседях (`NeighborNodeGUID`, `PeerPortNum`). Поэтому в нашем случае невозможно его построить, таких полей у нас нет.

**Текущая реализация — группировка по состоянию портов:**

Узел считается активным, если хотя бы один из его портов имеет `PortState == 4` (Active в InfiniBand). На этой основе формируются группы в ответе `/api/topology/{log_id}`:

| Группа | Условие |
|--------|---------|
| `host_active` | `node_type == 1` и есть порт с `port_state == 4` |
| `host_isolated` | `node_type == 1` и нет активных портов |
| `switch_active` | `node_type == 2` и есть порт с `port_state == 4` |
| `switch_isolated` | `node_type == 2` и нет активных портов |

---

## Структура проекта

```
cmd/
└── main.go                   # точка входа, миграции, DI
internal/
├── domain/                   # доменные типы (Log, Node, Port, NodeInfo, TopologyGroup)
├── parser/                   # парсер архивов и ibdiagnet2.db_csv
├── repository/
│   ├── repository.go         # интерфейс репозитория
│   └── postgres/             # реализация для PostgreSQL
├── service/                  # бизнес-логика
└── handler/
    ├── handler.go            # HTTP-хендлеры
    └── router.go             # маршрутизация + middleware логирования
migrations/
├── 001_init.up.sql
└── 001_init.down.sql
data/                         # монтируется в контейнер, сюда кладутся архивы
docker-compose.yml
Dockerfile
README.md
```

---

## Отказоустойчивость

- Невалидный архив (не zip/tar.gz, повреждённые байты) -> HTTP `400` с описанием ошибки.
- Ошибочный лог (отсутствующие обязательные секции, мусор в числовых полях, обрезанные строки) -> HTTP `422` с `log_id` и `error`. Запись сохраняется в таблице `logs` со статусом `error`, что позволяет отслеживать историю через `GET /api/log/{log_id}`.
- Сервис стартует только после успешного `healthcheck` базы данных (`pg_isready`).
- Путь к архиву проверяется на `path traversal` (`..`) — попытка выйти за пределы `data/` возвращает `400`.
- Все операции с БД при парсинге выполняются в транзакции (частичная запись невозможна).

## postman_collection

В проекте есть postman_collection.json, в котором можно проверить работу API и многие краевые случаи. Не зря же у нас в `data` 11 архивов :)

Положите архив log.zip в `data/`. Этот архив был запаролен в ТЗ, поэтому не буду его публично выкладывать.

