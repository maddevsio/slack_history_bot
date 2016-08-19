package conf

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWeCanGetConfiguration(t *testing.T) {
	cr := NewConfigurator()
	os.Clearenv()
	os.Setenv("SLACK_TOKEN", "Superpoper")
	cr.Run()
	conf := cr.Get()
	assert.Equal(t, conf.SlackToken, "Superpoper")
}
