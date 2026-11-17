# Etapa 1: compilación
FROM golang:1.22-alpine AS build
WORKDIR /src
COPY go.mod ./
COPY . .
RUN CGO_ENABLED=0 go build -o web-uoc .

# Etapa 2: imagen final
FROM alpine:3.20
WORKDIR /app
COPY --from=build /src/web-uoc .
EXPOSE 8080
CMD ["./web-uoc"]
