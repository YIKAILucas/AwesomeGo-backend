
![icon](https://vimeracke.oss-cn-shenzhen.aliyuncs.com/ima.jpeg)

[![Build Status](https://travis-ci.org/angular/angular.svg?branch=master)](https://travis-ci.org/angular/angular)


## Awesome is an awesome IOT platform

### Quickstart
   **docker build awesomeProject**



### Changelog
[查看往期版本迭代](http://157.122.146.233:88/G2/awesome.back-end/blob/master/CHANGELOG.md)


## 技术相关
1. 存储方案
	- MongoDB
	- TSDB
	- HBase
	- PostgreSQL
2. 消息中间件
	- RabbitMQ
	- EMQ
3. 数据缓存
	- Redis
4. 通信协议
	- MQTT
	
	
### 中间件技术细节
**EMQ开启上下线监听的方法：**
```
- docker exec -it emq /bin/sh
- vi etc/acl.conf
- 添加一行：{allow, {user, "golang-server"}, subscribe, ["$SYS/#"]}.
- 重启容器
```

**设置MySQL以支持中文：**

编辑：`/etc/mysql/mysql.cnf`，修改为以下配置：
```
!includedir /etc/mysql/conf.d/
!includedir /etc/mysql/mysql.conf.d/

# 设置字符集为UTF8,以支持中文
[mysqld]
character-set-server=utf8
collation-server=utf8_general_ci

[mysql]
default-character-set = utf8

[mysql.server]
default-character-set = utf8

[mysqld_safe]
default-character-set = utf8

[client]
default-character-set = utf8

```
