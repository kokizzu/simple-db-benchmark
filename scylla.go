package main

import "fmt"
import "github.com/gocql/gocql"
import (
	"log"
	"time"
)

func main() {
	clust := gocql.NewCluster(`127.0.0.1`)
	clust.Timeout = 8 * time.Second
	clust.RetryPolicy = &gocql.SimpleRetryPolicy{NumRetries:3}
	clust.Keyspace = `test4`
	db, err := clust.CreateSession()
	defer db.Close()
	if err != nil {
		log.Fatal(err)
		return
	}
	fmt.Println(`test4: scylladb`)
	
	if err := db.Query(`DROP TABLE test4`).Exec(); err != nil {
		log.Println(err)
	}
	if err := db.Query(`CREATE TABLE test4 (bucket text, k text, v TEXT, PRIMARY KEY(bucket,k))`).Exec(); err != nil {
		log.Fatal(err)
		return
	}
	
	BenchmarkInsert(func(x int) error {
		err = db.Query(fmt.Sprintf(`INSERT INTO test4(bucket,k,v)VALUES('foo','%05d','%05d')`, x, x)).Exec()
		return err
	})
	BenchmarkUpdate(func(x int) error {
		err = db.Query(fmt.Sprintf(`UPDATE test4 SET v = '%06d' WHERE bucket = 'foo' AND k = '%05d'`, x, x)).Exec()
		return err
	})
	BenchmarkSelect(func(x, lim int) error {
		sql := fmt.Sprintf(`SELECT k, v FROM test4 WHERE bucket = 'foo' AND  k >= '%05d' ORDER BY k ASC LIMIT %d`, x, lim)
		iter := db.Query(sql).Iter()
		defer iter.Close()
		tot := 0
		for {
			m := map[string]interface{}{}
			if !iter.MapScan(m) {
				break
			}
			tot++
		}
		if tot == 0 {
			return fmt.Errorf(`Empty result: `+sql)
		}
		return nil
	}, func(x, lim int) error {
		sql := fmt.Sprintf(`SELECT k, v FROM test4 WHERE bucket = 'foo' AND k <= '%05d' ORDER BY k DESC LIMIT %d`, x, lim)
		iter := db.Query(sql).Iter()
		defer iter.Close()
		tot := 0
		for {
			m := map[string]interface{}{}
			if !iter.MapScan(m) {
				break
			}
			tot++
		}
		if tot == 0 {
			return fmt.Errorf(`Empty result: `+sql)
		}
		return nil
	})
}

