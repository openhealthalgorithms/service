package models

// CarePlanContentMapping object
type CarePlanContentMapping struct {
    Meta                    interface{}               `json:"meta"`
    CarePlanContentGoals    CarePlanContentGoals      `json:"goals"`
    CarePlanContentActivity CarePlanContentActivities `json:"activities"`
}

type CarePlanContentGoals map[string]CarePlanContentGoal
type CarePlanContentGoal map[string]interface{}

type CarePlanContentActivities map[string]CarePlanContentActivity
type CarePlanContentActivity map[string]interface{}
