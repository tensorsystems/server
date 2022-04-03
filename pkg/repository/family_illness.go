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

// FamilyIllness ...
type FamilyIllness struct {
	gorm.Model
	ID               int    `gorm:"primaryKey" json:"id"`
	Title            string `json:"title"`
	Description      string `json:"description"`
	PatientHistoryID int    `json:"patientHistoryID"`
}

// Save ...
func (r *FamilyIllness) Save() error {
	err := DB.Create(&r).Error
	if err != nil {
		return err
	}

	return nil
}

// Get ...
func (r *FamilyIllness) Get(ID int) error {
	err := DB.Where("id = ?", ID).Take(&r).Error
	if err != nil {
		return err
	}

	return nil
}

// GetByPatientHistoryID ...
func (r *FamilyIllness) GetByPatientHistoryID(ID int) ([]*FamilyIllness, error) {
	var result []*FamilyIllness

	err := DB.Where("patient_history_id = ?", ID).Find(&result).Error

	if err != nil {
		return result, err
	}

	return result, nil
}

// Update ...
func (r *FamilyIllness) Update() error {
	return DB.Updates(&r).Error
}

// Delete ...
func (r *FamilyIllness) Delete(ID int) error {
	var e FamilyIllness
	err := DB.Where("id = ?", ID).Delete(&e).Error
	return err
}
