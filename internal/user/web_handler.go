package user

import (
	"marketplace_server/internal/common/logs"
	"marketplace_server/internal/servers/web/response"
	"marketplace_server/internal/user/model"

	application_product "marketplace_server/internal/product/application_layer"
	"net/http"

	"github.com/gin-gonic/gin"
)

// [interface層]
// 管理web使用的api
type UserHandler struct {
	UserApp    UserAppInterface
	ProductApp application_product.ProductAppInterface
}

func NewUserHandler(userApp UserAppInterface, productApp application_product.ProductAppInterface) *UserHandler {
	return &UserHandler{
		UserApp:    userApp,
		ProductApp: productApp,
	}
}

// 用戶登入
func (u *UserHandler) Login(c *gin.Context) {
	logPrefix := "Login"
	var err error
	req := &model.C2S_Login{}

	// 解析参数
	if err = c.ShouldBindJSON(req); err != nil {
		response.Err(c, http.StatusBadRequest, err.Error())
		return
	}

	// 轉換爲領域物件 + 參數驗證
	loginParams, err := req.ToDomain()
	if err != nil {
		logs.Errorf("%s verify failed, err: %+v", logPrefix, err)
		response.Err(c, http.StatusInternalServerError, err.Error())
		return
	}

	// 呼叫應用層
	user, err := u.UserApp.Login(loginParams)
	if err != nil {
		logs.Errorf("[Login] failed, err: %+v", err)
		response.Err(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Ok(c, user)
}

// 獲取用戶訊息
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

// 用戶注册
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

// 買商品 賣商品
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
	err = u.UserApp.TransactionProduct(transactionProductParams)
	if err != nil {
		response.Err(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Ok(c)
}

// 取消 買商品 賣商品
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
