package main

import (
	"github.com/bonreeapm/go"
	"github.com/bonreeapm/go/common"
	"github.com/bonreeapm/go/compatible/brmartini"
	"database/sql"
	"net/http"
	"strconv"
	"strings"

	"log"

	"github.com/go-martini/martini"
	_ "github.com/lib/pq"
)

var (
	id   int
	name string
)

func main() {
	app, err := bonree.NewApplication("martini")
	if err != nil {
		log.Fatalln(err)
		return
	}

	defer app.Release()

	m := martini.Classic()
	// bonree:添加中间件
	m.Use(brmartini.MiddlewareHandle(app))
	//普通的GET方式路由
	m.Get("/test", func(r *http.Request) string {
		// bonree:获取businessTransaction(业务)
		btn := brmartini.GetCurrentTransaction(r)

		if btn == nil {
			return "Get Transaction fail"
		}

		// bonree:创建快照方法
		snapshotFunc := btn.SnapshotFuncStart("main", "test")

		connStr := "postgres://apmtest:Apmceshi123@10.221.150.179:5444/edb?sslmode=disable"
		db, err := sql.Open("postgres", connStr)
		if err != nil {
			log.Fatal(err)
		}
		defer db.Close()

		// bonree:创建exitCall(后端)
		exitcall := btn.StartSQLExitCall(common.BR_SQL_TYPE_POSTGRESQL, "10.221.150.179", 5444, "postgres", "postgres", "")

		snapshotFunc.AddExitCall(exitcall)

		// bonree:结束后端
		defer exitcall.End()
		// bonree:结束快照方法
		defer btn.SnapshotFuncEnd(snapshotFunc)

		// bonree:设置后端详情
		exitcall.SetDetail("SELECT id, name FROM testtable", "SELECT id, name FROM testtable")

		rows, err := db.Query("SELECT id, name FROM testtable")

		if err != nil {
			// bonree:设置后端错误
			exitcall.AddError("POSTGRESQL Error", err.Error(), err.Error(), true)

			// bonree:如果使用Fatal会使程序退出，退出前需结束快照方法后端和业务
			btn.SnapshotFuncEnd(snapshotFunc)
			exitcall.End()
			btn.End()
			log.Fatal(err)
		}
		var array = make([]string, 0)
		for rows.Next() {
			err := rows.Scan(&id, &name)
			if err != nil {
				// bonree:设置后端错误
				exitcall.AddError("POSTGRESQL Error", err.Error(), err.Error(), true)

				// bonree:如果使用Fatal会使程序退出，退出前需结束后端和业务
				exitcall.End()
				btn.End()
				log.Fatal(err)
			}
			array = append(array, "("+strconv.Itoa(id)+","+name+")")
		}
		ret := "[" + strings.Join(array, ",") + "]"

		return ret
	})
	m.RunOnAddr(":58080")
}
