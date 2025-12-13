# Создание иконок для расширения

Для Firefox расширения требуются иконки двух размеров.

## Требования

- `icon48.png` - 48x48 пикселей
- `icon96.png` - 96x96 пикселей
- Формат: PNG с прозрачным фоном (рекомендуется)

## Способы создания

### 1. Онлайн-генераторы

Используйте любой онлайн генератор иконок:
- https://www.favicon-generator.org/
- https://www.websiteplanet.com/webtools/favicon-generator/
- https://favicon.io/

### 2. Графические редакторы

**GIMP (бесплатно):**
1. Создайте новое изображение 96x96
2. Добавьте текст "CT" (ClaudeCompanion) или логотип
3. Экспортируйте как PNG: `icon96.png`
4. Измените размер до 48x48
5. Экспортируйте как PNG: `icon48.png`

**Photoshop/Figma:**
1. Создайте artboard 96x96
2. Добавьте дизайн
3. Экспортируйте в обоих размерах

### 3. Командная строка (ImageMagick)

```bash
# Создать простую иконку с текстом
convert -size 96x96 xc:transparent -font Arial -pointsize 48 \
  -fill "#4A90E2" -gravity center -annotate +0+0 "CT" icon96.png

convert icon96.png -resize 48x48 icon48.png
```

### 4. Python (Pillow)

```python
from PIL import Image, ImageDraw, ImageFont

# Создать 96x96
img = Image.new('RGBA', (96, 96), (0, 0, 0, 0))
draw = ImageDraw.Draw(img)

# Добавить круг
draw.ellipse([10, 10, 86, 86], fill='#4A90E2')

# Добавить текст
font = ImageFont.truetype("arial.ttf", 40)
draw.text((48, 48), "CT", fill='white', anchor="mm", font=font)

img.save('icon96.png')

# Создать 48x48
img.resize((48, 48)).save('icon48.png')
```

## Рекомендации по дизайну

1. **Простота**: Иконка должна быть читаемой в маленьком размере
2. **Контраст**: Используйте контрастные цвета
3. **Узнаваемость**: Логотип Claude или буквы "CT"
4. **Прозрачность**: PNG с прозрачным фоном выглядит лучше

## Пример дизайна

Простая иконка может содержать:
- Круг синего цвета (#4A90E2)
- Белые буквы "CT" в центре
- Или: стилизованный логотип Claude (если есть права)

## После создания

1. Сохраните обе иконки в папку `extension/`
2. Проверьте, что файлы называются точно:
   - `icon48.png`
   - `icon96.png`
3. Перезагрузите расширение в Firefox
4. Иконка должна появиться в панели расширений
