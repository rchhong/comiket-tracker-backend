FROM golang

WORKDIR /app

RUN go install github.com/air-verse/air@latest

COPY go.mod go.sum ./
RUN go mod download && go mod verify

# COPY . .
# RUN go build -v -o /usr/local/bin/comiket-backend .

CMD ["air", "-c", ".air.toml"]

