package mysql

import (
	"database/sql"
	"fmt"
	"strconv"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"

	"github.com/webhkp/godft/internal/consts"
	"github.com/webhkp/godft/internal/driver/database"
)

type Mysql struct {
	*database.Database
	intputTaskName   string
	selectedDatabase *sql.DB
	host             string
	port             string
	protocol         string
	username         string
	password         string
}

func NewMysql(config consts.FieldsType) (mysql *Mysql) {
	mysql = &Mysql{
		Database: database.NewDatabase(config),
		host:     consts.DefaultMysqlHost,
		port:     strconv.Itoa(consts.DefaultMysqlPort),
		protocol: consts.DefaultMysqlProtocol,
		username: consts.DefaultMysqlUser,
		password: consts.DefaultMysqlPassword,
	}

	if mysql.Connection == "" {
		if _, ok := config[consts.InputDriverKey]; ok {
			mysql.intputTaskName = config[consts.InputDriverKey].(string)
		}

		if _, ok := config[consts.HostKey]; ok {
			mysql.host = config[consts.HostKey].(string)
		}

		if _, ok := config[consts.PortKey]; ok {
			mysql.port = config[consts.PortKey].(string)
		}

		if _, ok := config[consts.ProtocolKey]; ok {
			mysql.protocol = config[consts.ProtocolKey].(string)
		}

		if _, ok := config[consts.UserNameKey]; ok {
			mysql.username = config[consts.UserNameKey].(string)
		}

		if _, ok := config[consts.PasswordKey]; ok {
			mysql.password = config[consts.PasswordKey].(string)
		}

		// Construct connection from provided value
		mysql.Connection = fmt.Sprintf("%s:%s@%s(%s:%s)/%s", mysql.username, mysql.password, mysql.protocol, mysql.host, mysql.port, mysql.DatabaseName)
	}

	return
}

func (d *Mysql) Execute(data *consts.FlowDataSet) {
	startTime := time.Now()

	d.connect()
	d.Read(data)
	d.Write(data)

	fmt.Printf("Mysql operation time: %v\n", time.Since(startTime))
}

func (d *Mysql) connect() {
	selectedDatabase, err := sql.Open("mysql", d.Connection)

	if err != nil {
		fmt.Println("MySQL connection error: ", err)
		return
	}

	d.selectedDatabase = selectedDatabase
}

func (d *Mysql) Read(data *consts.FlowDataSet) {
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

func (d *Mysql) getCollectionData(collectionName string, collectionConfig database.Collection) []map[string]interface{} {
	collectionName = strings.ReplaceAll(collectionName, ";", "")
	rows, err := d.selectedDatabase.Query("SELECT * FROM " + collectionName)
	if err != nil {
		fmt.Println("Error processing MySQL query: ", err)
		return nil
	}
	defer rows.Close()

	return processRow(rows)
}

func processRow(list *sql.Rows) (rows []map[string]interface{}) {
	fields, _ := list.Columns()
	for list.Next() {
		scans := make([]interface{}, len(fields))
		row := make(map[string]interface{})

		for i := range scans {
			scans[i] = &scans[i]
		}
		list.Scan(scans...)
		for i, v := range scans {
			var value = ""
			if v != nil {
				value = fmt.Sprintf("%s", v)
			}
			row[fields[i]] = value
		}
		rows = append(rows, row)
	}
	return
}

func (d *Mysql) getAllCollections() (collections map[string]database.Collection) {
	collections = make(map[string]database.Collection)
	tableQueryRes, err := d.selectedDatabase.Query("SHOW TABLES")

	if err != nil {
		fmt.Println("In getall collection", err)
		return
	}

	var tableName string

	for tableQueryRes.Next() {
		tableQueryRes.Scan(&tableName)
		collections[tableName] = *database.NewDefaultCollection()
	}

	return
}

func (d *Mysql) Write(data *consts.FlowDataSet) {
	// @todo
}

func (d *Mysql) Validate() bool {
	return true
}

func (d *Mysql) GetInput() (string, bool) {
	if d.intputTaskName != "" {
		return d.intputTaskName, true
	}

	return "", false
}
