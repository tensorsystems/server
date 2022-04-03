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

// PastHospitalization ...
type PastHospitalization struct {
	gorm.Model
	ID               int        `gorm:"primaryKey"`
	Reason           string     `json:"reason"`
	Provider         string     `json:"provider"`
	From             *time.Time `json:"from"`
	To               *time.Time `json:"to"`
	PatientHistoryID int        `json:"patientHistoryID"`
}

// Save ...
func (r *PastHospitalization) Save() error {
	err := DB.Create(&r).Error
	if err != nil {
		return err
	}

	return nil
}

// Get ...
func (r *PastHospitalization) Get(ID int) error {
	err := DB.Where("id = ?", ID).Take(&r).Error
	if err != nil {
		return err
	}

	return nil
}

// GetByPatientHistoryID ...
func (r *PastHospitalization) GetByPatientHistoryID(ID int) ([]*PastHospitalization, error) {
	var result []*PastHospitalization

	err := DB.Where("patient_history_id = ?", ID).Find(&result).Error

	if err != nil {
		return result, err
	}

	return result, nil
}

// Update ...
func (r *PastHospitalization) Update() error {
	return DB.Updates(&r).Error
}

// Delete ...
func (r *PastHospitalization) Delete(ID int) error {
	var e PastHospitalization
	err := DB.Where("id = ?", ID).Delete(&e).Error
	return err
}
