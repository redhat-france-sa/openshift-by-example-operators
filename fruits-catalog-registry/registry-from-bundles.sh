
# Build a new Operator REgsitry by indexing operators OCI bundles.
# By default, opm is using podman. We added the `-c docker` ffag to force docker usage.
opm index add --bundles quay.io/lbroudoux/fruits-catalog-ansible-operator-bundle:v0.0.1 \
    --tag quay.io/lbroudoux/fruits-catalog-operators-registry:1.0.0 -c docker

# Finally push the container image to public registry.
docker push quay.io/lbroudoux/fruits-catalog-operators-registry:1.0.0