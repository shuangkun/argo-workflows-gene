package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/argoproj/argo-workflows/v3/cmd/argo/commands/common"
	"github.com/argoproj/pkg/errors"
	"github.com/shuangkun/argo-workflows-gene/pkg/wdl"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"os"

	"github.com/argoproj/argo-workflows/v3/cmd/argo/commands/client"
	workflowpkg "github.com/argoproj/argo-workflows/v3/pkg/apiclient/workflow"
	wfv1 "github.com/argoproj/argo-workflows/v3/pkg/apis/workflow/v1alpha1"
	"github.com/argoproj/argo-workflows/v3/workflow/util"
)

func NewSubmitCommand() *cobra.Command {
	var wdlFile string
	var parameterFile string

	command := &cobra.Command{
		Use:   "submit ",
		Short: "submit a wdl workflow",
		Run: func(cmd *cobra.Command, args []string) {

			ctx, apiClient := client.NewAPIClient(cmd.Context())
			serviceClient := apiClient.NewWorkflowServiceClient()
			namespace := client.Namespace()
			submitWorkflowsFromFile(ctx, serviceClient, namespace, wdlFile, parameterFile)

		},
	}
	command.Flags().StringVar(&wdlFile, "file", "", "wdl file")
	command.Flags().StringVar(&parameterFile, "parameter", "", "parameter file")
	return command
}

func submitWorkflowsFromFile(ctx context.Context, serviceClient workflowpkg.WorkflowServiceClient, namespace string, wdlFile string, parameterFile string) {

	var workflows []wfv1.Workflow
	wf := wdl.GenerateWorkflow(wdlFile, parameterFile)
	printWorkflow(&wf, common.GetFlags{Output: "json"})
	//fmt.Printf("Generated workflows from wdl-Workflows: %s\n", wf)
	workflows = append(workflows, wf)
	submitWorkflows(ctx, serviceClient, namespace, workflows, &wfv1.SubmitOpts{}, &common.CliSubmitOpts{})
}

func validateOptions(workflows []wfv1.Workflow, submitOpts *wfv1.SubmitOpts, cliOpts *common.CliSubmitOpts) {
	if cliOpts.Watch {
		if len(workflows) > 1 {
			log.Fatalf("Cannot watch more than one workflow")
		}
		if cliOpts.Wait {
			log.Fatalf("--wait cannot be combined with --watch")
		}
		if submitOpts.DryRun {
			log.Fatalf("--watch cannot be combined with --dry-run")
		}
		if submitOpts.ServerDryRun {
			log.Fatalf("--watch cannot be combined with --server-dry-run")
		}
	}

	if cliOpts.Wait {
		if submitOpts.DryRun {
			log.Fatalf("--wait cannot be combined with --dry-run")
		}
		if submitOpts.ServerDryRun {
			log.Fatalf("--wait cannot be combined with --server-dry-run")
		}
	}

	if submitOpts.DryRun {
		if cliOpts.Output == "" {
			log.Fatalf("--dry-run should have an output option")
		}
		if submitOpts.ServerDryRun {
			log.Fatalf("--dry-run cannot be combined with --server-dry-run")
		}
	}

	if submitOpts.ServerDryRun {
		if cliOpts.Output == "" {
			log.Fatalf("--server-dry-run should have an output option")
		}
	}
}

func submitWorkflows(ctx context.Context, serviceClient workflowpkg.WorkflowServiceClient, namespace string, workflows []wfv1.Workflow, submitOpts *wfv1.SubmitOpts, cliOpts *common.CliSubmitOpts) {
	validateOptions(workflows, submitOpts, cliOpts)

	if len(workflows) == 0 {
		log.Println("No Workflow found in given files")
		os.Exit(1)
	}

	var workflowNames []string

	for _, wf := range workflows {
		if wf.Namespace == "" {
			// This is here to avoid passing an empty namespace when using --server-dry-run
			wf.Namespace = namespace
		}
		err := util.ApplySubmitOpts(&wf, submitOpts)
		errors.CheckError(err)
		if cliOpts.Priority != nil {
			wf.Spec.Priority = cliOpts.Priority
		}
		options := &metav1.CreateOptions{}
		if submitOpts.DryRun {
			options.DryRun = []string{"All"}
		}
		created, err := serviceClient.CreateWorkflow(ctx, &workflowpkg.WorkflowCreateRequest{
			Namespace:     wf.Namespace,
			Workflow:      &wf,
			ServerDryRun:  submitOpts.ServerDryRun,
			CreateOptions: options,
		})
		if err != nil {
			log.Fatalf("Failed to submit workflow: %v", err)
		}

		printWorkflow(created, common.GetFlags{Output: cliOpts.Output, Status: cliOpts.GetArgs.Status})
		workflowNames = append(workflowNames, created.Name)
	}

	common.WaitWatchOrLog(ctx, serviceClient, namespace, workflowNames, *cliOpts)
}

func printWorkflow(wf *wfv1.Workflow, getArgs common.GetFlags) {
	switch getArgs.Output {
	case "name":
		fmt.Println(wf.ObjectMeta.Name)
	case "json":
		outBytes, _ := json.MarshalIndent(wf, "", "    ")
		fmt.Println(string(outBytes))
	case "yaml":
		outBytes, _ := yaml.Marshal(wf)
		fmt.Print(string(outBytes))
	case "short", "wide", "":
		fmt.Print(common.PrintWorkflowHelper(wf, getArgs))
	default:
		log.Fatalf("Unknown output format: %s", getArgs.Output)
	}
}
