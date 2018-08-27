package assessments

func GetExercise(current int) (ExerciseAssessment, error) {
	return GetExerciseWithTarget(current, 150)
}

func GetExerciseWithTarget(current, target int) (ExerciseAssessment, error) {
	targetTime := "150 minutes"
	resultCode := "PA-2"

	if current > target {
		resultCode = "PA-1"
	}

	exerciseObj := NewExerciseAssessment(current, resultCode, targetTime)

	return exerciseObj, nil
}

type ExerciseAssessment struct {
	Current int    `structs:"value"`
	Code    string `structs:"code"`
	Target  string `structs:"target"`
}

// NewExerciseAssessment returns a BP object.
func NewExerciseAssessment(current int, code, target string) ExerciseAssessment {
	return ExerciseAssessment{
		Current: current,
		Code:    code,
		Target:  target,
	}
}
