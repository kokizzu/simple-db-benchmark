# simple-db-benchmark

* [Cockroach 2.1.3](https://kokizzu.blogspot.com/2019/01/cockroachdb-213-benchmark.html)
* [PostgreSQL 9.6.2 vs ScyllaDB 1.6.4 vs CockroachDB 1.0](http://kokizzu.blogspot.sg/2017/05/postgresql-vs-cockroachdb-vs-scylladb.html)
* [PostgreSQL 9.6.2 vs ScyllaDB 1.7RC2](http://kokizzu.blogspot.co.id/2017/05/postgresql-962-vs-scylladb-17rc2.html) 

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
CREATE DATABASE test3;
GRANT ALL ON DATABASE test3 TO test3;

# test4: scylladb (under docker)
dir=`pwd`
x=1
mkdir -p $dir/scylla$x/commitlog $dir/scylla$x/data
docker stop scylla$x
docker rm scylla$x
docker run --volume $dir/scylla$x:/var/lib/scylla --name scylla$x \
  -d scylladb/scylla --developer-mode 1 --memory 4G --smp 4
docker logs scylla$x | tail
sleep 2;
docker exec -it scylla$x nodetool status
docker exec -it scylla$x cqlsh 
CREATE KEYSPACE test4 WITH REPLICATION = {'class':'SimpleStrategy', 'replication_factor':1};

# test4: scylladb (supported: ubuntu 17.04+xfs)
sudo wget -O /etc/apt/sources.list.d/scylla.list http://downloads.scylladb.com/deb/ubuntu/scylla-1.7-xenial.list
sudo apt-get update
sudo systemctl enable scylla-server
sudo scylla_setup
sudo sed -i 's|/usr/bin/scylla $SCYLLA_ARGS|/usr/bin/scylla -m 8G -c 8 $SCYLLA_ARGS|g' /lib/systemd/system/scylla-server.service
sudo systemctl daemon-reload
sudo systemctl start scylla-server # cqlsh 127.0.0.1
CREATE KEYSPACE test4 WITH REPLICATION = {'class':'SimpleStrategy', 'replication_factor':1};
```
