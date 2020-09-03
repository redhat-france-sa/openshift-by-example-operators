## Ansible Operator

This page is a walkthrough on how to develop an Operator using [Ansible](https://www.ansible.com/) with [Operator SDK](https://sdk.operatorframework.io).

It details how to do so using the new [`1.0.0` version](#operator-sdk-1.0.0) of the SDK and the older [`0.19.2` version](#operator-sdk-0.19.2).

> Make sure you have checked the [pre-requisites](../README.md#pre-requisites) before starting ;-) 

### operator-sdk 1.0.0

Start creating a new directory and initializing a new ansible operator project in it:
```
$ mkdir fruits-catalog-ansible-operator && cd fruits-catalog-ansible-operator
$ operator-sdk init --plugins=ansible --domain=com --group=redhat --version v1beta1 --kind FruitsCatalogA1
```

Check the noticeable resources within the scafolded project:
* `watches.yaml` is the file where you're gonna link your Custom Resource to the execution of an Ansible role or playbook,
* `requirements.yml` is the file allowing you to embed dependencies from the Ansible ecosystem,
* `roles/` is where we're gonna put specific Ansible tasks for our operator
* `config/` folder with your Custom Resource Definition and all the operator deployment resources,
* `Dockerfile` and `Makefile` at root.

#### Completing the operator code

Initiaze the following structure within `roles/` folder:

```
/roles
  /fruitscataloga1
    /defaults
      main.yml
    /tasks
      main.yml
    /templates
      ...
```

Edit the `watches.yaml` to replace the `FIXME` with a reference to our new role. Check the completed resources within the [`fruits-catalog-ansible-operator-1.0.0`](./fruits-catalog-ansible-operator-1.0.0) folder of this Git repository.

#### Testing in local

The Operator is running locally but is watching and realizing actions on a remote Kubernetes cluster so be sure to be connected to a cluster.

First install the CRD and then start running the operator:

```
$ make install
/Users/lbroudou/Development/go-workspace/bin/kustomize build config/crd | kubectl apply -f -
customresourcedefinition.apiextensions.k8s.io/fruitscataloga1s.redhat.com created

$ make run
/usr/local/bin/ansible-operator run
{"level":"info","ts":1597918056.000154,"logger":"cmd","msg":"Version","Go Version":"go1.13.11","GOOS":"darwin","GOARCH":"amd64","ansible-operator":"v1.0.0"}
{"level":"info","ts":1597918056.0064452,"logger":"cmd","msg":"WATCH_NAMESPACE environment variable not set. Watching all namespaces.","Namespace":""}
[...]
```

Within another terminal window, still being connected to your remote Kubernetes cluster:

```
$ kubectl create ns fruits-catalog-ansible
$ kubectl apply -f config/samples/redhat_v1beta1_fruitscataloga1.yaml -n fruits-catalog-ansible
```

#### Building and packaging the operator

Before launching `make docker-build docker-push`, export a `$IMG` variable customized to your own repository:

```
$ export IMG=quay.io/lbroudoux/fruits-catalog-ansible-operator:0.0.2
$ make docker-build docker-push IMG=$IMG
```

#### Deploying the operator

If you haven't run your operator locally first, you'll need to install the CRD within the connected Kubernetes cluster:

```
$ make install
/Users/lbroudou/Development/go-workspace/bin/kustomize build config/crd | kubectl apply -f -
customresourcedefinition.apiextensions.k8s.io/fruitscataloga1s.redhat.com created
```

> Before deploying, it may be necessary to edit the `config/default/kustomization.yaml` to change the values of `namespace` and `namePrefix`. Depending on how you name your operator, the generated names may have be longer than the 64 characters allowed in Kubernetes!

And deploy the operator image and manifests into the default namespace that will be `fruits-catalog-ans-operator-system`:

```
$ make deploy IMG=$IMG
[...]
deployment.apps/fruits-catalog-ans-operator-controller-manager created
```

#### Testing the operator on the cluster

Just deploy a new Custom Resource coming from the `samples` within another namespace (so that we keep things separated and ordered ;-)):

```
$ kubectl create ns fruits-catalog-ansible
$ kubectl apply -f config/samples/redhat_v1beta1_fruitscataloga1.yaml -n fruits-catalog-ansible
```

### operator-sdk 0.19.2

As this version of Operator SDK is now kind of obsolete, the commands are less detailed here.

Create the operator project:

```
$ mkdir fruits-catalog-ansible-operator && cd fruits-catalog-ansible-operator
$ operator-sdk new fruits-catalog-ansible-operator --api-version=redhat.com/v1alpha1 --kind=FruitsCatalogA --type=ansible
```

### Competing the operator code

Check the completed resources within the [`fruits-catalog-ansible-operator-0.19.2`](./fruits-catalog-ansible-operator-0.19.2) folder of this Git repository.

#### Testing in local

Ansible-runner module is requried for local testing. You can install and set it up with following commands:

```
$ /usr/local/Cellar/ansible/2.9.1/libexec/bin/pip install ansible-runner-http openshift 
[...]
$ ln -s /usr/local/Cellar/ansible/2.9.1/libexec/bin/ansible-runner /usr/local/bin/ansible-runner

$ ansible-runner --version
1.4.5
```

Then install the CRD in remote cluster and start operator in local:

```
$ kubectl create ns fruits-catalog-ansible
$ kubectl create -f deploy/crds/redhat.com_fruitscatalogas_crd.yaml
$ operator-sdk run local --watch-namespace fruits-catalog-ansible
```

Within another terminal window, still being connected to your remote Kubernetes cluster:

```
$ kubectl apply -f config/samples/redhat_v1beta1_fruitscataloga1.yaml -n fruits-catalog-ansible
```

#### Building and packaging the operator

```
$ operator-sdk build quay.io/lbroudoux/fruits-catalog-ansible-operator:0.0.1
$ docker push quay.io/lbroudoux/fruits-catalog-ansible-operator:0.0.1
```

#### Deploying the operator

```
kubectl create ns fruits-catalog-ansible
kubectl create -f deploy/crds/redhat.com_fruitscatalogas_crd.yaml -n fruits-catalog-ansible
kubectl create -f deploy/role.yaml -n fruits-catalog-ansible
kubectl create -f deploy/role_binding.yaml -n fruits-catalog-ansible
kubectl create -f deploy/service_account.yaml -n fruits-catalog-ansible
kubectl create -f deploy/operator.yaml -n fruits-catalog-ansible
```

#### Testing the operator on cluster

```
$ kubectl apply -f deploy/crds/redhat.com_v1alpha1_fruitscataloga_cr.yaml -n fruits-catalog-ansible 
```