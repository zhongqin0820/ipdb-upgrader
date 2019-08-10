# ipdb-upgrader
[![Build Status](https://travis-ci.org/zhongqin0820/ipdb-upgrader.svg?branch=master)](https://travis-ci.org/zhongqin0820/ipdb-upgrader)
[![Coverage Status](https://coveralls.io/repos/github/zhongqin0820/ipdb-upgrader/badge.svg?branch=master)](https://coveralls.io/github/zhongqin0820/ipdb-upgrader?branch=master)

定制化提升免费版<https://ipip.net>的IPDB格式地级市精度IP离线库展示的内容。

# 目录结构
```
.
├── utils                       # 存储提升内容文件
│   ├── codeCity.json           # 大陆城市代码
│   ├── codeCountry.json        # 国家代码
│   ├── codeRegion.json         # 大陆地区代码
│   └── ipipfree.ipdb           # ipip离线数据库文件
├── .coveralls.yml
├── .gitignore
├── .travis.yml
├── city.go
├── city_test.go
├── cityUpgrader.go             # upgrader
├── cityUpgrader_test.go        # upgrader
├── reader.go
├── go.mod
├── LICENSE
└── README.md
```

# 结果对比
提升内容包括：
- 所在大洲名
- 国家/地区电话号码前缀
- 中国大陆行政区划代码（省级）
- 中国大陆行政区划代码（市级）

## Upgrade前
省略输出部分空字段。
```json
{
    "country_name": "中国",
    "region_name": "辽宁",
    "city_name": "铁岭"
}
```

## Upgrade后
```json
{
    "continent_code": "亚洲",
    "country_name": "中国",
    "idd_code": "86",
    "region_name": "辽宁",
    "china_admin_code": "210000",
    "city_name": "铁岭",
    "china_city_code": "211200"
}
```

# 基准测试
```
# ipdb-go
BenchmarkCity_Find-4         5000000           352 ns/op
BenchmarkCity_FindMap-4      2000000           646 ns/op
BenchmarkCity_FindInfo-4     1000000          1681 ns/op
# ipdb-upgrader
BenchmarkFindInfo-4           500000          2086 ns/op
```

# 相关链接
## 官方仓库
- [ipipdotnet/ipdb-go](https://github.com/ipipdotnet/ipdb-go)

## 测试IP来源
- [世界各国IP列表](http://ip.yqie.com/world.aspx)
- [中国各省市IP地址列表](http://ip.yqie.com/china.aspx)