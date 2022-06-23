package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/tensoremr/server/pkg/graphql/graph/generated"
	"github.com/tensoremr/server/pkg/graphql/graph/model"
	"github.com/tensoremr/server/pkg/middleware"
	"github.com/tensoremr/server/pkg/repository"
)

func (r *mutationResolver) CreateTodo(ctx context.Context, input model.NewTodo) (*model.Todo, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) Todos(ctx context.Context) ([]*model.Todo, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) GetHealthCheck(ctx context.Context) (*model.HealthCheckReport, error) {
	var entity repository.User

	pingErr := entity.Ping()
	if pingErr != nil {
		return &model.HealthCheckReport{
			Health: "NOT",
			Db:     false,
		}, pingErr
	}

	return &model.HealthCheckReport{Health: "YES", Db: true}, nil
}

func (r *queryResolver) ReceptionHomeStats(ctx context.Context) (*model.HomeStats, error) {
	var entity repository.Appointment
	scheduled, checkedIn, checkedOut, err := entity.ReceptionHomeStats()

	return &model.HomeStats{
		Scheduled:  scheduled,
		CheckedIn:  checkedIn,
		CheckedOut: checkedOut,
	}, err
}

func (r *queryResolver) NurseHomeStats(ctx context.Context) (*model.HomeStats, error) {
	var entity repository.Appointment
	scheduled, checkedIn, checkedOut, err := entity.NurseHomeStats()

	return &model.HomeStats{
		Scheduled:  scheduled,
		CheckedIn:  checkedIn,
		CheckedOut: checkedOut,
	}, err
}

func (r *queryResolver) PhysicianHomeStats(ctx context.Context) (*model.HomeStats, error) {
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

	var entity repository.Appointment
	scheduled, checkedIn, checkedOut, err := entity.PhysicianHomeStats(user.ID)

	return &model.HomeStats{
		Scheduled:  scheduled,
		CheckedIn:  checkedIn,
		CheckedOut: checkedOut,
	}, err
}

func (r *queryResolver) CurrentDateTime(ctx context.Context) (*time.Time, error) {
	currentTime := time.Now()
	return &currentTime, nil
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
