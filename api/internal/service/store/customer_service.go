package store

import (
	"context"

	"github.com/yeftaz/susano.id/api/internal/domain/store"
	storeRepo "github.com/yeftaz/susano.id/api/internal/repository/store"
)

type CustomerService struct {
	customerRepo *storeRepo.CustomerRepository
}

func NewCustomerService(customerRepo *storeRepo.CustomerRepository) *CustomerService {
	return &CustomerService{
		customerRepo: customerRepo,
	}
}

// GetByID retrieves a customer by ID
func (s *CustomerService) GetByID(ctx context.Context, id string) (*store.Customer, error) {
	customer, err := s.customerRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Clear password before returning
	customer.Password = ""

	return customer, nil
}

// Update updates a customer profile
func (s *CustomerService) Update(ctx context.Context, id string, email, name string) (*store.Customer, error) {
	var emailPtr, namePtr *string

	if email != "" {
		emailPtr = &email
	}
	if name != "" {
		namePtr = &name
	}

	customer, err := s.customerRepo.Update(ctx, id, emailPtr, namePtr)
	if err != nil {
		return nil, err
	}

	// Clear password before returning
	customer.Password = ""

	return customer, nil
}

// Delete soft deletes a customer
func (s *CustomerService) Delete(ctx context.Context, id string) error {
	return s.customerRepo.Delete(ctx, id)
}

// UpdateAvatar updates customer avatar
func (s *CustomerService) UpdateAvatar(ctx context.Context, id, avatarPath string) error {
	return s.customerRepo.UpdateAvatarPath(ctx, id, avatarPath)
}
