package main

import "fmt"
import "github.com/gocql/gocql"
import "time"
import "log"

func main() {
	clust := gocql.NewCluster(`172.17.0.2`)
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
	
	max := 9999
	t := time.Now()
	for x := 1; x < max; x++ {
		err = db.Query(fmt.Sprintf(`INSERT INTO test4(bucket,k,v)VALUES('foo','%05d','%05d')`, x, x)).Exec()
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
		err = db.Query(fmt.Sprintf(`UPDATE test4 SET v = '%06d' WHERE bucket = 'foo' AND k = '%05d'`, x, x)).Exec()
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
			iter := db.Query(fmt.Sprintf(`SELECT k, v FROM test4 WHERE bucket = 'foo' AND  k >= '%05d' ORDER BY k ASC LIMIT 20`, x)).Consistency(gocql.One).Iter()
			for {
				m := map[string]interface{}{}
				if !iter.MapScan(m) {
					break
				}
			}
			iter.Close()
		}
		for x := 1; x < max; x += y {
			iter := db.Query(fmt.Sprintf(`SELECT k, v FROM test4 WHERE bucket = 'foo' AND k <= '%05d' ORDER BY k DESC LIMIT 20`, x)).Consistency(gocql.One).Iter()
			for {
				m := map[string]interface{}{}
				if !iter.MapScan(m) {
					break
				}
			}
			iter.Close()
		}
		fmt.Print(`.`)
	}
	fmt.Printf("SELECT: %v\n", time.Now().Sub(t))
	
}

