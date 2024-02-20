package user

import (
	"encoding/json"
	"errors"
	"fmt"
	"marketplace_server/internal/bill"
	model_bill "marketplace_server/internal/bill/model"
	"marketplace_server/internal/common/logs"
	"marketplace_server/internal/common/rabbitmqx"
	application_product "marketplace_server/internal/product/application_layer"
	model_product "marketplace_server/internal/product/model"
	"marketplace_server/internal/user/model"
	"time"

	"github.com/shopspring/decimal"
)

var (
	Error_UserAlreadyExists = errors.New("用户已存在")
	Error_VerifyFailed      = errors.New("验证失败")
)

// [應用層]
type UserAppInterface interface {
	Login(login *model.LoginParams) (*model.S2C_Login, error)
	GetAuthInfo(token string) (*model.AuthInfo, error)
	GetUserInfo(userID int64) (*model.S2C_UserInfo, error)
	Register(register *model.RegisterParams) (*model.S2C_Login, error)

	TransactionProduct(pirchase *model.ProductTransactionParams) error // 買 / 賣 商品
	CancelProduct(pirchase *model.ProductCancelParams) error           // 取消交易
}

// 用戶應用層物件
type UserApp struct {
	userRepo        UserRepo
	authRepo        AuthInterface
	transferService TransferService
	rateService     RateService
	//transactionApp  bill.TransactionAppInterface
	transactionRepo bill.TransactionRepo                    // 交易清單
	productAPP      application_product.ProductAppInterface // 產品應用層
}

func NewUserApp(userRepo UserRepo, authRepo AuthInterface, transactionRepo bill.TransactionRepo, productAPP application_product.ProductAppInterface) UserAppInterface {
	return &UserApp{
		userRepo:        userRepo,
		authRepo:        authRepo,
		transferService: NewTransferService(),
		rateService:     NewRateService(),
		//transactionApp:  bill.NewTransactionApp(transactionRepo),
		transactionRepo: transactionRepo,
		productAPP:      productAPP,
	}
}

// Login
func (u *UserApp) Login(login *model.LoginParams) (*model.S2C_Login, error) {
	// 登录
	user, err := u.userRepo.GetUserByLoginParams(login)
	if err != nil {
		return nil, err
	}

	token := ""
	auth, err := u.authRepo.GetAuthUser(user.UserID)
	if err != nil {

		if err.Error() != "redis: nil" {
			logs.Warnf("獲取用戶快取失敗 userID:%d, err:%v", user.UserID, err)
			return nil, err
		}

		// 快取找不到user, 生成 token
		authInfo := &model.AuthInfo{
			UserID: user.UserID,
			Amount: user.Amount,
		}
		token, err = u.authRepo.Set(authInfo)
		if err != nil {
			return nil, err
		}
	} else {
		token = u.authRepo.GetKey(user.UserID)
	}
	logs.Debugf("auth:%+v", auth)

	return user.ToLoginResp(token), nil
}

// GetAuthInfo 從 token 中 取得用戶資訊
func (u *UserApp) GetAuthInfo(token string) (*model.AuthInfo, error) {
	return u.authRepo.Get(token)
}

// GetUserInfo 取得用戶資訊
func (u *UserApp) GetUserInfo(userID int64) (*model.S2C_UserInfo, error) {
	// 持久層 取得用戶資訊
	user, err := u.userRepo.GetUserInfo(userID)
	if err != nil {
		return nil, err
	}

	// 領域層物件轉換
	return user.ToUserInfo(), nil
}

// Register 注册 + 自动登录
func (u *UserApp) Register(register *model.RegisterParams) (*model.S2C_Login, error) {
	// 检查是否已经注册
	getUser, err := u.userRepo.GetUserByRegisterParams(register)
	if getUser != nil || err == nil {
		return nil, Error_UserAlreadyExists
	}

	// 转换参数
	params, err := register.ToDomain()
	if err != nil {
		return nil, Error_UserAlreadyExists
	}

	// 注册
	user, err := u.userRepo.Save(params)
	if err != nil {
		return nil, err
	}

	// 生成 token
	authInfo := &model.AuthInfo{
		UserID: user.UserID,
	}
	token, err := u.authRepo.Set(authInfo)
	if err != nil {
		return nil, err
	}

	return user.ToLoginResp(token), nil
}

