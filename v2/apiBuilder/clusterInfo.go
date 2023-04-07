package apiBuilder

type ClusterInfo struct {
	Alloc        string `json:"Alloc"`
	ClientId     string `json:"ClientId"`
	ClusterIP    string `json:"ClusterIP"`
	CpuNum       int    `json:"CpuNum"`
	FreesTimes   int    `json:"FreesTimes"`
	GCTime       string `json:"GCTime"`
	GoArch       string `json:"GoArch"`
	GoVersion    string `json:"GoVersion"`
	GoroutineNum int    `json:"GoroutineNum"`
	LastGCTime   string `json:"LastGCTime"`
	MacInfo      string `json:"MacInfo"`
	MallocsTimes int    `json:"MallocsTimes"`
	NextGC       string `json:"NextGC"`
	Platform     string `json:"Platform"`
	QQUsers      []struct {
		MoneyCount    string `json:"MoneyCount"`
		OnlieTime     string `json:"OnlieTime"`
		QQ            string `json:"QQ"`
		ReceiveCount  int    `json:"ReceiveCount"`
		SendCount     int    `json:"SendCount"`
		TotalMoney    string `json:"TotalMoney"`
		TotalRecv     string `json:"TotalRecv"`
		TotalSend     string `json:"TotalSend"`
		UserLevelInfo string `json:"UserLevelInfo"`
	} `json:"QQUsers"`
	QQUsersCounts int    `json:"QQUsersCounts"`
	ServerRuntime string `json:"ServerRuntime"`
	Sys           string `json:"Sys"`
	TotalAlloc    string `json:"TotalAlloc"`
	Version       string `json:"Version"`
}

func (b *Builder) GetClusterInfo() DoApi {
	method := "GET"
	path := "/v1/clusterinfo"
	b.method = &method
	b.path = &path
	return b
}
