## Go Operator

This page is a walkthrough on how to develop an Operator using [Go language](https://golang.org/) with [Operator SDK](https://sdk.operatorframework.io).

It details how to do so using the new [`1.0.0` version](#operator-sdk-1.0.0) of the SDK and the older [`0.19.2` version](#operator-sdk-0.19.2).

> Make sure you have checked the [pre-requisites](../README.md#pre-requisites) before starting ;-) 

### operator-sdk 1.0.0

Start creating a new directory and initializing a new Go operator project in it. This time you'll need to specify a Git repository in order to generate coherent Go package names:

```
$ mkdir fruits-catalog-go-operator && cd fruits-catalog-go-operator
$ operator-sdk init --domain=com --repo=github.com/redhat-france-sa/openshift-by-example-operators/fruits-catalog-go-operator
[...]
Next: define a resource with:
$ operator-sdk create api
```

As you seen reading the output, this is a 2 steps process: you'll need now to generate at least an API (a CRD and a controller) to define your operator. Let's do that:

```
$ operator-sdk create api --group=redhat --version=v1beta1 --kind=FruitsCatalogG1 --resource=true --controller=true
```

Check the noticeable resources within the scafolded project:
* `api/` folder with your CRD Go types definitions,
* `config/` folder with your CRD manifests and all the operator deployment resources,
* `controllers/` with the Go code corresponding to your API controllers,
* `main.go` and `go.mod` at root,
* `Dockerfile` and `Makefile` at root.

#### Completing the operator code

Completing the operator code consists in 2 steps:
* Defining the Go structures that will handle your Custom Resource Definitions (both the `Spec` and the `Status` parts),
* Defining the reconcialiation loop logic withn a Controller.

Start editing the `api/v1beta1/fruitscatalogg1_types.go` file, completing the structure definitions to suit your CRD design. Once it is done - and each time you'll do a cahnge in the future, you will have to run the `make` command to generate some utility Go functions and register API within main program:

```
$ make
/Users/lbroudou/Development/go-workspace/bin/controller-gen object:headerFile="hack/boilerplate.go.txt" paths="./..."
go fmt ./...
go vet ./...
go build -o bin/manager main.go
```

Then develop your reconciliation loop logic within the `controllers/fruitscatalogg1_controller.go` file. This implies:
* Adjusting RBAC persmisions using `+kubebuilder` annotations,
* Puting some comments on the `Reconcile` function to ensure it is actually exported,
* Adding the logic within the fonction + all the Kubernetes resources management stuffs,
* Complete the `manager` definition with `Owns()` directives to add watches on dependant resources,
* Complete `main.go` file with all the scheme inclusions for dependencies. Ex: `utilruntime.Must(routev1.AddToScheme(scheme))`

Check the completed resources within the [`fruits-catalog-go-operator-1.0.0`](./fruits-catalog-go-operator-1.0.0) folder of this Git repository.

#### Testing in local

The Operator is running locally but is watching and realizing actions on a remote Kubernetes cluster so be sure to be connected to a cluster.

First install the CRD and then start running the operator:

```
$ make install                                                                                                               
/Users/lbroudou/Development/go-workspace/bin/controller-gen "crd:trivialVersions=true,crdVersions=v1" rbac:roleName=manager-role webhook paths="./..." output:crd:artifacts:config=config/crd/bases
/Users/lbroudou/Development/go-workspace/bin/kustomize build config/crd | kubectl apply -f -
customresourcedefinition.apiextensions.k8s.io/fruitscatalogg1s.redhat.com created
```

Then you can launth the operator using a local process:

```
$ make run                                                                                                                                  
/Users/lbroudou/Development/go-workspace/bin/controller-gen object:headerFile="hack/boilerplate.go.txt" paths="./..."
go fmt ./...
go vet ./...
/Users/lbroudou/Development/go-workspace/bin/controller-gen "crd:trivialVersions=true,crdVersions=v1" rbac:roleName=manager-role webhook paths="./..." output:crd:artifacts:config=config/crd/bases
go run ./main.go
I0820 15:15:39.808141   42998 request.go:621] Throttling request took 1.025789279s, request: GET:https://api.cluster-1286.1286.example.opentlc.com:6443/apis/apiregistration.k8s.io/v1?timeout=32s
2020-08-20T15:15:41.136+0200	INFO	controller-runtime.metrics	metrics server is starting to listen	{"addr": ":8080"}
2020-08-20T15:15:41.136+0200	INFO	setup	starting manager
2020-08-20T15:15:41.136+0200	INFO	controller-runtime.manager	starting metrics server	{"path": "/metrics"}
[...]
```

Within another terminal window, still being connected to your remote Kubernetes cluster:

```
$ kubectl create ns fruits-go-ansible
$ kubectl apply -f config/samples/redhat_v1beta1_fruitscatalogg1.yaml -n fruits-catalog-go
```

#### Building and packaging the operator

Before launching `make docker-build docker-push`, export a `$IMG` variable customized to your own repository:

```
$ export IMG=quay.io/lbroudoux/fruits-catalog-go-operator:0.0.2
$ make docker-build docker-push IMG=$IMG
```

#### Deploying the operator

If you haven't run your operator locally first, you'll need to install the CRD within the connected Kubernetes cluster:

```
$ make install
/Users/lbroudou/Development/go-workspace/bin/controller-gen "crd:trivialVersions=true,crdVersions=v1" rbac:roleName=manager-role webhook paths="./..." output:crd:artifacts:config=config/crd/bases
/Users/lbroudou/Development/go-workspace/bin/kustomize build config/crd | kubectl apply -f -
customresourcedefinition.apiextensions.k8s.io/fruitscatalogg1s.redhat.com created
```

And deploy the operator image and manifests into the default namespace that will be `fruits-catalog-go-operator-system`:

```
$ make deploy IMG=$IMG
```

### Testing the operator on the cluster

Just deploy a new Custom Resource coming from the `samples` within another namespace (so that we keep things separated and ordered ;-)):

```
$ kubectl create ns fruits-catalog-go
$ kubectl apply -f config/samples/redhat_v1beta1_fruitscatalogg1.yaml -n fruits-catalog-go
```

### operator-sdk 0.19.2

There's very few changes between `0.19.2` and `1.0.0` versions for the Go based operators. Check the completed resources within the [`fruits-catalog-go-operator-0.19.2`](./fruits-catalog-go-operator-0.19.2) folder of this Git repository for any details but process remains exactly the same. 