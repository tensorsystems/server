package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	"github.com/Tensor-Systems/tensoremr-server/pkg/graphql/graph/model"
	"github.com/Tensor-Systems/tensoremr-server/pkg/repository"
)

func (r *queryResolver) Notifs(ctx context.Context) (*model.Notif, error) {
	var diagnosticOrderRepo repository.DiagnosticProcedureOrder
	diagnostics := diagnosticOrderRepo.GetTodaysOrderedCount()

	var surgicalOrderRepo repository.SurgicalOrder
	surgical := surgicalOrderRepo.GetTodaysOrderedCount()

	var labOrderRepo repository.LabOrder
	labs := labOrderRepo.GetTodaysOrderedCount()

	var treatmentRepo repository.TreatmentOrder
	treatments := treatmentRepo.GetTodaysOrderedCount()

	var followUpRepo repository.FollowUpOrder
	followUps := followUpRepo.GetTodaysOrderedCount()

	var referralRepo repository.ReferralOrder
	referrals := referralRepo.GetTodaysOrderedCount()

	var paymentWaiver repository.PaymentWaiver
	waivers, err := paymentWaiver.GetApprovedCount()

	if err != nil {
		return nil, err
	}

	return &model.Notif{
		DiagnosticProcedureOrders: diagnostics,
		LabOrders:                 labs,
		TreatmentOrders:           treatments,
		SurgicalOrders:            surgical,
		ReferralOrders:            referrals,
		FollowUpOrders:            followUps,
		PaymentWaivers:            waivers,
	}, nil
}
