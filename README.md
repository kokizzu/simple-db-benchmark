# simple-db-benchmark

Testing PostgreSQL 9.6.2, ScyllaDB 1.6.4, CockroachDB 1.0, result: http://kokizzu.blogspot.sg/2017/05/postgresql-vs-cockroachdb-vs-scylladb.html

## Preparation

```
# test1: postgresql
sudo su - postgres <<EOF
createuser test1
createdb test1
psql -c 'GRANT ALL PRIVILEGES ON DATABASE test1 TO test1;'
EOF

# test2: postgresql jsonb
sudo su - postgres <<EOF
createuser test2
createdb test2
psql -c 'GRANT ALL PRIVILEGES ON DATABASE test2 TO test2;'
EOF

# test3: cockroachdb
cockroach start --insecure
cockroach sql --insecure
CREATE DATABASE test4;

# test4: scylladb
dir=`pwd`
x=1
mkdir -p $dir/scylla$x/commitlog $dir/scylla$x/data
docker stop scylla$x
docker rm scylla$x
docker run --volume $dir/scylla$x:/var/lib/scylla --name scylla$x \
  -d scylladb/scylla --developer-mode 1 --memory 1G --smp 2
docker logs scylla$x | tail
sleep 2;
docker exec -it scylla$x nodetool status
docker exec -it scylla$x cqlsh 
CREATE KEYSPACE test4 WITH REPLICATION = {'class':'SimpleStrategy', 'replication_factor':1};
```