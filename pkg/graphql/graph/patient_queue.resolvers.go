package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/Tensor-Systems/tensoremr-server/pkg/graphql/graph/generated"
	"github.com/Tensor-Systems/tensoremr-server/pkg/graphql/graph/model"
	"github.com/Tensor-Systems/tensoremr-server/pkg/middleware"
	"github.com/Tensor-Systems/tensoremr-server/pkg/repository"
	"gorm.io/datatypes"
)

func (r *mutationResolver) SubscribeQueue(ctx context.Context, userID int, patientQueueID int) (*repository.QueueSubscription, error) {
	var entity repository.QueueSubscription

	if err := entity.Subscribe(userID, patientQueueID); err != nil {
		return nil, err
	}

	return &entity, nil
}

func (r *mutationResolver) UnsubscribeQueue(ctx context.Context, userID int, patientQueueID int) (*repository.QueueSubscription, error) {
	var entity repository.QueueSubscription

	if err := entity.Unsubscribe(userID, patientQueueID); err != nil {
		return nil, err
	}

	return &entity, nil
}

func (r *mutationResolver) SavePatientQueue(ctx context.Context, input model.PatientQueueInput) (*repository.PatientQueue, error) {
	var entity repository.PatientQueue
	entity.QueueName = input.QueueName
	entity.QueueType = input.QueueType

	queue := datatypes.JSON([]byte("[" + strings.Join(input.Queue, ", ") + "]"))
	entity.Queue = queue

	if err := entity.GetByQueueName(input.QueueName); err != nil {
		if err := entity.Save(); err != nil {
			return nil, err
		}
	} else {
		if err := entity.UpdateQueue(input.QueueName, queue); err != nil {
			return nil, err
		}
	}

	return &entity, nil
}

func (r *mutationResolver) DeleteFromQueue(ctx context.Context, patientQueueID int, appointmentID int) (*repository.PatientQueue, error) {
	var entity repository.PatientQueue

	if err := entity.DeleteFromQueue(patientQueueID, appointmentID); err != nil {
		return nil, err
	}

	return &entity, nil
}

func (r *mutationResolver) CheckOutPatient(ctx context.Context, patientQueueID int, appointmentID int) (*repository.PatientQueue, error) {
	var appointment repository.Appointment
	err := appointment.Get(appointmentID)
	if err != nil {
		return nil, err
	}

	// Change appointment status to Checked Out
	var status repository.AppointmentStatus
	status.GetByTitle("Checked-Out")
	appointment.AppointmentStatusID = status.ID
	appointment.CheckedOutTime = time.Now()

	if err := appointment.Update(); err != nil {
		return nil, err
	}

	var patientQueue repository.PatientQueue

	if err := patientQueue.DeleteFromQueue(patientQueueID, appointmentID); err != nil {
		return nil, err
	}

	return &patientQueue, nil
}

func (r *mutationResolver) PushPatientQueue(ctx context.Context, patientQueueID int, appointmentID int, destination model.Destination) (*repository.PatientQueue, error) {
	var entity repository.PatientQueue

	if destination.String() == "PREEXAM" {
		if err := entity.MoveToQueueName(patientQueueID, "Pre-Exam", appointmentID, ""); err != nil {
			return nil, err
		}
	} else if destination.String() == "PREOPERATION" {
		if err := entity.MoveToQueueName(patientQueueID, "Pre-Operation", appointmentID, ""); err != nil {
			return nil, err
		}
	} else if destination.String() == "PHYSICIAN" {
		var appointment repository.Appointment
		if err := appointment.Get(appointmentID); err != nil {
			return nil, err
		}

		var provider repository.User
		if err := provider.Get(appointment.UserID); err != nil {
			return nil, err
		}

		if err := entity.MoveToQueueName(patientQueueID, "Dr. "+provider.FirstName+" "+provider.LastName, appointmentID, "USER"); err != nil {
			return nil, err
		}
	}

	return &entity, nil
}

func (r *mutationResolver) MovePatientQueue(ctx context.Context, appointmentID int, sourceQueueID int, destinationQueueID int) (*repository.PatientQueue, error) {
	var entity repository.PatientQueue

	if err := entity.Move(sourceQueueID, destinationQueueID, appointmentID); err != nil {
		return nil, err
	}

	return &entity, nil
}

