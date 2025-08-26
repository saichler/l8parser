package tests

import (
	"strings"
	"testing"

	"github.com/saichler/l8parser/go/parser/boot"
	"github.com/saichler/l8pollaris/go/types"
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

func TestPolling(t *testing.T) {
	m := CheckPollaris()
	if len(m) > 0 {
		t.Fail()
	}
}

func CheckPollaris() map[string]string {
	plrs := boot.GetAllPolarisModels()
	result := make(map[string]string)
	for _, p := range plrs {
		checkPollaris(p, result)
	}
	return result
}

func checkPollaris(p *types.Pollaris, invalid map[string]string) {
	for _, poll := range p.Polling {
		for _, attr := range poll.Attributes {
			for _, rule := range attr.Rules {
				if rule.Name == "Set" {
					from, ok := rule.Params["from"]
					if ok && !strings.HasPrefix(from.Value, poll.What) {
						invalid[attr.PropertyId] = poll.Name
					}
				}
			}
		}
	}
}
