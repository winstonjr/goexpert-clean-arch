FROM golang:1.23 AS build
WORKDIR /app
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o desafio-clean-arch ./cmd/ordersystem/main.go ./cmd/ordersystem/wire_gen.go

FROM scratch
WORKDIR /app
COPY --from=build /app/desafio-clean-arch .
COPY --from=build /app/cmd/ordersystem/.env .
ENTRYPOINT ["./desafio-clean-arch"]