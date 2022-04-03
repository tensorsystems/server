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

// Billing ...
type Billing struct {
	gorm.Model
	ID       int     `gorm:"primaryKey"`
	Item     string  `json:"item"`
	Code     string  `json:"code"`
	Price    float64 `json:"price"`
	Credit   bool    `json:"credit"`
	Remark   string  `json:"remark"`
	Document string  `gorm:"type:tsvector"`
	Count    int64   `json:"count"`
}

// Save ...
func (r *Billing) Save() error {
	err := DB.Create(&r).Error
	if err != nil {
		return err
	}

	return nil
}

// Get ...
func (r *Billing) Get(ID int) error {
	err := DB.Where("id = ?", ID).Take(&r).Error
	if err != nil {
		return err
	}

	return nil
}

// GetByIds ...
func (r *Billing) GetByIds(ids []*int) ([]Billing, error) {
	var result []Billing

	err := DB.Where("id IN ?", ids).Find(&result).Error

	if err != nil {
		return result, err
	}

	return result, nil
}

// GetAll ...
func (r *Billing) GetAll(p PaginationInput) ([]Billing, int64, error) {
	var result []Billing

	dbOp := DB.Scopes(Paginate(&p)).Select("*, count(*) OVER() AS count").Order("id ASC").Find(&result)

	var count int64
	if len(result) > 0 {
		count = result[0].Count
	}

	if dbOp.Error != nil {
		return result, 0, dbOp.Error
	}

	return result, count, dbOp.Error
}

// GetConsultationBillings ...
func (r *Billing) GetConsultationBillings() ([]*Billing, error) {
	var result []*Billing

	err := DB.Where("item ILIKE ?", "%Consultation%").Find(&result).Error

	if err != nil {
		return result, err
	}

	return result, err
}

// Search ...
func (r *Billing) Search(p PaginationInput, filter *Billing, searchTerm *string) ([]Billing, int64, error) {
	var result []Billing

	tx := DB.Scopes(Paginate(&p)).Select("*, count(*) OVER() AS count").Where(filter).Order("id ASC")

	if searchTerm != nil {
		tx.Where("document @@ plainto_tsquery(?)", *searchTerm)
	}

	err := tx.Find(&result).Error

	if err != nil {
		return result, 0, err
	}

	var count int64
	if len(result) > 0 {
		count = result[0].Count
	}

	return result, count, err
}

// Update ...
func (r *Billing) Update() (*Billing, error) {
	err := DB.Save(&r).Error
	if err != nil {
		return nil, err
	}
	return r, nil
}

// Delete ...
func (r *Billing) Delete(ID int) error {
	var e Billing
	err := DB.Where("id = ?", ID).Delete(&e).Error
	return err
}
