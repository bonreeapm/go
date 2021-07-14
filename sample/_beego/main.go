package main

import (
	"github.com/bonreeapm/go"
	"github.com/bonreeapm/go/common"
	"github.com/bonreeapm/go/compatible/brbeego"
	"github.com/bonreeapm/go/routineEngine"

	"log"
	"time"
	"net/http"
	"strconv"
	"fmt"	
	"github.com/astaxie/beego"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/garyburd/redigo/redis"
)

// MainController is the Controller
type MainController struct {
	beego.Controller
}

// Get is the func
func (mainController *MainController) Get() {
	btn := bonree.GetCurrentTransaction(mainController.Ctx.ResponseWriter.ResponseWriter)

	snapshotFunc := btn.SnapshotFuncStart("MainController", "Get")

	defer btn.SnapshotFuncEnd(snapshotFunc)

	mainController.Ctx.WriteString("hello world")

	time.Sleep(time.Duration(3)*time.Second)
}

func setURL(w http.ResponseWriter, r *http.Request) {	
	btn := bonree.GetCurrentTransaction(w)

	if btn == nil {
		fmt.Fprint(w, "Get Transaction fail")
		return
	}
	fmt.Fprint(w, "setURL")
	
	time.Sleep(time.Duration(3)*time.Second)
}

func routine(w http.ResponseWriter, r *http.Request) {	
	btn := bonree.GetRoutineTransaction()

	if btn == nil {
		fmt.Fprint(w, "Get Transaction fail")
		return
	}
	fmt.Fprint(w, "routine")

	snapshotFunc := btn.SnapshotFuncStart("main", "Routine")

	defer btn.SnapshotFuncEnd(snapshotFunc)
	
	time.Sleep(time.Duration(3)*time.Second)
}

func addError(w http.ResponseWriter, r *http.Request) {
	btn := bonree.GetCurrentTransaction(w)

	if btn == nil {
		fmt.Fprint(w, "Get Transaction fail")
		return
	}
	fmt.Fprint(w, "addError")

	btn.AddError("UnkownException","UnkownException: something wrong","UnkownException: something wrong in File xxx.cpp:12", true)

	time.Sleep(time.Duration(3)*time.Second)
}

func receiveCrossRequest(w http.ResponseWriter, r *http.Request) {	
	time.Sleep(time.Duration(3)*time.Second)
}

func sendCrossRequest(w http.ResponseWriter, r *http.Request) {
	host := "127.0.0.1"
	port := 9099

	btn := bonree.GetCurrentTransaction(w)

	if btn == nil {
		fmt.Fprint(w, "Get Transaction fail")
		return
	} 

	snapshotFunc := btn.SnapshotFuncStart("main", "sendCrossRequest")
	exitcall := btn.StartRPCExitCall(common.BACKEND_TYPE_HTTP, host, port)
	snapshotFunc.AddExitCall(exitcall)
	defer exitcall.End()
	defer btn.SnapshotFuncEnd(snapshotFunc)

	client := &http.Client{}
	//client.Transport = exitcall.RoundTripper()
	exitcall.SetDetail(host+":"+strconv.Itoa(port)+"/receiveCrossRequest", host+":"+strconv.Itoa(port)+"/receiveCrossRequest")
	_, err := client.Get("http://" + host + ":" + strconv.Itoa(port) + "/receiveCrossRequest")

	if err != nil {
		fmt.Fprint(w, err.Error())
		exitcall.AddError("Http Get Error", err.Error(), err.Error(), true)
		return
	}

	//exitcall.SetCrossResheader(resp.Header)

	fmt.Fprint(w, "Send Cross Request.")
	return
}

const mysqldb = "root:brxm@123@tcp(backend.br007.top:3306)/test"
func mysqlSelectHandler(w http.ResponseWriter, r *http.Request) {
	btn := bonree.GetCurrentTransaction(w)

	if btn == nil {
		fmt.Fprint(w, "Get Transaction fail")
		return
	}

	db, err := sql.Open("mysql", mysqldb)

	if err != nil {
		return
	}

	defer db.Close()

	snapshotFunc := btn.SnapshotFuncStart("main", "mysqlSelectHandler")
	exitcall := btn.StartSQLExitCall(common.BACKEND_TYPE_MYSQL, "backend.br007.top", 3306, "test", "PROC")
	snapshotFunc.AddExitCall(exitcall)
	defer exitcall.End()
	defer btn.SnapshotFuncEnd(snapshotFunc)

	exitcall.SetDetail("select * from go_test where id = ?", "select * from go_test where id = ?")
	stmt, err := db.Prepare("select * from go_test where id = ?")

	var id, prop3, prop4 int
	var name, name2 string

	row := stmt.QueryRow(1)

	if row != nil {
		err = row.Scan(&id, &name, &name2, &prop3, &prop4)

		if err != nil {
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

	c, err := redis.Dial("tcp", "backend.br007.top:6379")
	if err != nil {
		return
	}

	snapshotFunc := btn.SnapshotFuncStart("main", "redisGetHandler")
	exitcall := btn.StartNoSQLExitCall(common.BACKEND_TYPE_REDIS, "backend.br007.top", 6379, "StackExchangeRedis")
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
	app, err := bonree.NewApplication("GoAgent_beego")
	if err != nil {
		log.Fatalln(err)
		return
	}

	// 如果需要RoutineEngine支持，则加上此行代码
	bonree.RoutineEngineInit(routineEngine.Get())

	defer app.Release()

	beego.Handler(brbeego.WrapHandleFunc("/routine", routine))
	beego.Handler(brbeego.WrapHandleFunc("/setURL", setURL))
	beego.Handler(brbeego.WrapHandleFunc("/addError", addError))
	beego.Handler(brbeego.WrapHandleFunc("/sendCrossRequest", sendCrossRequest))
	beego.Handler(brbeego.WrapHandleFunc("/receiveCrossRequest", receiveCrossRequest))
	beego.Handler(brbeego.WrapHandleFunc("/mysql", mysqlSelectHandler))
	beego.Handler(brbeego.WrapHandleFunc("/redis", redisGetHandler))

	beego.Router("/", &MainController{})

	brbeego.InsertFilter(app)

	beego.Run(":9099")
}

