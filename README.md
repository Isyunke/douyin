# 测试
## 准备工作
### 安装mysql
https://www.runoob.com/mysql/mysql-install.html
### 安装redis
https://www.runoob.com/redis/redis-install.html
### 导入测试数据
进入MySQL命令行，执行以下语句，注意source中要用sql文件的完整路径，不能用相对路径
```sql
drop database douyin;
create database douyin;
source ${PROJECT_PATH}/data/test.sql
```
