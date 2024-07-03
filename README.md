# 用途

- 簡單交易搓合服務
- 使用 ddd 框架開發微服務
- marketplace_server 服務 負責 建立帳號 登入帳號 上架商品 取得市場價格.... 等 api
  - /v1/transaction_product 買賣商品 api
- transaction_server 服務 負責 接收 rabbit mq 的訊息, 將等待搓合訂單, 進入搓合系統, 配對成功後, 更新 db 或 redis
  - 搓合 cmd/transaction_server/main.go

# API List

- /auth/register 用戶註冊
- /auth/login 用戶登錄
- /v1/create_product 上架新商品
- /v1/get_market_price 取得市場行情 ( 並且儲存到 redis 快取上)
- /v1/transaction_product 買商品 或 賣商品

# DB Table List

範例儲存在 /sql/Dump_test_db

- backpack 用戶商品背包, 持有商品儲存在此
- transaction 用戶交易清單
- user 用戶資料表
- product 產品資料表

## 參考範例

依赖的环境

- golang
- docker

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

# Swagger

# api 文件 製作方式

go get -u github.com/swaggo/swag/cmd/swag
go install github.com/swaggo/swag/cmd/swag@latest

go get -u github.com/swaggo/gin-swagger
go get -u github.com/swaggo/files

go get github.com/swaggo/swag/example/celler/httputil
go get github.com/swaggo/swag/example/celler/model

go install github.com/swaggo/swag/cmd/swag@latest

swag init

swag 是執行檔 有問題就去設定
linux => PATH 例如: export PATH=$HOME/go/bin:$PATH
windos => 環境變數內設定

@Param：參數訊息，用空格分隔的參數。param name,param type,data type,is mandatory?,comment attribute(optional) 1.參數名稱 2.參數類型，可以有的值是 formData、query、path、body、header，formData 表示是 post 請求的數據， query 表示帶在 url 之後的參數，path 表示請求路徑上得參數，例如上面例子裡面的 key，body 表示是一個 raw 資料請求，header 表示帶在 header 資訊中得參數。 3.參數類型 4.是否必須 5.註釋
例如：

// @Param name query string true "用户姓名"

[常用註解格式]("https://blog.csdn.net/qq_38371367/article/details/123005909")

[swagger 教學]("https://igouist.github.io/post/2021/05/newbie-4-swagger/")

## 如果出現 Fetch error Internal Server Error http://localhost:8080/swagger/doc.json

請在專案上面 import \_ "bito_group/docs"

## 註解編譯 cannot find type definition: httputil.HTTPError

請 import 在編譯一次
"github.com/swaggo/swag/example/celler/httputil"
"github.com/swaggo/swag/example/celler/model"

## 教學文

[手把手詳細教你如何使用 go-swagger 文檔]("https://juejin.cn/post/7126802030944878600")
