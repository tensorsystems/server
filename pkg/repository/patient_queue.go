/*
  Copyright 2021 Kidus Tiliksew

  This file is part of Tensor EMR.

  Tensor EMR is free software: you can redistribute it and/or modify
  it under the terms of the version 2 of GNU General Public License as published by
  the Free Software Foundation.

  Tensor EMR is distributed in the hope that it will be useful,
  but WITHOUT ANY WARRANTY; without even the implied warranty of
  MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
  GNU General Public License for more details.

  You should have received a copy of the GNU General Public License
  along with this program.  If not, see <http://www.gnu.org/licenses/>.
*/

package repository

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

// QueueType ...
type QueueType string

// Queue Types ...
const (
	UserQueue       QueueType = "USER"
	DiagnosticQueue QueueType = "DIAGNOSTIC"
	LabQueue        QueueType = "LAB"
	TreatmentQueue  QueueType = "TREATMENT"
	SurgicalQueue   QueueType = "SURGICAL"
	PreExamQueue    QueueType = "PREEXAM"
	PreOperation    QueueType = "PREOPERATION"
)

// PatientQueue ...
type PatientQueue struct {
	gorm.Model
	ID        int            `gorm:"primaryKey" json:"id"`
	QueueName string         `json:"queueName"`
	Queue     datatypes.JSON `json:"queue"`
	QueueType QueueType      `json:"queueType"`
}

// Save
func (r *PatientQueue) Save() error {
	return DB.Create(&r).Error
}

// GetAll
func (r *PatientQueue) GetAll() ([]*PatientQueue, error) {
	var result []*PatientQueue
	err := DB.Order("queue_type DESC").Find(&result).Error

	return result, err
}

// Get ...
func (r *PatientQueue) Get(id int) error {
	return DB.Where("id = ?", id).Take(&r).Error
}

// GetByQueueName ...
func (r *PatientQueue) GetByQueueName(queueName string) error {
	return DB.Where("queue_name = ?", queueName).Take(&r).Error
}

// GetByQueueName ...
func (r *PatientQueue) UpdateQueue(queueName string, queue datatypes.JSON) error {
	return DB.Where("queue_name = ?", queueName).Updates(&PatientQueue{Queue: queue}).Error
}

// Move ...
func (r *PatientQueue) Move(fromQueueID int, toQueueID int, appointmentID int) error {
	return DB.Transaction(func(tx *gorm.DB) error {
		// Get source queue
		var sourceQueue PatientQueue
		if err := tx.Where("id = ?", fromQueueID).Take(&sourceQueue).Error; err != nil {
			return err
		}

		// Get destination queue
		var destinationQueue PatientQueue
		if err := tx.Where("id = ?", toQueueID).Take(&destinationQueue).Error; err != nil {
			return err
		}

		// Remove appointment id from source queue
		var sourceIds []int
		if err := json.Unmarshal([]byte(sourceQueue.Queue.String()), &sourceIds); err != nil {
			return err
		}

		for i, v := range sourceIds {
			if v == appointmentID {
				sourceIds = remove(sourceIds, i)
			}
		}

		var sourceAppointmentIds []string
		for _, v := range sourceIds {
			sourceAppointmentIds = append(sourceAppointmentIds, fmt.Sprint(v))
		}

		sourceQueue.Queue = datatypes.JSON([]byte("[" + strings.Join(sourceAppointmentIds, ", ") + "]"))
		if err := tx.Updates(&PatientQueue{ID: sourceQueue.ID, Queue: sourceQueue.Queue}).Error; err != nil {
			return err
		}

		var destinationIds []int
		if err := json.Unmarshal([]byte(destinationQueue.Queue.String()), &destinationIds); err != nil {
			return err
		}

		// Skip adding if it alread exists
		exists := false
		for _, e := range destinationIds {
			if e == appointmentID {
				exists = true
			}
		}

		if exists {
			return nil
		}

		var destinationAppointmentIds []string
		for _, v := range destinationIds {
			destinationAppointmentIds = append(destinationAppointmentIds, fmt.Sprint(v))
		}

		destinationAppointmentIds = append(destinationAppointmentIds, fmt.Sprint(appointmentID))
		destinationQueue.Queue = datatypes.JSON([]byte("[" + strings.Join(destinationAppointmentIds, ", ") + "]"))
		if err := tx.Updates(&PatientQueue{ID: destinationQueue.ID, Queue: destinationQueue.Queue}).Error; err != nil {
			return err
		}

		return nil
	})
}

