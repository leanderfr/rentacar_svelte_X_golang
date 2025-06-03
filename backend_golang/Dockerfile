

FROM golang:1.24.2 AS build

#ADD https://github.com/golang/go/raw/master/lib/time/zoneinfo.zip /zoneinfo.zip
#ENV ZONEINFO /zoneinfo.zip


WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . /app
COPY .env /app

RUN CGO_ENABLED=0 GOOS=linux go build -o api main.go

FROM scratch

WORKDIR /app

COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=build /app/api ./
COPY --from=build /app/.env ./

EXPOSE 8070

CMD [ "./api" ]

#   docker build --rm -t rentacar .

# docker run  -d -p 8070:8070 --name rentacar rentacar




