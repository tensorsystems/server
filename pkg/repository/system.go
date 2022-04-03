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

// System ...
type System struct {
	gorm.Model
	ID     int    `json:"id" gorm:"primaryKey"`
	Title  string `json:"title"`
	Count  int64  `json:"count"`
	Active bool   `json:"active"`
}

// Save ...
func (r *System) Save() error {
	err := DB.Create(&r).Error
	if err != nil {
		return err
	}

	return nil
}

// GetAll ...
func (r *System) GetAll(p PaginationInput, searchTerm *string) ([]System, int64, error) {
	var result []System

	dbOp := DB.Scopes(Paginate(&p)).Select("*, count(*) OVER() AS count")

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
func (r *System) Get(ID int) error {
	return DB.Where("id = ?", ID).Take(&r).Error
}

// GetByTitle ...
func (r *System) GetByTitle(title string) error {
	return DB.Where("title = ?", title).Take(&r).Error
}

// Update ...
func (r *System) Update() error {
	return DB.Updates(&r).Error
}

// Delete ...
func (r *System) Delete(ID int) error {
	return DB.Where("id = ?", ID).Delete(&r).Error
}
