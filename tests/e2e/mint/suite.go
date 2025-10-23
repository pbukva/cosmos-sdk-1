package mint

import (
	"fmt"
	tmcli "github.com/cometbft/cometbft/libs/cli"
	"strings"

	"github.com/stretchr/testify/suite"

	"github.com/cosmos/cosmos-sdk/client/flags"
	clitestutil "github.com/cosmos/cosmos-sdk/testutil/cli"
	"github.com/cosmos/cosmos-sdk/testutil/network"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/mint/client/cli"
	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"
)

type E2ETestSuite struct {
	suite.Suite

	cfg     network.Config
	network *network.Network
}

func NewE2ETestSuite(cfg network.Config) *E2ETestSuite {
	return &E2ETestSuite{cfg: cfg}
}

func (s *E2ETestSuite) SetupSuite() {
	s.T().Log("setting up e2e test suite")

	genesisState := s.cfg.GenesisState

	var mintData minttypes.GenesisState
	s.Require().NoError(s.cfg.Codec.UnmarshalJSON(genesisState[minttypes.ModuleName], &mintData))

	inflation := sdk.MustNewDecFromStr("1.0")
	mintData.Minter.Inflation = inflation
	mintData.Minter.MunicipalInflation = []*minttypes.MunicipalInflationPair{
		{Denom: "denom1", Inflation: minttypes.NewMunicipalInflation("cosmos12kdu2sy0zcmz84qymyj6zcfvwss3a703xgpczm", sdk.NewDecWithPrec(234, 4))},
		{Denom: "denom0", Inflation: minttypes.NewMunicipalInflation("cosmos1d9pzg5542spe4anjgu2zmk7wxhgh04ysn2phpq", sdk.NewDecWithPrec(123, 2))},
		{Denom: "denom3", Inflation: minttypes.NewMunicipalInflation("cosmos1ck73rpk6eqxtla4rv7rspsq7apl3740rgjfte4", sdk.NewDecWithPrec(456, 3))},
		{Denom: "denom2", Inflation: minttypes.NewMunicipalInflation("cosmos1ury8qn5w7m3xkl9pdd9ehazd2c9urx7qht2jly", sdk.NewDecWithPrec(345, 3))},
	}

	mintData.Params.InflationRate = sdk.MustNewDecFromStr("0.012345")

	mintDataBz, err := s.cfg.Codec.MarshalJSON(&mintData)
	s.Require().NoError(err)
	genesisState[minttypes.ModuleName] = mintDataBz
	s.cfg.GenesisState = genesisState

	s.network, err = network.New(s.T(), s.T().TempDir(), s.cfg)
	s.Require().NoError(err)

	s.Require().NoError(s.network.WaitForNextBlock())
}

func (s *E2ETestSuite) TearDownSuite() {
	s.T().Log("tearing down e2e test suite")
	s.network.Cleanup()
}

func (s *E2ETestSuite) TestGetCmdQueryParams() {
	val := s.network.Validators[0]

	testCases := []struct {
		name           string
		args           []string
		expectedOutput string
	}{
		{
			"json output",
			[]string{fmt.Sprintf("--%s=1", flags.FlagHeight), fmt.Sprintf("--%s=json", tmcli.OutputFlag)},
			`{"mint_denom":"stake","inflation_rate":"0.012345000000000000","blocks_per_year":"6311520"}`,
		},
		{
			"text output",
			[]string{fmt.Sprintf("--%s=1", flags.FlagHeight), fmt.Sprintf("--%s=text", tmcli.OutputFlag)},
			`blocks_per_year: "6311520"
inflation_rate: "0.012345000000000000"
mint_denom: stake`,
		},
	}

	for _, tc := range testCases {
		tc := tc

		s.Run(tc.name, func() {
			cmd := cli.GetCmdQueryParams()
			clientCtx := val.ClientCtx

			out, err := clitestutil.ExecTestCLICmd(clientCtx, cmd, tc.args)
			s.Require().NoError(err)
			s.Require().Equal(tc.expectedOutput, strings.TrimSpace(out.String()))
		})
	}
}

