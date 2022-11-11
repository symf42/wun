FROM golang:1.19.3-bullseye

WORKDIR /usr/src

COPY . .

ARG MYSQL_HOSTNAME=localhost
ARG MYSQL_PORT=3306
ARG MYSQL_USERNAME=root
ARG MYSQL_PASSWORD=secret
ARG MYSQL_DATABASE=wun

RUN echo MYSQL_HOSTNAME=${MYSQL_HOSTNAME} > /usr/src/.env
RUN echo MYSQL_PORT=${MYSQL_PORT} >> /usr/src/.env
RUN echo MYSQL_USERNAME=${MYSQL_USERNAME} >> /usr/src/.env
RUN echo MYSQL_PASSWORD=${MYSQL_PASSWORD} >> /usr/src/.env
RUN echo MYSQL_DATABASE=${MYSQL_DATABASE} >> /usr/src/.env

RUN go build -v

CMD ./wun