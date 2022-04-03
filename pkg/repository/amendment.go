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

// Amendment ...
type Amendment struct {
	gorm.Model
	ID             int    `gorm:"primaryKey"`
	PatientChartID int    `json:"patientChartId"`
	Note           string `json:"note"`
	Count          int64  `json:"count"`
}

// Create ...
func (r *Amendment) Create() error {
	return DB.Create(&r).Error
}

// Get ...
func (r *Amendment) Get(ID int) error {
	return DB.Where("id = ?", ID).Take(&r).Error
}

// GetAll ...
func (r *Amendment) GetAll(filter *Amendment) ([]*Amendment, error) {
	var result []*Amendment
	err := DB.Where(filter).Order("id ASC").Find(&result).Error
	return result, err
}

// Update ...
func (r *Amendment) Update() error {
	return DB.Updates(&r).Error
}

// Delete ...
func (r *Amendment) Delete(ID int) error {
	return DB.Where("id = ?", ID).Delete(&r).Error
}
