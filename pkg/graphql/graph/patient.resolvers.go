package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"time"

	graph_models "github.com/tensoremr/server/pkg/graphql/graph/model"
	"github.com/tensoremr/server/pkg/models"
	"github.com/tensoremr/server/pkg/repository"
	deepCopy "github.com/ulule/deepcopier"
)

func (r *mutationResolver) SavePatient(ctx context.Context, input graph_models.PatientInput) (*models.Patient, error) {
	// Copy
	var patient models.Patient
	deepCopy.Copy(&input).To(&patient)

	// Upload paper record document
	if input.PaperRecordDocument != nil {
		fileName, hashedFileName, hash, ext := HashFileName(input.PaperRecordDocument.Name)
		err := WriteFile(input.PaperRecordDocument.File.File, hashedFileName+"."+ext)
		if err != nil {
			return nil, err
		}

		patient.PaperRecordDocument = &models.File{
			ContentType: input.PaperRecordDocument.File.ContentType,
			Size:        input.PaperRecordDocument.File.Size,
			FileName:    fileName,
			Extension:   ext,
			Hash:        hash,
		}
	}

	// Upload other doucments
	for _, fileUpload := range input.Documents {
		fileName, hashedFileName, hash, ext := HashFileName(fileUpload.Name)
		err := WriteFile(fileUpload.File.File, hashedFileName+"."+ext)
		if err != nil {
			return nil, err
		}

		patient.Documents = append(patient.Documents, models.File{
			ContentType: fileUpload.File.ContentType,
			Size:        fileUpload.File.Size,
			FileName:    fileName,
			Extension:   ext,
			Hash:        hash,
		})
	}

	var patientRepository repository.PatientRepository

	// Save
	if err := patientRepository.Save(&patient); err != nil {
		return nil, err
	}

	// Return
	return &patient, nil
}

func (r *mutationResolver) SavePatientV2(ctx context.Context, input graph_models.PatientInputV2, dateOfBirthInput graph_models.DateOfBirthInput) (*models.Patient, error) {
	// Copy
	var patient models.Patient
	deepCopy.Copy(&input).To(&patient)

	var dateOfBirth time.Time

	now := time.Now()

	if dateOfBirthInput.InputType == graph_models.DateOfBirthInputTypeDate {
		dateOfBirth = *dateOfBirthInput.DateOfBirth
	} else if dateOfBirthInput.InputType == graph_models.DateOfBirthInputTypeAgeYear {
		dateOfBirth = now.AddDate(-*dateOfBirthInput.AgeInYears, 0, 0)
	} else if dateOfBirthInput.InputType == graph_models.DateOfBirthInputTypeAgeMonth {
		dateOfBirth = now.AddDate(0, -*dateOfBirthInput.AgeInMonths, 0)
	}

	patient.DateOfBirth = dateOfBirth

	// Upload paper record document
	if input.PaperRecordDocument != nil {
		fileName, hashedFileName, hash, ext := HashFileName(input.PaperRecordDocument.Name)
		err := WriteFile(input.PaperRecordDocument.File.File, hashedFileName+"."+ext)
		if err != nil {
			return nil, err
		}

		patient.PaperRecordDocument = &models.File{
			ContentType: input.PaperRecordDocument.File.ContentType,
			Size:        input.PaperRecordDocument.File.Size,
			FileName:    fileName,
			Extension:   ext,
			Hash:        hash,
		}
	}

	// Upload other doucments
	for _, fileUpload := range input.Documents {
		fileName, hashedFileName, hash, ext := HashFileName(fileUpload.Name)
		err := WriteFile(fileUpload.File.File, hashedFileName+"."+ext)
		if err != nil {
			return nil, err
		}

		patient.Documents = append(patient.Documents, models.File{
			ContentType: fileUpload.File.ContentType,
			Size:        fileUpload.File.Size,
			FileName:    fileName,
			Extension:   ext,
			Hash:        hash,
		})
	}

	var patientRepository repository.PatientRepository

	// Save
	if err := patientRepository.Save(&patient); err != nil {
		return nil, err
	}

	// Return
	return &patient, nil
}

