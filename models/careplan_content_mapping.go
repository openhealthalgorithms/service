package models

// CarePlanContentMapping object
type CarePlanContentMapping struct {
	Meta                    interface{}               `json:"meta"`
	CarePlanContentGoals    CarePlanContentGoals      `json:"goals"`
	CarePlanContentActivity CarePlanContentActivities `json:"activities"`
}

// CarePlanContentGoals object
type CarePlanContentGoals map[string]CarePlanContentGoal

// CarePlanContentGoal object
type CarePlanContentGoal map[string]interface{}

// CarePlanContentActivities object
type CarePlanContentActivities map[string]CarePlanContentActivity

// CarePlanContentActivity object
type CarePlanContentActivity map[string]interface{}
