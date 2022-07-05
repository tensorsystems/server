package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	graph_models "github.com/tensoremr/server/pkg/graphql/graph/model"
	"github.com/tensoremr/server/pkg/repository"
)

func (r *queryResolver) Notifs(ctx context.Context) (*graph_models.Notif, error) {
	var diagnosticOrderRepo repository.DiagnosticProcedureOrderRepository
	diagnostics := diagnosticOrderRepo.GetTodaysOrderedCount()

	var surgicalOrderRepo repository.SurgicalOrderRepository
	surgical := surgicalOrderRepo.GetTodaysOrderedCount()

	var labOrderRepo repository.LabOrderRepository
	labs := labOrderRepo.GetTodaysOrderedCount()

	var treatmentRepo repository.TreatmentOrderRepository
	treatments := treatmentRepo.GetTodaysOrderedCount()

	var followUpRepo repository.FollowUpOrderRepository
	followUps := followUpRepo.GetTodaysOrderedCount()

	var referralRepo repository.ReferralOrderRepository
	referrals := referralRepo.GetTodaysOrderedCount()

	var paymentWaiver repository.PaymentWaiverRepository
	waivers, err := paymentWaiver.GetApprovedCount()

	if err != nil {
		return nil, err
	}

	return &graph_models.Notif{
		DiagnosticProcedureOrders: diagnostics,
		LabOrders:                 labs,
		TreatmentOrders:           treatments,
		SurgicalOrders:            surgical,
		ReferralOrders:            referrals,
		FollowUpOrders:            followUps,
		PaymentWaivers:            waivers,
	}, nil
}
