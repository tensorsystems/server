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

// PatientDiagnosis ...
type PatientDiagnosis struct {
	gorm.Model
	ID                     int     `gorm:"primaryKey"`
	PatientChartID         int     `json:"patientChartId"`
	CategoryCode           *string `json:"categoryCode"`
	DiagnosisCode          *string `json:"diagnosisCode"`
	FullCode               *string `json:"fullCode"`
	AbbreviatedDescription *string `json:"abbreviatedDescription"`
	FullDescription        string  `json:"fullDescription"`
	CategoryTitle          *string `json:"categoryTitle"`
	Location               string  `json:"location"`
	Differential           bool    `json:"differential"`
	Count                  int64   `json:"count"`
}

// Save ...
func (r *PatientDiagnosis) Save(diagnosisID int) error {
	return DB.Transaction(func(tx *gorm.DB) error {
		var diagnosis Diagnosis
		if err := tx.Where("id = ?", diagnosisID).Take(&diagnosis).Error; err != nil {
			return err
		}

		r.CategoryCode = diagnosis.CategoryCode
		r.DiagnosisCode = diagnosis.DiagnosisCode
		r.FullCode = diagnosis.FullCode
		r.AbbreviatedDescription = diagnosis.AbbreviatedDescription
		r.FullDescription = diagnosis.FullDescription

		if err := tx.Create(&r).Error; err != nil {
			return err
		}

		return nil
	})
}

// GetByPatientChartID ...
func (r *PatientDiagnosis) GetByPatientChartID(ID int) error {
	return DB.Where("patient_chart_id = ?", ID).Take(&r).Error
}

// Get ...
func (r *PatientDiagnosis) Get(ID int) error {
	err := DB.Where("id = ?", ID).Take(&r).Error
	if err != nil {
		return err
	}

	return nil
}

// GetAll ...
func (r *PatientDiagnosis) GetAll(p PaginationInput, filter *PatientDiagnosis) ([]PatientDiagnosis, int64, error) {
	var result []PatientDiagnosis

	dbOp := DB.Scopes(Paginate(&p)).Select("*, count(*) OVER() AS count").Where(filter).Order("id ASC").Find(&result)

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
func (r *PatientDiagnosis) Update() error {
	return DB.Updates(&r).Error
}

// Delete ...
func (r *PatientDiagnosis) Delete(ID int) error {
	var e PatientDiagnosis
	err := DB.Where("id = ?", ID).Delete(&e).Error
	return err
}
