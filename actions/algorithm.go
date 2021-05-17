package actions

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"

	// imported for pgsql support
	_ "github.com/lib/pq"

	"github.com/openhealthalgorithms/service/algorithms"
	"github.com/openhealthalgorithms/service/config"
	"github.com/openhealthalgorithms/service/database"
	a "github.com/openhealthalgorithms/service/engines/assessments"
	"github.com/openhealthalgorithms/service/models"
	"github.com/openhealthalgorithms/service/tools"
)

var (
	dbFile          = filepath.Join(tools.GetCurrentDirectory(), "logs.db")
	sqlite          *database.SqliteDb
	currentSettings config.Settings
)

func load(c echo.Context) error {
	var err error
	currentSettings = c.Get("current_config").(config.Settings)
	dbFile = currentSettings.LogFile
	sqlite, err = database.InitDb(dbFile)
	if err != nil {
		return err
	}
	return nil
}

// AlgorithmHandler function
func AlgorithmHandler(c echo.Context) error {
	var err error
	err = load(c)
	if err != nil {
		return ErrorResponse(c, err, 500)
	}
	o := new(models.OHARequest)
	if err = c.Bind(o); err != nil {
		return ErrorResponse(c, err, 500)
	}
	if err = c.Validate(o); err != nil {
		return ErrorResponse(c, err, 400)
	}
	if *o.Config.Algorithm != "hearts" {
		return ErrorResponse(c, errors.New("algorithm not found"), 404)
	}

	if o.Config.RiskModelVersion == nil {
		rmVersion := "who_ish_2007"
		o.Config.RiskModelVersion = &rmVersion
	}

	if o.Config.LabBased == nil {
		lab := false
		o.Config.LabBased = &lab
	}

	colorChartPath := filepath.Join(currentSettings.ColorChart, *o.Config.RiskModel, *o.Config.RiskModelVersion, "charts.json")
	if _, e := os.Stat(colorChartPath); e != nil {
		return ErrorResponse(c, e, 400)
	}

	countriesPath := filepath.Join(currentSettings.ColorChart, *o.Config.RiskModel, *o.Config.RiskModelVersion, "countries.json")
	if _, e := os.Stat(countriesPath); e != nil {
		return ErrorResponse(c, e, 400)
	}

	guideFiles, err := tools.ParseGuidesFiles(c)
	if err != nil {
		return ErrorResponse(c, err, 500)
	}

	gd, gdc, gl, glc, cp, cpc, err := processGuides(guideFiles)
	if err != nil {
		return ErrorResponse(c, err, 500)
	}

	hearts := &algorithms.Hearts{
		Guideline:        *gd,
		GuidelineContent: *gdc,
		Goal:             *gl,
		GoalContent:      *glc,
	}

	hs, hg, hr, hd, hrs, err := hearts.Process(*o, colorChartPath, countriesPath)
	if err != nil {
		return ErrorResponse(c, err, 500)
	}

	output := models.NewOutput(*o.Config.Algorithm)
	output.Assessments = hs
	output.Goals = hg
	output.Referrals = hr
	output.Errors = make([]string, 0)
	output.Errors = append(output.Errors, hrs...)
	output.Meta.Debug = false
	output.Meta.CarePlan = false
	output.Meta.RiskModelVersion = *o.Config.RiskModelVersion
	output.Meta.LabBased = *o.Config.LabBased

	if o.Config.Debug != nil && *o.Config.Debug {
		output.Debug = make(map[string]interface{})
		for k, v := range hd {
			output.Debug[k] = v
		}
		output.Meta.Debug = true
	}

	if o.Config.CarePlan != nil && *o.Config.CarePlan {
		rCodes := make([]string, 0)
		if hs != nil {
			if hs.BloodPressure.Code != nil {
				rCodes = append(rCodes, *hs.BloodPressure.Code)
			}
			if hs.BodyComposition != nil && hs.BodyComposition.Components.BMI != nil && hs.BodyComposition.Components.BMI.Code != nil {
				rCodes = append(rCodes, *hs.BodyComposition.Components.BMI.Code)
			}
			if hs.BodyComposition != nil && hs.BodyComposition.Components.BodyFat != nil && hs.BodyComposition.Components.BodyFat.Code != nil {
				rCodes = append(rCodes, *hs.BodyComposition.Components.BodyFat.Code)
			}
			if hs.BodyComposition != nil && hs.BodyComposition.Components.WaistCirc != nil && hs.BodyComposition.Components.WaistCirc.Code != nil {
				rCodes = append(rCodes, *hs.BodyComposition.Components.WaistCirc.Code)
			}
			if hs.BodyComposition != nil && hs.BodyComposition.Components.WHR != nil && hs.BodyComposition.Components.WHR.Code != nil {
				rCodes = append(rCodes, *hs.BodyComposition.Components.WHR.Code)
			}
			if hs.Cholesterol != nil && hs.Cholesterol.Components.TChol != nil && hs.Cholesterol.Components.TChol.Code != nil {
				rCodes = append(rCodes, *hs.Cholesterol.Components.TChol.Code)
			}
			if hs.CVD.Code != nil {
				rCodes = append(rCodes, *hs.CVD.Code)
			}
			if hs.Diabetes.Code != nil {
				rCodes = append(rCodes, *hs.Diabetes.Code)
			}
			if hs.Lifestyle != nil && hs.Lifestyle.Components.Alcohol != nil && hs.Lifestyle.Components.Alcohol.Code != nil {
				rCodes = append(rCodes, *hs.Lifestyle.Components.Alcohol.Code)
			}
			if hs.Lifestyle != nil && hs.Lifestyle.Components.PhysicalActivity != nil && hs.Lifestyle.Components.PhysicalActivity.Code != nil {
				rCodes = append(rCodes, *hs.Lifestyle.Components.PhysicalActivity.Code)
			}
			if hs.Lifestyle != nil && hs.Lifestyle.Components.Smoking != nil && hs.Lifestyle.Components.Smoking.Code != nil {
				rCodes = append(rCodes, *hs.Lifestyle.Components.Smoking.Code)
			}
			if hs.Lifestyle != nil && hs.Lifestyle.Components.Diet != nil && hs.Lifestyle.Components.Diet.Components.Fruit != nil && hs.Lifestyle.Components.Diet.Components.Fruit.Code != nil {
				rCodes = append(rCodes, *hs.Lifestyle.Components.Diet.Components.Fruit.Code)
			}
			if hs.Lifestyle != nil && hs.Lifestyle.Components.Diet != nil && hs.Lifestyle.Components.Diet.Components.Vegetable != nil && hs.Lifestyle.Components.Diet.Components.Vegetable.Code != nil {
				rCodes = append(rCodes, *hs.Lifestyle.Components.Diet.Components.Vegetable.Code)
			}
			if hs.Lifestyle != nil && hs.Lifestyle.Components.Diet != nil && hs.Lifestyle.Components.Diet.Components.FruitVegetable != nil && hs.Lifestyle.Components.Diet.Components.FruitVegetable.Code != nil {
				rCodes = append(rCodes, *hs.Lifestyle.Components.Diet.Components.FruitVegetable.Code)
			}
		}

		goals := make([]models.CarePlanGoal, 0)
		activities := make([]string, 0)

		for _, m := range cp.CarePlanConditionMapping {
			for _, c := range m.CarePlanConditions {
				matchAND := true
				for _, s := range c {
					matchOR := false
					for _, k := range s {
						_, check := tools.SliceContainsString(rCodes, k)
						matchOR = matchOR || check
						if matchOR {
							break
						}
					}
					matchAND = matchAND && matchOR
				}
				if matchAND {
					goals = append(goals, m.CarePlanGoals...)
					activities = append(activities, getActivities(m.CarePlanActivities, o.Params.Components.Medications)...)
					break
				}
			}
		}

		outGoals := getGoalsOutput(cpc.CarePlanContentGoals, goals)
		outActivities := getActivitiesOutput(cpc.CarePlanContentActivity, activities)

		carePlan := models.CarePlanOutput{
			CarePlanOutputGoals:      outGoals,
			CarePlanOutputActivities: outActivities,
		}
		output.CarePlan = &carePlan
		output.Meta.CarePlan = true
	}

	tx, err := sqlite.DB.Begin()
	if err != nil {
		log.Println(err)
	}
	stmt, err := tx.Prepare("insert into logs(request, response) values(?, ?)")
	if err != nil {
		log.Println(err)
	}
	defer stmt.Close()

	requestObj, _ := json.Marshal(o)
	responseObj, _ := json.Marshal(output)
	_, err = stmt.Exec(string(requestObj), string(responseObj))
	if err != nil {
		tx.Rollback()
		log.Println(err)
	}
	tx.Commit()

	return c.JSON(http.StatusOK, output)
}

