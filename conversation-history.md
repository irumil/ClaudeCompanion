# История разработки ClaudeCompanion

Полная история диалога по разработке приложения ClaudeCompanion - системного трей-приложения для мониторинга квоты Claude.ai API.

---

## Сессия 1: Начальная разработка и основные фичи

### 1. Исправление меню "Открыть Claude.ai"

**Запрос пользователя:**
> пункт меню "открыть Claude.ai" если если куки то открывает /usage, нужно открывать просто claude.ai

**Решение:**
- Файл: `internal/tray/tray.go`
- Изменено: `t.openURL(t.targetURL)` → `t.openURL("https://claude.ai")`
- Результат: Меню теперь открывает базовый URL Claude.ai вместо endpoint /usage

---

### 2. Добавление фичи "Утренний привет Клоду"

**Запрос пользователя:**
> Добавим новую фичу. "Оптимизация потребления 5-ти часового лимита".
>
> Настройки:
> - greeting_cron - cron расписание когда отправлять
> - greeting_text - текст сообщения
> - greeting_chat_id - UUID чата куда отправлять
>
> При срабатывании cron отправлять curl POST запрос с сообщением.
> После отправки показывать уведомление.

**Реализация:**

#### 2.1 Добавлена зависимость Cron
- Файл: `go.mod`
- Добавлено: `github.com/robfig/cron/v3 v3.0.1`

#### 2.2 Обновлено браузерное расширение
- Файл: `extension/background.js`
- Изменения:
  - Функция `getOrgData()` теперь возвращает как `organizationId`, так и `usageUrl`
  - Отправляется `organizationId` в десктопное приложение
  - Добавлено логирование для отладки

```javascript
const orgData = await getOrgData(sessionKey);
const payload = {
  cookies: `sessionKey=${sessionKey}`,
  targetUrl: orgData.usageUrl,
  organizationId: orgData.organizationId
};
```

#### 2.3 Обновлён HTTP сервер
- Файл: `internal/server/server.go`
- Изменения:
  - Добавлено поле `OrganizationID` в структуру `ContextData`
  - Callback теперь принимает `organizationID` как параметр
  - Сервер передаёт `organizationID` в callback

```go
type ContextData struct {
    Cookies        string `json:"cookies"`
    TargetURL      string `json:"targetUrl"`
    OrganizationID string `json:"organizationId"`
}
```

#### 2.4 Обновлён API клиент
- Файл: `internal/api/client.go`
- Изменения:
  - Добавлено поле `organizationID` в структуру `Client`
  - Добавлен метод `SendGreeting(chatID, text string)` для отправки приветствия
  - Используется POST запрос к `/api/organizations/{orgId}/chat_conversations/{chatId}/completion`

```go
func (c *Client) SendGreeting(chatID, text string) error {
    if c.organizationID == "" {
        return fmt.Errorf("organization ID not set")
    }

    url := fmt.Sprintf("https://claude.ai/api/organizations/%s/chat_conversations/%s/completion",
        c.organizationID, chatID)

    payload := map[string]interface{}{
        "prompt": text,
        "timezone": "UTC",
        // ... другие поля
    }
    // ... curl POST запрос
}
```

#### 2.5 Добавлено уведомление о приветствии
- Файл: `internal/notifier/notifier.go`
- Добавлен метод `NotifyGreeting()` с иконкой солнца ☀️

```go
func (n *Notifier) NotifyGreeting() {
    title := "Утренний привет Клоду ☀️"
    notification := toast.Notification{
        AppID: "ClaudeCompanion",
        Title: title,
        Icon:  getIconPath(),
    }
    notification.Push()
}
```

#### 2.6 Добавлена структура конфигурации
- Файл: `internal/config/config.go`
- Добавлена структура `Greeting`:

```go
type Greeting struct {
    Cron   string `yaml:"greeting_cron"`
    Text   string `yaml:"greeting_text"`
    ChatID string `yaml:"greeting_chat_id"`
}
```

