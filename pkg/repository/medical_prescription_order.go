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

// MedicalPrescriptionOrder is a gorm struct for the prescription_order table
type MedicalPrescriptionOrder struct {
	gorm.Model
	ID                   int                   `gorm:"primaryKey"`
	PharmacyID           int                   `json:"pharmacyId"`
	Pharmacy             Pharmacy              `json:"pharmacy"`
	PatientChartID       int                   `json:"patientChartId"`
	OrderedByID          *int                  `json:"orderedById"`
	OrderedBy            *User                 `json:"orderedBy"`
	FirstName            string                `json:"firstName"`
	LastName             string                `json:"lastName"`
	PhoneNo              string                `json:"phoneNo"`
	UserName             string                `json:"userName"`
	Status               string                `json:"status"`
	MedicalPrescriptions []MedicalPrescription `json:"medicalPrescriptions"`
	Count                int64                 `json:"count"`
}

// SaveMedicalPrescription ...
func (r *MedicalPrescriptionOrder) SaveMedicalPrescription(medicalPrescription MedicalPrescription, patientID int) error {
	return DB.Transaction(func(tx *gorm.DB) error {
		var patient Patient
		if err := tx.Where("id = ?", patientID).Take(&patient).Error; err != nil {
			return err
		}

		var user User
		if err := tx.Where("id = ?", r.OrderedByID).Take(&user).Error; err != nil {
			return err
		}

		r.FirstName = patient.FirstName
		r.LastName = patient.LastName
		r.PhoneNo = patient.PhoneNo
		r.UserName = user.FirstName + " " + user.LastName
		r.Status = "ORDERED"

		var existing MedicalPrescriptionOrder
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

		medicalPrescription.Status = "Active"

		tx.Model(&r).Association("MedicalPrescriptions").Append(&medicalPrescription)

		return nil
	})
}

// Search ...
func (r *MedicalPrescriptionOrder) Search(p PaginationInput, filter *MedicalPrescriptionOrder, date *time.Time, searchTerm *string, ascending bool) ([]MedicalPrescriptionOrder, int64, error) {
	var result []MedicalPrescriptionOrder

	dbOp := DB.Scopes(Paginate(&p)).Select("*, count(*) OVER() AS count").Where(filter).Preload("MedicalPrescriptions").Preload("OrderedBy")

	if date != nil {
		prescribedDate := *date
		start := time.Date(prescribedDate.Year(), prescribedDate.Month(), prescribedDate.Day(), 0, 0, 0, 0, prescribedDate.Location())
		end := start.AddDate(0, 0, 1)
		dbOp.Where("created_at >= ?", start).Where("created_at <= ?", end)
	}

	if searchTerm != nil {
		dbOp.Where("first_name ILIKE ?", "%"+*searchTerm+"%").Or("last_name ILIKE ?", "%"+*searchTerm+"%").Or("phone_no ILIKE ?", "%"+*searchTerm+"%")
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
func (r *MedicalPrescriptionOrder) Get(ID int) error {
	return DB.Where("id = ?", ID).Take(&r).Error
}

// GetByPatientChartID ...
func (r *MedicalPrescriptionOrder) GetByPatientChartID(patientChartID int) error {
	return DB.Where("patient_chart_id = ?", patientChartID).Preload("MedicalPrescriptions").Take(&r).Error
}

// GetAll ...
func (r *MedicalPrescriptionOrder) GetAll(p PaginationInput, filter *MedicalPrescriptionOrder) ([]MedicalPrescriptionOrder, int64, error) {
	var result []MedicalPrescriptionOrder

	dbOp := DB.Scopes(Paginate(&p)).Select("*, count(*) OVER() AS count").Where(filter).Preload("MedicalPrescriptions").Order("id ASC").Find(&result)

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
func (r *MedicalPrescriptionOrder) Update() error {
	return DB.Updates(&r).Error
}

// Delete ...
func (r *MedicalPrescriptionOrder) Delete(ID int) error {
	return DB.Where("id = ?", ID).Delete(&r).Error
}
