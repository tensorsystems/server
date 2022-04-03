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

// FollowUpStatus ...
type FollowUpStatus string

// SurgicalProcedureOrder statuses ...
const (
	FollowUpStatusOrdered   FollowUpStatus = "ORDERED"
	FollowUpStatusCompleted FollowUpStatus = "COMPLETED"
)

// FollowUp ...
type FollowUp struct {
	gorm.Model
	ID              int            `gorm:"primaryKey"`
	FollowUpOrderID int            `json:"followUpOrderId"`
	PatientChartID  int            `json:"patientChartId"`
	Status          FollowUpStatus `json:"status"`
	ReceptionNote   string         `json:"receptionNote"`
	Count           int64          `json:"count"`
}

// Save ...
func (r *FollowUp) Save() error {
	return DB.Create(&r).Error
}

// Get ...
func (r *FollowUp) Get(ID int) error {
	return DB.Where("id = ?", ID).Take(&r).Error
}

// GetByPatientChart ...
func (r *FollowUp) GetByPatientChart(ID int) error {
	return DB.Where("patient_chart_id = ?", ID).Take(&r).Error
}

// GetAll ...
func (r *FollowUp) GetAll(p PaginationInput, filter *FollowUp) ([]FollowUp, int64, error) {
	var result []FollowUp

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
func (r *FollowUp) Update() error {
	return DB.Updates(&r).Error
}

// Delete ...
func (r *FollowUp) Delete(ID int) error {
	return DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("id = ?", ID).Take(&r).Error; err != nil {
			return err
		}

		if err := tx.Where("id = ?", ID).Delete(&r).Error; err != nil {
			return err
		}

		var followUpsCount int64
		if err := tx.Model(&r).Where("follow_up_order_id = ?", r.FollowUpOrderID).Count(&followUpsCount).Error; err != nil {
			return err
		}

		if followUpsCount == 0 {
			if err := tx.Where("id = ?", r.FollowUpOrderID).Delete(&FollowUpOrder{}).Error; err != nil {
				return err
			}
		}

		return nil
	})
}
