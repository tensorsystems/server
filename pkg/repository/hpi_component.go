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

// HpiComponent ...
type HpiComponent struct {
	gorm.Model
	ID                 int              `gorm:"primaryKey"`
	Title              string           `json:"title"`
	HpiComponentTypeID int              `json:"hpiComponentTypeId"`
	HpiComponentType   HpiComponentType `json:"hpiComponentType"`
	Count              int64            `json:"count"`
}

// Save ...
func (r *HpiComponent) Save() error {
	return DB.Create(&r).Error
}

// Get ...
func (r *HpiComponent) Get(ID int) error {
	err := DB.Where("id = ?", ID).Preload("HpiComponentType").Take(&r).Error
	if err != nil {
		return err
	}

	return nil
}

// GetByIds ...
func (r *HpiComponent) GetByIds(ids []*int) ([]HpiComponent, error) {
	var result []HpiComponent

	err := DB.Where("id IN ?", ids).Find(&result).Error

	if err != nil {
		return result, err
	}

	return result, nil
}

// Update ...
func (r *HpiComponent) Update() error {
	return DB.Updates(&r).Error
}

// GetAll ...
func (r *HpiComponent) GetAll(p PaginationInput, filter *HpiComponent) ([]HpiComponent, int64, error) {
	var result []HpiComponent

	dbOp := DB.Scopes(Paginate(&p)).Select("*, count(*) OVER() AS count").Where(filter).Preload("HpiComponentType").Order("id ASC").Find(&result)

	var count int64
	if len(result) > 0 {
		count = result[0].Count
	}

	if dbOp.Error != nil {
		return result, 0, dbOp.Error
	}

	return result, count, dbOp.Error
}

// Search ...
func (r *HpiComponent) Search(p PaginationInput, filter *HpiComponent, searchTerm *string) ([]HpiComponent, int64, error) {
	var result []HpiComponent

	tx := DB.Scopes(Paginate(&p)).Select("*, count(*) OVER() AS count").Where(filter)

	if searchTerm != nil {
		tx.Where("title ILIKE ?", "%"+*searchTerm+"%")
	}

	err := tx.Preload("HpiComponentType").Order("id ASC").Find(&result).Error

	if err != nil {
		return result, 0, err
	}

	var count int64
	if len(result) > 0 {
		count = result[0].Count
	}

	return result, count, err
}

// Delete ...
func (r *HpiComponent) Delete(ID int) error {
	var e HpiComponent
	err := DB.Where("id = ?", ID).Delete(&e).Error
	return err
}
