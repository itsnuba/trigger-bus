FROM golang:1.18-alpine

WORKDIR /app

ADD go.* ./

RUN go mod download

COPY . .

# build app
RUN go build -o trigger_bus -tags=nomsgpack .
