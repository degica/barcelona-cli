package api

import (
	"fmt"
	"strings"

	"github.com/degica/barcelona-cli/config"
)

type DistrictResponse struct {
	District    *District   `json:"district,omitempty"`
	Districts   []*District `json:"districts,omitempty"`
	Certificate string      `json:"certificate,omitempty"`
}

type DistrictRequest struct {
	Name                string `json:"name,omitempty"`
	Region              string `json:"region,omitempty"`
	ClusterSize         *int   `json:"cluster_size,omitempty"`
	ClusterInstanceType string `json:"cluster_instance_type,omitempty"`
	NatType             string `json:"nat_type,omitempty"`
	ClusterBackend      string `json:"cluster_backend,omitempty"`
	AwsAccessKeyId      string `json:"aws_access_key_id,omitempty"`
	AwsSecretAccessKey  string `json:"aws_secret_access_key,omitempty"`
}

type District struct {
	Name                string               `json:"name"`
	Region              string               `json:"region"`
	BastionIP           string               `json:"bastion_ip"`
	ClusterSize         int                  `json:"cluster_size"`
	ClusterInstanceType string               `json:"cluster_instance_type"`
	S3BucketName        string               `json:"s3_bucket_name"`
	StackStatus         string               `json:"stack_status"`
	NatType             string               `json:"nat_type"`
	ClusterBackend      string               `json:"cluster_backend"`
	CidrBlock           string               `json:"cidr_block"`
	StackName           string               `json:"stack_name"`
	AwsAccessKeyId      string               `json:"aws_access_key_id"`
	AwsRole             string               `json:"aws_role"`
	ContainerInstances  []*ContainerInstance `json:"container_instances"`
	Plugins             []*Plugin            `json:"plugins"`
	Heritages           []*Heritage          `json:"heritages"`
	Notifications       []*Notification      `json:"notifications"`
}

type ContainerInstance struct {
	ContainerInstanceArn string `json:"container_instance_arn"`
	EC2InstanceID        string `json:"ec2_instance_id"`
	PendingTasksCount    int    `json:"pending_tasks_count"`
	RunningTasksCount    int    `json:"running_tasks_count"`
	Status               string `json:"status"`
	PrivateIPAddress     string `json:"private_ip_address"`
	RemainingResources   []struct {
		Name         string `json:"name"`
		IntegerValue int    `json:"integer_value"`
	} `json:"remaining_resources"`
}

type Plugin struct {
	Name       string            `json:"name,omitempty"`
	Attributes map[string]string `json:"attributes,omitempty"`
}

type PluginResponse struct {
	Plugin  *Plugin   `json:"plugin"`
	Plugins []*Plugin `json:"plugins"`
}

type HeritageResponse struct {
	Heritage  *Heritage   `json:"heritage,omitempty"`
	Heritages []*Heritage `json:"heritages,omitempty"`
}

type EnvironmentPair struct {
	Name      string  `yaml:"name" json:"name"`
	Value     *string `yaml:"value,omitempty" json:"value,omitempty"`
	ValueFrom *string `yaml:"value_from,omitempty" json:"value_from,omitempty"`
	SsmPath   *string `yaml:"ssm_path,omitempty" json:"ssm_path,omitempty"`
}

type ReviewGroupRequest struct {
	Name         string `json:"name"`
	BaseDomain   string `json:"base_domain"`
	EndpointName string `json:"endpoint"`
}

type ReviewGroup struct {
	Name       string       `json:"name"`
	BaseDomain string       `json:"base_domain"`
	Endpoint   *Endpoint    `json:"endpoint"`
	Token      *string      `json:"token,omitempty"`
	ReviewApps []*ReviewApp `json:"review_apps,omitempty"`
}

type ReviewGroupResponse struct {
	ReviewGroup  *ReviewGroup   `json:"review_group,omitempty"`
	ReviewGroups []*ReviewGroup `json:"review_groups,omitempty"`
}

type ReviewAppService struct {
	Name        string      `yaml:"name" json:"name"`
	ServiceType string      `yaml:"service_type" json:"service_type"`
	Cpu         int         `yaml:"cpu" json:"cpu,omitempty"`
	Memory      int         `yaml:"memory" json:"memory"`
	Command     string      `yaml:"command" json:"command"`
	ForceSsl    bool        `yaml:"force_ssl" json:"force_ssl"`
	Listeners   []*Listener `yaml:"listeners" json:"listeners"`
}

