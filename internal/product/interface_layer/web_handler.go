package interface_layer

import (
	"marketplace_server/internal/common/logs"
	application_product "marketplace_server/internal/product/application_layer"
	"marketplace_server/internal/product/model"
	"marketplace_server/internal/servers/web/response"
	"net/http"

	"github.com/gin-gonic/gin"
)

// [interface層]
// 管理web使用的api
type ProductHandler struct {
	ProductApp application_product.ProductAppInterface
}

func NewProducHandler(productApp application_product.ProductAppInterface) *ProductHandler {
	return &ProductHandler{
		ProductApp: productApp,
	}
}

// PingExample godoc
// @Summary 建立新商品
// @Description create new products onto the market
// @Schemes
// @Tags user
// @Accept json
// @Produce json
// @Param			message	body	model.C2S_ProductCreate		true		"要上架的商品"
// @Failure     500		{object}	response.HTTPError
// @Failure     400		{object}	response.HTTPError
// @Router /v1/create_product [post]
func (u *ProductHandler) CreateProduct(c *gin.Context) {

	var err error
	req := &model.C2S_ProductCreate{}

	// 解析参数
	if err = c.ShouldBindJSON(req); err != nil {
		response.Err(c, http.StatusBadRequest, err.Error())
		return
	}

	// 转化为领域对象 + 参数验证
	registerParams, err := req.ToDomain()
	if err != nil {
		logs.Errorf("[Register] failed, err: %v", err)
		response.Err(c, http.StatusBadRequest, err.Error())
		return
	}

	// 呼叫應用層 建立商品
	err = u.ProductApp.CreateProduct(registerParams)
	if err != nil {
		response.Err(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Ok(c)
}

// 買商品
func (u *ProductHandler) PurchaseProduct(c *gin.Context) {

	var err error
	req := &model.C2S_PurchaseProduct{}

	// 解析参数
	if err = c.ShouldBindJSON(req); err != nil {
		response.Err(c, http.StatusBadRequest, err.Error())
		return
	}

	// 转化为领域对象 + 参数验证
	registerParams, err := req.ToDomain()
	if err != nil {
		logs.Errorf("[Register] failed, err: %+v", err)
		response.Err(c, http.StatusBadRequest, err.Error())
		return
	}

	// 呼叫應用層
	err = u.ProductApp.CreateProduct(registerParams)
	if err != nil {
		response.Err(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Ok(c)
}

/*
// PingExample godoc
// @Summary 取得市場價格
// @Description get the latest market price
// @Schemes
// @Tags user
// @Accept json
// @Produce json
// @Param  message	body				{object}  	model.C2S_MarketPrice 				true 		"用戶名"
// @Success 	200 	{object} 	model.S2C_MarketPrice
// @Failure     500		{object}	response.HTTPError
// @Failure     400		{object}	response.HTTPError
// @Router /v1/get_market_price [get]
*/
func (u *ProductHandler) GetMarketPrice(c *gin.Context) {

	var err error
	req := &model.C2S_MarketPrice{}

	// 解析参数
	if err = c.ShouldBindJSON(req); err != nil {
		response.Err(c, http.StatusBadRequest, err.Error())
		return
	}

	// 转化为领域对象 + 参数验证
	registerParams, err := req.ToDomain()
	if err != nil {
		logs.Errorf("[Register] failed, err: %+v", err)
		response.Err(c, http.StatusBadRequest, err.Error())
		return
	}

	// 呼叫應用層 取得市場價格
	marketPriceList, _, err := u.ProductApp.GetMarketPrice(registerParams)
	if err != nil {
		response.Err(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Ok(c, marketPriceList)
}
