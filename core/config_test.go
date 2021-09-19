package core

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDefaultConfigValues(t *testing.T) {
	uadminConfig := NewConfig("configs/test_sqlite.yml")
	assert.Equal(t, uadminConfig.D.Uadmin.Theme, "default")
}
