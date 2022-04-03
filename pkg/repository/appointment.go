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
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"gorm.io/gorm"
)

// Appointment ...
type Appointment struct {
	gorm.Model
	ID                  int               `gorm:"primaryKey" json:"id"`
	PatientID           int               `json:"patientId"`
	Patient             Patient           `json:"patient"`
	FirstName           string            `json:"firstName"`
	LastName            string            `json:"lastName"`
	PhoneNo             string            `json:"phoneNo"`
	CheckInTime         time.Time         `json:"checkInTime" gorm:"index:check_in_time_idx"`
	CheckedInTime       *time.Time        `json:"checkedInTime" gorm:"index:daily_appointment_idx"`
	CheckedOutTime      time.Time         `json:"checkedOutTime"`
	RoomID              int               `json:"roomId"`
	Room                Room              `json:"room"`
	VisitTypeID         int               `json:"visitTypeId"`
	VisitType           VisitType         `json:"visitType"`
	AppointmentStatusID int               `json:"appointmentStatusId" gorm:"index:daily_appointment_idx"`
	AppointmentStatus   AppointmentStatus `json:"appointmentStatus"`
	Emergency           *bool             `json:"emergency"`
	MedicalDepartment   string            `json:"medicalDepartment"`
	Credit              bool              `json:"credit"`
	Payments            []Payment         `json:"payments" gorm:"many2many:appointment_payments;"`
	Files               []File            `json:"files" gorm:"many2many:appointment_files"`
	UserID              int               `json:"userId"`
	ProviderName        string            `json:"providerName"`
	PatientChart        PatientChart      `json:"patientChart"`
	QueueID             int               `json:"queueId"`
	QueueName           string            `json:"queueName"`
	Document            string            `gorm:"type:tsvector"`
	Count               int64             `json:"count"`
}

// AppointmentSearchInput ...
type AppointmentSearchInput struct {
	SearchTerm          *string    `json:"searchTerm"`
	UserID              *int       `json:"userId"`
	PatientID           *int       `json:"patientId"`
	AppointmentStatusID *string    `json:"appointmentStatusId"`
	VisitTypeID         *string    `json:"visitTypeId"`
	CheckInTime         *time.Time `json:"checkInTime"`
}

// AfterCreate ...
func (r *Appointment) AfterCreate(tx *gorm.DB) error {
	var patient Patient
	err := tx.Where("id = ?", r.PatientID).Take(&patient).Error
	if err != nil {
		return err
	}

	r.FirstName = patient.FirstName
	r.LastName = patient.LastName
	r.PhoneNo = patient.PhoneNo

	var provider User
	tx.Where("id = ?", r.UserID).Take(&provider)
	r.ProviderName = provider.FirstName + " " + provider.LastName

	tx.Model(r).Updates(&r)

	return nil
}

// Save ...
func (r *Appointment) Save() error {
	err := DB.Create(&r).Error
	if err != nil {
		return err
	}

	return nil
}

// CreateNewAppointment ... Creates a new appointment along with PatientChart
func (r *Appointment) CreateNewAppointment(billingID *int, invoiceNo *string) error {
	return DB.Transaction(func(tx *gorm.DB) error {
		var sickVisit VisitType
		if err := tx.Where("title = ?", "Sick Visit").Take(&sickVisit).Error; err != nil {
			return err
		}

		if r.VisitTypeID == sickVisit.ID {
			var existingAppointment Appointment

			checkInTime := r.CheckInTime
			start := time.Date(checkInTime.Year(), checkInTime.Month(), checkInTime.Day(), 0, 0, 0, 0, checkInTime.Location())
			end := start.AddDate(0, 0, 1)

			var status AppointmentStatus
			if err := tx.Where("title = ?", "Checked-In").Take(&status).Error; err != nil {
				return err
			}

			if err := tx.Where("patient_id = ?", r.PatientID).Where("user_id = ?", r.UserID).Where("visit_type_id = ?", r.VisitTypeID).Where("check_in_time >= ?", start).Where("check_in_time < ?", end).Where("appointment_status_id = ?", status.ID).Take(&existingAppointment).Error; err == nil {
				if existingAppointment.ID != 0 {
					return errors.New("Appointment already exists")
				}
			}
		}

		var status AppointmentStatus
		if err := tx.Where("title = ?", "Scheduled").Take(&status).Error; err != nil {
			return err
		}

		r.AppointmentStatusID = status.ID

		if invoiceNo != nil && billingID != nil {
			var payment Payment
			payment.InvoiceNo = *invoiceNo
			payment.Status = PaidPaymentStatus
			payment.BillingID = *billingID

			if err := tx.Create(&payment).Error; err != nil {
				return err
			}

			r.Payments = append(r.Payments, payment)
		}

		if err := tx.Create(&r).Error; err != nil {
			return err
		}

		patientChart := &PatientChart{
			AppointmentID: r.ID,
		}

		if err := tx.Create(&patientChart).Error; err != nil {
			return err
		}

		return nil
	})
}

