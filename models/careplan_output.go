package models

// CarePlanOutput object
type CarePlanOutput struct {
	CarePlanOutputGoals      []CarePlanContentGoal     `json:"goals"`
	CarePlanOutputActivities []CarePlanContentActivity `json:"activities"`
}