- Добавлены дефолтные значения в `createDefaultConfig()`:

```go
Greeting: Greeting{
    Cron:   "0 8 * * *", // 8 AM every day
    Text:   "Ok",
    ChatID: "", // User must specify chat UUID
},
```

#### 2.7 Реализован Cron планировщик
- Файл: `cmd/claudecompanion/main.go`
- Изменения:
  - Добавлено поле `cronScheduler *cron.Cron` в структуру `App`
  - Добавлена функция `setupGreetingScheduler()` - создаёт и запускает cron
  - Добавлена функция `sendGreeting()` - отправляет приветствие и показывает уведомление
  - Планировщик запускается при получении context от расширения
  - При shutdown останавливается cron scheduler

```go
func (a *App) setupGreetingScheduler() {
    cfg := a.configMgr.Get()

    // Stop existing scheduler if any
    if a.cronScheduler != nil {
        a.cronScheduler.Stop()
    }

    // Check if greeting is configured
    if cfg.Greeting.Cron == "" || cfg.Greeting.ChatID == "" {
        return
    }

    // Create new cron scheduler
    a.cronScheduler = cron.New()

    // Add greeting job
    a.cronScheduler.AddFunc(cfg.Greeting.Cron, func() {
        a.sendGreeting()
    })

    a.cronScheduler.Start()
}
```

#### 2.8 Обновлена конфигурация
- Файл: `config.yaml.example`
- Добавлена секция:

```yaml
greeting:
  greeting_cron: "0 8 * * *"  # Cron schedule: 8 AM every day
  greeting_text: "Ok"         # Message to send
  greeting_chat_id: ""        # UUID of chat conversation (required)
```

---

### 3. Документация фичи приветствия

**Запрос пользователя:**
> добавь в ридми

**Решение:**
- Файлы: `README.md`, `README.ru.md`
- Добавлена секция "Morning Greeting to Claude" / "Утренний привет Клоду"
- Описано:
  - Как получить UUID чата (из URL: `https://claude.ai/chat/{UUID}`)
  - Примеры cron расписаний
  - Что показывается уведомление при отправке

---

### 4. Демонстрация приветствия в Demo режиме

**Запрос пользователя:**
> добавь демонстрацию новой фичи в demo режим

**Решение:**
- Файл: `cmd/claudecompanion/main.go`
- Изменения в `handleDemoMode()`:
  - Добавлено поле `demoGreetingShown bool` для отслеживания состояния
  - В начале каждого цикла (когда value >= 95 && <= 100) показывается уведомление о приветствии
  - Флаг сбрасывается при старте нового цикла

```go
// Show greeting notification at the start of each cycle (once per cycle)
if value >= 95 && value <= 100 && !a.demoGreetingShown {
    a.notifier.NotifyGreeting()
    a.demoGreetingShown = true
}
```

- Обновлена документация в README о том что демонстрируется в demo режиме

---

### 5. Инструкции по сборке

**Запрос пользователя:**
> добавь в ридми описание как собрать приложение

**Решение:**

#### 5.1 Создан build скрипт
- Файл: `build/build-all.bat`
- Функции:
  - Сборка release версии без консоли
  - Сборка debug версии с консолью
  - Генерация Windows ресурсов (иконка)
  - Копирование необходимых файлов
  - Упаковка браузерного расширения

#### 5.2 Создан debug build скрипт
- Файл: `build/build-debug.bat`
- Собирает версию с консолью для отладки
- Опция запуска сразу после сборки

#### 5.3 Документация
- Добавлена детальная секция "Building from Source" в README.md и README.ru.md
- Описаны:
  - Требования (Go 1.21+, Windows)
  - Опция 1: Автоматическая сборка (build-all.bat)
  - Опция 2: Быстрая ручная сборка
  - Опция 3: Сборка с встроенной иконкой (rsrc)
  - Debug сборка
  - Сборка только расширения
  - Проверка результата

