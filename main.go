// Example to demonstrate helm chart installation using helm client-go
// Most of the code is copied from https://github.com/helm/helm repo

package main

import (
	"context"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/gofrs/flock"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/chart"
	"helm.sh/helm/v3/pkg/chart/loader"
	"helm.sh/helm/v3/pkg/cli"
	"helm.sh/helm/v3/pkg/cli/values"
	"helm.sh/helm/v3/pkg/downloader"
	"helm.sh/helm/v3/pkg/getter"
	"helm.sh/helm/v3/pkg/repo"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
	"k8s.io/helm/pkg/strvals"
)

var settings *cli.EnvSettings
var err error

var (
	// url1 = "https://registry.fke.fptcloud.com/chartrepo/xplat-fke"
	url1 = "https://aquasecurity.github.io/helm-charts/"
	// repoName           = "xplat-fke"
	repoName = "aqua"
	// chartName          = "gpu-operator"
	chartName = "trivy-operator"
	// releaseName        = "operator"
	releaseName = "trivy-operator"
	// namespace          = "gpu-operator"
	namespace          = "trivy-system"
	prometheus_service = "prometheus-stack-kube-prom-prometheus"
	args               = map[string]string{
		// comma seperated values to set
		// "set": "database.volume.storageClassName=vsan-default-storage-policy,imagePullPolicy=Always",
		// "set": "prometheus.url=http://" + prometheus_service + ".prometheus.svc.fke-demo-lab2",
		"set": "targetNamespaces=kube-system,kube-public,fptcloud-runtime-security,default",
		// "set": "mig.strategy=mixed",
	}
	clusterName = "llbo913r"
	credential  = map[string]string{
		"endPoint":        "s3-stg09.fptcloud.net",
		"accessKeyID":     "00da3ac6e85660b4e9c7",
		"secretAccessKey": "hKPSk+kfUs3STgF0kF3H/Pt4CQ1U4fOMkWq8jNft",
		"bucketName":      "kubeconfig",
	}
	strategy = "mixed"
)

