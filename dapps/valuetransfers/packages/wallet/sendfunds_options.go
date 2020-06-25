package wallet

import (
	"errors"

	"github.com/iotaledger/goshimmer/dapps/valuetransfers/packages/address"
	"github.com/iotaledger/goshimmer/dapps/valuetransfers/packages/balance"
)

// Option is the type for the optional parameters for the SendFunds call.
type SendFundsOption func(*SendFundsOptions) error

// Destination is an option for the SendFunds call that defines a destination for funds that are supposed to be moved.
func Destination(addr address.Address, amount uint64, optionalColor ...balance.Color) SendFundsOption {
	// determine optional output color
	var outputColor balance.Color
	switch len(optionalColor) {
	case 0:
		outputColor = balance.ColorIOTA
	case 1:
		outputColor = optionalColor[0]
	default:
		return optionError(errors.New("providing more than one output color for the destination of funds is forbidden"))
	}

	// return an error if the amount is less
	if amount == 0 {
		return optionError(errors.New("the amount provided in the destinations needs to be larger than 0"))
	}

	// return Option
	return func(options *SendFundsOptions) error {
		// initialize destinations property
		if options.Destinations == nil {
			options.Destinations = make(map[address.Address]map[balance.Color]uint64)
		}

		// initialize address specific destination
		if _, addressExists := options.Destinations[addr]; !addressExists {
			options.Destinations[addr] = make(map[balance.Color]uint64)
		}

		// initialize color specific destination
		if _, colorExists := options.Destinations[addr][outputColor]; !colorExists {
			options.Destinations[addr][outputColor] = 0
		}

		// increase amount
		options.Destinations[addr][outputColor] += amount

		return nil
	}
}

// Remainder is an option for the SendsFunds call that allows us to specify the remainder address that is
// supposed to be used in the corresponding transaction.
func Remainder(addr Address) SendFundsOption {
	return func(options *SendFundsOptions) error {
		options.RemainderAddress = addr

		return nil
	}
}

// Options is a struct that is used to aggregate the optional parameters provided in the SendFunds call.
type SendFundsOptions struct {
	Destinations     map[address.Address]map[balance.Color]uint64
	RemainderAddress Address
}

func BuildOptions(options ...SendFundsOption) (result *SendFundsOptions, err error) {
	// create options to collect the arguments provided
	result = &SendFundsOptions{}

	// apply arguments to our options
	for _, option := range options {
		if err = option(result); err != nil {
			return
		}
	}

	// sanitize parameters
	if len(result.Destinations) == 0 {
		err = errors.New("you need to provide at least one Destination for a valid transfer to be issued")

		return
	}

	return
}

// optionError is a utility function that returns a Option that returns the error provided in the
// argument.
func optionError(err error) SendFundsOption {
	return func(options *SendFundsOptions) error {
		return err
	}
}