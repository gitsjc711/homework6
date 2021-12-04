package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"log"
)
var(
	name="root"
	pword=""
	host="localhost"
	port="3306"
	dbname="test"
)
var DB *sql.DB
type user struct {
	id int
	account string//账号
	password string//密码
	mibao string//密保
	tishici string//提示词
}
func IntoDB(){
	dns:=name+":"+pword+"@tcp("+host+":"+port+")/"+dbname
	db,err:=sql.Open("mysql",dns)
	if err!=nil{
		log.Fatal(err)
	}
	DB=db
}
func prepareInsertDemo(a,b,c,d string) {//注册
	sqlStr := "insert into user(account, password,mibao,tishici) values (?,?,?,?)"
	stmt, err := DB.Prepare(sqlStr)
	if err != nil {
		fmt.Printf("prepare failed, err:%v\n", err)
		return
	}
	defer stmt.Close()
	_, err = stmt.Exec(a,b,c,d)
	if err != nil {
		fmt.Printf("insert failed, err:%v\n", err)
		return
	}
	fmt.Println("insert success.")
}
func prepareQueryDemo()(a,b,c,d []string) {//查找是否重复，读取账号置入切片，最多存储100（应该够了）

	a = make([]string, 100)
	b = make([]string, 100)
	c = make([]string, 100)
	d = make([]string, 100)
	sqlStr := "select id, account,password,mibao,tishici from user where id > ?"
	stmt, err := DB.Prepare(sqlStr)
	if err != nil {
		fmt.Printf("prepare failed, err:%v\n", err)
		return
	}
	defer stmt.Close()
	rows, err := stmt.Query(0)
	if err != nil {
		fmt.Printf("query failed, err:%v\n", err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		var u user
		err := rows.Scan(&u.id,&u.account,&u.password,&u.mibao,&u.tishici)
		if err != nil {
			fmt.Printf("scan failed, err:%v\n", err)
			return
		}
		a[u.id]=u.account
		b[u.id]=u.password
		c[u.id]=u.mibao
		d[u.id]=u.tishici
	}
	return a,b,c,d
}
func main() {
	IntoDB()
	var account, password, mibao, tishici string

	for {
		var s string
		fmt.Println("请输入要进行的操作，注册，登入或找回密码")
		fmt.Scanln(&s)
		switch s {
		case "注册":
			fmt.Println("输入账号")
		again1:
			fmt.Scanln(&account)
			s1, _, _,_ := prepareQueryDemo()
			for _, accountin := range s1 { //验证账号是否存在
				if account == "" {
					fmt.Println("账号不能为空")
					goto again1
				}
				if account == accountin {
					fmt.Println("账号已经存在，请重新输入")
					goto again1
				}
			}
		again2:
			fmt.Println("输入密码、提示词、密保")
			fmt.Scanln(&password, &tishici, &mibao)
			if len(password) < 6 {
				fmt.Println("密码过短")
				goto again2
			}
			if tishici == "" || mibao == "" {
				fmt.Println("提示词或密保不能为空")
				goto again2
			}
			prepareInsertDemo(account, password, mibao, tishici)
		case "登入":
			s1, s2, _,_ := prepareQueryDemo()
			fmt.Println("登入输入账号密码")
		again3:
			fmt.Scanln(&account, &password)
			var key bool //判断是否登入成功
			for i := 1; i <100; i++ {
				if account == s1[i] && password == s2[i] {
					fmt.Println("登入成功")
					key = true
				}
			}
			if key == false {
				fmt.Println("账号密码不正确请重新输入")
				goto again3
			}
		case "找回密码":
			fmt.Println("输入账号")
			var missing string//需要找回的账号
			var i int
			var missingmibao string
			again4:fmt.Scanln(&missing)
			s1,s2, s3,s4 := prepareQueryDemo()
			for i = 1; i < 100; i++ {
				if missing == s1[i]  {
					fmt.Println("你的提示词为")
					fmt.Println(s4[i])
					break
				}
			}
			if i==100{
				fmt.Println("账号不存在")
				goto again4
			}
			fmt.Println("请输入密保")
			fmt.Scanln(&missingmibao)
			if missingmibao==s3[i]{
				fmt.Println("你的密码为")
				fmt.Println(s2[i])
			}else{
				fmt.Println("密保错误")
			}
		default:
			fmt.Println("输入非法")
		}
	}
}
