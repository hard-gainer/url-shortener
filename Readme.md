## API Usage Examples

**Endpoint:** `POST /api/shorten`
Создаёт новую укороченную ссылку.

**Request:**
```json
{
    "url": "https://www.google.com/search?q=%D0%BA%D0%BE%D1%88%D0%BA%D0%B8+%D0%BA%D0%B0%D1%80%D1%82%D0%B8%D0%BD%D0%BA%D0%B8&sca_esv=200a9dc461e1f367&sxsrf=AHTn8zpSCQzR-dORKTsTaOIG2QKPrZzEYQ%3A1741963706360&source=hp&ei=ukHUZ9bDE4PIwPAPiv2U2QE&iflsig=ACkRmUkAAAAAZ9RPymbtILAXDO48rDBFLi5VKouolWpE&ved=0ahUKEwjWiLW_6ImMAxUDJBAIHYo-JRsQ4dUDCBg&uact=5&oq=%D0%BA%D0%BE%D1%88%D0%BA%D0%B8+%D0%BA%D0%B0%D1%80%D1%82%D0%B8%D0%BD%D0%BA%D0%B8&gs_lp=Egdnd3Mtd2l6IhvQutC-0YjQutC4INC60LDRgNGC0LjQvdC60LgyBRAAGIAEMgUQABiABDIFEAAYgAQyBRAAGIAEMgUQABiABDIFEAAYgAQyBRAAGIAEMgUQABiABDIGEAAYFhgeMgYQABgWGB5IkSJQowxY6B9wAXgAkAEAmAGCAaABnwiqAQQxMy4xuAEDyAEA-AEBmAIPoALQCKgCCsICBxAjGCcY6gLCAgQQIxgnwgIKECMYgAQYJxiKBcICCxAAGIAEGLEDGIMBwgILEC4YgAQYsQMYgwHCAggQABiABBixA8ICERAuGIAEGLEDGNEDGIMBGMcBwgIIEC4YgAQYsQPCAg4QLhiABBixAxiDARiKBcICCxAuGIAEGLEDGNQCwgIOEAAYgAQYsQMYgwEYigXCAgsQLhiABBjHARivAcICBRAuGIAEwgIIEC4YgAQY1AKYAwnxBcbRpDpmqJDikgcEMTQuMaAHmK0B&sclient=gws-wiz#vhid=0qUTOpyBXYMBnM&vssid=_wEHUZ-e9N-38wPAPl86coQ0_36"
}
```

**Response:**
```json
{
    "short_url": "/e3Yc2CQVCJ",
    "original_url": "https://www.google.com/search?q=%D0%BA%D0%BE%D1%88%D0%BA%D0%B8+%D0%BA%D0%B0%D1%80%D1%82%D0%B8%D0%BD%D0%BA%D0%B8&sca_esv=200a9dc461e1f367&sxsrf=AHTn8zpSCQzR-dORKTsTaOIG2QKPrZzEYQ%3A1741963706360&source=hp&ei=ukHUZ9bDE4PIwPAPiv2U2QE&iflsig=ACkRmUkAAAAAZ9RPymbtILAXDO48rDBFLi5VKouolWpE&ved=0ahUKEwjWiLW_6ImMAxUDJBAIHYo-JRsQ4dUDCBg&uact=5&oq=%D0%BA%D0%BE%D1%88%D0%BA%D0%B8+%D0%BA%D0%B0%D1%80%D1%82%D0%B8%D0%BD%D0%BA%D0%B8&gs_lp=Egdnd3Mtd2l6IhvQutC-0YjQutC4INC60LDRgNGC0LjQvdC60LgyBRAAGIAEMgUQABiABDIFEAAYgAQyBRAAGIAEMgUQABiABDIFEAAYgAQyBRAAGIAEMgUQABiABDIGEAAYFhgeMgYQABgWGB5IkSJQowxY6B9wAXgAkAEAmAGCAaABnwiqAQQxMy4xuAEDyAEA-AEBmAIPoALQCKgCCsICBxAjGCcY6gLCAgQQIxgnwgIKECMYgAQYJxiKBcICCxAAGIAEGLEDGIMBwgILEC4YgAQYsQMYgwHCAggQABiABBixA8ICERAuGIAEGLEDGNEDGIMBGMcBwgIIEC4YgAQYsQPCAg4QLhiABBixAxiDARiKBcICCxAuGIAEGLEDGNQCwgIOEAAYgAQYsQMYgwEYigXCAgsQLhiABBjHARivAcICBRAuGIAEwgIIEC4YgAQY1AKYAwnxBcbRpDpmqJDikgcEMTQuMaAHmK0B&sclient=gws-wiz#vhid=0qUTOpyBXYMBnM&vssid=_wEHUZ-e9N-38wPAPl86coQ0_36"
}
```

**Endpoint:** `GET /api/info/{shortURL}`

**Response:**
Html code of original_url


## Задание (Стажер-разработчик)

Укорачиватель ссылок

Необходимо реализовать сервис, который должен предоставлять API по созданию сокращенных ссылок следующего формата:
- Ссылка должна быть уникальной и на один оригинальный URL должна ссылаться только одна сокращенная ссылка.
- Ссылка должна быть длинной 10 символов
- Ссылка должна состоять из символов латинского алфавита в нижнем и верхнем регистре, цифр и символа _ (подчеркивание)

Сервис должен быть написан на Go и принимать следующие запросы по http:
1. Метод Post, который будет сохранять оригинальный URL в базе и возвращать сокращённый
2. Метод Get, который будет принимать сокращённый URL и возвращать оригинальный URL

Решение должно быть предоставлено в «конечном виде», а именно:
- Сервис должен быть распространён в виде Docker-образа 
- В качестве хранилища ожидается использовать две реализации. Какое хранилище использовать, указывается параметром при запуске сервиса.  
    - Первое это postgresql.
    - Второе - самостоятельно написать пакет для хранения ссылок в памяти приложения.
- Покрыть реализованный функционал Unit-тестами

Результат предоставить в виде публичного репозитория на github.com

В процессе собеседования-ревью посмотрим:
- Как генерируются ссылки и почему предложенный алгоритм будет работать; насколько он соответствует заданию и прост в понимании.
- Как раскиданы типы по файлам, файлики по пакетам, пакеты по приложению: структуру проекта.
- Как обрабатываются ошибки в разных сценариях использования
- Насколько удобен и логичен сервис в использовании
- Как сервис будет себя вести, если им будут пользоваться одновременно сотни людей (как например youtu.be / ya.cc)
- Что будет, если сервис оставить работать на очень долгое время
- Общую чистоту кода

