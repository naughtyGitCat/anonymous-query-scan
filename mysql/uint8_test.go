/**
 * Created by zhangruizhi on 2024/7/22
 */

package mysql

import (
	"encoding/json"
	"testing"
)

// TestUint8Marshal
// []uint8{1, 2, 3, 4, 5, 6} and []int8{1, 2, 3, 4, 5, 6} after marshal given different result
func TestUint8Marshal(t *testing.T) {
	var x = []uint8{1, 2, 3, 4, 5, 6}
	var y = []int8{1, 2, 3, 4, 5, 6}
	uint8Bytes, err := json.Marshal(x)
	if err != nil {
		panic(err)
	}
	int8Bytes, err := json.Marshal(y)
	if err != nil {
		panic(err)
	}
	if string(uint8Bytes) != "AQIDBAUG" {
		t.Failed()
	}
	if string(int8Bytes) != "[1,2,3,4,5,6]" {
		t.Failed()
	}
}
