package common

type BR_SQL_TYPE int
const (
	BR_SQL_TYPE_MYSQL = iota
	BR_SQL_TYPE_ORACLE
	BR_SQL_TYPE_MSSQLSERVER
	BR_SQL_TYPE_IBM_DB2
	BR_SQL_TYPE_POSTGRESQL
	BR_SQL_TYPE_SHENTONG
	BR_SQL_TYPE_SYBASE
	BR_SQL_TYPE_IBM_INFOMIX
	BR_SQL_TYPE_MS_ACCESS
)

type BR_NOSQL_TYPE int
const (
	BR_NOSQL_TYPE_MONGODB = iota
	BR_NOSQL_TYPE_MEMCACHE
	BR_NOSQL_TYPE_REDIS
	BR_NOSQL_TYPE_COUCHBASE
	BR_NOSQL_TYPE_COUCHDB
	BR_NOSQL_TYPE_CASSANDRA
)

type BR_RPC_TYPE int
const (
	BR_RPC_TYPE_HTTP = iota      			//web展示拼接方式： Http-host:port
	BR_RPC_TYPE_THRIFT      		//web展示拼接方式： Thrift-host:port
	BR_RPC_TYPE_WEBSERVICE      	//web展示拼接方式： WebService-host:port
	BR_RPC_TYPE_WEBSOCKET       	//web展示拼接方式： WebSocket-host:port
	BR_RPC_TYPE_WCF      			//web展示拼接方式： WCF-host:port
	BR_RPC_TYPE_WEB_API_1       	//web展示拼接方式： Web-Api1-host:port
	BR_RPC_TYPE_WEB_API_2       	//web展示拼接方式： Web-Api2-host:port
	BR_RPC_TYPE_DUBBO     			//web展示拼接方式： Dubbo-host:port
	BR_RPC_TYPE_SUN_RMI      		//web展示拼接方式： RMI-host:port
	BR_RPC_TYPE_SOCKET      		//web展示拼接方式： Sock-host:port
	BR_RPC_TYPE_HSF      			//web展示拼接方式： HSF-host:port
	BR_RPC_TYPE_JMS     			//web展示拼接方式： JMS-host:port
	BR_RPC_TYPE_SPRING_JMS      	//web展示拼接方式： Spring-JMS-host:port
	BR_RPC_TYPE_RABBITMQ      		//web展示拼接方式： RabbitMQ-host:port
	BR_RPC_TYPE_CXF      			//web展示拼接方式： CXF-host:port
	BR_RPC_TYPE_WTC      			//web展示拼接方式： WTC-host:port
	BR_RPC_TYPE_GRPC      			//web展示拼接方式： GRPC-host:port
	BR_RPC_TYPE_AXIS      			//web展示拼接方式： AXIS-host:port
	BR_RPC_TYPE_SOAP      			//web展示拼接方式： SOAP-host:port
	BR_RPC_TYPE_SOLR      			//web展示拼接方式： SOLR-host:port
	BR_RPC_TYPE_MSMQ      			//web展示拼接方式： MSMQ-host:port
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