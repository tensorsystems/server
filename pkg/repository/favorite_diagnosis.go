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

// FavoriteDiagnosis ...
type FavoriteDiagnosis struct {
	gorm.Model
	ID          int       `gorm:"primaryKey" json:"id"`
	DiagnosisID int       `json:"diagnosisId"`
	Diagnosis   Diagnosis `json:"diagnosis"`
	UserID      int       `json:"user_id"`
	User        User      `json:"user"`
	Count       int64     `json:"count"`
}

// Save ...
func (r *FavoriteDiagnosis) Save() error {
	err := DB.Create(&r).Error
	return err
}

// Get ...
func (r *FavoriteDiagnosis) Get(ID int) error {
	return DB.Where("id = ?", ID).Take(&r).Error
}

// GetByUser ...
func (r *FavoriteDiagnosis) GetByUser(ID int) ([]*FavoriteDiagnosis, error) {
	var result []*FavoriteDiagnosis
	err := DB.Where("user_id = ?", ID).Find(&result).Error
	if err != nil {
		return result, err
	}

	return result, nil
}

// Update ...
func (r *FavoriteDiagnosis) Update() error {
	return DB.Updates(&r).Error
}

// Delete ...
func (r *FavoriteDiagnosis) Delete(id int) error {
	return DB.Where("id = ?", id).Delete(&r).Error
}
