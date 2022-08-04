# lubimyczytacrss

lubimyczytacrss is an RSS feed generator for [lubimyczytac.pl](https://lubimyczytac.pl).

## How to use?

### Option 1: Use Docker

```shell
$ docker run -d \
  -p 9234:9234 \
  --name lubimyczytacrss \
  kichiyaki/lubimyczytacrss:latest
```
Or via docker-compose:
```yaml
version: "3.6"

services:
  lubimyczytacrss:
    image: kichiyaki/lubimyczytacrss:latest
    restart: unless-stopped
    ports:
      - '9234:9234'
    networks:
      - default
```

## Available endpoints

### GET /api/v1/rss/authors/{authorID}

Returns the RSS feed of the newest books written by the author with the given ID.
