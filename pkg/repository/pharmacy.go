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

// Pharmacy is a repository for the pharmacy domain.
type Pharmacy struct {
	gorm.Model
	ID      int    `gorm:"primaryKey"`
	Title   string `json:"title" gorm:"uniqueIndex"`
	Address string `json:"address"`
	Region  string `json:"region"`
	Country string `json:"country"`
	Phone   string `json:"phone"`
	InHouse bool   `json:"inHouse"`
	Count   int64  `json:"count"`
	Active  bool   `json:"active"`
}

// Get ...
func (r *Pharmacy) Get(ID int) error {
	return DB.Where("id = ?", ID).Take(&r).Error
}

// GetAll ...
func (r *Pharmacy) GetAll(p PaginationInput, filter *Pharmacy) ([]Pharmacy, int64, error) {
	var result []Pharmacy

	dbOp := DB.Scopes(Paginate(&p)).Select("*, count(*) OVER() AS count").Where(filter).Order("id ASC").Find(&result)

	var count int64
	if len(result) > 0 {
		count = result[0].Count
	}

	if dbOp.Error != nil {
		return result, 0, dbOp.Error
	}

	return result, count, dbOp.Error
}

// Save ...
func (r *Pharmacy) Save() error {
	return DB.Create(&r).Error
}

// Update ...
func (r *Pharmacy) Update() error {
	return DB.Updates(&r).Error
}

// Delete ...
func (r *Pharmacy) Delete(ID int) error {
	return DB.Where("id = ?", ID).Delete(&r).Error
}
