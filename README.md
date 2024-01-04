# simple bank

## 简介
simple bank 项目是用来学习 Go 后端开发流程所产生的

### 涉及到的工具
1. db migration 
 数据库迁移工具，命令行形式安装
https://github.com/golang-migrate/migrate

2. 






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