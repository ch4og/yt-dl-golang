## yt-dl-golang
This telegram bot allows you to download videos from youtube by just sending it a link.

It's just simple bot that I wrote to kinda learn Go. 

### To use this bot you will need [local telegram API](https://core.telegram.org/bots/api#using-a-local-bot-api-server) 

Without it uploads are limited with 50MB uploads, so bot needs local API.

You will need to create .env file and specify:

`TELEGRAM_API_TOKEN` (get it from [@BotFather](https://t.me/BotFather) in telegram)

`TELEGRAM_API_URL` url on which you run your local API. By default it is http://localhost:8081
