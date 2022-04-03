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

// CoverTest ...
type CoverTest struct {
	gorm.Model
	ID             int     `gorm:"primaryKey"`
	RightCoverTest *string `json:"rightCoverTest"`
	LeftCoverTest  *string `json:"leftCoverTest"`
	Note           *string `json:"note"`
	PatientChartID int     `json:"patientChartId"`
}

// Save ...
func (r *CoverTest) Save() error {
	err := DB.Create(&r).Error
	if err != nil {
		return err
	}

	return nil
}

// SaveForPatientChart ...
func (r *CoverTest) SaveForPatientChart() error {
	return DB.Transaction(func(tx *gorm.DB) error {
		var existing CoverTest
		existingErr := tx.Where("patient_chart_id = ?", r.PatientChartID).Take(&existing).Error

		if existingErr != nil {
			if err := tx.Create(&r).Error; err != nil {
				return err
			}
		} else {
			r.ID = existing.ID
			if err := tx.Updates(&r).Error; err != nil {
				return err
			}
		}

		return nil
	})
}

// Get ...
func (r *CoverTest) Get(filter CoverTest) error {
	return DB.Where(filter).Take(&r).Error
}

// GetByPatientChart ...
func (r *CoverTest) GetByPatientChart(ID int) error {
	return DB.Where("patient_chart_id = ?", ID).Take(&r).Error
}

// Update ...
func (r *CoverTest) Update() error {
	return DB.Updates(&r).Error
}
