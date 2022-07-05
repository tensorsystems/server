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

	"github.com/tensoremr/server/pkg/graphql/graph/generated"
	graph_models "github.com/tensoremr/server/pkg/graphql/graph/model"
	"github.com/tensoremr/server/pkg/middleware"
	"github.com/tensoremr/server/pkg/models"
	"github.com/tensoremr/server/pkg/repository"
	"gorm.io/datatypes"
)

func (r *mutationResolver) SubscribeQueue(ctx context.Context, userID int, patientQueueID int) (*models.QueueSubscription, error) {
	var entity models.QueueSubscription

	var repository repository.QueueSubscriptionRepository
	if err := repository.Subscribe(&entity, userID, patientQueueID); err != nil {
		return nil, err
	}

	return &entity, nil
}

func (r *mutationResolver) UnsubscribeQueue(ctx context.Context, userID int, patientQueueID int) (*models.QueueSubscription, error) {
	var entity models.QueueSubscription

	var repository repository.QueueSubscriptionRepository
	if err := repository.Unsubscribe(&entity, userID, patientQueueID); err != nil {
		return nil, err
	}

	return &entity, nil
}

func (r *mutationResolver) SavePatientQueue(ctx context.Context, input graph_models.PatientQueueInput) (*models.PatientQueue, error) {
	var entity models.PatientQueue
	entity.QueueName = input.QueueName
	entity.QueueType = input.QueueType

	queue := datatypes.JSON([]byte("[" + strings.Join(input.Queue, ", ") + "]"))
	entity.Queue = queue

	var repository repository.PatientQueueRepository

	if err := repository.GetByQueueName(&entity, input.QueueName); err != nil {
		if err := repository.Save(&entity); err != nil {
			return nil, err
		}
	} else {
		if err := repository.UpdateQueue(input.QueueName, queue); err != nil {
			return nil, err
		}
	}

	return &entity, nil
}

func (r *mutationResolver) DeleteFromQueue(ctx context.Context, patientQueueID int, appointmentID int) (*models.PatientQueue, error) {
	var entity models.PatientQueue
	var repository repository.PatientQueueRepository

	if err := repository.DeleteFromQueue(&entity, patientQueueID, appointmentID); err != nil {
		return nil, err
	}

	return &entity, nil
}

func (r *mutationResolver) CheckOutPatient(ctx context.Context, patientQueueID int, appointmentID int) (*models.PatientQueue, error) {
	var appointment models.Appointment
	var appointmentRepository repository.AppointmentRepository

	if err := appointmentRepository.Get(&appointment, appointmentID); err != nil {
		return nil, err
	}

	// Change appointment status to Checked Out
	var status models.AppointmentStatus
	var appointmentStatusRepository repository.AppointmentStatusRepository
	if err := appointmentStatusRepository.GetByTitle(&status, "Checked-Out"); err != nil {
		return nil, err
	}

	appointment.AppointmentStatusID = status.ID
	appointment.CheckedOutTime = time.Now()

	if err := appointmentRepository.Update(&appointment); err != nil {
		return nil, err
	}

	var patientQueue models.PatientQueue
	var patientQueueRepository repository.PatientQueueRepository

	if err := patientQueueRepository.DeleteFromQueue(&patientQueue, patientQueueID, appointmentID); err != nil {
		return nil, err
	}

	return &patientQueue, nil
}

func (r *mutationResolver) PushPatientQueue(ctx context.Context, patientQueueID int, appointmentID int, destination graph_models.Destination) (*models.PatientQueue, error) {
	var entity models.PatientQueue
	var patientQueueRepository repository.PatientQueueRepository
	var appointmentRepository repository.AppointmentRepository

	if destination.String() == "PREEXAM" {
		if err := patientQueueRepository.MoveToQueueName(patientQueueID, "Pre-Exam", appointmentID, ""); err != nil {
			return nil, err
		}
	} else if destination.String() == "PREOPERATION" {
		if err := patientQueueRepository.MoveToQueueName(patientQueueID, "Pre-Operation", appointmentID, ""); err != nil {
			return nil, err
		}
	} else if destination.String() == "PHYSICIAN" {
		var appointment models.Appointment
		if err := appointmentRepository.Get(&appointment, appointmentID); err != nil {
			return nil, err
		}

		var userRepository repository.UserRepository
		var provider models.User
		if err := userRepository.Get(&provider, appointment.UserID); err != nil {
			return nil, err
		}

		if err := patientQueueRepository.MoveToQueueName(patientQueueID, "Dr. "+provider.FirstName+" "+provider.LastName, appointmentID, "USER"); err != nil {
			return nil, err
		}
	}

	return &entity, nil
}

func (r *mutationResolver) MovePatientQueue(ctx context.Context, appointmentID int, sourceQueueID int, destinationQueueID int) (*models.PatientQueue, error) {
	var entity models.PatientQueue
	var repository repository.PatientQueueRepository

	if err := repository.Move(&entity, sourceQueueID, destinationQueueID, appointmentID); err != nil {
		return nil, err
	}

	return &entity, nil
}

