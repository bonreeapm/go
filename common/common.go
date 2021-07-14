package common

//BackendTypes
type BR_BACKEND_TYPE int
const (
	BACKEND_TYPE_UNKNOWN = 0			//unknown类型
	BACKEND_TYPE_CUSTOM  = 300			//web展示拼接方式： Custom-xxxx
	//SQL
	BACKEND_TYPE_MYSQL = 1				//web展示拼接方式： 数据库名-MySQL-host
	BACKEND_TYPE_ORACLE = 2			//web展示拼接方式： 数据库名-Oracle-host
	BACKEND_TYPE_MSSQLSERVER = 3		//web展示拼接方式： 数据库名-SQLServer-host
	BACKEND_TYPE_IBM_DB2 = 4			//web展示拼接方式： 数据库名-DB2-host
	BACKEND_TYPE_POSTGRESQL = 5		//web展示拼接方式： 数据库名-PostgreSQL-host
	BACKEND_TYPE_SHENTONG = 6			//web展示拼接方式： 数据库名-ShenTong-host
	BACKEND_TYPE_SYBASE = 7			//web展示拼接方式： 数据库名-Sybase-host
	BACKEND_TYPE_IBM_INFOMIX = 8		//web展示拼接方式： 数据库名-IBM-Infomix-host
	BACKEND_TYPE_MS_ACCESS = 9			//web展示拼接方式： 数据库名-MSAccess-filename
	BACKEND_TYPE_SQLITE = 10			//web展示拼接方式： 数据库名-SQLite-filename
	BACKEND_TYPE_GENERIC_JDBC = 11     //web展示拼接方式： 数据库名-JDBC-host

	//NOSQL
	BACKEND_TYPE_MONGODB = 101			//web展示拼接方式： Mongo-host:port
	BACKEND_TYPE_MEMCACHE = 102		//web展示拼接方式： Memcached
	BACKEND_TYPE_REDIS = 103			//web展示拼接方式： Redis-host:port
	BACKEND_TYPE_COUCHBASE = 104		//web展示拼接方式： CouchBase-host
	BACKEND_TYPE_COUCHDB = 105			//web展示拼接方式： CouchDB-host:port
	BACKEND_TYPE_CASSANDRA = 106		//web展示拼接方式： CASSANDRA-host:port
	BACKEND_TYPE_SSDB = 107			//web展示拼接方式： SSDB-host:port

	//远程调用
	BACKEND_TYPE_HTTP = 201			//web展示拼接方式： Http-host:port
	BACKEND_TYPE_THRIFT = 202			//web展示拼接方式： Thrift-host:port
	BACKEND_TYPE_WEBSERVICE = 203      //web展示拼接方式： WebService-host:port
	BACKEND_TYPE_WEBSOCKET = 204		//web展示拼接方式： WebSocket-host:port
	BACKEND_TYPE_WCF = 205				//web展示拼接方式： WCF-host:port
	BACKEND_TYPE_WEB_API_1 = 206		//web展示拼接方式： Web-Api1-host:port
	BACKEND_TYPE_WEB_API_2 = 207		//web展示拼接方式： Web-Api2-host:port
	BACKEND_TYPE_DUBBO = 208			//web展示拼接方式： Dubbo-host:por
	BACKEND_TYPE_SUN_RMI = 209			//web展示拼接方式： RMI-host:port
	BACKEND_TYPE_SOCKET = 210			//web展示拼接方式： Sock-host:port
	BACKEND_TYPE_HSF = 211				//web展示拼接方式： HSF-host:port
	BACKEND_TYPE_JMS = 212				//web展示拼接方式： JMS-host:port
	BACKEND_TYPE_SPRING_JMS = 213      //web展示拼接方式： Spring-JMS-host:port
	BACKEND_TYPE_RABBITMQ = 214		//web展示拼接方式： RabbitMQ-host:port
	BACKEND_TYPE_CXF = 215				//web展示拼接方式： CXF-host:port
	BACKEND_TYPE_WTC = 216				//web展示拼接方式： WTC-host:port
	BACKEND_TYPE_GRPC = 217			//web展示拼接方式： GRPC-host:port
	BACKEND_TYPE_AXIS = 218			//web展示拼接方式： AXIS-host:port
	BACKEND_TYPE_SOAP = 219			//web展示拼接方式： SOAP-host:port
	BACKEND_TYPE_SOLR = 220			//web展示拼接方式： SOLR-host:port
	BACKEND_TYPE_MSMQ = 221			//web展示拼接方式： MSMQ-host:port
	BACKEND_TYPE_CIB_RPC = 222			//web展示拼接方式： CIB-host:port  (兴业证券rpc)
	BACKEND_TYPE_MOTAN = 223			//web展示拼接方式： MOTAN-host:port  (新浪微博rpc)
	BACKEND_TYPE_IBMMQ = 224			//web展示拼接方式： IBMMQ-host:port
	BACKEND_TYPE_ROCKETMQ = 225		//web展示拼接方式： RocketMQ-host:port
	BACKEND_TYPE_ACTIVEMQ = 226		//web展示拼接方式： ActiveMQ-host:port
	BACKEND_TYPE_HBASE = 227			//web展示拼接方式： Hbase-host:port
	BACKEND_TYPE_KAFKA = 228			//web展示拼接方式： Kafka-host:port
	BACKEND_TYPE_ES = 229				//web展示拼接方式： ES-host:port ( elasticsearch  )
)

type BR_ERROR_TYPE int
const (
	BR_ERROR_TYPE_HTTP = iota
	BR_ERROR_TYPE_EXCEPTION
)

// CrossHeader is represents the header value for cross request.
type CrossHeader string
const (
	// CrossRequestHeader is represents the request header value for cross request.
	CrossRequestHeader = "BrCsReq"
	// CrossResponseHeader is represents the response header value for cross request.
	CrossResponseHeader = "BrCsRes"
)