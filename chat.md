## Chat

- The rowkey of the table is the composite key made out of `id` and `created_at`.
- Every row are uniquely identifiable by (or partitioned by) `id + created_at` and ordered by `msg_sent_to_server` timestamp in descending order.

```CQL
CREATE KEYSPACE chat WITH replication = {'class': 'SimpleStrategy', 'replication_factor': 3};
CREATE TABLE chat.message(
  id uuid,
	members list<text>,
	sent_by ascii,
	created_at date,
	msg map<text, text>, // {'msgtype': 'private|normal|disappearing', 'msgdata': 'hi'}
	msg_sent_to_server timeuuid,
	msg_sent_to_recv map<text,timestamp>, // Key-value pair of the members and the time it was sent to them. {alice: toTimestamp(now()), bob: toTimestamp(now())}
	msg_read_by_recv map<text,timestamp>, // Key-value pair of the members and the time it was read by them.
	PRIMARY KEY ((id, created_at), msg_sent_to_server)
) WITH CLUSTERING ORDER BY (msg_sent_to_server DESC);
```

Insert:
```cql
INSERT INTO chat.message JSON '{
	"id": uuid(), // Not possible.
	"members": ["alice", "bob"],
	"sent_by": "alice",
	"created_at": "2020-01-01",
	"msg": {
		"alice": "hi"
	}
}';

INSERT INTO chat.message (id, members, sent_by, created_at, msg, msg_sent_to_server, msg_sent_to_recv)
VALUES (
	uuid(),
	['alice', 'bob'],
	'bob',
	toDate(now()),
	{'msgtype': 'message', 'msgdata': 'hi how are you'},
	now(),
	{'alice': toTimestamp(now()), 'bob': toTimestamp(now())}
);
```

Select:

```cql
SELECT *
FROM chat.message;
```

## References
- http://bytecontinnum.com/2016/09/wide-row-data-modelling-apache-cassandra/


Unanswered:
- how to query all messages by user?
- what if the user's name change?
- how is this better than relational database?
