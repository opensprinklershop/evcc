package vehicle

import (
	"fmt"
	"strings"
	"time"

	"github.com/andig/evcc/api"
	"github.com/andig/evcc/server/config"
)

const interval = 15 * time.Minute

var configTypes []config.Type

func registerConfig(typ, name string, defaults interface{}, rank ...int) {
	typeConfig := config.Type{
		Type:   typ,
		Name:   name,
		Config: defaults,
	}

	if len(rank) > 0 {
		typeConfig.Rank = rank[0]
	}

	configTypes = append(configTypes, typeConfig)
}

func ConfigTypes() []config.Type {
	return configTypes
}

type vehicleRegistry map[string]func(map[string]interface{}) (api.Vehicle, error)

func (r vehicleRegistry) Add(name string, factory func(map[string]interface{}) (api.Vehicle, error)) {
	if _, exists := r[name]; exists {
		panic(fmt.Sprintf("cannot register duplicate vehicle type: %s", name))
	}
	r[name] = factory
}

func (r vehicleRegistry) Get(name string) (func(map[string]interface{}) (api.Vehicle, error), error) {
	factory, exists := r[name]
	if !exists {
		return nil, fmt.Errorf("vehicle type not registered: %s", name)
	}
	return factory, nil
}

var registry vehicleRegistry = make(map[string]func(map[string]interface{}) (api.Vehicle, error))

// NewFromConfig creates vehicle from configuration
func NewFromConfig(typ string, other map[string]interface{}) (v api.Vehicle, err error) {
	factory, err := registry.Get(strings.ToLower(typ))
	if err == nil {
		if v, err = factory(other); err != nil {
			err = fmt.Errorf("cannot create type '%s': %w", typ, err)
		}
	} else {
		err = fmt.Errorf("invalid vehicle type: %s", typ)
	}

	return
}
