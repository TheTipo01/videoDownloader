# videoDownloader

[![Go Report Card](https://goreportcard.com/badge/github.com/TheTipo01/videoDownloader)](https://goreportcard.com/report/github.com/TheTipo01/videoDownloader)

This telegram bot downloads videos from the configured sites and sends them back in the chat.
It also supports inline queries.

## Installation

### Natively
Just grab the latest release from the [releases page](https://github.com/TheTipo01/videoDownloader/releases/), modify
the included `example_config.yml` file (adding your telegram token that you got from [@BotFather](https://t.me/BotFather)), rename it to `config.yml` and run the bot.

Make sure to have ffmpeg and yt-dlp installed and in your PATH.
### Docker
Clone the repo, modify the included `example_config.yml` file (adding your telegram token that you got from [@BotFather](https://t.me/BotFather)), rename it to `config.yml`, and do a `docker-compose up`.


The docker image is also available on [Docker hub](https://hub.docker.com/r/thetipo01/videodownloader), [Quay.io](https://quay.io/repository/thetipo01/videodownloader) and [Github packages](https://github.com/TheTipo01/videoDownloader/pkgs/container/videodownloader)