// 買商品 / 賣商品
func (u *UserApp) TransactionProduct(transactionParams *model.ProductTransactionParams) error {

	if transactionParams == nil {
		return fmt.Errorf("transactionParams == nil")
	}

	// 讀取db用戶數據 (來源)
	fromUser, err := u.userRepo.GetUserInfo(transactionParams.UserID)
	if err != nil {
		return err
	}

	// 讀取匯率
	rate, err := u.rateService.GetRate(fromUser.Currency, transactionParams.Currency)
	if err != nil {
		return err
	}

	// 讀取 redis 目前市場價格 ( 橫向調用了 )
	_, dataMap, err := u.productAPP.GetMarketPrice(nil)
	if err != nil {
		return err
	}
	// 解析 redis 資料
	var marketPriceRedis model_product.MarketPriceRedis
	err = json.Unmarshal([]byte(dataMap[transactionParams.ProductName]), &marketPriceRedis)
	if err != nil {
		logs.Errorf("productName:%v, json:%+v, err:%v",
			transactionParams.ProductName, dataMap[transactionParams.ProductName], err)
		return err
	}

	// 取得用戶緩存
	auth, err := u.authRepo.GetAuthUser(transactionParams.UserID)
	if err != nil {
		logs.Errorf("userID:%v err:%v", transactionParams.UserID, err)
		return err
	}

	logs.Debugf("productName:%v, marketPriceRedis:%v  rate:%v, auth:%+v",
		transactionParams.ProductName, marketPriceRedis, rate.Get().String(), auth)

	// 取得買或賣的數量
	operateCount := decimal.NewFromInt(int64(transactionParams.OperateCount))

	// 還沒到搓合階段, 無法知道真實成交價
	var productNeedPrice decimal.Decimal
	switch model.TransferMode(transactionParams.TransferMode) {
	case model.Purchase: // 買單
		// 計算 購買商品的價格 = redis 的商品價格 * 操作數量 * 匯率
		productNeedPrice = marketPriceRedis.Amount.Mul(operateCount).Mul(rate.Get())
		logs.Debugf("用戶的錢:%s, 操作數量:%v, 匯率:%v 購買商品的價格:%s, 商品名稱:%s",
			fromUser.Amount.String(), operateCount.String(), rate.Get().String(), productNeedPrice.String(), transactionParams.ProductName)

		//判斷用戶是否足夠錢買 (使用redis的緩存錢來判斷, db的用戶金額是真實交易時才會異動)
		if !auth.Amount.GreaterThan(productNeedPrice) {
			errMsg := fmt.Errorf("不夠錢買 %s < %s", auth.Amount.String(), productNeedPrice.String())
			logs.Errorf("err:%v", errMsg)
			return errMsg
		}
	case model.Sell: // 賣單
		// todo:撈取db 看賣家是否有足夠數量
	default:
		return fmt.Errorf("transferMode fail mode:%v", transactionParams.TransferMode)
	}

	// 時間戳
	transactionParams.TimeStamp = time.Now().UnixNano()

	// 產生交易ID 格式為 UserID + TransferMode(買或賣) + 流水id
	var id int64
	id, err = u.transactionRepo.GetLastInsterId()
	if err != nil {
		if err.Error() != "record not found" {
			logs.Errorf("getLastInsterId err:%v", err)
			return err
		}
	}
	id++
	transactionId := fmt.Sprintf("%d-%d-%012d", transactionParams.UserID, transactionParams.TransferMode, id)
	transactionParams.TransactionID = transactionId

	// 寫進message queue 給搓合微服務 transaction_engine
	var cmd model.Notify_Cmd
	switch model.TransferMode(transactionParams.TransferMode) {
	case model.Purchase: // 買
		cmd = model.Notify_Cmd_Purchase
	case model.Sell: // 賣
		cmd = model.Notify_Cmd_Sell
	}
	productTransactionNotify := model.ProductTransactionNotify{
		Cmd:  cmd,
		Data: transactionParams,
	}
	mqDataBytes, err := json.Marshal(productTransactionNotify)
	if err != nil {
		return fmt.Errorf("marshal fail err=%v", err)
	}
	err = rabbitmqx.GetMq().PutIntoQueue(model.TransactionExchange, model.BindKeyPurchaseProduct, mqDataBytes)
	if err != nil {
		logs.Errorf("putIntoQueue err:%v, exchange:%v, bindKey:%v",
			err, model.TransactionExchange, model.BindKeyPurchaseProduct)
		return err
	}

	logs.Debugf("成功發送到mq exchangeName:%s, routeKey:%s, transactionParams:%+v",
		model.TransactionExchange, model.BindKeyPurchaseProduct, transactionParams)

	// 寫入db或redis, 狀態設定為 wait 搓合
	transaction := &model_bill.Transaction{
		TransactionID: transactionId,                            // 交易單號
		FromUserID:    transactionParams.UserID,                 // 發起人的用戶ID
		ToUserID:      0,                                        // 交易對象的用戶ID (等交易完成後更新)
		ProductName:   transactionParams.ProductName,            // 產品名稱
		ProductCount:  transactionParams.OperateCount,           // 產品數量
		Amount:        decimal.NewFromFloat(0),                  // 金額 (等交易完成後更新)
		Currency:      transactionParams.Currency,               // 貨幣
		CreatedAt:     time.Now(),                               // 創建時間
		UodateAt:      time.Now(),                               // 更新時間
		Status:        int8(model_bill.Transaction_Status_Wait), // 交易狀態 0:未完成 1:已完成
	}
	if err = u.transactionRepo.Save(transaction); err != nil {
		logs.Errorf("transactionRepo save err:%v", err)
		return err
	}
	logs.Debugf("寫入transaction:%+v", transaction)

	// 更新用戶的緩存 (如果是買單 先預扣)
	auth.Amount = auth.Amount.Sub(productNeedPrice)
	if _, err = u.authRepo.Set(auth); err != nil {
		logs.Errorf("update user cache err:%v", err)
		return err
	}

	return nil
}

