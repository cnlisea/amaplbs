package amaplbs

import "testing"

func TestAmapLbs_GenLocationTrim(t *testing.T) {
	amapLbs := NewAmapLbs(&Config{
		Key: "",
	})

	t.Log(amapLbs.GenLocationTrim("3.532344", 6))
	t.Log(amapLbs.GenLocationTrim("3.53234", 6))
	t.Log(amapLbs.GenLocationTrim("3.53234888", 6))
}

func TestAmapLbs_ReGeoCode(t *testing.T) {
	amapLbs := NewAmapLbs(&Config{
		Key: "",
	})
	info, err := amapLbs.ReGeoCode("116.481488", "39.990464")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(info)
}
