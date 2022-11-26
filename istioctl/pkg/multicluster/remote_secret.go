// Copyright Istio Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package multicluster

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"strings"
	"time"
)

import (
	"github.com/cenkalti/backoff/v4"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"istio.io/pkg/log"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer/json"
	"k8s.io/apimachinery/pkg/runtime/serializer/versioning"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	_ "k8s.io/client-go/plugin/pkg/client/auth"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/clientcmd/api"
	"k8s.io/client-go/tools/clientcmd/api/latest"
)

import (
	"github.com/apache/dubbo-go-pixiu/operator/cmd/mesh"
	"github.com/apache/dubbo-go-pixiu/operator/pkg/helm"
	"github.com/apache/dubbo-go-pixiu/pkg/config/constants"
	"github.com/apache/dubbo-go-pixiu/pkg/config/labels"
	"github.com/apache/dubbo-go-pixiu/pkg/kube"
	"github.com/apache/dubbo-go-pixiu/pkg/kube/multicluster"
)

var (
	codec  runtime.Codec
	scheme *runtime.Scheme

	tokenWaitBackoff = time.Second
)

func init() {
	scheme = runtime.NewScheme()
	utilruntime.Must(v1.AddToScheme(scheme))
	opt := json.SerializerOptions{
		Yaml:   true,
		Pretty: false,
		Strict: false,
	}
	yamlSerializer := json.NewSerializerWithOptions(json.DefaultMetaFactory, scheme, scheme, opt)
	codec = versioning.NewDefaultingCodecForScheme(
		scheme,
		yamlSerializer,
		yamlSerializer,
		v1.SchemeGroupVersion,
		runtime.InternalGroupVersioner,
	)
}

const (
	remoteSecretPrefix = "istio-remote-secret-"
	configSecretName   = "istio-kubeconfig"
	configSecretKey    = "config"
)

func remoteSecretNameFromClusterName(clusterName string) string {
	return remoteSecretPrefix + clusterName
}

// NewCreateRemoteSecretCommand creates a new command for joining two contexts
// together in a multi-cluster mesh.
func NewCreateRemoteSecretCommand() *cobra.Command {
	opts := RemoteSecretOptions{
		AuthType:         RemoteSecretAuthTypeBearerToken,
		AuthPluginConfig: make(map[string]string),
		Type:             SecretTypeRemote,
	}
	c := &cobra.Command{
		Use:   "create-remote-secret",
		Short: "Create a secret with credentials to allow Istio to access remote Kubernetes apiservers",
		Example: `  # Create a secret to access cluster c0's apiserver and install it in cluster c1.
  istioctl --kubeconfig=c0.yaml x create-remote-secret --name c0 \
    | kubectl --kubeconfig=c1.yaml apply -f -

  # Delete a secret that was previously installed in c1
  istioctl --kubeconfig=c0.yaml x create-remote-secret --name c0 \
    | kubectl --kubeconfig=c1.yaml delete -f -

  # Create a secret access a remote cluster with an auth plugin
  istioctl --kubeconfig=c0.yaml x create-remote-secret --name c0 --auth-type=plugin --auth-plugin-name=gcp \
    | kubectl --kubeconfig=c1.yaml apply -f -`,
		Args: cobra.NoArgs,
		RunE: func(c *cobra.Command, args []string) error {
			if err := opts.prepare(c.Flags()); err != nil {
				return err
			}
			env, err := NewEnvironmentFromCobra(opts.Kubeconfig, opts.Context, c)
			if err != nil {
				return err
			}
			out, warn, err := CreateRemoteSecret(opts, env)
			if err != nil {
				_, _ = fmt.Fprintf(c.OutOrStderr(), "error: %v\n", err)
				return err
			}
			if warn != nil {
				_, _ = fmt.Fprintf(c.OutOrStderr(), "warn: %v\n", warn)
			}
			_, _ = fmt.Fprint(c.OutOrStdout(), out)
			return nil
		},
	}
	opts.addFlags(c.PersistentFlags())
	return c
}

