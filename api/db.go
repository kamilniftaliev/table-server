package main

import (
	"context"
	"os"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var host, _ = os.Hostname()

var IsLocalhost = host != "table"

// RequestContext - Current request's context for handling cancelation
var RequestContext = context.TODO()

func getMongoURI(isLocalhost bool) string {
	if isLocalhost {
		return "localhost:27017"
	} else {
		return "guest:kamil123@localhost:5355"
	}
}

// Client - MongoDB client for Database actions
var Client, _ = mongo.Connect(RequestContext, options.Client().ApplyURI("mongodb://"+getMongoURI(IsLocalhost)+"/table"))

// DB - Selected database for accessing data
var DB = Client.Database("table")
