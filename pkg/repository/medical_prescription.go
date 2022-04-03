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
	"time"

	"gorm.io/gorm"
)

// MedicalPrescription ...
type MedicalPrescription struct {
	gorm.Model
	ID                         int        `gorm:"primaryKey"`
	MedicalPrescriptionOrderID *int       `json:"medicalPrescriptionOrderId"`
	PatientID                  int        `json:"patientId"`
	Patient                    Patient    `json:"patient"`
	Medication                 string     `json:"medication"`
	RxCui                      *string    `json:"rxCui"`
	Synonym                    *string    `json:"synonym"`
	Tty                        *string    `json:"tty"`
	Language                   *string    `json:"language"`
	Sig                        *string    `json:"sig"`
	Refill                     *int       `json:"refill"`
	Generic                    *bool      `json:"generic"`
	SubstitutionAllowed        *bool      `json:"substitutionAllowed"`
	DirectionToPatient         *string    `json:"directionToPatient"`
	PrescribedDate             *time.Time `json:"prescribedDate"`
	History                    bool       `json:"history"`
	Status                     string     `json:"status"`
	Count                      int64      `json:"count"`
}

// Save ...
func (r *MedicalPrescription) Save() error {
	return DB.Create(&r).Error
}

// Get ...
func (r *MedicalPrescription) Get(ID int) error {
	err := DB.Where("id = ?", ID).Take(&r).Error
	if err != nil {
		return err
	}

	return nil
}

// GetAll ...
func (r *MedicalPrescription) GetAll(p PaginationInput, filter *MedicalPrescription) ([]MedicalPrescription, int64, error) {
	var result []MedicalPrescription

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

// Search ...
func (r *MedicalPrescription) Search(p PaginationInput, filter *MedicalPrescription, date *time.Time, searchTerm *string, ascending bool) ([]MedicalPrescription, int64, error) {
	var result []MedicalPrescription

	dbOp := DB.Scopes(Paginate(&p)).Select("*, count(*) OVER() AS count").Where(filter)

	if date != nil {
		prescribedDate := *date
		start := time.Date(prescribedDate.Year(), prescribedDate.Month(), prescribedDate.Day(), 0, 0, 0, 0, prescribedDate.Location())
		end := start.AddDate(0, 0, 1)
		dbOp.Where("prescribed_date >= ?", start).Where("prescribed_date <= ?", end)
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

// Update ...
func (r *MedicalPrescription) Update() error {
	return DB.Updates(&r).Error
}

// Delete ...
func (r *MedicalPrescription) Delete(ID int) error {
	return DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("id = ?", ID).Take(&r).Error; err != nil {
			return err
		}

		var order MedicalPrescriptionOrder
		if err := tx.Where("id = ?", r.MedicalPrescriptionOrderID).Take(&order).Error; err != nil {
			return err
		}

		var patientChart PatientChart
		if err := tx.Where("id = ?", order.PatientChartID).Take(&patientChart).Error; err != nil {
			return err
		}

		if patientChart.Locked != nil && *patientChart.Locked {
			return errors.New("This prescription cannot be deleted because it's respective chart has been locked")
		}

		if err := tx.Delete(&r).Error; err != nil {
			return err
		}

		return nil
	})
}
