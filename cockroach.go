package main

import "fmt"
import _ "github.com/lib/pq"
import "github.com/jmoiron/sqlx"
import "time"
import "log"

func main() {
	
	db := sqlx.MustConnect(`postgres`, `postgresql://test3@localhost:26257/test3?sslmode=disable`)
	var err error
	fmt.Println(`test3: cockroachdb`)
	_, err = db.Exec(`CREATE TABLE test3 (id BIGSERIAL PRIMARY KEY, k TEXT UNIQUE, v TEXT)`)
	if err != nil {
		_, err = db.Exec(`TRUNCATE TABLE test3`)
		if err != nil {
			log.Fatal(err)
			return
		}
	}
	
	max := 9999
	t := time.Now()
	for x := 1; x < max; x++ {
		_, err = db.Exec(fmt.Sprintf(`INSERT INTO test3(k,v)VALUES('%05d','%05d')`, x, x))
		if err != nil {
			log.Fatal(err)
			return
		}
		if x % 100 == 0 {
			fmt.Print(`.`)
		}
	}
	fmt.Printf("INSERT: %v\n", time.Now().Sub(t))
	
	t = time.Now()
	for x := 1; x < max; x++ {
		_, err = db.Exec(fmt.Sprintf(`UPDATE test3 SET v = '%06d' WHERE k = '%05d'`, x, x))
		if err != nil {
			log.Fatal(err)
			return
		}
		if x % 100 == 0 {
			fmt.Print(`.`)
		}
	}
	fmt.Printf("UPDATE: %v\n", time.Now().Sub(t))
	
	t = time.Now()
	for y := 2; y < 39; y++ {
		for x := max - 1; x > 0; x -= y {
			rows, err := db.Queryx(fmt.Sprintf(`SELECT id,k,v FROM test3 WHERE k >= '%05d' ORDER BY k ASC LIMIT 20`, x))
			if err != nil {
				log.Fatal(err)
				return
			}
			for rows.Next() {
				m := map[string]interface{}{}
				rows.MapScan(m)
			}
			rows.Close()
		}
		for x := 1; x < max; x += y {
			rows, err := db.Queryx(fmt.Sprintf(`SELECT id,k,v FROM test3 WHERE k <= '%05d' ORDER BY k DESC LIMIT 20`, x))
			if err != nil {
				log.Fatal(err)
				return
			}
			for rows.Next() {
				m := map[string]interface{}{}
				rows.MapScan(m)
			}
			rows.Close()
		}
		fmt.Print(`.`)
	}
	fmt.Printf("SELECT: %v\n", time.Now().Sub(t))
	
}

