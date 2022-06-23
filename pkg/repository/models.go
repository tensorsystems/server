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
	"fmt"
	"os"
	"reflect"
	"strings"

	"github.com/casbin/casbin/v2"
	"github.com/tensoremr/server/pkg/conf"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// Global variable
var DB *gorm.DB

// Model facilitate database interactions
type Model struct {
	models map[string]reflect.Value
	isOpen bool
	*gorm.DB
	*casbin.Enforcer
}

// NewModel returns a new Model without opening database connection
func NewModel() *Model {
	return &Model{
		models: make(map[string]reflect.Value),
	}
}

// IsOpen returns true if the Model has already established connection
// to the database
func (m *Model) IsOpen() bool {
	return m.isOpen
}

// OpenPostgres ...
func (m *Model) OpenPostgres() error {
	dbHost := os.Getenv("DB_HOST")
	// dbReplicationHost := os.Getenv("DB_REPLICATION_HOST")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	dbPort := os.Getenv("DB_PORT")

	DBURL := fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=disable TimeZone=UTC password=%s", dbHost, dbPort, dbUser, dbName, dbPassword)
	db, err := gorm.Open(postgres.Open(DBURL), &gorm.Config{})

	if err != nil {
		return err
	}

	DB = db

	m.DB = db
	m.isOpen = true

	return nil
}

// OpenWithConfigTest opens database connection with test database
func (m *Model) OpenWithConfigTest(cfg *conf.Configuration) error {
	dbHost := os.Getenv("TEST_DB_HOST")
	dbUser := os.Getenv("TEST_DB_USER")
	dbPassword := os.Getenv("TEST_DB_PASSWORD")
	dbName := os.Getenv("TEST_DB_NAME")
	dbPort := os.Getenv("TEST_DB_PORT")

	DBURL := fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=disable TimeZone=UTC password=%s", dbHost, dbPort, dbUser, dbName, dbPassword)
	db, err := gorm.Open(postgres.Open(DBURL), &gorm.Config{})

	if err != nil {
		return err
	}

	DB = db

	m.DB = db
	m.isOpen = true

	return nil
}

// Register adds the values to the models registry
func (m *Model) Register(values ...interface{}) error {

	// do not work on them.models first, this is like an insurance policy
	// whenever we encounter any error in the values nothing goes into the registry
	models := make(map[string]reflect.Value)
	if len(values) > 0 {
		for _, val := range values {
			rVal := reflect.ValueOf(val)
			if rVal.Kind() == reflect.Ptr {
				rVal = rVal.Elem()
			}
			switch rVal.Kind() {
			case reflect.Struct:
				models[getTypeName(rVal.Type())] = reflect.New(rVal.Type())
			default:
				return errors.New("models must be structs")
			}
		}
	}
	for k, v := range models {
		m.models[k] = v
	}
	return nil
}

// AutoMigrateAll runs migrations for all the registered models
func (m *Model) AutoMigrateAll() {
	for _, v := range m.models {
		m.AutoMigrate(v.Interface())
	}
}

// DropAll drops all tables
func (m *Model) DropAll() {
	for _, v := range m.models {
		m.Migrator().DropTable(v.Interface())
	}
}

