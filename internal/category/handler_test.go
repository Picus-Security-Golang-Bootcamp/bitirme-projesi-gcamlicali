package category

import (
	"github.com/gin-gonic/gin"
	"testing"
)

func Test_categoryHandler_addSingle(t *testing.T) {
	type fields struct {
		service Service
	}
	type args struct {
		c *gin.Context
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &categoryHandler{
				service: tt.fields.service,
			}
			h.addSingle(tt.args.c)
		})
	}
}
