package generator

import (
	"fmt"
	"math"
	"math/rand"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/go-faker/faker/v4"
	"github.com/schollz/progressbar/v3"

	"regexp"

	"github.com/webhkp/godft/internal/consts"
)

type fieldsType map[string]interface{}

type Generator struct {
	Collections map[string]Collection
}

func NewGenerator(config fieldsType) (generator *Generator) {
	generator = &Generator{}

	if _, ok := config[consts.CollectionKey]; ok {
		generator.Collections = make(map[string]Collection)

		for key, collection := range config[consts.CollectionKey].(map[interface{}]interface{}) {
			currentCollection := *NewCollection(collection)

			if len(currentCollection.Fields) > 0 {
				generator.Collections[key.(string)] = currentCollection
			}
		}
	}

	return
}

// You can also add your own generator function to your own defined tags.
func (g *Generator) Execute(data *consts.FlowDataSet) {
	startTime := time.Now()

	g.Read(data)

	fmt.Printf("Generator operation time: %v\n", time.Since(startTime))
}

func (g *Generator) Read(data *consts.FlowDataSet) {
	for key, collection := range g.Collections {
		bar := progressbar.Default(int64(collection.Limit))

		for i := 0; i < collection.Limit; i++ {
			(*data)[key] = append((*data)[key], g.traverseFields(collection.Fields))
			bar.Add(1)
		}
	}
}

func (g *Generator) Write(data *consts.FlowDataSet) {
	// Not applicaticable in case of generator
}

func (g *Generator) Validate() bool {
	return true
}

func (g *Generator) GetInput() (string, bool) {
	return "", false
}

func (g *Generator) traverseFields(fields consts.GeneratorCollectionFieldType) consts.FlowData {
	data := make(map[string]interface{})

	for k, v := range fields {
		// switch v.(type) {
		// case map[string]interface{}:
		// 	data[k] = g.traverseFields(v.(fieldsType))
		// default:
		data[k] = generateFieldData(v)
		// }
	}

	return data
}

func generateFieldData(fieldType string) interface{} {
	fieldType = strings.ToLower(fieldType)

	// Check if field type is a number
	numberRegex := regexp.MustCompile(`^(number|intiger|int|float|decimal|double)(:[0-9]*)?(:[0-9]*)?`)
	numberMatch := numberRegex.FindStringSubmatch(fieldType)

	if len(numberMatch) > 0 {
		numberType := numberMatch[1]
		numberMin := 0
		numberMax := math.MaxInt

		if 2 < len(numberMatch) {
			if parsedNum, err := strconv.Atoi(strings.ReplaceAll(numberMatch[2], ":", "")); err == nil {
				numberMin = parsedNum
			}
		}

		if 3 < len(numberMatch) {
			if parsedNum, err := strconv.Atoi(strings.ReplaceAll(numberMatch[3], ":", "")); err == nil {
				numberMax = parsedNum
			}
		}

		// fmt.Println(fieldType, numberMatch, numberType, numberMin, numberMax)
		if slices.Contains([]string{"float", "decimal", "double"}, numberType) {
			return (float64(numberMin) + rand.Float64()*(float64(numberMax-numberMin)))
		}

		return (numberMin + rand.Intn(numberMax-numberMin))
	}

	switch fieldType {

	// Address
	case "latitude", "lat":
		return faker.Latitude() // 81.12195
	case "longitude", "lng":
		return faker.Longitude() // -84.38158
	case "address":
		return faker.GetRealAddress() // {2755 Country Drive Fremont CA 94536 {37.557882 -121.986823}}

	// Datetime
	case "unixtime", "unixtimestamp":
		return faker.UnixTime() // 1197930901
	case "date":
		return faker.Date() // 1982-02-27
	case "time", "TimeString":
		return faker.TimeString() // 03:10:25
	case "monthname", "month:name", "month:string":
		return faker.MonthName() // February
	case "year", "yearstring":
		return faker.YearString() // 1994
	case "day", "dayofweek":
		return faker.DayOfWeek() // Sunday
	case "dayofmonth", "dateofmonth":
		return faker.DayOfMonth() // 20
	case "timestamp", "datetime":
		return faker.Timestamp() // 1973-06-21 14:50:46
	case "century":
		return faker.Century() // IV
	case "timezone", "tz":
		return faker.Timezone() // Asia/Jakarta
	case "timeperiod", "ampm":
		return faker.Timeperiod() // PM

	// Internet
	case "email":
		return faker.Email() // mJBJtbv@OSAaT.com
	case "macaddress", "mac":
		return faker.MacAddress() // cd:65:e1:d4:76:c6
	case "domainname", "domain":
		return faker.DomainName() // FWZcaRE.org
	case "url":
		return faker.URL() // https://www.oEuqqAY.org/QgqfOhd
	case "username", "user":
		return faker.Username() // lVxELHS
	case "password":
		return faker.Password() // dfJdyHGuVkHBgnHLQQgpINApynzexnRpgIKBpiIjpTP
	case "ipv4", "ipaddress":
		return faker.IPv4() // 99.23.42.63
	case "ipv6":
		return faker.IPv6() // 975c:fb2c:2133:fbdd:beda:282e:1e0a:ec7d

	// Words and Sentences
	case "word":
		return faker.Word()
	case "sentence":
		return faker.Sentence() // Consequatur perferendis voluptatem accusantium.
	case "paragraph":
		return faker.Paragraph() // Aut consequatur sit perferendis accusantium voluptatem. Accusantium perferendis consequatur voluptatem sit aut. Aut sit accusantium consequatur voluptatem perferendis. Perferendis voluptatem aut accusantium consequatur sit.

	// Payment
	case "cctype", "cardtype", "card:type":
		return faker.CCType() // American Express
	case "ccnumber", "cardnumber", "card:number":
		return faker.CCNumber() // 373641309057568

	case "currency":
		return faker.Currency() // USD
	case "payment:amount", "amoutwithcurrency", "amount":
		return faker.AmountWithCurrency() // USD 49257.100

	// Person
	case "titlemale", "maletitle":
		return faker.TitleMale() // Mr.
	case "titlefemale", "femaletitle":
		return faker.TitleFemale() // Mrs.
	case "firstname", "namefirst":
		return faker.FirstName() // Whitney
	case "firstnamemale", "malefirstname":
		return faker.FirstNameMale() // Kenny
	case "firstnamefemale", "femalefirstname":
		return faker.FirstNameFemale() // Jana
	case "lastname", "namelast":
		return faker.LastName() // Rohan
	case "name", "fullname":
		return faker.Name()

	//Phone
	case "phone", "phonenumber", "mobile", "mobilenumber":
		return faker.Phonenumber() // 201-886-0269
	case "phone:tf", "phone:tollfree", "phone:toll-free":
		return faker.TollFreePhoneNumber() // (777) 831-964572
	case "phone:e164":
		return faker.E164PhoneNumber() // +724891571063

	//  UUID
	case "uuid":
		return faker.UUIDHyphenated() // 8f8e4463-9560-4a38-9b0c-ef24481e4e27
	case "uuid2":
		return faker.UUIDDigit() // 90ea6479fd0e4940af741f0a87596b73

	}
	// faker.ResetUnique() // Forget all generated unique values

	// If not of the cases match
	// then return the same string
	return fieldType

}