---

### 6. Передача User-Agent из браузера

**Запрос пользователя:**
> доработай браузерное расширение пусть отправляет User-Agent в десктоп часть, десткоп подставляет переданный User-Agent во все curl запросы к claude api чтобы эмулировать что запросы идут из браузера

**Реализация:**

#### 6.1 Браузерное расширение
- Файл: `extension/background.js`
- Изменения:
  - Добавлено получение `navigator.userAgent`
  - User-Agent передаётся в payload вместе с cookies

```javascript
// Get User-Agent from browser
const userAgent = navigator.userAgent;
console.log('[ClaudeCompanion] ✅ User-Agent:', userAgent);

const payload = {
  cookies: `sessionKey=${sessionKey}`,
  targetUrl: orgData.usageUrl,
  organizationId: orgData.organizationId,
  userAgent: userAgent
};
```

#### 6.2 HTTP сервер
- Файл: `internal/server/server.go`
- Добавлено поле `UserAgent` в структуру `ContextData`
- Callback теперь принимает `userAgent` параметр

```go
type ContextData struct {
    Cookies        string `json:"cookies"`
    TargetURL      string `json:"targetUrl"`
    OrganizationID string `json:"organizationId"`
    UserAgent      string `json:"userAgent"`
}
```

#### 6.3 API клиент
- Файл: `internal/api/client.go`
- Изменения:
  - Добавлено поле `userAgent` в структуру `Client`
  - Метод `SetContext()` теперь принимает `userAgent`
  - Во всех curl запросах используется User-Agent из браузера
  - Fallback на дефолтный User-Agent если не передан

```go
// Use browser User-Agent if available, otherwise use default
userAgent := c.userAgent
if userAgent == "" {
    userAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36"
}

args = append(args, "-H", fmt.Sprintf("User-Agent: %s", userAgent))
```

#### 6.4 Главное приложение
- Файл: `cmd/claudecompanion/main.go`
- Обновлен callback сервера для передачи userAgent в API клиент

```go
app.httpServer.SetContextCallback(func(cookies, targetURL, organizationID, userAgent string) {
    app.apiClient.SetContext(cookies, targetURL, organizationID, userAgent)
    // ...
})
```

---

### 7. Исправление логики уведомлений

**Запрос пользователя:**
> уведомление о низкой квоте должно выводится один раз при достижении порога < 20, и уведомление о 0 тоже один раз

**Проблема:**
- Уведомления показывались несколько раз
- Сброс состояния уведомлений происходил слишком рано

**Решение:**
- Файл: `cmd/claudecompanion/main.go`
- Изменения в `handleDemoMode()`:
  - Удалён преждевременный вызов `ResetLowValueNotification()` в начале цикла
  - Сброс теперь происходит только когда значение поднимается выше порога
  - Это обеспечивается через `checkLowValueNotifications()` которая вызывает reset когда `value > threshold`

```go
// Removed premature reset:
// a.notifier.ResetLowValueNotification() // DON'T reset here

// Reset only happens in checkLowValueNotifications when value goes above threshold
if value <= cfg.LowValueNotifications.Threshold {
    a.notifier.NotifyLowValue(value, phrase)
} else {
    a.notifier.ResetLowValueNotification() // Only here
}
```

**Подтверждение пользователя:**
> в обчном режиме тоже корректно?

**Ответ:** Да, логика одинаковая для обычного и demo режима - оба используют `checkLowValueNotifications()`

---

### 8. Обсуждение рисков блокировки аккаунта

**Запрос пользователя:**
> может ли cluade.ai заблокировать аккуант если использовать это приложение?

**Ответ:**
Да, есть потенциальный риск, потому что:
1. Используется неофициальный API
2. Автоматические запросы каждые 30-60 секунд
3. Автоматическая отправка сообщений (greeting)
4. Может рассматриваться как bot активность

