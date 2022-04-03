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
	"time"

	"gorm.io/gorm"
)

// SurgicalOrderStatus ...
type SurgicalOrderStatus string

// SurgicalOrderStatus statuses ...
const (
	SurgicalOrderStatusOrdered   SurgicalOrderStatus = "ORDERED"
	SurgicalOrderStatusCompleted SurgicalOrderStatus = "COMPLETED"
)

// SurgicalOrder ...
type SurgicalOrder struct {
	gorm.Model
	ID                 int                 `gorm:"primaryKey"`
	PatientChartID     int                 `json:"patientChartId"`
	PatientID          int                 `json:"patientId"`
	FirstName          string              `json:"firstName"`
	LastName           string              `json:"lastName"`
	PhoneNo            string              `json:"phoneNo"`
	UserName           string              `json:"userName"`
	OrderedByID        int                 `json:"orderedById"`
	OrderedBy          User                `json:"orderedBy"`
	Status             SurgicalOrderStatus `json:"status"`
	SurgicalProcedures []SurgicalProcedure `json:"surgicalProcedures"`
	Emergency          *bool               `json:"emergency"`
	Document           string              `gorm:"type:tsvector"`
	Count              int64               `json:"count"`
}

// SaveOpthalmologyOrder ...
func (r *SurgicalOrder) SaveOpthalmologyOrder(surgicalProcedureTypeID int, patientChartID int, patientID int, billingID int, user User, performOnEye string, orderNote string, receptionNote string) error {
	return DB.Transaction(func(tx *gorm.DB) error {
		// Get Patient
		var patient Patient
		if err := tx.Model(&Patient{}).Where("id = ?", patientID).Take(&patient).Error; err != nil {
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
		r.OrderedByID = user.ID
		r.Status = SurgicalOrderStatusOrdered

		var existing SurgicalOrder
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

		// Surgical Procedure Type
		var surgicalProcedureType SurgicalProcedureType
		if err := tx.Model(&SurgicalProcedureType{}).Preload("Supplies.Billings").Where("id = ?", surgicalProcedureTypeID).Take(&surgicalProcedureType).Error; err != nil {
			return err
		}

		// Create payment
		var payments []Payment

		// Payment for procedure
		var payment Payment
		payment.Status = NotPaidPaymentStatus
		payment.BillingID = billingID
		payments = append(payments, payment)

		// Attach supply payments
		for _, supply := range surgicalProcedureType.Supplies {
			for _, billing := range supply.Billings {
				var payment Payment
				payment.Status = NotPaidPaymentStatus
				payment.BillingID = billing.ID
				payments = append(payments, payment)
			}
		}

		// Create surgical procedure
		var surgicalProcedure SurgicalProcedure
		surgicalProcedure.SurgicalProcedureTypeID = surgicalProcedureType.ID
		surgicalProcedure.SurgicalOrderID = r.ID
		surgicalProcedure.PatientChartID = patientChartID
		surgicalProcedure.Payments = payments
		surgicalProcedure.Status = SurgeryStatusOrdered
		surgicalProcedure.SurgicalProcedureTypeTitle = surgicalProcedureType.Title
		surgicalProcedure.PerformOnEye = performOnEye
		surgicalProcedure.OrderNote = orderNote
		surgicalProcedure.ReceptionNote = receptionNote

		if err := tx.Create(&surgicalProcedure).Error; err != nil {
			return err
		}

		return nil
	})
}

// GetTodaysOrderedCount ...
func (r *SurgicalOrder) GetTodaysOrderedCount() (count int) {
	now := time.Now()
	start := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	end := start.AddDate(0, 0, 1)

	var countTmp int64
	err := DB.Model(&r).Where("created_at >= ?", start).Where("created_at <= ?", end).Where("status = ?", SurgicalOrderStatusOrdered).Count(&countTmp).Error
	if err != nil {
		countTmp = 0
	}

	count = int(countTmp)

	return
}

// ConfirmOrder ...
func (r *SurgicalOrder) ConfirmOrder(surgicalOrderID int, surgicalProcedureID int, invoiceNo string, roomID int, checkInTime time.Time) error {
	return DB.Transaction(func(tx *gorm.DB) error {

		var surgicalProcedure SurgicalProcedure
		if err := tx.Where("id = ?", surgicalProcedureID).Preload("Payments").Take(&surgicalProcedure).Error; err != nil {
			return err
		}

		// Update all surgical procedure payments to paid
		var paymentIds []int
		for _, payment := range surgicalProcedure.Payments {
			paymentIds = append(paymentIds, payment.ID)
		}

		if err := tx.Model(&Payment{}).Where("id IN ?", paymentIds).Updates(map[string]interface{}{"invoice_no": invoiceNo, "status": "PAID"}).Error; err != nil {
			return err
		}

		// Get surgical order with payments
		if err := tx.Where("id = ?", surgicalOrderID).Preload("SurgicalProcedures.Payments").Take(&r).Error; err != nil {
			return err
		}

		var patientChart PatientChart
		if err := tx.Where("id = ?", r.PatientChartID).Take(&patientChart).Error; err != nil {
			return err
		}

		var previousAppointment Appointment
		if err := tx.Where("id = ?", patientChart.AppointmentID).Take(&previousAppointment).Error; err != nil {
			return err
		}

		var allPayments []Payment
		for _, surgicalProcedure := range r.SurgicalProcedures {
			allPayments = append(allPayments, surgicalProcedure.Payments...)
		}

		allPaid := true
		for _, payment := range allPayments {
			if payment.Status == NotPaidPaymentStatus {
				allPaid = false
			}
		}

		if allPaid {
			// Update surgical order to completed
			r.Status = SurgicalOrderStatusCompleted
			if err := tx.Updates(&r).Error; err != nil {
				return err
			}
		}

		// Create new appointment
		var appointment Appointment
		appointment.PatientID = r.PatientID
		appointment.RoomID = roomID
		appointment.CheckInTime = checkInTime
		appointment.UserID = r.OrderedByID
		appointment.Credit = false
		appointment.Payments = surgicalProcedure.Payments
		appointment.MedicalDepartment = previousAppointment.MedicalDepartment

		// Assign surgery visit type
		var visitType VisitType
		if err := tx.Where("title = ?", "Surgery").Take(&visitType).Error; err != nil {
			return err
		}
		appointment.VisitTypeID = visitType.ID

		// Assign scheduled status
		var status AppointmentStatus
		if err := tx.Where("title = ?", "Scheduled").Take(&status).Error; err != nil {
			return err
		}
		appointment.AppointmentStatusID = status.ID

		// Create appointment
		if err := tx.Create(&appointment).Error; err != nil {
			return err
		}

		// Create new patient chart
		var newPatientChart PatientChart
		newPatientChart.AppointmentID = appointment.ID
		if err := tx.Create(&newPatientChart).Error; err != nil {
			return err
		}

		surgicalProcedure.Status = SurgeryStatusOrdered
		surgicalProcedure.PatientChartID = newPatientChart.ID
		if err := tx.Updates(&surgicalProcedure).Error; err != nil {
			return err
		}

		return nil
	})
}

// GetCount ...
func (r *SurgicalOrder) GetCount(filter *SurgicalOrder, date *time.Time, searchTerm *string) (int64, error) {
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
func (r *SurgicalOrder) Search(p PaginationInput, filter *SurgicalOrder, date *time.Time, searchTerm *string, ascending bool) ([]SurgicalOrder, int64, error) {
	var result []SurgicalOrder

	dbOp := DB.Scopes(Paginate(&p)).Select("*, count(*) OVER() AS count").Where(filter).Preload("SurgicalProcedures.Payments.Billing").Preload("SurgicalProcedures.SurgicalProcedureType").Preload("OrderedBy.UserTypes")

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
func (r *SurgicalOrder) GetByPatientChartID(patientChartID int) error {
	return DB.Where("patient_chart_id = ?", patientChartID).Preload("SurgicalProcedures").Preload("SurgicalProcedures.Payments").Preload("SurgicalProcedures.SurgicalProcedureType").Take(&r).Error
}

// GetAll ...
func (r *SurgicalOrder) GetAll(p PaginationInput, filter *SurgicalOrder) ([]SurgicalOrder, int64, error) {
	var result []SurgicalOrder

	dbOp := DB.Scopes(Paginate(&p)).Select("*, count(*) OVER() AS count").Where(filter).Preload("SurgicalProcedures").Order("id ASC").Find(&result)

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
func (r *SurgicalOrder) Update() error {
	return DB.Updates(&r).Error
}

// Delete ...
func (r *SurgicalOrder) Delete(ID int) error {
	return DB.Where("id = ?", ID).Delete(&r).Error
}
