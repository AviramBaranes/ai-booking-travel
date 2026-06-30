package currency

// import (
// 	"context"

// 	"encore.app/internal/api_errors"
// 	dbadapters "encore.app/internal/db_adapters"
// 	"encore.app/internal/icount"
// 	"encore.app/services/accounts"
// 	"encore.app/services/booking/db"
// 	"encore.dev/rlog"
// )

// // UpdateCurrenciesRates fetches the latest currency rates from iCount and updates the database.
// func (s *CurrencyService) UpdateCurrenciesRates(ctx context.Context, cid string, user string) error {
// 	i := icount.NewIcount()
// 	res, err := i.FetchCurrencies()
// 	if err != nil {
// 		rlog.Error("failed to fetch currencies from iCount", "error", err)
// 		return api_errors.ErrInternalError
// 	}
// 	err = s.query.UpsertCurrencies(ctx, createUpsertParams(res))
// 	if err != nil {
// 		rlog.Error("failed to upsert currencies", "error", err)
// 		return api_errors.ErrInternalError
// 	}

// 	for isoName, rate := range res.Rates {
// 		err := accounts.CurrenciesRates.Set(ctx, isoName, rate)
// 		if err != nil {
// 			rlog.Error("failed to set currency rate in cache", "error", err, "currency_iso_name", isoName)
// 			return api_errors.ErrInternalError
// 		}
// 	}

// 	return nil
// }

// // createUpsertParams converts the iCount response to the parameters needed for the UpsertCurrencies query.
// func createUpsertParams(res *icount.GetCurrenciesRatesResponse) db.UpsertCurrenciesParams {
// 	currencyCodes := make([]string, 0, len(res.Rates))
// 	currencyIsoNames := make([]string, 0, len(res.Rates))
// 	rates := make([]dbadapters.Numeric, 0, len(res.Rates))

// 	for isoName, rate := range res.Rates {
// 		currencyCodes = append(currencyCodes, currencyIsoNameToCode(isoName))
// 		currencyIsoNames = append(currencyIsoNames, isoName)
// 		rates = append(rates, dbadapters.NumericFromFloat64(rate))
// 	}

// 	return db.UpsertCurrenciesParams{
// 		CurrencyCodes:    currencyCodes,
// 		CurrencyIsoNames: currencyIsoNames,
// 		Rates:            rates,
// 	}
// }

// // currencyIsoNameToCode maps currency ISO names to their corresponding symbols. If the ISO name is not recognized, it returns the ISO name itself.
// func currencyIsoNameToCode(isoName string) string {
// 	switch isoName {
// 	case "EUR":
// 		return "€"
// 	case "USD":
// 		return "$"
// 	case "GBP":
// 		return "£"
// 	case "ILS":
// 		return "₪"
// 	default:
// 		return isoName
// 	}
// }
