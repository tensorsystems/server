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
	"database/sql/driver"
	"time"

	"gorm.io/gorm"
)

// OrderStatus ...
type OrderStatus string

// OrderType ...
type OrderType string

// Enums ...
const (
	OrderedOrderStatus   OrderStatus = "ORDERED"
	CompletedOrderStatus OrderStatus = "COMPLETED"

	DiagnosticOrderType        OrderType = "DIAGNOSTIC_PROCEDURE"
	SurgicalOrderType          OrderType = "SURGICAL_PROCEDURE"
	FollowUpOrderType          OrderType = "FOLLOW_UP"
	InHouseReferralOrderType   OrderType = "PATIENT_IN_HOUSE_REFERRAL"
	OutsourceReferralOrderType OrderType = "PATIENT_OUTSOURCE_REFERRAL"
	TreatmentOrderType         OrderType = "TREATMENT"
	LabOrderType               OrderType = "LABRATORY"
)

// Scan ...
func (p *OrderStatus) Scan(value interface{}) error {
	*p = OrderStatus(value.(string))
	return nil
}

// Value ...
func (p OrderStatus) Value() (driver.Value, error) {
	return string(p), nil
}

// Scan ...
func (p *OrderType) Scan(value interface{}) error {
	*p = OrderType(value.(string))
	return nil
}

// Value ...
func (p OrderType) Value() (driver.Value, error) {
	return string(p), nil
}

