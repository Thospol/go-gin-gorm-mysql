# go-gin-gorm-mysql

## How to connect database mysql with docker (local).
```
docker pull mysql

docker run --name=my-local -e MYSQL_ROOT_PASSWORD=password -e MYSQL_DATABASE=local -p 3306:3306 -d mysql --character-set-server=utf8mb4 --collation-server=utf8mb4_unicode_ci 
```