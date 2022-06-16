# WB Zero

Веб-сервис на Go получающий заказы пользователей из подписки на канал в NATS
Streaming server. Полученные данные сохраняет в базу данных и в in memory кэш.
Ответы отдает из кэша, при (ре)старте кэш восстанавливает из базы.

Запускается командой `make up` и доступен по адресу <http://127.0.0.1:8080/>.

Команда `make produce` публикует несколько заказов в канал NATS Streaming server.

Остальные команды: `make help`

## Требуемое ПО

- make
- docker-compose
- [dotenv](https://github.com/theskumar/python-dotenv#command-line-interface)
(лучше устаналивать с помощью [pipx](https://pypa.github.io/pipx/):
`pipx install python-dotenv[cli]`)
