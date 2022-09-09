# 临时文本服务

[![codecov](https://codecov.io/gh/sixwaaaay/temp-text/branch/master/graph/badge.svg?token=UwTUzTcS2G)](https://codecov.io/gh/sixwaaaay/temp-text)

## 接口

- 共享文本 POST `/share`

| 参数名     | 位置   | 类型     | 说明    |
|---------|------|--------|-------|
| content | form | string | 待分享文本 |

- 查询文本 GET `/query`

| 参数名    | 位置     | 类型     | 说明     |
|--------|--------|--------|--------|
| tid    | query  | string | 文本ID   |

