FROM golang:1.21.5-alpine3.19

WORKDIR /app

COPY go.mod ./
COPY go.sum ./

RUN go mod download

COPY . .

RUN go build -o main .

RUN rm .env && touch .env

EXPOSE 3000

CMD [ "./main" ]
