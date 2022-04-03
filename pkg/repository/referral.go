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
	"gorm.io/gorm"
)

// ReferralStatus ...
type ReferralStatus string

// ReferralType ...
type ReferralType string

// SurgicalProcedureOrder statuses ...
const (
	ReferralStatusOrdered   ReferralStatus = "ORDERED"
	ReferralStatusCompleted ReferralStatus = "COMPLETED"

	ReferralTypeInHouse   ReferralType = "PATIENT_IN_HOUSE_REFERRAL"
	ReferralTypeOutsource ReferralType = "PATIENT_OUTSOURCE_REFERRAL"
)

// Referral ...
type Referral struct {
	gorm.Model
	ID              int            `gorm:"primaryKey"`
	ReferralOrderID int            `json:"referralOrderId"`
	PatientChartID  int            `json:"patientChartId"`
	Reason          string         `json:"reason"`
	ReferredToID    *int           `json:"referredToId"`
	ReferredToName  string         `json:"referredToName"`
	Status          ReferralStatus `json:"status"`
	Type            ReferralType   `json:"type"`
	ReceptionNote   string         `json:"receptionNote"`
	Count           int64          `json:"count"`
}

// Save ...
func (r *Referral) Save() error {
	return DB.Create(&r).Error
}

// Get ...
func (r *Referral) Get(ID int) error {
	return DB.Where("id = ?", ID).Take(&r).Error
}

// Get ...
func (r *Referral) GetByOrderID(ID int) error {
	return DB.Where("patient_chart_id = ?", ID).Take(&r).Error
}

// // ConfirmReferral ...
// func (r *Referral) ConfirmReferral(orderID int, checkInTime time.Time, roomId int) error {
// 	return DB.Transaction(func(tx *gorm.DB) error {
// 		var order Order
// 		if err := order.Get(orderID); err != nil {
// 			return err
// 		}

// 		var referral Referral
// 		if err := referral.GetByOrderID(order.ID); err != nil {
// 			return err
// 		}

// 		*r = referral

// 		var appointment Appointment
// 		if err := appointment.GetWithDetails(referral.AppointmentID); err != nil {
// 			return err
// 		}

// 		var visitType VisitType
// 		if err := tx.Where("title = ?", "Referral").Take(&visitType).Error; err != nil {
// 			return err
// 		}

// 		var status AppointmentStatus
// 		if err := tx.Where("title = ?", "Scheduled").Take(&status).Error; err != nil {
// 			return err
// 		}

// 		var newAppointment Appointment
// 		newAppointment.AppointmentStatusID = status.ID
// 		newAppointment.CheckInTime = checkInTime
// 		newAppointment.Payments = appointment.Payments
// 		newAppointment.PatientID = appointment.PatientID
// 		newAppointment.VisitTypeID = visitType.ID
// 		newAppointment.RoomID = roomId
// 		newAppointment.UserID = appointment.UserID
// 		newAppointment.PatientChart = appointment.PatientChart
// 		if err := tx.Create(&newAppointment).Error; err != nil {
// 			return err
// 		}

// 		order.Status = CompletedOrderStatus
// 		if err := order.Update(); err != nil {
// 			return err
// 		}

// 		return nil
// 	})
// }

// GetAll ...
func (r *Referral) GetAll(p PaginationInput, filter *Referral) ([]Referral, int64, error) {
	var result []Referral

	dbOp := DB.Scopes(Paginate(&p)).Select("*, count(*) OVER() AS count").Where(filter).Order("id ASC").Find(&result)

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
func (r *Referral) Update() error {
	return DB.Updates(&r).Error
}

// Delete ...
func (r *Referral) Delete(ID int) error {
	return DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("id = ?", ID).Take(&r).Error; err != nil {
			return err
		}

		if err := tx.Where("id = ?", ID).Delete(&r).Error; err != nil {
			return err
		}

		var referralsCount int64
		if err := tx.Model(&r).Where("referral_order_id = ?", r.ReferralOrderID).Count(&referralsCount).Error; err != nil {
			return err
		}

		if referralsCount == 0 {
			if err := tx.Where("id = ?", r.ReferralOrderID).Delete(&ReferralOrder{}).Error; err != nil {
				return err
			}
		}

		return nil
	})
}

// NewOrder ...
// func (r *Referral) NewOrder(appointmentID int, orderedByID int, orderedToID int, patientID int, reason *string) error {
// 	return DB.Transaction(func(tx *gorm.DB) error {
// 		// Get Patient
// 		var patient Patient
// 		if err := tx.Model(&Patient{}).Where("id = ?", patientID).Take(&patient).Error; err != nil {
// 			return err
// 		}

// 		// Get Appointment
// 		var appointment Appointment
// 		if err := tx.Model(&Appointment{}).Where("id = ?", appointmentID).Preload("PatientChart").Take(&appointment).Error; err != nil {
// 			return err
// 		}

// 		// Get referred to
// 		var referredTo User
// 		if err := tx.Model(&User{}).Where("id = ?", orderedToID).Take(&referredTo).Error; err != nil {
// 			return err
// 		}

// 		// Create order
// 		var order Order
// 		order.UserID = orderedByID
// 		order.Status = OrderedOrderStatus
// 		order.OrderType = InHouseReferralOrderType
// 		order.Note = *reason
// 		order.FirstName = patient.FirstName
// 		order.LastName = patient.LastName
// 		order.PhoneNo = patient.PhoneNo
// 		order.PatientID = patient.ID
// 		order.PatientChartID = appointment.PatientChart.ID

// 		if err := tx.Create(&order).Error; err != nil {
// 			return err
// 		}

// 		// Create referral
// 		r.Reason = *reason
// 		r.ReferredDate = time.Now()
// 		r.OrderedByID = orderedByID
// 		r.ReferredToID = referredTo.ID
// 		r.ReferredToName = referredTo.FirstName + " " + referredTo.LastName
// 		r.AppointmentID = appointmentID
// 		r.OrderID = order.ID
// 		if err := tx.Create(&r).Error; err != nil {
// 			return err
// 		}

// 		return nil
// 	})
// }
