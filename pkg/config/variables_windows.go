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
			Name: OPIIS,
			Endpoint: EPFile,
			SendMethod: http.MethodPut,
			Params: OpParams{
				Type: OpParamsEmpty,
				Required: opIISParamsRequired,
			},
			OutType: OpFileList,
			Timeout: opCollectFilesTimeout,
			SaveToFile: opIISSaveToFile,
			TaskMap: opIISTaskMap,
		},
		Operation{
			Name: OPRegistry,
			Endpoint: EPFile,
			SendMethod: http.MethodPut,
			Params: OpParams{
				Type: OpParamsEmpty,
				Required: opRegistryParamsRequired,
			},
			OutType: OpFileList,
			Timeout: opCollectFilesTimeout,
			SaveToFile: opRegistrySaveToFile,
			TaskMap: opRegistryTaskMap,
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
