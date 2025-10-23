package mint

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"cosmossdk.io/math"
	"github.com/cosmos/cosmos-sdk/testutil"
	grpctypes "github.com/cosmos/cosmos-sdk/types/grpc"

	"github.com/cosmos/gogoproto/proto"

	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"
)

func (s *E2ETestSuite) TestQueryGRPC() {
	val := s.network.Validators[0]
	baseURL := val.APIAddress
	testCases := []struct {
		name     string
		url      string
		headers  map[string]string
		respType proto.Message
		expected proto.Message
	}{
		{
			"gRPC request params",
			fmt.Sprintf("%s/cosmos/mint/v1beta1/params", baseURL),
			map[string]string{},
			&minttypes.QueryParamsResponse{},
			&minttypes.QueryParamsResponse{
				Params: minttypes.NewParams("stake", math.LegacyMustNewDecFromStr("0.012345000000000000"), 60*60*8766/5),
			},
		},
		{
			"gRPC request inflation",
			fmt.Sprintf("%s/cosmos/mint/v1beta1/inflation", baseURL),
			map[string]string{},
			&minttypes.QueryInflationResponse{},
			&minttypes.QueryInflationResponse{
				Inflation: math.LegacyMustNewDecFromStr("0.012345000000000000"),
			},
		},
		{
			"gRPC request municipal_inflation",
			fmt.Sprintf("%s/cosmos/mint/v1beta1/municipal_inflation", baseURL),
			map[string]string{},
			&minttypes.QueryMunicipalInflationResponse{},
			&minttypes.QueryMunicipalInflationResponse{
				Inflations: []*minttypes.MunicipalInflationPair{
					{Denom: "denom1", Inflation: minttypes.NewMunicipalInflation("cosmos12kdu2sy0zcmz84qymyj6zcfvwss3a703xgpczm", sdk.NewDecWithPrec(234, 4))},
					{Denom: "denom0", Inflation: minttypes.NewMunicipalInflation("cosmos1d9pzg5542spe4anjgu2zmk7wxhgh04ysn2phpq", sdk.NewDecWithPrec(123, 2))},
					{Denom: "denom3", Inflation: minttypes.NewMunicipalInflation("cosmos1ck73rpk6eqxtla4rv7rspsq7apl3740rgjfte4", sdk.NewDecWithPrec(456, 3))},
					{Denom: "denom2", Inflation: minttypes.NewMunicipalInflation("cosmos1ury8qn5w7m3xkl9pdd9ehazd2c9urx7qht2jly", sdk.NewDecWithPrec(345, 3))},
				},
			},
		},
		{
			"gRPC request annual provisions",
			fmt.Sprintf("%s/cosmos/mint/v1beta1/annual_provisions", baseURL),
			map[string]string{
				grpctypes.GRPCBlockHeightHeader: "1",
			},
			&minttypes.QueryAnnualProvisionsResponse{},
			&minttypes.QueryAnnualProvisionsResponse{
				AnnualProvisions: math.LegacyNewDec(6172500),
			},
		},
	}
	for _, tc := range testCases {
		resp, err := testutil.GetRequestWithHeaders(tc.url, tc.headers)
		s.Run(tc.name, func() {
			s.Require().NoError(err)
			s.Require().NoError(val.ClientCtx.Codec.UnmarshalJSON(resp, tc.respType))
			s.Require().Equal(tc.expected.String(), tc.respType.String())
		})
	}
}
