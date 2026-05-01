FROM golang:1.26

# Install goose
RUN go install github.com/pressly/goose/v3/cmd/goose@latest

CMD ["goose", "up"]
