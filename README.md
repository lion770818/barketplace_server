# 用途
* 簡單交易搓合服務
* 使用ddd框架開發微服務
* marketplace_server 服務 負責 建立帳號 登入帳號 上架商品 取得市場價格.... 等 api
    - /v1/transaction_product 買賣商品api
* transaction_server 服務 負責 接收 rabbit mq 的訊息, 將等待搓合訂單, 進入搓合系統, 配對成功後, 更新db或redis
    - 搓合 cmd/transaction_server/main.go

# API List
- /auth/register 用戶註冊
- /auth/login 用戶登錄
- /v1/create_product 上架新商品
- /v1/get_market_price 取得市場行情 ( 並且儲存到 redis快取上)
- /v1/transaction_product 買商品 或 賣商品

# DB Table List
範例儲存在 /sql/Dump_test_db
- backpack 用戶商品背包, 持有商品儲存在此
- transaction 用戶交易清單
- user 用戶資料表
- product 產品資料表

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
