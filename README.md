[![Build Status](https://travis-ci.org/angular/angular.svg?branch=master)](https://travis-ci.org/angular/angular)

# Awesome

Awesome is a Wechat Message Queue platform for notify the user

## Quickstart


## Changelog
[Learn about the latest improvements](https://github.com/vimeracke/awesomeProject/blob/master/CHANGELOG.md)


- 执行命令的接口：`/deviceControl`
  - 方法：POST
  - 参数：
    - device_id：字符串，设备号，必须
    - cmd：字符串，命令名，必须
    - arg：字符串，命令参数，非必须
    - timeout：整型，命令超时时间(默认为10s)，非必须
  - 返回：
    - error_code：整型，错误码
    - error_msg：字符串，错误信息
    - cmd_id：字符串，命令标识码

- 获取结果信息的接口：`/deviceInfo`
  - 方法：POST
  - 参数：
    - cmd_id : 字符串，命令标识码
  - 返回：
    - error_code：整型，错误码
    - error_msg：字符串，错误信息
    - result：字典，命令执行结果

- 回调：
  - MQTT订阅：`tf/Attendance/v1/devices/+/control` 话题，当收到消息时，根据消息ID，将命令存入库中。
  - MQTT订阅：`tf/Attendance/v1/devices/+/info` 话题，当收到消息时，根据消息ID，将命令执行结果存入库中。


- 数据结构：
```
{
    "id": 整型,
    "topic": 字符串,
    "cmd": 字符串,
    "arg": 字符串,
    "timeout": 整型,
    "result": {
        "success": 布尔型,
        "message": 字符串,
        "device_id": 字符串,
        "result": {
            ....
        }
    }
}
```
