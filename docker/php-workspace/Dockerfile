ARG DOCKER_PHP_VERSION

FROM php:${DOCKER_PHP_VERSION}-fpm

ARG TZ='UTC'

USER root

RUN rm /etc/apt/preferences.d/no-debian-php

RUN apt-get update && apt-get install -y \
        supervisor \
        imagemagick \
        libxslt-dev \
        libzip-dev \
        libpng-dev \
        libjpeg-dev \
        libfreetype6  \
        libfreetype6-dev \
        libjpeg62-turbo-dev \
        libonig-dev \
        libicu-dev \
        icu-devtools \
        patch \
        bash \
        htop \
        curl \
        vim \
        openssh-client \
        python3 \
        python3-pip \
        mariadb-client \
        npm

RUN apt-get install -y \
            wget \
            gzip \
            tar

#  Install php extensions
RUN php -m && \
    docker-php-ext-configure bcmath --enable-bcmath && \
    docker-php-ext-configure gd \
      --with-freetype=/usr/include/ \
      --with-jpeg=/usr/include/ && \
    docker-php-ext-configure intl && \
    docker-php-ext-configure pcntl --enable-pcntl && \
    docker-php-ext-configure soap && \
    docker-php-ext-install exif \
        mysqli \
        ftp \
        opcache \
        xsl \
        bcmath \
        gd \
        gettext \
        intl \
        pcntl \
        soap \
        sysvsem \
        zip \
        pdo_mysql && \
    pecl install redis && \
    docker-php-ext-enable redis

RUN apt-get update && apt-get install -y cron

RUN pecl install xdebug \
    && docker-php-ext-enable xdebug

# Clean
RUN rm -rf /var/cache/apk/* && docker-php-source delete

USER root

RUN curl -fsSL https://deb.nodesource.com/setup_23.x | bash - \
    && apt-get install -y nodejs

RUN npm install -g laravel-echo-server

RUN wget -q -O /usr/local/bin/composer https://getcomposer.org/download/latest-stable/composer.phar \
    && chmod +x /usr/local/bin/composer

RUN python3 -m pip install --upgrade supervisor --ignore-installed --break-system-packages

COPY supervisord.conf /etc/supervisord.conf
COPY crontab /etc/cron.d/crontab

RUN chmod 0644 /etc/cron.d/crontab
RUN crontab /etc/cron.d/crontab
RUN touch /var/log/cron.log

ENTRYPOINT ["/usr/bin/supervisord", "-n", "-c",  "/etc/supervisord.conf"]

WORKDIR /var/www/

USER www-data:www-data