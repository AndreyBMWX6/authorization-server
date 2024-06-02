package scratch

import (
	"encoding/json"
	"fmt"

	"github.com/peterbourgon/mergemap"
)

// swagJoiner glues up several Swagger definitions to one.
type swagJoiner struct {
	result map[string]interface{}
}

// AddDefinition adds another definition to the soup.
func (c *swagJoiner) AddDefinition(buf []byte) error {
	def := map[string]interface{}{}

	err := json.Unmarshal(buf, &def)
	if err != nil {
		return fmt.Errorf("couldn't unmarshal JSON def: %w", err)
	}
	if c.result == nil {
		c.result = def
		return nil
	}
	c.result = mustMergeSwaggers(c.result, def)
	return nil
}

// SumDefinitions returns a (hopefully) valid Swagger definition combined
// from everything that came up .AddDefinition().
func (c *swagJoiner) SumDefinitions() []byte {
	if c.result == nil {
		c.result = map[string]interface{}{}
	}
	ret, err := json.Marshal(c.result)
	if err != nil {
		panic(err)
	}
	return ret
}

func mustMergeSwaggers(dst, src map[string]any) map[string]any {
	return mergemap.Merge(dst, src)
}