// SchedulePostOp ...
func (r *Appointment) SchedulePostOp(appointment Appointment) error {
	return DB.Transaction(func(tx *gorm.DB) error {

		var room Room
		if err := tx.Where("title = ?", "Post-Op Room").Take(&room).Error; err != nil {
			room.Title = "Post-Op Room"

			if err := tx.Create(&room).Error; err != nil {
				return err
			}
		}

		var visitType VisitType
		if err := tx.Where("title = ?", "Post-Op").Take(&visitType).Error; err != nil {
			visitType.Title = "Post-Op"

			if err := tx.Create(&visitType).Error; err != nil {
				return err
			}
		}

		var status AppointmentStatus
		if err := tx.Where("title = ?", "Scheduled").Take(&status).Error; err != nil {
			status.Title = "Scheduled"

			if err := tx.Create(&status).Error; err != nil {
				return err
			}
		}

		now := time.Now()
		tomorrow := now.AddDate(0, 0, 1)

		r.CheckInTime = tomorrow
		r.RoomID = room.ID
		r.VisitTypeID = visitType.ID
		r.AppointmentStatus = status
		r.Credit = appointment.Credit
		r.MedicalDepartment = appointment.MedicalDepartment
		r.PatientID = appointment.PatientID
		r.FirstName = appointment.FirstName
		r.LastName = appointment.LastName
		r.PhoneNo = appointment.PhoneNo
		r.ProviderName = appointment.ProviderName
		r.UserID = appointment.UserID

		if err := tx.Create(&r).Error; err != nil {
			return err
		}

		patientChart := &PatientChart{
			AppointmentID: r.ID,
		}

		if err := tx.Create(&patientChart).Error; err != nil {
			return err
		}

		return nil
	})
}

// GetAll ...
func (r *Appointment) GetAll(p PaginationInput, filter *Appointment) ([]Appointment, int64, error) {
	var result []Appointment

	dbOp := DB.Scopes(Paginate(&p)).Select("*, count(*) OVER() AS count").Preload("VisitType").Preload("AppointmentStatus").Where(filter).Order("id ASC").Find(&result)

	var count int64
	if len(result) > 0 {
		count = result[0].Count
	}

	if dbOp.Error != nil {
		return result, 0, dbOp.Error
	}

	return result, count, dbOp.Error
}

// PayForConsultation ...
func (r *Appointment) PayForConsultation(patientID int, date *time.Time) (bool, error) {
	endDate := time.Now()
	if date != nil {
		endDate = *date
	}

	consultationStartDate := endDate.AddDate(0, 0, -16)
	surgeryStartDate := endDate.AddDate(0, 0, -31)

	var consultationCount int64
	var surgeryCount int64

	cErr := DB.Model(&Appointment{}).Where("patient_id = ?", patientID).Where("check_in_time >= ?", consultationStartDate).Where("check_in_time <= ?", endDate).Count(&consultationCount).Error
	if cErr != nil {
		return false, cErr
	}

	var surgicalVisitType VisitType
	if err := DB.Where("title = ?", "Surgery").Take(&surgicalVisitType).Error; err != nil {
		return false, err
	}

	sErr := DB.Model(&Appointment{}).Where("patient_id = ?", patientID).Where("check_in_time >= ?", surgeryStartDate).Where("check_in_time <= ?", endDate).Where("visit_type_id = ?", surgicalVisitType.ID).Count(&surgeryCount).Error
	if sErr != nil {
		return false, sErr
	}

	return consultationCount == 0 && surgeryCount == 0, nil
}

