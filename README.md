# cassandra-learn

Adding notes on cassandra best practices.


## Notes

- CQL does not support aggregation queries like max, min, avg
- CQL does not support group by, having queries
- CQL does not support joins
- CQL does not support OR queries
- CQL does not support wildcard queries
- CQL does not support union, intersection queries
- Table columns cannot be filtered without creating the index
- Greater than (>) and less than (<) are only supported on clustering column.
- CQL is not suitable for analytic purposes because it has so many limitations

Keyspace
- Cassandra keyspace is a namespace that defines how data is replicated on nodes.
- Typically a cluster has one keyspace per application
- When creating a keyspace, using the replication strategy class `SimpleStrategy` for development, in production use `NetworkTopologyStrategy`

## Basics

Get keyspaces info:
```cql
SELECT *
FROM system_schema.keyspaces;
```

Get tables info:

```cql
SELECT *
FROM system_schema.tables
WHERE keyspace_name = 'keyspace name';
```

Get table info:
```cql
SELECT *
FROM system_schema.columns
WHERE keyspace_name = 'keyspace name'
AND table_name = 'table name';
```


## Getting number of partitions


Use `nodetool` to get the `Number of partitions`:
```bash
$ make sh
$ cd /opt/bitnami/cassandra/bin
$ nodetool tablestats hotel.hotels_by_poi;
```

Output:
```bash
Total number of tables: 45
----------------
Keyspace : hotel
        Read Count: 0
        Read Latency: NaN ms
        Write Count: 1
        Write Latency: 3.026 ms
        Pending Flushes: 0
                Table: hotels_by_poi
                SSTable count: 1
                Old SSTable count: 0
                Space used (live): 5086
                Space used (total): 5086
                Space used by snapshots (total): 0
                Off heap memory used (total): 34
                SSTable Compression Ratio: 1.0476190476190477
                Number of partitions (estimate): 1
                Memtable cell count: 0
                Memtable data size: 0
                Memtable off heap memory used: 0
                Memtable switch count: 1
                Local read count: 0
                Local read latency: NaN ms
                Local write count: 1
                Local write latency: NaN ms
                Pending flushes: 0
                Percent repaired: 0.0
                Bytes repaired: 0.000KiB
                Bytes unrepaired: 0.041KiB
                Bytes pending repair: 0.000KiB
                Bloom filter false positives: 0
                Bloom filter false ratio: 0.00000
                Bloom filter space used: 16
                Bloom filter off heap memory used: 8
                Index summary off heap memory used: 18
                Compression metadata off heap memory used: 8
                Compacted partition minimum bytes: 36
                Compacted partition maximum bytes: 42
                Compacted partition mean bytes: 42
                Average live cells per slice (last five minutes): NaN
                Maximum live cells per slice (last five minutes): 0
                Average tombstones per slice (last five minutes): NaN
                Maximum tombstones per slice (last five minutes): 0
                Dropped Mutations: 0
                Droppable tombstone ratio: 0.00000

----------------
```

## Check partition ranges

```bash
$ nodetool ring
```

## Checking how the data is stored in Cassandra


Reference [here](https://stackoverflow.com/questions/35945636/cassandra-cli-list-in-cassandra-3-0). Cassandra used to have the `cassandra-cli` that shows how the data is stored. Before doing this, create a table and insert some data first.

```bash
# Flush data to disk.
$ /opt/bitnami/cassandra/bin/nodetool flush

# Go to the directory where the dump is produced
$ cd /bitnami/cassandra/data/data/hotel/hotels_by_poi-a52429b08bdc11eda5eb23123d1849fe
$ ls
'?'        nb-1-big-CompressionInfo.db   nb-1-big-Digest.crc32   nb-1-big-Index.db        nb-1-big-Summary.db
 backups   nb-1-big-Data.db              nb-1-big-Filter.db      nb-1-big-Statistics.db   nb-1-big-TOC.txt


# Run sstabledump
$ /opt/bitnami/cassandra/tools/bin/sstabledump nb-1-big-Data.db
```

Output:

```json
[
  {
    "partition" : {
      "key" : [ "sunway" ],
      "position" : 0
    },
    "rows" : [
      {
        "type" : "row",
        "position" : 20,
        "clustering" : [ "1" ],
        "liveness_info" : { "tstamp" : "2023-01-04T03:06:25.710079Z" },
        "cells" : [
          { "name" : "name", "value" : "sunway" },
          { "name" : "phone", "value" : "1234" }
        ]
      }
    ]
  }
]
```

## Use Expand ON

To show the data in better order.

## Pagination in Cassandra

There is no `offset` in Casssandra, and querying data from different partitions might require a different strategy. Pagination in cassandra can be done using `pageState`.

- http://www.inanzzz.com/index.php/post/t7fd/cassandra-pagination-example-with-golang
- https://medium.com/@shahsiddharth/cassandra-sorting-and-paging-across-multiple-partitions-for-rest-api-cecf452cbf96
