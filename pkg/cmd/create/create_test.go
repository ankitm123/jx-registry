package create_test

import (
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/jenkins-x-plugins/jx-registry/pkg/amazon/ecrs/fakeecr"
	"github.com/jenkins-x-plugins/jx-registry/pkg/cmd/create"
	jxcore "github.com/jenkins-x/jx-api/v4/pkg/apis/core/v4beta1"
	"github.com/stretchr/testify/require"
)

var (
	// generateTestOutput enable to regenerate the expected output
	generateTestOutput = false
)

func TestCreateForNonEKS(t *testing.T) {
	_, o := create.NewCmdCreate()

	o.Requirements = &jxcore.RequirementsConfig{
		Cluster: jxcore.ClusterConfig{
			Provider: "gke",
		},
	}

	err := o.Run()
	require.NoError(t, err, "failed to run")
}

func TestCreateForEKS(t *testing.T) {
	_, o := create.NewCmdCreate()

	o.Requirements = &jxcore.RequirementsConfig{
		Cluster: jxcore.ClusterConfig{
			Provider: "eks",
		},
	}

	o.AWSRegion = "dummy"
	o.Config = &aws.Config{}
	o.AppName = "myapp"
	fakeECR := fakeecr.NewFakeECR()
	o.ECRClient = fakeECR
	o.CacheSuffix = "/cache"
	o.RegistryID = "123456789012"

	err := o.Run()
	require.NoError(t, err, "failed to run")

	// lets check we have a repository
	require.Len(t, fakeECR.Repositories, 2, "should have created 2 repositories")

	for _, v := range fakeECR.Repositories {
		require.NotNil(t, v.RepositoryUri, "should have a repository URI")
		t.Logf("found ECR repository %s\n", *v.RepositoryUri)
	}
}