func main() {
	// config, err := buildConfig("target-kubeconfig.yaml")
	// if err != nil {
	// 	log.Fatalf("Error building kubeconfig: %v", err)
	// }
	// config, err := rest.InClusterConfig()
	// if err != nil {
	// 	// If running outside the cluster, use kubeconfig file
	// 	kubeconfig := filepath.Join(homedir.HomeDir(), ".kube", "target-kubeconfig.yaml ")
	// 	config, err = clientcmd.BuildConfigFromFlags("", kubeconfig)
	// 	if err != nil {
	// 		log.Fatalf("Failed to build Kubernetes client configuration: %v", err)
	// 	}
	// }
	// clientset, err := kubernetes.NewForConfig(config)
	// if err != nil {
	// 	log.Fatalf("Error building kubernetes clientset: %v", err)
	// }
	// yamlFile, err := os.ReadFile("/home/sondx12/xuanson2406/job.json")
	// if err != nil {
	// 	panic(err.Error())
	// }
	// // Decode the YAML file into a Pod object
	// job := batchv1.Job{}
	// if err := json.Unmarshal(yamlFile, &job); err != nil {
	// 	panic(err.Error())
	// }
	// // Create the Job
	// fmt.Println("Creating Job...")
	// fmt.Printf("name: %s\nrestartPolicy: %s", job.Name, job.Spec.Template.Spec.RestartPolicy)
	// result, err := clientset.BatchV1().Jobs("kube-system").Create(context.TODO(), &job, metav1.CreateOptions{})
	// if err != nil {
	// 	panic(err)
	// }
	// fmt.Printf("Created Job %q.\n", result.GetObjectMeta().GetName())

	// Read the YAML file
	// filePath := "/home/sondx12/xuanson2406/crd2.yaml"
	// crdYaml, err := ioutil.ReadFile(filePath)
	// if err != nil {
	// 	fmt.Printf("Error reading YAML file: %v\n", err)
	// 	os.Exit(1)
	// }
	// // Convert []byte to io.Reader
	// crdReader := bytes.NewReader(crdYaml)
	// // Decode the YAML into an unstructured object
	// // dec := yaml.New(unstructured.UnstructuredJSONScheme)
	// obj := &unstructured.Unstructured{}
	// decoder := k8syaml.NewYAMLOrJSONDecoder(crdReader, 10000)
	// if err := decoder.Decode(obj); err != nil {
	// 	fmt.Printf("Error decoding YAML: %v\n", err)
	// 	os.Exit(1)
	// }

	// // Apply the CRD to the cluster
	// dynamicClient, err := dynamic.NewForConfig(config)
	// if err != nil {
	// 	fmt.Printf("Error creating dynamic client: %v\n", err)
	// 	os.Exit(1)
	// }

	// gvr := schema.GroupVersionResource{
	// 	Group:    "apiextensions.k8s.io",
	// 	Version:  "v1",
	// 	Resource: "customresourcedefinitions",
	// }

	// // Use the dynamic client to create or update the CRD
	// _, err = dynamicClient.Resource(gvr).Create(context.Background(), obj, metav1.CreateOptions{})
	// if err != nil {
	// 	fmt.Printf("Error applying CRD: %v\n", err)
	// 	os.Exit(1)
	// }

	// fmt.Println("CRD applied successfully!")

	// nodeName := "fke-check125-lhogmw2l-worker-1y886uoh-5cb89-xxhr8"
	// node, _ := clientset.CoreV1().Nodes().Get(context.Background(), nodeName, metav1.GetOptions{})
	// if node.Labels["worker.fptcloud/system-components"] != "true" {
	// 	fmt.Printf("true")
	// }
	// os.Setenv("HELM_NAMESPACE", "gpu-operator")
	// settings := cli.New()
	// settings.KubeConfig = "/home/sondx12/go/src/github.com/xuanson2406/go-helm/kubecfg/test.yaml"
	// actionConfig := new(action.Configuration)

	// err := actionConfig.Init(settings.RESTClientGetter(), settings.Namespace(), os.Getenv("HELM_DRIVER"), log.Printf)
	// if err != nil {
	// 	panic(err)
	// }
	// // Load the Helm chart from the specified path
	// chartPath := "/home/sondx12/go/src/github.com/xuanson2406/go-helm/gpu-operator"
	// chart, err := loader.Load(chartPath)
	// if err != nil {
	// 	panic(err)
	// }

	// // Use actionConfig to perform Helm operations
	// // For example, you can install a Helm chart:
	// installChart := action.NewInstall(actionConfig)
	// installChart.Namespace = "gpu-operator"
	// installChart.ReleaseName = "sondx12"
	// installChart.CreateNamespace = true

	// _, err = installChart.Run(chart, nil)
	// if err != nil {
	// 	panic(err)
	// }

	// fmt.Println("Helm chart installed successfully!")

	// os.Setenv("HELM_NAMESPACE", namespace)
	// settings.KubeConfig, err = downloadKubecfg(clusterName, credential)
	// if err != nil {
	// 	fmt.Printf("Can not download file kubeconfig of cluster %s: %v", clusterName, err)
	// }
	// settings.KubeConfig = "/home/sondx12/.kube/shoot-config"
	// kubeconfigFile := os.Getenv("HOME") + "/.kube/target-kubeconfig.yaml"
	// settings = CreateSetting("", kubeconfigFile)
	// // Add helm repo
	// RepoAdd(repoName, url1)

	// // Update charts from the helm repo
	// RepoUpdate()

	// // Install charts
	// settings = CreateSetting(namespace, kubeconfigFile)
	// InstallChart(releaseName, repoName, chartName, namespace)

	rule :=
		`- rule: kube-system-ignore
	condition:
	  kubernetes.pod.namespace == "kube-system"
	priority: 100
	output:
	  skip: true`

	customRule := fmt.Sprintf(`
custom-rules.yaml: |-
  %s\n`, rule)
	fmt.Printf(customRule)
	// Install GPU Operator
	// os.Setenv("HELM_NAMESPACE", "gpu-operator")
	// err := InstallChart(releaseName, repoName, chartName, args)
	// if err != nil {
	// 	fmt.Printf("Unable to install chart gpu-operator: %v", err)
	// } else {
	// 	fmt.Println("Successfully install chart gpu-operator")
	// 	time.Sleep(40 * time.Second)
	// }

	// Install kube-prometheus-stack
	// releaseNameStack := "prometheus-stack"
	// os.Setenv("HELM_NAMESPACE", "prometheus")
	// err = InstallChart("prometheus-stack", repoName, "kube-prometheus-stack", nil, "prometheus")
	// if err != nil {
	// 	fmt.Printf("Unable to install chart kube-prometheus-stack: %v", err)
	// } else {
	// 	fmt.Println("Successfully install chart kube-prometheus-stack")
	// 	time.Sleep(40 * time.Second)
	// }

	// Install prometheus-adapter
	// releaseNameAdapter := "prometheus-adapter"
	// value = "prometheus.url=http://" + releaseNameStack + prometheus_service + ".prometheus.svc.cluster.local"
	// err = InstallChart("prometheus-adapter", repoName, "prometheus-adapter", args, "prometheus")
	// if err != nil {
	// 	fmt.Printf("Unable to install chart prometheus-adapter: %v", err)
	// } else {
	// 	fmt.Println("Successfully install chart prometheus-adapter")
	// 	time.Sleep(40 * time.Second)
	// }
	// UnInstall charts
	// UnInstallChart(settings, "sondx12")
	// UnInstallChart(releaseNameAdapter)
}

