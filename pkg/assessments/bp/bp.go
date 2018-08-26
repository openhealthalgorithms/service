package bp

import (
	"context"
	"github.com/fatih/structs"
	"github.com/pkg/errors"

	"github.com/openhealthalgorithms/service/pkg/tools"
)

// Data holds results of plugin.
type Data struct {
	BP `structs:"BP_ASSESSMENT"`
}

// New returns a ready to use instance of the plugin.
func New() *Data {
	return &Data{}
}

// Get fills the Data and returns error.
func (d *Data) Get(ctx context.Context) (err error) {
	defer func() {
		if r := recover(); r != nil {
			var ok bool
			err, ok = r.(error)
			if !ok {
				err = errors.Errorf("%v", r)
			}
		}
	}()

	return d.get(ctx)
}

// Output returns information gathered by the plugin.
func (d *Data) Output() (map[string]interface{}, error) {
	return structs.Map(d), nil
}

// get does all the job.
func (d *Data) get(ctx context.Context) error {
	inputs := tools.ParseParams(ctx)
	resultCode := ""

	sbp := inputs.Sbp
	dbp := inputs.Dbp

	currentBp := tools.JoinStrings(string(sbp), "/", string(dbp))
	target := currentBp

	if inputs.Diabetes {
		if sbp > 130 {
			resultCode = "BP-3B"
			target = "130/80"
		} else {
			resultCode = "BP-3A"
			target = "130/80"
		}
	} else {
		if sbp > 160 {
			resultCode = "BP-2"
			target = "140/90"
		} else if sbp > 140 {
			resultCode = "BP-1B"
			target = "140/90"
		} else if sbp <= 140 && sbp >= 120 {
			resultCode = "BP-1A"
			target = "140/90"
		} else {
			resultCode = "BP-0"
			target = "140/90"
		}
	}

	d.BP = NewBP(currentBp, sbp, dbp, resultCode, target)

	return nil
}

type BP struct {
	BP     string `structs:"bp"`
	SBP    int    `structs:"sbp"`
	DBP    int    `structs:"dbp"`
	Code   string `structs:"code"`
	Target string `structs:"target"`
}

// NewBP returns a BP object.
func NewBP(bp string, sbp, dbp int, code, target string) BP {
	return BP{
		BP:     bp,
		SBP:    sbp,
		DBP:    dbp,
		Code:   code,
		Target: target,
	}
}
