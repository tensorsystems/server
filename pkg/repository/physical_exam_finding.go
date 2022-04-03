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

// PhysicalExamFinding ...
type PhysicalExamFinding struct {
	gorm.Model
	ID             int          `gorm:"primaryKey" json:"id"`
	PatientChartID int          `json:"patientChartId" gorm:"uniqueIndex"`
	ExamCategoryID int          `json:"examCategoryId"`
	ExamCategory   ExamCategory `json:"examCategory"`
	Abnormal       bool         `string:"abnormal"`
	Note           *string      `json:"note"`
	Count          int64        `json:"count"`
}

// Save ...
func (r *PhysicalExamFinding) Save() error {
	return DB.Create(&r).Error
}

// GetAll ...
func (r *PhysicalExamFinding) GetAll(p PaginationInput, filter *PhysicalExamFinding) ([]PhysicalExamFinding, int64, error) {
	var result []PhysicalExamFinding

	dbOp := DB.Scopes(Paginate(&p)).Select("*, count(*) OVER() AS count").Preload("ExamCategory").Where(filter).Order("id ASC").Find(&result)

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
func (r *PhysicalExamFinding) Get(ID int) error {
	return DB.Where("id = ?", ID).Take(&r).Error
}

// GetByTitle ...
func (r *PhysicalExamFinding) GetByPatientChartID(id string) error {
	return DB.Where("patient_chart_id = ?", id).Take(&r).Error
}

// Update ...
func (r *PhysicalExamFinding) Update() error {
	return DB.Save(&r).Error
}

// Delete ...
func (r *PhysicalExamFinding) Delete(ID int) error {
	return DB.Where("id = ?", ID).Delete(&r).Error
}

// DeleteExamCategory
func (r *PhysicalExamFinding) DeleteExamCategory(physicalExamFindingID int, examCategoryID int) error {
	return DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("id = ?", physicalExamFindingID).Take(&r).Error; err != nil {
			return err
		}

		var examCategory ExamCategory
		if err := tx.Where("id = ?", examCategoryID).Take(&examCategory).Error; err != nil {
			return err
		}

		tx.Model(&r).Where("id = ?", physicalExamFindingID).Association("ExamCategory").Delete(&examCategory)

		return nil
	})
}