// 取消交易單
func (u *UserApp) CancelProduct(cancelParams *model.ProductCancelParams) error {
	if cancelParams == nil {
		return fmt.Errorf("cancel == nil")
	}

	// 讀取db用戶數據 (來源)
	user, err := u.userRepo.GetUserInfo(cancelParams.UserID)
	if err != nil {
		return err
	}
	// 讀取db是否有此交易單
	transaction, err := u.transactionRepo.GetTransactionInfo(cancelParams.TransactionID)
	if err != nil {
		return err
	}

	// 組合通知 mq (todo 放到底層)
	productTransactionNotify := model.ProductTransactionNotify{
		Cmd:  model.Notify_Cmd_Cancel,
		Data: cancelParams,
	}
	mqDataBytes, err := json.Marshal(productTransactionNotify)
	if err != nil {
		return fmt.Errorf("marshal fail err=%v", err)
	}
	err = rabbitmqx.GetMq().PutIntoQueue(model.TransactionExchange, model.BindKeyPurchaseProduct, mqDataBytes)
	if err != nil {
		logs.Errorf("putIntoQueue err:%v, exchange:%v, bindKey:%v",
			err, model.TransactionExchange, model.BindKeyPurchaseProduct)
		return err
	}

	logs.Debugf("成功發送到mq exchangeName:%s, routeKey:%s, cancelParams:%+v, Username:%v, ProductName:%v, Status:%v",
		model.TransactionExchange, model.BindKeyPurchaseProduct, cancelParams, user.Username, transaction.ProductName, transaction.Status)

	return nil
}
