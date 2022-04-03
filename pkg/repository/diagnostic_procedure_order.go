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

// DiagnosticProcedureOrderStatus ...
type DiagnosticProcedureOrderStatus string

// DiagnosticProcedureOrder statuses ...
const (
	DiagnosticProcedureOrderOrderedStatus   DiagnosticProcedureOrderStatus = "ORDERED"
	DiagnosticProcedureOrderCompletedStatus DiagnosticProcedureOrderStatus = "COMPLETED"
)

// DiagnosticProcedureOrder struct holds the order of the diagnostic procedures.
type DiagnosticProcedureOrder struct {
	gorm.Model
	ID                   int                            `gorm:"primaryKey"`
	PatientChartID       int                            `json:"patientChartId"`
	PatientID            int                            `json:"patientId"`
	FirstName            string                         `json:"firstName"`
	LastName             string                         `json:"lastName"`
	PhoneNo              string                         `json:"phoneNo"`
	UserName             string                         `json:"userName"`
	OrderedByID          *int                           `json:"orderedById"`
	OrderedBy            *User                          `json:"orderedBy"`
	Status               DiagnosticProcedureOrderStatus `json:"status"`
	DiagnosticProcedures []DiagnosticProcedure          `json:"diagnosticProcedures"`
	Document             string                         `gorm:"type:tsvector"`
	Count                int64                          `json:"count"`
}

// Save ...
func (r *DiagnosticProcedureOrder) Save(diagnosticProcedureTypeID int, patientChartID int, patientID int, billingID int, user User, orderNote string, receptionNote string) error {
	return DB.Transaction(func(tx *gorm.DB) error {
		// Get Patient
		var patient Patient
		if err := tx.Model(&Patient{}).Where("id = ?", patientID).Take(&patient).Error; err != nil {
			return err
		}

		// Diagnostic Procedure Type
		var diagnosticProcedureType DiagnosticProcedureType
		if err := tx.Model(&DiagnosticProcedureType{}).Where("id = ?", diagnosticProcedureTypeID).Take(&diagnosticProcedureType).Error; err != nil {
			return err
		}

		// Create payment
		var payment Payment
		payment.Status = NotPaidPaymentStatus
		payment.BillingID = billingID
		if err := tx.Create(&payment).Error; err != nil {
			return err
		}

		isPhysician := false
		for _, e := range user.UserTypes {
			if e.Title == "Physician" {
				isPhysician = true
			}
		}

		orderedByPrefix := ""
		if isPhysician {
			orderedByPrefix = "Dr. "
		}

		r.PatientChartID = patientChartID
		r.PatientID = patientID
		r.FirstName = patient.FirstName
		r.LastName = patient.LastName
		r.PhoneNo = patient.PhoneNo
		r.UserName = orderedByPrefix + user.FirstName + " " + user.LastName
		r.OrderedByID = &user.ID
		r.Status = DiagnosticProcedureOrderOrderedStatus

		var existing DiagnosticProcedureOrder
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

		// Create Diagnostic Procedure
		var diagnosticProcedure DiagnosticProcedure
		diagnosticProcedure.DiagnosticProcedureTypeID = diagnosticProcedureType.ID
		diagnosticProcedure.DiagnosticProcedureOrderID = r.ID
		diagnosticProcedure.PatientChartID = patientChartID
		diagnosticProcedure.Payments = append(diagnosticProcedure.Payments, payment)
		diagnosticProcedure.Status = DiagnosticProcedureOrderedStatus
		diagnosticProcedure.DiagnosticProcedureTypeTitle = diagnosticProcedureType.Title
		diagnosticProcedure.OrderNote = orderNote
		diagnosticProcedure.ReceptionNote = receptionNote

		if diagnosticProcedureType.Title == "Refraction" || diagnosticProcedureType.Title == "Refraction Advanced- VIP" {
			diagnosticProcedure.IsRefraction = true
		}

		if err := tx.Create(&diagnosticProcedure).Error; err != nil {
			return err
		}

		return nil
	})
}

// OrderAndConfirm ...
func (r *DiagnosticProcedureOrder) OrderAndConfirm() {

}

// GetTodaysOrderedCount ...
func (r *DiagnosticProcedureOrder) GetTodaysOrderedCount() (count int) {
	now := time.Now()
	start := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	end := start.AddDate(0, 0, 1)

	var countTmp int64
	err := DB.Model(&r).Where("created_at >= ?", start).Where("created_at <= ?", end).Where("status = ?", DiagnosticProcedureOrderOrderedStatus).Count(&countTmp).Error
	if err != nil {
		countTmp = 0
	}

	count = int(countTmp)

	return
}