// Order ...
type Order struct {
	gorm.Model
	ID             int         `gorm:"primaryKey"`
	UserID         int         `json:"userID"`
	User           User        `json:"user"`
	FirstName      string      `json:"firstName"`
	LastName       string      `json:"lastName"`
	PhoneNo        string      `json:"phoneNo"`
	UserName       string      `json:"userName"`
	PatientID      int         `json:"patientId"`
	PatientChartID int         `json:"patientChartId"`
	AppointmentID  int         `json:"appointmentId"`
	Emergency      *bool       `json:"emergency"`
	Note           string      `json:"note"`
	Status         OrderStatus `json:"status" sql:"order_status"`
	OrderType      OrderType   `json:"orderType" sql:"order_type"`
	Payments       []Payment   `json:"payments" gorm:"many2many:order_payments;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Count          int64       `json:"count"`
}

// OrderFilterInput ...
type OrderFilterInput struct {
	AppointmentID  *int       `json:"appointmentId"`
	PatientChartID *int       `json:"patientChartId"`
	UserID         *int       `json:"userId"`
	Status         *string    `json:"status"`
	OrderType      *string    `json:"orderType"`
	SearchTerm     *string    `json:"searchTerm"`
	Date           *time.Time `json:"date"`
}

// Save ...
func (r *Order) Save() error {
	return DB.Create(&r).Error
}

// Get ...
func (r *Order) Get(ID int) error {
	return DB.Where("id = ?", ID).Take(&r).Error
}

// GetAll ...
func (r *Order) Search(p PaginationInput, f *OrderFilterInput) ([]Order, int64, error) {
	var result []Order

	tx := DB.Scopes(Paginate(&p)).Select("*, count(*) OVER() AS count")

	if f.Date != nil {
		start := time.Date(f.Date.Year(), f.Date.Month(), f.Date.Day(), 0, 0, 0, 0, f.Date.Location())
		end := start.AddDate(0, 0, 1)
		tx.Where("created_at >= ?", start).Where("created_at <= ?", end)
	}

	if f.UserID != nil {
		tx.Where("user_id = ?", f.UserID)
	}

	if f.AppointmentID != nil {
		tx.Where("appointment_id = ?", f.AppointmentID)
	}

	if f.PatientChartID != nil {
		tx.Where("patient_chart_id = ?", f.PatientChartID)
	}

	if f.OrderType != nil {
		orderType := OrderType(*f.OrderType)
		tx.Where("order_type = ?", orderType)
	}

	if f.Status != nil {
		tx.Where("status = ?", f.Status)
	}

	if f.Status != nil {
		tx.Preload("Payments")
	}

	tx.Preload("Payments.Billing").Preload("User").Order("id ASC").Find(&result)

	var count int64
	if len(result) > 0 {
		count = result[0].Count
	}

	if tx.Error != nil {
		return result, 0, tx.Error
	}

	return result, count, tx.Error
}

// Update ...
func (r *Order) Update() error {
	return DB.Updates(&r).Error
}

// Delete ...
func (r *Order) Delete(ID int) error {
	return DB.Where("id = ?", ID).Delete(&r).Error
}

// ScheduleSurgery ...
func (r *Order) ScheduleSurgery(orderID int, roomID int, checkInTime time.Time, invoiceNo string) error {
	return DB.Transaction(func(tx *gorm.DB) error {
		var order Order

		// Get order
		if err := tx.Where("id = ?", orderID).Preload("Payments").Take(&order).Error; err != nil {
			return err
		}

		// Update order status
		order.Status = CompletedOrderStatus
		if err := tx.Updates(&order).Error; err != nil {
			return err
		}

		// Update every payment status to paid
		for _, payment := range order.Payments {
			payment.Status = PaidPaymentStatus
			payment.InvoiceNo = invoiceNo

			if err := tx.Updates(&payment).Error; err != nil {
				return err
			}
		}

		// Create new surgical appointment
		var appointment Appointment
		appointment.PatientID = order.PatientID
		appointment.RoomID = roomID
		appointment.CheckInTime = checkInTime
		appointment.UserID = order.UserID
		appointment.Credit = false
		appointment.Payments = order.Payments

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

		// Get surgical procedure
		var surgicalProcedure SurgicalProcedure
		if err := tx.Where("patient_chart_id = ?", order.PatientChartID).Take(&surgicalProcedure).Error; err != nil {
			return err
		}

		// Update surgical procedure to new patient chart
		surgicalProcedure.PatientChartID = newPatientChart.ID
		if err := tx.Updates(&surgicalProcedure).Error; err != nil {
			return err
		}

		return nil
	})
}

// ScheduleSurgery ...
func (r *Order) ScheduleTreatment(orderID int, roomID int, checkInTime time.Time, invoiceNo string) error {
	return DB.Transaction(func(tx *gorm.DB) error {
		var order Order

		// Get order
		if err := tx.Where("id = ?", orderID).Preload("Payments").Take(&order).Error; err != nil {
			return err
		}

		// Update order status
		order.Status = CompletedOrderStatus
		if err := tx.Updates(&order).Error; err != nil {
			return err
		}

		// Update every payment status to paid
		for _, payment := range order.Payments {
			payment.Status = PaidPaymentStatus
			payment.InvoiceNo = invoiceNo

			if err := tx.Updates(&payment).Error; err != nil {
				return err
			}
		}

		// Create new treatment
		var appointment Appointment
		appointment.PatientID = order.PatientID
		appointment.RoomID = roomID
		appointment.CheckInTime = checkInTime
		appointment.UserID = order.UserID
		appointment.Credit = false

		// Assign surgery visit type
		var visitType VisitType
		if err := tx.Where("title = ?", "Treatment").Take(&visitType).Error; err != nil {
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

		// Get treatment
		var treatment Treatment
		if err := tx.Where("patient_chart_id = ?", order.PatientChartID).Take(&treatment).Error; err != nil {
			return err
		}

		// Update treatment to new patient chart
		treatment.PatientChartID = newPatientChart.ID
		if err := tx.Updates(&treatment).Error; err != nil {
			return err
		}

		return nil
	})
}

// OrderFollowup ...
func (r *Order) OrderFollowup(appointmentID int, userID int, note *string) error {
	return DB.Transaction(func(tx *gorm.DB) error {
		var order Order

		var appointment Appointment
		if err := tx.Where("id = ?", appointmentID).Preload("PatientChart").Take(&appointment).Error; err != nil {
			return err
		}

		var patient Patient
		if err := tx.Where("id = ?", appointment.PatientID).Take(&patient).Error; err != nil {
			return err
		}

		var user User
		if err := tx.Where("id = ?", userID).Take(&user).Error; err != nil {
			return err
		}

		order.UserID = userID
		order.AppointmentID = appointment.ID
		order.PatientChartID = appointment.PatientChart.ID
		order.PatientID = patient.ID
		order.FirstName = patient.FirstName
		order.LastName = patient.LastName
		order.PhoneNo = patient.PhoneNo
		order.UserName = user.FirstName + " " + user.LastName
		order.Note = *note
		order.OrderType = FollowUpOrderType
		order.Status = OrderedOrderStatus

		if err := tx.Create(&order).Error; err != nil {
			return err
		}

		return nil
	})
}

// ConfirmFollowUpOrder ...
func (r *Order) ConfirmFollowUpOrder(orderID int, checkInTime time.Time, roomID int) error {
	return DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("id = ?", orderID).Take(&r).Error; err != nil {
			return err
		}

		var appointment Appointment
		if err := tx.Where("id = ?", r.AppointmentID).Take(&appointment).Error; err != nil {
			return err
		}

		var visitType VisitType
		if err := tx.Where("title = ?", "Follow Up").Take(&visitType).Error; err != nil {
			return err
		}

		var status AppointmentStatus
		if err := tx.Where("title = ?", "Scheduled").Take(&status).Error; err != nil {
			return err
		}

		var destination QueueDestination
		if err := tx.Where("title = ?", "Front Desk").Take(&destination).Error; err != nil {
			return err
		}

		var newAppointment Appointment
		newAppointment.AppointmentStatusID = status.ID
		newAppointment.CheckInTime = checkInTime
		newAppointment.Payments = appointment.Payments
		newAppointment.PatientID = appointment.PatientID
		newAppointment.VisitTypeID = visitType.ID
		newAppointment.RoomID = roomID
		newAppointment.UserID = appointment.UserID
		newAppointment.PatientChart = appointment.PatientChart

		if err := tx.Create(&newAppointment).Error; err != nil {
			return err
		}

		r.Status = CompletedOrderStatus
		if err := tx.Updates(&r).Error; err != nil {
			return err
		}
		return nil
	})
}

func (r *Order) Counts() (treatment int, referral int, followUps int, errs error) {
	var orders []Order

	now := time.Now()
	start := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	end := start.AddDate(0, 0, 1)

	err := DB.Select("order_type").Where("created_at >= ?", start).Where("created_at <= ?", end).Where("status = ?", OrderedOrderStatus).Find(&orders).Error
	if err != nil {
		return 0, 0, 0, err
	}

	treatmentOrders := 0
	referralOrders := 0
	followUpOrders := 0

	for _, e := range orders {
		if e.OrderType == TreatmentOrderType {
			treatmentOrders = treatmentOrders + 1
		}

		if e.OrderType == InHouseReferralOrderType {
			referralOrders = referralOrders + 1
		}

		if e.OrderType == FollowUpOrderType {
			followUpOrders = followUpOrders + 1
		}
	}

	return treatmentOrders, referralOrders, followUpOrders, nil
}

// ProviderOrders ...
func (r *Order) ProviderOrders(p PaginationInput, searchTerm *string, userId int) ([]Order, int64, error) {
	var orders []Order

	now := time.Now()
	start := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	end := start.AddDate(0, 0, 1)

	var err error

	if searchTerm != nil {
		err = DB.Scopes(Paginate(&p)).Select("*, count(*) OVER() AS count").Where(
			DB.Where("user_id = ?", userId).Where("created_at >= ?", start).Where("created_at <= ?", end).Where("order_type IN ?", []string{"DIAGNOSTIC_PROCEDURE", "SURGICAL_PROCEDURE", "TREATMENT", "LABRATORY"}),
		).Where(
			DB.Where("first_name ILIKE ?", "%"+*searchTerm+"%").Or("last_name ILIKE ?", "%"+*searchTerm+"%").Or("phone_no ILIKE ?", "%"+*searchTerm+"%"),
		).Preload("Payments.Billing").Order("id DESC").Find(&orders).Error
	} else {

		err = DB.Scopes(Paginate(&p)).Select("*, count(*) OVER() AS count").Where("user_id = ?", userId).Where("created_at >= ?", start).Where("created_at <= ?", end).Where("order_type IN ?", []string{"DIAGNOSTIC_PROCEDURE", "SURGICAL_PROCEDURE", "TREATMENT", "LABRATORY"}).Preload("Payments.Billing").Order("id DESC").Find(&orders).Error
	}

	if err != nil {
		return nil, 0, err
	}

	var count int64
	if len(orders) > 0 {
		count = orders[0].Count
	}

	return orders, count, nil
}

// Confirm ...
func (r *Order) Confirm(orderID int, invoiceNo string) error {
	return DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("id = ?", orderID).Preload("Payment").First(&r).Error; err != nil {
			return err
		}

		var paymentIds []int
		for _, payment := range r.Payments {
			paymentIds = append(paymentIds, payment.ID)
		}

		if err := tx.Model(&Payment{}).Where("id IN ?", paymentIds).Updates(map[string]interface{}{"invoice_no": invoiceNo, "status": "PAID"}).Error; err != nil {
			return err
		}

		r.Status = CompletedOrderStatus

		if err := tx.Updates(&r).Error; err != nil {
			return err
		}

		return nil
	})
}