func createRemoteServiceAccountSecret(kubeconfig *api.Config, clusterName, secName string) (*v1.Secret, error) { // nolint:interfacer
	var data bytes.Buffer
	if err := latest.Codec.Encode(kubeconfig, &data); err != nil {
		return nil, err
	}
	key := clusterName
	if secName == configSecretName {
		key = configSecretKey
	}
	out := &v1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name: secName,
			Annotations: map[string]string{
				clusterNameAnnotationKey: clusterName,
			},
			Labels: map[string]string{
				multicluster.MultiClusterSecretLabel: "true",
			},
		},
		Data: map[string][]byte{
			key: data.Bytes(),
		},
	}
	return out, nil
}

func createBaseKubeconfig(caData []byte, clusterName, server string) *api.Config {
	return &api.Config{
		Clusters: map[string]*api.Cluster{
			clusterName: {
				CertificateAuthorityData: caData,
				Server:                   server,
			},
		},
		AuthInfos: map[string]*api.AuthInfo{},
		Contexts: map[string]*api.Context{
			clusterName: {
				Cluster:  clusterName,
				AuthInfo: clusterName,
			},
		},
		CurrentContext: clusterName,
	}
}

func createBearerTokenKubeconfig(caData, token []byte, clusterName, server string) *api.Config {
	c := createBaseKubeconfig(caData, clusterName, server)
	c.AuthInfos[c.CurrentContext] = &api.AuthInfo{
		Token: string(token),
	}
	return c
}

func createPluginKubeconfig(caData []byte, clusterName, server string, authProviderConfig *api.AuthProviderConfig) *api.Config {
	c := createBaseKubeconfig(caData, clusterName, server)
	c.AuthInfos[c.CurrentContext] = &api.AuthInfo{
		AuthProvider: authProviderConfig,
	}
	return c
}

func createRemoteSecretFromPlugin(
	tokenSecret *v1.Secret,
	server, clusterName, secName string,
	authProviderConfig *api.AuthProviderConfig,
) (*v1.Secret, error) {
	caData, ok := tokenSecret.Data[v1.ServiceAccountRootCAKey]
	if !ok {
		return nil, errMissingRootCAKey
	}

	// Create a Kubeconfig to access the remote cluster using the auth provider plugin.
	kubeconfig := createPluginKubeconfig(caData, clusterName, server, authProviderConfig)
	if err := clientcmd.Validate(*kubeconfig); err != nil {
		return nil, fmt.Errorf("invalid kubeconfig: %v", err)
	}

	// Encode the Kubeconfig in a secret that can be loaded by Istio to dynamically discover and access the remote cluster.
	return createRemoteServiceAccountSecret(kubeconfig, clusterName, secName)
}

var (
	errMissingRootCAKey = fmt.Errorf("no %q data found", v1.ServiceAccountRootCAKey)
	errMissingTokenKey  = fmt.Errorf("no %q data found", v1.ServiceAccountTokenKey)
)

func createRemoteSecretFromTokenAndServer(client kube.ExtendedClient, tokenSecret *v1.Secret, clusterName, server, secName string) (*v1.Secret, error) {
	caData, token, err := waitForTokenData(client, tokenSecret)
	if err != nil {
		return nil, err
	}

	// Create a Kubeconfig to access the remote cluster using the remote service account credentials.
	kubeconfig := createBearerTokenKubeconfig(caData, token, clusterName, server)
	if err := clientcmd.Validate(*kubeconfig); err != nil {
		return nil, fmt.Errorf("invalid kubeconfig: %v", err)
	}

	// Encode the Kubeconfig in a secret that can be loaded by Istio to dynamically discover and access the remote cluster.
	return createRemoteServiceAccountSecret(kubeconfig, clusterName, secName)
}