// processGuides maps the guide files to specific types
func processGuides(m map[string][]byte) (*a.Guidelines, *a.GuideContents, *a.GoalGuidelines, *a.GoalGuideContents, *models.CarePlanConditionsMapping, *models.CarePlanContentMapping, error) {
	engineGuide := a.Guidelines{}
	if err := json.Unmarshal(m["guide"], &engineGuide); err != nil {
		return nil, nil, nil, nil, nil, nil, err
	}

	engineGuideContent := a.GuideContents{}
	if err := json.Unmarshal(m["guideContent"], &engineGuideContent); err != nil {
		return nil, nil, nil, nil, nil, nil, err
	}

	engineGoal := a.GoalGuidelines{}
	if err := json.Unmarshal(m["goal"], &engineGoal); err != nil {
		return nil, nil, nil, nil, nil, nil, err
	}

	engineGoalContent := a.GoalGuideContents{}
	if err := json.Unmarshal(m["goalContent"], &engineGoalContent); err != nil {
		return nil, nil, nil, nil, nil, nil, err
	}

	engineCarePlan := models.CarePlanConditionsMapping{}
	if err := json.Unmarshal(m["careplan"], &engineCarePlan); err != nil {
		return nil, nil, nil, nil, nil, nil, err
	}

	engineCarePlanContent := models.CarePlanContentMapping{}
	if err := json.Unmarshal(m["careplanContent"], &engineCarePlanContent); err != nil {
		return nil, nil, nil, nil, nil, nil, err
	}

	return &engineGuide, &engineGuideContent, &engineGoal, &engineGoalContent, &engineCarePlan, &engineCarePlanContent, nil
}

