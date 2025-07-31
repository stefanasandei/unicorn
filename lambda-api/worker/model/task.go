package model

import (
	"errors"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/common/model"
	"github.com/google/uuid"
)

type Task struct {
	ID       uuid.UUID              `json:"id"`
	Request  model.ExecutionRequest `json:"request"`
	Language model.Language         `json:"language"`

	WorkingDir string
}

type CommandSpec struct {
	EntryFilename string
	OutFilename   string
	ProcLimits    model.ProcessInfo
}

func NewTask(req model.ExecutionRequestWrapper, langsRepo []model.Language) (Task, error) {
	task := Task{
		ID:      req.Id,
		Request: req.Req,
	}

	// search for the language that matches the runtime name
	found := false

	for index, lang := range langsRepo {
		if lang.Name == req.Req.Runtime.Name {
			task.Language = langsRepo[index]
			found = true
			break
		}
	}

	if !found {
		return task, errors.New("language not found")
	}

	return task, nil
}

func (t *Task) Execute() (model.WorkerResponse, error) {
	execDir, err := t.initWorkingDir()
	if err != nil {
		_ = t.cleanup()
		return model.WorkerResponse{}, err
	}

	// write the other files (which are not be compiled, only resources)
	for _, file := range t.Request.Project.Files {
		err := writeFile(file.Name, file.Contents, execDir)
		if err != nil {
			_ = t.cleanup()
			return model.WorkerResponse{}, err
		}
	}

	// write the source code of the entry file to ./tasks/exec_id/main
	entryFilename := fmt.Sprintf("%s.%s", "main", t.Language.Extension)
	err = writeFile(entryFilename, t.Request.Project.Entry, execDir)
	if err != nil {
		_ = t.cleanup()
		return model.WorkerResponse{}, err
	}

	entryFullFilename := fmt.Sprintf("%s/%s", execDir, entryFilename)
	outputFilename := fmt.Sprintf("%s/%s", execDir, "main")

	// setup the cmd spec, used to run & compile
	cmdSpec := CommandSpec{
		EntryFilename: entryFullFilename,
		OutFilename:   outputFilename,
		ProcLimits:    t.Request.Process,
	}

	compileProcess, err := t.compileFile(&cmdSpec)

	if err != nil || compileProcess.ExitCode != 0 {
		_ = t.cleanup()
		return model.WorkerResponse{
			Compile: compileProcess,
		}, nil
	}

	// run the code
	cmdSpec.EntryFilename = entryFilename
	cmdSpec.OutFilename = "./main" // TODO
	cmdSpec.ProcLimits.WorkingDirectory = execDir
	runProcess, err := t.runFile(cmdSpec)

	//	return model.WorkerResponse{Compile: model.ProcessResult{ExitCode: 1}}, nil

	if err != nil || runProcess.ExitCode != 0 {
		log.Printf("error: %s\n", err)

		_ = t.cleanup()
		return model.WorkerResponse{
			Compile: compileProcess,
			Run:     runProcess,
		}, nil
	}

	err = t.cleanup()
	if err != nil {
		return model.WorkerResponse{}, err
	}

	return model.WorkerResponse{
		Compile: compileProcess,
		Run:     runProcess,
	}, err
}

func (t *Task) compileFile(spec *CommandSpec) (model.ProcessResult, error) {
	if len(t.Language.CompileCmds) == 0 {
		spec.OutFilename = spec.EntryFilename
		return model.ProcessResult{}, nil
	}

	compileCommands := t.processCommands(spec.EntryFilename, spec.OutFilename, t.Language.CompileCmds)

	result, err := ExecuteSystemCommand(compileCommands, DefaultLimits())

	return result, err
}

func (t *Task) runFile(spec CommandSpec) (model.ProcessResult, error) {
	runCommands := t.processCommands(spec.EntryFilename, spec.OutFilename, t.Language.RunCmds)

	result, err := ExecuteSystemCommand(runCommands, spec.ProcLimits)

	return result, err
}

// replace <entry> and <output> within the run/compile commands
func (t *Task) processCommands(entry, output string, commands []string) []string {
	resultCmds := make([]string, len(commands))

	for i := 0; i < len(commands); i++ {
		resultCmds[i] = commands[i]

		if strings.Contains(resultCmds[i], "<entry>") {
			resultCmds[i] = strings.ReplaceAll(resultCmds[i], "<entry>", entry)
		}

		if strings.Contains(resultCmds[i], "<output>") {
			resultCmds[i] = strings.ReplaceAll(resultCmds[i], "<output>", output)
		}
	}

	return resultCmds
}

func (t *Task) cleanup() error {
	// now delete the dir with all the task files
	err := os.RemoveAll(t.WorkingDir)
	if err != nil {
		return err
	}

	t.WorkingDir = ""

	return nil
}

func (t *Task) initWorkingDir() (string, error) {
	execId := generateExecutableID()

	execDir := fmt.Sprintf("./tasks/%s", execId)
	err := os.Mkdir(execDir, 0755)
	if err != nil {
		dir, _ := os.Getwd()
		log.Printf("Failed to create a task directory. CWD: %s", dir)
		return "", err
	}

	t.WorkingDir = execDir
	return execDir, nil
}

func generateExecutableID() string {
	unixEpoch := time.Now().UnixNano()
	randomDigits := rand.Intn(9999-1000+1) + 1000

	id := fmt.Sprintf("%d%d", randomDigits, unixEpoch)
	return id
}

func writeFile(name string, contents string, parentDir string) error {
	filename := fmt.Sprintf("%s/%s", parentDir, name)

	err := os.WriteFile(filename, []byte(contents), 0644)
	return err
}