func waitForTokenData(client kube.ExtendedClient, secret *v1.Secret) (ca, token []byte, err error) {
	ca, token, err = tokenDataFromSecret(secret)
	if err == nil {
		return
	}

	log.Infof("Waiting for data to be populated in %s", secret.Name)
	err = backoff.Retry(
		func() error {
			secret, err = client.Kube().CoreV1().Secrets(secret.Namespace).Get(context.TODO(), secret.Name, metav1.GetOptions{})
			if err != nil {
				return err
			}
			ca, token, err = tokenDataFromSecret(secret)
			return err
		},
		backoff.WithMaxRetries(backoff.NewConstantBackOff(tokenWaitBackoff), 5))
	return
}

func tokenDataFromSecret(tokenSecret *v1.Secret) (ca, token []byte, err error) {
	var ok bool
	ca, ok = tokenSecret.Data[v1.ServiceAccountRootCAKey]
	if !ok {
		err = errMissingRootCAKey
		return
	}
	token, ok = tokenSecret.Data[v1.ServiceAccountTokenKey]
	if !ok {
		err = errMissingTokenKey
		return
	}
	return
}

func getServiceAccountSecret(client kube.ExtendedClient, opt RemoteSecretOptions) (*v1.Secret, error) {
	// Create the service account if it doesn't exist.
	serviceAccount, err := getOrCreateServiceAccount(client, opt)
	if err != nil {
		return nil, err
	}

	if !kube.IsAtLeastVersion(client, 24) {
		return legacyGetServiceAccountSecret(serviceAccount, client, opt)
	}
	return getOrCreateServiceAccountSecret(serviceAccount, client, opt)
}

// In Kubernetes 1.24+ we can't assume the secrets will be referenced in the ServiceAccount or be created automatically.
// See https://github.com/istio/istio/issues/38246
func getOrCreateServiceAccountSecret(
	serviceAccount *v1.ServiceAccount,
	client kube.ExtendedClient,
	opt RemoteSecretOptions,
) (*v1.Secret, error) {
	ctx := context.TODO()

	// manually specified secret, make sure it references the ServiceAccount
	if opt.SecretName != "" {
		secret, err := client.Kube().CoreV1().Secrets(opt.Namespace).Get(ctx, opt.SecretName, metav1.GetOptions{})
		if err != nil {
			return nil, fmt.Errorf("could not get specified secret %s/%s: %v",
				opt.Namespace, opt.SecretName, err)
		}
		if err := secretReferencesServiceAccount(serviceAccount, secret); err != nil {
			return nil, err
		}
		return secret, nil
	}

	// first try to find an existing secret that references the SA
	// TODO will the SA have any reference to secrets anymore, can we avoid this list?
	allSecrets, err := client.Kube().CoreV1().Secrets(opt.Namespace).List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed listing secrets in %s: %v", opt.Namespace, err)
	}
	for _, item := range allSecrets.Items {
		secret := item
		if secretReferencesServiceAccount(serviceAccount, &secret) == nil {
			return &secret, nil
		}
	}

	// finally, create the sa token secret manually
	// https://kubernetes.io/docs/tasks/configure-pod-container/configure-service-account/#manually-create-a-service-account-api-token
	// TODO ephemeral time-based tokens are preferred; we should re-think this
	log.Infof("Creating token secret for service account %q", serviceAccount.Name)
	secretName := tokenSecretName(serviceAccount.Name)
	return client.Kube().CoreV1().Secrets(opt.Namespace).Create(ctx, &v1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:        secretName,
			Annotations: map[string]string{v1.ServiceAccountNameKey: serviceAccount.Name},
		},
		Type: v1.SecretTypeServiceAccountToken,
	}, metav1.CreateOptions{})
}

func tokenSecretName(saName string) string {
	return saName + "-istio-remote-secret-token"
}

func secretReferencesServiceAccount(serviceAccount *v1.ServiceAccount, secret *v1.Secret) error {
	if secret.Type != v1.SecretTypeServiceAccountToken ||
		secret.Annotations[v1.ServiceAccountNameKey] != serviceAccount.Name {
		return fmt.Errorf("secret %s/%s does not reference ServiceAccount %s",
			secret.Namespace, secret.Name, serviceAccount.Name)
	}
	return nil
}

