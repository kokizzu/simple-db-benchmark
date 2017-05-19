package main

import "fmt"
import _ "github.com/lib/pq"
import "github.com/jmoiron/sqlx"
import "log"

func main() {
	db := sqlx.MustConnect(`postgres`, `user=test2 dbname=test2 sslmode=disable`)
	var err error
	fmt.Println(`test2: postgresql jsonb`)
	_, err = db.Exec(`CREATE TABLE test2 (id BIGSERIAL PRIMARY KEY, k TEXT UNIQUE, v JSONB)`)
	if err != nil {
		_, err = db.Exec(`TRUNCATE TABLE test2`)
		if err != nil {
			log.Fatal(err)
			return
		}
	}
	
	BenchmarkInsert(func(x int) error {
		_, err = db.Exec(fmt.Sprintf(`INSERT INTO test2(k,v)VALUES('%05d','{"v":"%05d"}')`, x, x))
		return err
	})
	BenchmarkUpdate(func(x int) error {
		_, err = db.Exec(fmt.Sprintf(`UPDATE test2 SET v = v || '{"v":"%06d"}' WHERE k = '%05d'`, x, x))
		return err
	})
	BenchmarkSelect(func(x, lim int) error {
		rows, err := db.Queryx(fmt.Sprintf(`SELECT id, k, v->>'v' FROM test2 WHERE k >= '%05d' ORDER BY k ASC LIMIT %d`, x, lim))
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
		rows, err := db.Queryx(fmt.Sprintf(`SELECT id, k, v->>'v' FROM test2 WHERE k <= '%05d' ORDER BY k DESC LIMIT %d`, x, lim))
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

