package models

import (
	"time"

	"github.com/jinzhu/gorm"
)

//models : 데이터베이스 모델 코딩
type Product struct { //제품 데이터 모델
	//구조체 태그라는 문법으로 해당 필드가 json 문서에서 어떤 키에 해당하는 지 나타내는 문법
	gorm.Model
	Image       string  `json:"img"`
	ImagAlt     string  `json:"imgalt"`
	Price       float64 `json:"price"`
	Promotion   float64 `json:"promotion"` //sql.NullFloat64
	ProductName string  `gorm:"column:production" json:"productname"`
	Description string
}

func (Product) TableName() string {
	return "products"
}

type Customer struct { //고객 데이터 모델
	gorm.Model
	FirstName string  `gorm:"column:firstname" json:"firstname"`
	LastName  string  `gorm:"column:lastname" json:"lastname"`
	Email     string  `gorm:"column:email" json:"email"`
	Pass      string  `json:"password"`
	LoggedIn  bool    `gorm:"column:loggedin" json:"loggedin"`
	Orders    []Order `json:"orders"`
}
type Order struct { //주문 데이터 모델
	gorm.Model
	Product
	Customer
	CustomerID   int       `gorm:"customer_id"`
	ProductID    int       `gorm:"product_id"`
	Price        float64   `gorm:"price" json:"sell_price"`
	PurchaseDate time.Time `gorm:"column:purchase_date" json:"purchase_date"`
}

func (Order) TableName() string {
	return "orders"
}

/*

	사이트의 주요 기능
	1. 상품 목록 조회
	2. 프로모션 목록 조회
	3. 사용자 이름과 성으로 정보 검색
	4. 사용자 id로 고객 정보 검색
	5. 상품 id로 상품 정보 검색
	6. 신규 사용자 등록
	7. 데이터베이스에 로그인 계정 마킹
	8. 데이터베이스에 로그아웃 계정 마킹
	9. 사용자 id로 주문 내역 조회

*/
