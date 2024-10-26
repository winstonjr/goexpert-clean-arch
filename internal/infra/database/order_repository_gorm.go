package database

import (
	"github.com/winstonjr/goexpert-clean-arch/internal/entity"
	"gorm.io/gorm"
)

type GORMOrder struct {
	ID         string `gorm:"primary_key"`
	Price      float64
	Tax        float64
	FinalPrice float64
	gorm.Model
}

type OrderRepositoryGorm struct {
	Db *gorm.DB
}

func NewOrderRepositoryGorm(db *gorm.DB) *OrderRepositoryGorm {
	return &OrderRepositoryGorm{Db: db}
}

func (r *OrderRepositoryGorm) Save(order *entity.Order) error {
	result := r.Db.Create(&GORMOrder{
		ID:         order.ID,
		Price:      order.Price,
		Tax:        order.Tax,
		FinalPrice: order.FinalPrice,
	})

	return result.Error
}

func (r *OrderRepositoryGorm) GetTotal() (int64, error) {
	var count int64
	result := r.Db.Model(&GORMOrder{}).Count(&count)
	if result.Error != nil {
		return 0, result.Error
	}
	return count, nil
}

func (r *OrderRepositoryGorm) ListAll() ([]*entity.Order, error) {
	var gormorders []*GORMOrder
	result := r.Db.Find(&gormorders)
	if result.Error != nil {
		return nil, result.Error
	}
	var orders = make([]*entity.Order, len(gormorders))
	for i, order := range gormorders {
		orders[i] = &entity.Order{
			ID:         order.ID,
			Price:      order.Price,
			Tax:        order.Tax,
			FinalPrice: order.FinalPrice,
		}
	}
	return orders, nil
}
