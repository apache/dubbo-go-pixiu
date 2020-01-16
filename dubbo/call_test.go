package dubbo

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"testing"
)

var (
	testDealGerRespString1 = `{
	"devices": [{
		"activeTime": 1562521341000,
		"bv": [{
			"name": "xxxx"
		}]
	}],
	"limit": 12,
	"offset": 0,
	"pageNo": 1,
	"page_size": 12,
	"total": 17
}`
	exceptGerRespString2 = `{"devices":[{"active_time":1562521341000,"bv":[{"name":"xxxx"}]}],"limit":12,"offset":0,"page_no":1,"page_size":12,"total":17}`
)

func TestDealResp(t *testing.T) {
	func() {
		var testTmp = make(map[string]interface{})
		e1 := json.Unmarshal([]byte(testDealGerRespString1), &testTmp)
		assert.Equal(t, nil, e1)
		out, e2 := dealResp(testTmp, true)
		assert.Equal(t, nil, e2)
		out2, e3 := json.Marshal(out)
		assert.Equal(t, nil, e3)
		assert.Equal(t, exceptGerRespString2, string(out2))
	}()

	func() {
		type valueStu struct {
			V string
		}
		testTmp := [][]map[interface{}]interface{}{
			{
				{
					"xxx": valueStu{},
				},
			},
		}
		out, e2 := dealResp(testTmp, true)
		assert.Equal(t, nil, e2)
		out2, e3 := json.Marshal(out)
		assert.Equal(t, nil, e3)
		assert.Equal(t, "[[{\"xxx\":{\"v\":\"\"}}]]", string(out2))
	}()
	func() {
		type valueStu struct {
			V string
		}
		testTmp := []map[interface{}]interface{}{
			{
				"xxxx": valueStu{},
			},
		}
		out, e2 := dealResp(testTmp, true)
		assert.Equal(t, nil, e2)
		out2, e3 := json.Marshal(out)
		assert.Equal(t, nil, e3)
		assert.Equal(t, "[{\"xxxx\":{\"v\":\"\"}}]", string(out2))
	}()
}
func Test_Map2xx_yy(t *testing.T) {
	testMap := make(map[string]interface{})
	testData := `{"AaAa":"1","BaBa":"1","CaCa":{"BaBa":"2","AaAa":"2","XxYy":{"XxXx":"3","Xx":"3"}}}`
	e := json.Unmarshal([]byte(testData), &testMap)
	assert.Nil(t, nil)
	m := humpToLine(testMap)
	s, e := json.Marshal(m)
	assert.Equal(t, e, nil)
	assert.Equal(t, `{"aa_aa":"1","ba_ba":"1","ca_ca":{"aa_aa":"2","ba_ba":"2","xx_yy":{"xx":"3","xx_xx":"3"}}}`, string(s))
}
