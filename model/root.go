package model

type JsonCfg struct {
	Clusters []*ClusterCfg `json:"clusters"`
}

type ClusterCfg struct {
	ClusterURI  string        `json:"clusterURI"`
	DumpCfgs    []*DumpCfg    `json:"dumpCfgs"`
	RestoreCfgs []*RestoreCfg `json:"restoreCfgs"`
}

type DumpCfg struct {
	DBName  string `json:"dbName"`
	DownDir string `json:"downDir"`
}

type RestoreCfg struct {
	DBName       string `json:"dbName"`
	UploadDBName string `json:"uploaddbName"`
	UpDir        string `json:"upDir"`
}