// RegisterAllModels ...
func (m *Model) RegisterAllModels() {
	m.Register(AppointmentStatus{})
	m.Register(Appointment{})
	m.Register(Billing{})
	m.Register(ChiefComplaint{})
	m.Register(PatientDiagnosis{})
	m.Register(DiagnosticProcedure{})
	m.Register(EyewearPrescription{})
	m.Register(FamilyIllness{})
	m.Register(File{})
	m.Register(Funduscopy{})
	m.Register(HpiComponentType{})
	m.Register(HpiComponent{})
	m.Register(Lab{})
	m.Register(Lifestyle{})
	m.Register(MedicalPrescription{})
	m.Register(Order{})
	m.Register(PastSurgery{})
	m.Register(PastHospitalization{})
	m.Register(PastIllness{})
	m.Register(PastInjury{})
	m.Register(PastOptSurgery{})
	m.Register(PatientChart{})
	m.Register(PatientHistory{})
	m.Register(Patient{})
	m.Register(Payment{})
	m.Register(Pupils{})
	m.Register(Room{})
	m.Register(SurgicalProcedure{})
	m.Register(Treatment{})
	m.Register(UserType{})
	m.Register(User{})
	m.Register(VisitType{})
	m.Register(AppointmentQueue{})
	m.Register(QueueDestination{})
	m.Register(PastIllnessType{})
	m.Register(LifestyleType{})
	m.Register(ChiefComplaintType{})
	m.Register(Diagnosis{})
	m.Register(DiagnosticProcedureType{})
	m.Register(FavoriteMedication{})
	m.Register(SurgicalProcedureType{})
	m.Register(Supply{})
	m.Register(TreatmentType{})
	m.Register(LabType{})
	m.Register(ChatDelete{})
	m.Register(ChatMember{})
	m.Register(ChatMessage{})
	m.Register(ChatMute{})
	m.Register(ChatUnreadMessage{})
	m.Register(Chat{})
	m.Register(Allergy{})
	m.Register(Referral{})
	m.Register(FavoriteChiefComplaint{})
	m.Register(FavoriteDiagnosis{})
	m.Register(PaymentWaiver{})
	m.Register(PatientEncounterLimit{})
	m.Register(Pharmacy{})
	m.Register(EyewearShop{})
	m.Register(MedicalPrescriptionOrder{})
	m.Register(EyewearPrescriptionOrder{})
	m.Register(MedicalPrescriptionOrder{})
	m.Register(DiagnosticProcedureOrder{})
	m.Register(LabOrder{})
	m.Register(Amendment{})
	m.Register(OrganizationDetails{})
	m.Register(System{})
	m.Register(SystemSymptom{})
	m.Register(ReviewOfSystem{})
	m.Register(ExamCategory{})
	m.Register(ExamFinding{})
	m.Register(PhysicalExamFinding{})
	m.Register(PatientQueue{})
	m.Register(QueueSubscription{})
	m.Register(OpthalmologyExam{})
	m.Register(VitalSigns{})
	m.Register(SurgicalOrder{})
	m.Register(TreatmentOrder{})
	m.Register(FollowUp{})
	m.Register(FollowUpOrder{})
	m.Register(ReferralOrder{})
}

// SeedData ...
func (m *Model) SeedData() error {
	var appointmentStatus AppointmentStatus
	appointmentStatus.Seed()

	var userType UserType
	userType.Seed()

	var adminUser User
	adminUser.Seed()

	var visitType VisitType
	visitType.Seed()
	
	// Save default patient encounter limits
	var user User
	physicians, err := user.GetByUserTypeTitle("Physician")

	if err != nil {
		return err
	}

	for _, p := range physicians {
		user := *p

		var patientEncounterLimit PatientEncounterLimit
		if err := patientEncounterLimit.GetByUser(user.ID); err != nil {
			patientEncounterLimit.UserID = user.ID
			patientEncounterLimit.MondayLimit = 150
			patientEncounterLimit.TuesdayLimit = 150
			patientEncounterLimit.WednesdayLimit = 150
			patientEncounterLimit.ThursdayLimit = 150
			patientEncounterLimit.FridayLimit = 150
			patientEncounterLimit.SaturdayLimit = 150
			patientEncounterLimit.SundayLimit = 150
			patientEncounterLimit.Overbook = 150

			if err := patientEncounterLimit.Save(); err != nil {
				return err
			}
		}
	}



	// var queueDestination QueueDestination
	// queueDestination.Seed()

	return nil
}

func getTypeName(typ reflect.Type) string {
	if typ.Name() != "" {
		return typ.Name()
	}
	split := strings.Split(typ.String(), ".")
	return split[len(split)-1]
}
