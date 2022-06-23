package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"

	"github.com/tensoremr/server/pkg/graphql/graph/model"
	"github.com/tensoremr/server/pkg/repository"
	deepCopy "github.com/ulule/deepcopier"
)

func (r *mutationResolver) SavePatient(ctx context.Context, input model.PatientInput) (*repository.Patient, error) {
	// Copy
	var patient repository.Patient
	deepCopy.Copy(&input).To(&patient)

	// Upload paper record document
	if input.PaperRecordDocument != nil {
		fileName, hashedFileName, hash, ext := HashFileName(input.PaperRecordDocument.Name)
		err := WriteFile(input.PaperRecordDocument.File.File, hashedFileName+"."+ext)
		if err != nil {
			return nil, err
		}

		patient.PaperRecordDocument = &repository.File{
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

		patient.Documents = append(patient.Documents, repository.File{
			ContentType: fileUpload.File.ContentType,
			Size:        fileUpload.File.Size,
			FileName:    fileName,
			Extension:   ext,
			Hash:        hash,
		})
	}

	// Save
	if err := patient.Save(); err != nil {
		return nil, err
	}

	// Return
	return &patient, nil
}

func (r *mutationResolver) UpdatePatient(ctx context.Context, input model.PatientUpdateInput) (*repository.Patient, error) {
	var patient repository.Patient
	deepCopy.Copy(&input).To(&patient)

	// Upload paper record document
	if input.PaperRecordDocument != nil {
		fileName, hashedFileName, hash, ext := HashFileName(input.PaperRecordDocument.Name)
		err := WriteFile(input.PaperRecordDocument.File.File, hashedFileName+"."+ext)
		if err != nil {
			return nil, err
		}

		patient.PaperRecordDocument = &repository.File{
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

		patient.Documents = append(patient.Documents, repository.File{
			ContentType: fileUpload.File.ContentType,
			Size:        fileUpload.File.Size,
			FileName:    fileName,
			Extension:   ext,
			Hash:        hash,
		})
	}

	if err := patient.Update(); err != nil {
		return nil, err
	}

	return &patient, nil
}

func (r *mutationResolver) DeletePatient(ctx context.Context, id int) (bool, error) {
	var patient repository.Patient

	err := patient.Delete(id)
	if err != nil {
		return false, err
	}

	return true, nil
}

func (r *queryResolver) Patient(ctx context.Context, id int) (*repository.Patient, error) {
	var patient repository.Patient

	err := patient.Get(id)
	if err != nil {
		return nil, err
	}

	return &patient, err
}

func (r *queryResolver) SearchPatients(ctx context.Context, term string) ([]*repository.Patient, error) {
	var patient repository.Patient

	patients, err := patient.Search(term)
	if err != nil {
		return nil, err
	}

	return patients, nil
}

func (r *queryResolver) GetByCardNo(ctx context.Context, cardNo string) (*repository.Patient, error) {
	var patient repository.Patient

	if err := patient.FindByCardNo(cardNo); err != nil {
		return nil, err
	}

	return &patient, nil
}

func (r *queryResolver) GetProgressNotes(ctx context.Context, appointmentID int) (*model.ProgressNote, error) {
	var patient repository.Patient

	patientHistory, appointments, err := patient.GetAllProgressNotes(appointmentID)
	if err != nil {
		return nil, err
	}

	return &model.ProgressNote{
		PatientHistory: patientHistory,
		Appointments:   appointments,
	}, nil
}

func (r *queryResolver) GetAllPatientProgress(ctx context.Context, patientID int) (*model.ProgressNote, error) {
	var patient repository.Patient

	patientHistory, appointments, err := patient.GetAllProgress(patientID)
	if err != nil {
		return nil, err
	}

	return &model.ProgressNote{
		PatientHistory: patientHistory,
		Appointments:   appointments,
	}, nil
}

func (r *queryResolver) GetVitalSignsProgress(ctx context.Context, patientID int) (*model.VitalSignsProgress, error) {
	var repository repository.Patient

	appointments, err := repository.GetVitalSignsProgress(patientID)
	if err != nil {
		return nil, err
	}

	return &model.VitalSignsProgress{
		Appointments: appointments,
	}, nil
}

func (r *queryResolver) GetPatientDiagnosticProgress(ctx context.Context, patientID int, procedureTypeTitle string) ([]*repository.Appointment, error) {
	var repo repository.Patient

	appointments, err := repo.GetPatientDiagnosticProcedures(patientID, procedureTypeTitle)
	if err != nil {
		return nil, err
	}

	return appointments, nil
}

func (r *queryResolver) GetPatientDiagnosticProcedureTitles(ctx context.Context, patientID int) ([]string, error) {
	var repository repository.DiagnosticProcedureOrder
	return repository.GetPatientDiagnosticProcedureTitles(patientID)
}

func (r *queryResolver) Patients(ctx context.Context, page repository.PaginationInput) (*model.PatientConnection, error) {
	var patient repository.Patient
	patients, count, err := patient.GetAll(page)

	if err != nil {
		return nil, err
	}

	edges := make([]*model.PatientEdge, len(patients))

	for i, entity := range patients {
		e := entity

		edges[i] = &model.PatientEdge{
			Node: &e,
		}
	}

	pageInfo, totalCount := GetPageInfo(patients, count, page)
	return &model.PatientConnection{PageInfo: pageInfo, Edges: edges, TotalCount: totalCount}, nil
}

func (r *queryResolver) GetPatientOrderCount(ctx context.Context, patientID int) (*model.OrdersCount, error) {
	var diagnosticRepo repository.DiagnosticProcedureOrder
	diagnosticCount, err := diagnosticRepo.GetCount(&repository.DiagnosticProcedureOrder{PatientID: patientID}, nil, nil)
	if err != nil {
		return nil, err
	}

	var labRepo repository.LabOrder
	labCount, err := labRepo.GetCount(&repository.LabOrder{PatientID: patientID}, nil, nil)
	if err != nil {
		return nil, err
	}

	var surgeryRepo repository.SurgicalOrder
	surgeryCount, err := surgeryRepo.GetCount(&repository.SurgicalOrder{PatientID: patientID}, nil, nil)
	if err != nil {
		return nil, err
	}

	var treatmentRepo repository.TreatmentOrder
	treatmentCount, err := treatmentRepo.GetCount(&repository.TreatmentOrder{PatientID: patientID}, nil, nil)
	if err != nil {
		return nil, err
	}

	var followUpRepo repository.FollowUpOrder
	followUpCount, err := followUpRepo.GetCount(&repository.FollowUpOrder{PatientID: patientID}, nil, nil)
	if err != nil {
		return nil, err
	}

	var referralRepo repository.ReferralOrder
	referralCount, err := referralRepo.GetCount(&repository.ReferralOrder{PatientID: patientID}, nil, nil)
	if err != nil {
		return nil, err
	}

	return &model.OrdersCount{
		DiagnosticProcedureOrders: int(diagnosticCount),
		LabOrders:                 int(labCount),
		SurgicalOrders:            int(surgeryCount),
		TreatmentOrders:           int(treatmentCount),
		FollowUpOrders:            int(followUpCount),
		ReferralOrders:            int(referralCount),
	}, nil
}

func (r *queryResolver) GetPatientFiles(ctx context.Context, patientID int) ([]*repository.File, error) {
	var patient repository.Patient
	files, err := patient.GetPatientFiles(patientID)
	return files, err
}

func (r *queryResolver) FindSimilarPatients(ctx context.Context, input model.SimilarPatientsInput) (*model.SimilarPatients, error) {
	var repository repository.Patient

	byName, err := repository.FindByName(input.FirstName, input.LastName)
	if err != nil {
		return nil, err
	}

	byPhone, err := repository.FindByPhoneNo(input.PhoneNo)
	if err != nil {
		return nil, err
	}

	return &model.SimilarPatients{
		ByName:  byName,
		ByPhone: byPhone,
	}, nil
}

// !!! WARNING !!!
// The code below was going to be deleted when updating resolvers. It has been copied here so you have
// one last chance to move it out of harms way if you want. There are two reasons this happens:
//  - When renaming or deleting a resolver the old code will be put in here. You can safely delete
//    it when you're done.
//  - You have helper methods in this file. Move them out to keep these resolver files clean.
func (r *queryResolver) GetDiagnosticPatientProgress(ctx context.Context, patientID int, procedureTypeTitle string) ([]*repository.Appointment, error) {
	panic(fmt.Errorf("not implemented"))
}
