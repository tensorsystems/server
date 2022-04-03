package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	"github.com/Tensor-Systems/tensoremr-server/pkg/graphql/graph/model"
	"github.com/Tensor-Systems/tensoremr-server/pkg/repository"
	deepCopy "github.com/ulule/deepcopier"
)

func (r *mutationResolver) SaveChiefComplaint(ctx context.Context, input model.ChiefComplaintInput) (*repository.ChiefComplaint, error) {
	var chiefComplaint repository.ChiefComplaint
	deepCopy.Copy(&input).To(&chiefComplaint)

	var hpiComponent repository.HpiComponent
	hpiComponents, err := hpiComponent.GetByIds(input.HpiComponentIds)
	if err != nil {
		return nil, err
	}

	chiefComplaint.HPIComponents = hpiComponents

	if err := chiefComplaint.Save(); err != nil {
		return nil, err
	}

	return &chiefComplaint, nil
}

func (r *mutationResolver) UpdateChiefComplaint(ctx context.Context, input model.ChiefComplaintUpdateInput) (*repository.ChiefComplaint, error) {
	var chiefComplaint repository.ChiefComplaint
	deepCopy.Copy(&input).To(&chiefComplaint)

	var hpiComponent repository.HpiComponent
	hpiComponents, err := hpiComponent.GetByIds(input.HpiComponentIds)
	if err != nil {
		return nil, err
	}

	chiefComplaint.HPIComponents = hpiComponents

	if err := chiefComplaint.Update(); err != nil {
		return nil, err
	}

	chiefComplaint.Get(chiefComplaint.ID)

	return &chiefComplaint, nil
}

func (r *mutationResolver) DeleteChiefComplaint(ctx context.Context, id int) (bool, error) {
	var entity repository.ChiefComplaint

	if err := entity.Delete(id); err != nil {
		return false, err
	}

	return true, nil
}

func (r *mutationResolver) SavePatientChiefComplaint(ctx context.Context, input model.ChiefComplaintInput) (*repository.ChiefComplaint, error) {
	var patientChart repository.PatientChart

	if err := patientChart.Get(input.PatientChartID); err != nil {
		return nil, err
	}

	var chiefComplaint repository.ChiefComplaint
	chiefComplaint.Title = input.Title
	chiefComplaint.PatientChartID = input.PatientChartID
	if err := chiefComplaint.Save(); err != nil {
		return nil, err
	}

	return &chiefComplaint, nil
}

func (r *mutationResolver) DeletePatientChiefComplaint(ctx context.Context, id int) (bool, error) {
	var chiefComplaint repository.ChiefComplaint
	if err := chiefComplaint.Delete(id); err != nil {
		return false, err
	}

	return true, nil
}

func (r *queryResolver) ChiefComplaints(ctx context.Context, page repository.PaginationInput, filter *model.ChiefComplaintFilter) (*model.ChiefComplaintConnection, error) {
	var f repository.ChiefComplaint
	if filter != nil {
		deepCopy.Copy(filter).To(&f)
	}

	var entity repository.ChiefComplaint
	chiefComplaints, count, err := entity.GetAll(page, &f)

	if err != nil {
		return nil, err
	}

	edges := make([]*model.ChiefComplaintEdge, len(chiefComplaints))

	for i, entity := range chiefComplaints {
		e := entity

		edges[i] = &model.ChiefComplaintEdge{
			Node: &e,
		}
	}

	pageInfo, totalCount := GetPageInfo(chiefComplaints, count, page)
	return &model.ChiefComplaintConnection{PageInfo: pageInfo, Edges: edges, TotalCount: totalCount}, nil
}

func (r *queryResolver) SearchChiefComplaints(ctx context.Context, searchTerm string, page repository.PaginationInput) (*model.ChiefComplaintConnection, error) {
	var entity repository.ChiefComplaint
	chiefComplaints, count, err := entity.Search(page, searchTerm)

	if err != nil {
		return nil, err
	}

	edges := make([]*model.ChiefComplaintEdge, len(chiefComplaints))

	for i, entity := range chiefComplaints {
		e := entity

		edges[i] = &model.ChiefComplaintEdge{
			Node: &e,
		}
	}

	pageInfo, totalCount := GetPageInfo(chiefComplaints, count, page)
	return &model.ChiefComplaintConnection{PageInfo: pageInfo, Edges: edges, TotalCount: totalCount}, nil
}
