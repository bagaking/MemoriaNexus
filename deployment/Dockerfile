FROM golang:alpine
WORKDIR /app
COPY . .
RUN go build -o myapp ./cmd
CMD ["./myapp"]
