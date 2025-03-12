# Estágio de construção do Go
FROM golang:1.21-alpine AS builder
WORKDIR /app

# Instalar dependências do MySQL
RUN apk add --no-cache gcc musl-dev mysql-client

# Copiar e construir a aplicação
COPY . .
RUN go mod download
RUN go build -o /myapp

# Estágio de execução final
FROM alpine:latest
WORKDIR /app

# Instalar dependências do MySQL e ferramentas
RUN apk add --no-cache mariadb mariadb-client

# Copiar binário da aplicação e arquivos necessários
COPY --from=builder /myapp /app/myapp
COPY .env /app/.env
COPY entrypoint.sh /app/entrypoint.sh

# Configurar diretórios e permissões
RUN mkdir -p /run/mysqld /var/lib/mysql && \
    chown -R mysql:mysql /run/mysqld /var/lib/mysql && \
    chmod +x /app/entrypoint.sh

EXPOSE 8080 3306

ENTRYPOINT ["/app/entrypoint.sh"]