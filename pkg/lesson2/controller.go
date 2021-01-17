package lesson2

import (
	uuid "github.com/satori/go.uuid"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// Controller runs on the EventLoop go routine mostly,
// except a number of Async methods that take a channel to sync back with (with optional error!)
type Controller struct {
	//Note: Controller doesn't own the current round any more - event loop does that
	PlanID uuid.UUID
	Plan *Plan // This is the fixed recipe for the lesson
	Model Model // This is the dynamically evaluated model as questions are answered etc.

	conn *gorm.DB
}

func NewController(planID uuid.UUID) *Controller {
	//TODO need to make the gorm connection some 'ow
	//TODO and it makes sense to make this here...?
	return &Controller{
		planID,
		nil,
		Model{},
		nil, //TODO make this here...
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
