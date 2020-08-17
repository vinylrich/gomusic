package rest

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/ajtwoddltka/gomusic/backend/src/dblayer"
	"github.com/ajtwoddltka/gomusic/backend/src/models"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

type HandlerInterface interface { //핸들러 모음
	GetProducts(c *gin.Context)
	GetPromos(c *gin.Context)
	AddUser(c *gin.Context)
	SignIn(c *gin.Context)
	SignOut(c *gin.Context)
	GetOrders(c *gin.Context)
	Charge(c *gin.Context)
}
type Handler struct {
	db dblayer.DBLayer
}

func NewHandler() (*Handler, error) { //handler 생성자에 데이터베이스 연결 코드 추가
	db, err := dblayer.NewORM("mysql", "user:toor@localhost:3300/gomusic")
	if err != nil {
		return nil, err
	}
	return &Handler{
		db: db,
	}, nil
}

func (h *Handler) GetMainPage(c *gin.Context) {
	log.Println("Main page....")
	fmt.Fprintf(c.Writer, "Main page for secure API!!")
}

func (h *Handler) GetProducts(c *gin.Context) { //상품 항목을 불러옴

	if h.db == nil { //핸들러의 db가 비어있으면 반환 할 값이 없음
		return
	}
	products, err := h.db.GetAllProducts()
	if err != nil { //error가 발생했다면 오류를 출력하고 종료

		//첫 번째 인자는 HTTP 상태 코드,두 번째는 응답의 바디
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	//에러가 발생하지 않았다면 데이터베이스에서 읽은 상품 목록을 반환.
	c.JSON(http.StatusOK, products)
}

func (h *Handler) GetPromos(c *gin.Context) { //프로모션 항목을 불러옴

	if h.db == nil { //핸들러의 db가 비어있으면 반환 할 값이 없음
		return
	}
	promos, err := h.db.GetPromos()
	if err != nil { //error가 발생했다면 오류를 출력하고 종료

		//첫 번째 인자는 HTTP 상태 코드,두 번째는 응답의 바디
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	//에러가 발생하지 않았다면 데이터베이스에서 읽은 상품 목록을 반환.
	c.JSON(http.StatusOK, promos)
}

//로그인
func (h *Handler) SignIn(c *gin.Context) {
	if h.db == nil {
		return
	}
	var customer models.Customer
	err := c.ShouldBindJSON(&customer)
	customer, err = h.db.SignInUser(customer)
	if err != nil { //StatusInternalServerError = http 서버 에러 처리
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	customer, err = h.db.SignInUser(customer.Email, customer.Pass)
	if err != nil {
		//잘못된 패스워드인 경우 forbidden http 에러 반환
		if err == dblayer.ErrINVALIDPASSWORD {
			return
		}
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, customer)
}

//신규 사용자 추가
func (h *Handler) AddUser(c *gin.Context) {
	if h.db == nil {
		return
	}
	var customer models.Customer
	err := c.ShouldBindJSON(&customer)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	customer, err = h.db.AddUser(customer)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, customer)
}

//로그아웃 요청
func (h *Handler) SignOut(c *gin.Context) {
	if h.db == nil {
		return
	}
	p := c.Param("id")
	//id를 문자->정수형으로 변환
	id, err := strconv.Atoi(p)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err = h.db.SignOutUserById(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
}
func (h *Handler) GetOrders(c *gin.Context) {
	if h.db == nil {
		return
	}
	p := c.Param("id")
	//id를 문자->정수형으로 변환
	id, err := strconv.Atoi(p)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	orders, err = h.db.GetCustomerOrdersByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, orders)
}

//신용카드 결제
func (h *Handler) Charge(c *gin.Context) {
	if h.db == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"body": "server database error"})
		return
	}
	request := struct {
		models.Order
		Remember    bool   `json:"rememberCard"`
		UseExisting bool   `json:"useExisting"`
		Token       string `json:"token"`
	}{}
	err := c.ShouldBindJSON(&request)

	if err != nil {
		c.JSON(http.StatusBadRequest, request)
		return
	}
	stripe.Key = "sk_test_4eC39HqLyjWDarjtT1zdp7dc"
	chargeP := &stripe.ChargeParams{

		Amount:      stripe.Int64(int64(request.Price)),
		Currency:    stripe.String("usd"),
		Description: stripe.String("GoMusic charge..."),
	}
	stripeCustomerID := ""
	if request.UseExisting {
		//지정된 카드를 사용하는 경우라면
		log.Println("Getting creadit card id...")
		//스트라이프 아이디 db조회 메서드
		stripeCustomerID, err = h.db.GetCreditCardCID(request.CustomerID) //스트라이프 사용자 id 조회
		if err != nil {
			log.Println(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	} else {
		cp := &stripe.CustomerParams{}
		cp.SetSource(request.Token)
		customer, err := customer.New(cp)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error", err.Error()})
			return
		}
	}
	stripeCustomerID = customer.ID
	if request.Remember {
		//스트라이브 사용자 id를 저장하고 데이터베이스에 저장된 사용자 id와 연결
		err.h.db.SaveCreditCardForCustomer(request.CustomerID, stripeCustomerID) //스트라이프 사용자 id저장
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error", error.Error()})
			return
		}
	}
	return
}

//middleware
func MyCustomMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		//요청을 처리하기 전에 실행할 미들웨어

		c.Set("v", "123")
		//c.Get("v") 변수 값을 확인할 수 있다
		c.Next()

		status := c.Writer.Status()

	}
}
func MyCustomLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		fmt.Println("****************************")
		c.Next()
		fmt.Println("****************************")

	}
}
func RunAPIWithHandler(address string, h HandlerInterface) error {
	r := gin.Default()
	r.Use(MyCustomLogger())
	// r:=gin.New()//두 개 이상의 미들웨어가 필요하면 USE()메스드의 인자에 추가
	// r.Use(MyCustomLogger()...)
}
