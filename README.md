# openshift-by-example-operators

## What is it?

This is a companion project to [`openshift-by-example`](https://github.com/redhat-france-sa/openshift-by-example) and the associated article that now tackle the topic of Operators development and packaging using tools and practices from the [Operator Framework](https://operatorframework.io).

![Operator SDK](https://master.sdk.operatorframework.io/build/images/logo.svg)

More precisely, this repository contains samples on how to develop Operators using the [Operator SDK](https://sdk.operatorframework.io) and the 3 different technologies embedded: [Helm](#helm-operator), [Ansible](#ansible-operator) and [Go language](#go-perator). We hope it will help you understand the pro and cons of different approaches and how they can map on the Operators Maturity Model phases.

![Operators Maturity Model](./assets/operators-maturity-model.png)

## The fruits-catalog application

The operators we have developed here

## Pre-requisites

* `operator-sdk` should be installed and present into your `$PATH`. We use latest `1.0.0` version and previous `0.19.2`version. Just check https://master.sdk.operatorframework.io/docs/installation/install-operator-sdk/ on how to install it,
* `docker` (or another tool compatible with multi-stage Dockerfiles) should be installed and present into your `$PATH`. Minimum version is `17.03+`,
* `go` should be installed and present into your `$PATH`. Minimum version is Go `1.13`,
* `kubectl` shoud be installed and present into your `$PATH`. `v1.16.0+` is the minimum version,

## Helm Operator

Operator SDK allows creating Operator using/re-using [Helm Charts](https://helm.sh). Whilst they're great starting points, we tend to think that Helm only allows you to support the first phases of the maturity model.

Discover how to develop such Operator on the [dedicated Helm page](HELM.md).

| Pros                     | Cons                     |
| ------------------------ | ------------------------ |
| | |

## Ansible Operator

Discover how to develop such Operator on the [dedicated Ansible page](ANSIBLE.md).

| Pros                     | Cons                     |
| ------------------------ | ------------------------ |
| | |

## Go Operator

Discover how to develop such Operator on the [dedicated Go page](GO.md).

| Pros                     | Cons                     |
| ------------------------ | ------------------------ |
| | |

## OLM manifests, Bundles and Scorecard

The Operator Lifecycle Manager (OLM) is a set of cluster resources that manage the lifecycle of an Operator. The Operator SDK supports both creating manifests for OLM deployment, and testing your Operator on an OLM-enabled Kubernetes cluster.

The OLM defines some mandatory metadata manifests that can be generated and then completed before having them packages into what will define an [Operator bundle](https://github.com/operator-framework/operator-registry/blob/v1.12.6/docs/design/operator-bundle.md).

So start, generating this famous manifests and everything we'll need to create a bundle! 

### Generating manifests metadata

As an example, we did this from within the [`fruits-catalog-ansible-operator-1.0.0`](./fruits-catalog-ansible-operator-1.0.0) folder

```
$ export IMG=quay.io/lbroudoux/fruits-catalog-ansible-operator:0.0.2
$ make bundle IMG=$IMG
operator-sdk generate kustomize manifests -q

Display name for the operator (required): 
> fruits-catalog-ansible-operator-1.0.0

Description for the operator (required): 
> Operator for the FruitsCatalog app using Operator Ansible SDK 1.0.0

Provider's name for the operator (required): 
> Laurent Broudoux

Any relevant URL for the provider name (optional): 
> https://github.com/redhat-france-sa/openshift-by-example-operators

Comma-separated list of keywords for your operator (required): 
> sample,operator-sdk,fruits-catalog

Comma-separated list of maintainers and their emails (e.g. 'name1:email1, name2:email2') (required): 
> lbroudoux:laurent.broudoux@redhat.com
cd config/manager && /Users/lbroudou/Development/go-workspace/bin/kustomize edit set image controller=quay.io/lbroudoux/fruits-catalog-ansible-operator:0.0.2
/Users/lbroudou/Development/go-workspace/bin/kustomize build config/manifests | operator-sdk generate bundle -q --overwrite --version 0.0.1  
INFO[0000] Building annotations.yaml                    
INFO[0000] Writing annotations.yaml in /Users/lbroudou/Development/github/openshift-by-example-operators/fruits-catalog-ansible-operator-1.0.0/bundle/metadata 
INFO[0000] Building Dockerfile                          
INFO[0000] Writing bundle.Dockerfile in /Users/lbroudou/Development/github/openshift-by-example-operators/fruits-catalog-ansible-operator-1.0.0 
operator-sdk bundle validate ./bundle
INFO[0000] Found annotations file                        bundle-dir=bundle container-tool=docker
INFO[0000] Could not find optional dependencies file     bundle-dir=bundle container-tool=docker
INFO[0000] All validation tests have completed successfully 
```

The main metadata manifests generated is a `Cluster Service Version`. It is indeed an aggregate of your Operator CRD, deployment informations, RBAC permissions + a number of metadata allowing to identify owner, provider and categorize operators. More information on CSV [here](https://olm.operatorframework.io/docs/tasks/packaging-an-operator/).

The previous command generates 2 things:
* Aggregation rules of elements using `kustomize` within the `config/manifests` folder of the project,
* Result of aggregation of elements within the `bundle/` folder of the project.

If you need to make adjustements to the generated CSV, you'll have to edit the `config/manifests/base` file and then re-run the following:

```
$ operator-sdk generate bundle -q --overwrite --version 0.0.1
```

> `0.0.1` being here the version of the generated CSV.

### Running the Scorecard tool

The [Scorecard](https://master.sdk.operatorframework.io/docs/advanced-topics/scorecard/scorecard/) command within Operator SDK executes tests on your operator manifests and future bundle based upon a configuration file and test images. Tests are implemented within test images that are configured and constructed to be executed by scorecard.

In order to run, Scorecard need you to be connected to a Kubernetes cluster. You just ask Scorecard to eevaluate the manifest that have been previously generated and are present into the `bundle/` folder of your project like this: 

```
$ operator-sdk scorecard bundle                                                                                                              
--------------------------------------------------------------------------------
Image:      quay.io/operator-framework/scorecard-test:master
Entrypoint: [scorecard-test olm-crds-have-resources]
Labels:
	"suite":"olm"
	"test":"olm-crds-have-resources-test"
Results:
	Name: olm-crds-have-resources
	State: pass

	Log:
		Loaded ClusterServiceVersion: fruits-catalog-ansible-operator.v0.0.1


--------------------------------------------------------------------------------
Image:      quay.io/operator-framework/scorecard-test:master
Entrypoint: [scorecard-test olm-crds-have-validation]
Labels:
	"test":"olm-crds-have-validation-test"
	"suite":"olm"
Results:
	Name: olm-crds-have-validation
	State: fail

	Suggestions:
		Add CRD validation for spec field `name` in FruitsCatalogA1/v1beta1
		Add CRD validation for spec field `webapp` in FruitsCatalogA1/v1beta1
	Log:
		Loaded 1 Custom Resources from alm-examples
		Loaded CustomresourceDefinitions: [&CustomResourceDefinition{ObjectMeta:{fruitscataloga1s.redhat.com      0 0001-01-01 00:00:00 +0000 UTC <nil> <nil> map[] map[] [] []  []},Spec:CustomResourceDefinitionSpec{Group:redhat.com,Names:CustomResourceDefinitionNames{Plural:fruitscataloga1s,Singular:fruitscataloga1,ShortNames:[],Kind:FruitsCatalogA1,ListKind:FruitsCatalogA1List,Categories:[],},Scope:Namespaced,Versions:[]CustomResourceDefinitionVersion{CustomResourceDefinitionVersion{Name:v1beta1,Served:true,Storage:true,Schema:&CustomResourceValidation{OpenAPIV3Schema:&JSONSchemaProps{ID:,Schema:,Ref:nil,Description:FruitsCatalogA1 is the Schema for the fruitscataloga1s API,Type:object,Format:,Title:,Default:nil,Maximum:nil,ExclusiveMaximum:false,Minimum:nil,ExclusiveMinimum:false,MaxLength:nil,MinLength:nil,Pattern:,MaxItems:nil,MinItems:nil,UniqueItems:false,MultipleOf:nil,Enum:[]JSON{},MaxProperties:nil,MinProperties:nil,Required:[],Items:nil,AllOf:[]JSONSchemaProps{},OneOf:[]JSONSchemaProps{},AnyOf:[]JSONSchemaProps{},Not:nil,Properties:map[string]JSONSchemaProps{apiVersion: {  <nil> APIVersion defines the versioned schema of this representation of an object. Servers should convert recognized schemas to the latest internal value, and may reject unrecognized values. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#resources string   nil <nil> false <nil> false <nil> <nil>  <nil> <nil> false <nil> [] <nil> <nil> [] nil [] [] [] nil map[] nil map[] map[] nil map[] nil nil false <nil> false false [] <nil> <nil>},kind: {  <nil> Kind is a string value representing the REST resource this object represents. Servers may infer this from the endpoint the client submits requests to. Cannot be updated. In CamelCase. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#types-kinds string   nil <nil> false <nil> false <nil> <nil>  <nil> <nil> false <nil> [] <nil> <nil> [] nil [] [] [] nil map[] nil map[] map[] nil map[] nil nil false <nil> false false [] <nil> <nil>},metadata: {  <nil>  object   nil <nil> false <nil> false <nil> <nil>  <nil> <nil> false <nil> [] <nil> <nil> [] nil [] [] [] nil map[] nil map[] map[] nil map[] nil nil false <nil> false false [] <nil> <nil>},spec: {  <nil> Spec defines the desired state of FruitsCatalogA1 object   nil <nil> false <nil> false <nil> <nil>  <nil> <nil> false <nil> [] <nil> <nil> [] nil [] [] [] nil map[] nil map[] map[] nil map[] nil nil false 0xc000236f10 false false [] <nil> <nil>},status: {  <nil> Status defines the observed state of FruitsCatalogA1 object   nil <nil> false <nil> false <nil> <nil>  <nil> <nil> false <nil> [] <nil> <nil> [] nil [] [] [] nil map[] nil map[] map[] nil map[] nil nil false 0xc000236f11 false false [] <nil> <nil>},},AdditionalProperties:nil,PatternProperties:map[string]JSONSchemaProps{},Dependencies:JSONSchemaDependencies{},AdditionalItems:nil,Definitions:JSONSchemaDefinitions{},ExternalDocs:nil,Example:nil,Nullable:false,XPreserveUnknownFields:nil,XEmbeddedResource:false,XIntOrString:false,XListMapKeys:[],XListType:nil,XMapType:nil,},},Subresources:&CustomResourceSubresources{Status:&CustomResourceSubresourceStatus{},Scale:nil,},AdditionalPrinterColumns:[]CustomResourceColumnDefinition{},},},Conversion:nil,PreserveUnknownFields:false,},Status:CustomResourceDefinitionStatus{Conditions:[]CustomResourceDefinitionCondition{},AcceptedNames:CustomResourceDefinitionNames{Plural:,Singular:,ShortNames:[],Kind:,ListKind:,Categories:[],},StoredVersions:[],},}]


--------------------------------------------------------------------------------
Image:      quay.io/operator-framework/scorecard-test:master
Entrypoint: [scorecard-test olm-status-descriptors]
Labels:
	"suite":"olm"
	"test":"olm-status-descriptors-test"
Results:
	Name: olm-status-descriptors
	State: fail

	Log:
		Loaded ClusterServiceVersion: fruits-catalog-ansible-operator.v0.0.1
		Loaded 1 Custom Resources from alm-examples


--------------------------------------------------------------------------------
Image:      quay.io/operator-framework/scorecard-test:master
Entrypoint: [scorecard-test basic-check-spec]
Labels:
	"suite":"basic"
	"test":"basic-check-spec-test"
Results:
	Name: basic-check-spec
	State: pass



--------------------------------------------------------------------------------
Image:      quay.io/operator-framework/scorecard-test:master
Entrypoint: [scorecard-test olm-spec-descriptors]
Labels:
	"suite":"olm"
	"test":"olm-spec-descriptors-test"
Results:
	Name: olm-spec-descriptors
	State: fail

	Suggestions:
		Add a spec descriptor for webapp
	Errors:
		webapp does not have a spec descriptor
	Log:
		Loaded ClusterServiceVersion: fruits-catalog-ansible-operator.v0.0.1
		Loaded 1 Custom Resources from alm-examples


--------------------------------------------------------------------------------
Image:      quay.io/operator-framework/scorecard-test:master
Entrypoint: [scorecard-test olm-bundle-validation]
Labels:
	"suite":"olm"
	"test":"olm-bundle-validation-test"
Results:
	Name: olm-bundle-validation
	State: pass

	Log:
		time="2020-08-31T07:52:00Z" level=debug msg="Found manifests directory" name=bundle-test
		time="2020-08-31T07:52:00Z" level=debug msg="Found metadata directory" name=bundle-test
		time="2020-08-31T07:52:00Z" level=debug msg="Getting mediaType info from manifests directory" name=bundle-test
		time="2020-08-31T07:52:00Z" level=info msg="Found annotations file" name=bundle-test
		time="2020-08-31T07:52:00Z" level=info msg="Could not find optional dependencies file" name=bundle-test
```

The Scorecard output gives you advices and suggestions for making your `Cluster Service Version` and aggregated manifests compliant with best practices.

### Building and validating the Bundle

Once you're happy with the Scorecard results, it's time to package and distribute your Operator bundle. The previous commands have also generated a `bundle.Dockerfile` at project root and you can simply use it for building a contaienr image for your bundle:

```
$ docker build -f bundle.Dockerfile -t quay.io/lbroudoux/fruits-catalog-ansible-operator-bundle:v0.0.1 .
$ docker push quay.io/lbroudoux/fruits-catalog-ansible-operator-bundle:v0.0.1
```

Last step before going is to validate the published bundle with this final command: 

```
$ operator-sdk bundle validate quay.io/lbroudoux/fruits-catalog-ansible-operator-bundle:v0.0.1                                              
INFO[0000] Unpacking image layers                       
INFO[0000] running /usr/local/bin/docker pull quay.io/lbroudoux/fruits-catalog-ansible-operator-bundle:v0.0.1 
INFO[0001] running docker create                        
INFO[0001] running docker cp                            
INFO[0002] running docker rm                            
INFO[0002] Found annotations file                        bundle-dir=/var/folders/wy/lb2qd3m51ld_k8ds2b6xkz180000gn/T/bundle-845300542 container-tool=docker
INFO[0002] Could not find optional dependencies file     bundle-dir=/var/folders/wy/lb2qd3m51ld_k8ds2b6xkz180000gn/T/bundle-845300542 container-tool=docker
INFO[0002] All validation tests have completed successfully
```

## Operators Registry and Catalog