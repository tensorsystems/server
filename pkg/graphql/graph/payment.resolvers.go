package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"errors"
	"fmt"

	"github.com/tensoremr/server/pkg/graphql/graph/generated"
	"github.com/tensoremr/server/pkg/graphql/graph/model"
	"github.com/tensoremr/server/pkg/middleware"
	"github.com/tensoremr/server/pkg/repository"
	deepCopy "github.com/ulule/deepcopier"
)

func (r *mutationResolver) SavePayment(ctx context.Context, input model.PaymentInput) (*repository.Payment, error) {
	var entity repository.Payment

	deepCopy.Copy(&input).To(&entity)

	err := entity.Save()
	if err != nil {
		return nil, err
	}

	return &entity, nil
}

func (r *mutationResolver) UpdatePayment(ctx context.Context, input model.PaymentInput) (*repository.Payment, error) {
	var entity repository.Payment

	deepCopy.Copy(&input).To(&entity)

	if err := entity.Update(); err != nil {
		return nil, err
	}

	return &entity, nil
}

func (r *mutationResolver) DeletePayment(ctx context.Context, id int) (bool, error) {
	var entity repository.Payment

	err := entity.Delete(id)
	if err != nil {
		return false, err
	}

	return true, nil
}

func (r *mutationResolver) ConfirmPayment(ctx context.Context, id int, invoiceNo string) (*repository.Payment, error) {
	var entity repository.Payment
	entity.ID = id
	entity.Status = repository.PaidPaymentStatus
	entity.InvoiceNo = invoiceNo

	if err := entity.Update(); err != nil {
		return nil, err
	}

	return &entity, nil
}

func (r *mutationResolver) ConfirmPayments(ctx context.Context, ids []int, invoiceNo string) (bool, error) {
	var entity repository.Payment
	if err := entity.BatchUpdate(ids, repository.Payment{Status: repository.PaidPaymentStatus, InvoiceNo: invoiceNo}); err != nil {
		return false, err
	}
	return true, nil
}

func (r *mutationResolver) RequestPaymentWaiver(ctx context.Context, paymentID int, patientID int) (*repository.Payment, error) {
	gc, err := middleware.GinContextFromContext(ctx)
	if err != nil {
		return nil, err
	}

	email := gc.GetString("email")
	if len(email) == 0 {
		return nil, errors.New("Cannot find user")
	}

	var user repository.User
	err = user.GetByEmail(email)
	if err != nil {
		return nil, err
	}

	var entity repository.Payment
	if err := entity.RequestWaiver(paymentID, patientID, user.ID); err != nil {
		return nil, err
	}

	return &entity, nil
}

func (r *mutationResolver) RequestPaymentWaivers(ctx context.Context, ids []int, patientID int) (bool, error) {
	gc, err := middleware.GinContextFromContext(ctx)
	if err != nil {
		return false, err
	}

	email := gc.GetString("email")
	if len(email) == 0 {
		return false, errors.New("Cannot find user")
	}

	var user repository.User
	err = user.GetByEmail(email)
	if err != nil {
		return false, err
	}

	var entity repository.Payment
	if err := entity.RequestWaiverBatch(ids, patientID, user.ID); err != nil {
		return false, err
	}

	return true, nil
}

func (r *paymentResolver) Status(ctx context.Context, obj *repository.Payment) (string, error) {
	return string(obj.Status), nil
}

func (r *queryResolver) Payments(ctx context.Context, page repository.PaginationInput) (*model.PaymentConnection, error) {
	panic(fmt.Errorf("not implemented"))
}

// Payment returns generated.PaymentResolver implementation.
func (r *Resolver) Payment() generated.PaymentResolver { return &paymentResolver{r} }

type paymentResolver struct{ *Resolver }
