package main

import "fmt"
import _ "github.com/lib/pq"
import "github.com/jmoiron/sqlx"
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
	
	BenchmarkInsert(func(x int) error {
		_, err = db.Exec(fmt.Sprintf(`INSERT INTO test3(k,v)VALUES('%05d','%05d')`, x, x))
		return err
	})
	BenchmarkUpdate(func(x int) error {
		_, err = db.Exec(fmt.Sprintf(`UPDATE test3 SET v = '%06d' WHERE k = '%05d'`, x, x))
		return err
	})
	BenchmarkSelect(func(x, lim int) error {
		rows, err := db.Queryx(fmt.Sprintf(`SELECT id,k,v FROM test3 WHERE k >= '%05d' ORDER BY k ASC LIMIT %d`, x, lim))
		if err != nil {
			return err
		}
		defer rows.Close()
		for rows.Next() {
			m := map[string]interface{}{}
			rows.MapScan(m)
		}
		return nil
	}, func(x, lim int) error {
		rows, err := db.Queryx(fmt.Sprintf(`SELECT id,k,v FROM test3 WHERE k <= '%05d' ORDER BY k DESC LIMIT %d`, x, lim))
		if err != nil {
			return err
		}
		defer rows.Close()
		for rows.Next() {
			m := map[string]interface{}{}
			rows.MapScan(m)
		}
		return nil
	})
	
}

