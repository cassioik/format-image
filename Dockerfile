FROM golang:1.21-alpine
WORKDIR /app
COPY . .
RUN go build format-image.go
CMD ["/app/format-image"]