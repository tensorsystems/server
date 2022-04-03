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
	"errors"
	"time"

	"github.com/lib/pq"
	"gorm.io/gorm"
)

// Patient ...
type Patient struct {
	gorm.Model
	ID                     int            `gorm:"primaryKey"`
	FirstName              string         `json:"firstName" gorm:"not null;"`
	LastName               string         `json:"lastName" gorm:"not null;"`
	FullName               string         `json:"fullName"`
	Gender                 string         `json:"gender"`
	PhoneNo                string         `json:"phoneNo" gorm:"size:100;not null;"`
	PhoneNo2               string         `json:"phoneNo2" gorm:"size:100;not null;"`
	TelNo                  string         `json:"telNo"`
	HomePhone              string         `json:"homePhone"`
	Email                  string         `json:"email" gorm:"size:100;not null;"`
	DateOfBirth            time.Time      `json:"dateOfBirth"`
	IDNo                   string         `json:"idNo"`
	IDType                 string         `json:"idType"`
	MartialStatus          string         `json:"martialStatus"`
	Occupation             string         `json:"occupation"`
	Credit                 *bool          `json:"credit"`
	CreditCompany          *string        `json:"creditCompany"`
	EmergencyContactName   string         `json:"emergencyContactName"`
	EmergencyContactRel    string         `json:"emergencyContactRel"`
	EmergencyContactPhone  string         `json:"emergencyContactPhone"`
	EmergencyContactPhone2 string         `json:"emergencyContactPhone2"`
	EmergencyContactMemo   string         `json:"emergencyContactMemo"`
	City                   string         `json:"city"`
	SubCity                string         `json:"subCity"`
	Region                 string         `json:"region"`
	Woreda                 string         `json:"woreda"`
	Zone                   string         `json:"zone"`
	Kebele                 string         `json:"kebele"`
	HouseNo                string         `json:"houseNo"`
	Memo                   string         `json:"memo"`
	CardNo                 string         `json:"cardNo"`
	PaperRecord            bool           `json:"paperRecord"`
	PaperRecordDocumentID  *int           `json:"paperRecordDocumentId"`
	PaperRecordDocument    *File          `json:"paperRecordDocument"`
	Documents              []File         `json:"documents" gorm:"many2many:patient_documents"`
	PatientHistory         PatientHistory `json:"patientHistory"`
	Appointments           []Appointment  `json:"appointments"`
	Document               string         `gorm:"type:tsvector"`
	Count                  int64          `json:"count"`
}

// AfterCreate ...
func (r *Patient) AfterCreate(tx *gorm.DB) error {
	r.FullName = r.FirstName + " " + r.LastName

	if err := tx.Model(r).Save(&r).Error; err != nil {
		return err
	}

	return nil
}

// AfterUpdate ...
func (r *Patient) AfterUpdate(tx *gorm.DB) (err error) {
	for _, appointment := range r.Appointments {
		appointment.FirstName = r.FirstName
		appointment.LastName = r.LastName
		appointment.PhoneNo = r.PhoneNo

		var provider User
		tx.Where("id = ?", appointment.UserID).Take(&provider)
		appointment.ProviderName = provider.FirstName + " " + provider.LastName

		tx.Model(appointment).Updates(&appointment)
	}

	return
}

// Save ...
func (r *Patient) Save() error {
	return DB.Transaction(func(tx *gorm.DB) error {
		var existingPatient Patient
		if err := tx.Where("trim(first_name) = ?", r.FirstName).Where("trim(last_name) = ?", r.LastName).Where("trim(phone_no) = ?", r.PhoneNo).Take(&existingPatient).Error; err == nil {
			if existingPatient.ID != 0 {
				return errors.New("Patient already exists")
			}
		}

		if err := tx.Create(&r).Error; err != nil {
			return err
		}

		if err := tx.Create(&PatientHistory{PatientID: r.ID}).Error; err != nil {
			return err
		}

		return nil
	})
}

// Get ...
func (r *Patient) Get(ID int) error {
	err := DB.Where("id = ?", ID).Preload("PatientHistory").Preload("PaperRecordDocument").Preload("Documents").Take(&r).Error
	if err != nil {
		return err
	}

	return nil
}

