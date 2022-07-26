# coordinate

获取世界所有城市经纬度数据，然后持久化在存储系统中。

目前支持的地区有：

|地区|数据来源|历史快照|
|----|----|----|
|中国|[高德开放平台](https://lbs.amap.com/api/webservice/guide/api/district)|[`region/china.json`](https://github.com/gogroup/coordinate/blob/main/region/china.json)

目前支持的存储系统以及使用样例：

- `mysql`

  ```shell
  # 拉取中国地区数据并且将数据存入 mysql
  $ coordinate --storage.type mysql --storage.mysql.dsn {YOUR_DSN} --amap.key {YOUR_AMAP_KEY}
  ```

## TODO

- 国外城市数据；
- 台湾没有二级城市数据；
- 项目用法；

## 可能有用的资料

- 中国城市经纬度查询: http://www.hao828.com/ChaXun/ZhongGuoChengShiJingWeiDu
- 国内国外邮编查询: https://www.nowmsg.com/