func (s *E2ETestSuite) TestGetCmdQueryMunicipalInflation() {
	val := s.network.Validators[0]

	testCases := []struct {
		name           string
		args           []string
		expectedOutput string
	}{
		{
			"full - json output",
			[]string{fmt.Sprintf("--%s=1", flags.FlagHeight), fmt.Sprintf("--%s=json", tmcli.OutputFlag)},
			`{"inflations":[{"denom":"denom1","inflation":{"target_address":"cosmos12kdu2sy0zcmz84qymyj6zcfvwss3a703xgpczm","value":"0.023400000000000000"}},{"denom":"denom0","inflation":{"target_address":"cosmos1d9pzg5542spe4anjgu2zmk7wxhgh04ysn2phpq","value":"1.230000000000000000"}},{"denom":"denom3","inflation":{"target_address":"cosmos1ck73rpk6eqxtla4rv7rspsq7apl3740rgjfte4","value":"0.456000000000000000"}},{"denom":"denom2","inflation":{"target_address":"cosmos1ury8qn5w7m3xkl9pdd9ehazd2c9urx7qht2jly","value":"0.345000000000000000"}}]}`,
		},
		{
			"full - text output",
			[]string{fmt.Sprintf("--%s=1", flags.FlagHeight), fmt.Sprintf("--%s=text", tmcli.OutputFlag)},
			`inflations:
- denom: denom1
  inflation:
    target_address: cosmos12kdu2sy0zcmz84qymyj6zcfvwss3a703xgpczm
    value: "0.023400000000000000"
- denom: denom0
  inflation:
    target_address: cosmos1d9pzg5542spe4anjgu2zmk7wxhgh04ysn2phpq
    value: "1.230000000000000000"
- denom: denom3
  inflation:
    target_address: cosmos1ck73rpk6eqxtla4rv7rspsq7apl3740rgjfte4
    value: "0.456000000000000000"
- denom: denom2
  inflation:
    target_address: cosmos1ury8qn5w7m3xkl9pdd9ehazd2c9urx7qht2jly
    value: "0.345000000000000000"`,
		},
		{
			"selected denom - json output",
			[]string{"denom3", fmt.Sprintf("--%s=1", flags.FlagHeight), fmt.Sprintf("--%s=json", tmcli.OutputFlag)},
			`{"inflations":[{"denom":"denom3","inflation":{"target_address":"cosmos1ck73rpk6eqxtla4rv7rspsq7apl3740rgjfte4","value":"0.456000000000000000"}}]}`,
		},
		{
			"selected denom - text output",
			[]string{"denom3", fmt.Sprintf("--%s=1", flags.FlagHeight), fmt.Sprintf("--%s=text", tmcli.OutputFlag)},
			`inflations:
- denom: denom3
  inflation:
    target_address: cosmos1ck73rpk6eqxtla4rv7rspsq7apl3740rgjfte4
    value: "0.456000000000000000"`,
		},
	}

	for _, tc := range testCases {
		tc := tc

		s.Run(tc.name, func() {
			cmd := cli.GetCmdQueryMunicipalInflation()
			clientCtx := val.ClientCtx

			out, err := clitestutil.ExecTestCLICmd(clientCtx, cmd, tc.args)
			s.Require().NoError(err)
			s.Require().Equal(tc.expectedOutput, strings.TrimSpace(out.String()))
		})
	}
}

func (s *E2ETestSuite) TestGetCmdQueryInflation() {
	val := s.network.Validators[0]

	testCases := []struct {
		name           string
		args           []string
		expectedOutput string
	}{
		{
			"json output",
			[]string{fmt.Sprintf("--%s=1", flags.FlagHeight), fmt.Sprintf("--%s=json", flags.FlagOutput)},
			`0.012345000000000000`,
		},
		{
			"text output",
			[]string{fmt.Sprintf("--%s=1", flags.FlagHeight), fmt.Sprintf("--%s=text", flags.FlagOutput)},
			`0.012345000000000000`,
		},
	}

	for _, tc := range testCases {
		tc := tc

		s.Run(tc.name, func() {
			cmd := cli.GetCmdQueryInflation()
			clientCtx := val.ClientCtx

			out, err := clitestutil.ExecTestCLICmd(clientCtx, cmd, tc.args)
			s.Require().NoError(err)
			s.Require().Equal(tc.expectedOutput, strings.TrimSpace(out.String()))
		})
	}
}

func (s *E2ETestSuite) TestGetCmdQueryAnnualProvisions() {
	val := s.network.Validators[0]

	testCases := []struct {
		name           string
		args           []string
		expectedOutput string
	}{
		{
			"json output",
			[]string{fmt.Sprintf("--%s=1", flags.FlagHeight), fmt.Sprintf("--%s=json", flags.FlagOutput)},
			`6172500.000000000000000000`,
		},
		{
			"text output",
			[]string{fmt.Sprintf("--%s=1", flags.FlagHeight), fmt.Sprintf("--%s=text", flags.FlagOutput)},
			`6172500.000000000000000000`,
		},
	}

	for _, tc := range testCases {
		tc := tc

		s.Run(tc.name, func() {
			cmd := cli.GetCmdQueryAnnualProvisions()
			clientCtx := val.ClientCtx

			out, err := clitestutil.ExecTestCLICmd(clientCtx, cmd, tc.args)
			s.Require().NoError(err)
			s.Require().Equal(tc.expectedOutput, strings.TrimSpace(out.String()))
		})
	}
}