// GetPatientFiles ...
func (r *Patient) GetPatientFiles(patientID int) ([]*File, error) {
	var files []File

	err := DB.Transaction(func(tx *gorm.DB) error {
		var diagnosticProcedureOrder DiagnosticProcedureOrder
		tx.Model(&DiagnosticProcedureOrder{}).Where("patient_id = ?", patientID).Preload("DiagnosticProcedures.Images").Preload("DiagnosticProcedures.Documents").Take(&diagnosticProcedureOrder)

		for _, diagnosticProcedure := range diagnosticProcedureOrder.DiagnosticProcedures {
			files = append(files, diagnosticProcedure.Images...)
			files = append(files, diagnosticProcedure.Documents...)
		}

		var labOrder LabOrder
		tx.Model(&LabOrder{}).Where("patient_id = ?", patientID).Preload("Labs.Images").Preload("Labs.Documents").Take(&labOrder)

		for _, lab := range labOrder.Labs {
			files = append(files, lab.Images...)
			files = append(files, lab.Documents...)
		}

		return nil
	})

	var f []*File
	for _, file := range files {
		item := file
		f = append(f, &item)
	}

	return f, err
}

// GetAll ...
func (r *Patient) GetAll(p PaginationInput) ([]Patient, int64, error) {
	var result []Patient

	dbOp := DB.Scopes(Paginate(&p)).Select("*, count(*) OVER() AS count").Order("id ASC").Find(&result)

	var count int64
	if len(result) > 0 {
		count = result[0].Count
	}

	if dbOp.Error != nil {
		return result, 0, dbOp.Error
	}

	return result, count, dbOp.Error
}

