# simple bank

## 简介
simple bank 项目是用来学习 Go 后端开发流程所产生的

### 涉及到的工具及包
1. db migration 
 数据库迁移工具，命令行形式安装
https://github.com/golang-migrate/migrate

2. sqlc
	根据 SQL 生成数据库读取代码

3. gin 
	实现 RESTful API

4. viper
	读取配置

5. testify
	测试断言

6.  mock 
	mock 测试
	https://github.com/uber-go/mock

	命令 `$mockgen -package mockdb  -destination db/mock/store.go github.com/xiaowuzai/simplebank/db/sqlc Store` 





## 数据库迁移

### 创建迁移文件
```
$ migrate create -ext sql -dir db/migration -seq init_schema
/Users/zly/workspace/go/src/github.com/xiaowuzai/simplebank/db/migration/000001_init_schema.up.sql
/Users/zly/workspace/go/src/github.com/xiaowuzai/simplebank/db/migration/000001_init_schema.down.sql
```
### 执行迁移命令

```
migrate -path db/migration -database "postgresql://root:123456@localhost:5432/simple_bank?sslmode=disable" -verbose up
```


## 注意
1. 编写单元测试时，要保持用例相互独立，减少彼此的耦合