func CreateSetting(namespace string, kubeconfig string) *cli.EnvSettings {
	os.Setenv("HELM_NAMESPACE", namespace)
	settings := cli.New()
	settings.KubeConfig = kubeconfig
	return settings
}

// RepoAdd adds repo with given name and url
func RepoAdd(name, url string) {
	repoFile := settings.RepositoryConfig

	//Ensure the file directory exists as it is required for file locking
	err := os.MkdirAll(filepath.Dir(repoFile), os.ModePerm)
	if err != nil && !os.IsExist(err) {
		log.Fatal(err)
	}

	// Acquire a file lock for process synchronization
	fileLock := flock.New(strings.Replace(repoFile, filepath.Ext(repoFile), ".lock", 1))
	lockCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	locked, err := fileLock.TryLockContext(lockCtx, time.Second)
	if err == nil && locked {
		defer fileLock.Unlock()
	}
	if err != nil {
		log.Fatal(err)
	}

	b, err := ioutil.ReadFile(repoFile)
	if err != nil && !os.IsNotExist(err) {
		log.Fatal(err)
	}

	var f repo.File
	if err := yaml.Unmarshal(b, &f); err != nil {
		log.Fatal(err)
	}

	if f.Has(name) {
		fmt.Printf("repository name (%s) already exists\n", name)
		return
	}

	c := repo.Entry{
		Name: name,
		URL:  url,
	}

	r, err := repo.NewChartRepository(&c, getter.All(settings))
	if err != nil {
		log.Fatal(err)
	}

	if _, err := r.DownloadIndexFile(); err != nil {
		err := errors.Wrapf(err, "looks like %q is not a valid chart repository or cannot be reached", url)
		log.Fatal(err)
	}

	f.Update(&c)

	if err := f.WriteFile(repoFile, 0644); err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%q has been added to your repositories\n", name)
}

// RepoUpdate updates charts for all helm repos
func RepoUpdate() {
	repoFile := settings.RepositoryConfig

	f, err := repo.LoadFile(repoFile)
	if os.IsNotExist(errors.Cause(err)) || len(f.Repositories) == 0 {
		log.Fatal(errors.New("no repositories found. You must add one before updating"))
	}
	var repos []*repo.ChartRepository
	for _, cfg := range f.Repositories {
		r, err := repo.NewChartRepository(cfg, getter.All(settings))
		if err != nil {
			log.Fatal(err)
		}
		repos = append(repos, r)
	}

	fmt.Printf("Hang tight while we grab the latest from your chart repositories...\n")
	var wg sync.WaitGroup
	for _, re := range repos {
		wg.Add(1)
		go func(re *repo.ChartRepository) {
			defer wg.Done()
			if _, err := re.DownloadIndexFile(); err != nil {
				fmt.Printf("...Unable to get an update from the %q chart repository (%s):\n\t%s\n", re.Config.Name, re.Config.URL, err)
			} else {
				fmt.Printf("...Successfully got an update from the %q chart repository\n", re.Config.Name)
			}
		}(re)
	}
	wg.Wait()
	fmt.Printf("Update Complete. ⎈ Happy Helming!⎈\n")
}

