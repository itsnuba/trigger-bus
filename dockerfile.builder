FROM golang:1.16-alpine

WORKDIR /app

ADD go.* ./

RUN go mod download

COPY . .

# build app
RUN go build -o trigger_bus -tags=nomsgpack .
