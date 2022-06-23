package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/tensoremr/server/pkg/graphql/graph/model"
	"github.com/tensoremr/server/pkg/middleware"
	"github.com/tensoremr/server/pkg/repository"
	deepCopy "github.com/ulule/deepcopier"
)

func (r *mutationResolver) NewAppointment(ctx context.Context, input model.AppointmentInput) (*repository.Appointment, error) {
	var appointment repository.Appointment
	deepCopy.Copy(&input).To(&appointment)

	if err := appointment.CreateNewAppointment(input.BillingID, input.InvoiceNo); err != nil {
		return nil, err
	}

	return &appointment, nil
}

func (r *mutationResolver) UpdateAppointment(ctx context.Context, input model.AppointmentUpdateInput) (*repository.Appointment, error) {
	var entity repository.Appointment
	deepCopy.Copy(&input).To(&entity)

	if err := entity.Update(); err != nil {
		return nil, err
	}

	return &entity, nil
}

func (r *mutationResolver) DeleteAppointment(ctx context.Context, id int) (bool, error) {
	var entity repository.Appointment

	err := entity.Delete(id)
	if err != nil {
		return false, err
	}

	return true, nil
}

func (r *mutationResolver) SaveAppointmentStatus(ctx context.Context, input model.AppointmentStatusInput) (*repository.AppointmentStatus, error) {
	var entity repository.AppointmentStatus
	deepCopy.Copy(&input).To(&entity)

	err := entity.Save()
	if err != nil {
		return nil, err
	}

	return &entity, nil
}

func (r *mutationResolver) UpdateAppointmentStatus(ctx context.Context, input model.AppointmentStatusInput, id int) (*repository.AppointmentStatus, error) {
	var entity repository.AppointmentStatus

	deepCopy.Copy(&input).To(&entity)

	entity.ID = id

	_, err := entity.Update()
	if err != nil {
		return nil, err
	}

	return &entity, nil
}

func (r *mutationResolver) DeleteAppointmentStatus(ctx context.Context, id int) (bool, error) {
	var entity repository.AppointmentStatus

	err := entity.Delete(id)
	if err != nil {
		return false, err
	}

	return true, nil
}

func (r *queryResolver) Appointment(ctx context.Context, id int) (*repository.Appointment, error) {
	var entity repository.Appointment
	if err := entity.GetWithDetails(id); err != nil {
		return nil, err
	}

	var history repository.PatientHistory
	if err := history.GetByPatientID(entity.Patient.ID); err != nil {
		return nil, err
	}

	entity.Patient.PatientHistory = history

	return &entity, nil
}

func (r *queryResolver) Appointments(ctx context.Context, page repository.PaginationInput, filter *model.AppointmentFilter) (*model.AppointmentConnection, error) {
	var f repository.Appointment

	if filter != nil {
		deepCopy.Copy(filter).To(&f)
	}

	appointments, count, err := f.GetAll(page, &f)

	if err != nil {
		return nil, err
	}

	edges := make([]*model.AppointmentEdge, len(appointments))

	for i, entity := range appointments {
		e := entity

		edges[i] = &model.AppointmentEdge{
			Node: &e,
		}
	}

	pageInfo, totalCount := GetPageInfo(appointments, count, page)
	return &model.AppointmentConnection{PageInfo: pageInfo, Edges: edges, TotalCount: totalCount}, nil
}

func (r *queryResolver) AppointmentStatuses(ctx context.Context, page repository.PaginationInput) (*model.AppointmentStatusConnection, error) {
	var entity repository.AppointmentStatus
	entities, count, err := entity.GetAll(page)

	if err != nil {
		return nil, err
	}

	edges := make([]*model.AppointmentStatusEdge, len(entities))

	for i, entity := range entities {
		e := entity

		edges[i] = &model.AppointmentStatusEdge{
			Node: &e,
		}
	}

	pageInfo, totalCount := GetPageInfo(entities, count, page)
	return &model.AppointmentStatusConnection{PageInfo: pageInfo, Edges: edges, TotalCount: totalCount}, nil
}

func (r *queryResolver) FindAppointmentsByPatientAndRange(ctx context.Context, patientID int, start time.Time, end time.Time) ([]*repository.Appointment, error) {
	var entity repository.Appointment
	entities, err := entity.FindAppointmentsByPatientAndRange(patientID, start, end)

	if err != nil {
		return entities, err
	}

	return entities, err
}

func (r *queryResolver) PatientsAppointmentToday(ctx context.Context, patientID int, checkedIn bool) (*repository.Appointment, error) {
	var entity repository.Appointment
	appointment, _ := entity.PatientsAppointmentToday(patientID, &checkedIn)
	return &appointment, nil
}

