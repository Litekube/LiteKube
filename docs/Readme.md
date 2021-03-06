<h1 align="center">Tutorial</h1>

# Catalogue
- [Catalogue](#catalogue)
- [Build](#build)
  - [Simple](#simple)
  - [Cross-compile](#cross-compile)
  - [By dockerfile](#by-dockerfile)
- [Deployment](#deployment)
- [Usage](#usage)
- [Test LiteKube](#test-litekube)
# Build

## Simple
you can build  simplely by `go build .` [simple-script](https://github.com/Litekube/LiteKube/blob/main/scripts/build/build.sh) is provided to build binaries for you. It will auto build all components for your local-enviroment into `build/outputs`. Of course, `golang` and `gcc` environment are needed.

## Cross-compile
`LiteKube`need to set `CGO_ENABLED=1` . If you are compiling for arm architecture, set `GOARM=7` additionally when necessary and `GOARM=6` is `golang-default`.

## By [dockerfile](../build/Dockerfile)

We also provide a [Dockerfile](../build/Dockerfile) to help simplify operations or as a reference, you can run by:

> assum you start your work in /mywork/

1. download code from github

    ```shell
    cd /mywork
    git clone https://github.com/Litekube/LiteKube.git 
    ```

2. build image by docker

    ```shell
    cd /mywork/LiteKube/build/
    docker build -t litekube/centos-go:v1 .
    ```

    if you need proxy, you can use proxy of your host-device and run:

    ```shell
    cd /mywork/LiteKube/build/
    export http_proxy="your proxy"
    export https_proxy="your proxy"
    docker build --network=host -t litekube/centos-go:v1 .
    ```

3. start to build binaries for LiteKube

    ```shell
    chmod +x /mywork/LiteKube/scripts/build/build.sh
    docker run -v /mywork/LiteKube:/LiteKube --name=compile-litekube litekube/centos-go:v1 /LiteKube/scripts/build/build.sh
    ```

    now, you can view binaries in `/mywork/LiteKube/build/outputs/`. 
    
    > we only provide two version in this container. 
    >
    > - the same arch with your machine for Linux
    > - `Armv7l` for Linux

# Deployment

**Notice**

- `network-controller`and `kine` can run in `leader` for default. They can also run in separate nodes or replace kine with `ETCD Cluster` by set `global.run-network-manager=false` and `global.run-kine=false` . As a cost, you need to set corresponding parameters for them.
- `build-in worker` for `leader` is also allowed but we set it disabled, you can enable by set `global.enable-worker=true`. Note that you will additionally need to provide `leader` with the same running environment as the `worker` if you do this.

**Components**
> If you have proxy configured on your nodes, be sure to delete the proxy rule for 10.1.1.0/24.

- [network-controller](https://github.com/Litekube/network-controller)
- [ncadm](https://github.com/Litekube/network-controller/blob/main/docs/ncadm-explain.md)
- [Kine](https://github.com/Litekube/kine) (you can also use `ETCD` cluster instead)
- [leader](leader/deploy.md)
- [worker](worker/deploy.md)
- [kubectl](kubectl/deploy.md) (no change to kubectl in kubernetes)
- [likuadm](likuadm/deploy.md)

# Usage
- [network-controller](https://github.com/Litekube/network-controller/blob/main/docs/demo-usage.md)
- [ncadm](https://github.com/Litekube/network-controller/blob/main/docs/ncadm-explain.md)
- [leader](leader/usage.md)
- [worker](worker/usage.md)
- [kubectl](https://github.com/kubernetes/kubectl)
- [likuadm](likuadm/usage.md)

# Test LiteKube

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  labels:
    app: test-litekube
  name: test-deployment
spec:
  replicas: 1 
  selector:
    matchLabels:
      app: test-litekube
  template:
    metadata:
      labels:
        app: test-litekube
    spec:
      containers:
      - name: nginx
        image: nginx:1.17.1
        imagePullPolicy: IfNotPresent
        ports:
        - containerPort: 80
      - name: centos 
        image: centos:7
        command: ["sh","-c",'sleep 3600s']
# test: kubectl exec -it test-deployment-*  -c centos  -- /usr/bin/bash

---

apiVersion: v1
kind: Service
metadata:
  labels:
    app: test-litekube
  name: test-nginx-service-nodeport
spec:
  ports:
  - port: 80
    protocol: TCP
    nodePort: 30001
  selector:
    app: test-litekube
  type: NodePort

---

apiVersion: v1
kind: Service
metadata:
  labels:
    app: test-litekube
  name: test-nginx-service-clusterip
spec:
  ports:
  - port: 80
    protocol: TCP
    targetPort: 80
  selector:
    app: test-litekube
  type: ClusterIP

# test: 
# - curl in host: http://{node-ip}:30001
# - curl in pod: http://{cluster-ip}:80
```
