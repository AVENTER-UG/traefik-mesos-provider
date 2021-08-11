package traefik_mesos_provider

type MesosTasks struct {
	Tasks []MesosTask
}

type MesosTask struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	FrameworkID string `json:"framework_id"`
	ExecutorID  string `json:"executor_id"`
	SlaveID     string `json:"slave_id"`
	State       string `json:"state"`
	Resources   struct {
		Disk  int     `json:"disk"`
		Mem   int     `json:"mem"`
		Gpus  int     `json:"gpus"`
		Cpus  float64 `json:"cpus"`
		Ports string  `json:"ports"`
	} `json:"resources"`
	Role     string `json:"role"`
	Statuses []struct {
		State           string  `json:"state"`
		Timestamp       float64 `json:"timestamp"`
		ContainerStatus struct {
			ContainerID struct {
				Value string `json:"value"`
			} `json:"container_id"`
			NetworkInfos []struct {
				IPAddresses []struct {
					Protocol  string `json:"protocol"`
					IPAddress string `json:"ip_address"`
				} `json:"ip_addresses"`
			} `json:"network_infos"`
		} `json:"container_status"`
		Healthy bool `json:"healthy,omitempty"`
	} `json:"statuses"`
	Labels    []MesosLabels `json:"labels"`
	Discovery struct {
		Visibility string `json:"visibility"`
		Name       string `json:"name"`
		Ports      struct {
			Ports []MesosPorts `json:"ports"`
		} `json:"ports"`
	} `json:"discovery"`
	Container struct {
		Type   string `json:"type"`
		Docker struct {
			Image        string `json:"image"`
			Network      string `json:"network"`
			PortMappings []struct {
				HostPort      int    `json:"host_port"`
				ContainerPort int    `json:"container_port"`
				Protocol      string `json:"protocol"`
			} `json:"port_mappings"`
			Privileged bool `json:"privileged"`
			Parameters []struct {
				Key   string `json:"key"`
				Value string `json:"value"`
			} `json:"parameters"`
			ForcePullImage bool `json:"force_pull_image"`
		} `json:"docker"`
	} `json:"container"`
	HealthCheck struct {
		DelaySeconds        int    `json:"delay_seconds"`
		IntervalSeconds     int    `json:"interval_seconds"`
		TimeoutSeconds      int    `json:"timeout_seconds"`
		ConsecutiveFailures int    `json:"consecutive_failures"`
		GracePeriodSeconds  int    `json:"grace_period_seconds"`
		Type                string `json:"type"`
		HTTP                struct {
			Protocol string `json:"protocol"`
			Scheme   string `json:"scheme"`
			Port     int    `json:"port"`
			Path     string `json:"path"`
		} `json:"http"`
	} `json:"health_check"`
}

type MesosLabels struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type MesosPorts struct {
	Number   int    `json:"number"`
	Name     string `json:"name"`
	Protocol string `json:"protocol"`
	Labels   struct {
		Labels []struct {
			Key   string `json:"key"`
			Value string `json:"value"`
		} `json:"labels"`
	} `json:"labels"`
}
