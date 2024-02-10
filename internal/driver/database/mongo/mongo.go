package mongo

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/webhkp/godft/internal/consts"
	"github.com/webhkp/godft/internal/driver/database"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Mongo struct {
	*database.Database
	intputTaskName   string
	ctx              context.Context
	selectedDatabase *mongo.Database
}

func NewMongo(config consts.FieldsType) (mongo *Mongo) {
	mongo = &Mongo{
		Database: database.NewDatabase(config),
	}

	if mongo.Connection == "" {
		panic("Connection is required for mongo driver")
	}

	if _, ok := config[consts.InputDriverKey]; ok {
		mongo.intputTaskName = config[consts.InputDriverKey].(string)
	}

	return
}

func (d *Mongo) Execute(data *consts.FlowDataSet) {
	startTime := time.Now()

	d.connect()
	d.Read(data)
	d.Write(data)

	fmt.Printf("Mongo operation time: %v\n", time.Since(startTime))
}

func (d *Mongo) connect() {
	mongoClient, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(d.Connection))

	if err != nil {
		fmt.Println("MongoDB connection error: ", err)
		return
	}

	d.ctx, _ = context.WithTimeout(context.Background(), 60*time.Second)
	d.selectedDatabase = mongoClient.Database(d.DatabaseName)
}

func (d *Mongo) Read(data *consts.FlowDataSet) {
	if !d.AllCollection && len(d.Collections) == 0 {
		return
	}

	if d.AllCollection {
		d.Collections = d.getAllCollections()
	}

	for collectionName, collectionConfig := range d.Collections {
		(*data)[collectionName] = (d.getCollectionData(collectionName, collectionConfig))
	}
}

func (d *Mongo) getCollectionData(collectionName string, collectionConfig database.Collection) []map[string]interface{} {
	collectionObj := d.selectedDatabase.Collection(collectionName)

	findOptions := options.Find()

	if collectionConfig.Limit != -1 {
		findOptions.SetLimit(int64(collectionConfig.Limit))
	}

	if len(collectionConfig.Sort) > 0 {
		sortValues := bson.D{}

		for key, sort := range collectionConfig.Sort {
			sortValues = append(sortValues, primitive.E{key, sort})
		}

		findOptions.SetSort(sortValues)
	}

	if len(collectionConfig.Fields) > 0 {
		projectValues := bson.D{}

		for _, projection := range collectionConfig.Fields {
			projectValues = append(projectValues, primitive.E{projection, true})
		}

		findOptions.SetProjection(projectValues)
	}

	var findResult []map[string]interface{}
	cursor, err := collectionObj.Find(d.ctx, bson.M{}, findOptions)

	if err != nil {
		fmt.Println("Error while getting result from collection: " + collectionName)
		return findResult
	}

	defer cursor.Close(d.ctx)
	cursor.All(d.ctx, &findResult)

	return findResult

}

func (d *Mongo) getAllCollections() (collections map[string]database.Collection) {
	collections = make(map[string]database.Collection)

	result, err := d.selectedDatabase.ListCollectionNames(
		d.ctx,
		bson.D{},
	)

	if err != nil {
		return
	}

	for _, val := range result {
		collections[val] = *database.NewDefaultCollection()
	}

	return
}

func (d *Mongo) Write(data *consts.FlowDataSet) {
	if d.ReadOnly || len(*data) == 0 {
		return
	}

	for collectionName, collectionData := range *data {
		collection := d.selectedDatabase.Collection(collectionName)
		data, _ := json.Marshal(collectionData)
		var jsonData []interface{}
		json.Unmarshal(data, &jsonData)

		_, err := collection.InsertMany(d.ctx, jsonData)

		if err != nil {
			fmt.Println("Data insert error: ", err)
		}
	}
}

func (d *Mongo) Validate() bool {
	return true
}

func (d *Mongo) GetInput() (string, bool) {
	if d.intputTaskName != "" {
		return d.intputTaskName, true
	}

	return "", false
}
