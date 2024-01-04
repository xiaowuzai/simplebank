## 启动 mysql

docker run --name mysql8.2 -p 3306:3306  -e MYSQL_ROOT_PASSWORD=123456 -e MYSQL_DATABASE=simple_bank -d mysql:8.2