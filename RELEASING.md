# Releasing ClaudeCompanion

Это руководство описывает процесс создания нового релиза.

## Автоматический релиз через GitHub Actions

### Шаг 1: Обновите версию расширения

Если изменялось расширение, обновите версию в `extension/manifest.json`:

```json
{
  "version": "1.0.1"
}
```

### Шаг 2: Закоммитьте изменения

```bash
git add .
git commit -m "Release v1.0.1"
git push origin main
```

### Шаг 3: Создайте тег

```bash
# Формат: v{MAJOR}.{MINOR}.{PATCH}
git tag v1.0.1
git push origin v1.0.1
```

### Шаг 4: GitHub Actions автоматически

После пуша тега GitHub Actions автоматически:

1. ✅ Соберёт приложение для всех платформ:
   - Windows x64
   - Windows x86
   - Linux x64
   - Linux ARM64
   - macOS Intel
   - macOS Apple Silicon

2. ✅ Встроит иконку в Windows-версии (через rsrc)

3. ✅ Упакует браузерное расширение Firefox (.xpi)

4. ✅ Создаст архивы:
   - `.zip` для Windows
   - `.tar.gz` для Linux/macOS
   - `.xpi` для Firefox

5. ✅ Сгенерирует `checksums.txt` с SHA256

6. ✅ Создаст GitHub Release с описанием

7. ✅ Загрузит все артефакты в релиз

### Шаг 5: Проверьте релиз

Перейдите на https://github.com/irumil/ClaudeCompanion/releases и проверьте новый релиз.

## Что включено в каждый архив

### Windows (zip)
```
claudecompanion-{version}-windows-amd64.zip
├── claudecompanion.exe     # Приложение (со встроенной иконкой)
├── icon.ico                # Иконка для уведомлений
├── config.yaml             # Пример конфигурации
├── README.md               # English
├── README.ru.md            # Русский
└── README.tt.md            # Татарча
```

### Linux/macOS (tar.gz)
```
claudecompanion-{version}-linux-amd64.tar.gz
├── claudecompanion         # Исполняемый файл
├── config.yaml             # Пример конфигурации
├── README.md               # English
├── README.ru.md            # Русский
└── README.tt.md            # Татарча
```

### Firefox Extension
```
claudecompanion-extension-{version}.xpi
├── manifest.json
├── background.js
├── options.html
├── options.js
├── icon48.png
├── icon96.png
└── PRIVACY.md
```

## Ручная сборка (без GitHub Actions)

Если нужно собрать локально:

### Все платформы сразу

```bash
cd build
./build-all.bat   # Windows
./build-all.sh    # Linux/macOS
```

### Отдельно для конкретной платформы

```bash
# Windows amd64
GOOS=windows GOARCH=amd64 go build -ldflags "-H windowsgui -s -w" -o dist/claudecompanion.exe ./cmd/claudecompanion

# Linux amd64
GOOS=linux GOARCH=amd64 go build -ldflags "-s -w" -o dist/claudecompanion ./cmd/claudecompanion

# macOS Apple Silicon
GOOS=darwin GOARCH=arm64 go build -ldflags "-s -w" -o dist/claudecompanion ./cmd/claudecompanion
```

## Версионирование (Semantic Versioning)

Используем формат: `v{MAJOR}.{MINOR}.{PATCH}`

- **MAJOR** - несовместимые изменения API
- **MINOR** - новый функционал (обратная совместимость)
- **PATCH** - исправления багов

Примеры:
- `v1.0.0` - первый стабильный релиз
- `v1.1.0` - добавлена новая функция
- `v1.1.1` - исправлен баг
- `v2.0.0` - breaking changes

## Проверка перед релизом

Перед созданием тега проверьте:

- [ ] Все тесты проходят: `go test ./...`
- [ ] Приложение собирается: `go build ./cmd/claudecompanion`
- [ ] Расширение упаковывается: `cd build && package-extension.bat`
- [ ] Обновлена документация (README.md, CHANGELOG.md)
- [ ] Обновлена версия в extension/manifest.json (если изменялось)
- [ ] Нет незакоммиченных изменений: `git status`

## Откат релиза

Если релиз создан по ошибке:

```bash
# Удалить тег локально
git tag -d v1.0.1

# Удалить тег на GitHub
git push origin :refs/tags/v1.0.1
```

Затем удалите релиз вручную через GitHub UI.

## Changelog

Рекомендуется вести CHANGELOG.md с описанием изменений в каждой версии.

Пример:

```markdown
# Changelog

## [1.1.0] - 2025-12-20

### Added
- Встроенная иконка в Windows exe
- Поддержка Linux ARM64
- Татарский README

### Fixed
- Уведомления теперь показывают иконку

### Changed
- Обновлены build-скрипты

## [1.0.0] - 2025-12-01

### Added
- Первый стабильный релиз
```

## Публикация расширения на AMO

После создания релиза на GitHub, расширение нужно отдельно загрузить на Mozilla Add-ons:

1. Скачайте `claudecompanion-extension-{version}.xpi` из релиза
2. Перейдите на https://addons.mozilla.org/developers/
3. Загрузите новую версию расширения
4. Дождитесь проверки Mozilla

## Проблемы и решения

### GitHub Actions не запускается

**Причина:** Тег не в формате `v*.*.*`

**Решение:** Используйте формат `v1.0.0`, а не `1.0.0` или `release-1.0.0`

### Сборка для Windows без иконки

**Причина:** rsrc не установлен или не может найти icon.ico

**Решение:** Убедитесь что `icon.ico` в корне репозитория

### Checksums не генерируются

**Причина:** Архивы не созданы

**Решение:** Проверьте логи сборки, убедитесь что все архивы созданы перед генерацией checksums

---

**Вопросы?** Создайте issue: https://github.com/irumil/ClaudeCompanion/issues