func legacyGetServiceAccountSecret(
	serviceAccount *v1.ServiceAccount,
	client kube.ExtendedClient,
	opt RemoteSecretOptions,
) (*v1.Secret, error) {
	if len(serviceAccount.Secrets) == 0 {
		return nil, fmt.Errorf("no secret found in the service account: %s", serviceAccount)
	}

	secretName := ""
	secretNamespace := ""
	if opt.SecretName != "" {
		found := false
		for _, secret := range serviceAccount.Secrets {
			if secret.Name == opt.SecretName {
				found = true
				secretName = secret.Name
				secretNamespace = secret.Namespace
				break
			}
		}
		if !found {
			return nil, fmt.Errorf("provided secret does not exist: %s", opt.SecretName)
		}
	} else {
		if len(serviceAccount.Secrets) == 1 {
			secretName = serviceAccount.Secrets[0].Name
			secretNamespace = serviceAccount.Secrets[0].Namespace
		} else {
			return nil, fmt.Errorf("wrong number of secrets (%v) in serviceaccount %s/%s, please use --secret-name to specify one",
				len(serviceAccount.Secrets), opt.Namespace, opt.ServiceAccountName)
		}
	}

	if secretNamespace == "" {
		secretNamespace = opt.Namespace
	}
	return client.Kube().CoreV1().Secrets(secretNamespace).Get(context.TODO(), secretName, metav1.GetOptions{})
}

func getOrCreateServiceAccount(client kube.ExtendedClient, opt RemoteSecretOptions) (*v1.ServiceAccount, error) {
	if sa, err := client.Kube().CoreV1().ServiceAccounts(opt.Namespace).Get(
		context.TODO(), opt.ServiceAccountName, metav1.GetOptions{}); err == nil {
		return sa, nil
	} else if !opt.CreateServiceAccount {
		// User chose not to automatically create the service account.
		return nil, fmt.Errorf("failed retrieving service account %s.%s required for creating "+
			"the remote secret (hint: try installing a minimal Istio profile on the cluster first, "+
			"or run with '--create-service-account=true'): %v",
			opt.ServiceAccountName,
			opt.Namespace,
			err)
	}

	if err := createServiceAccount(client, opt); err != nil {
		return nil, err
	}

	// Return the newly created service account.
	sa, err := client.Kube().CoreV1().ServiceAccounts(opt.Namespace).Get(
		context.TODO(), opt.ServiceAccountName, metav1.GetOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed retrieving service account %s.%s after creating it: %v",
			opt.ServiceAccountName, opt.Namespace, err)
	}
	return sa, nil
}

func createServiceAccount(client kube.ExtendedClient, opt RemoteSecretOptions) error {
	yaml, err := generateServiceAccountYAML(opt)
	if err != nil {
		return err
	}

	// Before we can apply the yaml, we have to ensure the system namespace exists.
	if err := createNamespaceIfNotExist(client, opt.Namespace); err != nil {
		return err
	}

	// Apply the YAML to the cluster.
	return applyYAML(client, yaml, opt.Namespace)
}

func generateServiceAccountYAML(opt RemoteSecretOptions) (string, error) {
	// Create a renderer for the base installation.
	baseRenderer := helm.NewHelmRenderer(opt.ManifestsPath, "base", "Base", opt.Namespace, nil)
	discoveryRenderer := helm.NewHelmRenderer(opt.ManifestsPath, "istio-control/istio-discovery", "Pilot", opt.Namespace, nil)

	baseTemplates := []string{"reader-serviceaccount.yaml"}
	discoveryTemplates := []string{"clusterrole.yaml", "clusterrolebinding.yaml"}

	if err := baseRenderer.Run(); err != nil {
		return "", fmt.Errorf("failed running base Helm renderer: %w", err)
	}
	if err := discoveryRenderer.Run(); err != nil {
		return "", fmt.Errorf("failed running base discovery Helm renderer: %w", err)
	}

	values := fmt.Sprintf(`
global:
  istioNamespace: %s
`, opt.Namespace)

	// Render the templates required for the service account and role bindings.
	baseContent, err := baseRenderer.RenderManifestFiltered(values, func(template string) bool {
		for _, t := range baseTemplates {
			if strings.Contains(template, t) {
				return true
			}
		}
		return false
	})
	if err != nil {
		return "", fmt.Errorf("failed rendering base manifest: %w", err)
	}
	discoveryContent, err := discoveryRenderer.RenderManifestFiltered(values, func(template string) bool {
		for _, t := range discoveryTemplates {
			if strings.Contains(template, t) {
				return true
			}
		}
		return false
	})
	if err != nil {
		return "", fmt.Errorf("failed rendering discovery manifest: %w", err)
	}

	aggregateContent := fmt.Sprintf(`
%s
---
%s
`, baseContent, discoveryContent)
	return aggregateContent, nil
}

