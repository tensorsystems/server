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

// ChiefComplaint ...
type ChiefComplaint struct {
	gorm.Model
	ID             int            `gorm:"primaryKey"`
	Title          string         `json:"title"`
	PatientChartID int            `json:"patientChartId"`
	HPIComponents  []HpiComponent `gorm:"many2many:complaint_hpi_components"`
	Count          int64          `json:"count"`
	OldId          int            `json:"oldId"`
}

// Save ...
func (r *ChiefComplaint) Save() error {
	var existing ChiefComplaint
	DB.Unscoped().Where("title = ?", r.Title).Where("patient_chart_id = ?", r.PatientChartID).Take(&existing)

	if existing.ID != 0 {
		DB.Model(&ChiefComplaint{}).Unscoped().Where("id = ?", existing.ID).Update("deleted_at", nil)
		r = &existing
		return nil
	}

	err := DB.Create(&r).Error
	return err
}

// Get ...
func (r *ChiefComplaint) Get(ID int) error {
	err := DB.Where("id = ?", ID).Preload("HPIComponents").Take(&r).Error
	if err != nil {
		return err
	}

	return nil
}

// GetAll ...
func (r *ChiefComplaint) GetAll(p PaginationInput, filter *ChiefComplaint) ([]ChiefComplaint, int64, error) {
	var result []ChiefComplaint

	dbOp := DB.Scopes(Paginate(&p)).Select("*, count(*) OVER() AS count").Where(filter).Preload("HPIComponents").Order("id ASC").Find(&result)

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
func (r *ChiefComplaint) Search(p PaginationInput, searchTerm string) ([]ChiefComplaint, int64, error) {
	var result []ChiefComplaint
	dbOp := DB.Scopes(Paginate(&p)).Select("*, count(*) OVER() AS count").Where("title LIKE ?", "%"+searchTerm+"%").Order("id ASC").Preload("HPIComponents").Find(&result)

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
func (r *ChiefComplaint) Update() error {
	return DB.Transaction(func(tx *gorm.DB) error {
		tx.Model(&r).Association("HPIComponents").Replace(&r.HPIComponents)

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
func (r *ChiefComplaint) Delete(ID int) error {
	return DB.Where("id = ?", ID).Delete(&r).Error
}
