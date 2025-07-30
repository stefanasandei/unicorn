package application

import (
	"binomeway.com/common/broker"
	common "binomeway.com/common/model"
	"binomeway.com/worker/model"
	"encoding/json"
	"log"
)

func (app *App) HandleQueueMessage(msg broker.DeliveryMessage) error {
	var execReq common.ExecutionRequestWrapper

	// get the struct out of the json
	err := json.Unmarshal([]byte(msg.Body), &execReq)
	if err != nil {
		return err
	}

	languages, err := model.QuerySupportedLanguages(app.config.RuntimesDir)
	if err != nil {
		log.Printf("Failed to query languages: %s", err)
		return err
	}

	task, err := model.NewTask(execReq, languages)
	if err != nil {
		errMsg := common.WorkerResponseWrapper{
			Id: execReq.Id,
			Res: common.WorkerResponse{
				Compile: common.ProcessResult{
					Output:   "Language not found",
					ExitCode: 1,
				},
			},
		}

		msg, _ := json.Marshal(errMsg)

		err = app.broker.SendMessageToQueue("reply", string(msg))
		if err != nil {
			return err
		}

		return nil
	}

	result, err := task.Execute()
	if err != nil {
		return err
	}

	resultWrapper := common.WorkerResponseWrapper{
		Id:  task.ID,
		Res: result,
	}

	jsonResult, err := json.Marshal(resultWrapper)
	if err != nil {
		return err
	}

	err = app.broker.SendMessageToQueue("reply", string(jsonResult))
	if err != nil {
		return err
	}

	return nil
}