func applyYAML(client kube.ExtendedClient, yamlContent, ns string) error {
	yamlFile, err := writeToTempFile(yamlContent)
	if err != nil {
		return fmt.Errorf("failed creating manifest file: %v", err)
	}

	// Apply the YAML to the cluster.
	if err := client.ApplyYAMLFiles(ns, yamlFile); err != nil {
		return fmt.Errorf("failed applying manifest %s: %v", yamlFile, err)
	}
	return nil
}

func createNamespaceIfNotExist(client kube.Client, ns string) error {
	if _, err := client.Kube().CoreV1().Namespaces().Get(context.TODO(), ns, metav1.GetOptions{}); err != nil {
		if _, err := client.Kube().CoreV1().Namespaces().Create(context.TODO(), &v1.Namespace{
			ObjectMeta: metav1.ObjectMeta{
				Name: ns,
			},
		}, metav1.CreateOptions{}); err != nil {
			return fmt.Errorf("failed creating namespace %s: %v", ns, err)
		}
	}
	return nil
}

func writeToTempFile(content string) (string, error) {
	outFile, err := os.CreateTemp("", "remote-secret-manifest-*")
	if err != nil {
		return "", fmt.Errorf("failed creating temp file for manifest: %v", err)
	}
	defer func() { _ = outFile.Close() }()

	if _, err := outFile.Write([]byte(content)); err != nil {
		return "", fmt.Errorf("failed writing manifest file: %v", err)
	}
	return outFile.Name(), nil
}

func getServerFromKubeconfig(context string, config *api.Config) (string, Warning, error) {
	if context == "" {
		context = config.CurrentContext
	}

	configContext, ok := config.Contexts[context]
	if !ok {
		return "", nil, fmt.Errorf("could not find cluster for context %q", context)
	}
	cluster, ok := config.Clusters[configContext.Cluster]
	if !ok {
		return "", nil, fmt.Errorf("could not find server for context %q", context)
	}
	if strings.Contains(cluster.Server, "127.0.0.1") || strings.Contains(cluster.Server, "localhost") {
		return cluster.Server, fmt.Errorf(
			"server in Kubeconfig is %s. This is likely not reachable from inside the cluster, "+
				"if you're using Kubernetes in Docker, pass --server with the container IP for the API Server",
			cluster.Server), nil
	}
	return cluster.Server, nil, nil
}

const (
	outputHeader  = "# This file is autogenerated, do not edit.\n"
	outputTrailer = "---\n"
)

func writeEncodedObject(out io.Writer, in runtime.Object) error {
	if _, err := fmt.Fprint(out, outputHeader); err != nil {
		return err
	}
	if err := codec.Encode(in, out); err != nil {
		return err
	}
	if _, err := fmt.Fprint(out, outputTrailer); err != nil {
		return err
	}
	return nil
}

type writer interface {
	io.Writer
	String() string
}

func makeOutputWriter() writer {
	return &bytes.Buffer{}
}

var makeOutputWriterTestHook = makeOutputWriter

// RemoteSecretAuthType is a strongly typed authentication type suitable for use with pflags.Var().
type (
	RemoteSecretAuthType string
	SecretType           string
)

var _ pflag.Value = (*RemoteSecretAuthType)(nil)

