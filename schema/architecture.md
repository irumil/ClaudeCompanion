# ClaudeCompanion - Архитектура приложения

## Общая схема

```mermaid
graph TB
    subgraph "Браузер Firefox"
        EXT[Browser Extension<br/>background.js]
        CLAUDE[Claude.ai Website<br/>https://claude.ai]
    end

    subgraph "Desktop Application"
        TRAY[System Tray Icon<br/>tray.go]
        HTTP[HTTP Server<br/>:8383<br/>server.go]
        API[API Client<br/>curl wrapper<br/>api/client.go]
        CFG[Config Manager<br/>config.yaml<br/>hot-reload]
        NOTIF[Notifier<br/>Windows Toast<br/>notifier.go]
        LOGGER[Logger<br/>file + console<br/>logger.go]
        ICON[Icon Generator<br/>Dynamic PNG→ICO<br/>icon.go]
    end

    subgraph "External"
        PROXY[Proxy Server<br/>Optional<br/>your-proxy:port]
        CLAUDEAPI[Claude.ai API<br/>/api/organizations/{uuid}/usage]
    end

    %% Browser Extension Flow
    CLAUDE -->|1. User visits| EXT
    EXT -->|2. Extract sessionKey cookie| EXT
    EXT -->|3. Fetch org UUID| CLAUDEAPI
    EXT -->|4. POST /set-context<br/>cookies + targetUrl| HTTP

    %% Desktop App Flow
    HTTP -->|5. Store context| API
    HTTP -.->|Update URL| TRAY

    CFG -->|Configure| API
    CFG -->|Configure| HTTP
    CFG -->|Configure| NOTIF
    CFG -->|Configure| LOGGER
    CFG -->|Browser path| TRAY

    API -->|6. Every 30s<br/>curl with proxy| PROXY
    PROXY -->|7. Forward request| CLAUDEAPI
    CLAUDEAPI -->|8. Return usage data<br/>five_hour + seven_day| PROXY
    PROXY -->|9. Return response| API

    API -->|10. Parse usage| TRAY
    TRAY -->|11. Generate icon<br/>with percentage| ICON
    ICON -->|12. Update tray| TRAY

    API -->|13. Check thresholds| NOTIF
    NOTIF -->|14. Show Toast<br/>if low/zero| NOTIF

    TRAY -->|15. Click menu| TRAY
    TRAY -->|Open settings| CFG
    TRAY -->|Open Claude.ai| CLAUDE

    API -.->|All operations| LOGGER
    HTTP -.->|All operations| LOGGER
    TRAY -.->|All operations| LOGGER
    NOTIF -.->|All operations| LOGGER

    style EXT fill:#4CAF50
    style TRAY fill:#2196F3
    style API fill:#FF9800
    style CLAUDEAPI fill:#9C27B0
    style PROXY fill:#607D8B
    style CFG fill:#FFC107
    style NOTIF fill:#E91E63
```

## Компоненты системы

### 1. Browser Extension (Firefox)
**Файл:** `extension/background.js`

**Функции:**
- Перехватывает cookie `sessionKey` с сайта Claude.ai
- Получает UUID организации через API `/api/organizations`
- Формирует URL для получения usage: `/api/organizations/{uuid}/usage`
- Отправляет данные на локальный HTTP сервер

**Технологии:**
- Firefox WebExtensions API
- browser.cookies, browser.tabs, browser.storage

### 2. Desktop Application (Go)
**Исполняемый файл:** `dist/claudecompanion.exe`

#### 2.1 HTTP Server (`internal/server/server.go`)
- Слушает на порту 8383 (конфигурируемо)
- Endpoint: `POST /set-context` - получает cookies и targetURL от расширения
- Endpoint: `GET /health` - проверка работоспособности
- Передает данные в API Client

#### 2.2 API Client (`internal/api/client.go`)
- Хранит контекст (cookies, targetURL)
- Выполняет запросы через системный curl с поддержкой прокси
- Скрывает консольное окно curl (Windows-specific)
- Парсит JSON ответ с данными usage
- Возвращает структуру UsageResponse с полями:
  - `five_hour` - использование за 5 часов
  - `seven_day` - использование за 7 дней

#### 2.3 System Tray Manager (`internal/tray/tray.go`)
- Создает иконку в системном трее
- Генерирует динамическую иконку с процентами (0-100)
- Цветовая индикация:
  - 🟢 Зеленый: > 40%
  - 🟡 Желтый: 20-40%
  - 🔴 Красный: < 20%
  - ⚪ Серый: ошибка подключения
- Контекстное меню:
  - "Открыть Claude.ai" - открывает браузер
  - "Открыть настройки" - открывает config.yaml в notepad
  - "Выход" - закрывает приложение
- Tooltip с информацией о квоте

#### 2.4 Icon Generator (`internal/icon/generator.go`)
- Генерирует иконки 48x48 пикселей
- Рисует цифры пиксель-арт стилем (6x6 блоки)
- Прозрачный фон
- Конвертирует PNG → ICO формат для Windows

#### 2.5 Notifier (`internal/notifier/notifier.go`)
- Windows Toast уведомления через go-toast
- Типы уведомлений:
  - **Низкая квота**: при пороге ≤20% (конфигурируемо)
  - **Квота исчерпана**: при 0% с временем восстановления
  - **Ошибка авторизации**: после N неудачных запросов
