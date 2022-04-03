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
	"errors"

	"github.com/lib/pq"
	"gorm.io/gorm"
)

// Room ...
type Room struct {
	gorm.Model
	ID    int    `gorm:"primaryKey"`
	Title string `json:"title" gorm:"unique"`
}

// Save ...
func (r *Room) Save() error {
	err := DB.Create(&r).Error

	if err, ok := err.(*pq.Error); ok && err.Code.Name() == "unique_violation" {
		var existing Room
		existingErr := DB.Unscoped().Where("title = ?", r.Title).Take(&existing).Error

		if existingErr == nil {
			DB.Model(&Room{}).Unscoped().Where("id = ?", existing.ID).Update("deleted_at", nil)
			r = &existing
			return nil
		}

		return errors.New("Duplicate, " + err.Detail)
	}

	if err != nil {
		return err
	}

	return nil
}

// Get ...
func (r *Room) Get(ID int) error {
	err := DB.Where("id = ?", ID).Take(&r).Error
	if err != nil {
		return err
	}

	return nil
}

// GetByTitle ...
func (r *Room) GetByTitle(title string) error {
	return DB.Where("title = ?", title).Take(&r).Error
}

// Count ...
func (r *Room) Count(dbString string) (int64, error) {
	var count int64

	err := DB.Model(&Room{}).Count(&count).Error
	return count, err
}

// Update ...
func (r *Room) Update() error {
	err := DB.Save(&r).Error

	if err, ok := err.(*pq.Error); ok && err.Code.Name() == "unique_violation" {
		return errors.New("Duplicate, " + err.Detail)
	}

	return nil
}

// Delete ...
func (r *Room) Delete(ID int) error {
	var e Room
	err := DB.Where("id = ?", ID).Delete(&e).Error
	return err
}

// GetAll ...
func (r *Room) GetAll(p PaginationInput) ([]Room, int64, error) {
	var result []Room

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
