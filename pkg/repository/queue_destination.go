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

// QueueDestination ...
type QueueDestination struct {
	gorm.Model
	ID    int    `gorm:"primaryKey" json:"id"`
	Title string `json:"title" gorm:"uniqueIndex"`
	Count int64  `json:"count"`
}

// Seed ...
func (r *QueueDestination) Seed() {
	DB.Create(&QueueDestination{Title: "Front Desk"})
	DB.Create(&QueueDestination{Title: "Exam Room"})
	DB.Create(&QueueDestination{Title: "Operating Room"})
	DB.Create(&QueueDestination{Title: "Emergency Room"})
	DB.Create(&QueueDestination{Title: "Pre-Exam"})
	DB.Create(&QueueDestination{Title: "Optometry"})
}

// Save ...
func (r *QueueDestination) Save() error {
	err := DB.Create(&r).Error
	if err != nil {
		return err
	}

	return nil
}

// GetAll ...
func (r *QueueDestination) GetAll(p PaginationInput) ([]QueueDestination, int64, error) {
	var result []QueueDestination

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

// Get ...
func (r *QueueDestination) Get(ID int) error {
	err := DB.Where("id = ?", ID).Take(&r).Error
	if err != nil {
		return err
	}

	return nil
}

// Update ...
func (r *QueueDestination) Update() (*QueueDestination, error) {
	err := DB.Save(&r).Error
	if err != nil {
		return nil, err
	}
	return r, nil
}

// Delete ...
func (r *QueueDestination) Delete(ID int) error {
	var e QueueDestination
	err := DB.Where("id = ?", ID).Delete(&e).Error
	return err
}

// GetByTitle ...
func (r *QueueDestination) GetByTitle(title string) error {
	return DB.Where("title = ?", title).Take(&r).Error
}

// GetUserTypeFromDestination ...
func (r *QueueDestination) GetUserTypeFromDestination(destination string) string {
	if destination == "PRE_EXAM" {
		return "Nurse"
	}

	if destination == "OPTOMETRY" {
		return "Optometrist"
	}

	if destination == "EXAM_ROOM" {
		return "Physician"
	}

	if destination == "OPERATING_ROOM" {
		return "Physician"
	}

	if destination == "DIAGNOSTIC_PROCEDURE" {
		return "Nurse"
	}

	return ""
}
