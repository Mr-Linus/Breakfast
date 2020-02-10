# Breakfast

[![Go Report Card](https://goreportcard.com/badge/github.com/NJUPT-ISL/Breakfast)](https://goreportcard.com/report/github.com/NJUPT-ISL/Breakfast)

Breakfast is a kubernetes custom resource operator.
It manages the life cycle of running Machine 
Learning pods through custom controllers.

## Get Started

### Preparations

- Kubernetes v1.15+
- [SCV](https://github.com/NJUPT-ISL/SCV) has been deployed.
- [Yoda-Scheduler](https://github.com/NJUPT-ISL/Yoda-Scheduler) has been deployed.

### Deploy Breakfast
```shell
kubectl apply -f https://raw.githubusercontent.com/NJUPT-ISL/Breakfast/master/deploy/breakfast.yaml
```

### Create Machine Learning Task
- Create a ML task with Tensorflow Framework
```yaml
apiVersion: core.run-linux.com/v1alpha1
kind: Bread
metadata:
  name: test
  namespace: root
spec:
  scv:
    gpu: "1"
    memory: "4000"
    level: "Medium"
  framework:
    name: "tensorflow"
    version: "2.0"
  task:
    type: train
    path: "/root"
    command: "python /root/test.py"
```
- Create a ML task with ssh service
```yaml
apiVersion: core.run-linux.com/v1alpha1
kind: Bread
metadata:
  name: test
spec:
  scv:
    gpu: "1"
    memory: "4000"
    level: "Medium"
  framework:
    name: "tensorflow"
    version: "2.0"
  task:
    type: ssh
    path: "/root"
    command: ""
```

### How to Dev
- Build the Breakfast
```shell script
make 
```
- Build the Docker Image
```shell script
make docker-build
```
- Debug the Breakfast Controller
```shell script
make run
```

## Contact us
![QQ Group](https://github.com/NJUPT-ISL/Breakfast/blob/master/img/qrcode_1581334380545.jpg)