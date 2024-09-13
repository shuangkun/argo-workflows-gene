package wdl

import (
	"fmt"
	"github.com/antlr/antlr4/runtime/Go/antlr/v4"
	wfv1 "github.com/argoproj/argo-workflows/v3/pkg/apis/workflow/v1alpha1"
	"github.com/shuangkun/argo-workflows-gene/pkg/parsers/v1_0"
	log "github.com/sirupsen/logrus"
	apiv1 "k8s.io/api/core/v1"
	"strings"
)

func GenerateWorkflow(wdlFile string, parameterFile string) wfv1.Workflow {
	fmt.Printf("wdlFile: %v\n", wdlFile)
	builder, err := NewWdlBuilder(wdlFile)
	if err != nil {
		fmt.Printf("Error occurred when importing \"%v\": %v\n", wdlFile, err)
	}
	fmt.Printf("builder: %v\n", builder.Version)
	document, err := builder.ParseDocument()
	if err != nil {
		fmt.Printf("Failed to parse the WDL document. Reason: %v\n", err)
	}
	wf := GetArgoWorkflow(document)
	return wf
}

func NewWdlV1_0Parser(data string) *v1_0.WdlV1Parser {
	input := antlr.NewInputStream(data)
	lexer := v1_0.NewWdlV1Lexer(input)
	stream := antlr.NewCommonTokenStream(lexer, 0)
	p := v1_0.NewWdlV1Parser(stream)
	return p
}

func GetArgoWorkflow(document *Document) wfv1.Workflow {
	log.Infof("document: %v\n", document)

	wf := wfv1.Workflow{}
	if document.Workflow == nil {
		return wf
	}
	wdlwf := document.Workflow

	wf.APIVersion = "argoproj.io/v1alpha1"
	wf.Kind = "Workflow"
	wf.Name = string(wdlwf.Name)
	spec := wfv1.WorkflowSpec{}

	var entryPoint wfv1.Template
	entryPoint.Name = "main"

	spec.Entrypoint = entryPoint.Name

	// wdl inputs
	spec.Arguments = GetInputsFromDocument(wdlwf)

	// wdl tasks
	tasks := GetTaskFromDocument(document)
	for _, task := range tasks {
		spec.Templates = append(spec.Templates, task)
	}

	var callSteps []wfv1.ParallelSteps
	for _, innerWorkflow := range wdlwf.InnerWorkflows {
		var tmpParralel wfv1.ParallelSteps
		for _, call := range innerWorkflow.Calls {
			log.Infof("call: %v", call)
			if call.IsValid() {
				tmpParralel = GetParralelStepFromCalls(string(call))
			}
			log.Infof("tmpParralel: %v", tmpParralel)
		}
		callSteps = append(callSteps, tmpParralel)
	}
	entryPoint.Steps = callSteps
	spec.Templates = append(spec.Templates, entryPoint)
	wf.Spec = spec
	return wf
}

func GetTaskFromDocument(document *Document) map[string]wfv1.Template {
	tasks := make(map[string]wfv1.Template)

	log.Infof("document tasks: %v\n", document.Tasks)
	for _, task := range document.Tasks {
		var template wfv1.Template
		template.Name = string(task.Name)
		log.Infof("task runtime: %v", task.Runtime)
		str, _ := task.Runtime["container"].GetValue().(string)
		var container apiv1.Container
		container.Image = strings.TrimSpace(str)
		str = task.Command
		str = strings.TrimSpace(str)
		str = strings.ReplaceAll(str, "\n", "")
		command := str
		container.Command = []string{"sh", "-c"}
		container.Command = append(container.Command, command)
		template.Container = &container
		template.Inputs = GetInputsFromTask(task)
		template.Outputs = GetOutputsFromTask(task)
		tasks[string(task.Name)] = template
	}

	log.Infof("tasks: %v\n", tasks)
	return tasks
}

func GetInputsFromTask(task Task) wfv1.Inputs {
	var inputs wfv1.Inputs
	var parameters []wfv1.Parameter
	for _, declaration := range task.Input {
		var parameter wfv1.Parameter
		parameter.Name = string(declaration.Identifier)
		if declaration.Expr != nil {
			str := IExpression(*declaration.Expr).GetValue()
			parameter.Value = wfv1.AnyStringPtr(str)
		}
		parameters = append(parameters, parameter)
	}
	inputs.Parameters = parameters
	return inputs
}

func GetOutputsFromTask(task Task) wfv1.Outputs {
	var outputs wfv1.Outputs
	var parameters []wfv1.Parameter
	for _, declaration := range task.Output {
		var parameter wfv1.Parameter
		parameter.Name = string(declaration.Identifier)
		if declaration.Expr != nil {
			str := IExpression(*declaration.Expr).GetValue()
			parameter.Value = wfv1.AnyStringPtr(str)
		}
		parameters = append(parameters, parameter)
	}
	return outputs
}

func GetInputsFromDocument(wdlWorkflow *Workflow) wfv1.Arguments {
	inputs := wfv1.Arguments{}
	var parameters []wfv1.Parameter
	for _, declaration := range wdlWorkflow.Inputs {
		var parameter wfv1.Parameter
		parameter.Name = string(declaration.Identifier)
		if declaration.Expr != nil {
			str := IExpression(*declaration.Expr).GetValue()
			parameter.Value = wfv1.AnyStringPtr(str)
		} else {
			parameter.Value = wfv1.AnyStringPtr("")
		}
		parameters = append(parameters, parameter)
	}
	inputs.Parameters = parameters
	return inputs
}

func GetParralelStepFromCalls(call string) wfv1.ParallelSteps {
	parallelSteps := wfv1.ParallelSteps{}
	var steps []wfv1.WorkflowStep
	var oneStep wfv1.WorkflowStep
	oneStep.Name = call
	oneStep.Template = call
	steps = append(steps, oneStep)
	parallelSteps.Steps = steps
	return parallelSteps
}
