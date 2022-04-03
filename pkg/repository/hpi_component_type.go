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

// HpiComponentType ...
type HpiComponentType struct {
	gorm.Model
	ID          int    `gorm:"primaryKey"`
	Title       string `json:"title"`
	Description string `json:"description"`
}

// Save ...
func (r *HpiComponentType) Save() error {
	return DB.Create(&r).Error
}

// Get ...
func (r *HpiComponentType) Get(ID int) error {
	return DB.Where("id = ?", ID).Take(&r).Error
}

// Update ...
func (r *HpiComponentType) Update() error {
	return DB.Updates(&r).Error
}

// Count ...
func (r *HpiComponentType) Count(dbString string) (int64, error) {
	var count int64

	err := DB.Model(&HpiComponent{}).Count(&count).Error
	return count, err
}

// GetAll ...
func (r *HpiComponentType) GetAll(p PaginationInput) ([]HpiComponentType, int64, error) {
	var result []HpiComponentType

	var count int64
	count, countErr := r.Count("")
	if countErr != nil {
		return result, 0, countErr
	}

	err := DB.Scopes(Paginate(&p)).Order("id ASC").Find(&result).Error
	if err != nil {
		return result, 0, err
	}

	return result, count, err
}

// Delete ...
func (r *HpiComponentType) Delete(ID int) error {
	var e HpiComponent
	err := DB.Where("id = ?", ID).Delete(&e).Error
	return err
}