**Рекомендации для снижения риска:**
1. Увеличить интервал опроса до 60-120 секунд
2. Использовать greeting максимум 1-2 раза в день
3. Личное использование только на одном аккаунте
4. Не запускать 24/7, только в рабочие часы

**Запрос пользователя:**
> да

(подтверждение на добавление предупреждений в README)

**Действия:**

#### 8.1 Обновлены дефолтные значения
- Файл: `internal/config/config.go`
- Изменено: `PollIntervalSeconds` с 30 на 60 секунд для безопасности
- Файл: `config.yaml.example`
- Изменено: `poll_interval_seconds: 60` с комментарием о рекомендации

#### 8.2 Добавлены предупреждения в README
- Файлы: `README.md`, `README.ru.md`
- Добавлена большая секция "⚠️ Disclaimer and Risks" / "⚠️ Дисклеймер и риски"
- Содержание:
  - Информация об использовании неофициального API
  - Потенциальный риск блокировки аккаунта
  - Нарушение Terms of Service
  - Рекомендации для безопасного использования
  - Оценка рисков (LOW vs HIGH)
  - Дисклеймер об ответственности

---

### 9. Функция "Рабочие часы"

**Запрос пользователя:**
> давай внедрим настройку рабочего дня, например получать автоматически делать запрос только с 8 утра до 20 вечера, типа pool_start pool_end

**Реализация:**

#### 9.1 Структура конфигурации
- Файл: `internal/config/config.go`
- Добавлена структура `WorkHours`:

```go
type WorkHours struct {
    Enabled bool   `yaml:"enabled"`
    Start   string `yaml:"start"` // Format: "08:00"
    End     string `yaml:"end"`   // Format: "20:00"
}
```

- Добавлен метод проверки времени:

```go
func (wh *WorkHours) IsWithinWorkHours() bool {
    if !wh.Enabled {
        return true // Always allow if work hours not enabled
    }

    now := time.Now()
    currentTime := now.Format("15:04")

    // Simple string comparison works for HH:MM format
    if start <= end {
        // Normal case: 08:00 - 20:00
        return currentTime >= start && currentTime < end
    } else {
        // Overnight case: 20:00 - 08:00 (next day)
        return currentTime >= start || currentTime < end
    }
}
```

- Добавлены дефолтные значения:

```go
WorkHours: WorkHours{
    Enabled: false,      // Disabled by default
    Start:   "08:00",    // 8 AM
    End:     "20:00",    // 8 PM
},
```

#### 9.2 Проверка в poll loop
- Файл: `cmd/claudecompanion/main.go`
- Добавлена проверка в функцию `poll()`:

```go
// Check work hours
if !cfg.WorkHours.IsWithinWorkHours() {
    logger.Debug("Outside work hours, skipping poll")
    return
}
```

#### 9.3 Обновлена конфигурация
- Файл: `config.yaml.example`
- Добавлена секция:

```yaml
work_hours:
  enabled: false              # Enable to limit polling to work hours only
  start: "08:00"              # Start time (HH:MM format)
  end: "20:00"                # End time (HH:MM format)
```

#### 9.4 Документация
- Файлы: `README.md`, `README.ru.md`
- Добавлена секция "Work Hours" / "Рабочие часы"
- Описано:
  - Как работает функция
  - Поддержка 24-часового формата
  - Поддержка ночных интервалов (overnight)
  - Примеры использования:
    - Типичный рабочий день (08:00-20:00)
    - Офисные часы (09:00-17:00)
    - Ночная смена (20:00-08:00)
  - Польза для снижения риска обнаружения

---

## Итоговые изменения файлов

### Созданные файлы:
1. `build/build-all.bat` - Скрипт полной сборки
2. `build/build-debug.bat` - Скрипт debug сборки

### Изменённые файлы:

1. **extension/background.js**
   - Отправка organizationId
   - Отправка User-Agent

2. **internal/server/server.go**
   - Приём organizationId и userAgent
   - Обновлён callback signature