// FindAppointmentsByPatientAndRange ...
func (r *Appointment) FindAppointmentsByPatientAndRange(patientID int, start time.Time, end time.Time) ([]*Appointment, error) {
	var result []*Appointment

	err := DB.Where("patient_id = ?", patientID).Where("check_in_time >= ?", start).Where("check_in_time <= ?", end).Preload("VisitType").Preload("Room").Find(&result).Error

	return result, err
}

// PatientsAppointmentsToday ...
func (r *Appointment) PatientsAppointmentToday(patientID int, checkedIn *bool) (Appointment, error) {
	now := time.Now()
	start := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	end := start.AddDate(0, 0, 1)

	var appointment Appointment
	tx := DB.Model(&Appointment{}).Where("patient_id = ?", patientID).Where("check_in_time >= ?", start).Where("check_in_time <= ?", end)

	if checkedIn != nil {
		var checkInCondition string
		if *checkedIn {
			checkInCondition = "checked_in_time IS NOT NULL"
		} else {
			checkInCondition = "checked_in_time IS NULL"
		}

		tx.Where(checkInCondition)
	}

	err := tx.Preload("VisitType").Preload("Room").Preload("AppointmentStatus").Preload("Patient").Take(&appointment).Error

	return appointment, err
}

// FindTodaysAppointments ...
func (r *Appointment) FindTodaysAppointments(p PaginationInput, searchTerm *string) ([]Appointment, int64, error) {
	var result []Appointment

	now := time.Now()
	start := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	end := start.AddDate(0, 0, 1)

	tx := DB.Scopes(Paginate(&p)).Select("*, count(*) OVER() AS count").Where("check_in_time >= ?", start).Where("check_in_time < ?", end).Preload("VisitType").Preload("Room").Preload("Patient").Preload("AppointmentStatus").Preload("Payments.Billing")

	if searchTerm != nil {
		tx.Where("document @@ plainto_tsquery(?)", *searchTerm)
	}

	tx.Order("check_in_time ASC").Find(&result)

	var count int64
	if len(result) > 0 {
		count = result[0].Count
	}

	if tx.Error != nil {
		return result, 0, tx.Error
	}

	return result, count, tx.Error
}

// FindTodaysCheckedInAppointments ...
func (r *Appointment) FindTodaysCheckedInAppointments(p PaginationInput, searchTerm *string, visitTypes []string) ([]Appointment, int64, error) {
	var result []Appointment

	now := time.Now()
	start := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	end := start.AddDate(0, 0, 1)

	var status AppointmentStatus
	if err := status.GetByTitle("Checked-In"); err != nil {
		return result, 0, err
	}

	tx := DB.Scopes(Paginate(&p)).Select("*, count(*) OVER() AS count").Preload("VisitType").Preload("Room").Preload("Patient").Preload("AppointmentStatus").Where("checked_in_time >= ?", start).Where("checked_in_time < ?", end).Where("appointment_status_id = ?", status.ID)

	if searchTerm != nil {
		tx.Where("document @@ plainto_tsquery(?)", *searchTerm)
	}

	if len(visitTypes) > 0 {
		var appointmentVisitType VisitType
		v, err := appointmentVisitType.GetByTitles(visitTypes)
		if err != nil {
			return nil, 0, err
		}

		var visitTypeIds []int
		for _, visitType := range v {
			visitTypeIds = append(visitTypeIds, visitType.ID)
		}

		tx.Where("visit_type_id IN ?", visitTypeIds)
	}

	err := tx.Order("check_in_time ASC").Find(&result).Error

	var count int64
	if len(result) > 0 {
		count = result[0].Count
	}

	return result, count, err
}

// GetByIds ...
func (r *Appointment) GetByIds(ids []int, p PaginationInput) ([]Appointment, int64, error) {
	var result []Appointment

	dbOp := DB.Scopes(Paginate(&p)).Select("*, count(*) OVER() AS count").Where("id IN ?", ids).Find(&result)

	var count int64
	if len(result) > 0 {
		count = result[0].Count
	}

	if dbOp.Error != nil {
		return result, 0, dbOp.Error
	}

	return result, count, dbOp.Error
}