func (r *mutationResolver) UpdatePatient(ctx context.Context, input graph_models.PatientUpdateInput) (*models.Patient, error) {
	var patient models.Patient
	deepCopy.Copy(&input).To(&patient)

	// Upload paper record document
	if input.PaperRecordDocument != nil {
		fileName, hashedFileName, hash, ext := HashFileName(input.PaperRecordDocument.Name)
		err := WriteFile(input.PaperRecordDocument.File.File, hashedFileName+"."+ext)
		if err != nil {
			return nil, err
		}

		patient.PaperRecordDocument = &models.File{
			ContentType: input.PaperRecordDocument.File.ContentType,
			Size:        input.PaperRecordDocument.File.Size,
			FileName:    fileName,
			Extension:   ext,
			Hash:        hash,
		}
	}

	// Upload other doucments
	for _, fileUpload := range input.Documents {
		fileName, hashedFileName, hash, ext := HashFileName(fileUpload.Name)
		err := WriteFile(fileUpload.File.File, hashedFileName+"."+ext)
		if err != nil {
			return nil, err
		}

		patient.Documents = append(patient.Documents, models.File{
			ContentType: fileUpload.File.ContentType,
			Size:        fileUpload.File.Size,
			FileName:    fileName,
			Extension:   ext,
			Hash:        hash,
		})
	}

	var patientRepository repository.PatientRepository
	if err := patientRepository.Update(&patient); err != nil {
		return nil, err
	}

	return &patient, nil
}

func (r *mutationResolver) DeletePatient(ctx context.Context, id int) (bool, error) {
	var repository repository.PatientRepository

	if err := repository.Delete(id); err != nil {
		return false, err
	}

	return true, nil
}

func (r *queryResolver) Patient(ctx context.Context, id int) (*models.Patient, error) {
	var patient models.Patient

	var repository repository.PatientRepository
	if err := repository.Get(&patient, id); err != nil {
		return nil, err
	}

	return &patient, nil
}

func (r *queryResolver) SearchPatients(ctx context.Context, term string) ([]*models.Patient, error) {
	var repository repository.PatientRepository
	patients, err := repository.Search(term)
	if err != nil {
		return nil, err
	}

	return patients, nil
}

func (r *queryResolver) GetByCardNo(ctx context.Context, cardNo string) (*models.Patient, error) {
	var patient models.Patient
	var repository repository.PatientRepository
	if err := repository.FindByCardNo(&patient, cardNo); err != nil {
		return nil, err
	}

	return &patient, nil
}

func (r *queryResolver) GetProgressNotes(ctx context.Context, appointmentID int) (*graph_models.ProgressNote, error) {
	var repository repository.PatientRepository

	patientHistory, appointments, err := repository.GetAllProgressNotes(appointmentID)
	if err != nil {
		return nil, err
	}

	return &graph_models.ProgressNote{
		PatientHistory: patientHistory,
		Appointments:   appointments,
	}, nil
}

func (r *queryResolver) GetAllPatientProgress(ctx context.Context, patientID int) (*graph_models.ProgressNote, error) {
	var repository repository.PatientRepository

	patientHistory, appointments, err := repository.GetAllProgress(patientID)
	if err != nil {
		return nil, err
	}

	return &graph_models.ProgressNote{
		PatientHistory: patientHistory,
		Appointments:   appointments,
	}, nil
}

func (r *queryResolver) GetVitalSignsProgress(ctx context.Context, patientID int) (*graph_models.VitalSignsProgress, error) {
	var repository repository.PatientRepository

	appointments, err := repository.GetVitalSignsProgress(patientID)
	if err != nil {
		return nil, err
	}

	return &graph_models.VitalSignsProgress{
		Appointments: appointments,
	}, nil
}

