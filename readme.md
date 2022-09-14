# 临时文本服务

[![codecov](https://codecov.io/gh/sixwaaaay/temp-text/branch/master/graph/badge.svg?token=UwTUzTcS2G)](https://codecov.io/gh/sixwaaaay/temp-text)
[![Test](https://github.com/sixwaaaay/temp-text/workflows/Test/badge.svg)](https://github.com/sixwaaaay/temp-text/workflows/Test/badge.svg)

## 接口

- 共享文本 POST `/share`

| 参数名  | 位置 | 类型   | 说明       |
| ------- | ---- | ------ | ---------- |
| content | form | string | 待分享文本 |

- 查询文本 GET `/query`

| 参数名 | 位置  | 类型   | 说明    |
| ------ | ----- | ------ | ------- |
| tid    | query | string | 文本 ID |

## 运行

config.example 为配置文件示例，修改配置并将文件类型重命名为yaml即可运行

依赖：

1. redis