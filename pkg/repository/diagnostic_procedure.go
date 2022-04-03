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

// DiagnosticProcedureStatus ...
type DiagnosticProcedureStatus string

// DiagnosticProcedureOrder statuses ...
const (
	DiagnosticProcedureOrderedStatus   DiagnosticProcedureStatus = "ORDERED"
	DiagnosticProcedureCompletedStatus DiagnosticProcedureStatus = "COMPLETED"
)

// DiagnosticProcedure ...
type DiagnosticProcedure struct {
	gorm.Model
	ID                           int                       `gorm:"primaryKey"`
	DiagnosticProcedureOrderID   int                       `json:"diagnosticProcedureOrderId"`
	PatientChartID               int                       `json:"patientChartId"`
	GeneralText                  *string                   `json:"generalText"`
	Images                       []File                    `json:"images" gorm:"many2many:diagnostic_images"`
	Documents                    []File                    `json:"documents" gorm:"many2many:diagnostic_documents"`
	IsRefraction                 bool                      `json:"isRefraction"`
	RightDistanceSubjectiveSph   *string                   `json:"rightDistanceSubjectiveSph"`
	LeftDistanceSubjectiveSph    *string                   `json:"leftDistanceSubjectiveSph"`
	RightDistanceSubjectiveCyl   *string                   `json:"rightDistanceSubjectiveCyl"`
	LeftDistanceSubjectiveCyl    *string                   `json:"leftDistanceSubjectiveCyl"`
	RightDistanceSubjectiveAxis  *string                   `json:"rightDistanceSubjectiveAxis"`
	LeftDistanceSubjectiveAxis   *string                   `json:"leftDistanceSubjectiveAxis"`
	RightNearSubjectiveSph       *string                   `json:"rightNearSubjectiveSph"`
	LeftNearSubjectiveSph        *string                   `json:"leftNearSubjectiveSph"`
	RightNearSubjectiveCyl       *string                   `json:"rightNearSubjectiveCyl"`
	LeftNearSubjectiveCyl        *string                   `json:"leftNearSubjectiveCyl"`
	RightNearSubjectiveAxis      *string                   `json:"rightNearSubjectiveAxis"`
	LeftNearSubjectiveAxis       *string                   `json:"leftNearSubjectiveAxis"`
	RightDistanceObjectiveSph    *string                   `json:"rightDistanceObjectiveSph"`
	LeftDistanceObjectiveSph     *string                   `json:"leftDistanceObjectiveSph"`
	RightDistanceObjectiveCyl    *string                   `json:"rightDistanceObjectiveCyl"`
	LeftDistanceObjectiveCyl     *string                   `json:"leftDistanceObjectiveCyl"`
	RightDistanceObjectiveAxis   *string                   `json:"rightDistanceObjectiveAxis"`
	LeftDistanceObjectiveAxis    *string                   `json:"leftDistanceObjectiveAxis"`
	RightNearObjectiveSph        *string                   `json:"rightNearObjectiveSph"`
	LeftNearObjectiveSph         *string                   `json:"leftNearObjectiveSph"`
	RightNearObjectiveCyl        *string                   `json:"rightNearObjectiveCyl"`
	LeftNearObjectiveCyl         *string                   `json:"leftNearObjectiveCyl"`
	RightNearObjectiveAxis       *string                   `json:"rightNearObjectiveAxis"`
	LeftNearObjectiveAxis        *string                   `json:"leftNearObjectiveAxis"`
	RightDistanceFinalSph        *string                   `json:"rightDistanceFinalSph"`
	LeftDistanceFinalSph         *string                   `json:"leftDistanceFinalSph"`
	RightDistanceFinalCyl        *string                   `json:"rightDistanceFinalCyl"`
	LeftDistanceFinalCyl         *string                   `json:"leftDistanceFinalCyl"`
	RightDistanceFinalAxis       *string                   `json:"rightDistanceFinalAxis"`
	LeftDistanceFinalAxis        *string                   `json:"leftDistanceFinalAxis"`
	RightNearFinalSph            *string                   `json:"rightNearFinalSph"`
	LeftNearFinalSph             *string                   `json:"leftNearFinalSph"`
	RightNearFinalCyl            *string                   `json:"rightNearFinalCyl"`
	LeftNearFinalCyl             *string                   `json:"leftNearFinalCyl"`
	RightNearFinalAxis           *string                   `json:"rightNearFinalAxis"`
	LeftNearFinalAxis            *string                   `json:"leftNearFinalAxis"`
	RightVisualAcuity            *string                   `json:"rightVisualAcuity"`
	LeftVisualAcuity             *string                   `json:"leftVisualAcuity"`
	FarPd                        *string                   `json:"farPd"`
	NearPd                       *string                   `json:"nearPd"`
	DiagnosticProcedureTypeID    int                       `json:"diagnosticProcedureTypeId"`
	DiagnosticProcedureType      DiagnosticProcedureType   `json:"diagnosticProcedureType"`
	DiagnosticProcedureTypeTitle string                    `json:"diagnosticProcedureTypeTitle"`
	Payments                     []Payment                 `json:"payments" gorm:"many2many:diagnostic_procedure_payments;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Status                       DiagnosticProcedureStatus `json:"status"`
	OrderNote                    string                    `json:"orderNote"`
	ReceptionNote                string                    `json:"receptionNote"`
	Count                        int64                     `json:"count"`
	OldId                        int                       `json:"oldId"`
}

// Save ...
func (r *DiagnosticProcedure) Save() error {
	return DB.Create(&r).Error
}

// Get ...
func (r *DiagnosticProcedure) Get(ID int) error {
	return DB.Where("id = ?", ID).Take(&r).Error
}

// GetRefraction ...
func (r *DiagnosticProcedure) GetRefraction(patientChartID int) error {
	return DB.Transaction(func(tx *gorm.DB) error {
		var diagnosticProcedureType DiagnosticProcedureType
		if err := tx.Where("id = ?", "4").Take(&diagnosticProcedureType).Error; err != nil {
			return err
		}

		if err := tx.Where("is_refraction = ?", true).Where("patient_chart_id = ?", patientChartID).Take(&r).Error; err != nil {
			return err
		}

		return nil
	})
}

// GetAll ...
func (r *DiagnosticProcedure) GetAll(p PaginationInput, filter *DiagnosticProcedure) ([]DiagnosticProcedure, int64, error) {
	var result []DiagnosticProcedure

	dbOp := DB.Scopes(Paginate(&p)).Select("*, count(*) OVER() AS count").Where(filter).Preload("DiagnosticProcedureType").Preload("Images").Preload("RightEyeImages").Preload("LeftEyeImages").Preload("RightEyeSketches").Preload("LeftEyeSketches").Preload("Documents").Order("id ASC").Find(&result)

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
func (r *DiagnosticProcedure) Update() error {
	return DB.Updates(&r).Preload("Images").Preload("Documents").Error
}

// DeleteFile ...
func (r *DiagnosticProcedure) DeleteFile(association string, diagnosticProcedureID int, fileID int) error {
	return DB.Model(&DiagnosticProcedure{ID: diagnosticProcedureID}).Association(association).Delete(&File{ID: fileID})
}

// ClearAssociation ...
func (r *DiagnosticProcedure) ClearAssociation(association string, diagnosticProcedureID int) error {
	return DB.Model(&DiagnosticProcedure{ID: diagnosticProcedureID}).Association(association).Clear()
}

// Delete ...
func (r *DiagnosticProcedure) Delete(ID int) error {
	return DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("id = ?", ID).Take(&r).Error; err != nil {
			return err
		}

		if err := tx.Where("id = ?", ID).Delete(&r).Error; err != nil {
			return err
		}

		var diagnosticProceduresCount int64
		if err := tx.Model(&r).Where("diagnostic_procedure_order_id = ?", r.DiagnosticProcedureOrderID).Count(&diagnosticProceduresCount).Error; err != nil {
			return err
		}

		if diagnosticProceduresCount == 0 {
			if err := tx.Where("id = ?", r.DiagnosticProcedureOrderID).Delete(&DiagnosticProcedureOrder{}).Error; err != nil {
				return err
			}
		}

		return nil
	})
}

// GetByPatientChartID ...
func (r *DiagnosticProcedure) GetByPatientChartID(ID int) error {
	return DB.Where("patient_chart_id = ?", ID).Take(&r).Error
}
