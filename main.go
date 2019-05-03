package main

import (
	"database/sql"
	"fmt"
    "net/http"
    "html/template"
	_ "github.com/go-sql-driver/mysql"
	"mux"
)

type done struct {
	d bool
}

type cal struct {
	Income int
	Expenditure int
	Mode string
	Balance int
}

type user struct {
	Name string
	Email string
	Pass string
}

func dbConn() (db *sql.DB) {
    dbDriver := "mysql"
    dbUser := "root"
    dbPass := "root"
    dbName := "users"
    db, err := sql.Open(dbDriver, dbUser+":"+dbPass+"@tcp(127.0.0.1:3306)/"+dbName)
    if err != nil {
        panic(err.Error())
    }
    return db
}


//This is done
func UserInfo(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("userinfo.html"))
	db := dbConn()
    selDB, err := db.Query("SELECT * FROM info")
    if err != nil {
        panic(err.Error())
    }
	u := user{}
	res := []user{}
    for selDB.Next() {
        var name, email, pass string
        err = selDB.Scan(&name,&email,&pass)
        if err != nil {
            panic(err.Error())
        }
        u.Name = name
        u.Email = email
		u.Pass = pass
		fmt.Println(u.Name+" "+u.Email+" "+u.Pass)
		res = append(res,u)
		
    }
    tmpl.Execute(w,res)
    defer db.Close()
}

//Done
func Signup(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("signup.html"))
	dope := done{}
	
	db := dbConn()
    u := user{
		Name : r.FormValue("name"),
		Email : r.FormValue("mail"),
		Pass : r.FormValue("pass"),
	}
	fmt.Println(u)
	insVal, err1 := db.Prepare("INSERT INTO info VALUES(?,?,?)")
	if err1 != nil {
		panic(err1.Error())
	}
	insVal.Exec(u.Name,u.Email,u.Pass)
	stmt := "CREATE TABLE "+u.Email+" (income int(100), expenditure int(100), mode varchar(100),balance int(100))"
	fmt.Println(stmt)
	creaDb, err := db.Prepare(stmt)
	if err != nil {
		panic(err.Error())
	}
	creaDb.Exec()
	dope.d = true
    defer db.Close()
    tmpl.Execute(w,dope)
}

//User authentication =====>  Done

func AuthUser(w http.ResponseWriter, r *http.Request){
	db := dbConn()
	email := r.FormValue("usermail")
	pass := r.FormValue("userpass")
	_, err := db.Query(`SELECT * FROM info WHERE email="`+email+`" AND pass="`+pass+`"`)
	if err != nil{
		fmt.Fprintf(w,"<h2>Hey, I guess you have entered wrong email or password</h2>")
	}else {
		calDb,err2 := db.Query("SELECT * FROM "+email)
		if err2 != nil {
			panic(err2.Error())
		}
		ca := cal{}
		res := []cal{}
		for calDb.Next() {
			var income,expenditure,balance int
			var mode string
			err1 := calDb.Scan(&income,&expenditure,&mode,&balance)
			if err1 != nil {
				panic(err1.Error())
			}
			ca.Income = income
			ca.Expenditure = expenditure
			ca.Mode = mode
			ca.Balance = balance
			res = append(res,ca)
		}
		tmpl := template.Must(template.ParseFiles("afterLogin.html"))
		tmpl.Execute(w,res)
	}
	defer db.Close()
}

/*func Calculate(w http.ResponseWriter, r *http.Request) {
	inc := int(r.ParseForm("income"))
	exp := int(r.ParseForm("expenditure"))
	mode := r.ParseForm("mode")

	db := dbConn()
	insVal, err := db.Prepare("INSERT INTO "+u.Email+" VALUES(?,?,?)")
	insVal.Exec(inc,exp,mode)

	newVal, err := db.Query("SELECT * FROM "+u.Email)

}
*/


//===================Delete===========================
func Delete(w http.ResponseWriter, r *http.Request) {
    db := dbConn()
    email := r.URL.Query().Get("email")
    delForm, err := db.Prepare("DELETE FROM info WHERE email=?")
    if err != nil {
        panic(err.Error())
    }
    delForm.Exec(email)
	defer db.Close()
	u := user{}
	res := []user{}
	selDb, err1 := db.Query("SELECT * FROM info")
	if err1 != nil {
		panic(err1.Error())
	}
	for selDb.Next() {
		var name, email, pass string
		err2 := selDb.Scan(&name,&email,&pass)
		if err2 != nil {
			panic(err2.Error())
		}
		u.Name = name
		u.Email = email
		u.Pass = pass
		res = append(res,u)
	}
	tmpl := template.Must(template.ParseFiles("userinfo.html"))
	tmpl.Execute(w,res)
}

//index ============================================================
func index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w,`<h1>Cash Memo App</h1>
	<h2>Log in for Admin</h2>
<form action="/UserInfo" method="POST">
    <label>Email:</label><br>
    <input type="email" name="adminmail"><br>
    <label>Password:</label><br>
    <input type="password" name="adminpass"><br>
    <input type="submit" value="sign in">
</form>

<h2>Log in for User</h2>
<form action="/AuthUser" method="POST">
    <label>Email:</label><br>
    <input type="email" name="usermail"><br>
    <label>Password:</label><br>
    <input type="password" name="userpass"><br>
    <input type="submit" value="sign in">
</form>

<form action="/Signup">
<input type="submit" value="sign up">
</form>
	`)
}

func main() {
	r := mux.NewRouter()
	r.Handle('/images/{rest}',http.StripPrefix("/fold/",http.FileServer(http.Dir(HomeFolder + "fold/"))))
	fs := http.FileServer(http.Dir("./fold/"))
	r.pathPrefix("/fold/").Handler(http.StripPrefix("/fold/",fs))
	r.HandleFunc("/",index)
	r.HandleFunc("/UserInfo/",UserInfo)
	r.HandleFunc("/Signup/",Signup)
	r.HandleFunc("/AuthUser/",AuthUser)
	r.HandleFunc("/Delete/",Delete)
	r.ListenAndServe(":8080",r)
}