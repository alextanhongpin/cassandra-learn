package main

import (
	"context"
	"encoding/hex"
	"fmt"

	"github.com/gocql/gocql"
)

type User struct {
	Name  string
	Age   int
	Hobby []string
}

func main() {
	cluster := gocql.NewCluster("127.0.0.1:9042")
	cluster.Keyspace = "greet"
	//cluster.Timeout = 5 * time.Second (default 600ms)
	cluster.Consistency = gocql.LocalOne // Does not work with consistency level Quorum for single node.
	cluster.ProtoVersion = 4
	cluster.Authenticator = gocql.PasswordAuthenticator{
		Username: "cassandra",
		Password: "cassandra",
	}
	session, err := cluster.CreateSession()
	if err != nil {
		panic(err)
	}
	defer session.Close()

	if false {
		if err := createUser(session, User{
			Name:  "alice",
			Age:   123,
			Hobby: []string{"programming", "swimming", "fishing"},
		}); err != nil {
			panic(err)
		}
	}

	users, pageState, err := queryUsers(session, "alice", 1, nil)
	if err != nil {
		panic(err)
	}

	fmt.Printf("1. users: %+v, pageState: %q\n", users, hex.EncodeToString(pageState))
	encodedPageState := hex.EncodeToString(pageState)
	decodedPageState, err := hex.DecodeString(encodedPageState)
	if err != nil {
		panic(err)
	}
	users, pageState, err = queryUsers(session, "alice", 2, decodedPageState)
	if err != nil {
		panic(err)
	}
	fmt.Printf("2. users: %+v, pageState: %q\n", users, hex.EncodeToString(pageState))

	users, nextPageState, nextPartitionKey, err := queryUsersDifferentPartition(session, []string{"alice", "bob"}, 3, nil)
	if err != nil {
		panic(err)
	}
	fmt.Println("3. users", users)
	fmt.Println("pageState", hex.EncodeToString(nextPageState))
	fmt.Println("partitionKey", nextPartitionKey)
}

func createUser(session *gocql.Session, u User) error {
	return session.Query("insert into greet.name (name, age, hobby) values (?, ?, ?)", u.Name, u.Age, u.Hobby).Exec()
}

func queryUsersDifferentPartition(session *gocql.Session, partitionKeys []string, pageSize int, pageState []byte) (users []User, nextPageState []byte, nextPartitionKey string, err error) {
	for _, key := range partitionKeys {
		partitionUsers, partitionPageState, err := queryUsers(session, key, pageSize, pageState)
		if err != nil {
			return nil, nil, "", err
		}
		users = append(users, partitionUsers...)
		nextPartitionKey = key
		nextPageState = partitionPageState
		if len(users) == pageSize {
			return users, nextPageState, nextPartitionKey, nil
		}
	}

	return
}

func queryUsers(session *gocql.Session, partitionKey string, pageSize int, pageState []byte) (users []User, nextPageState []byte, err error) {
	itr := session.
		Query("select name, age, hobby from greet.name where name = ?", partitionKey).
		WithContext(context.Background()).
		PageSize(pageSize).
		PageState(pageState).
		Iter()
	defer itr.Close()

	fmt.Println("numRows", itr.NumRows())
	users = make([]User, 0, itr.NumRows())

	nextPageState = itr.PageState()

	scanner := itr.Scanner()

	for scanner.Next() {
		var u User
		err = scanner.Scan(&u.Name, &u.Age, &u.Hobby)
		if err != nil {
			return
		}

		users = append(users, u)
	}

	if err = scanner.Err(); err != nil {
		return
	}

	return users, nextPageState, err
}
