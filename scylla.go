package main

import "fmt"
import "github.com/gocql/gocql"
import "time"
import "log"

func main() {
	clust := gocql.NewCluster(`127.0.0.1`)
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
	if err := db.Query(`CREATE TABLE test4 (bucket text, k text, v TEXT, PRIMARY KEY(bucket,k))`).RetryPolicy(nil).Exec(); err != nil {
		log.Fatal(err)
		return
	}
	
	const max = 9999
	const jump = 40
	t := time.Now()
	for x := 1; x <= max; x++ {
		err = db.Query(fmt.Sprintf(`INSERT INTO test4(bucket,k,v)VALUES('foo','%05d','%05d')`, x, x)).Exec()
		if err != nil {
			log.Fatal(err)
			return
		}
		if x % 200 == 0 {
			fmt.Print(`.`)
		}
	}
	dur := time.Now().Sub(t)
	fmt.Printf("INSERT: %v (%.2f ms/op)\n", dur, float64(dur.Nanoseconds()) / 1000000 / max)
	
	t = time.Now()
	for x := 1; x <= max; x++ {
		err = db.Query(fmt.Sprintf(`UPDATE test4 SET v = '%06d' WHERE bucket = 'foo' AND k = '%05d'`, x, x)).Exec()
		if err != nil {
			log.Fatal(err)
			return
		}
		if x % 200 == 0 {
			fmt.Print(`.`)
		}
	}
	dur = time.Now().Sub(t)
	fmt.Printf("UPDATE: %v (%.2f ms/op)\n", dur, float64(dur.Nanoseconds()) / 1000000 / max)
	
	t = time.Now()
	ops := int64(0)
	for y := 2; y < jump; y++ {
		for x := max - 1; x > 0; x -= y {
			ops++
			iter := db.Query(fmt.Sprintf(`SELECT k, v FROM test4 WHERE bucket = 'foo' AND  k >= '%05d' ORDER BY k ASC LIMIT %d`, x, y * y)).Iter()
			for {
				m := map[string]interface{}{}
				if !iter.MapScan(m) {
					break
				}
			}
			iter.Close()
			if ops % 500 == 0 {
				fmt.Print(`.`)
			}
		}
		for x := 1; x < max; x += y {
			ops++
			iter := db.Query(fmt.Sprintf(`SELECT k, v FROM test4 WHERE bucket = 'foo' AND k <= '%05d' ORDER BY k DESC LIMIT %d`, x, y * y)).Iter()
			for {
				m := map[string]interface{}{}
				if !iter.MapScan(m) {
					break
				}
			}
			iter.Close()
			if ops % 500 == 0 {
				fmt.Print(`.`)
			}
		}
	}
	dur = time.Now().Sub(t)
	fmt.Printf("SELECT: %v (%.2f ms/op)\n", dur, float64(dur.Nanoseconds()) / 1000000 / float64(ops))
	
}