3. **internal/api/client.go**
   - Хранение organizationId и userAgent
   - Метод SendGreeting()
   - Использование User-Agent в запросах

4. **internal/config/config.go**
   - Структура Greeting
   - Структура WorkHours
   - Метод IsWithinWorkHours()
   - Дефолтные значения обновлены

5. **internal/notifier/notifier.go**
   - Метод NotifyGreeting()

6. **internal/tray/tray.go**
   - Исправлен URL для "Открыть Claude.ai"

7. **cmd/claudecompanion/main.go**
   - Поле cronScheduler
   - Функция setupGreetingScheduler()
   - Функция sendGreeting()
   - Проверка work hours в poll()
   - Исправлена логика уведомлений в demo режиме
   - Обновлён callback для userAgent

8. **config.yaml.example**
   - Изменён poll_interval_seconds на 60
   - Добавлена секция greeting
   - Добавлена секция work_hours

9. **README.md**
   - Секция "⚠️ Disclaimer and Risks"
   - Секция "Building from Source"
   - Секция "Morning Greeting to Claude"
   - Секция "Work Hours"
   - Обновлена документация Demo Mode

10. **README.ru.md**
    - Секция "⚠️ Дисклеймер и риски"
    - Секция "Сборка из исходников"
    - Секция "Утренний привет Клоду"
    - Секция "Рабочие часы"
    - Обновлена документация Demo Mode

11. **go.mod**
    - Добавлена зависимость github.com/robfig/cron/v3

---

## Технические детали

### Архитектура приложения:
```
Browser Extension (Firefox)
    ↓ (sessionKey, organizationId, userAgent)
HTTP Server (:8383)
    ↓
API Client
    ↓ (with browser User-Agent)
Claude.ai API

Cron Scheduler → SendGreeting() → Claude.ai API
```

### Основные компоненты:
1. **Browser Extension** - Извлекает cookies и метаданные
2. **HTTP Server** - Принимает данные от расширения
3. **API Client** - Выполняет запросы к Claude.ai с curl
4. **Cron Scheduler** - Планирует автоматические задачи
5. **Tray Manager** - Управляет системным треем
6. **Icon Generator** - Создаёт динамические иконки
7. **Notifier** - Показывает Toast уведомления
8. **Config Manager** - Hot-reload конфигурации (каждые 2 сек)
9. **Logger** - Логирование в файл/консоль

### Безопасность и риски:
- **Низкий риск:** polling ≥ 60s, greeting 1-2 раза в день, личное использование
- **Высокий риск:** частые запросы, множественные аккаунты, 24/7 работа
- **Рекомендация:** Использовать work_hours для ограничения времени работы

---

## Зависимости

```go
require (
    github.com/getlantern/systray v1.2.2
    github.com/go-toast/toast v0.0.0-20190211030409-01e6764cf0a4
    github.com/robfig/cron/v3 v3.0.1
    gopkg.in/yaml.v3 v3.0.1
)
```

---

## Конфигурация по умолчанию

```yaml
server_port: 8383
poll_interval_seconds: 60
use_curl_fallback: true
gray_mode_threshold: 5
notification_threshold: 10
proxy: ""
enable_file_logging: true
browser_path: ""

low_value_notifications:
  enabled: true
  threshold: 20
  phrases:
    - "Пора идти домой! 🏡"
    - "Система устала. Вы — тоже. 😴"
    # ... и другие
  zero_phrases:
    - "Всё, капут! 💥"
    - "0 — это не число, это приговор. 🛌"
    # ... и другие

demo_mode:
  enabled: false
  duration_seconds: 60

greeting:
  greeting_cron: "0 8 * * *"
  greeting_text: "Ok"
  greeting_chat_id: ""

work_hours:
  enabled: false
  start: "08:00"
  end: "20:00"
```

---

## Дата разработки
Декабрь 2024

## Разработчик
Совместная разработка с Claude AI (Anthropic)

---

**Конец истории диалога**
