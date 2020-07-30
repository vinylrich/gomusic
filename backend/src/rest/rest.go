package rest

import (
	"github.com/gin-gonic/gin"
)

func RunAPI(address string) error {
	r := gin.Default()
	r.GET("/relativepath/to/url", func(c *gin.Context) { //익명함수 gin.Context는 요청 확인과 처리, 응답에 필요한 기능 제공
		//logic
	})
	r.GET("/product", func(c *gin.Context) {
		//상품 목록 반환
	})
	r.GET("/promos", func(c *gin.Context) {
		//프로모션 목록 반환
	})
	r.POST("/users/signin", func(c *gin.Context) {
		//사용자 로그인 POST
	})
	r.POST("/users", func(c *gin.Context) {
		//사용자 추가
	})
	//사용자
	r.POST("/users/:id/signout", func(c *gin.Context) {
		//사용자 로그아웃 요청
	})

	r.GET("/user/:id/orders", func(c *gin.Context) {
		//
	})
}