func (at *RemoteSecretAuthType) String() string { return string(*at) }
func (at *RemoteSecretAuthType) Type() string   { return "RemoteSecretAuthType" }
func (at *RemoteSecretAuthType) Set(in string) error {
	*at = RemoteSecretAuthType(in)
	return nil
}

func (at *SecretType) String() string { return string(*at) }
func (at *SecretType) Type() string   { return "SecretType" }
func (at *SecretType) Set(in string) error {
	*at = SecretType(in)
	return nil
}

const (
	// Use a bearer token for authentication to the remote kubernetes cluster.
	RemoteSecretAuthTypeBearerToken RemoteSecretAuthType = "bearer-token"

	// Use a custom authentication plugin for the remote kubernetes cluster.
	RemoteSecretAuthTypePlugin RemoteSecretAuthType = "plugin"

	// Secret generated from remote cluster
	SecretTypeRemote SecretType = "remote"

	// Secret generated from config cluster
	SecretTypeConfig SecretType = "config"
)

// RemoteSecretOptions contains the options for creating a remote secret.
type RemoteSecretOptions struct {
	KubeOptions

	// Name of the local cluster whose credentials are stored in the secret. Must be
	// DNS1123 label as it will be used for the k8s secret name.
	ClusterName string

	// Create a secret with this service account's credentials.
	ServiceAccountName string

	// CreateServiceAccount if true, the service account specified by ServiceAccountName
	// will be created if it doesn't exist.
	CreateServiceAccount bool

	// Authentication method for the remote Kubernetes cluster.
	AuthType RemoteSecretAuthType
	// Authenticator plugin configuration
	AuthPluginName   string
	AuthPluginConfig map[string]string

	// Type of the generated secret
	Type SecretType

	// ManifestsPath is a path to a manifestsPath and profiles directory in the local filesystem,
	// or URL with a release tgz. This is only used when no reader service account exists and has
	// to be created.
	ManifestsPath string

	// ServerOverride overrides the server IP/hostname field from the Kubeconfig
	ServerOverride string

	// SecretName selects a specific secret from the remote service account, if there are multiple
	SecretName string
}

func (o *RemoteSecretOptions) addFlags(flagset *pflag.FlagSet) {
	flagset.StringVar(&o.ServiceAccountName, "service-account", "",
		"Create a secret with this service account's credentials. Default value is \""+
			constants.DefaultServiceAccountName+"\" if --type is \"remote\", \""+
			constants.DefaultConfigServiceAccountName+"\" if --type is \"config\".")
	flagset.BoolVar(&o.CreateServiceAccount, "create-service-account", true,
		"If true, the service account needed for creating the remote secret will be created "+
			"if it doesn't exist.")
	flagset.StringVar(&o.ClusterName, "name", "",
		"Name of the local cluster whose credentials are stored "+
			"in the secret. If a name is not specified the kube-system namespace's UUID of "+
			"the local cluster will be used.")
	flagset.StringVar(&o.ServerOverride, "server", "",
		"The address and port of the Kubernetes API server.")
	flagset.StringVar(&o.SecretName, "secret-name", "",
		"The name of the specific secret to use from the service-account. Needed when there are multiple secrets in the service account.")
	var supportedAuthType []string
	for _, at := range []RemoteSecretAuthType{RemoteSecretAuthTypeBearerToken, RemoteSecretAuthTypePlugin} {
		supportedAuthType = append(supportedAuthType, string(at))
	}
	var supportedSecretType []string
	for _, at := range []SecretType{SecretTypeRemote, SecretTypeConfig} {
		supportedSecretType = append(supportedSecretType, string(at))
	}

	flagset.Var(&o.AuthType, "auth-type",
		fmt.Sprintf("Type of authentication to use. supported values = %v", supportedAuthType))
	flagset.StringVar(&o.AuthPluginName, "auth-plugin-name", o.AuthPluginName,
		fmt.Sprintf("Authenticator plug-in name. --auth-type=%v must be set with this option",
			RemoteSecretAuthTypePlugin))
	flagset.StringToString("auth-plugin-config", o.AuthPluginConfig,
		fmt.Sprintf("Authenticator plug-in configuration. --auth-type=%v must be set with this option",
			RemoteSecretAuthTypePlugin))
	flagset.Var(&o.Type, "type",
		fmt.Sprintf("Type of the generated secret. supported values = %v", supportedSecretType))
	flagset.StringVarP(&o.ManifestsPath, "manifests", "d", "", mesh.ManifestsFlagHelpStr)
}

