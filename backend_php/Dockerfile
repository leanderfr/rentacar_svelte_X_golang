

FROM php:8.1-apache
RUN docker-php-ext-install mysqli

WORKDIR /var/www/html

COPY . /var/www/html/
RUN mkdir -p /var/www/html/tmp
RUN chmod -R 0777 /var/www/html/tmp


RUN echo "ServerName localhost" >> /etc/apache2/apache2.conf

RUN a2enmod rewrite

EXPOSE 80

# docker build -t rentacar_php .

# docker run -p 8072:80 rentacar_php