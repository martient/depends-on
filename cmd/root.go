package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/martient/depends-on/pkg"
	"github.com/martient/golang-utils/utils"
	"github.com/spf13/cobra"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

type arrayFlags []string

func (i *arrayFlags) String() string {
	return strings.Join(*i, ",")
}
func (i *arrayFlags) Set(value string) error {
	*i = append(*i, strings.TrimSpace(value))
	return nil
}

// Add this method
func (i *arrayFlags) Type() string {
	return "string"
}

var jobs arrayFlags
var services arrayFlags
var checkInterval int
var BEMversion string

var rootCmd = &cobra.Command{
	Use:   "depends-on",
	Short: "depends-on the wait helper for Kubernetes jobs and/or services to be ready !",
	Run: func(cmd *cobra.Command, args []string) {
		if len(jobs) == 0 && len(services) == 0 {
			fmt.Println("No jobs or services provided. Please follow the docs")
			cmd.Help()
			os.Exit(1)
		}

		nsb, err := os.ReadFile("/var/run/secrets/kubernetes.io/serviceaccount/namespace")
		if err != nil {
			panic(err.Error())
		}
		ns := string(nsb)
		utils.LogInfo("Determined namespace: %s\n", "CMD", ns)
		utils.LogInfo("Creating client...", "CMD")
		config, err := rest.InClusterConfig()
		if err != nil {
			panic(err.Error())
		}
		clientset, err := kubernetes.NewForConfig(config)
		if err != nil {
			panic(err.Error())
		}

		// Wait for jobs
		for _, jobName := range jobs {
			if err := pkg.WaitForJob(clientset, ns, jobName, checkInterval); err != nil {
				panic(err)
			}
		}

		// Wait for services
		for _, serviceName := range services {
			if err := pkg.WaitForService(clientset, ns, serviceName, checkInterval); err != nil {
				panic(err)
			}
		}
	},
}

// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute(versionFormated string, version string) {
	rootCmd.Version = versionFormated
	BEMversion = version
	cobra.CheckErr(rootCmd.Execute())
}

func init() {
	rootCmd.Flags().Var(&jobs, "job", "Job, which successful completion to wait for. Can be specified multiple times")
	rootCmd.Flags().Var(&services, "service", "Service, which pods to wait for. Can be specified multiple times")
	rootCmd.Flags().IntVar(&checkInterval, "check-interval", 5, "Seconds to wait between each check attempts")
	// Add this line
	rootCmd.Flags().SortFlags = false
}
