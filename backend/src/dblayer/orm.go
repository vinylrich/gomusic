package dblayer

import (
	"errors"

	"github.com/ajtwoddltka/gomusic/backend/src/models"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
)

type DBORM struct {
	*gorm.DB
}

func NewORM(dbname, con string) (*DBORM, error) {
	db, err := gorm.Open(dbname, con)
	return &DBORM{
		DB: db,
	}, err
}

func (db *DBORM) GetAllProducts() (products []models.Product, err error) {
	return products, db.Find(&products).Error //db.Find=select * from "인자"
}

func (db *DBORM) GetPromos() (products []models.Product, err error) {
	return products, db.Where("promotion IS NOT NULL").Find(&products).Error //select * from products where promotion IS NOT NULL
}

func (db *DBORM) GetCustomerByName(firstname, lastname string) (customer models.Customer, err error) {
	return customer, db.Where(&models.Customer{FirstName: firstname, LastName: lastname}).Find(&customer).Error //select * from customers where firstname and lastname
}

func (db *DBORM) GetCustomerByID(id int) (customer models.Customer, err error) {
	return customer, db.First(&customer, id).Error //쿼리의 조건을 만족하는 첫 번째 결과만 반환하는 First 메서드
}

func (db *DBORM) GetProduct(id int) (product models.Product, err error) {
	return product, db.First(&product, id).Error
}

func (db *DBORM) AddUser(customer models.Customer) (models.Customer, error) {
	hashPassword(&customer.Pass)
	customer.LoggedIn = true
	err := db.Create(&customer).Error
	customer.Pass = "" //보내주고 문자열을 보안을 위해 지움
	return customer, err
}
func checkPassword(existingHash, incomingPass string) bool {
	//해시와 패스워드 문자열이 일치하지 않으면 에러 반환
	return bcrypt.CompareHashAndPassword([]byte(existingHash), []byte(incomingPass)) == nil
}

func (db *DBORM) SignInUser(email, pass string) (customer models.Customer, err error) {

	//사용자 행을 나타내는 *gorm.DB 타입
	result := db.Table("Customers").Where(&models.Customer{Email: email}) //질의 결과를 나타내는 구조체 반환
	//loggedin 필드 업데이트
	err = db.First(&customer).Error
	if err != nil {
		return customer, err
	}
	if !checkPassword(customer.Pass, pass) {
		return customer, ErrINVALIDPASSWORD
	}
	customer.Pass = ""
	err = result.Update("loggedin", 1).Error //loggedin==1 -> login complete
	if err != nil {
		return customer, err
	}

	return customer, result.Find(&customer).Error
}

func (db *DBORM) SignOutUserByID(email, pass string) (customer models.Customer, err error) {
	customer = models.Customer{
		Model: gorm.Model{
			ID: uint(id),
		},
	}
	return db.Table("Customers").Where(&customer).Update("loggedin", 0), db.Table("Customers").Where(&customer).Update("loggedin", 0).Error
}
func (db *DBORM) GetCustomerOrdersByID(id int) (orders []models.Order, err error) {
	return orders, db.Table("orders").Select("*").Joins("join customers on customers.id = customer_id").Where("customer_id=?,id").Scan(&orders).Error
	//reference.Joins("join products on products.id")
}
func hashPassword(s *string) error {
	if s == nil {
		return errors.New("Reference provided for hashing password is nil")
	}
	sBytes := []byte(*s)
	//GenerateFromPassword() 메서드는 패스워드 해시를 반환
	hashedBytes, err := bcrypt.GenerateFromPassword(sBytes, bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	*s = string(hashedBytes[:])
	return nil
}
func (db *DBORM) AddOrder(order models.Order) error {
	return db.Create(&order).Error
}
func (db *DBORM) GetCreditCardCID(id int) (string, error) {
	customerWithCCID := struct {
		models.Customer
		CCID string `gorm:"column:cc_customerid"`
	}{}
	return customerWithCCID.CCID, db.First(&customerWithCCID, id).Error
}
func (db *DBORM) SaveCreditCardForCustomer(id int, ccid string) error {
	result := db.Table("customers").Where("id=?", id) //select * from customers where id=id
	return result.Update("cc_customerid", ccid).Error
}
