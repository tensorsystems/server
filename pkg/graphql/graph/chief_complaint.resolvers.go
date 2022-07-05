package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	graph_models "github.com/tensoremr/server/pkg/graphql/graph/model"
	"github.com/tensoremr/server/pkg/models"
	"github.com/tensoremr/server/pkg/repository"
	deepCopy "github.com/ulule/deepcopier"
)

func (r *mutationResolver) SaveChiefComplaint(ctx context.Context, input graph_models.ChiefComplaintInput) (*models.ChiefComplaint, error) {
	var chiefComplaint models.ChiefComplaint
	deepCopy.Copy(&input).To(&chiefComplaint)

	var hpiComponentRepository repository.HpiComponentRepository
	hpiComponents, err := hpiComponentRepository.GetByIds(input.HpiComponentIds)
	if err != nil {
		return nil, err
	}

	chiefComplaint.HPIComponents = hpiComponents

	var chiefComplaintRepository repository.ChiefComplaintRepository
	if err := chiefComplaintRepository.Save(&chiefComplaint); err != nil {
		return nil, err
	}

	return &chiefComplaint, nil
}

func (r *mutationResolver) UpdateChiefComplaint(ctx context.Context, input graph_models.ChiefComplaintUpdateInput) (*models.ChiefComplaint, error) {
	var chiefComplaint models.ChiefComplaint
	deepCopy.Copy(&input).To(&chiefComplaint)

	var hpiComponentRepository repository.HpiComponentRepository
	hpiComponents, err := hpiComponentRepository.GetByIds(input.HpiComponentIds)
	if err != nil {
		return nil, err
	}

	var chiefComplaintRepository repository.ChiefComplaintRepository
	chiefComplaint.HPIComponents = hpiComponents

	if err := chiefComplaintRepository.Update(&chiefComplaint); err != nil {
		return nil, err
	}

	if err := chiefComplaintRepository.Get(&chiefComplaint, chiefComplaint.ID); err != nil {
		return nil, err
	}

	return &chiefComplaint, nil
}

func (r *mutationResolver) DeleteChiefComplaint(ctx context.Context, id int) (bool, error) {
	var repository repository.ChiefComplaintRepository

	if err := repository.Delete(id); err != nil {
		return false, err
	}

	return true, nil
}

func (r *mutationResolver) SavePatientChiefComplaint(ctx context.Context, input graph_models.ChiefComplaintInput) (*models.ChiefComplaint, error) {
	var patientChartRepository repository.PatientChartRepository

	var patientChart models.PatientChart
	if err := patientChartRepository.Get(&patientChart, input.PatientChartID); err != nil {
		return nil, err
	}

	var chiefComplaintRepository repository.ChiefComplaintRepository
	var chiefComplaint models.ChiefComplaint
	chiefComplaint.Title = input.Title
	chiefComplaint.PatientChartID = input.PatientChartID
	if err := chiefComplaintRepository.Save(&chiefComplaint); err != nil {
		return nil, err
	}

	return &chiefComplaint, nil
}

func (r *mutationResolver) DeletePatientChiefComplaint(ctx context.Context, id int) (bool, error) {
	var repository repository.ChiefComplaintRepository
	if err := repository.Delete(id); err != nil {
		return false, err
	}

	return true, nil
}

func (r *queryResolver) ChiefComplaints(ctx context.Context, page models.PaginationInput, filter *graph_models.ChiefComplaintFilter) (*graph_models.ChiefComplaintConnection, error) {
	var f models.ChiefComplaint
	if filter != nil {
		deepCopy.Copy(filter).To(&f)
	}

	var repository repository.ChiefComplaintRepository
	chiefComplaints, count, err := repository.GetAll(page, &f)

	if err != nil {
		return nil, err
	}

	edges := make([]*graph_models.ChiefComplaintEdge, len(chiefComplaints))

	for i, entity := range chiefComplaints {
		e := entity

		edges[i] = &graph_models.ChiefComplaintEdge{
			Node: &e,
		}
	}

	pageInfo, totalCount := GetPageInfo(chiefComplaints, count, page)
	return &graph_models.ChiefComplaintConnection{PageInfo: pageInfo, Edges: edges, TotalCount: totalCount}, nil
}

func (r *queryResolver) SearchChiefComplaints(ctx context.Context, searchTerm string, page models.PaginationInput) (*graph_models.ChiefComplaintConnection, error) {
	var repository repository.ChiefComplaintRepository
	chiefComplaints, count, err := repository.Search(page, searchTerm)

	if err != nil {
		return nil, err
	}

	edges := make([]*graph_models.ChiefComplaintEdge, len(chiefComplaints))

	for i, entity := range chiefComplaints {
		e := entity

		edges[i] = &graph_models.ChiefComplaintEdge{
			Node: &e,
		}
	}

	pageInfo, totalCount := GetPageInfo(chiefComplaints, count, page)
	return &graph_models.ChiefComplaintConnection{PageInfo: pageInfo, Edges: edges, TotalCount: totalCount}, nil
}
