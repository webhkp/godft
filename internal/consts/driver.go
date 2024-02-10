package consts

const (
	// Driver names
	MongoDriverName = "mongo"
	MysqlDriverName = "mysql"
	GeneratorDriverName = "generator"
	JsonDriverName = "json"

	// Common key identifiers in drivers
	InputDriverKey = "input"
	ConnectionKey = "connection"
	OutputPathKey = "outputPath"
	InputPathKey = "inputPath"

	ReadOnlyKey = "readOnly"
	DatabaseKey = "database"
	CollectionKey = "collection"

	LimitKey = "limit"
	SortKey = "sort"
	FieldKey = "field"

	AllCollectionValue = "all"

	NoLimit = -1 // means no limit

	// Database config keys
	HostKey = "host"
	PortKey = "port"
	ProtocolKey = "protocol"
	UserNameKey = "user"
	PasswordKey = "password"

	// MySQL Defaults
	DefaultMysqlHost = "127.0.0.1"
	DefaultMysqlPort = 3306
	DefaultMysqlUser = "root"
	DefaultMysqlPassword = ""
	DefaultMysqlProtocol = "tcp"

)