func (r *queryResolver) GetPatientDiagnosticProgress(ctx context.Context, patientID int, procedureTypeTitle string) ([]*models.Appointment, error) {
	var repository repository.PatientRepository

	appointments, err := repository.GetPatientDiagnosticProcedures(patientID, procedureTypeTitle)
	if err != nil {
		return nil, err
	}

	return appointments, nil
}

func (r *queryResolver) GetPatientDiagnosticProcedureTitles(ctx context.Context, patientID int) ([]string, error) {
	var repository repository.DiagnosticProcedureOrderRepository
	return repository.GetPatientDiagnosticProcedureTitles(patientID)
}

func (r *queryResolver) Patients(ctx context.Context, page models.PaginationInput) (*graph_models.PatientConnection, error) {
	var repository repository.PatientRepository
	patients, count, err := repository.GetAll(page)

	if err != nil {
		return nil, err
	}

	edges := make([]*graph_models.PatientEdge, len(patients))

	for i, entity := range patients {
		e := entity

		edges[i] = &graph_models.PatientEdge{
			Node: &e,
		}
	}

	pageInfo, totalCount := GetPageInfo(patients, count, page)
	return &graph_models.PatientConnection{PageInfo: pageInfo, Edges: edges, TotalCount: totalCount}, nil
}

func (r *queryResolver) GetPatientOrderCount(ctx context.Context, patientID int) (*graph_models.OrdersCount, error) {
	var diagnosticRepo repository.DiagnosticProcedureOrderRepository
	diagnosticCount, err := diagnosticRepo.GetCount(&models.DiagnosticProcedureOrder{PatientID: patientID}, nil, nil)
	if err != nil {
		return nil, err
	}

	var labRepo repository.LabOrderRepository
	labCount, err := labRepo.GetCount(&models.LabOrder{PatientID: patientID}, nil, nil)
	if err != nil {
		return nil, err
	}

	var surgeryRepo repository.SurgicalOrderRepository
	surgeryCount, err := surgeryRepo.GetCount(&models.SurgicalOrder{PatientID: patientID}, nil, nil)
	if err != nil {
		return nil, err
	}

	var treatmentRepo repository.TreatmentOrderRepository
	treatmentCount, err := treatmentRepo.GetCount(&models.TreatmentOrder{PatientID: patientID}, nil, nil)
	if err != nil {
		return nil, err
	}

	var followUpRepo repository.FollowUpOrderRepository
	followUpCount, err := followUpRepo.GetCount(&models.FollowUpOrder{PatientID: patientID}, nil, nil)
	if err != nil {
		return nil, err
	}

	var referralRepo repository.ReferralOrderRepository
	referralCount, err := referralRepo.GetCount(&models.ReferralOrder{PatientID: patientID}, nil, nil)
	if err != nil {
		return nil, err
	}

	return &graph_models.OrdersCount{
		DiagnosticProcedureOrders: int(diagnosticCount),
		LabOrders:                 int(labCount),
		SurgicalOrders:            int(surgeryCount),
		TreatmentOrders:           int(treatmentCount),
		FollowUpOrders:            int(followUpCount),
		ReferralOrders:            int(referralCount),
	}, nil
}

func (r *queryResolver) GetPatientFiles(ctx context.Context, patientID int) ([]*models.File, error) {
	var repository repository.PatientRepository
	files, err := repository.GetPatientFiles(patientID)
	return files, err
}

func (r *queryResolver) FindSimilarPatients(ctx context.Context, input graph_models.SimilarPatientsInput) (*graph_models.SimilarPatients, error) {
	var repository repository.PatientRepository

	byName, err := repository.FindByName(input.FirstName, input.LastName)
	if err != nil {
		return nil, err
	}

	byPhone, err := repository.FindByPhoneNo(input.PhoneNo)
	if err != nil {
		return nil, err
	}

	return &graph_models.SimilarPatients{
		ByName:  byName,
		ByPhone: byPhone,
	}, nil
}
