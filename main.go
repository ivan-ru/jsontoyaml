package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"io/ioutil"
	"log"
	"os"
)

// Deploy ...
type Deploy struct {
	APIVersion   string `json:"apiVersion"`
	Domain       string `json:"domain"`
	Name         string `json:"name"`
	Image        string `json:"image"`
	Environments []struct {
		Name  string `json:"name"`
		Value string `json:"value"`
	} `json:"environments"`
}

// Service ...
type Service struct {
	APIVersion string `json:"apiVersion"`
	Name       string `json:"name"`
	Ports      []struct {
		Port       string `json:"port"`
		TargetPort string `json:"targetPort"`
		NodePort   string `json:"nodePort"`
		Name       string `json:"name"`
	} `json:"ports"`
}

func init() {
	flag.Parse()
}

func main() {
	var byteDeploy, byteService []byte
	var dataDeploy Deploy
	var dataService Service
	var err error

	byteDeploy, _ = ioutil.ReadFile("deploy.json")
	byteService, _ = ioutil.ReadFile("service.json")

	err = json.Unmarshal(byteDeploy, &dataDeploy)
	if err != nil {
		log.Fatal("error unmarshal json", err)
		return
	}
	err = json.Unmarshal(byteService, &dataService)
	if err != nil {
		log.Fatal("error unmarshal json", err)
		return
	}
	var deployYAMLString, serviceYAMLString string
	newline := "\n"

	// generate deploy.yaml string
	deployYAMLString += "# this file for deployment " + dataDeploy.Name + newline
	deployYAMLString += "apiVersion: " + dataDeploy.APIVersion + newline
	deployYAMLString += "kind: Deployment" + newline
	deployYAMLString += "metadata:" + newline
	deployYAMLString += "  name: " + dataDeploy.Name + newline
	deployYAMLString += "  labels:" + newline
	deployYAMLString += "    app: " + dataDeploy.Name + newline
	deployYAMLString += "spec:" + newline
	deployYAMLString += "  replicas: 1" + newline
	deployYAMLString += "  template:" + newline
	deployYAMLString += "    metadata:" + newline
	deployYAMLString += "      labels:" + newline
	deployYAMLString += "        app: " + dataDeploy.Name + newline
	deployYAMLString += "    spec:" + newline
	deployYAMLString += "      containers:" + newline
	deployYAMLString += "      - name: " + dataDeploy.Name + newline
	deployYAMLString += "        image: " + dataDeploy.Image + newline
	deployYAMLString += "        imagePullPolicy: Always" + newline
	deployYAMLString += "        volumeMounts:" + newline
	deployYAMLString += "        - mountPath: " + "/go/src/" + dataDeploy.Name + "/logs" + newline
	deployYAMLString += "          name: " + dataDeploy.Name + "-logs" + newline
	deployYAMLString += "        env:" + newline
	for _, v := range dataDeploy.Environments {
		deployYAMLString += "        - name: " + v.Name + newline
		deployYAMLString += "          value: \"" + v.Value + "\"" + newline
	}
	deployYAMLString += "      restartPolicy: Always" + newline
	deployYAMLString += "      volumes:" + newline
	deployYAMLString += "      - name: " + dataDeploy.Name + "-logs" + newline
	deployYAMLString += "        hostPath:" + newline
	deployYAMLString += "          path: /data/" + dataDeploy.Domain + "_" + dataDeploy.Name + "/logs" + newline
	deployYAMLString += "      imagePullSecrets:" + newline
	deployYAMLString += "      - name: regsecret" + newline
	deployYAMLString += "  strategy:" + newline
	deployYAMLString += "    type: RollingUpdate" + newline
	deployYAMLString += "    rollingUpdate:" + newline
	deployYAMLString += "      maxUnavailable: 1" + newline
	deployYAMLString += "      maxSurge: 1" + newline

	// write deploy
	fDeploy, err := os.Create("result/deploy.yaml")
	if err != nil {
		log.Fatal("error create file", err)
		return
	}

	defer fDeploy.Close()
	wDeploy := bufio.NewWriter(fDeploy)
	_, err = wDeploy.WriteString(deployYAMLString)
	if err != nil {
		log.Fatal("error write to file deploy.yaml", err)
		return
	}

	wDeploy.Flush()

	// generate service.yaml string
	serviceYAMLString += "apiVersion: " + dataService.APIVersion + newline
	serviceYAMLString += "kind: Service" + newline
	serviceYAMLString += "metadata:" + newline
	serviceYAMLString += "  name: " + dataService.Name + newline
	serviceYAMLString += "  labels:" + newline
	serviceYAMLString += "    app: " + dataService.Name + newline
	serviceYAMLString += "spec:" + newline
	serviceYAMLString += "  ports:" + newline

	for _, v := range dataService.Ports {
		serviceYAMLString += "  - port: " + v.Port + newline
		serviceYAMLString += "    targetPort: " + v.TargetPort + newline
		serviceYAMLString += "    nodePort: " + v.NodePort + newline
		serviceYAMLString += "    name: " + v.Name + newline
	}

	serviceYAMLString += "  type: LoadBalancer" + newline
	serviceYAMLString += "  selector:" + newline
	serviceYAMLString += "    app: " + dataService.Name + newline

	// write service
	fService, err := os.Create("result/service.yaml")
	if err != nil {
		log.Fatal("error create file service.yaml", err)
		return
	}

	defer fService.Close()
	wService := bufio.NewWriter(fService)
	_, err = wService.WriteString(serviceYAMLString)
	if err != nil {
		log.Fatal("error write to file service.yaml", err)
		return
	}

	wService.Flush()
}
