package lesson2

import (
	uuid "github.com/satori/go.uuid"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type Controller struct {
	PlanID uuid.UUID
	Plan *Plan // This is the fixed recipe for the lesson
	Model Model // This is the dynamically evaluated model as questions are answered etc.

	conn *gorm.DB
	InChannel <-chan ControllerEvent
	OutChannel chan<- EventLoopEvent
}

func NewController(planID uuid.UUID) *Controller {
	//TODO deffo need some channels.
	//TODO need to make the gorm connection
	return &Controller{
		planID,
		nil,
		Model{},
		nil,
	}
}

func (ctrl *Controller) LoadPlan() error {
	zap.L().Info("Loading Plan")
	result:= ctrl.conn.First(ctrl.Plan, ctrl.PlanID)
	if result.Error != nil {
		ctrl.Plan = nil // Reset required to make PlanIsLoaded work
		return result.Error
	}
	return nil
}
