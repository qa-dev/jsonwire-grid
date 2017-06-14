package main

import (
	"flag"
	"fmt"
	log "github.com/Sirupsen/logrus"
	"os"
	"os/exec"
	"os/signal"
	"time"
)

func main() {
	hubUrl := flag.String("hub", "http://127.0.0.1:4444", "address of hub, default http://127.0.0.1:4444")
	startPort := flag.Int("startPort", 5000, "start port")
	maxDuration := flag.Int("maxDuration", 100, "request duration [0 <=duration], default 0")
	countNodes := flag.Int("countNodes", 100, "count nodes")
	flag.Parse()
	log.Infof("hub url: %v", *hubUrl)
	log.Infof("startPort: %v", *startPort)
	log.Infof("maxDuration: %v", *maxDuration)
	log.Infof("countNodes: %v", *countNodes)

	stop := make(chan os.Signal)
	signal.Notify(stop, os.Interrupt)
	done := make(chan struct{})
	isAlive := true

	go func() {
		<-stop
		isAlive = false
		done <- struct{}{}
	}()

	for port := *startPort; port < *startPort+*countNodes && isAlive; port++ {
		time.Sleep(time.Millisecond * 5)
		cmd := exec.Command(
			"webdriver-node-mock",
			fmt.Sprintf("-hub=%v", *hubUrl),
			fmt.Sprintf("-port=%v", port),
			fmt.Sprintf("-maxDuration=%v", *maxDuration),
		)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		err := cmd.Start()
		if err != nil {
			log.Error(err)
		}

		//wait interrupt child process
		defer func(cmd *exec.Cmd) {
			cmd.Process.Signal(os.Interrupt)
			cmd.Process.Wait()
		}(cmd)
		log.Info("Created instance #", port-*startPort+1)
	}
	<-done

}