// getActivities returns the appropriate activity from a list of activities according to additional rules
func getActivities(activities []models.CarePlanActivity, medications []models.ORMedication) []string {
	a := make([]string, 0)

	mCats := make(map[string]string)
	for _, m := range medications {
		mCats[*m.Class] = *m.Category
	}

	for _, act := range activities {
		rules := act["rules"]
		if len(rules) > 1 {
			activity := ""
			for _, r := range rules {
				done := false
				for k, v := range mCats {
					if r[k] == v {
						activity = strings.ToUpper(r["activity"])
						done = true
						break
					}
				}
				if done {
					break
				}
			}
			if len(activity) == 0 {
				activity = strings.ToUpper(rules[len(rules)-1]["activity"])
			}
			a = append(a, activity)
		} else {
			activity := strings.ToUpper(rules[0]["activity"])
			a = append(a, activity)
		}
	}

	return a
}

// getGoalsOutput returns the Careplan goal from a list of goals
func getGoalsOutput(g models.CarePlanContentGoals, gs []models.CarePlanGoal) []models.CarePlanContentGoal {
	goals := make([]models.CarePlanContentGoal, 0)

	for _, goal := range gs {
		s := string(goal)
		if value, ok := g[s]; ok {
			goals = append(goals, value)
		}
	}

	return goals
}

// getActivitiesOutput returns the Careplan activity from a list of activities
func getActivitiesOutput(a models.CarePlanContentActivities, as []string) []models.CarePlanContentActivity {
	activities := make([]models.CarePlanContentActivity, 0)

	for _, s := range as {
		if value, ok := a[s]; ok {
			activities = append(activities, value)
		}
	}

	return activities
}

// checkAPIToken returns the API tokens or error
func checkAPIToken(token, host, port, dbname, user, password string) (string, error) {
	psqlInfo := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		user,
		password,
		host,
		port,
		dbname)
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		return "", err
	}
	defer db.Close()

	sqlStatement := `SELECT projects.project_id AS projectname FROM integrations
LEFT JOIN projects ON (integrations.project_id = projects.id) WHERE integrations.api_key = $1
AND integrations.deleted_at IS null`
	projectName := ""
	err = db.QueryRow(sqlStatement, token).Scan(&projectName)
	if err != nil {
		return "", errors.New("no project found")
	}

	if len(projectName) > 0 {
		return projectName, nil
	}

	return "", errors.New("no project found")
}
