package interface_layer

import (
	"marketplace_server/internal/common/logs"
	"marketplace_server/internal/servers/web/response"
	"marketplace_server/internal/user/model"

	application_product "marketplace_server/internal/product/application_layer"
	application_user "marketplace_server/internal/user/application_layer"

	model_bill "marketplace_server/internal/bill/model"
	"net/http"

	"github.com/gin-gonic/gin"
)

// [interface層]
// 管理web使用的api
type UserHandler struct {
	UserApp    application_user.UserAppInterface
	ProductApp application_product.ProductAppInterface
}

func NewUserHandler(userApp application_user.UserAppInterface, productApp application_product.ProductAppInterface) *UserHandler {
	return &UserHandler{
		UserApp:    userApp,
		ProductApp: productApp,
	}
}

// PingExample godoc
// @Summary 用戶登入
// @Description user logsin this system, returns user token
// @Schemes
// @Tags user
// @Accept json
// @Produce json
// @Param			message	body	model.C2S_Login		true		"要登入的帳號"
// @Success 	200 	{object} 	model.S2C_Login
// @Failure     500		{object}	response.HTTPError
// @Failure     400		{object}	response.HTTPError
// @Router /auth/login [post]
func (u *UserHandler) Login(c *gin.Context) {
	logPrefix := "Login"
	var err error
	req := &model.C2S_Login{}

	// 解析参数
	if err = c.ShouldBindJSON(req); err != nil {
		//response.Err(c, http.StatusBadRequest, err.Error())
		response.ErrFromSwagger(c, http.StatusBadRequest, err.Error())
		return
	}

	// 轉換爲領域物件 + 參數驗證
	loginParams, err := req.ToDomain()
	if err != nil {
		logs.Errorf("%s verify failed, err: %+v", logPrefix, err)
		//response.Err(c, http.StatusInternalServerError, err.Error())
		response.ErrFromSwagger(c, http.StatusInternalServerError, err.Error())
		return
	}

	// 呼叫應用層
	user, err := u.UserApp.Login(loginParams)
	if err != nil {
		logs.Errorf("[Login] failed, err: %+v", err)
		//response.Err(c, http.StatusInternalServerError, err.Error())
		response.ErrFromSwagger(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Ok(c, user)
}

// PingExample godoc
// @Summary 獲取用戶訊息
// @Description get user info from this system, returns user token
// @Schemes
// @Tags user
// @Accept json
// @Produce json
// @Param 				username path  	string 				true 		"用戶名"
// @Success 	200 	{object} 	model.S2C_UserInfo
// @Failure     500		{object}	response.HTTPError
// @Failure     400		{object}	response.HTTPError
// @Router /v1/login/{username} [get]
func (u *UserHandler) UserInfo(c *gin.Context) {
	logPrefix := "Register"
	userID := c.GetInt64(UserIDKey)

	logs.Debugf("userID:%v", userID)

	// 應用層 取得用戶資訊
	userInfo, err := u.UserApp.GetUserInfo(userID)
	if err != nil {
		logs.Errorf("%s failed, err: %+v", logPrefix, err)
		response.Err(c, http.StatusInternalServerError, err.Error())
		return
	}

	// 返回用户信息
	response.Ok(c, userInfo)
}

// PingExample godoc
// @Summary 用戶注册
// @Description user register this system, returns user token
// @Schemes
// @Tags user
// @Accept json
// @Produce json
// @Param			message	body	model.C2S_Register		true		"要註冊的帳號"
// @Success 	200 	{object} 	model.S2C_Login
// @Failure     500		{object}	response.HTTPError
// @Failure     400		{object}	response.HTTPError
// @Router /auth/register [post]
func (u *UserHandler) Register(c *gin.Context) {
	logPrefix := "Register"
	var err error
	req := &model.C2S_Register{}

	// 解析参数
	if err = c.ShouldBindJSON(req); err != nil {
		response.Err(c, http.StatusBadRequest, err.Error())
		return
	}

	// 转化为领域对象 + 参数验证
	registerParams, err := req.ToDomain()
	if err != nil {
		logs.Errorf("%s failed, err: %+v", logPrefix, err)
		response.Err(c, http.StatusBadRequest, err.Error())
		return
	}

	// 呼叫應用層
	user, err := u.UserApp.Register(registerParams)
	if err != nil {
		response.Err(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Ok(c, user)
}

// PingExample godoc
// @Summary 買商品 賣商品
// @Description buy or sell product
// @Schemes
// @Tags user
// @Accept json
// @Produce json
// @Param			message	body	model.C2S_TransactionProduct		true		"要交易的商品"
// @Success 	200 	{object} 	model_bill.Transaction
// @Failure     500		{object}	response.HTTPError
// @Failure     400		{object}	response.HTTPError
// @Router /auth/transaction_product [post]
func (u *UserHandler) TransactionProduct(c *gin.Context) {

	logPrefix := "transactionProduct"
	req := &model.C2S_TransactionProduct{}
	var err error

	// 解析参数
	if err = c.ShouldBindJSON(req); err != nil {
		response.Err(c, http.StatusBadRequest, err.Error())
		return
	}

	// 转化为领域对象 + 参数验证
	transactionProductParams, err := req.ToDomain()
	if err != nil {
		logs.Errorf("%s failed, err: %+v", logPrefix, err)
		response.Err(c, http.StatusBadRequest, err.Error())
		return
	}

	// todo 後續 增加 同一用戶封包太頻繁交易就阻擋

	// 呼叫應用層 買商品 / 賣商品
	var transaction *model_bill.Transaction
	transaction, err = u.UserApp.TransactionProduct(transactionProductParams)
	if err != nil {
		response.Err(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Ok(c, transaction)
}

// PingExample godoc
// @Summary 取消 買商品 賣商品
// @Description buy or sell product
// @Schemes
// @Tags user
// @Accept json
// @Produce json
// @Param			message	body	model.C2S_CancelProduct		true		"要交易的商品"
// @Success 	200 	{object} 	model_bill.Transaction
// @Failure     500		{object}	response.HTTPError
// @Failure     400		{object}	response.HTTPError
// @Router /auth/cancel_product [post]
func (u *UserHandler) CancelProduct(c *gin.Context) {

	logPrefix := "cancelProduct"
	req := &model.C2S_CancelProduct{}
	var err error

	// 解析参数
	if err = c.ShouldBindJSON(req); err != nil {
		response.Err(c, http.StatusBadRequest, err.Error())
		return
	}

	// 转化为领域对象 + 参数验证
	transactionProductParams, err := req.ToDomain()
	if err != nil {
		logs.Errorf("%s failed, err: %+v", logPrefix, err)
		response.Err(c, http.StatusBadRequest, err.Error())
		return
	}

	// todo 後續 增加 同一用戶封包太頻繁交易就阻擋

	// 呼叫應用層 取消交易
	err = u.UserApp.CancelProduct(transactionProductParams)
	if err != nil {
		response.Err(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Ok(c)
}
