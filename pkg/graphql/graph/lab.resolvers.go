package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"errors"
	"time"

	graph_models "github.com/tensoremr/server/pkg/graphql/graph/model"
	"github.com/tensoremr/server/pkg/middleware"
	"github.com/tensoremr/server/pkg/models"
	"github.com/tensoremr/server/pkg/repository"
	deepCopy "github.com/ulule/deepcopier"
)

func (r *mutationResolver) SaveLab(ctx context.Context, input graph_models.LabInput) (*models.Lab, error) {
	var entity models.Lab
	deepCopy.Copy(&input).To(&entity)

	var repository repository.LabRepository
	if err := repository.Save(&entity); err != nil {
		return nil, err
	}

	return &entity, nil
}

func (r *mutationResolver) UpdateLab(ctx context.Context, input graph_models.LabUpdateInput) (*models.Lab, error) {
	var entity models.Lab
	deepCopy.Copy(&input).To(&entity)

	if input.Status != nil {
		entity.Status = models.LabStatus(*input.Status)
	}
	// Images ...
	for _, fileUpload := range input.Images {
		fileName, hashedFileName, hash, ext := HashFileName(fileUpload.Name)
		err := WriteFile(fileUpload.File.File, hashedFileName+"."+ext)
		if err != nil {
			return nil, err
		}

		entity.Images = append(entity.RightEyeImages, models.File{
			ContentType: fileUpload.File.ContentType,
			Size:        fileUpload.File.Size,
			FileName:    fileName,
			Extension:   ext,
			Hash:        hash,
		})
	}

	// Documents
	for _, fileUpload := range input.Documents {
		fileName, hashedFileName, hash, ext := HashFileName(fileUpload.Name)
		err := WriteFile(fileUpload.File.File, hashedFileName+"."+ext)
		if err != nil {
			return nil, err
		}

		entity.Documents = append(entity.Documents, models.File{
			ContentType: fileUpload.File.ContentType,
			Size:        fileUpload.File.Size,
			FileName:    fileName,
			Extension:   ext,
			Hash:        hash,
		})
	}

	var repository repository.LabRepository
	if err := repository.Update(&entity); err != nil {
		return nil, err
	}

	return &entity, nil
}

func (r *mutationResolver) DeleteLab(ctx context.Context, id int) (bool, error) {
	var repository repository.LabRepository

	if err := repository.Delete(id); err != nil {
		return false, err
	}

	return true, nil
}

func (r *mutationResolver) SaveLabType(ctx context.Context, input graph_models.LabTypeInput) (*models.LabType, error) {
	var entity models.LabType
	deepCopy.Copy(&input).To(&entity)

	var billingRepository repository.BillingRepository
	billings, err := billingRepository.GetByIds(input.BillingIds)
	if err != nil {
		return nil, err
	}

	entity.Billings = billings

	var labTypeRepository repository.LabTypeRepository
	if err := labTypeRepository.Save(&entity); err != nil {
		return nil, err
	}

	return &entity, nil
}

func (r *mutationResolver) UpdateLabType(ctx context.Context, input graph_models.LabTypeUpdateInput) (*models.LabType, error) {
	var entity models.LabType
	deepCopy.Copy(&input).To(&entity)

	var billingRepository repository.BillingRepository
	billings, err := billingRepository.GetByIds(input.BillingIds)
	if err != nil {
		return nil, err
	}

	entity.Billings = billings

	var labTypeRepository repository.LabTypeRepository
	if err := labTypeRepository.Update(&entity); err != nil {
		return nil, err
	}

	return &entity, nil
}

func (r *mutationResolver) DeleteLabType(ctx context.Context, id int) (bool, error) {
	var repository repository.LabTypeRepository

	if err := repository.Delete(id); err != nil {
		return false, err
	}

	return true, nil
}

