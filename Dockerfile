FROM golang:1.19
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
ENV GOARCH=amd64
RUN go build -o main ./src/nthreads
CMD ["./main"]
