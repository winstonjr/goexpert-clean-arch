package database

import (
	"gorm.io/gorm"
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/winstonjr/goexpert-clean-arch/internal/entity"

	// sqlite3
	"gorm.io/driver/sqlite"
)

type OrderRepositoryGormTestSuite struct {
	suite.Suite
	Db *gorm.DB
}

func (suite *OrderRepositoryGormTestSuite) SetupTest() {
	dsn := "file::memory:?cache=shared"
	db, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{})
	suite.NoError(err)
	if err != nil {
		panic(err)
	}
	err = db.AutoMigrate(&GORMOrder{})
	if err != nil {
		panic(err)
	}
	suite.Db = db
}

func (suite *OrderRepositoryGormTestSuite) TearDownTest() {
	sqlDB, err := suite.Db.DB()
	suite.NoError(err)
	sqlDB.Close()
}

func TestSuiteOrderGorm(t *testing.T) {
	suite.Run(t, new(OrderRepositoryGormTestSuite))
}

func (suite *OrderRepositoryGormTestSuite) TestGivenAnOrder_WhenSave_ThenShouldSaveOrder() {
	order, err := entity.NewOrder("456", 20.0, 4.0)
	suite.NoError(err)
	suite.NoError(order.CalculateFinalPrice())
	repo := NewOrderRepositoryGorm(suite.Db)
	err = repo.Save(order)
	suite.NoError(err)

	orderResult := &GORMOrder{ID: order.ID}
	result := suite.Db.First(orderResult)

	suite.NoError(result.Error)
	suite.Equal(int64(1), result.RowsAffected)
	suite.Equal(order.ID, orderResult.ID)
	suite.Equal(order.Price, orderResult.Price)
	suite.Equal(order.Tax, orderResult.Tax)
	suite.Equal(order.FinalPrice, orderResult.FinalPrice)
}

func (suite *OrderRepositoryGormTestSuite) TestGetAllRecordsInDatabase() {
	baseOrder := &GORMOrder{ID: "a", Price: 1, Tax: 1, FinalPrice: 2}
	result := suite.Db.Create(baseOrder)
	suite.NoError(result.Error)
	baseOrder = &GORMOrder{ID: "b", Price: 1, Tax: 2, FinalPrice: 3}
	result = suite.Db.Create(baseOrder)
	suite.NoError(result.Error)
	baseOrder = &GORMOrder{ID: "c", Price: 2, Tax: 3, FinalPrice: 5}
	result = suite.Db.Create(baseOrder)
	suite.NoError(result.Error)

	repo := NewOrderRepositoryGorm(suite.Db)
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