type ReviewAppDefinition struct {
	GroupName   string              `yaml:"group" json:"group_name"`
	ImageName   string              `yaml:"image_name" json:"image_name"`
	Environment []*EnvironmentPair  `yaml:"environment" json:"environment"`
	Services    []*ReviewAppService `yaml:"services" json:"services"`
}

type ReviewAppRequest struct {
	*ReviewAppDefinition
	Subject   string `json:"subject"`
	Retention int    `json:"retention"`
	ImageTag  string `json:"image_tag"`
}

type ReviewApp struct {
	Template    Heritage     `yaml:"template" json:"template,omitempty"`
	Heritage    Heritage     `json:"heritage"`
	Subject     string       `json:"subject"`
	Tag         string       `json:"tag"`
	Domain      string       `json:"domain"`
	ReviewGroup *ReviewGroup `json:"review_group"`
}

type ReviewAppResponse struct {
	ReviewApp  *ReviewApp   `json:"review_app,omitempty"`
	ReviewApps []*ReviewApp `json:"review_apps,omitempty"`
}

type RunEnv struct {
	Vars map[string]string `yaml:"vars" json:"vars"`
}

type Heritage struct {
	Name           string             `yaml:"name" json:"name"`
	ImageName      string             `yaml:"image_name" json:"image_name"`
	ImageTag       string             `yaml:"image_tag,omitempty" json:"image_tag,omitempty"`
	BeforeDeploy   *string            `yaml:"before_deploy" json:"before_deploy"`
	Version        int                `yaml:"version,omitempty" json:"version,omitempty"`
	ScheduledTasks []*ScheduledTask   `yaml:"scheduled_tasks" json:"scheduled_tasks"`
	Services       []*Service         `yaml:"services" json:"services"`
	EnvVars        map[string]string  `json:"env_vars,omitempty"`
	Environment    []*EnvironmentPair `yaml:"environment" json:"environment"`
	Token          string             `json:"token,omitempty"`
	RunEnv         *RunEnv            `yaml:"run_env,omitempty" json:"run_env,omitempty"`
}

func (h *Heritage) FillinDefaults() {
	if h.ScheduledTasks == nil {
		h.ScheduledTasks = []*ScheduledTask{}
	}
	if h.Services == nil {
		h.Services = []*Service{}
	}
	for _, service := range h.Services {
		service.FillinDefaults()
	}
}

type Service struct {
	Public           *bool          `yaml:"public,omitempty" json:"public,omitempty"`
	Name             string         `yaml:"name" json:"name"`
	Cpu              int            `yaml:"cpu" json:"cpu,omitempty"`
	Memory           int            `yaml:"memory" json:"memory"`
	Command          string         `yaml:"command" json:"command"`
	ServiceType      string         `yaml:"service_type" json:"service_type"`
	WebContainerPort int            `yaml:"web_container_port,omitempty" json:"web_container_port,omitempty"`
	ForceSsl         bool           `yaml:"force_ssl" json:"force_ssl"`
	PortMappings     []*PortMapping `yaml:"port_mappings,omitempty" json:"port_mappings,omitempty"`
	Hosts            []*Host        `yaml:"hosts" json:"hosts"`
	Listeners        []*Listener    `yaml:"listeners,omitempty" json:"listeners,omitempty"`
	// Response only parameters
	Status       string `json:"string,omitempty"`
	RunningCount int    `json:"running_count,omitempty"`
	PendingCount int    `json:"pending_count,omitempty"`
	DesiredCount int    `json:"desired_count,omitempty"`
}

func (s *Service) FillinDefaults() {
	if s.Memory == 0 {
		s.Memory = 512
	}

	if s.ServiceType == "" {
		s.ServiceType = "default"
	}

	if s.Hosts == nil {
		s.Hosts = []*Host{}
	}

	if s.Listeners == nil {
		s.Listeners = []*Listener{}
	}
}

