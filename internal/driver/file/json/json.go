package json

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"time"

	"github.com/webhkp/godft/internal/consts"
	_ "github.com/webhkp/godft/internal/driver"
	jsonUtil "github.com/webhkp/godft/pkg/json"
	"github.com/webhkp/godft/pkg/util"
)

type Json struct {
	intputTaskName string
	readOnly       bool
	inputPath      []string
	outputPath     string
}

func NewJson(config consts.FieldsType) (json *Json) {
	json = &Json{
		readOnly:   true,
		outputPath: "",
		inputPath:  []string{},
	}

	if _, ok := config[consts.InputPathKey]; ok {
		json.inputPath = make([]string, len(config[consts.InputPathKey].([]interface{})))

		for i, inputPath := range config[consts.InputPathKey].([]interface{}) {
			json.inputPath[i] = inputPath.(string)
		}
	}

	if _, ok := config[consts.OutputPathKey]; ok {
		json.outputPath = config[consts.OutputPathKey].(string) + "/"
	}

	if _, ok := config[consts.InputDriverKey]; ok {
		json.intputTaskName = config[consts.InputDriverKey].(string)
	}

	if _, ok := config[consts.ReadOnlyKey]; ok {
		json.readOnly = config[consts.ReadOnlyKey].(bool)
	}

	return
}

func (d *Json) Execute(data *consts.FlowDataSet) {
	startTime := time.Now()

	d.Read(data)
	d.Write(data)

	fmt.Printf("JSON process time: %v\n", time.Since(startTime))
}

func (d *Json) Read(data *consts.FlowDataSet) {
	if len(d.inputPath) == 0 {
		return
	}

	allJsonPath := d.getInputFileList()

	// pp.Println(allJsonPath)

	for k, v := range allJsonPath {
		type jsonDataType []consts.FlowData
		jsonData := jsonUtil.ReadJsonFile[jsonDataType](v)

		for _, jsonVal := range jsonData {
			(*data)[k] = append((*data)[k], jsonVal)
		}
	}

	// pp.Println(data)

}

func (d *Json) getInputFileList() map[string]string {
	// Get lsit of allfiles
	allFiles := make(map[string]string)

	for _, filePath := range d.inputPath {
		fileInfo, err := os.Stat(filePath)

		if err != nil {
			fmt.Println("Error processing: " + filePath)
			continue
		}

		// Check if it is a file or directory
		if fileInfo.IsDir() {
			f, err := os.Open(filePath)
			if err != nil {
				fmt.Println(err)
				continue
			}

			filesInDir, err := f.Readdir(0)
			if err != nil {
				fmt.Println(err)
				continue
			}

			for _, v := range filesInDir {
				if v.IsDir() {
					// @todo
				} else {
					allFiles[util.GetFileNameWithoutExtension(v.Name())] = filePath + "/" + v.Name()
				}
			}
		} else {
			allFiles[util.GetFileNameWithoutExtension(filePath)] = filePath
		}
	}

	return allFiles
}

func (d *Json) Write(data *consts.FlowDataSet) {
	if d.readOnly || len(d.outputPath) == 0 {
		return
	}

	for key, collection := range *data {
		jsonBytes, _ := json.Marshal(collection)

		err := os.WriteFile(d.outputPath+key+consts.DefaultJsonOutputPathSuffix, jsonBytes, fs.FileMode(0644))

		if err != nil {
			fmt.Println(err)
		}
	}
}

func (d *Json) Validate() bool {
	return true
}

func (d *Json) GetInput() (string, bool) {
	if d.intputTaskName != "" {
		return d.intputTaskName, true
	}

	return "", false
}