func (o *RemoteSecretOptions) prepare(flags *pflag.FlagSet) error {
	o.KubeOptions.prepare(flags)

	if o.ClusterName != "" {
		if !labels.IsDNS1123Label(o.ClusterName) {
			return fmt.Errorf("%v is not a valid DNS 1123 label", o.ClusterName)
		}
	}
	return nil
}

type Warning error

func createRemoteSecret(opt RemoteSecretOptions, client kube.ExtendedClient, env Environment) (*v1.Secret, Warning, error) {
	// generate the clusterName if not specified
	if opt.ClusterName == "" {
		uid, err := clusterUID(client)
		if err != nil {
			return nil, nil, err
		}
		opt.ClusterName = string(uid)
	}

	var secretName string
	switch opt.Type {
	case SecretTypeRemote:
		secretName = remoteSecretNameFromClusterName(opt.ClusterName)
		if opt.ServiceAccountName == "" {
			opt.ServiceAccountName = constants.DefaultServiceAccountName
		}
	case SecretTypeConfig:
		secretName = configSecretName
		if opt.ServiceAccountName == "" {
			opt.ServiceAccountName = constants.DefaultConfigServiceAccountName
		}
	default:
		return nil, nil, fmt.Errorf("unsupported type: %v", opt.Type)
	}
	tokenSecret, err := getServiceAccountSecret(client, opt)
	if err != nil {
		return nil, nil, fmt.Errorf("could not get access token to read resources from local kube-apiserver: %v", err)
	}

	var server string
	var warn Warning
	if opt.ServerOverride != "" {
		server = opt.ServerOverride
	} else {
		server, warn, err = getServerFromKubeconfig(opt.Context, env.GetConfig())
		if err != nil {
			return nil, warn, err
		}
	}

	var remoteSecret *v1.Secret
	switch opt.AuthType {
	case RemoteSecretAuthTypeBearerToken:
		remoteSecret, err = createRemoteSecretFromTokenAndServer(client, tokenSecret, opt.ClusterName, server, secretName)
	case RemoteSecretAuthTypePlugin:
		authProviderConfig := &api.AuthProviderConfig{
			Name:   opt.AuthPluginName,
			Config: opt.AuthPluginConfig,
		}
		remoteSecret, err = createRemoteSecretFromPlugin(tokenSecret, server, opt.ClusterName, secretName,
			authProviderConfig)
	default:
		err = fmt.Errorf("unsupported authentication type: %v", opt.AuthType)
	}
	if err != nil {
		return nil, warn, err
	}

	remoteSecret.Namespace = opt.Namespace
	return remoteSecret, warn, nil
}

// CreateRemoteSecret creates a remote secret with credentials of the specified service account.
// This is useful for providing a cluster access to a remote apiserver.
func CreateRemoteSecret(opt RemoteSecretOptions, env Environment) (string, Warning, error) {
	client, err := env.CreateClient(opt.Context)
	if err != nil {
		return "", nil, err
	}

	remoteSecret, warn, err := createRemoteSecret(opt, client, env)
	if err != nil {
		return "", warn, err
	}

	// convert any binary data to the string equivalent for easier review. The
	// kube-apiserver will convert this to binary before it persists it to storage.
	remoteSecret.StringData = make(map[string]string, len(remoteSecret.Data))
	for k, v := range remoteSecret.Data {
		remoteSecret.StringData[k] = string(v)
	}
	remoteSecret.Data = nil

	w := makeOutputWriterTestHook()
	if err := writeEncodedObject(w, remoteSecret); err != nil {
		return "", warn, err
	}
	return w.String(), warn, nil
}