type Listener struct {
	Endpoint                string          `yaml:"endpoint" json:"endpoint"`
	HealthCheckInterval     int             `yaml:"health_check_interval,omitempty" json:"health_check_interval,omitempty"`
	HealthCheckPath         string          `yaml:"health_check_path,omitempty" json:"health_check_path,omitempty"`
	HealthCheckTimeout      int             `yaml:"health_check_timeout,omitempty" json:"health_check_timeout,omitempty"`
	HealthyThresholdCount   int             `yaml:"healthy_threshold_count,omitempty" json:"healthy_threshold_count,omitempty"`
	UnhealthyThresholdCount int             `yaml:"unhealthy_threshold_count,omitempty" json:"unhealthy_threshold_count,omitempty"`
	RulePriority            int             `yaml:"rule_priority,omitempty" json:"rule_priority,omitempty"`
	RuleConditions          []RuleCondition `yaml:"rule_conditions,omitempty" json:"rule_conditions,omitempty"`
}

type RuleCondition struct {
	Type  string `yaml:"type" json:"type"`
	Value string `yaml:"value" json:"value"`
}

type ScheduledTask struct {
	Schedule string `json:"schedule"`
	Command  string `json:"command"`
}

type PortMapping struct {
	LbPort              int    `yaml:"lb_port" json:"lb_port"`
	HostPort            int    `yaml:"host_port,omitempty" json:"host_port,omitempty"`
	ContainerPort       int    `yaml:"container_port" json:"container_port"`
	Protocol            string `yaml:"protocol,omitempty" json:"protocol,omitempty"`
	EnableProxyProtocol bool   `yaml:"enable_proxy_protocol,omitempty" json:"enable_proxy_protocol,omitempty"`
}

type Host struct {
	Hostname    string `yaml:"hostname" json:"hostname"`
	SslCertPath string `yaml:"ssl_cert_path" json:"ssl_cert_path"`
	SslKeyPath  string `yaml:"ssl_key_path" json:"ssl_key_path"`
}

type OneoffResponse struct {
	Oneoff      *Oneoff `json:"oneoff"`
	Oneoffs     *Oneoff `json:"oneoffs"`
	Certificate string  `json:"certificate"`
}

type Oneoff struct {
	ID                    int       `json:"id"`
	TaskARN               string    `json:"task_arn"`
	Command               string    `json:"command"`
	Status                string    `json:"status"`
	ExitCode              string    `json:"exit_code"`
	Reason                string    `json:"reason"`
	InteractiveRunCommand string    `json:"interactive_run_command"`
	ContainerInstanceARN  string    `json:"container_instance_arn"`
	ContainerName         string    `json:"container_name"`
	Memory                int       `json:"memory"`
	District              *District `json:"district"`
}

type EndpointResponse struct {
	Endpoint  *Endpoint   `json:"endpoint"`
	Endpoints []*Endpoint `json:"endpoints"`
}

type Endpoint struct {
	Name          string `json:"name,omitempty"`
	Public        *bool  `json:"public,omitempty"`
	CertificateID string `json:"certificate_id,omitempty"`
	SslPolicy     string `json:"ssl_policy,omitempty"`
	// Response only
	DNSName  string    `json:"dns_name,omitempty"`
	District *District `json:"district,omitempty"`
}

type User struct {
	Token     string `json:"token,omitempty"`
	Name      string `json:"name"`
	PublicKey string `json:"public_key"`
}

type UserResponse struct {
	User  *User   `json:"user"`
	Users []*User `json:"users"`
}

type Notification struct {
	ID       int    `json:"id,omitempty"`
	Target   string `json:"target,omitempty"`
	Endpoint string `json:"endpoint,omitempty"`
}

type NotificationResponse struct {
	Notification  *Notification   `json:"notification,omitempty"`
	Notifications []*Notification `json:"notifications,omitempty"`
}

type APIError struct {
	Message      string   `json:"error"`
	DebugMessage string   `json:"debug_message"`
	Backtrace    []string `json:"backtrace"`
}

type VaultAuthResponseAuth struct {
	ClientToken string `json:"client_token"`
}

type VaultAuthResponse struct {
	Auth VaultAuthResponseAuth `json:"auth"`
}

func (err *APIError) Error() string {
	if config.Get().IsDebug() {
		return fmt.Sprintf("%s\n\n%s\n%s\n", err.Message, err.DebugMessage, strings.Join(err.Backtrace, "\n"))
	} else {
		return err.Message
	}
}
