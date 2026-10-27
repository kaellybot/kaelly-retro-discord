# kaelly-retro-discord

[![CI](https://github.com/kaellybot/kaelly-retro-discord/actions/workflows/ci.yml/badge.svg)](https://github.com/kaellybot/kaelly-retro-discord/actions/workflows/ci.yml)
[![codecov](https://codecov.io/gh/kaellybot/kaelly-retro-discord/branch/main/graph/badge.svg)](https://codecov.io/gh/kaellybot/kaelly-retro-discord) 


Application to interact with Discord written in Go, dedicated to **DOFUS Retro**.

It is a stripped-down fork of [kaelly-discord](https://github.com/kaellybot/kaelly-discord): it only serves
`/about`, `/help`, `/align`, `/job` and `/config` (`get`, `rss`, `server` and `twitter`), and tags every AMQP message with
`DOFUS_RETRO`.

## Local Development

You will probably need these docker images to make it work properly:

```Bash
docker run --name mysql --restart=always -p 3306:3306 -e MYSQL_ROOT_HOST=% -e MYSQL_ROOT_PASSWORD=password -d mysql/mysql-server:latest 
docker run --name phpmyadmin --restart=always -d --link mysql:db -p 9001:80 phpmyadmin/phpmyadmin:latest
docker run --name rabbitmq --restart=always -p 15672:15672 -p 5672:5672 -d rabbitmq:3-management:latest
docker run -p 6379:6379 --name redis -d redis
```