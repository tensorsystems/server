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

// LabOrderStatus
type LabOrderStatus string

// LabOrder statuses
const (
	LabOrderOrderedStatus   LabOrderStatus = "ORDERED"
	LabOrderCompletedStatus LabOrderStatus = "COMPLETED"
)

// LabOrder ...
type LabOrder struct {
	gorm.Model
	ID             int            `gorm:"primaryKey"`
	PatientChartID int            `json:"patientChartId"`
	PatientID      int            `json:"patientId"`
	FirstName      string         `json:"firstName"`
	LastName       string         `json:"lastName"`
	PhoneNo        string         `json:"phoneNo"`
	UserName       string         `json:"userName"`
	OrderedByID    *int           `json:"orderedById"`
	OrderedBy      *User          `json:"orderedBy"`
	Status         LabOrderStatus `json:"status"`
	Labs           []Lab          `json:"labs"`
	Document       string         `gorm:"type:tsvector"`
	Count          int64          `json:"count"`
}

// NewOrder ...
func (r *LabOrder) Save(labTypeID int, patientChartID int, patientID int, billingIds []int, user User, orderNote string, receptionNote string) error {
	return DB.Transaction(func(tx *gorm.DB) error {
		// Get Patient
		var patient Patient
		if err := tx.Model(&Patient{}).Where("id = ?", patientID).Take(&patient).Error; err != nil {
			return err
		}

		// Lab Type
		var labType LabType
		if err := tx.Model(&LabType{}).Where("id = ?", labTypeID).Take(&labType).Error; err != nil {
			return err
		}

		// Create payments
		var payments []Payment
		for _, billingId := range billingIds {
			var payment Payment
			payment.Status = NotPaidPaymentStatus
			payment.BillingID = billingId
			payments = append(payments, payment)
		}

		r.PatientChartID = patientChartID
		r.PatientID = patientID
		r.FirstName = patient.FirstName
		r.LastName = patient.LastName
		r.PhoneNo = patient.PhoneNo
		r.UserName = user.FirstName + " " + user.LastName
		r.OrderedByID = &user.ID
		r.Status = LabOrderOrderedStatus

		var existing LabOrder
		existingErr := tx.Where("patient_chart_id = ?", r.PatientChartID).Take(&existing).Error

		if existingErr != nil {
			if err := tx.Create(&r).Error; err != nil {
				return err
			}
		} else {
			r.ID = existing.ID
			if err := tx.Updates(&r).Error; err != nil {
				return err
			}
		}

		// Create Lab
		var lab Lab
		lab.LabTypeID = labType.ID
		lab.LabOrderID = r.ID
		lab.PatientChartID = patientChartID
		lab.Payments = payments
		lab.Status = LabOrderedStatus
		lab.LabTypeTitle = labType.Title
		lab.OrderNote = orderNote
		lab.ReceptionNote = receptionNote

		if err := tx.Create(&lab).Error; err != nil {
			return err
		}

		return nil
	})
}

// GetTodaysOrderedCount ...
func (r *LabOrder) GetTodaysOrderedCount() (count int) {
	now := time.Now()
	start := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	end := start.AddDate(0, 0, 1)

	var countTmp int64
	err := DB.Model(&r).Where("created_at >= ?", start).Where("created_at <= ?", end).Where("status = ?", LabOrderOrderedStatus).Count(&countTmp).Error
	if err != nil {
		countTmp = 0
	}

	count = int(countTmp)

	return
}

