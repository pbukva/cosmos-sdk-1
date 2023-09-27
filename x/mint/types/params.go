package types

import (
	"errors"
	"fmt"
	"strings"

	"cosmossdk.io/math"
	"sigs.k8s.io/yaml"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// NewParams returns Params instance with the given values.
func NewParams(mintDenom string, inflationRate sdk.Dec, blocksPerYear uint64) Params {
	return Params{
		MintDenom:     mintDenom,
		InflationRate: inflationRate,
		BlocksPerYear: blocksPerYear,
	}
}

// DefaultParams returns default x/mint module parameters.
func DefaultParams() Params {
	return Params{
		MintDenom:     sdk.DefaultBondDenom,
		InflationRate: sdk.NewDecWithPrec(3, 2),
		BlocksPerYear: uint64(60 * 60 * 8766 / 5), // assuming 5 second block times
	}
}

// Validate does the sanity check on the params.
func (p Params) Validate() error {
	if err := validateMintDenom(p.MintDenom); err != nil {
		return err
	}
	if err := validateInflationRate(p.InflationRate); err != nil {
		return err
	}
	if err := validateBlocksPerYear(p.BlocksPerYear); err != nil {
		return err
	}
	return nil
}

// String implements the Stringer interface.
func (p Params) String() string {
	out, _ := yaml.Marshal(p)
	return string(out)
}

func validateMintDenom(i interface{}) error {
	v, ok := i.(string)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if strings.TrimSpace(v) == "" {
		return errors.New("mint denom cannot be blank")
	}
	if err := sdk.ValidateDenom(v); err != nil {
		return err
	}

	return nil
}

func validateInflationRate(i interface{}) error {
	v, ok := i.(sdk.Dec)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v.IsNil() {
		return fmt.Errorf("inflation rate cannot be nil: %s", v)
	}
	if v.IsNegative() {
		return fmt.Errorf("inflation rate cannot be negative: %s", v)
	}
	if v.GT(math.LegacyOneDec()) {
		return fmt.Errorf("inflation rate too large: %s", v)
	}

	return nil
}

func validateBlocksPerYear(i interface{}) error {
	v, ok := i.(uint64)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v == 0 {
		return fmt.Errorf("blocks per year must be positive: %d", v)
	}

	return nil
}
