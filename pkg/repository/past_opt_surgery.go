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

// PastOptSurgery ...
type PastOptSurgery struct {
	gorm.Model
	ID               int    `gorm:"primaryKey"`
	Title            string `json:"title"`
	Description      string `json:"description"`
	PatientHistoryID int    `json:"patientHistoryID"`
}

// Save ...
func (r *PastOptSurgery) Save() error {
	err := DB.Create(&r).Error
	if err != nil {
		return err
	}

	return nil
}

// Get ...
func (r *PastOptSurgery) Get(ID int) error {
	err := DB.Where("id = ?", ID).Take(&r).Error
	if err != nil {
		return err
	}

	return nil
}

// Update ...
func (r *PastOptSurgery) Update() (*PastOptSurgery, error) {
	err := DB.Save(&r).Error
	if err != nil {
		return nil, err
	}
	return r, nil
}

// Delete ...
func (r *PastOptSurgery) Delete(ID int) error {
	var e PastOptSurgery
	err := DB.Where("id = ?", ID).Delete(&e).Error
	return err
}
