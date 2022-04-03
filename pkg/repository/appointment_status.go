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

// AppointmentStatus ...
type AppointmentStatus struct {
	gorm.Model
	ID    int    `gorm:"primaryKey"`
	Title string `json:"title" gorm:"uniqueIndex"`
}

//Seed ...
func (r *AppointmentStatus) Seed() {
	DB.Create(&AppointmentStatus{Title: "Scheduled"})
	DB.Create(&AppointmentStatus{Title: "Checked-In"})
	DB.Create(&AppointmentStatus{Title: "Checked-Out"})
	DB.Create(&AppointmentStatus{Title: "Cancelled"})
}

// Save ...
func (r *AppointmentStatus) Save() error {
	err := DB.Create(&r).Error
	if err != nil {
		return err
	}

	return nil
}

// Count ...
func (r *AppointmentStatus) Count(dbString string) (int64, error) {
	var count int64

	err := DB.Model(&AppointmentStatus{}).Count(&count).Error
	return count, err
}

// GetAll ...
func (r *AppointmentStatus) GetAll(p PaginationInput) ([]AppointmentStatus, int64, error) {
	var result []AppointmentStatus

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
func (r *AppointmentStatus) Get(ID int) error {
	return DB.Where("id = ?", ID).Take(&r).Error
}

// GetByTitle ...
func (r *AppointmentStatus) GetByTitle(title string) error {
	return DB.Where("title = ?", title).Take(&r).Error
}

// Update ...
func (r *AppointmentStatus) Update() (*AppointmentStatus, error) {
	err := DB.Save(&r).Error
	if err != nil {
		return nil, err
	}
	return r, nil
}

// Delete ...
func (r *AppointmentStatus) Delete(ID int) error {
	var e AppointmentStatus
	err := DB.Where("id = ?", ID).Delete(&e).Error
	return err
}
