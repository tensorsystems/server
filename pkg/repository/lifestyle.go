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

// Lifestyle ...
type Lifestyle struct {
	gorm.Model
	ID               int    `gorm:"primaryKey"`
	Title            string `json:"title"`
	Description      string `json:"description"`
	Note             string `json:"note"`
	PatientHistoryID int    `json:"patientHistoryID"`
}

// Save ...
func (r *Lifestyle) Save() error {
	err := DB.Create(&r).Error
	if err != nil {
		return err
	}

	return nil
}

// Get ...
func (r *Lifestyle) Get(ID int) error {
	return DB.Where("id = ?", ID).Take(&r).Error
}

// GetByPatientHistoryID ...
func (r *Lifestyle) GetByPatientHistoryID(ID int) ([]*Lifestyle, error) {
	var result []*Lifestyle

	err := DB.Where("patient_history_id = ?", ID).Find(&result).Error

	if err != nil {
		return result, err
	}

	return result, nil
}

// Update ...
func (r *Lifestyle) Update() error {
	return DB.Updates(&r).Error
}

// Delete ...
func (r *Lifestyle) Delete(ID int) error {
	var e Lifestyle
	err := DB.Where("id = ?", ID).Delete(&e).Error
	return err
}
