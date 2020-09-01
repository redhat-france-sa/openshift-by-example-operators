In this project, we are going in the details of building a custom Operator Catalog in order to enrich the default marketplace being present within a Kubernetes / OpenShift cluster.

In the case of default setting of OpenShift, you may find some additional, helpful information in [Managing custom catalogs](https://docs.openshift.com/container-platform/4.5/operators/olm-managing-custom-catalogs.html#olm-managing-custom-catalogs-bundle-format).

## Pre-requisites

Depending on how you plan to build your Operator Registry and deploy, you may need `opm` or `operator-courier`. The Operator Package Manager is necessary for building registry using pre-existing Operator OCI bundles. The Operator Courier is necessary for uploading your Operator manisfests to Quay.io if you want to use it a global registry.

* You may find `opm` at https://github.com/operator-framework/operator-registry/releases. Download the binary corresponding to your platform and place it somewhere into your `$PATH`.

Tests were done using following version:
```
$ opm version                                                                                      
Version: version.Version{OpmVersion:"v1.13.8", GitCommit:"0deaced", BuildDate:"2020-08-25T12:44:21Z", GoOs:"darwin", GoArch:"amd64"}
```

* You may find `operator-courier` at https://github.com/operator-framework/operator-courier. Install the binary using the `pip3` package manager; it will be placed into your `$PATH`.

Tests were done using the following version:
```
$ operator-courier -v
2.1.10
```

### Build an Operator Registry


#### Option 1: Build a registry from Manifests

The [operator-registry](https://github.com/operator-framework/operator-registry) project defines a format for storing sets of operators and exposing them to make them available on a cluster. To create a catalog that includes your package, simply build a container image that uses the `upstream-registry-builder` tool to generate a registry and serve it. For example, create a file in the root of your project called `registry.Dockerfile`.

Then just use your favorite container tooling to build the container image and push it to a registry:

```
$ docker build -t quay.io/lbroudoux/fruits-catalog-operators-registry:latest -f registry.Dockerfile .
docker push quay.io/lbroudoux/fruits-catalog-operators-registry:latest
```

#### Option 2: Build a registry from Operator OCI Bundles

Another option for building such a registry container image is to reuse the OCI bundles containing Operators metadata that you may have created within your Operator project using the `operator-sdk`. This option uses the `opm` CLI tool that provides a specific `opm index` command for creating new registry from existings OCI bundles. 

The command to execute can be simply found into `registry-from-bundes.sh` script at the root of project.

#### Option 3: Reuse Quay.io internal registry

```
$ operator-courier push bundles lbroudoux fruits-catalog-operators 0.0.1 "basic bGJyb3Vkb3V4OkJsYWNrUmVkTGlnaHQ3Mg=="
```

### Deploy your catalog on Kubernetes

#### Option 1: Deploy a CatalogSource

#### Option 2: Deploy an OperatorSource