func (r *mutationResolver) DeleteLabRightEyeImage(ctx context.Context, input graph_models.LabDeleteFileInput) (bool, error) {
	var repository repository.LabRepository

	if err := repository.DeleteFile("RightEyeImages", input.LabID, input.FileID); err != nil {
		return false, err
	}

	return true, nil
}

func (r *mutationResolver) DeleteLabLeftEyeImage(ctx context.Context, input graph_models.LabDeleteFileInput) (bool, error) {
	var repository repository.LabRepository

	if err := repository.DeleteFile("LeftEyeImages", input.LabID, input.FileID); err != nil {
		return false, err
	}

	return true, nil
}

func (r *mutationResolver) DeleteLabRightEyeSketch(ctx context.Context, input graph_models.LabDeleteFileInput) (bool, error) {
	var repository repository.LabRepository

	if err := repository.DeleteFile("RightEyeSketches", input.LabID, input.FileID); err != nil {
		return false, err
	}

	return true, nil
}

func (r *mutationResolver) DeleteLabLeftEyeSketch(ctx context.Context, input graph_models.LabDeleteFileInput) (bool, error) {
	var repository repository.LabRepository

	if err := repository.DeleteFile("LeftEyeSketches", input.LabID, input.FileID); err != nil {
		return false, err
	}

	return true, nil
}

func (r *mutationResolver) DeleteLabImage(ctx context.Context, input graph_models.LabDeleteFileInput) (bool, error) {
	var repository repository.LabRepository

	if err := repository.DeleteFile("Images", input.LabID, input.FileID); err != nil {
		return false, err
	}

	return true, nil
}

func (r *mutationResolver) DeleteLabDocument(ctx context.Context, input graph_models.LabDeleteFileInput) (bool, error) {
	var repository repository.LabRepository

	if err := repository.DeleteFile("Documents", input.LabID, input.FileID); err != nil {
		return false, err
	}

	return true, nil
}

func (r *mutationResolver) OrderLab(ctx context.Context, input graph_models.OrderLabInput) (*models.LabOrder, error) {
	// Get current user
	gc, err := middleware.GinContextFromContext(ctx)
	if err != nil {
		return nil, err
	}

	email := gc.GetString("email")
	if len(email) == 0 {
		return nil, errors.New("Cannot find user")
	}

	var userRepository repository.UserRepository
	var user models.User
	if err := userRepository.GetByEmail(&user, email); err != nil {
		return nil, err
	}

	var repository repository.LabOrderRepository

	// Save lab order
	var labOrder models.LabOrder
	if err := repository.Save(&labOrder, input.LabTypeID, input.PatientChartID, input.PatientID, input.BillingIds, user, input.OrderNote, input.ReceptionNote); err != nil {
		return nil, err
	}

	return &labOrder, nil
}

func (r *mutationResolver) ConfirmLabOrder(ctx context.Context, id int, invoiceNo string) (*models.LabOrder, error) {
	var entity models.LabOrder

	var repository repository.LabOrderRepository
	if err := repository.Confirm(&entity, id, invoiceNo); err != nil {
		return nil, err
	}

	return &entity, nil
}

func (r *mutationResolver) UpdateLabOrder(ctx context.Context, input graph_models.LabOrderUpdateInput) (*models.LabOrder, error) {
	var entity models.LabOrder
	deepCopy.Copy(&input).To(&entity)

	if input.Status != nil {
		entity.Status = models.LabOrderStatus(*input.Status)
	}

	var repository repository.LabOrderRepository
	if err := repository.Update(&entity); err != nil {
		return nil, err
	}

	return &entity, nil
}

