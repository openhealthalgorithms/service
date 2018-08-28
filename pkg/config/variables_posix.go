// +build linux darwin

package config

import (
	"net/http"
)

var (
	operationsList = []Operation{
		Operation{
			Name: OPReport,
			Endpoint: EPReport,
			SendMethod: http.MethodPost,
			Params: OpParams{
				Type: OpParamsEmpty,
				Required: opReportParamsRequired,
			},
			OutType: OpJSON,
			Timeout: opReportTimeout,
			SaveToFile: opReportSaveToFile,
			TaskMap: opReportTaskMap,
		},
		Operation{
			Name: OPETCDir,
			Endpoint: EPFile,
			SendMethod: http.MethodPut,
			Params: OpParams{
				Type: OpParamsEmpty,
				Required: opETCDirParamsRequired,
			},
			OutType: OpFileList,
			Timeout: opCollectFilesTimeout,
			SaveToFile: opETCDirSaveToFile,
			TaskMap: opETCDirTaskMap,
		},
		Operation{
			Name: OPFileSearch,
			Endpoint: EPFileSearch,
			SendMethod: http.MethodPost,
			Params: OpParams{
				Type: OpParamsAPI,
				Required: opFileSearchParamsRequired,
				Key: EPCompanyConfiguration,
			},
			OutType: OpJSON,
			Timeout: opFileSearchTimeout,
			SaveToFile: opFileSearchSaveToFile,
			TaskMap: opFileSearchTaskMap,
		},
		Operation{
			Name: OPFileSystem,
			Endpoint: EPFileSystem,
			SendMethod: http.MethodPost,
			Params: OpParams{
				Type: OpParamsEmpty,
				Required: opFileSystemParamsRequired,
			},
			OutType: OpRawBytes,
			Timeout: opFileSearchTimeout,
			SaveToFile: opFileSystemSaveToFile,
			TaskMap: opFileSystemTaskMap,
		},
		Operation{
			Name: OPConnMap,
			Endpoint: EPConnections,
			SendMethod: http.MethodPost,
			Params: OpParams{
				Type: OpParamsConfig,
				Required: opConnMapParamsRequired,
				Key: OPConnMap,
				Template: opConnMapTemplate,
			},
			OutType: OpJSON,
			Timeout: opConnMapTimeout,
			SaveToFile: opConnMapSaveToFile,
			TaskMap: opConnMapTaskMap,
		},
	}
)
