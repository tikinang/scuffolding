package skelet

import (
	"strings"

	"github.com/fatih/structs"
	"github.com/jeremywohl/flatten"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.uber.org/dig"
)

const emptyString = ""

type FlavorProvider interface {
	GetFlavors() []any
}

// TODO(mpavlicek): return default T
func AssembleSkelet[T any](
	cmd *cobra.Command,
	name string,
	flavors FlavorProvider,
	invokee T,
	bones ...any,
) (T, error) {
	ctn, err := registerProviders(cmd, name, flavors, bones...)
	if err != nil {
		return invokee, err
	}

	if err := ctn.Invoke(func(x T) { invokee = x }); err != nil {
		return invokee, errors.Wrap(err, "invoke invokee")
	}

	return invokee, nil
}

func AssembleRunner[T any](
	cmd *cobra.Command,
	name string,
	flavors FlavorProvider,
	skelet T,
	bones ...any,
) (*Runner, error) {
	ctn, err := registerProviders(cmd, name, flavors, bones...)
	if err != nil {
		return nil, err
	}

	if err := ctn.Provide(NewRunner); err != nil {
		return nil, errors.Wrap(err, "provide runner")
	}

	if err := ctn.Invoke(func(x T) { skelet = x }); err != nil {
		return nil, errors.Wrap(err, "invoke skelet")
	}

	var runner *Runner
	if err := ctn.Invoke(func(x *Runner) { runner = x }); err != nil {
		return nil, errors.Wrap(err, "invoke runner")
	}

	return runner, nil
}

func registerProviders(
	cmd *cobra.Command,
	name string,
	flavors FlavorProvider,
	bones ...any,
) (*dig.Container, error) {
	v := viper.New()
	v.AutomaticEnv()
	v.SetEnvPrefix(name)
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	flavorMap := structs.Map(flavors)
	flat, err := flatten.Flatten(flavorMap, emptyString, flatten.DotStyle)
	if err != nil {
		return nil, errors.Wrap(err, "flatten flavor")
	}

	// Bind each conf fields to environment vars.
	for key := range flat {
		if err := v.BindEnv(key); err != nil {
			return nil, errors.Wrapf(err, "bind env var: %s", key)
		}
		// TODO(mpavlicek): BindPFlag
	}

	if err := viper.BindPFlags(cmd.PersistentFlags()); err != nil {
		return nil, errors.Wrap(err, "bind persistent flags from cobra to viper")
	}

	if err := v.Unmarshal(&flavors); err != nil {
		return nil, errors.Wrap(err, "unmarshal flavors")
	}

	// Construct app with dig.
	ctn := dig.New(dig.DeferAcyclicVerification())
	for _, bone := range bones {
		if err := ctn.Provide(bone); err != nil {
			return nil, errors.Wrap(err, "provide bone")
		}
	}

	for _, flavor := range flavors.GetFlavors() {
		if err := ctn.Provide(flavor); err != nil {
			return nil, errors.Wrap(err, "provide flavor")
		}
	}

	return ctn, nil
}
