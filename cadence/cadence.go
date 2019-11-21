package cadence

import (
	"cadencedemo/workflow"
	"context"
	"fmt"
	"log"
	"time"

	"github.com/uber-go/tally"
	"go.uber.org/cadence/.gen/go/cadence/workflowserviceclient"
	"go.uber.org/cadence/client"
	"go.uber.org/cadence/worker"
	"go.uber.org/yarpc"
	"go.uber.org/yarpc/transport/tchannel"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var hostPort = "192.168.99.101:7933"
var domain = "samples-domain"
var taskListName = "SampleTaskList"
var clientName = "example-worker"
var cadenceService = "cadence-frontend"

func buildLogger() *zap.Logger {
	config := zap.NewDevelopmentConfig()
	config.Level.SetLevel(zapcore.InfoLevel)

	var err error
	logger, err := config.Build()
	if err != nil {
		panic("Failed to setup logger")
	}

	return logger
}

func buildWorkflowServiceClient() workflowserviceclient.Interface {
	ch, err := tchannel.NewChannelTransport(tchannel.ServiceName(clientName))
	if err != nil {
		panic("Failed to setup tchannel")
	}
	dispatcher := yarpc.NewDispatcher(yarpc.Config{
		Name: clientName,
		Outbounds: yarpc.Outbounds{
			cadenceService: {Unary: ch.NewSingleOutbound(hostPort)},
		},
	})
	if err := dispatcher.Start(); err != nil {
		fmt.Println(err)
		panic("Failed to start dispatcher")
	}

	return workflowserviceclient.New(dispatcher.ClientConfig(cadenceService))
}

func getCadenceClient() client.Client {
	clientOptions := &client.Options{}
	return client.NewClient(
		buildWorkflowServiceClient(), domain, clientOptions,
	)
}

func StartCadenceWorker() {
	// TaskListName - identifies set of client workflows, activities and workers.
	// it could be your group or client or application name.
	var err error
	//_, err := zap.NewDevelopment(zap.AddCallerSkip(1))
	//if err != nil {
	//panic(err)
	//}
	workerOptions := worker.Options{
		Logger:       buildLogger(),
		MetricsScope: tally.NewTestScope(taskListName, map[string]string{}),
	}

	worker := worker.New(
		buildWorkflowServiceClient(),
		domain,
		taskListName,
		workerOptions)
	err = worker.Start()
	if err != nil {
		panic("Failed to start worker")
	}

	//logger.Info("Started Worker.", zap.String("worker", TaskListName))
}

func StartWorkflow() {
	ctx := context.Background()
	cadenceClient := getCadenceClient()

	workflowOptions := client.StartWorkflowOptions{
		TaskList:                     taskListName,
		ExecutionStartToCloseTimeout: time.Minute,
		//DecisionTaskStartToCloseTimeout: time.Minute,
	}

	run, err := cadenceClient.ExecuteWorkflow(ctx, workflowOptions, workflow.DemoWorkFlow)
	fmt.Println(run)
	if err != nil {
		fmt.Println(err)
		panic("Failed to start wrokflow")
	}

	log.Printf("workflow=%q, run=%q", run.GetID(), run.GetRunID())

	err = run.Get(ctx, nil)
	if err != nil {
		panic(err)
	}

	log.Println("done")
}
