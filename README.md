# GoDFT - Go Data Flow Tunnel 

GoDFT is a Data Pipeline Multi-Tool, combining seamless bulk data generation, intelligent export-import capabilities,
and precision data mapping. 

Elevate your data management with a tool that offers flexibility, automation, and reliability, ensuring your data workflows are efficient and powerful.

**Full Documentation**: <a href="https://webhkp.com/godft">https://webhkp.com/godft</a>

# Installation
Run the following command to install GoDFT (make sure you have golang installed)-

```shell
go install github.com/webhkp/godft@latest
```

Or, you can install the binaries from the release page-
<a href="https://github.com/webhkp/godft/releases">https://github.com/webhkp/godft/releases</a>

# Usage

```shell
godft [command]
```


# Available Commands

```
explain     Explain task flow
help        Help about any command
run         Run all task flows
validate    Validate task flow configuration
```

# Drivers

The following drivers provided read and/or write functionality-

<table>
    <tr>
        <th>Driver</th>
        <th>Description</th>
        <th>Read</th>
        <th>Write</th>
    </tr>
    <tr>
        <td>MongoDB</td>
        <td>Perform operations on MongoDB</td>
        <td>Yes</td>
        <td>Yes</td>
    </tr>
    <tr>
        <td>MySQL</td>
        <td>Perform operations on MongoDB</td>
        <td>Yes</td>
        <td>No</td>
    </tr>
    <tr>
        <td>JSON</td>
        <td>Driver for JSON operations</td>
        <td>Yes</td>
        <td>Yes</td>
    </tr>
    <tr>
        <td>Generator</td>
        <td>Fake data generator</td>
        <td>Yes</td>
        <td>No</td>
    </tr>
</table>

# Example
Here are a few examples of how we can configure and perform options using GoDFT-

## Example #1: MongoDB to JSON

### Configuration
Create a configuration file with ".yml" extention. We are naming the file **mongo_to_json_test.yml**

```yaml
# Source service
my-mongo-input:
  driver: mongo
  connection: mongodb://webhkp:secretpass@localhost:27017/?maxPoolSize=20&w=majority
  database: bigboxcode
  collection:
    customer:
      limit: 4
      sort:
        age: -1
      field:
        - address.city
        - age
        - name
    product:
      limit: 10
      sort:
        id: -1
    order:
      field:
        - price
        - customer
        - quantity
      sort:
        price: 1
        quantity: -10
      limit: 22

# Destination service
json-child:
  driver: json
  input: my-mongo-input # Define the source, from which this service will get data
  outputPath: ./docs/testres/
  readOnly: false # Make this instance writable

```

### Run
Run the following command-

```shell
godft run mongo_to_json_test.yml
```

## Example #2: MongoDB to MongoDB

### Configuration
Create a configuration file with ".yml" extention. We are naming the file **mongo_to_mongo_test.yml**

```yaml
# Source MongoDB service
my-mongo-input:
  driver: mongo
  connection: mongodb://webhkp:secretpass@localhost:27017/?maxPoolSize=20&w=majority
  database: bigboxcode
  collection: all # copy all collections from this mongo instance


# Destination MongoDB service
my-mongo-output:
  driver: mongo
  connection: mongodb://webhkp:secretpass@localhost:27017/?maxPoolSize=20&w=majority
  database: godft
  readOnly: false # make this instance
  input: my-mongo-input # Define the service from which this one gets data

```

### Run
Run the following command-

```shell
godft run mongo_to_mongo_test.yml
```

## Example #3: JSON to MongoDB
Copy data from JSON file(s) to MongoDB instace.

### Configuration
Create a configuration file with ".yml" extention. We are naming the file **json_to_mongo_test.yml**

```yaml
json-parent:
  driver: json
  inputPath:
    - ./docs/testres/ # All json file from this directory
    - ./docs/customer_sample.json

my-mongo:
  driver: mongo
  connection: mongodb://webhkp:secretpass@localhost:27017/?maxPoolSize=20&w=majority
  database: godft
  readOnly: false
  input: json-parent # Receive data from "json_parent" service

```

### Run
Run the following command-

```shell
godft run json_to_mongo_test.yml
```

## Example #4: MySQL to JSON

### Configuration
Create a configuration file with ".yml" extention. We are naming the file **mysql_to_json_test.yml**

```yaml
mysql-source:
  driver: mysql
  host: 127.0.0.1
  port: 3306
  user: root
  password: root
  database: issue_tracker
  collection: all

json-child:
  driver: json
  input: mysql-source
  outputPath: ./docs/testres/
  readOnly: false
```

### Run
```shell
godft run mysql_to_json_test.yml
```

## Example #5: Generator to MongoDB

### Configuration
Create a configuration file with ".yml" extention. We are naming the file **generator_to_mongo_test.yml**

```yaml
# Source service
my-generator:
  driver: generator
  collection:
    customer:
      field:
        full_name: name
        age: number:10:90 # Minimum value 10, maximum value 90
        address: address
        customerid: number
        phone: phone
        number1: number:100:200
        number2: number:100
        number3: number::200
        decimal1: decimal:20 # Minimum 20
        id: number::400 # Max vlaue 400
        sku: uuid
        price: decimal::99
    product:
      limit: 3
      
    ordersample:
      limit: 50

      field:
        price: number
        product: uuid
        quantity: number:1:25


# Destination service
my-mongo-output:
  driver: mongo
  connection: mongodb://webhkp:secretpass@localhost:27017/?maxPoolSize=20&w=majority
  database: godft
  readOnly: false # Make this service writable
  input: my-generator # define the input service
```

### Run

```shell
godft run generator_to_mongo_test.yml
```

## Example #6: Generator to JSON

### Configuration
Create a configuration file with ".yml" extention. We are naming the file **generator_to_json_test.yml**


```yaml
# Source service
my-generator:
  driver: generator
  collection:
    customer:
      field:
        full_name: name
        age: number:10:90
        address: address
        customerid: number
        phone: phone
        number1: number:100:200
        number2: number:100
        number3: number::200
        decimal1: decimal:20
        id: number::400
        sku: uuid
        price: decimal::99
    product:
      limit: 3
      
    order:
      limit: 5

      field:
        price: number
        product: uuid
        quantity: number:1:25

# Destination service
fake-data-json-child:
  driver: json
  input: my-generator
  readOnly: false
```


### Run

```shell
godft run generator_to_json_test.yml
```

# Run/Build from Source Code

- Clone soruce code using- ``` git clone https://github.com/webhkp/godft ```
- Run project using- ``` go run main.go run your-yml-file-here.yml ```
- Build using - ``` go build main.go  ``` 


# Detail Documentation

<table>
    <tr>
        <td>Full Documentation</td>
        <td>
            <a href="https://webhkp.com/godft">
              https://webhkp.com/godft
            </a>
        </td>
    </tr>
    <tr>
        <td>Configuration Details</td>
        <td>
            <a href="https://webhkp.com/godft/godft-configuration">
              https://webhkp.com/godft/godft-configuration
            </a>
        </td>
    </tr>
    <tr>
        <td>Internals</td>
        <td>
            <a href="https://webhkp.com/godft/godft-internals">
              https://webhkp.com/godft/godft-internals
            </a>
        </td>
    </tr>
    <tr>
        <td>Source Code Explanation</td>
        <td>
            <a href="https://webhkp.com/godft/godft-source-code">
              https://webhkp.com/godft/godft-source-code
            </a>
        </td>
    </tr>
</table>