// InstallChart
func InstallChart(name, repo, chart string, namespace string) error {
	actionConfig := new(action.Configuration)
	if err := actionConfig.Init(settings.RESTClientGetter(), settings.Namespace(), os.Getenv("HELM_DRIVER"), debug); err != nil {
		// log.Fatal(err)
		return err
	}
	client := action.NewInstall(actionConfig)

	if client.Version == "" && client.Devel {
		client.Version = ">0.0.0-0"
	}
	//name, chart, err := client.NameAndChart(args)
	client.ReleaseName = name
	client.CreateNamespace = true

	cp, err := client.ChartPathOptions.LocateChart(fmt.Sprintf("%s/%s", repo, chart), settings)
	if err != nil {
		// log.Fatal(err)
		return err
	}

	debug("CHART PATH: %s\n", cp)

	p := getter.All(settings)
	valueOpts := &values.Options{}
	vals, err := valueOpts.MergeValues(p)
	if err != nil {
		// log.Fatal(err)
		return err
	}
	fmt.Printf("vals: %v", vals)
	// Add args
	args := []string{
		"operator.scannerReportTTL=4m",
		"targetNamespaces=kube-system\\,kube-public",
	}
	// mergedValues := make(map[string]interface{})

	// Merge maps using the mergo library function
	// values := mergeMaps(vals, args)
	// if err != nil {
	// 	return err
	// }
	if args != nil {
		for _, a := range args {
			if err := strvals.ParseInto(a, vals); err != nil {
				log.Fatal(errors.Wrap(err, "failed parsing --set data"))
				return err
			}
		}
	}

	// Check chart dependencies to make sure all are present in /charts
	chartRequested, err := loader.Load(cp)
	if err != nil {
		// log.Fatal(err)
		return err
	}

	validInstallableChart, err := isChartInstallable(chartRequested)
	if !validInstallableChart {
		// log.Fatal(err)
		return err
	}

	if req := chartRequested.Metadata.Dependencies; req != nil {
		// If CheckDependencies returns an error, we have unfulfilled dependencies.
		// As of Helm 2.4.0, this is treated as a stopping condition:
		// https://github.com/helm/helm/issues/2209
		if err := action.CheckDependencies(chartRequested, req); err != nil {
			if client.DependencyUpdate {
				man := &downloader.Manager{
					Out:              os.Stdout,
					ChartPath:        cp,
					Keyring:          client.ChartPathOptions.Keyring,
					SkipUpdate:       false,
					Getters:          p,
					RepositoryConfig: settings.RepositoryConfig,
					RepositoryCache:  settings.RepositoryCache,
				}
				if err := man.Update(); err != nil {
					// log.Fatal(err)
					return err
				}
			} else {
				// log.Fatal(err)
				return err
			}
		}
	}

	client.Namespace = namespace
	release, err := client.Run(chartRequested, vals)
	if err != nil {
		if err.Error() == "cannot re-use a name that is still in use" {
			fmt.Println("already installed")
		} else {
			// log.Fatal(err)
			return err
		}
	} else {
		fmt.Println(release.Manifest)
	}
	return nil
}
func mergeMaps(map1, map2 map[string]interface{}) map[string]interface{} {
	merged := make(map[string]interface{})

	// Add all key-value pairs from map1 to merged
	for key, value := range map1 {
		merged[key] = value
	}

	// Add all key-value pairs from map2 to merged (overwriting existing keys)
	for key, value := range map2 {
		merged[key] = value
	}

	return merged
}
func UnInstallChart(settings *cli.EnvSettings, name string) {
	actionConfig := new(action.Configuration)
	if err := actionConfig.Init(settings.RESTClientGetter(), settings.Namespace(), os.Getenv("HELM_DRIVER"), debug); err != nil {
		log.Fatal(err)
	}
	client := action.NewUninstall(actionConfig)
	_, err := client.Run(name)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Printf("Release %s has been uninstalled\n", name)

}

func isChartInstallable(ch *chart.Chart) (bool, error) {
	switch ch.Metadata.Type {
	case "", "application":
		return true, nil
	}
	return false, errors.Errorf("%s charts are not installable", ch.Metadata.Type)
}

