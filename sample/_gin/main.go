package main

import(
	"github.com/bonreeapm/go"
	"github.com/bonreeapm/go/common"
	"github.com/bonreeapm/go/compatible/brgin"
	"time"
	"log"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/garyburd/redigo/redis"

	"github.com/gin-gonic/gin"
)

func setURL(c *gin.Context) {	
	btn := brgin.GetCurrentTransaction(c)

	if btn == nil {
		c.Writer.WriteString("Get Transaction fail")
		return
	}
	
	time.Sleep(time.Duration(3)*time.Second)
}

func addError(c *gin.Context) {
	btn := brgin.GetCurrentTransaction(c)

	if btn == nil {
		c.Writer.WriteString("Get Transaction fail")
		return
	}

	btn.AddError("UnkownException","UnkownException: something wrong","UnkownException: something wrong in File xxx.cpp:12", true)

	time.Sleep(time.Duration(3)*time.Second)
}

const mysqldb = "root:111111@tcp(192.168.0.201:3306)/test"
func mysqlSelectHandler(c *gin.Context) {
	btn := brgin.GetCurrentTransaction(c)

	if btn == nil {
		c.Writer.WriteString("Get Transaction fail")
		return
	}

	db, err := sql.Open("mysql", mysqldb)

	if err != nil {
		return
	}

	defer db.Close();

	snapshotFunc := btn.SnapshotFuncStart("main", "mysqlSelectHandler")
	exitcall := btn.StartSQLExitCall(common.BR_SQL_TYPE_MYSQL, "192.168.0.201", 3306, "test", "mysql", "")
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

	c.Writer.WriteString("mysqlSelectHandler.")
	return
}

func redisGetHandler(c *gin.Context) {
	btn := brgin.GetCurrentTransaction(c)

	if btn == nil {
		c.Writer.WriteString("Get Transaction fail")
		return
	}

	conn, err := redis.Dial("tcp", "192.168.0.201:6379")
	if err != nil {
		return
	}

	snapshotFunc := btn.SnapshotFuncStart("main", "redisGetHandler")
	exitcall := btn.StartNoSQLExitCall(common.BR_NOSQL_TYPE_REDIS, "192.168.0.201", 6379, "redis")
	snapshotFunc.AddExitCall(exitcall)
	defer exitcall.End()
	defer btn.SnapshotFuncEnd(snapshotFunc)

	exitcall.SetDetail("GET", "GET test_value")

	_, err = conn.Do("GET", "test_value")

	if err != nil {
		exitcall.AddError("Redis Error", err.Error(), err.Error(), true)
	}

	c.Writer.WriteString("redisGetHandler.")
	return
}

func main() {
	app, err := bonree.NewApplication("gin")
	if err != nil {
		log.Fatalln(err)
		return
	}

	defer app.Release()

	router := gin.Default()
	router.Use(brgin.Middleware(app))

	router.GET("/setURL", setURL)
	router.GET("/addError", addError)
	router.GET("/mysql", mysqlSelectHandler)
	router.GET("/redis", redisGetHandler)

	router.Run(":9099")
}