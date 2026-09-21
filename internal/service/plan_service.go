package service

import (
	"github.com/iamrichmon/subscription-api/internal/model"
)

type PlanService struct {
	plans []model.Plan
}

func NewPlanService(plans []model.Plan) *PlanService {
	return &PlanService{
		plans: plans,
	}
}

func (s *PlanService) ListPlans() []model.Plan {
	return s.plans
}

func (s *PlanService) GetPlanByID(id model.SubscriptionStatus) (model.Plan, bool) {
	for _, plan := range s.plans {
		if plan.ID == id {
			return plan, true
		}
	}
	return model.Plan{}, false
}
