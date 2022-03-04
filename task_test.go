package workertask

import (
	"encoding/json"
	"testing"
	"time"

	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/types/known/structpb"
)

func TestTask(t *testing.T) {
	testTable := []struct {
		name          string
		securityEvent *SecurityEvent
		expected      bool
	}{
		{
			name: "ok",
			securityEvent: &SecurityEvent{
				ID:        "001",
				CreatedAt: time.Now(),
				Tags:      map[string]string{"message": "some mess"},
				Port:      8080,
			},
			expected: true,
		},
		{
			name: "wrongID",
			securityEvent: &SecurityEvent{
				ID:        "my id",
				CreatedAt: time.Now(),
				Tags:      map[string]string{"message": "some mess"},
				Port:      8080,
			},
			expected: false,
		},
		{
			name: "no such key",
			securityEvent: &SecurityEvent{
				ID:        "001",
				CreatedAt: time.Now(),
				Tags:      map[string]string{"someTag": "some mess"},
				Port:      8080,
			},
			expected: false,
		},
	}
	for _, test := range testTable {
		t.Run(test.name, func(t *testing.T) {
			b, err := json.Marshal(test.securityEvent)
			if err != nil {
				t.Fatal(err)
			}
			var spb structpb.Struct
			if err := protojson.Unmarshal(b, &spb); err != nil {
				t.Fatal(err)
			}

			val, _, err := program.Eval(map[string]interface{}{"securityEvent": &spb})
			if err != nil {
				t.Fatal(err)
			}

			gotted, ok := val.Value().(bool)
			if !ok {
				t.Fatalf("failed to convert %+v to bool", val)
			}

			if gotted != test.expected {
				t.Errorf("gotted result: %v, want: %v", gotted, test.expected)
			}
		})
	}
}
