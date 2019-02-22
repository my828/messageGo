package main

import (
	"fmt"

	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type Contact struct {
	ID        bson.ObjectId `bson:"_id"` //saved to mongo as `_id`
	Email     string
	FirstName string
	LastName  string
}

func main() {
	sess, err := mgo.Dial("127.0.0.1")
	if err != nil {
		fmt.Printf("error dialing mongo: %v\n", err)
	} else {
		fmt.Printf("connected successfully!\n")
    }
    c := &Contact{
        ID:        bson.NewObjectId(),
        Email:     "test@test.com",
        FirstName: "Test",
        LastName:  "Tester",
    }
    //get a reference to the "contacts" collection
    //in the "demo" database
    coll := sess.DB("demo").C("contacts")

    //insert struct into that collection
    if err := coll.Insert(c); err != nil {
        fmt.Printf("error inserting document: %v\n", err)
    } else {
        fmt.Printf("inserted document with ID %s\n", c.ID.Hex())
    }
    contacts := []*Contact{}

    //find all documents where the lastname property contains "Tester"
    //and decode them into the slice
    //(using same `coll` variable from earlier code snippet)
    coll.Find(bson.M{"lastname": "Tester"}).All(&contacts)

    //iterate over the slice, printing the data to std out
    for _, c := range contacts {
        fmt.Printf("%s: %s, %s, %s\n", c.ID.Hex(), c.Email, c.FirstName, c.LastName)
    }
}
