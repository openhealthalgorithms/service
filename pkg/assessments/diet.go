package assessments

func GetDiet(c_fruit, c_veg int) (DietAssessment, error) {
	return GetDietWithTarget(c_fruit, c_veg, 2, 5)
}

func GetDietWithTarget(c_fruit, c_veg, t_fruit, t_veg int) (DietAssessment, error) {
	resultCode := "NUT-1"

	if c_fruit < t_fruit && c_veg < t_veg {
		resultCode = "NUT-3"
	} else if (c_fruit < t_fruit && c_veg >= t_veg) || (c_fruit >= t_fruit && c_veg < t_veg) {
		resultCode = "NUT-2"
	}

	v := Values{c_fruit, c_veg}

	dietObj := NewDietAssessment(v, resultCode)

	return dietObj, nil
}

type Values struct {
	Fruit      int `structs:"fruit"`
	Vegetables int `structs:"vegetables"`
}

type DietAssessment struct {
	Values `structs:"value"`
	Code   string `structs:"code"`
}

// NewDietAssessment returns a BP object.
func NewDietAssessment(current Values, code string) DietAssessment {
	return DietAssessment{
		Values: current,
		Code:   code,
	}
}
