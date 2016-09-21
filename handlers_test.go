package main

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestPaginationOptions(t *testing.T) {
	req1, _ := http.NewRequest("GET", "/test?pageSize=4&page=3", nil)
	req2, _ := http.NewRequest("GET", "/test?pageSize=1&page=1", nil)
	req3, _ := http.NewRequest("GET", "/test", nil)

	p1 := GetPaginationFromRequest(req1)
	p2 := GetPaginationFromRequest(req2)
	p3 := GetPaginationFromRequest(req3)

	assert.Equal(t, 50, p3.Limit)
	assert.Equal(t, 0, p3.Offset)

	assert.Equal(t, 1, p2.Limit)
	assert.Equal(t, 0, p2.Offset)

	assert.Equal(t, 4, p1.Limit)
	assert.Equal(t, 8, p1.Offset)
}
