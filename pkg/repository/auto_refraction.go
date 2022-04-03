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

// AutoRefraction ...
type AutoRefraction struct {
	gorm.Model
	ID                 int    `gorm:"primaryKey"`
	RightDistanceSph   *string `json:"rightDistanceSph"`
	LeftDistanceSph    *string `json:"leftDistanceSph"`
	RightDistanceAxis  *string `json:"rightDistanceAxis"`
	LeftDistanceAxis   *string `json:"leftDistanceAxis"`
	RightDistanceCyl   *string `json:"rightDistanceCyl"`
	LeftDistanceCyl    *string `json:"leftDistanceCyl"`
	RightNearSph       *string `json:"rightNearSph"`
	LeftNearSph        *string `json:"leftNearSph"`
	RightNearCyl       *string `json:"rightNearCyl"`
	LeftNearCyl        *string `json:"leftNearCyl"`
	RightNearAxis      *string `json:"rightNearAxis"`
	LeftNearAxis       *string `json:"leftNearAxis"`
	RightLensMeterSph  *string `json:"rightLensMeterSph"`
	LeftLensMeterSph   *string `json:"leftLensMeterSph"`
	RightLensMeterAxis *string `json:"rightLensMeterAxis"`
	LeftLensMeterAxis  *string `json:"leftLensMeterAxis"`
	RightLensMeterCyl  *string `json:"rightLensMeterCyl"`
	LeftLensMeterCyl   *string `json:"leftLensMeterCyl"`
	PatientChartID     int    `json:"patientChartId" gorm:"unique"`
}

// Save ...
func (r *AutoRefraction) Save() error {
	return DB.Create(&r).Error
}

// SaveForPatientChart ...
func (r *AutoRefraction) SaveForPatientChart() error {
	return DB.Transaction(func(tx *gorm.DB) error {
		var existing AutoRefraction
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
func (r *AutoRefraction) Get(filter AutoRefraction) error {
	return DB.Where(filter).Take(&r).Error
}

// GetByPatientChart ...
func (r *AutoRefraction) GetByPatientChart(ID int) error {
	return DB.Where("patient_chart_id = ?", ID).Take(&r).Error
}

// Update ...
func (r *AutoRefraction) Update() error {
	return DB.Updates(&r).Error
}
