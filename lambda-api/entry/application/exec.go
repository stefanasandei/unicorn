package application

import (
	"encoding/json"
	"errors"
	"fmt"

	common "github.com/common/model"
	"github.com/google/uuid"
)

func (a *App) ExecuteCode(req common.ExecutionRequest) (common.ResponseTask, error) {
	workerId, err := ChooseWorker(a.redisDB, a.ctx.backgroundCtx)
	if err != nil {
		fmtErr := fmt.Errorf("failed to choose a worker: %s", err)
		return common.ResponseTask{}, fmtErr
	}

	// send task to worker
	wrapperMsg := common.ExecutionRequestWrapper{
		Id:  uuid.New(),
		Req: req,
	}

	brokerMsg, err := json.Marshal(wrapperMsg)
	err = a.broker.SendMessageToQueue(workerId, string(brokerMsg))
	if err != nil {
		fmtErr := fmt.Errorf("failed to send message to the queue: %s", err)
		return common.ResponseTask{}, fmtErr
	}

	// listen to the reply queue and wait for the right message to come
	a.ctx.messageMutex.Lock()
	workerResponse, err := a.getBackMessage(wrapperMsg.Id)
	a.ctx.messageMutex.Unlock()

	if err != nil {
		fmtErr := fmt.Errorf("failed to get back the worker message: %s", err)
		return common.ResponseTask{
			Status: common.StatusFailed,
			Output: common.WorkerResponse{},
		}, fmtErr
	}

	// write the response
	task := common.ResponseTask{
		Status: common.StatusDone,
		Output: workerResponse,
	}

	if workerResponse.Run.ExitCode != 0 || workerResponse.Compile.ExitCode != 0 {
		task.Status = common.StatusError
	}

	return task, nil
}

// TODO: use generics
func (a *App) getBackMessage(id uuid.UUID) (common.WorkerResponse, error) {
	for _, msg := range a.ctx.incomingMessages {
		if msg.Id == id {
			return a.unwrapWorkerMessage(msg)
		}
	}

	// loop until we get the right message
	for limit := 0; limit < 5; limit++ {
		for _, msg := range a.ctx.incomingMessages {
			if msg.Id == id {
				return a.unwrapWorkerMessage(msg)
			}
		}

		msg, err := a.reply.GetNewMessage()
		if err != nil {
			fmtErr := fmt.Errorf("failed to read a new message from the queue: %s", err)
			return common.WorkerResponse{}, fmtErr
		}

		// deserialize into a wrapper and add that to the msg array of the context
		workerResponseWrapper := common.WorkerResponseWrapper{}
		err = json.Unmarshal([]byte(msg), &workerResponseWrapper)
		if err != nil {
			fmtErr := fmt.Errorf("failed to deserialize the worker response: %s", err)
			return common.WorkerResponse{}, fmtErr
		}

		if workerResponseWrapper.Id == id {
			return a.unwrapWorkerMessage(workerResponseWrapper)
		}

		a.ctx.incomingMessages = append(a.ctx.incomingMessages, workerResponseWrapper)
	}

	return common.WorkerResponse{}, errors.New("the message ain't coming")
}

func (a *App) unwrapWorkerMessage(msg common.WorkerResponseWrapper) (common.WorkerResponse, error) {
	// return the response within the wrapper
	return msg.Res, nil
}
