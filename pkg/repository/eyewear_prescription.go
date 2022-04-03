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
	"time"

	"gorm.io/gorm"
)

// EyewearPrescription ...
type EyewearPrescription struct {
	gorm.Model
	ID                         int        `gorm:"primaryKey"`
	EyewearPrescriptionOrderID int        `json:"eyewearPrescriptionOrderId"`
	PatientID                  int        `json:"patientId"`
	Patient                    Patient    `json:"patient"`
	Glass                      *bool      `json:"glass"`
	Plastic                    *bool      `json:"plastic"`
	SingleVision               *bool      `json:"singleVision"`
	PhotoChromatic             *bool      `json:"photoChromatic"`
	GlareFree                  *bool      `json:"glareFree"`
	ScratchResistant           *bool      `json:"scratchResistant"`
	Bifocal                    *bool      `json:"bifocal"`
	Progressive                *bool      `json:"progressive"`
	TwoSeparateGlasses         *bool      `json:"twoSeparateGlasses"`
	HighIndex                  *bool      `json:"highIndex"`
	Tint                       *bool      `json:"tint"`
	BlueCut                    *bool      `json:"blueCut"`
	PrescribedDate             *time.Time `json:"prescribedDate"`
	History                    bool       `json:"history"`
	Status                     string     `json:"status"`
	Count                      int64      `json:"count"`
}

// Save ...
func (r *EyewearPrescription) Save() error {
	return DB.Create(&r).Error
}

// Get ...
func (r *EyewearPrescription) Get(ID int) error {
	return DB.Where("id = ?", ID).Take(&r).Error
}

// GetAll ...
func (r *EyewearPrescription) GetAll(p PaginationInput, filter *EyewearPrescription) ([]EyewearPrescription, int64, error) {
	var result []EyewearPrescription

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
func (r *EyewearPrescription) Update() error {
	return DB.Updates(&r).Error
}

// Delete ...
func (r *EyewearPrescription) Delete(ID int) error {
	return DB.Where("id = ?", ID).Delete(&r).Error
}
