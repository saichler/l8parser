package tests

import (
	"testing"
)
import "google.golang.org/protobuf/encoding/protojson"

func TestDevices(t *testing.T) {
	deviceList := Devices()
	_, err := protojson.Marshal(deviceList)
	if err != nil {
		t.Error(err)
		t.Fail()
		return
	}
	//fmt.Println(string(devices))
}