// MoveToQueueName ...
func (r *PatientQueue) MoveToQueueName(fromQueueID int, toQueueName string, appointmentID int, queueType string) error {
	return DB.Transaction(func(tx *gorm.DB) error {
		// Get source queue
		var sourceQueue PatientQueue
		if err := tx.Where("id = ?", fromQueueID).Take(&sourceQueue).Error; err != nil {
			return err
		}

		// Remove appointment id from source queue
		var sourceIds []int
		if err := json.Unmarshal([]byte(sourceQueue.Queue.String()), &sourceIds); err != nil {
			return err
		}

		for i, v := range sourceIds {
			if v == appointmentID {
				sourceIds = remove(sourceIds, i)
			}
		}

		var sourceAppointmentIds []string
		for _, v := range sourceIds {
			sourceAppointmentIds = append(sourceAppointmentIds, fmt.Sprint(v))
		}

		sourceQueue.Queue = datatypes.JSON([]byte("[" + strings.Join(sourceAppointmentIds, ", ") + "]"))
		if err := tx.Updates(&PatientQueue{ID: sourceQueue.ID, Queue: sourceQueue.Queue}).Error; err != nil {
			return err
		}

		// Add appointment to destination queue
		var destinationQueue PatientQueue

		// Check if queue exists, create new one if it doesn't
		if err := tx.Where("queue_name = ?", toQueueName).Take(&destinationQueue).Error; err != nil {
			destinationQueue.QueueName = toQueueName
			destinationQueue.Queue = datatypes.JSON([]byte("[" + fmt.Sprint(appointmentID) + "]"))

			if len(queueType) != 0 {
				destinationQueue.QueueType = QueueType(queueType)
			}

			if err := tx.Create(&destinationQueue).Error; err != nil {
				return err
			}

			return nil
		}

		var destinationIds []int
		if err := json.Unmarshal([]byte(destinationQueue.Queue.String()), &destinationIds); err != nil {
			return err
		}

		// Skip adding if it alread exists
		exists := false
		for _, e := range destinationIds {
			if e == appointmentID {
				exists = true
			}
		}

		if exists {
			return nil
		}

		var destinationAppointmentIds []string
		for _, v := range destinationIds {
			destinationAppointmentIds = append(destinationAppointmentIds, fmt.Sprint(v))
		}

		destinationAppointmentIds = append(destinationAppointmentIds, fmt.Sprint(appointmentID))
		destinationQueue.Queue = datatypes.JSON([]byte("[" + strings.Join(destinationAppointmentIds, ", ") + "]"))
		if err := tx.Updates(&PatientQueue{ID: destinationQueue.ID, Queue: destinationQueue.Queue}).Error; err != nil {
			return err
		}

		return nil
	})
}

// AddToQueue
func (r *PatientQueue) AddToQueue(toQueueName string, appointmentID int, queueType string) error {
	return DB.Transaction(func(tx *gorm.DB) error {

		// Check if queue exists, create new one if it doesn't
		if err := tx.Where("queue_name = ?", toQueueName).Take(&r).Error; err != nil {
			r.QueueName = toQueueName
			r.Queue = datatypes.JSON([]byte("[" + fmt.Sprint(appointmentID) + "]"))
			if len(queueType) != 0 {
				r.QueueType = QueueType(queueType)
			}

			if err := tx.Create(&r).Error; err != nil {
				return err
			}

			return nil
		}

		var destinationIds []int
		if err := json.Unmarshal([]byte(r.Queue.String()), &destinationIds); err != nil {
			return err
		}

		// Skip adding if it alread exists
		exists := false
		for _, e := range destinationIds {
			if e == appointmentID {
				exists = true
			}
		}

		if exists {
			return nil
		}

		var destinationAppointmentIds []string
		for _, v := range destinationIds {
			destinationAppointmentIds = append(destinationAppointmentIds, fmt.Sprint(v))
		}

		destinationAppointmentIds = append(destinationAppointmentIds, fmt.Sprint(appointmentID))
		r.Queue = datatypes.JSON([]byte("[" + strings.Join(destinationAppointmentIds, ", ") + "]"))
		if err := tx.Updates(&PatientQueue{ID: r.ID, Queue: r.Queue}).Error; err != nil {
			return err
		}

		return nil
	})
}

// DeleteFromQueue ...
func (r *PatientQueue) DeleteFromQueue(patientQueueID int, appointmentID int) error {
	return DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("id = ?", patientQueueID).Take(&r).Error; err != nil {
			return err
		}

		var ids []int
		if err := json.Unmarshal([]byte(r.Queue.String()), &ids); err != nil {
			return err
		}

		for i, v := range ids {
			if v == appointmentID {
				ids = remove(ids, i)
			}
		}

		var appointmentIds []string
		for _, v := range ids {
			appointmentIds = append(appointmentIds, fmt.Sprint(v))
		}

		r.Queue = datatypes.JSON([]byte("[" + strings.Join(appointmentIds, ", ") + "]"))
		if err := tx.Updates(&PatientQueue{ID: r.ID, Queue: r.Queue}).Error; err != nil {
			return err
		}

		return nil
	})
}

// ClearExpired ...
func (r *PatientQueue) ClearExpired() error {
	return DB.Transaction(func(tx *gorm.DB) error {
		var patientQueues []PatientQueue

		if err := tx.Find(&patientQueues).Error; err != nil {
			return err
		}

		for _, patientQueue := range patientQueues {
			var ids []int
			if err := json.Unmarshal([]byte(patientQueue.Queue.String()), &ids); err != nil {
				return err
			}

			now := time.Now()
			start := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())

			var expiredAppointments []Appointment

			if err := tx.Model(Appointment{}).Select("id").Where("id IN ?", ids).Where("checked_in_time < ?", start).Find(&expiredAppointments).Error; err != nil {
				return err
			}

			if len(expiredAppointments) != 0 {
				for index, id := range ids {
					for _, appointment := range expiredAppointments {
						if appointment.ID == id {
							ids = remove(ids, index)
						}
					}
				}
			}

			var appointmentIds []string
			for _, v := range ids {
				appointmentIds = append(appointmentIds, fmt.Sprint(v))
			}

			queue := datatypes.JSON([]byte("[" + strings.Join(appointmentIds, ", ") + "]"))
			patientQueue.Queue = queue

			if err := tx.Updates(&patientQueue).Error; err != nil {
				return err
			}
		}

		return nil
	})
}

func remove(slice []int, s int) []int {
	if len(slice) > 0 && s >= len(slice) {
		return slice[:len(slice)-1]
	}

	return append(slice[:s], slice[s+1:]...)
}
