package kubernetes

import (
	"errors"
	log "github.com/Sirupsen/logrus"
	"github.com/qa-dev/jsonwire-grid/jsonwire"
	"k8s.io/client-go/kubernetes"
	apiV1 "k8s.io/client-go/pkg/api/v1"
	"net"
	"strconv"
	"time"
)

type kubernetesProviderInterface interface {
	Create(podName string, nodeParams nodeParams) error
	Destroy(podName string) error
}

type kubernetesProvider struct {
	clientset     *kubernetes.Clientset
	namespace     string
	clientFactory jsonwire.ClientFactoryInterface
}

func (p *kubernetesProvider) Create(podName string, nodeParams nodeParams) error {
	pod := &apiV1.Pod{}
	pod.ObjectMeta.Name = podName
	pod.ObjectMeta.Labels = map[string]string{"name": podName}
	pod.Spec.Hostname = podName
	container := apiV1.Container{}
	container.Name = podName
	container.Image = nodeParams.Image
	port, err := strconv.Atoi(nodeParams.Port)
	if err != nil {
		return errors.New("convert to int nodeParams.Port, " + err.Error())
	}
	container.Ports = []apiV1.ContainerPort{{ContainerPort: int32(port)}}
	pod.Spec.Containers = append(pod.Spec.Containers, container)
	_, err = p.clientset.CoreV1Client.Pods(p.namespace).Create(pod)
	if err != nil {
		return errors.New("send command pod/create to k8s, " + err.Error())
	}

	service := &apiV1.Service{}
	service.ObjectMeta.Name = podName
	service.Spec.ClusterIP = "None"
	service.Spec.Ports = []apiV1.ServicePort{{Port: int32(port)}}
	service.Spec.Selector = map[string]string{"name": podName}
	_, err = p.clientset.CoreV1Client.Services(p.namespace).Create(service)
	if err != nil {
		return errors.New("send command service/create to k8s, " + err.Error())
	}

	// todo: пока так ожидаем поднятие ноды, так как не понятно что конкретно означают статусы возвращаемые через апи
	client := p.clientFactory.Create(net.JoinHostPort(podName, nodeParams.Port))
	stop := time.After(40 * time.Second)
	log.Debugln("start waiting")
Loop:
	for {
		select {
		case <-stop:
			return errors.New("wait stopped by timeout")
		default:
			time.Sleep(time.Second)
			log.Debugln("start request")
			message, err := client.Health()
			if err != nil {
				log.Debugf("fail request, %v", err)
				continue
			}
			log.Debugf("done request, status: %v", message.Status)
			if message.Status == 0 {
				break Loop
			}
		}
	}

	return nil
}

func (p *kubernetesProvider) Destroy(podName string) error {
	err := p.clientset.CoreV1Client.Pods(p.namespace).Delete(podName, &apiV1.DeleteOptions{})
	if err != nil {
		err = errors.New("send command pod/delete to k8s, " + err.Error())
		return err
	}
	err = p.clientset.CoreV1Client.Services(p.namespace).Delete(podName, &apiV1.DeleteOptions{})
	if err != nil {
		err = errors.New("send command service/delete to k8s, " + err.Error())
		return err
	}
	return nil
}
