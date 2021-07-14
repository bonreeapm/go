package main

import(
	"github.com/bonreeapm/go"
	"github.com/bonreeapm/go/compatible/brgorilla"
	"github.com/bonreeapm/go/common"
	"time"
	"net/http"
	"strconv"
	"fmt"	
	"log"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/garyburd/redigo/redis"

	"github.com/gorilla/mux"
)

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

	defer db.Close();

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
	app, err := bonree.NewApplication("gorilla")
	if err != nil {
		log.Fatalln(err)
		return
	}

	defer app.Release()

	r := mux.NewRouter()

	r.Handle("/setURL", http.HandlerFunc(setURL))
	r.Handle("/addError", http.HandlerFunc(addError))
	r.Handle("/sendCrossRequest", http.HandlerFunc(sendCrossRequest))
	r.Handle("/receiveCrossRequest", http.HandlerFunc(receiveCrossRequest))
	r.Handle("/mysql", http.HandlerFunc(mysqlSelectHandler))
	r.Handle("/redis", http.HandlerFunc(redisGetHandler))
	
	http.ListenAndServe(":9099", brgorilla.InstrumentRoutes(r, app))
}