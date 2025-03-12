#!/bin/sh
set -e

echo "Configurando diretórios..."
mkdir -p /run/mysqld
chown -R mysql:mysql /run/mysqld /var/lib/mysql

# Configuração do MariaDB
cat > /etc/my.cnf << EOF
[mysqld]
user = mysql
bind-address = 0.0.0.0
port = 3306
datadir = /var/lib/mysql
socket = /run/mysqld/mysqld.sock
skip-networking = 0
skip-name-resolve

[client]
socket = /run/mysqld/mysqld.sock

[mysql]
socket = /run/mysqld/mysqld.sock
EOF

# Inicializa o banco de dados se necessário
if [ ! -d "/var/lib/mysql/mysql" ]; then
    echo "Inicializando MariaDB..."
    mysql_install_db --user=mysql --datadir=/var/lib/mysql

    echo "Iniciando MariaDB temporariamente..."
    mysqld --user=mysql --datadir=/var/lib/mysql --skip-networking &
    MYSQL_PID=$!
    
    # Espera o processo iniciar
    while ! mysqladmin ping --silent; do
        sleep 1
    done

    echo "Configurando usuários e banco de dados..."
    mysql << EOF
DELETE FROM mysql.user WHERE User='';
DELETE FROM mysql.user WHERE User='root' AND Host NOT IN ('localhost', '%');
CREATE DATABASE IF NOT EXISTS people_db;
CREATE USER 'root'@'%' IDENTIFIED BY 'root';
GRANT ALL PRIVILEGES ON . TO 'root'@'%' WITH GRANT OPTION;
ALTER USER 'root'@'localhost' IDENTIFIED BY 'root';
FLUSH PRIVILEGES;
EOF

    echo "Parando MariaDB temporário..."
    kill $MYSQL_PID
    wait $MYSQL_PID || true
fi

# Inicia o MariaDB em background
echo "Iniciando MariaDB..."
mysqld --user=mysql --datadir=/var/lib/mysql &
MYSQL_PID=$!

# Espera o MariaDB ficar pronto
echo "Aguardando MariaDB ficar pronto..."
while ! mysqladmin -u root -proot ping --silent; do
    sleep 1
done

echo "MariaDB está pronto!"

# Inicia a aplicação Go
echo "Iniciando aplicação Go..."
cd /app && ./myapp

# Mantém o container rodando
wait $MYSQL_PID