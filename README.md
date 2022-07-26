# coordinate

获取世界城市经纬度，并且持久化在存储系统中。

目前支持的地区有：

|地区|数据来源|历史快照|
|----|----|----|
|中国|[高德开放平台](https://lbs.amap.com/api/webservice/guide/api/district)|[`region/china.json`](https://github.com/gogroup/coordinate/blob/main/region/china.json)

目前支持的存储系统以及命令样例：

- `mysql`

  ```shell
  # 拉取中国地区数据并且将数据存入 mysql
  $ coordinate --storage.type mysql --storage.mysql.dsn {YOUR_DSN} --amap.key {YOUR_AMAP_KEY}
  ```

支持从快照中获取数据，而不是从线上获取，使用 `--region.from-snapshots` 开启此功能。

## TODO

- 国外城市数据；
- 台湾没有二级城市数据；

## 可能有用的资料

- 中国城市经纬度查询: http://www.hao828.com/ChaXun/ZhongGuoChengShiJingWeiDu
- 国内国外邮编查询: https://www.nowmsg.com/
