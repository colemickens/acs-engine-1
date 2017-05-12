package cmd

import (
	"github.com/Azure/acs-engine/pkg/armhelpers"
	"github.com/Azure/go-autorest/autorest/azure"
	log "github.com/Sirupsen/logrus"
	"github.com/spf13/cobra"
	flag "github.com/spf13/pflag"
)

const (
	rootName             = "acs-engine"
	rootShortDescription = "ACS-Engine deploys and manages container orchestrators in Azure"
	rootLongDescription  = "ACS-Engine deploys and manages Kubernetes, Swarm Mode, and DC/OS clusters in Azure"
)

var (
	debug bool
)

func NewRootCmd() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   rootName,
		Short: rootShortDescription,
		Long:  rootLongDescription,
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			if debug {
				log.SetLevel(log.DebugLevel)
			}
		},
	}

	p := rootCmd.PersistentFlags()
	p.BoolVar(&debug, "debug", false, "enable verbose debug logs")

	rootCmd.AddCommand(NewVersionCmd())
	rootCmd.AddCommand(NewGenerateCmd())
	rootCmd.AddCommand(NewUpgradeCmd())

	return rootCmd
}

type authArgs struct {
	RawAzureEnvironment string
	SubscriptionID      string
	AuthMethod          string
	ClientID            string
	ClientSecret        string
	CertificatePath     string
	PrivateKeyPath      string
}

func addAuthFlags(authArgs *authArgs, f *flag.FlagSet) {
	f.StringVar(&authArgs.RawAzureEnvironment, "azure-env", "AzurePublicCloud", "the target Azure cloud")
	f.StringVar(&authArgs.SubscriptionID, "subscription-id", "", "azure subscription id")
	f.StringVar(&authArgs.AuthMethod, "auth-method", "device", "auth method (default:`device`, `client_secret`, `client_certificate`)")
	f.StringVar(&authArgs.ClientID, "client-id", "", "client id (used with --auth-method=[client_secret|client_certificate])")
	f.StringVar(&authArgs.ClientSecret, "client-secret", "", "client secret (used with --auth-mode=client_secret)")
	f.StringVar(&authArgs.CertificatePath, "certificate-path", "", "path to client certificate (used with --auth-method=client_certificate)")
	f.StringVar(&authArgs.PrivateKeyPath, "private-key-path", "", "path to private key (used with --auth-method=client_certificate)")
}

func (authArgs *authArgs) getClient() (*armhelpers.AzureClient, error) {
	// TODO: ensure subid is specified...

	if authArgs.AuthMethod == "client_secret" {
		if authArgs.ClientID == "" || authArgs.ClientSecret == "" {
			log.Fatal("--client-id and --client-secret must be specified when --auth-method=\"client_secret\".")
		}
	} else if authArgs.AuthMethod == "client_certificate" {
		if authArgs.ClientID == "" || authArgs.CertificatePath == "" || authArgs.PrivateKeyPath == "" {
			log.Fatal("--client-id and --certificate-path, and --private-key-path must be specified when --auth-method=\"client_certificate\".")
		}
	}

	env, err := azure.EnvironmentFromName(authArgs.RawAzureEnvironment)
	if err != nil {
		log.Fatal("failed to parse --azure-env as a valid target Azure cloud environment")
	}

	switch authArgs.AuthMethod {
	case "device":
		return armhelpers.NewAzureClientWithDeviceAuth(env, authArgs.SubscriptionID)
	case "client_secret":
		return armhelpers.NewAzureClientWithClientSecret(env, authArgs.SubscriptionID, authArgs.ClientID, authArgs.ClientSecret)
	case "client_certificate":
		return armhelpers.NewAzureClientWithClientCertificate(env, authArgs.SubscriptionID, authArgs.ClientID, authArgs.CertificatePath, authArgs.PrivateKeyPath)
	default:
		log.Fatalf("--auth-method: ERROR: method unsupported. method=%q.", authArgs.AuthMethod)
	}

	return nil, nil // unreachable
}
