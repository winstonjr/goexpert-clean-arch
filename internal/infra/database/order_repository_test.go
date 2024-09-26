package database

import (
	"database/sql"
	"testing"

	"github.com/devfullcycle/20-CleanArch/internal/entity"
	"github.com/stretchr/testify/suite"

	// sqlite3
	_ "github.com/mattn/go-sqlite3"
)

type OrderRepositoryTestSuite struct {
	suite.Suite
	Db *sql.DB
}

func (suite *OrderRepositoryTestSuite) SetupSuite() {
	db, err := sql.Open("sqlite3", ":memory:")
	suite.NoError(err)
	db.Exec("CREATE TABLE orders (id varchar(255) NOT NULL, price float NOT NULL, tax float NOT NULL, final_price float NOT NULL, PRIMARY KEY (id))")
	suite.Db = db
}

func (suite *OrderRepositoryTestSuite) TearDownTest() {
	suite.Db.Close()
}

func TestSuite(t *testing.T) {
	suite.Run(t, new(OrderRepositoryTestSuite))
}

func (suite *OrderRepositoryTestSuite) TestGivenAnOrder_WhenSave_ThenShouldSaveOrder() {
	order, err := entity.NewOrder("123", 10.0, 2.0)
	suite.NoError(err)
	suite.NoError(order.CalculateFinalPrice())
	repo := NewOrderRepository(suite.Db)
	err = repo.Save(order)
	suite.NoError(err)

	var orderResult entity.Order
	err = suite.Db.QueryRow("Select id, price, tax, final_price from orders where id = ?", order.ID).
		Scan(&orderResult.ID, &orderResult.Price, &orderResult.Tax, &orderResult.FinalPrice)

	suite.NoError(err)
	suite.Equal(order.ID, orderResult.ID)
	suite.Equal(order.Price, orderResult.Price)
	suite.Equal(order.Tax, orderResult.Tax)
	suite.Equal(order.FinalPrice, orderResult.FinalPrice)
}

func (suite *OrderRepositoryTestSuite) TestGetAllRecordsInDatabase() {
	stmt, err := suite.Db.Prepare("INSERT INTO orders (id, price, tax, final_price) VALUES (?, ?, ?, ?)")
	suite.NoError(err)
	_, err = stmt.Exec("a", 1, 1, 2)
	suite.NoError(err)
	_, err = stmt.Exec("b", 1, 2, 3)
	suite.NoError(err)
	_, err = stmt.Exec("c", 2, 3, 5)
	suite.NoError(err)

	repo := NewOrderRepository(suite.Db)
	orders, err := repo.ListAll()
	suite.NoError(err)
	suite.Equal(3, len(orders))

	suite.Equal(orders[0].ID, "a")
	suite.Equal(orders[0].Price, float64(1))
	suite.Equal(orders[0].Tax, float64(1))
	suite.Equal(orders[0].FinalPrice, float64(2))

	suite.Equal(orders[1].ID, "b")
	suite.Equal(orders[1].Price, float64(1))
	suite.Equal(orders[1].Tax, float64(2))
	suite.Equal(orders[1].FinalPrice, float64(3))

	suite.Equal(orders[2].ID, "c")
	suite.Equal(orders[2].Price, float64(2))
	suite.Equal(orders[2].Tax, float64(3))
	suite.Equal(orders[2].FinalPrice, float64(5))
}
