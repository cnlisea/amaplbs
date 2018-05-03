package amaplbs

import (
	"testing"
)

func TestAmapLbsBeeLineDistance(t *testing.T) {
	amapLbs := NewAmapLbs(&Config{
		Key: "",
	})
	distances, err := amapLbs.BeeLineDistance([]string{"114.421432,30.47276"}, "114.410237,30.476385")
	if err != nil {
		t.Fatal(err)
	}

	t.Log("distances:", distances)
}
