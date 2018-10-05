package main

import (
	"bufio"
	"flag"
	"log"
	"os"
)

var (
	domainName  = flag.String("domain", "", "Domain Name")
	serviceName = flag.String("svc", "", "Service Name")
	servicePort = flag.String("port", "", "Service Port")
)

func init() {
	flag.Parse()
}

func main() {
	var err error
	var deployYAMLString, serviceYAMLString string
	newline := "\n"

	// generate deploy.yaml string
	deployYAMLString += "# this file for deployment " + *serviceName + newline
	deployYAMLString += "apiVersion: extensions/v1beta1" + newline
	deployYAMLString += "kind: Deployment" + newline
	deployYAMLString += "metadata:" + newline
	deployYAMLString += "  name: " + *serviceName + newline
	deployYAMLString += "  labels:" + newline
	deployYAMLString += "    app: " + *serviceName + newline
	deployYAMLString += "spec:" + newline
	deployYAMLString += "  replicas: 1" + newline
	deployYAMLString += "  template:" + newline
	deployYAMLString += "    metadata:" + newline
	deployYAMLString += "      labels:" + newline
	deployYAMLString += "        app: " + *serviceName + newline
	deployYAMLString += "    spec:" + newline
	deployYAMLString += "      containers:" + newline
	deployYAMLString += "      - name: " + *serviceName + newline
	deployYAMLString += "        image: [replace]" + newline
	deployYAMLString += "        imagePullPolicy: Always" + newline
	deployYAMLString += "        volumeMounts:" + newline
	deployYAMLString += "        - mountPath: " + "/go/src/" + *serviceName + "/logs" + newline
	deployYAMLString += "          name: " + *serviceName + "-logs" + newline
	deployYAMLString += "        env:" + newline
	deployYAMLString += "        - name: [replace]" + newline
	deployYAMLString += "          value: \"[replace]\"" + newline
	deployYAMLString += "      restartPolicy: Always" + newline
	deployYAMLString += "      volumes:" + newline
	deployYAMLString += "      - name: " + *serviceName + "-logs" + newline
	deployYAMLString += "        hostPath:" + newline
	deployYAMLString += "          path: /data/" + *domainName + "_" + *serviceName + "/logs" + newline
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
	serviceYAMLString += "apiVersion: extensions/v1beta1" + newline
	serviceYAMLString += "kind: Service" + newline
	serviceYAMLString += "metadata:" + newline
	serviceYAMLString += "  name: " + *serviceName + newline
	serviceYAMLString += "  labels:" + newline
	serviceYAMLString += "    app: " + *serviceName + newline
	serviceYAMLString += "spec:" + newline
	serviceYAMLString += "  ports:" + newline

	// port HTTP
	serviceYAMLString += "  - port: " + *servicePort + newline
	serviceYAMLString += "    targetPort: " + *servicePort + newline
	serviceYAMLString += "    # nodePort: [replace]" + newline
	serviceYAMLString += "    name: app" + newline
	// port GRPC
	serviceYAMLString += "  - port: 5" + *servicePort + newline
	serviceYAMLString += "    targetPort: 5" + *servicePort + newline
	serviceYAMLString += "    # nodePort: [replace]" + newline
	serviceYAMLString += "    name: grpc" + newline

	serviceYAMLString += "  type: LoadBalancer" + newline
	serviceYAMLString += "  selector:" + newline
	serviceYAMLString += "    app: " + *serviceName + newline

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
