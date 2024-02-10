package util

func ProcessInterfaceMap(data map[interface{}]interface{}) (result map[string]interface{}) {
	result = make(map[string]interface{})

	for key, val := range data {
		result[key.(string)] = val
	}

	return
}
