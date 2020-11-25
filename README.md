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
