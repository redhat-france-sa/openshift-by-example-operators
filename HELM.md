## Helm Operator

### operator-sdk 1.0.0

Start creating a new directory and initializing a new helm operator project in it:
```
$ mkdir fruits-catalog-helm-operator && cd fruits-catalog-helm-operator
$ operator-sdk init --plugins=helm --domain=com --group=redhat --version v1beta1 --kind FruitsCatalogH1
```

#### Completing the operator code

* Edit the `helm-charts/fruitscatalogh1/values.yaml` to reflect your custom resource default values,
* Edit all the resources within `helm-charts/fruitscatalogh1/templates` resources accordingly,
* Adapt the required permsissions within `config/rbac/role.yaml`

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

```
$ kubectl create ns fruits-catalog-helm-1
$ kubectl apply -f config/samples/redhat_v1beta1_fruitscatalogh1.yaml -n fruits-catalog-helm-1
```

> Does not work! Why ? Normal ;-) 

```
$ kubectl apply -f config/samples/redhat_v1beta1_fruitscatalogh1-username+password.yaml -n fruits-catalog-helm-1
```

### operator-sdk 0.19.2