
## Pre-requisites

Depending on how you plan to build your Operator Registry, you may need `opm`. The Operator Package Manager is necessary for building registry using pre-existing Operator OCI bundles.

* You may find `opm` at https://github.com/operator-framework/operator-registry/releases. Download the bunary corresponding to your platform and place it somewhere into your `$PATH`.

Tests were done using following version; 
```
$ opm version                                                                                      
Version: version.Version{OpmVersion:"v1.13.8", GitCommit:"0deaced", BuildDate:"2020-08-25T12:44:21Z", GoOs:"darwin", GoArch:"amd64"}
```

```
$ operator-courier -v
2.1.10
```

```
operator-courier push bundles lbroudoux fruits-catalog-operators 0.0.1 "basic bGJyb3Vkb3V4OkJsYWNrUmVkTGlnaHQ3Mg=="
```


```
curl -sH "Content-Type: application/json" -XPOST https://quay-2884.apps.shared-na4.na4.openshift.opentlc.com/cnr/api/v1/users/login -d '{"user": {"username": "quayadmin", "password": "kcMQYWPZzXFFJDwe"}}' | jq -r '.token'


curl -sH "Content-Type: application/json" -XPOST https://quay-2884.apps.shared-na4.na4.openshift.opentlc.com/cnr/api/v1/users/login -d '{"user": {"username": "lbroudoux", "password": "N9TOY7IOAOY6DKCT3Q75BYETQT7LRIM4"}}' | jq -r '.token'

N9TOY7IOAOY6DKCT3Q75BYETQT7LRIM4
```

### Build an Operator Registry

#### Option 1: Build a registry from Manifests

#### Option 2: Build a regsitry from Operator OCI Bundles

This option uses 

### Deploy your catalog on Kubernetes
