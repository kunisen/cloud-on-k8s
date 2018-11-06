package kibana

import (
	"github.com/elastic/stack-operators/pkg/controller/stack/common"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

const (
	//HTTPPort is the (default) port used by Kibana
	HTTPPort = 5601

	defaultImageRepositoryAndName        string = "docker.elastic.co/kibana/kibana"
	defaultTerminationGracePeriodSeconds int64  = 20

	defaultRestartPolicy            = "Always"
	defaultDNSPolicy                = "ClusterFirst"
	defaultSchedulerName            = "default-scheduler"
	defaultTerminationMessagePolicy = "File"
	defaultTerminationMessagePath   = "/dev/termination-log"
)

type PodSpecParams struct {
	Version          string
	ElasticsearchUrl string
	CustomImageName  string
}

func imageWithVersion(image string, version string) string {
	return common.Concat(image, ":", version)
}

func NewPodSpec(p PodSpecParams) corev1.PodSpec {
	imageName := p.CustomImageName
	if p.CustomImageName == "" {
		imageName = imageWithVersion(defaultImageRepositoryAndName, p.Version)
	}

	probe := &corev1.Probe{
		FailureThreshold:    3,
		InitialDelaySeconds: 10,
		PeriodSeconds:       10,
		SuccessThreshold:    1,
		TimeoutSeconds:      5,
		Handler: corev1.Handler{
			HTTPGet: &corev1.HTTPGetAction{
				Port:   intstr.FromInt(HTTPPort),
				Path:   "/",
				Scheme: corev1.URISchemeHTTP,
			},
		},
	}

	var terminationGracePeriod = defaultTerminationGracePeriodSeconds
	return corev1.PodSpec{
		TerminationGracePeriodSeconds: &terminationGracePeriod,
		RestartPolicy:                 defaultRestartPolicy,
		DNSPolicy:                     defaultDNSPolicy,
		SecurityContext:               new(corev1.PodSecurityContext),
		SchedulerName:                 defaultSchedulerName,
		Containers: []corev1.Container{{
			TerminationMessagePath:   defaultTerminationMessagePath,
			TerminationMessagePolicy: defaultTerminationMessagePolicy,
			Env: []corev1.EnvVar{
				{Name: "ELASTICSEARCH_URL", Value: p.ElasticsearchUrl},
			},
			Image:           imageName,
			ImagePullPolicy: corev1.PullIfNotPresent,
			Name:            "kibana",
			Ports: []corev1.ContainerPort{
				{Name: "http", ContainerPort: int32(HTTPPort), Protocol: corev1.ProtocolTCP},
			},
			LivenessProbe:  probe,
			ReadinessProbe: probe,
		}},
	}

}
