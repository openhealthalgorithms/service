package tools

import (
	"context"
	"io/ioutil"
	"os"
	"strings"

	"cloud.google.com/go/storage"
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
	"google.golang.org/api/iterator"

	"github.com/openhealthalgorithms/service/config"
)

// ParseGuidesFiles function
func ParseGuidesFiles(c echo.Context) (map[string][]byte, error) {
	var guide, guideContent, goal, goalContent, careplan, careplanContent []byte
	var err error

	currentSettings := c.Get("current_config").(config.Settings)

	if currentSettings.CloudEnable {
		projectName := strings.Replace(c.Request().URL.Path, "/api/algorithm", "", 1)
		if len(projectName) > 1 {
			projectName = projectName[1:]
		}
		if len(projectName) == 0 {
			projectName = "default-json"
		}

		err := os.Setenv("GOOGLE_APPLICATION_CREDENTIALS", currentSettings.CloudConfigFile)
		if err != nil {
			return nil, err
		}

		ctxBack := context.Background()

		// Creates a client.
		client, err := storage.NewClient(ctxBack)
		if err != nil {
			return nil, err
		}
		bucket := client.Bucket(currentSettings.CloudBucket)
		objs := bucket.Objects(ctxBack, &storage.Query{
			Prefix:    projectName + "/",
			Delimiter: "/",
		})
		i := 0
		for {
			attrs, err := objs.Next()
			if err == iterator.Done {
				break
			}
			if err != nil {
				return nil, err
			}
			name := strings.ToLower(attrs.Name)
			if strings.Contains(name, "guideline_hearts.json") {
				i++
				guide, err = readStorageContent(ctxBack, bucket, attrs.Name)
				if err != nil {
					return nil, err
				}
			} else if strings.Contains(name, "guideline_hearts_content.json") {
				i++
				guideContent, err = readStorageContent(ctxBack, bucket, attrs.Name)
				if err != nil {
					return nil, err
				}
			} else if strings.Contains(name, "goals_hearts.json") {
				i++
				goal, err = readStorageContent(ctxBack, bucket, attrs.Name)
				if err != nil {
					return nil, err
				}
			} else if strings.Contains(name, "goals_hearts_content.json") {
				i++
				goalContent, err = readStorageContent(ctxBack, bucket, attrs.Name)
				if err != nil {
					return nil, err
				}
			} else if strings.Contains(name, "careplan_conditions.json") {
				i++
				careplan, err = readStorageContent(ctxBack, bucket, attrs.Name)
				if err != nil {
					return nil, err
				}
			} else if strings.Contains(name, "careplan_content.json") {
				i++
				careplanContent, err = readStorageContent(ctxBack, bucket, attrs.Name)
				if err != nil {
					return nil, err
				}
			}
		}

		if i != 6 {
			return nil, errors.New("guideline files for the project are missing")
		}
	} else {
		guide, err = ioutil.ReadFile(currentSettings.GuidelineFile)
		if err != nil {
			return nil, errors.New("Read file error: " + err.Error())
		}

		guideContent, err = ioutil.ReadFile(currentSettings.GuidelineContentFile)
		if err != nil {
			return nil, errors.New("Read file error: " + err.Error())
		}

		goal, err = ioutil.ReadFile(currentSettings.GoalFile)
		if err != nil {
			return nil, errors.New("Read file error: " + err.Error())
		}

		goalContent, err = ioutil.ReadFile(currentSettings.GoalContentFile)
		if err != nil {
			return nil, errors.New("Read file error: " + err.Error())
		}

		careplan, err = ioutil.ReadFile(currentSettings.CareplanConditionsFile)
		if err != nil {
			return nil, errors.New("Read file error: " + err.Error())
		}

		careplanContent, err = ioutil.ReadFile(currentSettings.CareplanContentFile)
		if err != nil {
			return nil, errors.New("Read file error: " + err.Error())
		}
	}

	r := make(map[string][]byte)
	r["guide"] = guide
	r["guideContent"] = guideContent
	r["goal"] = goal
	r["goalContent"] = goalContent
	r["careplan"] = careplan
	r["careplanContent"] = careplanContent

	return r, nil
}

func readStorageContent(ctx context.Context, bucket *storage.BucketHandle, name string) ([]byte, error) {
	rc, err := bucket.Object(name).NewReader(ctx)
	if err != nil {
		return nil, err
	}
	defer rc.Close()

	data, err := ioutil.ReadAll(rc)
	if err != nil {
		return nil, err
	}
	return data, nil
}
