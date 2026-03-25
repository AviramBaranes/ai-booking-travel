package pricing

import "math"

// roundToInt rounds a number to the nearest integer.
func RoundToInt[T float32 | float64](f T) int {
	return int(math.Round(float64(f)))
}

// calculateDiscountedPrice calculates the price after applying a discount percentage.
func CalculateDiscountedPrice(priceBeforeDesc float64, discountPercentage int) float64 {
	discountAmount := priceBeforeDesc * float64(discountPercentage) / 100
	return priceBeforeDesc - discountAmount
}

// ApplyMarkup applies a markup percentage to the base price and returns the new price.
func ApplyMarkup(basePrice, markupPct float64) float64 {
	return basePrice * (1 + markupPct/100)
}

// calculateTotalPrice calculates the total price by applying markup and discount to the purchase price.
func CalculateTotalPrice(purchasePrice, markupPercentage, brokerErp float64, btErp, discountPercentage int) int {
	carPriceWithMarkup := ApplyMarkup(purchasePrice, markupPercentage)
	brokerErpWithMarkup := ApplyMarkup(brokerErp, markupPercentage)
	discountedPrice := CalculateDiscountedPrice(carPriceWithMarkup+brokerErpWithMarkup, discountPercentage)
	return RoundToInt(discountedPrice) + btErp
}
