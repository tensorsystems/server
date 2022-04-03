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

// UserType ...
type UserType struct {
	gorm.Model
	ID    int    `gorm:"primaryKey"`
	Title string `json:"title" gorm:"uniqueIndex"`
	Users []User `gorm:"many2many:user_type_roles;"`
	Count int64  `json:"count"`
}

// Seed ...
func (r *UserType) Seed() {
	DB.Create(&UserType{Title: "Admin"})
	DB.Create(&UserType{Title: "Nurse"})
	DB.Create(&UserType{Title: "Pharmacist"})
	DB.Create(&UserType{Title: "Optical Assistant"})
	DB.Create(&UserType{Title: "Physician"})
	DB.Create(&UserType{Title: "Optometrist"})
	DB.Create(&UserType{Title: "Receptionist"})
}

// Save ...
func (r *UserType) Save() error {
	err := DB.Create(&r).Error
	if err != nil {
		return err
	}

	return nil
}

// GetAll ...
func (r *UserType) GetAll(p PaginationInput) ([]UserType, int64, error) {
	var result []UserType

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
func (r *UserType) Get(ID int) error {
	err := DB.Where("id = ?", ID).Take(&r).Error
	if err != nil {
		return err
	}

	return nil
}

// Update ...
func (r *UserType) Update() (*UserType, error) {
	err := DB.Save(&r).Error
	if err != nil {
		return nil, err
	}
	return r, nil
}

// Delete ...
func (r *UserType) Delete(ID int) error {
	var e UserType
	err := DB.Where("id = ?", ID).Delete(&e).Error
	return err
}

// GetByTitle ...
func (r *UserType) GetByTitle(title string) error {
	return DB.Where("title = ?", title).Take(&r).Error
}

// GetByIds
func (r *UserType) GetByIds(ids []*int) ([]UserType, error) {
	var result []UserType
	err := DB.Where("id IN ?", ids).Find(&result).Error
	return result, err
}
