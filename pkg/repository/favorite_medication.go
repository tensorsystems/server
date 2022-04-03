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

// FavoriteMedication ...
type FavoriteMedication struct {
	gorm.Model
	ID                  int     `gorm:"primaryKey" json:"id"`
	Medication          string  `json:"medication"`
	Sig                 string  `json:"sig"`
	RxCui               *string `json:"rxCui"`
	Synonym             *string `json:"synonym"`
	Tty                 *string `json:"tty"`
	Language            *string `json:"language"`
	Refill              int     `json:"refill"`
	Generic             bool    `json:"generic"`
	SubstitutionAllowed bool    `json:"substitutionAllowed"`
	DirectionToPatient  string  `json:"directionToPatient"`
	UserID              int     `json:"user_id"`
	User                User    `json:"user"`
	Count               int64   `json:"count"`
}

// Save ...
func (r *FavoriteMedication) Save() error {
	if err := DB.Create(&r).Error; err != nil {
		return err
	}

	return nil
}

// Get ...
func (r *FavoriteMedication) Get(ID int) error {
	return DB.Where("id = ?", ID).Take(&r).Error
}

// GetAll ...
func (r *FavoriteMedication) GetAll(p PaginationInput, filter *FavoriteMedication, searchTerm *string) ([]FavoriteMedication, int64, error) {
	var result []FavoriteMedication

	dbOp := DB.Scopes(Paginate(&p)).Select("*, count(*) OVER() AS count").Where(filter)

	if searchTerm != nil {
		dbOp.Where("medication ILIKE ?", "%"+*searchTerm+"%")
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

// Search ...
func (r *FavoriteMedication) Search(p PaginationInput, searchTerm string) ([]FavoriteMedication, int64, error) {
	var result []FavoriteMedication
	dbOp := DB.Scopes(Paginate(&p)).Select("*, count(*) OVER() AS count").Where("medication LIKE ?", "%"+searchTerm+"%").Order("id ASC").Find(&result)

	var count int64
	if len(result) > 0 {
		count = result[0].Count
	}

	if dbOp.Error != nil {
		return result, 0, dbOp.Error
	}

	return result, count, dbOp.Error
}

// Update ...
func (r *FavoriteMedication) Update() error {
	return DB.Transaction(func(tx *gorm.DB) error {
		err := tx.Updates(&r).Error

		if err, ok := err.(*pq.Error); ok && err.Code.Name() == "unique_violation" {
			return errors.New("Duplicate, " + err.Detail)
		}

		if err != nil {
			return err
		}
		return nil
	})
}

// Delete ...
func (r *FavoriteMedication) Delete(ID int) error {
	return DB.Where("id = ?", ID).Delete(&r).Error
}
