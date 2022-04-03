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

// VisitType ...
type VisitType struct {
	gorm.Model
	ID    int    `gorm:"primaryKey"`
	Title string `json:"title" gorm:"unique"`
}

// Seed ...
func (r *VisitType) Seed() {
	DB.Create(&VisitType{Title: "Sick visit"})
	DB.Create(&VisitType{Title: "Follow up"})
	DB.Create(&VisitType{Title: "Check up"})
	DB.Create(&VisitType{Title: "Surgery"})
	DB.Create(&VisitType{Title: "Treatment"})
	DB.Create(&VisitType{Title: "Post-Op"})
	DB.Create(&VisitType{Title: "Referral"})
}

// Save ...
func (r *VisitType) Save() error {
	err := DB.Create(&r).Error

	if err, ok := err.(*pq.Error); ok && err.Code.Name() == "unique_violation" {
		var existing VisitType
		existingErr := DB.Unscoped().Where("title = ?", r.Title).Take(&existing).Error

		if existingErr == nil {
			DB.Model(&VisitType{}).Unscoped().Where("id = ?", existing.ID).Update("deleted_at", nil)
			r = &existing
			return nil
		}

		return errors.New("Duplicate, " + err.Detail)
	}

	if err != nil {
		return err
	}

	return nil
}

// Get ...
func (r *VisitType) Get(ID int) error {
	err := DB.Where("id = ?", ID).Take(&r).Error
	if err != nil {
		return err
	}

	return nil
}

// GetByTitle ...
func (r *VisitType) GetByTitle(title string) error {
	err := DB.Where("title = ?", title).Take(&r).Error
	if err != nil {
		return err
	}

	return nil
}

// GetByTitles ...
func (r *VisitType) GetByTitles(titles []string) ([]VisitType, error) {
	var visitTypes []VisitType
	err := DB.Where("title IN ?", titles).Find(&visitTypes).Error
	return visitTypes, err
}

// Count ...
func (r *VisitType) Count(dbString string) (int64, error) {
	var count int64

	err := DB.Model(&VisitType{}).Count(&count).Error
	return count, err
}

// GetAll ...
func (r *VisitType) GetAll(p PaginationInput) ([]VisitType, int64, error) {
	var result []VisitType

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

// Update ...
func (r *VisitType) Update() error {
	return DB.Updates(&r).Error
}

// Delete ...
func (r *VisitType) Delete(ID int) error {
	var e VisitType
	err := DB.Where("id = ?", ID).Delete(&e).Error
	return err
}