// GetPatientDiagnosticOrderTitles ...
func (r *DiagnosticProcedureOrder) GetPatientDiagnosticProcedureTitles(patientID int) ([]string, error) {
	var procedureTitles []string

	var result []map[string]interface{}

	if err := DB.Raw(`
		SELECT DISTINCT(diagnostic_procedure_type_title) 
			FROM diagnostic_procedure_orders 
			INNER JOIN diagnostic_procedures
			ON diagnostic_procedure_orders.id = diagnostic_procedures.diagnostic_procedure_order_id
		WHERE patient_id = ?`, patientID).Find(&result).Error; err != nil {
		return procedureTitles, err
	}

	for _, order := range result {
		procedureTitles = append(procedureTitles, order["diagnostic_procedure_type_title"].(string))
	}

	return procedureTitles, nil
}

// Confirm ...
func (r *DiagnosticProcedureOrder) Confirm(id int, invoiceNo string) error {
	return DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("id = ?", id).Preload("DiagnosticProcedures.Payments").Take(&r).Error; err != nil {
			return err
		}

		var payments []Payment
		for _, diagnosticProcedure := range r.DiagnosticProcedures {
			payments = append(payments, diagnosticProcedure.Payments...)
		}

		var paymentIds []int
		for _, payment := range payments {
			paymentIds = append(paymentIds, payment.ID)
		}

		if err := tx.Model(&Payment{}).Where("id IN ?", paymentIds).Updates(map[string]interface{}{"invoice_no": invoiceNo, "status": "PAID"}).Error; err != nil {
			return err
		}

		r.Status = DiagnosticProcedureOrderCompletedStatus

		if err := tx.Updates(&r).Error; err != nil {
			return err
		}

		// Add to Diagnostic Queue
		var patientChart PatientChart

		if err := tx.Where("id = ?", r.PatientChartID).Take(&patientChart).Error; err != nil {
			return err
		}

		for _, diagnosticProcedure := range r.DiagnosticProcedures {
			var patientQueue PatientQueue

			// Create new patient queue if it doesn't exists
			if err := tx.Where("queue_name = ?", diagnosticProcedure.DiagnosticProcedureTypeTitle).Take(&patientQueue).Error; err != nil {
				patientQueue.QueueName = diagnosticProcedure.DiagnosticProcedureTypeTitle
				patientQueue.QueueType = DiagnosticQueue
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

				if err := tx.Where("queue_name = ?", diagnosticProcedure.DiagnosticProcedureTypeTitle).Updates(&PatientQueue{Queue: queue}).Error; err != nil {
					return err
				}
			}
		}

		return nil
	})
}

// GetCount ...
func (r *DiagnosticProcedureOrder) GetCount(filter *DiagnosticProcedureOrder, date *time.Time, searchTerm *string) (int64, error) {
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
func (r *DiagnosticProcedureOrder) Search(p PaginationInput, filter *DiagnosticProcedureOrder, date *time.Time, searchTerm *string, ascending bool) ([]DiagnosticProcedureOrder, int64, error) {
	var result []DiagnosticProcedureOrder

	dbOp := DB.Scopes(Paginate(&p)).Select("*, count(*) OVER() AS count").Where(filter).Preload("DiagnosticProcedures.Payments.Billing").Preload("DiagnosticProcedures.DiagnosticProcedureType").Preload("OrderedBy.UserTypes")

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

	err := dbOp.Find(&result).Error

	var count int64
	if len(result) > 0 {
		count = result[0].Count
	}

	return result, count, err
}

// GetByPatientChartID ...
func (r *DiagnosticProcedureOrder) GetByPatientChartID(patientChartID int) error {
	return DB.Where("patient_chart_id = ?", patientChartID).Preload("DiagnosticProcedures").Preload("DiagnosticProcedures.Payments").Preload("DiagnosticProcedures.DiagnosticProcedureType").Preload("DiagnosticProcedures.Images").Preload("DiagnosticProcedures.Documents").Take(&r).Error
}

// GetAll ...
func (r *DiagnosticProcedureOrder) GetAll(p PaginationInput, filter *DiagnosticProcedureOrder) ([]DiagnosticProcedureOrder, int64, error) {
	var result []DiagnosticProcedureOrder

	dbOp := DB.Scopes(Paginate(&p)).Select("*, count(*) OVER() AS count").Where(filter).Preload("DiagnosticProcedures").Order("id ASC").Find(&result)

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
func (r *DiagnosticProcedureOrder) Update() error {
	return DB.Updates(&r).Error
}

// Delete ...
func (r *DiagnosticProcedureOrder) Delete(ID int) error {
	return DB.Where("id = ?", ID).Delete(&r).Error
}
