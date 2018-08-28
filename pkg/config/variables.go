package config

import (
	"time"
)

const (
	EPCompanyConfiguration = "companyConfiguration"
	EPConnections = "connections"
	EPFile = "file"
	EPFileSearch = "fileSearch"
	EPFileSystem = "fileSystem"
	EPLog = "log"
	EPRegister = "register"
	EPReport = "report"

	OPReport = "Report"
	OPETCDir = "ETCDir"
	OPIIS = "IIS"
	OPRegistry = "Registry"
	OPFileSearch = "FileSearch"
	OPFileSystem = "FileSystem"
	OPConnMap = "ConnMap"

	OPConnMapRunOnce = false
)

const (
	opReportTimeout = 5 * time.Minute
	opFileSearchTimeout = 10 * time.Hour
	opFileSystemTimeout = 10 * time.Hour
	opCollectFilesTimeout = 10 * time.Minute
	opConnMapTimeout = 12 * time.Hour

	opReportParamsRequired = false
	opETCDirParamsRequired = false
	opIISParamsRequired = false
	opRegistryParamsRequired = false
	opFileSearchParamsRequired = true
	opFileSystemParamsRequired = false
	opConnMapParamsRequired = true

	opReportSaveToFile = true
	opETCDirSaveToFile = true
	opIISSaveToFile = true
	opRegistrySaveToFile = true
	opFileSearchSaveToFile = true
	opFileSystemSaveToFile = false
	opConnMapSaveToFile = true

	opConnMapTemplate = `{"duration":%d,"runOnce":%t}`
)	

var (
	opReportTaskMap = map[string]string{
		"AWS":               "AWS",
		"CronJobs":          "CronJobs",
		"Diskspaces":        "Diskspaces",
		"DNS":               "DNS",
		"Docker":            "Docker",
		"ENV":               "Environment",
		"Firewall":          "Firewall",
		"Fstab":             "FSTab",
		"Hostname":          "Hostname",
		"Hosts":             "Hosts",
		"LinuxRepositories": "LinRepos",
		"ListenPorts":       "ListenConn",
		"Mounts":            "Mounts",
		"ActiveConnections": "NetConn",
		"Network":           "Network",
		"Os":                "OS",
		"Packages":          "Packages",
		"Platform":          "Platform",
		"Processes":         "Process",
		"Routes":            "Routes",
		"Services":          "Services",
		"SysInfo":           "SysInfo",
		"TrafRules":         "TrafRules",
		"UserInfo":          "UserInfo",
	}

	opETCDirTaskMap = map[string]string{
		"ETCDir": "ETCDir",
	}

	opIISTaskMap = map[string]string{
		"IIS": "IIS",
	}

	opRegistryTaskMap = map[string]string{
		"Registry": "Registry",
	}

	opFileSearchTaskMap = map[string]string{
		"FileSearch": "FileSearch",
	}

	opFileSystemTaskMap = map[string]string{
		"FileSystem": "FileSystem",
	}

	opConnMapTaskMap = map[string]string{
		"ConnMap": "ConnMap",
	}
)
