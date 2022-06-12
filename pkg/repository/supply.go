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

// Supply ...
type Supply struct {
	gorm.Model
	ID       int       `gorm:"primaryKey" json:"id"`
	Title    string    `json:"title"`
	Active   bool      `json:"active"`
	Billings []Billing `json:"billings" gorm:"many2many:supply_billings"`
	Count    int64     `json:"count"`
}

// Save ...
func (r *Supply) Save() error {
	err := DB.Create(&r).Error
	if err != nil {
		return err
	}

	return nil
}

// GetAll ...
func (r *Supply) GetAll(p PaginationInput, searchTerm *string) ([]Supply, int64, error) {
	var result []Supply

	dbOp := DB.Scopes(Paginate(&p)).Select("*, count(*) OVER() AS count")

	if searchTerm != nil {
		dbOp.Where("title ILIKE ?", "%"+*searchTerm+"%")
	}

	dbOp.Order("id ASC").Preload("Billings").Find(&result)

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
func (r *Supply) Get(ID int) error {
	return DB.Where("id = ?", ID).Take(&r).Error
}

// GetByIds ...
func (r *Supply) GetByIds(ids []*int) ([]Supply, error) {
	var result []Supply

	err := DB.Where("id IN ?", ids).Find(&result).Error

	if err != nil {
		return result, err
	}

	return result, nil
}

// GetByTitle ...
func (r *Supply) GetByTitle(title string) error {
	return DB.Where("title = ?", title).Take(&r).Error
}

// Update ...
func (r *Supply) Update() error {
	return DB.Transaction(func(tx *gorm.DB) error {
		tx.Model(&r).Association("Billings").Replace(&r.Billings)

		return tx.Updates(&r).Error
	})
}

// Delete ...
func (r *Supply) Delete(ID int) error {
	return DB.Where("id = ?", ID).Delete(&r).Error
}
