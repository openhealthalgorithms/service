package actions

import (
    "database/sql"
    "fmt"
    "net/http"

    "github.com/labstack/echo/v4"
    "github.com/pkg/errors"

    "github.com/openhealthalgorithms/service/models"

    _ "github.com/lib/pq"
)

// AlgorithmHandler function
func AlgorithmHandler(c echo.Context) error {
    var err error
    o := new(models.OHARequest)
    if err = c.Bind(o); err != nil {
        return ErrorResponse(c, err, 500)
    }
    if err = c.Validate(o); err != nil {
        return ErrorResponse(c, err, 400)
    }



    // rCodes := make([]string, 0)
    // if o.Response.Hearts != nil {
    //     a := *o.Response.Hearts.Assessments
    //     rCodes = append(rCodes, *a.BloodPressure.Code)
    //     rCodes = append(rCodes, *a.BodyComposition.Components.BMI.Code)
    //     rCodes = append(rCodes, *a.BodyComposition.Components.BodyFat.Code)
    //     rCodes = append(rCodes, *a.BodyComposition.Components.WaistCirc.Code)
    //     rCodes = append(rCodes, *a.BodyComposition.Components.WHR.Code)
    //     rCodes = append(rCodes, *a.Cholesterol.Components.TChol.Code)
    //     rCodes = append(rCodes, *a.CVD.Code)
    //     rCodes = append(rCodes, *a.Diabetes.Code)
    //     rCodes = append(rCodes, *a.Lifestyle.Components.Alcohol.Code)
    //     rCodes = append(rCodes, *a.Lifestyle.Components.PhysicalActivity.Code)
    //     rCodes = append(rCodes, *a.Lifestyle.Components.Smoking.Code)
    //     rCodes = append(rCodes, *a.Lifestyle.Components.Diet.Components.Fruit.Code)
    //     rCodes = append(rCodes, *a.Lifestyle.Components.Diet.Components.Vegetable.Code)
    //     rCodes = append(rCodes, *a.Lifestyle.Components.Diet.Components.FruitVegetable.Code)
    // }
    //
    // mappings, err := getMapping()
    // if err != nil {
    //     return ErrorResponse(c, err, 500)
    // }
    //
    // contents, err := getContent()
    // if err != nil {
    //     return ErrorResponse(c, err, 500)
    // }
    //
    // goals := make([]models.Goal, 0)
    // activities := make([]string, 0)
    //
    // for _, m := range mappings.Mapping {
    //     for _, c := range m.Conditions {
    //         matchAND := true
    //         for _, s := range c {
    //             matchOR := false
    //             for _, k := range s {
    //                 _, check := tools.SliceContainsString(rCodes, k)
    //                 matchOR = matchOR || check
    //                 if matchOR {
    //                     break
    //                 }
    //             }
    //             matchAND = matchAND && matchOR
    //         }
    //         if matchAND {
    //             break
    //         }
    //     }
    //     goals = append(goals, m.Goals...)
    //     activities = append(activities, getActivities(m.Activities, o.Request.Params.Components.Medications)...)
    // }
    //
    // outGoals := getGoalsOutput(contents.Goals, goals)
    // outActivities := getActivitiesOutput(contents.Activity, activities)
    //
    // out := models.IOutput{}
    // out.Goals = outGoals
    // out.Activities = outActivities

    return c.JSON(http.StatusOK, o)
}

// func getMapping() (*models.IMapping, error) {
//     mappings := new(models.IMapping)
//
//     mappingFile := envy.Get("MAPPING_FILE", "")
//     if _, err := os.Stat(mappingFile); err != nil {
//         return mappings, err
//     }
//
//     mappingContent, err := ioutil.ReadFile(mappingFile)
//     if err != nil {
//         return mappings, err
//     }
//
//     if err := json.Unmarshal(mappingContent, &mappings); err != nil {
//         return mappings, err
//     }
//
//     return mappings, nil
// }
//
// func getContent() (*models.IContent, error) {
//     content := new(models.IContent)
//
//     contentFile := envy.Get("CONTENT_FILE", "")
//     if _, err := os.Stat(contentFile); err != nil {
//         return content, err
//     }
//
//     contentContent, err := ioutil.ReadFile(contentFile)
//     if err != nil {
//         return content, err
//     }
//
//     if err := json.Unmarshal(contentContent, &content); err != nil {
//         return content, err
//     }
//
//     return content, nil
// }
//
// func getActivities(activities []models.Activity, medications []models.ORMedication) []string {
//     a := make([]string, 0)
//
//     mCats := make(map[string]string)
//     for _, m := range medications {
//         mCats[*m.Category] = *m.Generic
//     }
//
//     for _, act := range activities {
//         rules := act["rules"]
//         if len(rules) > 1 {
//             activity := ""
//             for _, r := range rules {
//                 for k, v := range mCats {
//                     if r[k] == v {
//                         activity = strings.ToUpper(r["activity"])
//                     }
//                 }
//             }
//             if len(activity) == 0 {
//                 activity = strings.ToUpper(rules[len(rules)-1]["activity"])
//             }
//             a = append(a, activity)
//         } else {
//             activity := strings.ToUpper(rules[0]["activity"])
//             a = append(a, activity)
//         }
//     }
//
//     return a
// }
//
// func getGoalsOutput(g models.IContentGoals, gs []models.Goal) []models.IContentGoal {
//     goals := make([]models.IContentGoal, 0)
//
//     for _, goal := range gs {
//         s := string(goal)
//         if value, ok := g[s]; ok {
//             goals = append(goals, value)
//         }
//     }
//
//     return goals
// }
//
// func getActivitiesOutput(a models.IContentActivities, as []string) []models.IContentActivity {
//     activities := make([]models.IContentActivity, 0)
//
//     for _, s := range as {
//         if value, ok := a[s]; ok {
//             activities = append(activities, value)
//         }
//     }
//
//     return activities
// }

func checkAPIToken(token, host, dbname, user, password string) (string, error) {
    psqlInfo := fmt.Sprintf("postgres://%s:%s@%s:5432/%s?sslmode=disable",
        user,
        password,
        host,
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
