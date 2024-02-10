package generator

import (
	"reflect"
	"strings"
	"testing"

	"github.com/webhkp/godft/internal/consts"
)

var config = map[string]interface{}{
	"collection": map[interface{}]interface{}{
		"customer": map[interface{}]interface{}{
			"limit": 5,
			"field": map[interface{}]interface{}{
				"full_name": "name",
				"age":       "int:10:90",
			},
		},
	},
}

func TestNewGenerator(t *testing.T) {
	generator := NewGenerator(config)
	customer, customerKeyOk := generator.Collections["customer"]

	if !customerKeyOk {
		t.Error("\"customer\" key should be in present in the generator")
	}

	if customer.Limit != 5 {
		t.Errorf("Generator row limit is wrong: customer.Limit = %d, want %d", customer.Limit, 5)
	}

	if len(customer.Fields) != 2 {
		t.Errorf("Number of customer fields is wrong: length of customer.Fields is %d, want %d", len(customer.Fields), 2)
	}
}

func TestRead(t *testing.T) {
	generator := NewGenerator(config)
	flowDataSet := make(consts.FlowDataSet)

	generator.Read(&flowDataSet)

	customers, customerOk := flowDataSet["customer"]

	if !customerOk {
		t.Error("customer not present in the generated data")
	}

	if len(customers) != 5 {
		t.Errorf("Count of customer data is %d, want %d", len(customers), 5)
	}

	for _, val := range customers {
		if len(val) != 2 {
			t.Errorf("Count of customer field is %d, want %d", len(val), 2)
		}

		if age, ageOk := val["age"]; ageOk {
			if reflect.TypeOf(age).Kind() != reflect.Int {
				t.Errorf("Type of age is %q, want %q", reflect.TypeOf(age).Kind(), reflect.Int)
			}

			if age.(int) < 10 {
				t.Errorf("Customer age is %d, want > %d", age.(int), 10)
			}

			if age.(int) > 90 {
				t.Errorf("Customer age is %d, want < %d", age.(int), 90)
			}
		} else {
			t.Error("age is not present in customer data")
		}

		if _, nameOk := val["full_name"]; !nameOk {
			t.Error("full_name is not present in customer data")
		}

	}
}

func TestGenerateFieldData(t *testing.T) {
	wrongTypeVal := generateFieldData("wrong_type")

	if wrongTypeVal != "wrong_type" {
		t.Errorf("wrong_type = %s, want %s", wrongTypeVal, "wrong_type")
	}

	nameTypeVal := generateFieldData("name").(string)

	if !strings.Contains(nameTypeVal, " ") {
		t.Errorf("name = %s, want space in the name", nameTypeVal)
	}

	intTypeVal := generateFieldData("number:100:200").(int)

	if intTypeVal < 100 {
		t.Errorf("intTypeVal = %d, want >= %d", intTypeVal, 100)
	}

	if intTypeVal > 200 {
		t.Errorf("intTypeVal = %d, want <= %d", intTypeVal, 200)
	}

	decimalTypeVal := generateFieldData("decimal::99").(float64)

	if decimalTypeVal > 99 {
		t.Errorf("decimalTypeVal = %f, want < %f", decimalTypeVal, 99.0)
	}
}
