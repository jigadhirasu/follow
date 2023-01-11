package uuid_test

import (
	"fmt"
	"testing"

	"github.com/jigadhirasu/follow/uuid"
)

func TestUUID(t *testing.T) {

	id := uuid.Gen()

	fmt.Println(id.Hex(), id.String())

	u := uuid.FromHex(id.Hex())
	fmt.Println(u.Hex(), u.String())
}
