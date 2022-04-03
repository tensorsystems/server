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

// PatientEncounterLimit ...
type PatientEncounterLimit struct {
	gorm.Model
	ID             int   `gorm:"primaryKey"`
	UserID         int   `json:"userId" gorm:"uniqueIndex"`
	User           User  `json:"user"`
	MondayLimit    int   `json:"mondayLimit"`
	TuesdayLimit   int   `json:"tuesdayLimit"`
	WednesdayLimit int   `json:"wednesdayLimit"`
	ThursdayLimit  int   `json:"thursdayLimit"`
	FridayLimit    int   `json:"fridayLimit"`
	SaturdayLimit  int   `json:"saturdayLimit"`
	SundayLimit    int   `json:"sundayLimit"`
	Overbook       int   `json:"overbook"`
	Count          int64 `json:"count"`
}

// Save ...
func (r *PatientEncounterLimit) Save() error {
	return DB.Create(&r).Error
}

// Get ...
func (r *PatientEncounterLimit) Get(ID int) error {
	return DB.Where("id = ?", ID).Take(&r).Error
}

// GetByUser ...
func (r *PatientEncounterLimit) GetByUser(userID int) error {
	return DB.Where("user_id = ?", userID).Take(&r).Error
}

// GetAll ...
func (r *PatientEncounterLimit) GetAll(p PaginationInput) ([]PatientEncounterLimit, int64, error) {
	var result []PatientEncounterLimit

	dbOp := DB.Scopes(Paginate(&p)).Select("*, count(*) OVER() AS count").Preload("User").Order("id DESC").Find(&result)

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
func (r *PatientEncounterLimit) Update() error {
	return DB.Updates(&r).Error
}

// Delete ...
func (r *PatientEncounterLimit) Delete(ID int) error {
	return DB.Where("id = ?", ID).Delete(&r).Error
}
