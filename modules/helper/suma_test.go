package helper

import (
	"testing"
)

var test1 = []struct {
	z      int
	x      int
	result int
}{
	{
		2, 3, 5,
	},
	{
		2, 4, 6,
	},
	{
		5, 5, 10,
	},
}

//error opening Url's file

func TestSumma(t *testing.T) {
	for _, test := range test1 {
		if out := Summa(test.x, test.z); out != test.result {
			t.Errorf("Wanted %v, got %v \n", test.result, out)

		}

	}

}
