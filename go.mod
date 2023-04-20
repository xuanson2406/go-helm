module github.com/xuanson2406/helm-clientgo-example

go 1.16

require (
	github.com/gofrs/flock v0.8.1
	github.com/minio/minio-go/v7 v7.0.50
	github.com/pkg/errors v0.9.1
	gopkg.in/yaml.v2 v2.4.0
	helm.sh/helm/v3 v3.11.2
	k8s.io/helm v2.17.0+incompatible
)

replace github.com/docker/docker => github.com/moby/moby v0.7.3-0.20190826074503-38ab9da00309
