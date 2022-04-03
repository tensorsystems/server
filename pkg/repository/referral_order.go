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

// ReferralOrderStatus ...
type ReferralOrderStatus string

// ReferralOrderStatus statuses ...
const (
	ReferralOrderStatusOrdered   ReferralOrderStatus = "ORDERED"
	ReferralOrderStatusCompleted ReferralOrderStatus = "COMPLETED"
)

// ReferralOrder ...
type ReferralOrder struct {
	gorm.Model
	ID             int                 `gorm:"primaryKey"`
	PatientChartID int                 `json:"patientChartId"`
	PatientID      int                 `json:"patientId"`
	FirstName      string              `json:"firstName"`
	LastName       string              `json:"lastName"`
	PhoneNo        string              `json:"phoneNo"`
	UserName       string              `json:"userName"`
	OrderedByID    int                 `json:"orderedById"`
	OrderedBy      User                `json:"orderedBy"`
	Status         ReferralOrderStatus `json:"status"`
	Referrals      []Referral          `json:"referrals"`
	Emergency      *bool               `json:"emergency"`
	Document       string              `gorm:"type:tsvector"`
	Count          int64               `json:"count"`
}

// Save ...
func (r *ReferralOrder) Save(patientChartID int, patientID int, orderedToID *int, referralType ReferralType, user User, receptionNote *string, reason string, providerName *string) error {
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
		r.Status = ReferralOrderStatusOrdered

		var existing ReferralOrder
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

		var referral Referral
		referral.ReferralOrderID = r.ID
		referral.PatientChartID = patientChartID
		referral.Reason = reason
		referral.Type = referralType

		if referral.Type == ReferralTypeInHouse {
			referral.Status = ReferralStatusOrdered
		} else if referral.Type == ReferralTypeOutsource {
			referral.Status = ReferralStatusCompleted
		}

		if receptionNote != nil {
			referral.ReceptionNote = *receptionNote
		}

		if orderedToID != nil {
			var referredTo User
			if err := tx.Model(&User{}).Where("id = ?", *orderedToID).Take(&referredTo).Error; err != nil {
				return err
			}

			referral.ReferredToID = &referredTo.ID
			referral.ReferredToName = "Dr. " + referredTo.FirstName + " " + referredTo.LastName
		}

		if err := tx.Create(&referral).Error; err != nil {
			return err
		}

		return nil
	})
}

// GetTodaysOrderedCount ...
func (r *ReferralOrder) GetTodaysOrderedCount() (count int) {
	now := time.Now()
	start := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	end := start.AddDate(0, 0, 1)

	var countTmp int64
	err := DB.Model(&r).Where("created_at >= ?", start).Where("created_at <= ?", end).Where("status = ?", ReferralOrderStatusOrdered).Count(&countTmp).Error
	if err != nil {
		countTmp = 0
	}

	count = int(countTmp)

	return
}

// ConfirmOrder ...
func (r *ReferralOrder) ConfirmOrder(referralOrderID int, referralID int, billingID *int, invoiceNo *string, roomID *int, checkInTime *time.Time) error {
	return DB.Transaction(func(tx *gorm.DB) error {
		var referral Referral
		if err := tx.Where("id = ?", referralID).Take(&referral).Error; err != nil {
			return err
		}

		if err := tx.Where("id = ?", referralOrderID).Take(&r).Error; err != nil {
			return err
		}

		if referral.Type == ReferralTypeInHouse {
			var patientChart PatientChart
			if err := tx.Where("id = ?", r.PatientChartID).Take(&patientChart).Error; err != nil {
				return err
			}

			var previousAppointment Appointment
			if err := tx.Where("id = ?", patientChart.AppointmentID).Take(&previousAppointment).Error; err != nil {
				return err
			}

			// Create new appointment
			var appointment Appointment
			appointment.PatientID = r.PatientID
			appointment.RoomID = *roomID
			appointment.CheckInTime = *checkInTime
			appointment.UserID = *referral.ReferredToID
			appointment.Credit = false
			appointment.MedicalDepartment = previousAppointment.MedicalDepartment

			if billingID != nil {
				var payment Payment

				payment.Status = PaidPaymentStatus
				payment.BillingID = *billingID

				if invoiceNo != nil {
					payment.InvoiceNo = *invoiceNo
				}

				appointment.Payments = append(appointment.Payments, payment)
			}

			// Assign treatment visit type
			var visitType VisitType
			if err := tx.Where("title = ?", "Referral").Take(&visitType).Error; err != nil {
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
		}

		referral.Status = ReferralStatusCompleted
		if err := tx.Updates(&referral).Error; err != nil {
			return err
		}

		var referrals []Referral
		if err := tx.Where("referral_order_id = ?", r.ID).Find(&referrals).Error; err != nil {
			return err
		}

		allConfirmed := true
		for _, referral := range referrals {
			if referral.Type == ReferralTypeInHouse && referral.Status == ReferralStatusOrdered {
				allConfirmed = false
			}
		}

		if allConfirmed {
			r.Status = ReferralOrderStatusCompleted
			if err := tx.Updates(&r).Error; err != nil {
				return err
			}
		}

		return nil
	})
}

// GetCount ...
func (r *ReferralOrder) GetCount(filter *ReferralOrder, date *time.Time, searchTerm *string) (int64, error) {
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
func (r *ReferralOrder) Search(p PaginationInput, filter *ReferralOrder, date *time.Time, searchTerm *string, ascending bool) ([]ReferralOrder, int64, error) {
	var result []ReferralOrder

	dbOp := DB.Scopes(Paginate(&p)).Select("*, count(*) OVER() AS count").Where(filter).Preload("Referrals", "type = ?", ReferralTypeInHouse).Preload("OrderedBy.UserTypes")

	if date != nil {
		createdAt := *date
		start := time.Date(createdAt.Year(), createdAt.Month(), createdAt.Day(), 0, 0, 0, 0, createdAt.Location())
		end := start.AddDate(0, 0, 1)
		dbOp.Where("created_at >= ?", start).Where("referral_orders.created_at <= ?", end)
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
func (r *ReferralOrder) GetByPatientChartID(patientChartID int) error {
	return DB.Where("patient_chart_id = ?", patientChartID).Preload("Referrals").Preload("OrderedBy.UserTypes").Take(&r).Error
}

// GetAll ...
func (r *ReferralOrder) GetAll(p PaginationInput, filter *FollowUpOrder) ([]ReferralOrder, int64, error) {
	var result []ReferralOrder

	dbOp := DB.Scopes(Paginate(&p)).Select("*, count(*) OVER() AS count").Where(filter).Preload("Referrals").Order("id ASC").Find(&result)

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
func (r *ReferralOrder) Update() error {
	return DB.Updates(&r).Error
}

// Delete ...
func (r *ReferralOrder) Delete(ID int) error {
	return DB.Where("id = ?", ID).Delete(&r).Error
}