- Случайные фразы из конфигурации
- Иконка приложения в уведомлениях
- Состояние уведомлений (не показывать дубликаты)

#### 2.6 Config Manager (`internal/config/config.go`)
- Загрузка конфигурации из `config.yaml`
- Hot-reload: отслеживание изменений файла каждые 2 секунды
- Настройки:
  - Порт HTTP сервера
  - Интервал опроса API
  - Адрес прокси сервера
  - Пороги для уведомлений и серого режима
  - Путь к браузеру
  - Включение/выключение логирования в файл
  - Настройки demo mode
  - Фразы для уведомлений

#### 2.7 Logger (`internal/logger/logger.go`)
- Логирование в файл и/или консоль (конфигурируемо)
- Файлы с ротацией по дате: `claudecompanion-YYYY-MM-DD.log`
- Уровни: INFO, DEBUG, WARNING, ERROR, FATAL
- Расположение: рядом с exe файлом

### 3. External Services

#### 3.1 Proxy Server
- Опционально: настраивается в config.yaml
- Используется для всех curl запросов к Claude.ai API

#### 3.2 Claude.ai API
- **Organizations endpoint:** `GET https://claude.ai/api/organizations`
  - Возвращает список организаций с UUID
- **Usage endpoint:** `GET https://claude.ai/api/organizations/{uuid}/usage`
  - Требует cookie: `sessionKey`
  - Возвращает JSON с utilization и resets_at

## Поток данных

### Инициализация (один раз)
```
1. User opens Claude.ai in Firefox
2. Extension detects page load
3. Extension extracts sessionKey cookie
4. Extension calls /api/organizations to get UUID
5. Extension constructs usage URL
6. Extension sends POST to localhost:8383/set-context
7. Desktop app stores cookies and URL
```

### Polling Loop (каждые 30 секунд)
```
1. Desktop app calls curl with:
   - URL: https://claude.ai/api/organizations/{uuid}/usage
   - Cookie: sessionKey={value}
   - Proxy: optional (if configured)
2. Curl executes (hidden window on Windows)
3. Response parsed as JSON
4. Calculate remaining quota: 100 - utilization
5. Update tray icon with color and percentage
6. Update tooltip with detailed info
7. Check thresholds and show notifications if needed
8. Log all operations (if enabled)
```

### User Actions
```
- Right-click tray icon → Show context menu
- Click "Открыть Claude.ai" → Open browser with targetURL
- Click "Открыть настройки" → Open config.yaml in notepad
- Click "Выход" → Graceful shutdown
```

## Конфигурация

### config.yaml
```yaml
server_port: 8383
poll_interval_seconds: 30
use_curl_fallback: true
gray_mode_threshold: 5        # N errors before gray icon
notification_threshold: 10    # N errors before notification
proxy: ""                     # Optional: "http://your-proxy:port"
enable_file_logging: false    # true = file+console, false = console only
browser_path: "C:\\Program Files\\Mozilla Firefox\\firefox.exe"

low_value_notifications:
  enabled: true
  threshold: 20               # Show notification when ≤20%
  phrases:
    - "Пора идти домой! 🏡"
    - "Время отдохнуть! ☕"
  zero_phrases:
    - "Всё, капут! 💥"
    - "Game over! 🎮"

demo_mode:
  enabled: false              # For testing: simulates declining quota
  duration_seconds: 60        # Full cycle: 100% → 0% in 60 seconds
```

## Технологии

### Desktop App
- **Язык:** Go 1.21+
- **Библиотеки:**
  - `github.com/getlantern/systray` - system tray
  - `github.com/go-toast/toast` - Windows Toast notifications
  - `gopkg.in/yaml.v3` - YAML parsing
  - `image/*` - icon generation
- **Build:** `-ldflags "-H windowsgui"` - no console window
- **Resources:** rsrc for embedding icon.ico

### Browser Extension
- **Платформа:** Firefox WebExtensions
- **Manifest:** v2
- **API:** browser.* (Firefox-specific)

### External Tools
- **curl.exe** - системный curl для HTTP запросов
- **notepad.exe** - открытие настроек (Windows)

## Deployment

### Файловая структура
```
dist/
├── claudecompanion.exe     # Main executable (GUI, no console)
├── config.yaml             # Configuration (created on first run)
├── icon.ico                # Icon for notifications (optional)
└── claudecompanion-YYYY-MM-DD.log  # Log files (if enabled)
```

### Установка
1. Скопировать `claudecompanion.exe` в любую папку
2. Установить расширение в Firefox
3. Запустить exe файл
4. Открыть Claude.ai в Firefox
5. Приложение автоматически получит cookies и начнет работу

## Особенности Windows

1. **Скрытие консоли:**
   - Приложение: `-H windowsgui`
   - Curl: `CREATE_NO_WINDOW` flag

2. **Иконка в exe:**
   - Встроена через rsrc.syso
   - Отображается в проводнике и на панели задач

3. **Toast уведомления:**
   - Нативные Windows 10/11 Toast
   - С иконкой приложения
   - Кликабельные

4. **Многострочные tooltip:**
   - Используют `\r\n` вместо `\n`

5. **Пути к файлам:**
   - Относительно exe файла
   - Поддержка пробелов в путях (через кавычки)
