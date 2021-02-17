package create

import (
	"fmt"
	"strings"

	"github.com/jenkins-x-plugins/jx-registry/pkg/amazon/ecrs"
	"github.com/jenkins-x-plugins/jx-registry/pkg/rootcmd"
	jxcore "github.com/jenkins-x/jx-api/v4/pkg/apis/core/v4beta1"
	"github.com/jenkins-x/jx-api/v4/pkg/client/clientset/versioned"
	"github.com/jenkins-x/jx-gitops/pkg/variablefinders"
	"github.com/jenkins-x/jx-helpers/v3/pkg/cmdrunner"
	"github.com/jenkins-x/jx-helpers/v3/pkg/cobras/helper"
	"github.com/jenkins-x/jx-helpers/v3/pkg/cobras/templates"
	"github.com/jenkins-x/jx-helpers/v3/pkg/gitclient"
	"github.com/jenkins-x/jx-helpers/v3/pkg/gitclient/cli"
	"github.com/jenkins-x/jx-helpers/v3/pkg/kube/jxclient"
	"github.com/jenkins-x/jx-helpers/v3/pkg/options"
	"github.com/jenkins-x/jx-helpers/v3/pkg/termcolor"
	"github.com/jenkins-x/jx-logging/v3/pkg/log"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

var (
	info = termcolor.ColorInfo

	cmdLong = templates.LongDesc(`
		Lazy create a container registry for ECR
`)

	cmdExample = templates.Examples(`
		# lets ensure we have an ECR registry setup
		%s create
	`)
)

// Options the options for this command
type Options struct {
	options.BaseOptions
	ecrs.Options

	CacheSuffix   string `env:"CACHE_SUFFIX"`
	ECRSuffix     string
	Namespace     string
	JXClient      versioned.Interface
	GitClient     gitclient.Interface
	CommandRunner cmdrunner.CommandRunner
	Requirements  *jxcore.RequirementsConfig
}

// NewCmdCreate creates a command object for the command
func NewCmdCreate() (*cobra.Command, *Options) {
	o := &Options{}

	cmd := &cobra.Command{
		Use:     "create",
		Short:   "Lazy create a container registry for ECR",
		Long:    cmdLong,
		Example: fmt.Sprintf(cmdExample, rootcmd.BinaryName),
		Run: func(cmd *cobra.Command, args []string) {
			err := o.Run()
			helper.CheckErr(err)
		},
	}

	if o.Context == nil {
		o.Context = cmd.Context()
	}
	o.Options.EnvProcess()

	o.Options.AddFlags(cmd)

	cmd.Flags().StringVarP(&o.Namespace, "namespace", "n", "", "The namespace. Defaults to the current namespace")
	cmd.Flags().StringVarP(&o.ECRSuffix, "ecr-registry-suffix", "", ".amazonaws.com", "The registry suffix to check if we are using ECR")
	cmd.Flags().StringVarP(&o.CacheSuffix, "cache-suffix", "", o.CacheSuffix, "If specified (or enabled via $CACHE_SUFFIX) we will make sure an ECR is created for the cache image too")

	o.BaseOptions.AddBaseFlags(cmd)
	return cmd, o
}

func (o *Options) Validate() error {
	if o.GitClient == nil {
		o.GitClient = cli.NewCLIClient("", o.CommandRunner)
	}
	if o.Requirements == nil {
		var err error
		o.JXClient, o.Namespace, err = jxclient.LazyCreateJXClientAndNamespace(o.JXClient, o.Namespace)
		if err != nil {
			return errors.Wrapf(err, "failed to create jxClient")
		}
		o.Requirements, err = variablefinders.FindRequirements(o.GitClient, o.JXClient, o.Namespace, "")
		if err != nil {
			return errors.Wrapf(err, "failed to load requirements from dev environment")
		}
	}
	if o.Requirements == nil {
		return errors.Errorf("no requirements found for dev environment")
	}

	if o.AWSRegion == "" {
		o.AWSRegion = o.Requirements.Cluster.Region
	}
	if o.Registry == "" {
		o.Registry = o.Requirements.Cluster.Registry
	}
	return nil

}
func (o *Options) Run() error {
	err := o.Validate()
	if err != nil {
		return errors.Wrapf(err, "failed to validate options")
	}
	if o.Requirements.Cluster.Provider != "eks" {
		log.Logger().Infof("no ECR code necessary as using provider %s", o.Requirements.Cluster.Provider)
		return nil
	}
	registry := o.Requirements.Cluster.Registry
	if registry != "" && registry != "ecr.io" && !strings.HasSuffix(registry, o.ECRSuffix) {
		log.Logger().Infof("ignoring registry %s ", registry)
		return nil
	}

	log.Logger().Infof("verifying that container registry %s with organisation %s and app name %s has an ECR associated with it", info(registry), info(o.RegistryOrganisation), info(o.AppName))

	images := []string{o.AppName}
	if o.CacheSuffix != "" {
		images = append(images, o.AppName+o.CacheSuffix)
	}
	for _, image := range images {
		err = o.Options.LazyCreateRegistry(image)
		if err != nil {
			return errors.Wrapf(err, "failed to lazy create the ECR registry for %s", image)
		}
	}
	return nil
}
