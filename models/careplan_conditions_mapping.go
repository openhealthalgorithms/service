package models

// CarePlanConditionsMapping object
type CarePlanConditionsMapping struct {
	CarePlanConditionMeta    interface{}                `json:"meta"`
	CarePlanConditionMapping []CarePlanConditionMapping `json:"mappings"`
}

type (
	CarePlanConditionMapping struct {
		CarePlanConditions []CarePlanCondition `json:"conditions"`
		CarePlanGoals      []CarePlanGoal      `json:"goals"`
		CarePlanActivities []CarePlanActivity  `json:"activities"`
	}

	CarePlanCondition map[string][]string
	CarePlanGoal      string
	CarePlanActivity  map[string][]CarePlanRule
	CarePlanRule      map[string]string
)