// Confirm ...
func (r *LabOrder) Confirm(id int, invoiceNo string) error {
	return DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("id = ?", id).Preload("Labs.Payments").Take(&r).Error; err != nil {
			return err
		}

		var payments []Payment
		for _, lab := range r.Labs {
			payments = append(payments, lab.Payments...)
		}

		var paymentIds []int
		for _, payment := range payments {
			paymentIds = append(paymentIds, payment.ID)
		}

		if err := tx.Model(&Payment{}).Where("id IN ?", paymentIds).Updates(map[string]interface{}{"invoice_no": invoiceNo, "status": "PAID"}).Error; err != nil {
			return err
		}

		r.Status = LabOrderCompletedStatus

		if err := tx.Updates(&r).Error; err != nil {
			return err
		}

		// Add to Lab Queue
		var patientChart PatientChart
		if err := tx.Where("id = ?", r.PatientChartID).Take(&patientChart).Error; err != nil {
			return err
		}

		for _, lab := range r.Labs {
			var patientQueue PatientQueue

			// Create new patient queue if it doesn't exists
			if err := tx.Where("queue_name = ?", lab.LabTypeTitle).Take(&patientQueue).Error; err != nil {
				patientQueue.QueueName = lab.LabTypeTitle
				patientQueue.QueueType = LabQueue
				patientQueue.Queue = datatypes.JSON([]byte("[" + fmt.Sprint(patientChart.AppointmentID) + "]"))

				if err := tx.Create(&patientQueue).Error; err != nil {
					return err
				}
			} else {
				var ids []int
				if err := json.Unmarshal([]byte(patientQueue.Queue.String()), &ids); err != nil {
					return err
				}

				exists := false
				for _, e := range ids {
					if e == patientChart.AppointmentID {
						exists = true
					}
				}

				if exists {
					return nil
				}

				ids = append(ids, patientChart.AppointmentID)

				var appointmentIds []string
				for _, v := range ids {
					appointmentIds = append(appointmentIds, fmt.Sprint(v))
				}

				queue := datatypes.JSON([]byte("[" + strings.Join(appointmentIds, ", ") + "]"))

				if err := tx.Where("queue_name = ?", lab.LabTypeTitle).Updates(&PatientQueue{Queue: queue}).Error; err != nil {
					return err
				}
			}
		}

		return nil
	})
}

// GetCount ...
func (r *LabOrder) GetCount(filter *LabOrder, date *time.Time, searchTerm *string) (int64, error) {
	dbOp := DB.Model(&r).Where(filter)

	if date != nil {
		createdAt := *date
		start := time.Date(createdAt.Year(), createdAt.Month(), createdAt.Day(), 0, 0, 0, 0, createdAt.Location())
		end := start.AddDate(0, 0, 1)
		dbOp.Where("created_at >= ?", start).Where("created_at <= ?", end)
	}

	if searchTerm != nil {
		dbOp.Where("document @@ plainto_tsquery(?)", *searchTerm)
	}

	var count int64
	err := dbOp.Count(&count).Error

	return count, err
}

// Search ...
func (r *LabOrder) Search(p PaginationInput, filter *LabOrder, date *time.Time, searchTerm *string, ascending bool) ([]LabOrder, int64, error) {
	var result []LabOrder

	dbOp := DB.Scopes(Paginate(&p)).Select("*, count(*) OVER() AS count").Where(filter).Preload("Labs.Payments.Billing").Preload("Labs.LabType").Preload("OrderedBy.UserTypes")

	if date != nil {
		createdAt := *date
		start := time.Date(createdAt.Year(), createdAt.Month(), createdAt.Day(), 0, 0, 0, 0, createdAt.Location())
		end := start.AddDate(0, 0, 1)
		dbOp.Where("created_at >= ?", start).Where("created_at <= ?", end)
	}

	if searchTerm != nil {
		dbOp.Where("document @@ plainto_tsquery(?)", *searchTerm)
	}

	if ascending {
		dbOp.Order("id ASC")
	} else {
		dbOp.Order("id DESC")
	}

	dbOp.Find(&result)

	var count int64
	if len(result) > 0 {
		count = result[0].Count
	}

	if dbOp.Error != nil {
		return result, 0, dbOp.Error
	}

	return result, count, dbOp.Error
}

// Get ...
func (r *LabOrder) Get(ID int) error {
	return DB.Where("id = ?", ID).Take(&r).Error
}

// GetByPatientChartID ...
func (r *LabOrder) GetByPatientChartID(patientChartID int) error {
	return DB.Where("patient_chart_id = ?", patientChartID).Preload("Labs.Payments").Preload("Labs.LabType").Preload("Labs.RightEyeImages").Preload("Labs.LeftEyeImages").Preload("Labs.RightEyeSketches").Preload("Labs.LeftEyeSketches").Preload("Labs.Images").Preload("Labs.Documents").Take(&r).Error
}

// GetAll ...
func (r *LabOrder) GetAll(p PaginationInput, filter *LabOrder) ([]LabOrder, int64, error) {
	var result []LabOrder

	dbOp := DB.Scopes(Paginate(&p)).Select("*, count(*) OVER() AS count").Where(filter).Preload("Labs").Order("id ASC").Find(&result)

	var count int64
	if len(result) > 0 {
		count = result[0].Count
	}

	if dbOp.Error != nil {
		return result, 0, dbOp.Error
	}

	return result, count, dbOp.Error
}

// Update ...
func (r *LabOrder) Update() error {
	return DB.Updates(&r).Error
}

// Delete ...
func (r *LabOrder) Delete(ID int) error {
	return DB.Where("id = ?", ID).Delete(&r).Error
}
