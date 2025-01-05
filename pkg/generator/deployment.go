package generator

import (
	"fmt"
	"mimic/pkg/config/mimic"

	//"gopkg.in/yaml.v3"
	v1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/yaml"
)

type Deployment struct {
	Name    string          `json:"name,omitempty"`
	Ingress []mimic.Ingress `json:"ingress,omitempty"`
	Egress  []mimic.Egress  `json:"egress,omitempty"`
}

func (d *Deployment) Generate(namespace string, image string) []byte {
	deployment := v1.Deployment{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "apps/v1",
			Kind:       "Deployment",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      d.Name,
			Namespace: namespace,
		},
		Spec: v1.DeploymentSpec{
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"app": d.Name,
				},
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"app": d.Name,
					},
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						corev1.Container{
							Name:    d.Name,
							Image:   image,
							Command: []string{},
							Env: []corev1.EnvVar{
								{
									Name:  "CONFIG",
									Value: "/config/mimic.yml",
								},
							},
							VolumeMounts: []corev1.VolumeMount{
								{
									Name:      "config",
									MountPath: "/config",
								},
							},
						},
					},
					Volumes: []corev1.Volume{
						{
							Name: "config",
							VolumeSource: corev1.VolumeSource{
								ConfigMap: &corev1.ConfigMapVolumeSource{
									LocalObjectReference: corev1.LocalObjectReference{
										Name: d.Name,
									},
									Items: []corev1.KeyToPath{
										{
											Key:  "mimic.yml",
											Path: "mimic.yml",
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}

	for _, ingress := range d.Ingress {
		if ingress.HTTP != nil {
			ports := deployment.Spec.Template.Spec.Containers[0].Ports
			ports = append(ports, corev1.ContainerPort{
				Name:          "http",
				ContainerPort: ingress.HTTP.Port,
			})
			deployment.Spec.Template.Spec.Containers[0].Ports = ports
		}
	}

	data, err := yaml.Marshal(&deployment)
	if err != nil {
		fmt.Println(err)
	}
	return data
}

func (d *Deployment) GenerateConfigMap(ns string) []byte {
	conf := mimic.Config{
		Mimic: mimic.Mimic{
			Ingress: d.Ingress,
			Egress:  d.Egress,
		},
	}

	data, err := yaml.Marshal(&conf)
	if err != nil {
		fmt.Println(err)
	}

	cm := corev1.ConfigMap{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "v1",
			Kind:       "ConfigMap",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      d.Name,
			Namespace: ns,
		},
		Data: map[string]string{
			"mimic.yml": string(data),
		},
	}
	data, err = yaml.Marshal(&cm)
	if err != nil {
		fmt.Println(err)
	}
	return data
}

func (d *Deployment) GenerateService(ns string) []byte {
	svc := corev1.Service{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "v1",
			Kind:       "Service",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      d.Name,
			Namespace: ns,
		},
		Spec: corev1.ServiceSpec{
			Selector: map[string]string{
				"app": d.Name,
			},
		},
	}

	for _, ingress := range d.Ingress {
		if ingress.HTTP != nil {
			svc.Spec.Ports = append(svc.Spec.Ports, corev1.ServicePort{
				Name:     "http",
				Protocol: "TCP",
				Port:     ingress.HTTP.Port,
			})
		}
	}

	data, err := yaml.Marshal(&svc)
	if err != nil {
		fmt.Println(err)
	}
	return data
}
