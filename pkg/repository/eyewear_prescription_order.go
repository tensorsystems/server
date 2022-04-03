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

// EyewearPrescriptionOrder ...
type EyewearPrescriptionOrder struct {
	gorm.Model
	ID                   int                   `gorm:"primaryKey"`
	EyewearShopID        int                   `json:"eyewearShopId"`
	EyewearShop          EyewearShop           `json:"eyewearShop"`
	PatientChartID       int                   `json:"patientChartId"`
	FirstName            string                `json:"firstName"`
	LastName             string                `json:"lastName"`
	PhoneNo              string                `json:"phoneNo"`
	UserName             string                `json:"userName"`
	OrderedByID          *int                  `json:"orderedById"`
	OrderedBy            *User                 `json:"orderedBy"`
	Status               string                `json:"status"`
	EyewearPrescriptions []EyewearPrescription `json:"eyewearPrescription"`
	Count                int64                 `json:"count"`
}

// SaveEyewearPrescription ...
func (r *EyewearPrescriptionOrder) SaveEyewearPrescription(eyewearPrescription EyewearPrescription, patientID int) error {
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

		var existing EyewearPrescriptionOrder
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

		eyewearPrescription.Status = "Active"

		tx.Model(&r).Association("EyewearPrescriptions").Append(&eyewearPrescription)

		return nil
	})
}

// Search ...
func (r *EyewearPrescriptionOrder) Search(p PaginationInput, filter *EyewearPrescriptionOrder, date *time.Time, searchTerm *string, ascending bool) ([]EyewearPrescriptionOrder, int64, error) {
	var result []EyewearPrescriptionOrder

	dbOp := DB.Scopes(Paginate(&p)).Select("*, count(*) OVER() AS count").Where(filter).Preload("EyewearPrescriptions").Preload("OrderedBy")

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
func (r *EyewearPrescriptionOrder) Get(ID int) error {
	return DB.Where("id = ?", ID).Take(&r).Error
}

// GetByPatientChartID ...
func (r *EyewearPrescriptionOrder) GetByPatientChartID(patientChartID int) error {
	return DB.Where("patient_chart_id = ?", patientChartID).Preload("EyewearPrescriptions").Take(&r).Error
}

// GetAll ...
func (r *EyewearPrescriptionOrder) GetAll(p PaginationInput, filter *EyewearPrescriptionOrder) ([]EyewearPrescriptionOrder, int64, error) {
	var result []EyewearPrescriptionOrder

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

// Save ...
func (r *EyewearPrescriptionOrder) Save() error {
	return DB.Create(&r).Error
}

// Update ...
func (r *EyewearPrescriptionOrder) Update() error {
	return DB.Updates(&r).Error
}

// Delete ...
func (r *EyewearPrescriptionOrder) Delete(ID int) error {
	return DB.Where("id = ?", ID).Delete(&r).Error
}
