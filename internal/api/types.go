package api

type Pod struct {
	ID            string   `json:"id"`
	Name          string   `json:"name"`
	DesiredStatus string   `json:"desiredStatus"`
	CostPerHr     float64  `json:"costPerHr"`
	GpuCount      int      `json:"gpuCount"`
	MemoryInGb    float64  `json:"memoryInGb"`
	VcpuCount     int      `json:"vcpuCount"`
	UptimeSeconds int      `json:"uptimeSeconds"`
	Machine       Machine  `json:"machine"`
	Runtime       *Runtime `json:"runtime"`
}

type Machine struct {
	GpuDisplayName string `json:"gpuDisplayName"`
	Location       string `json:"location"`
}

type Runtime struct {
	UptimeInSeconds int       `json:"uptimeInSeconds"`
	Gpus            []GPU     `json:"gpus"`
	Container       Container `json:"container"`
	Ports           []Port    `json:"ports"`
}

type GPU struct {
	ID                string  `json:"id"`
	GpuUtilPercent    float64 `json:"gpuUtilPercent"`
	MemoryUtilPercent float64 `json:"memoryUtilPercent"`
}

type Container struct {
	CpuPercent    float64 `json:"cpuPercent"`
	MemoryPercent float64 `json:"memoryPercent"`
}

type Port struct {
	IP          string `json:"ip"`
	IsIpPublic  bool   `json:"isIpPublic"`
	PrivatePort int    `json:"privatePort"`
	PublicPort  int    `json:"publicPort"`
	Type        string `json:"type"`
}

type graphQLResponse struct {
	Data struct {
		Myself struct {
			Pods []Pod `json:"pods"`
		} `json:"myself"`
	} `json:"data"`
	Errors []graphQLError `json:"errors"`
}

type graphQLError struct {
	Message string `json:"message"`
}