func (r *mutationResolver) CheckInPatient(ctx context.Context, appointmentID int, destination model.Destination) (*repository.Appointment, error) {
	var appointment repository.Appointment

	if err := appointment.Get(appointmentID); err != nil {
		return nil, err
	}

	// Change appointment status to Checked In
	var status repository.AppointmentStatus
	if err := status.GetByTitle("Checked-In"); err != nil {
		return nil, err
	}

	var visitType repository.VisitType
	if err := visitType.Get(appointment.VisitTypeID); err != nil {
		return nil, err
	}

	if visitType.Title == "Surgery" {
		var postOpAppointent repository.Appointment

		if err := postOpAppointent.SchedulePostOp(appointment); err != nil {
			return nil, err
		}
	}

	appointment.AppointmentStatusID = status.ID
	checkedInTime := time.Now()
	appointment.CheckedInTime = &checkedInTime

	if err := appointment.Update(); err != nil {
		return nil, err
	}

	// Add to queue
	var patientQueue repository.PatientQueue
	if destination.String() == "PREEXAM" {
		if err := patientQueue.AddToQueue("Pre-Exam", appointmentID, "PREEXAM"); err != nil {
			return nil, err
		}

	} else if destination.String() == "PREOPERATION" {
		if err := patientQueue.AddToQueue("Pre-Operation", appointmentID, "PREOPERATION"); err != nil {
			return nil, err
		}
	} else if destination.String() == "PHYSICIAN" {
		var provider repository.User
		if err := provider.Get(appointment.UserID); err != nil {
			return nil, err
		}

		if err := patientQueue.AddToQueue("Dr. "+provider.FirstName+" "+provider.LastName, appointmentID, "USER"); err != nil {
			return nil, err
		}
	}

	return &appointment, nil
}

func (r *mutationResolver) UpdatePatientQueue(ctx context.Context, appointmentID int, destination *model.Destination) (*repository.PatientQueue, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *patientQueueResolver) Queue(ctx context.Context, obj *repository.PatientQueue) (string, error) {
	return obj.Queue.String(), nil
}

func (r *queryResolver) PatientQueues(ctx context.Context) ([]*model.PatientQueueWithAppointment, error) {
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

	isPhysician := false
	for _, e := range user.UserTypes {
		if e.Title == "Physician" {
			isPhysician = true
		}
	}

	var entity repository.PatientQueue

	var patientQueues []*repository.PatientQueue

	if !isPhysician {
		p, err := entity.GetAll()
		if err != nil {
			return nil, err
		}

		patientQueues = p
	} else {
		p, err := entity.GetAll()
		if err != nil {
			return nil, err
		}

		var t []*repository.PatientQueue

		for _, patientQueue := range p {
			if patientQueue.QueueType == repository.UserQueue && patientQueue.QueueName != "Dr. "+user.FirstName+" "+user.LastName {
				continue
			}

			t = append(t, patientQueue)
		}

		patientQueues = t
	}

	var result []*model.PatientQueueWithAppointment

	var appointmentRepo repository.Appointment
	for _, patientQueue := range patientQueues {
		var ids []int

		if err := json.Unmarshal([]byte(patientQueue.Queue.String()), &ids); err != nil {
			return nil, err
		}

		page := repository.PaginationInput{Page: 0, Size: 1000}

		appointments, _, _ := appointmentRepo.GetByIds(ids, page)
		var orderedAppointments []*repository.Appointment

		for _, id := range ids {
			for _, appointment := range appointments {
				if appointment.ID == id {
					a := appointment
					orderedAppointments = append(orderedAppointments, &a)
				}
			}
		}

		result = append(result, &model.PatientQueueWithAppointment{
			ID:        int(patientQueue.ID),
			QueueName: patientQueue.QueueName,
			QueueType: patientQueue.QueueType,
			Queue:     orderedAppointments,
		})
	}

	return result, err
}

func (r *queryResolver) UserSubscriptions(ctx context.Context, userID int) (*repository.QueueSubscription, error) {
	var entity repository.QueueSubscription

	if err := entity.GetByUserId(userID); err != nil {
		return nil, err
	}

	return &entity, nil
}

// PatientQueue returns generated.PatientQueueResolver implementation.
func (r *Resolver) PatientQueue() generated.PatientQueueResolver { return &patientQueueResolver{r} }

type patientQueueResolver struct{ *Resolver }