func (r *mutationResolver) CheckInPatient(ctx context.Context, appointmentID int, destination graph_models.Destination) (*models.Appointment, error) {
	var appointment models.Appointment
	var userRepository repository.UserRepository
	var appointmentRepository repository.AppointmentRepository
	var patientQueueRepository repository.PatientQueueRepository

	if err := appointmentRepository.Get(&appointment, appointmentID); err != nil {
		return nil, err
	}

	// Change appointment status to Checked In
	var status models.AppointmentStatus
	var appointmentStatusRepository repository.AppointmentStatusRepository
	if err := appointmentStatusRepository.GetByTitle(&status, "Checked-In"); err != nil {
		return nil, err
	}

	var visitType models.VisitType
	var visitTypeRepository repository.VisitTypeRepository
	if err := visitTypeRepository.Get(&visitType, appointment.VisitTypeID); err != nil {
		return nil, err
	}

	if visitType.Title == "Surgery" {
		var postOpAppointent models.Appointment

		if err := appointmentRepository.SchedulePostOp(&postOpAppointent, appointment); err != nil {
			return nil, err
		}
	}

	appointment.AppointmentStatusID = status.ID
	checkedInTime := time.Now()
	appointment.CheckedInTime = &checkedInTime

	if err := appointmentRepository.Update(&appointment); err != nil {
		return nil, err
	}

	// Add to queue
	var patientQueue models.PatientQueue
	if destination.String() == "PREEXAM" {
		if err := patientQueueRepository.AddToQueue(&patientQueue, "Pre-Exam", appointmentID, "PREEXAM"); err != nil {
			return nil, err
		}

	} else if destination.String() == "PREOPERATION" {
		if err := patientQueueRepository.AddToQueue(&patientQueue, "Pre-Operation", appointmentID, "PREOPERATION"); err != nil {
			return nil, err
		}
	} else if destination.String() == "PHYSICIAN" {
		var provider models.User
		if err := userRepository.Get(&provider, appointment.UserID); err != nil {
			return nil, err
		}

		if err := patientQueueRepository.AddToQueue(&patientQueue, "Dr. "+provider.FirstName+" "+provider.LastName, appointmentID, "USER"); err != nil {
			return nil, err
		}
	}

	return &appointment, nil
}

func (r *mutationResolver) UpdatePatientQueue(ctx context.Context, appointmentID int, destination *graph_models.Destination) (*models.PatientQueue, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *patientQueueResolver) Queue(ctx context.Context, obj *models.PatientQueue) (string, error) {
	return obj.Queue.String(), nil
}

func (r *queryResolver) PatientQueues(ctx context.Context) ([]*graph_models.PatientQueueWithAppointment, error) {
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

	isPhysician := false
	for _, e := range user.UserTypes {
		if e.Title == "Physician" {
			isPhysician = true
		}
	}

	var patientQueueRepository repository.PatientQueueRepository
	var patientQueues []*models.PatientQueue

	if !isPhysician {
		p, err := patientQueueRepository.GetAll()
		if err != nil {
			return nil, err
		}

		patientQueues = p
	} else {
		p, err := patientQueueRepository.GetAll()
		if err != nil {
			return nil, err
		}

		var t []*models.PatientQueue

		for _, patientQueue := range p {
			if patientQueue.QueueType == models.UserQueue && patientQueue.QueueName != "Dr. "+user.FirstName+" "+user.LastName {
				continue
			}

			t = append(t, patientQueue)
		}

		patientQueues = t
	}

	var result []*graph_models.PatientQueueWithAppointment

	var appointmentRepo repository.AppointmentRepository
	for _, patientQueue := range patientQueues {
		var ids []int

		if err := json.Unmarshal([]byte(patientQueue.Queue.String()), &ids); err != nil {
			return nil, err
		}

		page := repository.PaginationInput{Page: 0, Size: 1000}

		appointments, _, _ := appointmentRepo.GetByIds(ids, page)
		var orderedAppointments []*models.Appointment

		for _, id := range ids {
			for _, appointment := range appointments {
				if appointment.ID == id {
					a := appointment
					orderedAppointments = append(orderedAppointments, &a)
				}
			}
		}

		result = append(result, &graph_models.PatientQueueWithAppointment{
			ID:        int(patientQueue.ID),
			QueueName: patientQueue.QueueName,
			QueueType: patientQueue.QueueType,
			Queue:     orderedAppointments,
		})
	}

	return result, err
}

func (r *queryResolver) UserSubscriptions(ctx context.Context, userID int) (*models.QueueSubscription, error) {
	var entity models.QueueSubscription
	var repository repository.QueueSubscriptionRepository

	if err := repository.GetByUserId(&entity, userID); err != nil {
		return nil, err
	}

	return &entity, nil
}

// PatientQueue returns generated.PatientQueueResolver implementation.
func (r *Resolver) PatientQueue() generated.PatientQueueResolver { return &patientQueueResolver{r} }

type patientQueueResolver struct{ *Resolver }
