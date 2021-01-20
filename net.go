package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	age "github.com/bearbin/go-age"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

type customer struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	DOB  string `json:"dob"`
	Addr Address
}
type Address struct {
	ID         int
	Streetname string
	City       string
	State      string
	Customerid int
}

func getDBConnection() *sql.DB {
	db, err := sql.Open("mysql", "root:password@/customer_service")
	if err != nil {
		log.Fatal(err)
	}
	return db
}
func getDOB(year, month, day int) time.Time {
	dob := time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC)
	return dob
}
func GetName(w http.ResponseWriter, r *http.Request) {
	db := getDBConnection()
	var names []interface{}
	query := r.URL.Query()
	name := query.Get("name")
	q := `SELECT * FROM customer INNER JOIN Address ON customer.ID=Address.Customerid`
	if len(name) != 0 {
		q = `SELECT * FROM customer INNER JOIN Address ON customer.ID=Address.Customerid WHERE customer.Name=?`
		names = append(names, name)
	}
	rows, err := db.Query(q, names...)
	if err != nil {
		panic(err.Error())
	}

	var ans []customer
	for rows.Next() {
		var c customer
		if err := rows.Scan(&c.ID, &c.Name, &c.DOB, &c.Addr.ID, &c.Addr.Streetname, &c.Addr.City, &c.Addr.State, &c.Addr.Customerid); err != nil {
			log.Fatal(err)
		}
		ans = append(ans, c)
	}
	Byte, _ := json.Marshal(ans)
	_, err = io.WriteString(w, string(Byte))
	if err != nil {
		panic(err)
	}
}
func GetID(w http.ResponseWriter, r *http.Request) {
	db := getDBConnection()
	param := mux.Vars(r)
	rows, err := db.Query("SELECT * FROM customer INNER JOIN Address ON customer.ID=Address.Customerid WHERE customer.ID=?", param["id"])
	if err != nil {
		panic(err)
	}
	var c customer
	for rows.Next() {
		if err := rows.Scan(&c.ID, &c.Name, &c.DOB, &c.Addr.ID, &c.Addr.Streetname, &c.Addr.City, &c.Addr.State, &c.Addr.Customerid); err != nil {
			panic(err)
		}
	}
	if len(c.Name) != 0 {
		Byte, _ := json.Marshal(c)
		io.WriteString(w, string(Byte))
	} else {
		w.WriteHeader(http.StatusNoContent)
	}
	//fmt.Println("ok")
	//param := mux.Vars(r)
	//fmt.Println(param["id"])
	//var ids []interface{}
	//ids = append(ids, param["id"])
	//query := `SELECT * FROM customer INNER JOIN Address ON customer.ID=Address.Customerid WHERE customer.ID=?;`
	//rows, err := db.Query(query, ids...)
	//if err != nil {
	//	panic(err.Error())
	//}
	//if !rows.Next() {
	//	w.WriteHeader(http.StatusBadRequest)
	//	err = json.NewEncoder(w).Encode([]customer(nil))
	//	if err != nil {
	//		panic(err)
	//	}
	//} else {
	//	fmt.Println("in else", rows.Next())
	//	var c customer
	//	//var ans []customer
	//	for rows.Next() {
	//		fmt.Println("in for")
	//		if err := rows.Scan(&c.ID, &c.Name, &c.DOB, &c.Addr.ID, &c.Addr.Streetname, &c.Addr.City, &c.Addr.State, &c.Addr.Customerid); err != nil {
	//			fmt.Println("in error c")
	//			log.Fatal(err)
	//		}
	//		fmt.Println("c first is", c)
	//		//ans = append(ans, c)
	//	}
	//	fmt.Println("c is", c)
	//	//if len(c.Name) == 0 {
	//	//	w.WriteHeader(http.StatusNoContent)
	//	//	err = json.NewEncoder(w).Encode([]customer(nil))
	//	//	if err != nil {
	//	//		panic(err)
	//	//	}
	//	//} else {
	//	Byte, _ := json.Marshal(c)
	//	fmt.Println("Byte is:", string(Byte))
	//	//json.NewEncoder(w).Encode(c)
	//	io.WriteString(w, string(Byte))
	//	//if err != nil {
	//	//	panic(err)
	//	//}
	//}
}

func putcustomer(w http.ResponseWriter, r *http.Request) {
	db := getDBConnection()
	body, _ := ioutil.ReadAll(r.Body)
	var c customer
	err := json.Unmarshal(body, &c)
	if err != nil {
		fmt.Println("in err ")
		log.Fatal(err)
	}
	if c.DOB == "" {
		param := mux.Vars(r)
		id := param["id"]
		if c.Name != "" {
			_, err := db.Exec("update customer set Name=? where ID=?", c.Name, id)
			if err != nil {
				panic(err.Error())
				err = json.NewEncoder(w).Encode(customer{})
				if err != nil {
					panic(err)
				}
			}
		}
		var data []interface{}
		query := "update Address set "
		if c.Addr.City != "" {
			query += "City = ? ,"
			data = append(data, c.Addr.City)
		}
		if c.Addr.State != "" {
			query += "State = ? ,"
			data = append(data, c.Addr.State)
		}
		if c.Addr.Streetname != "" {
			query += "StreetNumber = ? ,"
			data = append(data, c.Addr.Streetname)
		}
		query = query[:len(query)-1]
		query += "where Customerid = ? and ID = ?"
		data = append(data, id)
		data = append(data, c.Addr.ID)
		_, err = db.Exec(query, data...)

		//if err != nil {
		//	log.Fatal(err)
		//}
		nid, _ := strconv.Atoi(param["id"])
		c.ID = nid
		err = json.NewEncoder(w).Encode(c)
		if err != nil {
			panic(err)
		}
	} else {
		GetID(w, r)
	}
}

