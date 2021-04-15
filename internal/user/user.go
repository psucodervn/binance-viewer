package user

import (
	"github.com/adshao/go-binance/v2/futures"
	"github.com/rs/zerolog/log"

	"copytrader/internal/binance"
	"copytrader/internal/model"
	"copytrader/internal/util"
)

func GetUserClients(u model.User) []*binance.AccountClient {
	clients := make([]*binance.AccountClient, len(u.Accounts))
	i := 0
	for _, acc := range u.Accounts {
		clients[i] = binance.NewAccountClient(acc.ApiKey, acc.SecretKey)
		i++
	}
	return clients
}

func GetUserClient(acc model.Account) *binance.AccountClient {
	return binance.NewAccountClient(acc.ApiKey, acc.SecretKey)
}

func TotalUserIncome(incomes []*futures.IncomeHistory) float64 {
	res := 0.0
	for _, tr := range incomes {
		v := util.ParseFloat(tr.Income)
		switch tr.IncomeType {
		case binance.TypeRealizedPNL:
		case binance.TypeCommission:
		case binance.TypeTransfer:
		case binance.TypeCommissionRebate:
		case binance.TypeFundingFee:
		case binance.TypeReferralKickback:
		case binance.TypeInsuranceClear:
		default:
			log.Warn().Float64("value", v).Msg("unknown type " + tr.IncomeType)
			continue
		}
		res += v
	}
	return res
}
