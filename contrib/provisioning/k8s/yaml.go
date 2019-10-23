package k8s

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"os"
	"strconv"
)

type Deployment struct {
	ApiVersion     string `yaml:"apiVersion"`
	Kind           string
	MetaData       DeploymentMetadata `yaml:"metadata"`
	DeploymentSpec DeploymentSpec     `yaml:"spec"`
}
type HttpGet struct {
	Port int
	Path string
}
type LivenessProbe struct {
	HttpGet HttpGet `yaml:"httpGet"`
}
type ReadinessProbe struct {
	HttpGet             HttpGet `yaml:"httpGet"`
	InitialDelaySeconds int     `yaml:"initialDelaySeconds"`
	PeriodSeconds       int     `yaml:"periodSeconds"`
}
type SecurityContext struct {
	Privileged bool
}
type Container struct {
	Env                      []Env `yaml:"env"`
	Image                    string
	LivenessProbe            LivenessProbe   `yaml:"livenessProbe"`
	ReadinessProbe           ReadinessProbe  `yaml:"readinessProbe"`
	SecurityContext          SecurityContext `yaml:"securityContext"`
	ImagePullPolicy          string          `yaml:"imagePullPolicy"`
	Name                     string
	Ports                    []Port        `yaml:"ports"`
	TerminationMessagePath   string        `yaml:"terminationMessagePath"`
	TerminationMessagePolicy string        `yaml:"terminationMessagePolicy"`
	VolumeMounts             []VolumeMount `yaml:"volumeMounts"`
}
type Env struct {
	Name  string
	Value string
}
type Port struct {
	ContainerPort int    `yaml:"containerPort"`
	Protocol      string `yaml:"protocol"`
}
type VolumeMount struct {
	MountPath string `yaml:"mountPath"`
	Name      string `yaml:"name"`
}
type PodMatchLabels struct {
	ValidatorOrder string `yaml:"validator-order"`
	P2PPort        string `yaml:"p2p-port"`
	ChainId        string `yaml:"chain-id"`
}
type DeploymentSelector struct {
	MatchLabels PodMatchLabels `yaml:"matchLabels"`
}
type NodeSelector struct {
	ValidatorOrder string `yaml:"validator-order"`
}
type DeploymentSpecTemplate struct {
	Metadata PodMetadata `yaml:"metadata"`
	Spec     struct {
		Containers    []Container
		NodeSelector  NodeSelector `yaml:"nodeSelector"`
		HostNetwork   bool         `yaml:"hostNetwork"`
		DnsPolicy     string       `yaml:"dnsPolicy"`
		RestartPolicy string       `yaml:"restartPolicy"`
		Volumes       []Volume     `yaml:"volumes"`
	}
}
type DeploymentStrategy map[string]string

type DeploymentSpec struct {
	Replicas               int
	DeploymentSelector     DeploymentSelector     `yaml:"selector"`
	DeploymentStrategy     DeploymentStrategy     `yaml:"strategy"`
	DeploymentSpecTemplate DeploymentSpecTemplate `yaml:"template"`
}
type Volume struct {
	HostPath map[string]string `yaml:"hostPath"`
	Name     string
}
type PodMetadata struct {
	Labels struct {
		ValidatorOrder string `yaml:"validator-order"`
		P2PPort        string `yaml:"p2p-port"`
		ChainId        string `yaml:"chain-id"`
	}
}
type DeploymentLabels struct {
	ValidatorOrder string `yaml:"validator-order"`
	P2PPort        string `yaml:"p2p-port"`
	ABCIPort       string `yaml:"abci-port"`
	RestAPIPort    string `yaml:"restapi-port"`
	ChainId        string `yaml:"chain-id"`
}
type DeploymentMetadata struct {
	Labels DeploymentLabels `yaml:"labels"`
	Name   string
}

const dirPermission = 0700

type DeploymentTemplate struct {
	Node *Node
}

