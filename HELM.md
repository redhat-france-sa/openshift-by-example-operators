## Helm Operator

This page is a walkthrough on how to develop an Operator using [Helm Charts](https://helm.sh) with [Operator SDK](https://sdk.operatorframework.io).

It details how to do so using the new [`1.0.0` version](#operator-sdk-1.0.0) of the SDK and the older [`0.19.2` version](#operator-sdk-0.19.2).

> Make sure you have checked the [pre-requisites](../README.md#pre-requisites) before starting ;-) 

### operator-sdk 1.0.0

Start creating a new directory and initializing a new helm operator project in it:
```
$ mkdir fruits-catalog-helm-operator && cd fruits-catalog-helm-operator
$ operator-sdk init --plugins=helm --domain=com --group=redhat --version v1beta1 --kind FruitsCatalogH1
```

Check the noticeable resources within the scafolded project:
* `watches.yaml` is the file where you're gonna link your Custome Resource to the execution of a Helm Chart,
* `helm-charts/fruitscatalogh1/` is where you'll define your Chart resources. You'll find typical `Chart.yaml` and `values.yaml` with all the Kubernetes resources templates,
* `config/` folder with your Custom Resource Definition and all the operator deployment resources,
* `Dockerfile` and `Makefile` at root.

#### Completing the operator code

* Edit the `helm-charts/fruitscatalogh1/values.yaml` to reflect your custom resource default values,
* Edit all the resources within `helm-charts/fruitscatalogh1/templates` resources accordingly,
* Adapt the required permsissions within `config/rbac/role.yaml`

Check the completed resources within the [`fruits-catalog-helm-operator-1.0.0`](./fruits-catalog-helm-operator-1.0.0) folder of this Git repository.

#### Building and packaging the operator

Before launching `make docker-build docker-push`, export a `$IMG` variable customized to your own repository:
```
$ export IMG=quay.io/lbroudoux/fruits-catalog-helm-operator:0.0.2
$ make docker-build docker-push IMG=$IMG
```

#### Installing the operator

Now, connected to the Kubernetes cluster, install the CRD within cluster:
```
$ make install
[...]
customresourcedefinition.apiextensions.k8s.io/fruitscataloghs.redhat.com configured
```

And deploy the operator image and manifests into the default namespace that will be `fruits-catalog-helm-operator-system`:
```
$ make deploy IMG=$IMG
[...]
deployment.apps/fruits-catalog-helm-operator-controller-manager created
```

#### Testing the operator on cluster

Just deploy a new Custom Resource coming from the `samples` within another namespace (so that we keep things separated and ordered ;-)):

```
$ kubectl create ns fruits-catalog-helm
$ kubectl apply -f config/samples/redhat_v1beta1_fruitscatalogh1.yaml -n fruits-catalog-helm
```

> Does not work! Why ? Normal. Because Helm is not that smart ;-) With this first sample CR, we are asking Helm to generate a database user/password for us because we want the operator to be simple to use. However, for each and every reconcilitation loop the operator is recreating the Secret!! So there's big chances the dabatase and the application Pod are not using the same credentials and thus cannot talk to each.

Using Helm as the Operator implementation technology, you cannot have random value generation. Helm is not able of discovering what has already been deployed to the cluster and only applying changes when required. It re-applies all the resources systematically.

So let's change our CR to another sample having username and password embedded and now it works:

```
$ kubectl delete fruitscatalogh1.redhat.com example-fruitscatalogh -n fruits-catalog-helm
$ kubectl apply -f config/samples/redhat_v1beta1_fruitscatalogh1-username+password.yaml -n fruits-catalog-helm
```

### operator-sdk 0.19.2

As this version of Operator SDK is now kind of obsolete, the commands are less detailed here.

Create the operator project:

```
$ mkdir fruits-catalog-helm-operator && cd fruits-catalog-helm-operator
$ operator-sdk new fruits-catalog-helm-operator --api-version=redhat.com/v1alpha1 --kind=FruitsCatalogH --type=helm
```

### Competing the operator code

Check the completed resources within the [`fruits-catalog-helm-operator-0.19.2`](./fruits-catalog-helm-operator-0.19.2) folder of this Git repository.

#### Building and packaging the operator

```
$ operator-sdk build quay.io/lbroudoux/fruits-catalog-helm-operator:0.0.1
$ docker push quay.io/lbroudoux/fruits-catalog-helm-operator:0.0.1
```

#### Deploying the operator

```
kubectl create ns fruits-catalog-helm
kubectl create -f deploy/crds/redhat.com_fruitscataloghs_crd.yaml -n fruits-catalog-helm
kubectl create -f deploy/role.yaml -n fruits-catalog-helm
kubectl create -f deploy/role_binding.yaml -n fruits-catalog-helm
kubectl create -f deploy/service_account.yaml -n fruits-catalog-helm
kubectl create -f deploy/operator.yaml -n fruits-catalog-helm
```

#### Testing the operator on cluster

```
$ kubectl create -f deploy/crds/redhat.com_v1alpha1_fruitscatalogh_cr-username+password.yaml -n fruits-catalog-helm
```