func postcustomer(w http.ResponseWriter, r *http.Request) {
	db := getDBConnection()
	var c customer
	bytevalue, _ := ioutil.ReadAll(r.Body)
	if len(bytevalue) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		err := json.NewEncoder(w).Encode("Give inputs in body to post")
		if err != nil {
			panic(err)
		}
	} else {
		err := json.Unmarshal(bytevalue, &c)
		//fmt.Println(c.DOB)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			err = json.NewEncoder(w).Encode("Incorrect Format")
			if err != nil {
				panic(err)
			}
		}
		if len(c.Name) == 0 || len(c.DOB) == 0 || len(c.Addr.State) == 0 || len(c.Addr.City) == 0 || len(c.Addr.Streetname) == 0 {
			w.WriteHeader(http.StatusBadRequest)
			err = json.NewEncoder(w).Encode("Missing Fields")
			if err != nil {
				panic(err)
			}
		} else {
			dob := c.DOB
			dob1 := strings.Split(dob, "/")
			y, _ := strconv.Atoi(dob1[2])
			m, _ := strconv.Atoi(dob1[1])
			d, _ := strconv.Atoi(dob1[0])
			fmt.Println("C :", c)
			getAge := getDOB(y, m, d)
			fmt.Println(age.Age(getAge))
			if age.Age(getAge) >= 18 {
				var value []interface{}

				query := `INSERT INTO customer VALUES(?,?,?)`
				value = append(value, 0)
				value = append(value, c.Name)
				value = append(value, c.DOB)
				rows, err := db.Exec(query, value...)
				if err != nil {
					panic(err.Error())
				}
				idAddr, _ := rows.LastInsertId()
				var Addrvalue []interface{}
				query = `INSERT INTO Address Values(?,?,?,?,?)`
				Addrvalue = append(Addrvalue, c.Addr.ID)
				Addrvalue = append(Addrvalue, c.Addr.Streetname)
				Addrvalue = append(Addrvalue, c.Addr.City)
				Addrvalue = append(Addrvalue, c.Addr.State)
				Addrvalue = append(Addrvalue, idAddr)
				rows2, err := db.Exec(query, Addrvalue...)
				//GetID(w,r)
				idAddr1, _ := rows2.LastInsertId()
				row, err := db.Query("SELECT * FROM customer INNER JOIN Address ON customer.ID=Address.Customerid WHERE customer.ID=?;", idAddr)
				if err != nil {
					panic(err.Error())
				}

				var ans []customer
				for row.Next() {
					//var c customer
					c.ID = int(idAddr)
					c.Addr.ID = int(idAddr1)
					c.Addr.Customerid = int(idAddr)
					if err := row.Scan(&c.ID, &c.Name, &c.DOB, &c.Addr.ID, &c.Addr.Streetname, &c.Addr.City, &c.Addr.State, &c.Addr.Customerid); err != nil {
						log.Fatal(err)
					}
					ans = append(ans, c)
				}
				Byte, _ := json.Marshal(ans)
				_, err = io.WriteString(w, string(Byte))
				if err != nil {
					panic(err)
				}
			} else {
				w.WriteHeader(http.StatusBadRequest)
				err = json.NewEncoder(w).Encode("Age is less than 18")
				if err != nil {
					panic(err)
				}
			}
		}
	}
}
func delcust(w http.ResponseWriter, r *http.Request) {
	db := getDBConnection()
	params := mux.Vars(r)
	GetID(w, r)
	_, err1 := db.Query("DELETE FROM Address WHERE Customerid=?;", params["id"])
	if err1 != nil {
		panic(err1.Error())
	}
	//fmt.Printf("%T, %d", params["id"], params["id"])
	//_, err1 = stmt1.Exec(params["id"])
	//db, err = sql.Open("mysql", "sumit:1234@/Cust_Service")
	//defer db.Close()

	if err1 != nil {
		panic(err1)
	}
	if err1 != nil {
		panic(err1.Error())
	}
	_, err := db.Query("DELETE FROM customer WHERE ID=?;", params["id"])
	if err != nil {
		panic(err.Error())
	}
	//fmt.Printf("%T, %d", params["id"], params["id"])
	//_, err = stmt.Exec(params["id"])
	//db, err = sql.Open("mysql", "sumit:1234@/Cust_Service")
	//defer db.Close()

	if err != nil {
		panic(err)
	}
	if err != nil {
		panic(err.Error())
	}

	//fmt.Fprintf(w, "Post with ID = %s was deleted", params["id"])
}
func main() {
	//db,err=sql.Open("mysql","root:password@/customer_service")
	//if err != nil {
	//	panic(err)
	//}
	//defer db.Close()
	router := mux.NewRouter()
	router.HandleFunc("/customer", GetName).Methods("GET")
	router.HandleFunc("/customer/{id:[0-9]+}", GetID).Methods(http.MethodGet)
	router.HandleFunc("/customer", postcustomer).Methods("POST")
	router.HandleFunc("/customer/{id:[0-9]+}", putcustomer).Methods("PUT")
	router.HandleFunc("/customer/{id:[0-9]+}", delcust).Methods("DELETE")
	err := http.ListenAndServe(":8080", router)
	if err != nil {
		panic(err)
	}
}
