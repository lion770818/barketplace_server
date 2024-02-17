package user

import (
	"encoding/json"
	"errors"
	"fmt"
	"marketplace_server/internal/bill"
	bill_model "marketplace_server/internal/bill/model"
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

	Transfer(fromUserID, toUserID int64, amount decimal.Decimal, currencyStr string) error
	TransactionProduct(pirchase *model.ProductTransactionParams) error // 買商品
}

type UserApp struct {
	userRepo        UserRepo
	authRepo        AuthInterface
	transferService TransferService
	rateService     RateService
	transactionApp  bill.TransactionAppInterface
	productAPP      application_product.ProductAppInterface // 產品應用層
}

func NewUserApp(userRepo UserRepo, authRepo AuthInterface, transactionRepo bill.TransactionRepo, productAPP application_product.ProductAppInterface) UserAppInterface {
	return &UserApp{
		userRepo:        userRepo,
		authRepo:        authRepo,
		transferService: NewTransferService(),
		rateService:     NewRateService(),
		transactionApp:  bill.NewTransactionApp(transactionRepo),
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

// 轉帳(兩人互轉) 等待廢棄
func (u *UserApp) Transfer(fromUserID, toUserID int64, amount decimal.Decimal, toCurrency string) error {
	// 讀取db用戶數據 (來源)
	fromUser, err := u.userRepo.GetUserInfo(fromUserID)
	if err != nil {
		return err
	}

	// 讀取db用戶數據 (目的)
	toUser, err := u.userRepo.GetUserInfo(toUserID)
	if err != nil {
		return err
	}

	// 讀取匯率
	rate, err := u.rateService.GetRate(fromUser.Currency, toCurrency)
	if err != nil {
		return err
	}

	//判斷

	// 轉帳
	err = u.transferService.Transfer(fromUser, toUser, amount, rate)
	if err != nil {
		return err
	}

	// 保存轉帳號金幣回DB
	u.userRepo.Save(fromUser)
	u.userRepo.Save(toUser)

	// 建立交易單
	transaction := &bill_model.Transaction{
		TransactionID: fmt.Sprintf("%d-%d-%s-%d", fromUser.UserID, toUser.UserID, toCurrency, time.Now().UnixNano()), // 交易單號
		FromUserID:    fromUser.UserID,
		ToUserID:      toUser.UserID,
		Amount:        amount,
		Currency:      toCurrency,
	}
	err = u.transactionApp.CreateTransaction(transaction)
	if err != nil {
		return err
	}

	return nil
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

	// 讀取 redis 目前市場價格
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

	logs.Debugf("productName:%v, marketPriceRedis:%v  rate:%v",
		transactionParams.ProductName, marketPriceRedis, rate.Get().String())

	// 取得買或賣的數量
	operateCount := decimal.NewFromInt(int64(transactionParams.OperateCount))

	switch model.TransferMode(transactionParams.TransferMode) {
	case model.Purchase: // 買
		// 計算 購買商品的價格 = redis 的商品價格 * 操作數量 * 匯率
		productNeedPrice := marketPriceRedis.Amount.Mul(operateCount).Mul(rate.Get())
		logs.Debugf("用戶的錢:%s, 操作數量:%v, 匯率:%v 購買商品的價格:%s",
			fromUser.Amount.String(), operateCount, rate.Get().String(), productNeedPrice.String())
		//判斷用戶是否足夠錢買
		if !fromUser.Amount.GreaterThan(productNeedPrice) {
			errMsg := fmt.Errorf("不夠錢買 %s < %s", fromUser.Amount.String(), productNeedPrice.String())
			logs.Warnf("err:%v", errMsg)
			return errMsg
		}
	case model.Sell: // 賣
		// todo:撈取db 看賣家是否有足夠數量
	default:
		return fmt.Errorf("transferMode fail:%v", transactionParams.TransferMode)
	}

	// 時間戳
	transactionParams.TimeStamp = time.Now().UnixNano()

	var cmd model.Notify_Cmd
	switch model.TransferMode(transactionParams.TransferMode) {
	case model.Purchase: // 買
		cmd = model.Notify_Cmd_Purchase
	case model.Sell:
		cmd = model.Notify_Cmd_Sell
	}

	// 寫進message queue 給搓合微服務 transaction_engine
	productTransactionNotify := model.ProductTransactionNotify{
		//Cmd:  model.Notify_Cmd_Purchase,
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
		return nil
	}

	logs.Debugf("成功發送到mq exchangeName:%s, routeKey:%s",
		model.TransactionExchange, model.BindKeyPurchaseProduct)

	// 轉帳
	// err = u.transferService.Transfer(fromUser, toUser, amount, rate)
	// if err != nil {
	// 	return err
	// }

	// // 保存轉帳號金幣回DB
	// u.userRepo.Save(fromUser)
	// u.userRepo.Save(toUser)

	// // 建立交易單
	// bill := &bill_model.Transaction{
	// 	TransactionID: fmt.Sprintf("%d-%d-%s-%d", fromUser.ID, toUser.ID, toCurrency, time.Now().UnixNano()), // 交易單號
	// 	FromUserID:    fromUser.ID,
	// 	ToUserID:      toUser.ID,
	// 	Amount:        amount,
	// 	Currency:      toCurrency,
	// }
	// err = u.billApp.CreateBill(bill)
	// if err != nil {
	// 	return err
	// }

	return nil
}
