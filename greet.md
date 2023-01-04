
```cql
CREATE KEYSPACE greet WITH replication = {'class': 'SimpleStrategy', 'replication_factor': 3};

USE greet;
CREATE TABLE greet.name (
	name text,
	age int,
	hobby set<text>,
	PRIMARY KEY (name, age)
) WITH CLUSTERING ORDER BY (age DESC);

INSERT INTO greet.name (name, age, hobby) values ('alice', 10, {'code'});
INSERT INTO greet.name (name, age, hobby) values ('alice', 20, {'programming'});
INSERT INTO greet.name (name, age, hobby) values ('bob', 13, {'code'});

EXPAND ON;
SELECT * FROM greet.name;
```

To view how the data is stored, particularly when using _partition key_ with _clustering key_.

```bash
# Flush data to disk.
$ /opt/bitnami/cassandra/bin/nodetool flush

# Go to the directory where the dump is produced
$ cd /bitnami/cassandra/data/data/greet/name-349bc6808be111eda5eb23123d1849fe

$ ls
'?'        nb-1-big-CompressionInfo.db   nb-1-big-Digest.crc32   nb-1-big-Index.db        nb-1-big-Summary.db
 backups   nb-1-big-Data.db              nb-1-big-Filter.db      nb-1-big-Statistics.db   nb-1-big-TOC.txt


# Run sstabledump
$ /opt/bitnami/cassandra/tools/bin/sstabledump nb-1-big-Data.db
```
The output:

```json
[
  {
    "partition" : {
      "key" : [ "bob" ],
      "position" : 0
    },
    "rows" : [
      {
        "type" : "row",
        "position" : 17,
        "clustering" : [ 13 ],
        "liveness_info" : { "tstamp" : "2023-01-04T03:50:49.950938Z" },
        "cells" : [
          { "name" : "hobby", "deletion_info" : { "marked_deleted" : "2023-01-04T03:50:49.950937Z", "local_delete_time" : "2023-01-04T03:50:49Z" } },
          { "name" : "hobby", "path" : [ "code" ], "value" : "" }
        ]
      }
    ]
  },
  {
    "partition" : {
      "key" : [ "alice" ],
      "position" : 42
    },
    "rows" : [
      {
        "type" : "row",
        "position" : 61,
        "clustering" : [ 20 ],
        "liveness_info" : { "tstamp" : "2023-01-04T03:50:49.939226Z" },
        "cells" : [
          { "name" : "hobby", "deletion_info" : { "marked_deleted" : "2023-01-04T03:50:49.939225Z", "local_delete_time" : "2023-01-04T03:50:49Z" } },
          { "name" : "hobby", "path" : [ "programming" ], "value" : "" }
        ]
      },
      {
        "type" : "row",
        "position" : 92,
        "clustering" : [ 10 ],
        "liveness_info" : { "tstamp" : "2023-01-04T03:50:26.277032Z" },
        "cells" : [
          { "name" : "hobby", "deletion_info" : { "marked_deleted" : "2023-01-04T03:50:26.277031Z", "local_delete_time" : "2023-01-04T03:50:26Z" } },
          { "name" : "hobby", "path" : [ "code" ], "value" : "" }
        ]
      }
    ]
  }
]
```
