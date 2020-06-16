package models

// CarePlanConditionsMapping object
type CarePlanConditionsMapping struct {
	CarePlanConditionMeta    interface{}                `json:"meta"`
	CarePlanConditionMapping []CarePlanConditionMapping `json:"mappings"`
}

type (

	// CarePlanConditionMapping object
	CarePlanConditionMapping struct {
		CarePlanConditions []CarePlanCondition `json:"conditions"`
		CarePlanGoals      []CarePlanGoal      `json:"goals"`
		CarePlanActivities []CarePlanActivity  `json:"activities"`
	}

	// CarePlanCondition object
	CarePlanCondition map[string][]string

	// CarePlanGoal object
	CarePlanGoal string

	// CarePlanActivity object
	CarePlanActivity map[string][]CarePlanRule

	// CarePlanRule object
	CarePlanRule map[string]string
)