// GetByIds ...
func (r *Appointment) FindByUserSubscriptions(ids []int, searchTerm *string, visitTypes []string, p PaginationInput) ([]Appointment, int64, error) {
	var result []Appointment

	var i []string
	for _, id := range ids {
		i = append(i, strconv.Itoa(id))
	}

	idsString := strings.Join(i, ",")

	join := fmt.Sprintf("JOIN unnest('{%s}'::int[]) WITH ORDINALITY t(id, ord) USING (id)", idsString)

	tx := DB.Scopes(Paginate(&p)).Select("*, count(*) OVER() AS count").Joins(join)

	if searchTerm != nil {
		tx.Where("document @@ plainto_tsquery(?)", *searchTerm)
	}

	if len(visitTypes) > 0 {
		var appointmentVisitType VisitType
		v, err := appointmentVisitType.GetByTitles(visitTypes)
		if err != nil {
			return nil, 0, err
		}

		var visitTypeIds []int
		for _, visitType := range v {
			visitTypeIds = append(visitTypeIds, visitType.ID)
		}

		tx.Where("visit_type_id IN ?", visitTypeIds)
	}

	err := tx.Preload("VisitType").Preload("Room").Preload("Patient").Preload("AppointmentStatus").Order("t.ord").Find(&result).Error

	var count int64
	if len(result) > 0 {
		count = result[0].Count
	}

	return result, count, err
}

// FindByProvider ...
func (r *Appointment) FindByProvider(p PaginationInput, searchTerm *string, visitTypes []string, userID int) ([]Appointment, int64, error) {
	var result []Appointment

	now := time.Now()
	start := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	end := start.AddDate(0, 0, 1)

	var status AppointmentStatus
	if err := status.GetByTitle("Checked-In"); err != nil {
		return result, 0, err
	}

	tx := DB.Scopes(Paginate(&p)).Select("*, count(*) OVER() AS count").Preload("VisitType").Preload("Room").Preload("Patient").Preload("AppointmentStatus").Where("check_in_time >= ?", start).Where("check_in_time < ?", end).Where("appointment_status_id = ?", status.ID).Where("user_id = ?", userID)

	if searchTerm != nil {
		tx.Where("document @@ plainto_tsquery(?)", *searchTerm)
	}

	if len(visitTypes) > 0 {
		var appointmentVisitType VisitType
		v, err := appointmentVisitType.GetByTitles(visitTypes)
		if err != nil {
			return nil, 0, err
		}

		var visitTypeIds []int
		for _, visitType := range v {
			visitTypeIds = append(visitTypeIds, visitType.ID)
		}

		tx.Where("visit_type_id IN ?", visitTypeIds)
	}

	err := tx.Order("check_in_time ASC").Find(&result).Error

	var count int64
	if len(result) > 0 {
		count = result[0].Count
	}

	return result, count, err
}

// SearchAppointments ...
func (r *Appointment) SearchAppointments(page PaginationInput, p AppointmentSearchInput) ([]Appointment, int64, error) {
	var count int64
	var appointments []Appointment

	tx := DB.Scopes(Paginate(&page)).Select("*, count(*) OVER() AS count").Preload("VisitType").Preload("Room").Preload("Patient").Preload("AppointmentStatus").Preload("Payments.Billing")

	if p.SearchTerm != nil {
		tx.Where("document @@ plainto_tsquery(?)", p.SearchTerm)
	}

	if p.AppointmentStatusID != nil {
		tx.Where("appointment_status_id = ?", p.AppointmentStatusID)
	}

	if p.VisitTypeID != nil {
		tx.Where("visit_type_id = ?", p.VisitTypeID)
	}

	if p.UserID != nil {
		tx.Where("user_id = ?", p.UserID)
	}

	if p.PatientID != nil {
		tx.Where("patient_id = ?", p.PatientID)
	}

	if p.CheckInTime != nil {
		checkInTime := p.CheckInTime
		start := time.Date(checkInTime.Year(), checkInTime.Month(), checkInTime.Day(), 0, 0, 0, 0, checkInTime.Location())
		end := time.Date(checkInTime.Year(), checkInTime.Month(), checkInTime.Day(), 24, 0, 0, 0, checkInTime.Location())

		tx.Where("check_in_time >= ?", start).Where("check_in_time <= ?", end)
	}

	tx.Order("check_in_time DESC")

	err := tx.Find(&appointments).Error

	if len(appointments) > 0 {
		count = appointments[0].Count
	}

	return appointments, count, err
}

