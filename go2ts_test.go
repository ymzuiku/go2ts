package go2ts_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/ymzuiku/go2ts"
)

type Hello struct {
	Name string `json:"name" validate:"required"`
	Age  string `json:"age"`
	Vip  bool   `json:"vip,omitempty"`
}

type World struct {
	Dog  string `json:"dog" ts_type:"any"`
	Fish string `json:"fish,omitempty" validate:"required"`
}

func TestGo2TsInterface(t *testing.T) {
	code := go2ts.New().Add(Hello{}).Add(World{}).Write("temp/temp1.ts")
	assert.Empty(t, code)
}

func GetWord(hello *Hello) (*World, error) {
	return nil, nil
}

func TestGo2TsFunction(t *testing.T) {
	code := go2ts.New().Add(GetWord).Write("temp/temp2.ts")
	assert.Empty(t, code)
}

func TestGo2TsApi(t *testing.T) {
	code := go2ts.New().AddApi("POST", "/v1/world", GetWord).Write("temp/temp3.ts")
	assert.Empty(t, code)
}