func (r *queryResolver) FindTodaysAppointments(ctx context.Context, page repository.PaginationInput, searchTerm *string) (*model.AppointmentConnection, error) {
	var entity repository.Appointment
	appointments, count, err := entity.FindTodaysAppointments(page, searchTerm)

	if err != nil {
		return nil, err
	}

	edges := make([]*model.AppointmentEdge, len(appointments))

	for i, entity := range appointments {
		e := entity

		edges[i] = &model.AppointmentEdge{
			Node: &e,
		}
	}

	pageInfo, totalCount := GetPageInfo(appointments, count, page)
	return &model.AppointmentConnection{PageInfo: pageInfo, Edges: edges, TotalCount: totalCount}, nil
}

func (r *queryResolver) FindTodaysCheckedInAppointments(ctx context.Context, page repository.PaginationInput, searchTerm *string) (*model.AppointmentConnection, error) {
	var entity repository.Appointment

	appointments, count, err := entity.FindTodaysCheckedInAppointments(page, searchTerm, []string{})
	if err != nil {
		return nil, err
	}

	edges := make([]*model.AppointmentEdge, len(appointments))

	for i, entity := range appointments {
		e := entity

		edges[i] = &model.AppointmentEdge{
			Node: &e,
		}
	}

	pageInfo, totalCount := GetPageInfo(appointments, count, page)
	return &model.AppointmentConnection{PageInfo: pageInfo, Edges: edges, TotalCount: totalCount}, nil
}

func (r *queryResolver) SearchAppointments(ctx context.Context, page repository.PaginationInput, input repository.AppointmentSearchInput) (*model.AppointmentConnection, error) {
	var entity repository.Appointment
	appointments, count, err := entity.SearchAppointments(page, input)

	if err != nil {
		return nil, err
	}

	edges := make([]*model.AppointmentEdge, len(appointments))

	for i, entity := range appointments {
		e := entity

		edges[i] = &model.AppointmentEdge{
			Node: &e,
		}
	}

	pageInfo, totalCount := GetPageInfo(appointments, count, page)
	return &model.AppointmentConnection{PageInfo: pageInfo, Edges: edges, TotalCount: totalCount}, nil
}

func (r *queryResolver) GetUserAppointments(ctx context.Context, page repository.PaginationInput, searchTerm *string, visitType *string, subscriptions *bool) (*model.AppointmentConnection, error) {
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

	var appointments []repository.Appointment
	var count int64
	var aErr error

	isPhysician := false
	for _, e := range user.UserTypes {
		if e.Title == "Physician" {
			isPhysician = true
		}
	}

	var visitTypes []string

	if visitType != nil {
		if *visitType == "Outpatient" {
			visitTypes = append(visitTypes, "Sick Visit", "Follow-Up", "Check-Up", "Referral")
		} else if *visitType == "Surgeries" {
			visitTypes = append(visitTypes, "Surgery")
		} else if *visitType == "Treatments" {
			visitTypes = append(visitTypes, "Treatment")
		} else if *visitType == "Post-Ops" {
			visitTypes = append(visitTypes, "Post-Op")
		}
	}

	if *subscriptions {
		var queueSubscription repository.QueueSubscription

		if err := queueSubscription.GetByUserId(user.ID); err != nil {
			return nil, err
		}

		patientQueues := queueSubscription.Subscriptions

		var appointmentRepo repository.Appointment
		for _, patientQueue := range patientQueues {
			var ids []int

			if err := json.Unmarshal([]byte(patientQueue.Queue.String()), &ids); err != nil {
				return nil, err
			}

			a, c, e := appointmentRepo.FindByUserSubscriptions(ids, searchTerm, visitTypes, page)

			var orderedAppointments []repository.Appointment
			for _, appointment := range a {
				aTemp := appointment

				aTemp.QueueName = patientQueue.QueueName
				aTemp.QueueID = patientQueue.ID
				orderedAppointments = append(orderedAppointments, aTemp)
			}

			appointments = append(appointments, orderedAppointments...)
			count += c
			aErr = e
		}
	} else {
		if isPhysician {
			appointments, count, aErr = entity.FindByProvider(page, searchTerm, visitTypes, user.ID)
			if aErr != nil {
				return nil, aErr
			}
		} else {
			appointments, count, aErr = entity.FindTodaysCheckedInAppointments(page, searchTerm, visitTypes)
			if aErr != nil {
				return nil, aErr
			}
		}
	}

	edges := make([]*model.AppointmentEdge, len(appointments))

	for i, entity := range appointments {
		e := entity

		edges[i] = &model.AppointmentEdge{
			Node: &e,
		}
	}

	pageInfo, totalCount := GetPageInfo(appointments, count, page)
	return &model.AppointmentConnection{PageInfo: pageInfo, Edges: edges, TotalCount: totalCount}, nil
}

func (r *queryResolver) PayForConsultation(ctx context.Context, patientID int, date *time.Time) (bool, error) {
	var entity repository.Appointment
	shouldPay, err := entity.PayForConsultation(patientID, date)
	return shouldPay, err
}
