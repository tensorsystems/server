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

// PatientChart ...
type PatientChart struct {
	gorm.Model
	ID                        int                      `gorm:"primaryKey"`
	AppointmentID             int                      `json:"appointmentId"`
	VitalSigns                VitalSigns               `json:"vitalSigns"`
	PhysicalExamFindings      []PhysicalExamFinding    `json:"physicalExamFindings"`
	PhysicalExamFindingNote   *string                  `json:"physicalExamFindingNote"`
	OpthalmologyExam          OpthalmologyExam         `json:"opthalmologyExam"`
	SurgicalProcedure         SurgicalProcedure        `json:"surgicalProcedure"`
	Treatment                 Treatment                `json:"treatment"`
	ChiefComplaints           []ChiefComplaint         `json:"chiefComplaints"`
	ChiefComplaintsNote       *string                  `json:"chiefComplaintNote"`
	BloodPressure             *string                  `json:"bloodPressure"`
	HpiNote                   *string                  `json:"hpiNote"`
	DiagnosisNote             *string                  `json:"diagnosisNote"`
	DifferentialDiagnosisNote *string                  `json:"differentialDiagnosisNote"`
	RightSummarySketch        *string                  `json:"rightSummarySketch"`
	LeftSummarySketch         *string                  `json:"leftSummarySketch"`
	SummaryNote               *string                  `json:"summaryNote"`
	StickieNote               *string                  `json:"stickieNote"`
	MedicalRecommendation     *string                  `json:"medicalRecommendation"`
	SickLeave                 *string                  `json:"sickLeave"`
	MedicalPrescriptionOrder  MedicalPrescriptionOrder `json:"medicalPrescriptionOrder"`
	EyewearPrescriptionOrder  EyewearPrescriptionOrder `json:"eyewearPrescriptionOrder"`
	DiagnosticProcedureOrder  DiagnosticProcedureOrder `json:"diagnosticProcedureOrder"`
	SurgicalOrder             SurgicalOrder            `json:"surgicalOrder"`
	TreatmentOrder            TreatmentOrder           `json:"treatmentOrder"`
	ReferralOrder             ReferralOrder            `json:"referralOrder"`
	FollowUpOrder             FollowUpOrder            `json:"followUpOrder"`
	LabOrder                  LabOrder                 `json:"labOrder"`
	Diagnoses                 []PatientDiagnosis       `json:"diagnoses"`
	Locked                    *bool                    `json:"locked"`
	LockedDate                *time.Time               `json:"lockedDate"`
	LockedByID                *int                     `json:"lockedById"`
	LockedBy                  *User                    `json:"lockedBy"`
	Amendments                []Amendment              `json:"amendments"`
	OldPatientChartId         int                      `json:"oldPatientChartId"`
}

// AfterCreate ...
func (r *PatientChart) AfterCreate(tx *gorm.DB) error {

	if err := tx.Create(&VitalSigns{PatientChartID: r.ID}).Error; err != nil {
		return err
	}

	if err := tx.Create(&OpthalmologyExam{PatientChartID: r.ID}).Error; err != nil {
		return err
	}

	return nil
}

// Save ...
func (r *PatientChart) Save() error {
	return DB.Create(&r).Error
}

// SignAndLock ...
func (r *PatientChart) SignAndLock(patientChartID int, userID *int) error {
	return DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("id = ?", patientChartID).Take(&r).Error; err != nil {
			return err
		}

		if r.Locked != nil && *r.Locked == true {
			return nil
		}

		locked := true
		lockedDate := time.Now()

		r.Locked = &locked
		r.LockedDate = &lockedDate
		r.LockedByID = userID

		if err := tx.Updates(&r).Error; err != nil {
			return err
		}

		var checkedOut AppointmentStatus
		if err := tx.Where("title = ?", "Checked-Out").Take(&checkedOut).Error; err != nil {
			return err
		}

		if err := tx.Table("appointments").Where("id = ?", r.AppointmentID).Updates(map[string]interface{}{"appointment_status_id": checkedOut.ID}).Error; err != nil {
			return err
		}

		return nil
	})

}

// Get ...
func (r *PatientChart) Get(ID int) error {
	return DB.Where("id = ?", ID).Take(&r).Error
}

// GetByAppointmentID ...
func (r *PatientChart) GetByAppointmentID(appointmentID int) error {
	return DB.Where("appointment_id = ?", appointmentID).Take(&r).Error
}

// Get ...
func (r *PatientChart) GetWithDetails(ID int) error {
	return DB.Where("id = ?", ID).Preload("ChiefComplaints.HPIComponents.HpiComponentType").Preload("MedicalPrescriptionOrder.MedicalPrescriptions").Preload("EyewearPrescriptionOrder.EyewearPrescriptions").Preload("VitalSigns").Preload("PhysicalExamFindings.ExamCategory").Preload("OpthalmologyExam").Preload("Diagnoses").Preload("LabOrder.Labs.LabType").Preload("LabOrder.Labs.RightEyeImages").Preload("LabOrder.Labs.LeftEyeImages").Preload("LabOrder.Labs.Documents").Preload("DiagnosticProcedureOrder.DiagnosticProcedures.DiagnosticProcedureType").Preload("DiagnosticProcedureOrder.DiagnosticProcedures.Images").Preload("DiagnosticProcedureOrder.DiagnosticProcedures.Documents").Preload("SurgicalProcedure.SurgicalProcedureType").Preload("Treatment.TreatmentType").Preload("Amendments").Take(&r).Error
}

// Update ...
func (r *PatientChart) Update() error {
	return DB.Updates(&r).Error
}