// GetAllProgressNotes ...
func (r *Patient) GetAllProgressNotes(appointmentID int) (*PatientHistory, []*Appointment, error) {
	var patientHistory *PatientHistory
	var appointments []*Appointment

	err := DB.Transaction(func(tx *gorm.DB) error {
		var appointment Appointment

		if err := tx.Where("id = ?", appointmentID).Take(&appointment).Error; err != nil {
			return err
		}

		if err := DB.Model(Appointment{}).Where("patient_id = ?", appointment.PatientID).Where("id != ?", appointmentID).Preload("Patient").Preload("VisitType").Preload("PatientChart.VitalSigns").Preload("PatientChart.Diagnoses").Preload("PatientChart.MedicalPrescriptionOrder.MedicalPrescriptions").Preload("PatientChart.LabOrder.Labs").Preload("PatientChart.DiagnosticProcedureOrder.DiagnosticProcedures").Preload("PatientChart.SurgicalOrder.SurgicalProcedures").Preload("PatientChart.TreatmentOrder.Treatments").Preload("PatientChart.ReferralOrder.Referrals").Preload("PatientChart.FollowUpOrder.FollowUps").Preload("PatientChart.SurgicalProcedure").Preload("PatientChart.Treatment").Order("id ASC").Find(&appointments).Error; err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, nil, err
	}

	return patientHistory, appointments, nil
}

// GetAllProgress ...
func (r *Patient) GetAllProgress(patientID int) (*PatientHistory, []*Appointment, error) {
	var patientHistory *PatientHistory
	var appointments []*Appointment

	err := DB.Transaction(func(tx *gorm.DB) error {

		if err := DB.Model(Appointment{}).Where("patient_id = ?", patientID).Preload("Patient").Preload("VisitType").Preload("PatientChart.VitalSigns").Preload("PatientChart.Diagnoses").Preload("PatientChart.MedicalPrescriptionOrder.MedicalPrescriptions").Preload("PatientChart.LabOrder.Labs").Preload("PatientChart.DiagnosticProcedureOrder.DiagnosticProcedures").Preload("PatientChart.SurgicalOrder.SurgicalProcedures").Preload("PatientChart.TreatmentOrder.Treatments").Preload("PatientChart.ReferralOrder.Referrals").Preload("PatientChart.FollowUpOrder.FollowUps").Preload("PatientChart.SurgicalProcedure").Preload("PatientChart.Treatment").Order("id ASC").Find(&appointments).Error; err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, nil, err
	}

	return patientHistory, appointments, nil
}

// GetVisionProgress ...
func (r *Patient) GetVitalSignsProgress(patientID int) ([]*Appointment, error) {
	var appointments []*Appointment

	if err := DB.Model(Appointment{}).Where("patient_id = ?", patientID).Order("id ASC").Preload("PatientChart.VitalSigns").Find(&appointments).Error; err != nil {
		return nil, err
	}

	return appointments, nil
}

// GetPatientDiagnosticProcedures ...
func (r *Patient) GetPatientDiagnosticProcedures(patientID int, diagnosticProcedureTypeTitle string) ([]*Appointment, error) {
	var appointments []*Appointment

	if err := DB.Model(Appointment{}).Where("patient_id = ?", patientID).Order("id ASC").Preload("PatientChart.DiagnosticProcedureOrder.DiagnosticProcedures", "diagnostic_procedure_type_title ILIKE ?", diagnosticProcedureTypeTitle).Preload("PatientChart.DiagnosticProcedureOrder.DiagnosticProcedures.Images").Preload("PatientChart.DiagnosticProcedureOrder.DiagnosticProcedures.Documents").Preload("PatientChart.DiagnosticProcedureOrder.DiagnosticProcedures.Payments").Find(&appointments).Error; err != nil {
		return nil, err
	}

	return appointments, nil
}

// Search ...
func (r *Patient) Search(term string) ([]*Patient, error) {
	var patients []*Patient
	err := DB.Raw("SELECT * FROM patients WHERE document @@ plainto_tsquery(?) AND deleted_at IS NULL LIMIT 20", term).Find(&patients).Error
	return patients, err
}

// FindByCardNo ...
func (r *Patient) FindByCardNo(cardNo string) error {
	err := DB.Where("card_no = ?", cardNo).Take(&r).Error
	if err != nil {
		return err
	}

	return nil
}

// FindByName ...
func (r *Patient) FindByName(firstName string, lastName string) ([]*Patient, error) {
	var patients []*Patient
	err := DB.Where("trim(first_name) ILIKE ?", firstName).Where("trim(last_name) ILIKE ?", lastName).Find(&patients).Error
	if err != nil {
		return patients, err
	}

	return patients, nil
}

// FindByPhoneNo ...
func (r *Patient) FindByPhoneNo(phoneNo string) ([]*Patient, error) {
	var patients []*Patient

	err := DB.Where("trim(phone_no) ILIKE ?", phoneNo).Find(&patients).Error
	if err != nil {
		return patients, err
	}

	return patients, nil
}

// Update ...
func (r *Patient) Update() error {
	err := DB.Updates(r).Error

	if err, ok := err.(*pq.Error); ok && err.Code.Name() == "unique_violation" {
		return errors.New("Duplicate, " + err.Detail)
	}

	return err
}

// Delete ...
func (r *Patient) Delete(ID int) error {
	var e Patient
	err := DB.Where("id = ?", ID).Delete(&e).Error
	if err != nil {
		return err
	}
	return nil
}

// Clean ...
func (r *Patient) Clean() error {
	return DB.Transaction(func(tx *gorm.DB) error {
		patients := []map[string]interface{}{}

		if err := tx.Raw("SELECT first_name, last_name, phone_no, count(*) FROM patients GROUP BY first_name, last_name, phone_no HAVING count(*) > 1").Find(&patients).Error; err != nil {
			return err
		}

		for _, e := range patients {
			var duplicatePatients []Patient

			if err := tx.Select("id, first_name, last_name, phone_no").Where("first_name = ?", e["first_name"].(string)).Where("last_name = ?", e["last_name"].(string)).Where("phone_no = ?", e["phone_no"].(string)).Preload("Appointments").Preload("PatientHistory.PastIllnesses").Preload("PatientHistory.PastInjuries").Preload("PatientHistory.PastHospitalizations").Preload("PatientHistory.PastSurgeries").Preload("PatientHistory.FamilyIllnesses").Preload("PatientHistory.Lifestyles").Preload("PatientHistory.Allergies").Preload("PatientHistory.PastHospitalizations").Order("id desc").Find(&duplicatePatients).Error; err != nil {
				return err
			}

			primaryPatient := duplicatePatients[0]

			// Attach paper record document to primary record
			if primaryPatient.PaperRecordDocumentID == nil {
				for _, e := range duplicatePatients {
					if e.PaperRecordDocumentID != nil {
						primaryPatient.PaperRecordDocumentID = e.PaperRecordDocumentID
						primaryPatient.CardNo = e.CardNo
					}
				}
			}

			// Update medical prescriptions
			for index, e := range duplicatePatients {
				if index != 0 {
					tx.Model(&MedicalPrescription{}).Where("patient_id = ?", e.ID).Update("patient_id", primaryPatient.ID)
				}
			}

			// Update eyewear prescriptions
			for index, e := range duplicatePatients {
				if index != 0 {
					tx.Model(&EyewearPrescription{}).Where("patient_id = ?", e.ID).Update("patient_id", primaryPatient.ID)
				}
			}

			// Update diagnostic orders
			for index, e := range duplicatePatients {
				if index != 0 {
					tx.Model(&DiagnosticProcedureOrder{}).Where("patient_id = ?", e.ID).Update("patient_id", primaryPatient.ID)
				}
			}

			// Update surgical orders
			for index, e := range duplicatePatients {
				if index != 0 {
					tx.Model(&SurgicalOrder{}).Where("patient_id = ?", e.ID).Update("patient_id", primaryPatient.ID)
				}
			}

			// Update lab
			for index, e := range duplicatePatients {
				if index != 0 {
					tx.Model(&LabOrder{}).Where("patient_id = ?", e.ID).Update("patient_id", primaryPatient.ID)
				}
			}

			// Update treatments
			for index, e := range duplicatePatients {
				if index != 0 {
					tx.Model(&TreatmentOrder{}).Where("patient_id = ?", e.ID).Update("patient_id", primaryPatient.ID)
				}
			}

			// Update follow-up
			for index, e := range duplicatePatients {
				if index != 0 {
					tx.Model(&FollowUpOrder{}).Where("patient_id = ?", e.ID).Update("patient_id", primaryPatient.ID)
				}
			}

			// Update referral
			for index, e := range duplicatePatients {
				if index != 0 {
					tx.Model(&ReferralOrder{}).Where("patient_id = ?", e.ID).Update("patient_id", primaryPatient.ID)
				}
			}

			// Attach histories
			for _, e := range duplicatePatients {
				for _, p := range e.PatientHistory.PastIllnesses {
					tx.Model(&PastIllness{}).Where("id = ?", p.ID).Update("patient_history_id", primaryPatient.PatientHistory.ID)
				}

				for _, p := range e.PatientHistory.PastHospitalizations {
					tx.Model(&PastHospitalization{}).Where("id = ?", p.ID).Update("patient_history_id", primaryPatient.PatientHistory.ID)
				}

				for _, p := range e.PatientHistory.FamilyIllnesses {
					tx.Model(&FamilyIllness{}).Where("id = ?", p.ID).Update("patient_history_id", primaryPatient.PatientHistory.ID)
				}

				for _, p := range e.PatientHistory.Lifestyles {
					tx.Model(&Lifestyle{}).Where("id = ?", p.ID).Update("patient_history_id", primaryPatient.PatientHistory.ID)
				}

				for _, p := range e.PatientHistory.Allergies {
					tx.Model(&Allergy{}).Where("id = ?", p.ID).Update("patient_history_id", primaryPatient.PatientHistory.ID)
				}
			}

			// Attach appointments
			for _, e := range duplicatePatients {
				for _, a := range e.Appointments {
					tx.Model(&Appointment{}).Where("id = ?", a.ID).Update("patient_id", primaryPatient.ID)
				}
			}

			// Delete non primary records
			for index, e := range duplicatePatients {
				if index != 0 {
					tx.Delete(&e)
				}
			}

			tx.Updates(&primaryPatient)
		}

		return nil
	})

}
