# 用途
* 簡單交易搓合服務
* 使用ddd框架開發微服務
* marketplace_server 服務 負責 建立帳號 登入帳號 上架商品 取得市場價格.... 等 api
* transaction_engine 服務 負責 接收 rabbit mq 的訊息, 將等待搓合訂單, 進入搓合系統, 配對成功後, 更新db或redis


## 參考範例

依赖的环境

* golang
* docker

```bash
# 下载项目
git clone git@github.com:dengjiawen8955/ddd_demo.git  && cd ddd_demo
# 准备环境 (启动mysql, redis)
docker-compose up -d
# 准备数据库 (创建数据库, 创建表)
make exec.sql
# 启动项目
make
```
