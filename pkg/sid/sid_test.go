package sid

import (
	"fmt"
	"log/slog"
	"strconv"
	"testing"

	"github.com/bwmarrin/snowflake"
)

func TestSid_GenString(t *testing.T) {
	s := NewSid()
	for i := range 100 {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			got, err := s.GenString()
			if err != nil {
				t.Errorf("GenString() error = %v", err)
			}
			slog.Info("Generated SID", "value", got)
		})

	}
}

func TestBwmarrinSnowflake(t *testing.T) {

	// Create a new Node with a Node number of 1
	node, err := snowflake.NewNode(1)
	if err != nil {
		fmt.Println(err)
		return
	}

	// Generate a snowflake ID.
	id := node.Generate()

	// Print out the ID in a few different ways.
	fmt.Printf("Int64  ID: %d\n", id)
	fmt.Printf("String ID: %s\n", id)
	fmt.Printf("Base2  ID: %s\n", id.Base2())
	fmt.Printf("Base64 ID: %s\n", id.Base64())

	// Print out the ID's timestamp
	fmt.Printf("ID Time  : %d\n", id.Time())

	// Print out the ID's node number
	fmt.Printf("ID Node  : %d\n", id.Node())

	// Print out the ID's sequence number
	fmt.Printf("ID Step  : %d\n", id.Step())

	fmt.Printf("ID Base62 : %s\n", IntToBase62(int(id.Int64())))

	// Generate and print, all in one.
	fmt.Printf("ID       : %d\n", node.Generate().Int64())
}
