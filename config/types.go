package config

// Settings represents settings for the application.
type Settings struct {
    Port                   string
    GuidelineFile          string
    GuidelineContentFile   string
    GoalFile               string
    GoalContentFile        string
    CareplanConditionsFile string
    CareplanContentFile    string
    LogFile                string
    CloudEnable            bool
    CloudBucket            string
    CloudConfigFile        string
    CloudDBHost            string
    CloudDBPort            string
    CloudDBName            string
    CloudDBUser            string
    CloudDBPassword        string
}
