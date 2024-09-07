FROM golang:1.22-alpine3.20

WORKDIR /app

RUN go install github.com/air-verse/air@latest

COPY go.* ./
RUN go mod download

COPY . .

EXPOSE 3333

CMD [ "air", "-c", ".air.toml" ]
