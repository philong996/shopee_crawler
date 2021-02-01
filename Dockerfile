FROM ubuntu:latest

WORKDIR /crawler_shopee

COPY . .


RUN apt-get update && apt-get -y install cron
RUN docker-php-ext-install mysqli && docker-php-ext-enable mysqli
RUN apk add --no-cache tini

RUN ls -l
