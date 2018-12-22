[![Build Status](https://travis-ci.org/angular/angular.svg?branch=master)](https://travis-ci.org/angular/angular)

# Awesome

Awesome is a Wechat Message Queue platform for notify the user

## Quickstart


## Changelog
[Learn about the latest improvements](https://github.com/vimeracke/awesomeProject/blob/master/CHANGELOG.md)





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