func (k *DeploymentTemplate) BinaryHome() string {
	return fmt.Sprintf("/%s/%s", k.Node.MetaData.ChainID, "binary/linkd")
}
func (k *DeploymentTemplate) LinkdHome() string {
	return fmt.Sprintf("/linkd/%s/node%d/linkd/", k.Node.MetaData.ChainID, k.Node.Idx)
}
func (k *DeploymentTemplate) DeployPath() string {
	return fmt.Sprintf("./build/%s/deployments/", k.Node.MetaData.ChainID)
}
func (k *DeploymentTemplate) DeployName() string {
	return fmt.Sprintf("validator-%s-%s-%s-%s",
		strconv.Itoa(k.Node.Idx),
		prefixForP2PPort+strconv.Itoa(k.Node.MetaData.NodeP2PPort),
		prefixPortRestAPIPort+strconv.Itoa(k.Node.MetaData.NodeRestAPIPort),
		prefixABCIPort+strconv.Itoa(k.Node.MetaData.NodeABCIPort),
	)
}

func (k *DeploymentTemplate) outputFileName() string {
	return fmt.Sprintf("validator-%d.yaml", k.Node.Idx)
}
func (k *DeploymentTemplate) Write() (*Deployment, error) {

	deploy := &Deployment{}
	bytes, err := ioutil.ReadFile(k.Node.MetaData.k8STemplateFilePath)
	if err != nil {
		panic(err)
	}
	err = yaml.Unmarshal(bytes, &deploy)
	if err != nil {
		return deploy, err
	}
	deploy.MetaData.Labels.ValidatorOrder = strconv.Itoa(k.Node.Idx)
	deploy.MetaData.Labels.P2PPort = strconv.Itoa(k.Node.MetaData.NodeP2PPort)
	deploy.MetaData.Labels.ABCIPort = strconv.Itoa(k.Node.MetaData.NodeABCIPort)
	deploy.MetaData.Labels.RestAPIPort = strconv.Itoa(k.Node.MetaData.NodeRestAPIPort)
	deploy.MetaData.Labels.ChainId = k.Node.MetaData.ChainID
	deploy.MetaData.Name = k.DeployName()

	deploy.DeploymentSpec.DeploymentSelector.MatchLabels.ValidatorOrder = strconv.Itoa(k.Node.Idx)
	deploy.DeploymentSpec.DeploymentSelector.MatchLabels.P2PPort = strconv.Itoa(k.Node.MetaData.NodeP2PPort)
	deploy.DeploymentSpec.DeploymentSelector.MatchLabels.ChainId = k.Node.MetaData.ChainID

	deploy.DeploymentSpec.DeploymentSpecTemplate.Metadata.Labels.ValidatorOrder = strconv.Itoa(k.Node.Idx)
	deploy.DeploymentSpec.DeploymentSpecTemplate.Metadata.Labels.P2PPort = strconv.Itoa(k.Node.MetaData.NodeP2PPort)
	deploy.DeploymentSpec.DeploymentSpecTemplate.Metadata.Labels.ChainId = k.Node.MetaData.ChainID

	deploy.DeploymentSpec.DeploymentSpecTemplate.Spec.NodeSelector = NodeSelector{strconv.Itoa(k.Node.Idx)}
	container := &deploy.DeploymentSpec.DeploymentSpecTemplate.Spec.Containers[0]
	container.Image = k.Node.MetaData.linkDockerImageUrl
	container.LivenessProbe.HttpGet.Port, _ = strconv.Atoi(deploy.MetaData.Labels.RestAPIPort)
	container.ReadinessProbe.HttpGet.Port, _ = strconv.Atoi(deploy.MetaData.Labels.RestAPIPort)
	container.Env[0].Value = k.BinaryHome()
	container.Env[1].Value = k.LinkdHome()
	container.Env[2].Value = strconv.Itoa(k.Node.Idx)
	container.Ports[0].ContainerPort = k.Node.MetaData.NodeP2PPort
	container.Ports[1].ContainerPort = k.Node.MetaData.NodeRestAPIPort
	container.Ports[2].ContainerPort = k.Node.MetaData.NodeABCIPort
	d, err := yaml.Marshal(&deploy)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	err = os.MkdirAll(k.DeployPath(), dirPermission)
	if err != nil {
		return deploy, fmt.Errorf(fmt.Sprintf("Could not prepare persist path - %s, reason - %s", k.DeployPath(), err))

	}
	k8sTemplateFilePath := fmt.Sprintf("%s/%s", k.DeployPath(), k.outputFileName())
	err = ioutil.WriteFile(k8sTemplateFilePath, d, dirPermission)
	if err != nil {
		return deploy, fmt.Errorf(fmt.Sprintf("Could not write k8s template file - %s, reason - %s", k8sTemplateFilePath, err))
	}
	return deploy, nil
}
