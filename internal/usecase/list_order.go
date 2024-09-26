package usecase

import (
	"github.com/devfullcycle/20-CleanArch/internal/entity"
)

type ListOrderUseCase struct {
	OrderRepository entity.OrderRepositoryInterface
}

func NewListOrderUseCase(OrderRepository entity.OrderRepositoryInterface) *ListOrderUseCase {
	return &ListOrderUseCase{OrderRepository: OrderRepository}
}

func (c *ListOrderUseCase) Execute() ([]OrderOutputDTO, error) {
	rows, err := c.OrderRepository.ListAll()
	if err != nil {
		return []OrderOutputDTO{}, err
	}

	var result []OrderOutputDTO
	for _, row := range rows {
		result = append(result, OrderOutputDTO{
			ID:         row.ID,
			Price:      row.Price,
			Tax:        row.Tax,
			FinalPrice: row.FinalPrice,
		})
	}

	return result, nil
}