func debug(format string, v ...interface{}) {
	format = fmt.Sprintf("[debug] %s\n", format)
	log.Output(2, fmt.Sprintf(format, v...))
}
func downloadKubecfg(name string, credential map[string]string) (string, error) {
	ctx := context.Background()
	minioClient, err := minio.New(credential["endPoint"], &minio.Options{
		Creds:  credentials.NewStaticV4(credential["accessKeyID"], credential["secretAccessKey"], ""),
		Secure: true,
		Region: "us-east-1",
	})
	if err != nil {
		return "", err
	}
	pathDownload := "kubecfg/" + name
	objectName := name + "-kubeconfig"
	err = minioClient.FGetObject(ctx, credential["bucketName"], objectName, pathDownload, minio.GetObjectOptions{})
	if err != nil {
		return "", err
	}
	return pathDownload, nil
}
func buildConfig(kubeconfigFile string) (*rest.Config, error) {
	var kubeconfig *string
	if home := homedir.HomeDir(); home != "" {
		kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", kubeconfigFile), "Path to the kubeconfig file")
	} else {
		kubeconfig = flag.String("kubeconfig", "", "Path to the kubeconfig file")
	}
	flag.Parse()

	if *kubeconfig == "" {
		if _, err := os.Stat("/var/run/secrets/kubernetes.io/serviceaccount/token"); err == nil {
			return rest.InClusterConfig()
		}
	}

	return clientcmd.BuildConfigFromFlags("", *kubeconfig)
}

// func installChart(value string, repoName string, chartName string, releaseName string, namespace string) error {
// 	actionConfig := new(action.Configuration)
// 	if err := actionConfig.Init(settings.RESTClientGetter(), namespace, os.Getenv("HELM_DRIVER"), debug); err != nil {
// 		log.Fatal(err)
// 	}
// 	client := action.NewInstall(actionConfig)

// 	if client.Version == "" && client.Devel {
// 		client.Version = ">0.0.0-0"
// 	}
// 	//name, chart, err := client.NameAndChart(args)
// 	client.ReleaseName = releaseName
// 	client.CreateNamespace = true
// 	// client.Wait = true
// 	client.Namespace = namespace

// 	cp, err := client.ChartPathOptions.LocateChart(fmt.Sprintf("%s/%s", repoName, chartName), settings)
// 	if err != nil {
// 		return err
// 	}

// 	debug("CHART PATH: %s\n", cp)

// 	p := getter.All(settings)
// 	valueOpts := &values.Options{}
// 	vals, err := valueOpts.MergeValues(p)
// 	if err != nil {
// 		return err
// 	}

// 	// Add args
// 	if value != "" {
// 		args := map[string]string{
// 			"set": value,
// 		}
// 		if err := strvals.ParseInto(args["set"], vals); err != nil {
// 			return errors.Wrap(err, "failed parsing --set data")
// 		}
// 	}

// 	// Check chart dependencies to make sure all are present in /charts
// 	chartRequested, err := loader.Load(cp)
// 	if err != nil {
// 		return err
// 	}

// 	validInstallableChart, err := isChartInstallable(chartRequested)
// 	if !validInstallableChart {
// 		return err
// 	}

// 	if req := chartRequested.Metadata.Dependencies; req != nil {
// 		// If CheckDependencies returns an error, we have unfulfilled dependencies.
// 		// As of Helm 2.4.0, this is treated as a stopping condition:
// 		// https://github.com/helm/helm/issues/2209
// 		if err := action.CheckDependencies(chartRequested, req); err != nil {
// 			if client.DependencyUpdate {
// 				man := &downloader.Manager{
// 					Out:              os.Stdout,
// 					ChartPath:        cp,
// 					Keyring:          client.ChartPathOptions.Keyring,
// 					SkipUpdate:       false,
// 					Getters:          p,
// 					RepositoryConfig: settings.RepositoryConfig,
// 					RepositoryCache:  settings.RepositoryCache,
// 				}
// 				if err := man.Update(); err != nil {
// 					return err
// 				}
// 			} else {
// 				return err
// 			}
// 		}
// 	}

// 	_, err = client.Run(chartRequested, vals)
// 	if err != nil {
// 		if err.Error() == "cannot re-use a name that is still in use" {
// 			return fmt.Errorf("already installed")
// 		} else {
// 			return err
// 		}
// 	} else {
// 		// fmt.Println(release.Manifest)
// 		return nil
// 	}
// }
