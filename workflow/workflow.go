package workflow

import (
	"context"
	"io/ioutil"
	"time"

	"go.uber.org/cadence/activity"
	"go.uber.org/cadence/workflow"
)

func init() {
	workflow.Register(DemoWorkFlow)
	activity.Register(getNameActivity)
	activity.Register(sayHello)
	//activity.Register(persisteResult)
}

func getNameActivity() (string, error) {
	return "cadence", nil
}

func sayHello(name string) (string, error) {
	return "Helloooooooooooooooooooooooooooo " + name + "!!!!!!!!!!!!!!!!!!!!!!!", nil
}

func persisteResult(ctx context.Context, data string) error {
	runID := activity.GetInfo(ctx).WorkflowExecution.RunID
	fileName := "/Users/bala/temp/" + runID
	return ioutil.WriteFile(fileName, []byte(runID+"_"+data), 0666)
}

//DemoWorkFlow comment
func DemoWorkFlow(ctx workflow.Context) error {
	var err error
	ao := workflow.ActivityOptions{
		StartToCloseTimeout:    time.Minute,
		ScheduleToStartTimeout: time.Minute,
		HeartbeatTimeout:       time.Second * 20,
	}
	ctx = workflow.WithActivityOptions(ctx, ao)

	var name string
	err = workflow.ExecuteActivity(ctx, getNameActivity).Get(ctx, &name)
	if err != nil {
		return nil
	}

	var result string
	err = workflow.ExecuteActivity(ctx, sayHello, name).Get(ctx, &result)
	if err != nil {
		return nil
	}

	// err = retry(func() error {
	// 	return workflow.ExecuteActivity(ctx, persisteResult, result).Get(ctx, nil)
	// })
	// if err != nil {
	// 	return nil
	// }

	workflow.GetLogger(ctx).Info("Result " + result)

	return nil
}

func retry(op func() error) error {
	var err error
	for i := 0; i < 2; i++ {
		if err = op(); err == nil {
			return nil
		}
	}
	return nil
}
