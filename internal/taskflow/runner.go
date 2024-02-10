package taskflow

import (
	"fmt"
	"os"

	"github.com/k0kubun/pp/v3"

	"github.com/webhkp/godft/internal/consts"
	mongoDriver "github.com/webhkp/godft/internal/driver/database/mongo"
	mysqlDriver "github.com/webhkp/godft/internal/driver/database/mysql"
	jsondriver "github.com/webhkp/godft/internal/driver/file/json"
	"github.com/webhkp/godft/internal/driver/generator"
	"github.com/webhkp/godft/pkg/config"
)

type TaskFlowRunnter struct {
	filePaths     *[]string
	taskFlows     map[string]map[string]*Task
	rootTaskFlows []*Task
}

func NewTaskFlowRunner(files *[]string) (taskFlow *TaskFlowRunnter) {
	taskFlow = &TaskFlowRunnter{}
	taskFlow.filePaths = files
	taskFlow.taskFlows = make(map[string]map[string]*Task)

	return
}

func (tfr *TaskFlowRunnter) Process() {
	configFiles := tfr.getAllFileList()

	tfr.combine(configFiles)
	tfr.createTaskFlowTree()
	tfr.runAll()

}

func (tfr *TaskFlowRunnter) getAllFileList() (allTaskFlowFiles []string) {
	for _, filePath := range *tfr.filePaths {
		fileInfo, err := os.Stat(filePath)

		if err != nil {
			fmt.Println("Error processing: " + filePath)
			continue
		}

		// Check if it is a file or directory
		if fileInfo.IsDir() {
			// fmt.Println(arg + "-dir")
		} else {
			allTaskFlowFiles = append(allTaskFlowFiles, filePath)
		}
	}

	return
}

func (tfr *TaskFlowRunnter) combine(files []string) {
	for _, file := range files {
		config := config.Read(file)

		currentNamespace := consts.DefaultNamespace

		for key, val := range config {
			if key == consts.NamespaceKey {
				currentNamespace = val.(string)

				continue
			}

			var currentDriver string
			parsedVal := make(map[string]interface{})

			for iKey, iVal := range val.(map[interface{}]interface{}) {
				if iKey.(string) == consts.DriverKey {
					currentDriver = iVal.(string)
					continue
				}

				if iKey.(string) == consts.NamespaceKey {
					currentNamespace = iVal.(string)
					continue
				}

				parsedVal[iKey.(string)] = iVal
			}

			if _, ok := tfr.taskFlows[currentNamespace]; !ok {
				tfr.taskFlows[currentNamespace] = make(map[string]*Task)
			}

			switch currentDriver {
			case "generator":
				tfr.taskFlows[currentNamespace][key] = NewTask(currentNamespace, key, generator.NewGenerator(parsedVal))
			case "json":
				tfr.taskFlows[currentNamespace][key] = NewTask(currentNamespace, key, jsondriver.NewJson(parsedVal))
			case "mongo":
				tfr.taskFlows[currentNamespace][key] = NewTask(currentNamespace, key, mongoDriver.NewMongo(parsedVal))
			case "mysql":
				tfr.taskFlows[currentNamespace][key] = NewTask(currentNamespace, key, mysqlDriver.NewMysql(parsedVal))
			}
		}
	}
}

func (tfr *TaskFlowRunnter) createTaskFlowTree() {
	for namespace, tasks := range tfr.taskFlows {
		for _, task := range tasks {
			input, hasInput := task.Driver.GetInput()

			if !hasInput {
				tfr.rootTaskFlows = append(tfr.rootTaskFlows, task)
				continue
			}

			if _, ok := tfr.taskFlows[namespace][input]; ok {
				tfr.taskFlows[namespace][input].AddNextTasks(task)
			}
		}
	}
}

func (tfr *TaskFlowRunnter) runAll() {
	for _, taskFlow := range tfr.rootTaskFlows {
		var data consts.FlowDataSet = make(consts.FlowDataSet)

		pp.Println("Processing namespace: " + taskFlow.Namespace)
		tfr.run(taskFlow, &data)
		pp.Println("Completed processing namespace: " + taskFlow.Namespace)
	}
}

func (tfr *TaskFlowRunnter) run(task *Task, data *consts.FlowDataSet) {
	pp.Printf("\tProcessing task: %s\n", task.Key)

	task.Driver.Execute(data)

	pp.Printf("\tCompleted processing task: %s\n\n", task.Key)

	for _, childTask := range task.GetNextTasks() {
		tfr.run(childTask, data)
	}

}
