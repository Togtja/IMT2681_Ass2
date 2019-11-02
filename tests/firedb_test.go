package tests

import (
	"RESTGvkGitLab/firedb"
	"RESTGvkGitLab/globals"
	"testing"
)

func TestFiredb(t *testing.T) {
	firedb.InitDataBase()
	//Doc("0") is reserved for status checks
	if firedb.Client == nil {
		t.Errorf("Could not get client")
		return

	}
	_, err := firedb.Client.Collection(globals.WebhookF).Doc("0").Get(firedb.Ctx)
	if err != nil {
		//Can not get to server for unknown reason
		//Gives a Service Unavailable error
		t.Errorf("Unable to get collection from db")
	}
}
