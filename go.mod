module github.com/jenkins-x-plugins/jx-registry

require (
	github.com/aws/aws-sdk-go v1.37.10
	github.com/gophercloud/gophercloud v0.15.0 // indirect
	github.com/jenkins-x/jx-api/v4 v4.0.24
	github.com/jenkins-x/jx-gitops v0.1.5
	github.com/jenkins-x/jx-helpers/v3 v3.0.77
	github.com/jenkins-x/jx-kube-client/v3 v3.0.2
	github.com/jenkins-x/jx-logging/v3 v3.0.3
	github.com/pkg/errors v0.9.1
	github.com/sethvargo/go-envconfig v0.3.2
	github.com/spf13/cobra v1.1.1
	github.com/stretchr/testify v1.6.1
	gopkg.in/yaml.v3 v3.0.0-20210107192922-496545a6307b // indirect
	k8s.io/client-go v11.0.1-0.20190805182717-6502b5e7b1b5+incompatible
)

replace (
	k8s.io/api => k8s.io/api v0.20.2
	k8s.io/apimachinery => k8s.io/apimachinery v0.20.2
	k8s.io/client-go => k8s.io/client-go v0.20.2
)

go 1.15
