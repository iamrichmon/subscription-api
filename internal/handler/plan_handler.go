package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/iamrichmon/subscription-api/internal/model"
	"github.com/iamrichmon/subscription-api/internal/service"
	"github.com/iamrichmon/subscription-api/internal/utils"
)

type PlanHandler struct {
	planService *service.PlanService
}

func NewPlanHandler(planService *service.PlanService) *PlanHandler {
	return &PlanHandler{
		planService: planService,
	}
}

func (h *PlanHandler) ListPlans(c *gin.Context) {
	plans := h.planService.ListPlans()

	c.JSON(http.StatusOK, plans)
}

func (h *PlanHandler) GetPlanByID(c *gin.Context) {
	id := c.Param("id")

	plan, found := h.planService.GetPlanByID(model.SubscriptionStatus(id))

	if !found {
		c.JSON(http.StatusNotFound, gin.H{"error": utils.ErrPlanNotFound.Error()})
		return
	}

	c.JSON(http.StatusOK, plan)
}