// Get ...
func (r *Appointment) Get(ID int) error {
	err := DB.Where("id = ?", ID).Take(&r).Error
	if err != nil {
		return err
	}

	return nil
}

// GetWithDetails ...
func (r *Appointment) GetWithDetails(ID int) error {
	err := DB.Where("id = ?", ID).Preload("VisitType").Preload("Room").Preload("Patient").Preload("AppointmentStatus").Preload("PatientChart.SurgicalProcedure.SurgicalProcedureType").Preload("PatientChart.SurgicalProcedure.PreanestheticDocuments").Preload("PatientChart.Treatment.TreatmentType").Preload("Patient.PaperRecordDocument").Preload("Patient.Documents").Take(&r).Error
	if err != nil {
		return err
	}

	return nil
}

// ReceptionHomeStats ...
func (r *Appointment) ReceptionHomeStats() (int, int, int, error) {
	var appointments []Appointment

	now := time.Now()
	start := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	end := start.AddDate(0, 0, 1)

	err := DB.Select("*, count(*) OVER() AS count").Where("check_in_time >= ?", start).Where("check_in_time < ?", end).Preload("AppointmentStatus").Find(&appointments).Error
	if err != nil {
		return 0, 0, 0, nil
	}

	var scheduled int
	var checkedIn int
	var checkedOut int
	for _, e := range appointments {
		if e.AppointmentStatus.Title == "Scheduled" {
			scheduled = scheduled + 1
		} else if e.AppointmentStatus.Title == "Checked-In" {
			checkedIn = checkedIn + 1
		} else if e.AppointmentStatus.Title == "Checked-Out" {
			checkedOut = checkedOut + 1
		}
	}

	return scheduled, checkedIn, checkedOut, nil
}

// NurseHomeStats ...
func (r *Appointment) NurseHomeStats() (int, int, int, error) {
	var appointments []Appointment

	now := time.Now()
	start := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	end := start.AddDate(0, 0, 1)

	err := DB.Select("*, count(*) OVER() AS count").Where("check_in_time >= ?", start).Where("check_in_time < ?", end).Preload("AppointmentStatus").Find(&appointments).Error
	if err != nil {
		return 0, 0, 0, nil
	}

	var scheduled int
	var checkedIn int
	var checkedOut int
	for _, e := range appointments {
		if e.AppointmentStatus.Title == "Scheduled" {
			scheduled = scheduled + 1
		} else if e.AppointmentStatus.Title == "Checked-In" {
			checkedIn = checkedIn + 1
		} else if e.AppointmentStatus.Title == "Checked-Out" {
			checkedOut = checkedOut + 1
		}
	}

	return scheduled, checkedIn, checkedOut, nil
}

// PhysicianHomeStats ...
func (r *Appointment) PhysicianHomeStats(userId int) (int, int, int, error) {
	var appointments []Appointment

	now := time.Now()
	start := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	end := start.AddDate(0, 0, 1)

	err := DB.Select("*, count(*) OVER() AS count").Where("check_in_time >= ?", start).Where("check_in_time < ?", end).Where("user_id = ?", userId).Preload("AppointmentStatus").Find(&appointments).Error
	if err != nil {
		return 0, 0, 0, nil
	}

	var scheduled int
	var checkedIn int
	var checkedOut int
	for _, e := range appointments {
		if e.AppointmentStatus.Title == "Scheduled" {
			scheduled = scheduled + 1
		} else if e.AppointmentStatus.Title == "Checked-In" {
			checkedIn = checkedIn + 1
		} else if e.AppointmentStatus.Title == "Checked-Out" {
			checkedOut = checkedOut + 1
		}
	}

	return len(appointments), checkedIn, checkedOut, nil
}

// Update ...
func (r *Appointment) Update() error {
	return DB.Updates(&r).Error
}

// Delete ...
func (r *Appointment) Delete(ID int) error {
	var e Appointment
	err := DB.Where("id = ?", ID).Delete(&e).Error
	return err
}
