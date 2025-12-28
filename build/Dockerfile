FROM golang

WORKDIR /app

RUN go install github.com/air-verse/air@latest

COPY . .
RUN go mod download && go mod verify

# RUN go build -v -o /usr/local/bin/comiket-backend ./...

CMD ["air", "-c", "/app/build/.air.toml"]