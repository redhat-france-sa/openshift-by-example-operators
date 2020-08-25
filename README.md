# openshift-by-example-operators

## Pre-requisites

* `operator-sdk` should be installed and present into your `$PATH`. We use latest `1.0.0` version and previous `0.19.2`version. Just check https://master.sdk.operatorframework.io/docs/installation/install-operator-sdk/ on how to install it,
* `docker` (or another tool compatible with multi-stage Dockerfiles) should be installed and present into your `$PATH`. Minimum version is `17.03+`,
* `go` should be installed and present into your `$PATH`. Minimum version is Go `1.13`,
* `kubectl` shoud be installed and present into your `$PATH`. `v1.16.0+` is the minimum version,

## Helm Operator

### operator-sdk 1.0.0

```
$ mkdir fruits-catalog-helm-operator && cd fruits-catalog-helm-operator
$ operator-sdk init --plugins=helm --domain=com --group=redhat --version v1beta1 --kind FruitsCatalogH1
```

```
$ export IMG=quay.io/lbroudoux/fruits-catalog-helm-operator:0.0.2
$ make docker-build docker-push IMG=$IMG
```

#### Installing the operator

```
$ make install
[...]
customresourcedefinition.apiextensions.k8s.io/fruitscataloghs.redhat.com configured
```

```
$ make deploy IMG=$IMG
[...]
deployment.apps/fruits-catalog-helm-operator-controller-manager created
```

#### Testing the operator on cluster

```
$ kubectl create ns fruits-catalog-helm-1
$ kubectl apply -f config/samples/redhat_v1beta1_fruitscatalogh1.yaml -n fruits-catalog-helm-1
```

> Does not work! Why ? Normal ;-) 

```
$ kubectl apply -f config/samples/redhat_v1beta1_fruitscatalogh1-username+password.yaml -n fruits-catalog-helm-1
```

### operator-sdk 0.19.2

## Ansible Operator

### operator-sdk 1.0.0

### operator-sdk 0.19.2

## Go Operator

### operator-sdk 1.0.0

### operator-sdk 0.19.2
