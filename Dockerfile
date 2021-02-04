FROM php:7.2-alpine
WORKDIR /parsers-php
COPY . .
RUN ls
ENTRYPOINT ["php", "JobReceiver.php"]

