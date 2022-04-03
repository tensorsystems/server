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

// LifestyleType ...
type LifestyleType struct {
	gorm.Model
	ID    int    `gorm:"primaryKey" json:"id"`
	Title string `json:"title"`
}

// Save ...
func (r *LifestyleType) Save() error {
	err := DB.Create(&r).Error
	if err != nil {
		return err
	}

	return nil
}

// Count ...
func (r *LifestyleType) Count(dbString string) (int64, error) {
	var count int64

	err := DB.Model(&LifestyleType{}).Count(&count).Error
	return count, err
}

// GetAll ...
func (r *LifestyleType) GetAll(p PaginationInput) ([]LifestyleType, int64, error) {
	var result []LifestyleType

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

// Get ...
func (r *LifestyleType) Get(ID int) error {
	return DB.Where("id = ?", ID).Take(&r).Error
}

// GetByTitle ...
func (r *LifestyleType) GetByTitle(title string) error {
	return DB.Where("title = ?", title).Take(&r).Error
}

// Update ...
func (r *LifestyleType) Update() error {
	return DB.Updates(&r).Error
}

// Delete ...
func (r *LifestyleType) Delete(ID int) error {
	return DB.Where("id = ?", ID).Delete(&r).Error
}
