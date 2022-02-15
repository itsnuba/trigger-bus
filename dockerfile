FROM nabihubadah/trigger_bus-builder AS builder
FROM alpine

WORKDIR /app

COPY --from=builder /app/trigger_bus* ./
COPY --from=builder /app/*.env ./

EXPOSE 8080

ENTRYPOINT ["./trigger_bus"]
