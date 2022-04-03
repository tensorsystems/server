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
	"gorm.io/gorm"
)

// SystemSymptom
type SystemSymptom struct {
	gorm.Model
	ID       int    `gorm:"primaryKey" json:"id"`
	Title    string `json:"title"`
	Count    int64  `json:"count"`
	Active   bool   `json:"active"`
	SystemID int    `json:"systemId"`
	System   System `json:"system"`
}

// Save ...
func (r *SystemSymptom) Save() error {
	err := DB.Create(&r).Error
	if err != nil {
		return err
	}

	return nil
}

// GetAll ...
func (r *SystemSymptom) GetAll(p PaginationInput, searchTerm *string) ([]SystemSymptom, int64, error) {
	var result []SystemSymptom

	dbOp := DB.Scopes(Paginate(&p)).Select("*, count(*) OVER() AS count").Preload("System")

	if searchTerm != nil {
		dbOp.Where("title ILIKE ?", "%"+*searchTerm+"%")
	}

	dbOp.Order("id ASC").Find(&result)

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
func (r *SystemSymptom) Get(ID int) error {
	return DB.Where("id = ?", ID).Take(&r).Error
}

// GetByTitle ...
func (r *SystemSymptom) GetByTitle(title string) error {
	return DB.Where("title = ?", title).Take(&r).Error
}

// Update ...
func (r *SystemSymptom) Update() error {
	return DB.Updates(&r).Error
}

// Delete ...
func (r *SystemSymptom) Delete(ID int) error {
	return DB.Where("id = ?", ID).Delete(&r).Error
}
