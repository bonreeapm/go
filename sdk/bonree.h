/*
 * Copyright (c) Bonree, Inc., and its affiliates
 * 2017
 * All Rights Reserved
 */

#ifndef INCLUDE_BONREE_H_
#define INCLUDE_BONREE_H_

#ifdef __cplusplus
extern "C" {
#endif

#ifdef _WIN32
#define BR_API __declspec(dllexport)
#else
#define BR_API __attribute__((visibility("default")))
#endif

/**
 *  跨容器调用请求头 头部 key 名称
 */
#define BR_CS_REQ_HEADER_NAME "br-cs-req"

/**
 * 跨容器调用响应头 头部 key 名称
 */
#define BR_CS_RSP_HEADER_NAME "br-cs-rsp"

typedef const void* br_app_t;
typedef const void* br_bt_t;
typedef const void* br_backend_t;
typedef const void* br_snapshot_thread_t;
typedef const void* br_snapshot_func_t;
typedef const void* br_exitcall_t;
typedef const void* br_error_t;

typedef struct br_app_config_s
{
    const char* app_name;
    const char* tier_name;
    const char* cluster_name;
    const char* agent_name;
} br_app_config_t;

typedef enum
{
    BR_SQL_TYPE_MYSQL,
    BR_SQL_TYPE_ORACLE,
    BR_SQL_TYPE_MSSQLSERVER,
    BR_SQL_TYPE_IBM_DB2,
    BR_SQL_TYPE_POSTGRESQL,
    BR_SQL_TYPE_SHENTONG,
    BR_SQL_TYPE_SYBASE,
    BR_SQL_TYPE_IBM_INFOMIX,
    BR_SQL_TYPE_MS_ACCESS,
} br_sql_type;

typedef enum
{
    BR_NOSQL_TYPE_MONGODB,
    BR_NOSQL_TYPE_MEMCACHE,
    BR_NOSQL_TYPE_REDIS,
    BR_NOSQL_TYPE_COUCHBASE,
    BR_NOSQL_TYPE_COUCHDB,
    BR_NOSQL_TYPE_CASSANDRA,
} br_nosql_type;

typedef enum
{
    BR_RPC_TYPE_HTTP,               //web展示拼接方式： Http-host:port
    BR_RPC_TYPE_THRIFT,             //web展示拼接方式： Thrift-host:port
    BR_RPC_TYPE_WEBSERVICE,         //web展示拼接方式： WebService-host:port
    BR_RPC_TYPE_WEBSOCKET,          //web展示拼接方式： WebSocket-host:port
    BR_RPC_TYPE_WCF,                //web展示拼接方式： WCF-host:port
    BR_RPC_TYPE_WEB_API_1,          //web展示拼接方式： Web-Api1-host:port
    BR_RPC_TYPE_WEB_API_2,          //web展示拼接方式： Web-Api2-host:port
    BR_RPC_TYPE_DUBBO,              //web展示拼接方式： Dubbo-host:port
    BR_RPC_TYPE_SUN_RMI,            //web展示拼接方式： RMI-host:port
    BR_RPC_TYPE_SOCKET,             //web展示拼接方式： Sock-host:port
    BR_RPC_TYPE_HSF,                //web展示拼接方式： HSF-host:port
    BR_RPC_TYPE_JMS,                //web展示拼接方式： JMS-host:port
    BR_RPC_TYPE_SPRING_JMS,         //web展示拼接方式： Spring-JMS-host:port
    BR_RPC_TYPE_RABBITMQ,           //web展示拼接方式： RabbitMQ-host:port
    BR_RPC_TYPE_CXF,                //web展示拼接方式： CXF-host:port
    BR_RPC_TYPE_WTC,                //web展示拼接方式： WTC-host:port
    BR_RPC_TYPE_GRPC,               //web展示拼接方式： GRPC-host:port
    BR_RPC_TYPE_AXIS,               //web展示拼接方式： AXIS-host:port
    BR_RPC_TYPE_SOAP,               //web展示拼接方式： SOAP-host:port
    BR_RPC_TYPE_SOLR,               //web展示拼接方式： SOLR-host:port
    BR_RPC_TYPE_MSMQ,               //web展示拼接方式： MSMQ-host:port
} br_rpc_type;

typedef enum
{
    BR_ERROR_TYPE_HTTP,
    BR_ERROR_TYPE_EXCEPTION,
} br_error_type;

typedef unsigned char br_bool;

/**
 * sdk 初始化
 * @return 0:失败   1:成功     失败原因可查看 log 记录
 */
BR_API int br_sdk_init();

/**
 * sdk 销毁
 */
BR_API void br_sdk_release();

/**
 * app 初始化
 * @return app_handle
 */
BR_API br_app_t br_app_init();

/**
 *  app 自定义初始化
 * @param app_cfg
 * @return app_handle
 */
BR_API br_app_t br_app_init_with_cfg(br_app_config_t *app_cfg);

/**
 * app 销毁
 * @param app_handle
 */
BR_API void br_app_release(br_app_t app_handle);

/**
 * 业务开始
 * @param app_handle
 * @param name 业务名称
 * @return bt_handle
 */
BR_API br_bt_t br_bt_begin(br_app_t app_handle, const char* name);

/**
 * 业务开始
 * @param app_handle
 * @param name 业务名称
 * @param cross_request_header  跨容器请求头
 * @return bt_handle
 */
BR_API br_bt_t br_bt_begin_ex(br_app_t app_handle, const char* name, const char* cross_request_header);

/**
 * 构造跨容器响应头部
 * @return 跨容器响应头部
 */
BR_API const char* br_bt_generate_cross_resheader(br_bt_t bt_handle);

/**
 * 业务结束
 * @param bt_handle
 */
BR_API void br_bt_end(br_bt_t bt_handle);

/**
 * 设置业务 url
 * @param bt_handle
 * @param url
 */
BR_API void br_bt_set_url(br_bt_t bt_handle, const char* url);

/**
 * 错误添加
 * @param bt_handle
 * @param error_type
 * @param error_name
 * @param summary
 * @param details
 * @param mark_bt_as_error
 */
BR_API void br_bt_add_error(br_bt_t bt_handle, br_error_type error_type, const char* error_name, const char* summary, const char* details, int mark_bt_as_error);

/**
 * 定义 sql 后端
 * @param type
 * @param host
 * @param port
 * @param dbschema
 * @param vendor
 * @param version
 * @return backend_handle
 */
BR_API br_backend_t br_backend_declare_sql(br_sql_type type, const char* host, int port, const char* dbschema, const char* vendor, const char* version);

/**
 * 定义 nosql 后端
 * @param type
 * @param server_pool
 * @param port
 * @param vendor
 * @return backend_handle
 */
BR_API br_backend_t br_backend_declare_nosql(br_nosql_type type, const char* server_pool, int port, const char* vendor);

/**
 * 定义 rpc 后端
 * @param type
 * @param host
 * @param port
 * @return backend_handle
 */
BR_API br_backend_t br_backend_declare_rpc(br_rpc_type type, const char* host, int port);

/**
 * 后端调用开始
 * @param bt_handle
 * @param backend
 * @return
 */
BR_API br_exitcall_t br_exitcall_begin(br_bt_t bt_handle, br_backend_t backend);


/**
 * 设置后端调用详情
 * @param exitcall
 * @param cmd
 * @param details
 * @return 0:失败   1:成功  失败原因可查看 log 记录
 */
BR_API int br_exitcall_set_detail(br_exitcall_t exitcall, const char* cmd, const char* details);

/**
 * 后端调用错误添加
 * @param exitcall
 * @param error_type
 * @param error_name
 * @param summary
 * @param details
 * @param mark_as_error
 */
BR_API void br_exitcall_add_error(br_exitcall_t exitcall, br_error_type error_type, const char* error_name, const char* summary, const char* details, int mark_as_error);

/**
 * 后端调用结束
 * @param exitcall
 */
BR_API void br_exitcall_end(br_exitcall_t exitcall);

/**
 * 构造跨容器调用请求头部
 * @param exitcall
 * @return
 */
BR_API const char* br_exitcall_generate_cross_reqheader(br_exitcall_t exitcall);

/**
 * 设置跨容器响应头部
 * @param exitcall
 * @param cross_response_header
 */
BR_API void br_exitcall_set_cross_resheader(br_exitcall_t exitcall, const char* cross_response_header);

/**
 * 业务是否正在采集快照
 * @param bt_handle
 * @return 0:没有   1:有
 */
BR_API br_bool br_bt_is_snapshotting(br_bt_t bt_handle);

/**
 * 添加业务快照数据
 *
 * 该方法是添加数据到快照，当且仅当业务在采集快照的时候。
 * 业务没有采集，调用该方法会立即返回。
 * 如果构造参数数据比较耗时，建议先调用 br_bt_is_snapshotting 做一下判断。
 *
 * @param bt_handle
 * @param key
 * @param value
 */
BR_API void br_bt_snapshot_data(br_bt_t bt_handle, const char* key, const char* value);

/**
 * 快照线程开始
 * @param bt_handle
 * @param thread_name
 * @return 快照线程 handle
 */
BR_API br_snapshot_thread_t br_bt_snapshot_thread_start(br_bt_t bt_handle);

/**
 * 快照线程结束
 * @param thread_handle
 */
BR_API void br_bt_snapshot_thread_end(br_snapshot_thread_t thread_handle);

/**
 * 快照方法开始
 * @param thread_handle  线程 handle
 * @param func_name      方法名称
 * @param file_name      代码文件
 * @param lineno         代码行号
 * @return  快照方法 handle
 */
BR_API br_snapshot_func_t br_bt_snapshot_func_start(br_snapshot_thread_t thread_handle, const char* class_name, const char* func_name, const char* file_name, int lineno);

/**
 * 快照方法结束
 * @param br_snapshot_func_t
 */
BR_API void br_bt_snapshot_func_end(br_snapshot_func_t func_handle);

/**
 * 快照后端调用添加
 * @param br_snapshot_func_t
 * @param exitcall
 */
BR_API void br_snapshot_exitcall_add(br_snapshot_func_t func_handle, br_exitcall_t exitcall);

/**
 * 快照错误添加
 * @param br_snapshot_thread_t
 * @param error_name
 * @param summary
 * @param details
 */
BR_API void br_snapshot_error_add(br_snapshot_thread_t thread_handle, const char* error_name, const char* summary, const char* details);

#ifdef __cplusplus
} /* extern "C" */
#endif

#endif /* INCLUDE_BONREE_H_ */