func (r *mutationResolver) OrderAndConfirmLab(ctx context.Context, input graph_models.OrderAndConfirmLabInput) (*models.LabOrder, error) {
	// Get current user
	gc, err := middleware.GinContextFromContext(ctx)
	if err != nil {
		return nil, err
	}

	email := gc.GetString("email")
	if len(email) == 0 {
		return nil, errors.New("Cannot find user")
	}

	var userRepository repository.UserRepository
	var user models.User
	if err := userRepository.GetByEmail(&user, email); err != nil {
		return nil, err
	}

	var appointmentRepository repository.AppointmentRepository
	var appointment models.Appointment
	if err := appointmentRepository.Get(&appointment, input.AppointmentID); err != nil {
		return nil, err
	}

	var patientChartRepository repository.PatientChartRepository
	var patientChart models.PatientChart
	if err := patientChartRepository.GetByAppointmentID(&patientChart, appointment.ID); err != nil {
		return nil, err
	}

	var labOrderRepository repository.LabOrderRepository
	var labOrder models.LabOrder
	if err := labOrderRepository.Save(&labOrder, input.LabTypeID, patientChart.ID, input.PatientID, input.BillingIds, user, input.OrderNote, ""); err != nil {
		return nil, err
	}

	if err := labOrderRepository.Confirm(&labOrder, labOrder.ID, input.InvoiceNo); err != nil {
		return nil, err
	}

	return &labOrder, nil
}

func (r *queryResolver) Labs(ctx context.Context, page models.PaginationInput, filter *graph_models.LabFilter) (*graph_models.LabConnection, error) {
	var f models.Lab
	if filter != nil {
		deepCopy.Copy(filter).To(&f)
	}

	var repository repository.LabRepository
	entities, count, err := repository.GetAll(page, &f)

	if err != nil {
		return nil, err
	}

	edges := make([]*graph_models.LabEdge, len(entities))

	for i, entity := range entities {
		e := entity

		edges[i] = &graph_models.LabEdge{
			Node: &e,
		}
	}

	pageInfo, totalCount := GetPageInfo(entities, count, page)
	return &graph_models.LabConnection{PageInfo: pageInfo, Edges: edges, TotalCount: totalCount}, nil
}

func (r *queryResolver) LabTypes(ctx context.Context, page models.PaginationInput, searchTerm *string) (*graph_models.LabTypeConnection, error) {
	var repository repository.LabTypeRepository
	result, count, err := repository.GetAll(page, searchTerm)
	if err != nil {
		return nil, err
	}

	edges := make([]*graph_models.LabTypeEdge, len(result))

	for i, entity := range result {
		e := entity

		edges[i] = &graph_models.LabTypeEdge{
			Node: &e,
		}
	}

	pageInfo, totalCount := GetPageInfo(result, count, page)
	return &graph_models.LabTypeConnection{PageInfo: pageInfo, Edges: edges, TotalCount: totalCount}, nil
}

func (r *queryResolver) LabOrder(ctx context.Context, patientChartID int) (*models.LabOrder, error) {
	var entity models.LabOrder

	var repository repository.LabOrderRepository
	if err := repository.GetByPatientChartID(&entity, patientChartID); err != nil {
		return nil, err
	}

	return &entity, nil
}

func (r *queryResolver) SearchLabOrders(ctx context.Context, page models.PaginationInput, filter *graph_models.LabOrderFilter, date *time.Time, searchTerm *string) (*graph_models.LabOrderConnection, error) {
	var f models.LabOrder
	if filter != nil {
		deepCopy.Copy(filter).To(&f)
	}

	if filter.Status != nil {
		f.Status = models.LabOrderStatus(*filter.Status)
	}

	var repository repository.LabOrderRepository
	result, count, err := repository.Search(page, &f, date, searchTerm, false)

	if err != nil {
		return nil, err
	}

	edges := make([]*graph_models.LabOrderEdge, len(result))

	for i, entity := range result {
		e := entity

		edges[i] = &graph_models.LabOrderEdge{
			Node: &e,
		}
	}

	pageInfo, totalCount := GetPageInfo(result, count, page)
	return &graph_models.LabOrderConnection{PageInfo: pageInfo, Edges: edges, TotalCount: totalCount}, nil
}
