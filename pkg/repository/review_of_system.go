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

import "gorm.io/gorm"

// ReviewOfSystem ...
type ReviewOfSystem struct {
	gorm.Model
	ID               int           `gorm:"primaryKey" json:"id"`
	PatientHistoryID int           `json:"patientHistoryId"`
	SystemSymptomID  int           `json:"systemSymptomId"`
	SystemSymptom    SystemSymptom `json:"systemSymptom"`
	Note             *string       `json:"note"`
	Count            int64         `json:"count"`
}

// Save ...
func (r *ReviewOfSystem) Save() error {
	return DB.Create(&r).Error
}

// GetAll ...
func (r *ReviewOfSystem) GetAll(p PaginationInput, filter *ReviewOfSystem) ([]ReviewOfSystem, int64, error) {
	var result []ReviewOfSystem

	dbOp := DB.Scopes(Paginate(&p)).Select("*, count(*) OVER() AS count").Preload("SystemSymptom.System").Where(filter).Order("id ASC").Find(&result)

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
func (r *ReviewOfSystem) Get(ID int) error {
	return DB.Where("id = ?", ID).Take(&r).Error
}

// GetByTitle ...
func (r *ReviewOfSystem) GetByPatientHistoryID(id string) error {
	return DB.Where("patient_history_id = ?", id).Take(&r).Error
}

// Update ...
func (r *ReviewOfSystem) Update() error {
	return DB.Updates(&r).Error
}

// Delete ...
func (r *ReviewOfSystem) Delete(ID int) error {
	return DB.Where("id = ?", ID).Delete(&r).Error
}
