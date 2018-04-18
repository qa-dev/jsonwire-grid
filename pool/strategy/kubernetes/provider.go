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
	"strings"
	"fmt"
)

type kubernetesProviderInterface interface {
	Create(podName string, nodeParams nodeParams) (nodeAddress string, err error)
	// idempotent operation
	Destroy(podName string) error
}

type kubDnsProvider struct {
	clientset     *kubernetes.Clientset
	namespace     string
	clientFactory jsonwire.ClientFactoryInterface
}

func (p *kubDnsProvider) Create(podName string, nodeParams nodeParams) (nodeAddress string, err error) {
	pod := &apiV1.Pod{}
	pod.ObjectMeta.Name = podName
	pod.ObjectMeta.Labels = map[string]string{"name": podName}
	pod.Spec.Hostname = podName
	container := apiV1.Container{}
	container.Name = podName
	container.Image = nodeParams.Image
	port, err := strconv.Atoi(nodeParams.Port)
	if err != nil {
		return "", errors.New("convert to int nodeParams.Port, " + err.Error())
	}
	container.Ports = []apiV1.ContainerPort{{ContainerPort: int32(port)}}
	pod.Spec.Containers = append(pod.Spec.Containers, container)
	_, err = p.clientset.CoreV1Client.Pods(p.namespace).Create(pod)
	if err != nil {
		return "", errors.New("send command pod/create to k8s, " + err.Error())
	}

	stopWaitIP := time.After(40 * time.Second)
	log.Debugf("start waiting pod ip")
	var createdPodIP string
LoopWaitIP:
	for {
		select {
		case <-stopWaitIP:
			return "", fmt.Errorf("wait podIP stopped by timeout, %v", podName)
		default:
			time.Sleep(time.Second)
			createdPod, err := p.clientset.CoreV1Client.Pods(p.namespace).Get(podName)
			if err != nil {
				log.Debugf("fail get created pod, %v, %v",podName, err)
				continue
			}
			if createdPod.Status.PodIP == "" {
				log.Debugf("empty pod ip, %v", podName)
				continue
			}
			createdPodIP = createdPod.Status.PodIP
			break LoopWaitIP
		}
	}

	// todo: пока так ожидаем поднятие ноды, так как не понятно что конкретно означают статусы возвращаемые через апи
	nodeAddress = net.JoinHostPort(createdPodIP, nodeParams.Port)
	client := p.clientFactory.Create(nodeAddress)
	stop := time.After(40 * time.Second)
	log.Debugln("start waiting selenium")
LoopWaitSelenium:
	for {
		select {
		case <-stop:
			return "", fmt.Errorf("wait selenium stopped by timeout, %v", podName)
		default:
			time.Sleep(time.Second)
			message, err := client.Health()
			if err != nil {
				log.Debugf("fail request, %v", err)
				continue
			}
			log.Debugf("done request, status: %v", message.Status)
			if message.Status == 0 {
				break LoopWaitSelenium
			}
		}
	}

	return nodeAddress, nil
}

//Destroy - destroy all pod data (idempotent operation)
func (p *kubDnsProvider) Destroy(podName string) error {
	err := p.clientset.CoreV1Client.Pods(p.namespace).Delete(podName, &apiV1.DeleteOptions{})
	switch {
	case err != nil && strings.Contains(err.Error(), "not found"):
		// pod already deleted
	case err != nil:
		err = errors.New("send command pod/delete to k8s, " + err.Error())
		return err
	}
	return nil
}
