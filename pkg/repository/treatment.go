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

// TreatmentStatus ...
type TreatmentStatus string

// SurgicalProcedureOrder statuses ...
const (
	TreatmentStatusOrdered   TreatmentStatus = "ORDERED"
	TreatmentStatusCompleted TreatmentStatus = "COMPLETED"
)

// Treatment ...
type Treatment struct {
	gorm.Model
	ID                 int             `gorm:"primaryKey"`
	TreatmentOrderID   int             `json:"treatmentOrderId"`
	PatientChartID     int             `json:"patientChartId"`
	Note               string          `json:"note"`
	Result             string          `json:"result"`
	RightEyeText       string          `json:"rightEyeText"`
	LeftEyeText        string          `json:"leftEyeText"`
	GeneralText        string          `json:"generalText"`
	TreatmentTypeID    int             `json:"treatmentTypeId"`
	TreatmentType      TreatmentType   `json:"treatmentType"`
	TreatmentTypeTitle string          `json:"treatmentTypeTitle"`
	Payments           []Payment       `json:"payments" gorm:"many2many:treatment_payments;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	ReceptionNote      string          `json:"receptionNote"`
	Status             TreatmentStatus `json:"status"`
	Count              int64           `json:"count"`
}

// Save ...
func (r *Treatment) Save() error {
	return DB.Create(&r).Error
}

// Get ...
func (r *Treatment) Get(ID int) error {
	return DB.Where("id = ?", ID).Take(&r).Error
}

// GetByPatientChart ...
func (r *Treatment) GetByPatientChart(ID int) error {
	return DB.Where("patient_chart_id = ?", ID).Preload("TreatmentType").Take(&r).Error
}

// GetAll ...
func (r *Treatment) GetAll(p PaginationInput, filter *Treatment) ([]Treatment, int64, error) {
	var result []Treatment

	dbOp := DB.Scopes(Paginate(&p)).Select("*, count(*) OVER() AS count").Where(filter).Preload("Order").Preload("Order.Payments").Preload("TreatmentType").Order("id ASC").Find(&result)

	var count int64
	if len(result) > 0 {
		count = result[0].Count
	}

	if dbOp.Error != nil {
		return result, 0, dbOp.Error
	}

	return result, count, dbOp.Error
}

// GetByPatient ...
func (r *Treatment) GetByPatient(p PaginationInput, patientID int) ([]Treatment, int64, error) {
	var result []Treatment

	dbOp := DB.Scopes(Paginate(&p)).Joins("INNER JOIN orders ON orders.id = treatments.order_id").Where("orders.patient_id = ?", patientID).Preload("Order").Preload("Order.Payments").Preload("Order.User").Preload("TreatmentType").Order("treatments.id ASC").Find(&result)

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
func (r *Treatment) Update() error {
	return DB.Updates(&r).Error
}

// Delete ...
func (r *Treatment) Delete(ID int) error {
	return DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("id = ?", ID).Take(&r).Error; err != nil {
			return err
		}

		if err := tx.Where("id = ?", ID).Delete(&r).Error; err != nil {
			return err
		}

		var treatmentsCount int64
		if err := tx.Model(&r).Where("treatment_order_id = ?", r.TreatmentOrderID).Count(&treatmentsCount).Error; err != nil {
			return err
		}

		if treatmentsCount == 0 {
			if err := tx.Where("id = ?", r.TreatmentOrderID).Delete(&TreatmentOrder{}).Error; err != nil {
				return err
			}
		}

		return nil
	})
}
