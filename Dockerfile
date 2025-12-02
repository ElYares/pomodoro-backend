# ---- Stage 1: Build ----
FROM golang:1.24 AS builder

# Directorio de trabajo dentro del contenedor
WORKDIR /app

# Copiamos los archivos de dependencias
COPY go.mod go.sum ./

# Descargamos los módulos sin compilar aún
RUN go mod download

# Copiamos el código fuente
COPY . .

# Compilamos binario optimizado
RUN CGO_ENABLED=0 GOOS=linux go build -o pomodoro .

# ---- Stage 2: Run ----
FROM alpine:3.19

WORKDIR /app

# Copiar el binario compilado
COPY --from=builder /app/pomodoro .

# Variables por defecto (se pueden sobreescribir en docker-compose)
ENV HTTP_PORT=8080

EXPOSE 8080

CMD ["./pomodoro"]
