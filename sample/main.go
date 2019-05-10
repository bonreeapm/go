package main

import(
	"github.com/bonreeapm/go"
	"github.com/bonreeapm/go/common"
	"github.com/bonreeapm/go/routineEngine"
	"time"
	"net/http"
	"strconv"
	"fmt"	
	"log"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/garyburd/redigo/redis"
)

func routine(w http.ResponseWriter, r *http.Request) {	
	btn := bonree.GetRoutineTransaction()

	if btn == nil {
		fmt.Fprint(w, "Get Transaction fail")
		return
	}

	snapshotFunc := btn.SnapshotFuncStart("main", "Routine")

	defer btn.SnapshotFuncEnd(snapshotFunc)
	
	time.Sleep(time.Duration(3)*time.Second)
}

func setURL(w http.ResponseWriter, r *http.Request) {	
	btn := bonree.GetCurrentTransaction(w)

	if btn == nil {
		fmt.Fprint(w, "Get Transaction fail")
		return
	}
	
	time.Sleep(time.Duration(3)*time.Second)
}

func addError(w http.ResponseWriter, r *http.Request) {
	btn := bonree.GetCurrentTransaction(w)

	if btn == nil {
		fmt.Fprint(w, "Get Transaction fail")
		return
	}

	btn.AddError("UnkownException","UnkownException: something wrong","UnkownException: something wrong in File xxx.cpp:12", true)

	time.Sleep(time.Duration(3)*time.Second)
}

func receiveCrossRequest(w http.ResponseWriter, r *http.Request) {	
	time.Sleep(time.Duration(3)*time.Second)
}

func sendCrossRequest(w http.ResponseWriter, r *http.Request) {
	host := "127.0.0.1"
	port := 9099

	// bonree:获取businessTransaction(业务)
	btn := bonree.GetCurrentTransaction(w)

	if btn == nil {
		fmt.Fprint(w, "Get Transaction fail")
		return
	} 

	// bonree:创建快照方法
	snapshotFunc := btn.SnapshotFuncStart("main", "sendCrossRequest")

	// bonree:创建exitCall(后端)
	exitcall := btn.StartRPCExitCall(common.BR_RPC_TYPE_HTTP, host, port)

	snapshotFunc.AddExitCall(exitcall)

	// bonree:结束后端
	defer exitcall.End()
	// bonree:结束快照方法
	defer btn.SnapshotFuncEnd(snapshotFunc)

	client := &http.Client{}
	client.Transport = exitcall.RoundTripper()
	resp, err := client.Get("http://" + host + ":" + strconv.Itoa(port) + "/receiveCrossRequest")

	if err != nil {
		fmt.Fprint(w, err.Error())

		// bonree:设置后端错误
		exitcall.AddError("Http Get Error", err.Error(), err.Error(), true)
		return
	}

	// bonree:设置后端跨容器响应头
	exitcall.SetCrossResheader(resp.Header)

	fmt.Fprint(w, "Send Cross Request.")
	return
}

const mysqldb = "root:111111@tcp(192.168.0.201:3306)/test"
func mysqlSelectHandler(w http.ResponseWriter, r *http.Request) {
	// bonree:获取businessTransaction(业务)
	btn := bonree.GetCurrentTransaction(w)

	if btn == nil {
		fmt.Fprint(w, "Get Transaction fail")
		return
	}

	db, err := sql.Open("mysql", mysqldb)

	if err != nil {
		return
	}

	defer db.Close();

	// bonree:创建快照方法
	snapshotFunc := btn.SnapshotFuncStart("main", "mysqlSelectHandler")

	// bonree:创建exitCall(后端)
	exitcall := btn.StartSQLExitCall(common.BR_SQL_TYPE_MYSQL, "192.168.0.201", 3306, "test", "mysql", "")

	snapshotFunc.AddExitCall(exitcall)

	// bonree:结束后端
	defer exitcall.End()

	// bonree:结束快照方法
	defer btn.SnapshotFuncEnd(snapshotFunc)

	// bonree:设置后端详情
	exitcall.SetDetail("select * from go_test where id = ?", "select * from go_test where id = ?")
	stmt, err := db.Prepare("select * from go_test where id = ?")

	var id, prop3, prop4 int
	var name, name2 string

	row := stmt.QueryRow(1)

	if row != nil {
		err = row.Scan(&id, &name, &name2, &prop3, &prop4)

		if err != nil {
			// bonree:设置后端错误
			exitcall.AddError("MySQL Error", err.Error(), err.Error(), true)
		}
	}

	fmt.Fprint(w, "mysqlSelectHandler.")
	return
}

func redisGetHandler(w http.ResponseWriter, r *http.Request) {
	btn := bonree.GetCurrentTransaction(w)

	if btn == nil {
		fmt.Fprint(w, "Get Transaction fail")
		return
	}

	c, err := redis.Dial("tcp", "192.168.0.201:6379")
	if err != nil {
		return
	}

	snapshotFunc := btn.SnapshotFuncStart("main", "redisGetHandler")

	exitcall := btn.StartNoSQLExitCall(common.BR_NOSQL_TYPE_REDIS, "192.168.0.201", 6379, "redis")

	snapshotFunc.AddExitCall(exitcall)
	
	defer exitcall.End()
	defer btn.SnapshotFuncEnd(snapshotFunc)

	exitcall.SetDetail("GET", "GET test_value")

	_, err = c.Do("GET", "test_value")

	if err != nil {
		exitcall.AddError("Redis Error", err.Error(), err.Error(), true)
	}

	fmt.Fprint(w, "redisGetHandler.")
	return
}

func main() {
	app, err := bonree.NewApplication("sample")
	if err != nil {
		log.Fatalln(err)
		return
	}

	// 如果需要RoutineEngine支持，则加上此行代码
	bonree.RoutineEngineInit(routineEngine.Get())

	api, err := bonree.NewAPI(app)

	if err != nil {
		log.Fatalln(err)
		return
	}

	defer api.Stop()

	http.HandleFunc(api.WrapHandleFunc("/routine", routine))
	http.HandleFunc(api.WrapHandleFunc("/setURL", setURL))
	http.HandleFunc(api.WrapHandleFunc("/addError", addError))
	http.HandleFunc(api.WrapHandleFunc("/sendCrossRequest", sendCrossRequest))
	http.HandleFunc(api.WrapHandleFunc("/receiveCrossRequest", receiveCrossRequest))
	http.HandleFunc(api.WrapHandleFunc("/mysql", mysqlSelectHandler))
	http.HandleFunc(api.WrapHandleFunc("/redis", redisGetHandler))
	
	http.ListenAndServe(":9099", nil)	
}