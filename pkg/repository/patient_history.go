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

// PatientHistory ...
type PatientHistory struct {
	gorm.Model
	ID                   int                   `gorm:"primaryKey" json:"id"`
	PatientID            int                   `json:"patientId" gorm:"uniqueIndex"`
	ReviewOfSystems      []ReviewOfSystem      `json:"reviewOfSystems"`
	ReviewOfSystemsNote  *string               `json:"reviewOfSystemsNote"`
	PastIllnesses        []PastIllness         `json:"pastIllnesses"`
	PastInjuries         []PastInjury          `json:"pastInjuries"`
	PastHospitalizations []PastHospitalization `json:"pastHospitalizations"`
	PastSurgeries        []PastSurgery         `json:"pastSurgeries"`
	FamilyIllnesses      []FamilyIllness       `json:"familyIllnesses"`
	Lifestyles           []Lifestyle           `json:"lifestyles"`
	Allergies            []Allergy             `json:"allergies"`
}

// Save ...
func (r *PatientHistory) Save() error {
	err := DB.Create(&r).Error
	if err != nil {
		return err
	}

	return nil
}

// Get ...
func (r *PatientHistory) Get(ID int) error {
	return DB.Where("id = ?", ID).Take(&r).Error
}

// GetByPatientID ...
func (r *PatientHistory) GetByPatientID(ID int) error {
	return DB.Where("patient_id = ?", ID).Take(&r).Error
}

// GetByPatientIDWithDetails ...
func (r *PatientHistory) GetByPatientIDWithDetails(ID int) error {
	return DB.Where("patient_id = ?", ID).Preload("PastIllnesses").Preload("PastInjuries").Preload("PastHospitalizations").Preload("PastSurgeries").Preload("FamilyIllnesses").Preload("Lifestyles").Preload("Allergies").Take(&r).Error
}

// Update ...
func (r *PatientHistory) Update() error {
	return DB.Updates(&r).Error
}

// Delete ...
func (r *PatientHistory) Delete(ID int) error {
	var e PatientHistory
	err := DB.Where("id = ?", ID).Delete(&e).Error
	return err
}
