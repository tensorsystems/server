package server

import (
	"github.com/google/wire"
	"github.com/tensoremr/server/pkg/repository"
	"gorm.io/gorm"
)

func injectDb(db *gorm.DB) *gorm.DB {
	wire.Build(repository.ProvideAllergyRepository)
	wire.Build(repository.ProvideAmendmentRepository)
	wire.Build(repository.ProvideAppointmentQueueRepository)
	wire.Build(repository.ProvideAppointmentStatusRepository)
	wire.Build(repository.ProvideAppointmentRepository)
	wire.Build(repository.ProvideAutoRefractionRepository)
	wire.Build(repository.ProvideBillingRepository)
	wire.Build(repository.ProvideChatDeleteRepository)
	wire.Build(repository.ProvideChatMemberRepository)
	wire.Build(repository.ProvideChatMessageRepository)
	wire.Build(repository.ProvideChatMuteRepository)
	wire.Build(repository.ProvideChatUnreadRepository)
	wire.Build(repository.ProvideChatRepository)
	wire.Build(repository.ProvideChiefComplaintTypeRepository)
	wire.Build(repository.ProvideChiefComplaintRepository)
	wire.Build(repository.ProvideCoverTestRepository)
	wire.Build(repository.ProvideDiagnosisRepository)
	wire.Build(repository.ProvideExamCategoryRepository)
	wire.Build(repository.ProvideExternalExamRepository)
	wire.Build(repository.ProvideEyewearPrescriptionOrderRepository)
	wire.Build(repository.ProvideEyewearPrescriptionRepository)
	wire.Build(repository.ProvideEyewearShopRepository)
	wire.Build(repository.ProvideFamilyIllnessRepository)
	wire.Build(repository.ProvideFavoriteChiefComplaintRepository)
	wire.Build(repository.ProvideFavoriteDiagnosisRepository)
	wire.Build(repository.ProvideFavoriteMedicationRepository)
	wire.Build(repository.ProvideFileRepository)
	wire.Build(repository.ProvideFollowUpOrderRepository)
	wire.Build(repository.ProvideFollowUpRepository)
	wire.Build(repository.ProvideFunduscopyRepository)
	wire.Build(repository.ProvideHpiComponentTypeRepository)
	wire.Build(repository.ProvideHpiComponentRepository)
	wire.Build(repository.ProvideIopRepository)
	wire.Build(repository.ProvideLabOrderRepository)
	wire.Build(repository.ProvideLabTypeRepository)
	wire.Build(repository.ProvideLabRepository)
	wire.Build(repository.ProvideLifestyleTypeRepository)
	wire.Build(repository.ProvideLifestyleRepository)
	wire.Build(repository.ProvideMedicalPrescriptionOrderRepository)
	wire.Build(repository.ProvideMedicalPrescriptionRepository)
	wire.Build(repository.ProvideOcularMotilityRepository)
	wire.Build(repository.ProvideOpthalmologyExamRepository)
	wire.Build(repository.ProvideOpticDiscRepository)
	wire.Build(repository.ProvideOrganizationDetailsRepository)
	wire.Build(repository.ProvidePastHospitalizationRepository)
	wire.Build(repository.ProvidePastIllnessTypeRepository)
	wire.Build(repository.ProvidePastIllnessRepository)
	wire.Build(repository.ProvidePastInjuryRepository)
	wire.Build(repository.ProvidePastOptSurgeryRepository)
	wire.Build(repository.ProvidePastSurgeryRepository)
	wire.Build(repository.ProvidePatientChartRepository)
	wire.Build(repository.ProvidePatientDiagnosisRepository)
	wire.Build(repository.ProvidePatientEncounterLimitRepository)
	wire.Build(repository.ProvidePatientHistoryRepository)
	wire.Build(repository.ProvidePatientQueueRepository)
	wire.Build(repository.ProvidePatientRepository)
	wire.Build(repository.ProvidePaymentWaiverRepository)
	wire.Build(repository.ProvidePaymentRepository)
	wire.Build(repository.ProvidePharmacyRepository)
	wire.Build(repository.ProvidePhysicalExamFindingRepository)
	wire.Build(repository.ProvidePupilsRepository)
	wire.Build(repository.ProvideQueueDestinationRepository)
	wire.Build(repository.ProvideQueueSubscriptionRepository)
	wire.Build(repository.ProvideReferralOrderRepository)
	wire.Build(repository.ProvideReferralRepository)
	wire.Build(repository.ProvideReviewOfSystemRepository)
	wire.Build(repository.ProvideRoomRepository)
	wire.Build(repository.ProvideSlitLampExamRepository)
	wire.Build(repository.ProvideSupplyRepository)
	wire.Build(repository.ProvideSurgicalOrderRepository)
	wire.Build(repository.ProvideSurgicalProcedureTypeRepository)
	wire.Build(repository.ProvideSurgicalProcedureRepository)
	wire.Build(repository.ProvideSystemSymptomRepository)
	wire.Build(repository.ProvideSystemRepository)
	wire.Build(repository.ProvideTreatmentOrderRepository)
	wire.Build(repository.ProvideTreatmentTypeRepository)
	wire.Build(repository.ProvideTreatmentRepository)
	wire.Build(repository.ProvideUserTypeRepository)
	wire.Build(repository.ProvideUserRepository)
	wire.Build(repository.ProvideVisitTypeRepository)
	wire.Build(repository.ProvideVisualAcuityRepository)
	wire.Build(repository.ProvideVitalSignsRepository)

	return db
}