package user

import (
	"marketplace_server/internal/common/logs"
	"marketplace_server/internal/servers/web/response"
	"marketplace_server/internal/user/model"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	UserApp UserAppInterface
}

func NewUserHandler(userApp UserAppInterface) *UserHandler {
	return &UserHandler{
		UserApp: userApp,
	}
}

// 用戶登入
func (u *UserHandler) Login(c *gin.Context) {
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
		logs.Errorf("[Login] verify failed, err: %w", err)
		response.Err(c, http.StatusInternalServerError, err.Error())
		return
	}

	// 调用应用层
	user, err := u.UserApp.Login(loginParams)
	if err != nil {
		logs.Errorf("[Login] failed, err: %w", err)
		response.Err(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Ok(c, user)
}

// 獲取用戶訊息
func (u *UserHandler) UserInfo(c *gin.Context) {
	userID := c.GetInt64(UserIDKey)

	logs.Debugf("userID:%v", userID)

	userInfo, err := u.UserApp.Get(userID)
	if err != nil {
		logs.Errorf("[UserInfo] failed, err: %w", err)
		response.Err(c, http.StatusInternalServerError, err.Error())
		return
	}

	// 返回用户信息
	response.Ok(c, userInfo)
}

// 用戶注册
func (u *UserHandler) Register(c *gin.Context) {
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
		logs.Errorf("[Register] failed, err: %w", err)
		response.Err(c, http.StatusBadRequest, err.Error())
		return
	}

	// 调用应用层
	user, err := u.UserApp.Register(registerParams)
	if err != nil {
		response.Err(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Ok(c, user)
}

// 用戶交易
func (u *UserHandler) Transfer(c *gin.Context) {
	var err error
	req := &model.C2S_Transfer{}

	// 解析参数
	if err = c.ShouldBindJSON(req); err != nil {
		response.Err(c, http.StatusBadRequest, err.Error())
		return
	}

	// 轉化為領域對象 + 參數驗證
	err = req.Verify()
	if err != nil {
		response.Err(c, http.StatusInternalServerError, err.Error())
		return
	}

	// 取得來源用戶ID
	fromUserID := c.GetInt64(UserIDKey)

	// 調用應用層
	err = u.UserApp.Transfer(fromUserID, req.ToUserID, req.Amount, req.Currency)
	if err != nil {
		response.Err(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Ok(c)
}
