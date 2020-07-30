package rest

import (
	"github.com/gin-gonic/gin"
)

type HandlerInterface interface {
	GetProducts(c *gin.Context)
	GetPromos(c *gin.Context)
	AddUser(c *gin.Context)
	SignIn(c *gin.Context)
	SignOut(c *gin.Context)
	GetOrders(c *gin.Context)
	Charge(c *gin.Context)
}
type Handler struct {
	db dblayer.DBlayer
}

func NewHandler() (*Handler, error) {
	return new(Handler), nil //Handler 객체에 대한 포인터
}
