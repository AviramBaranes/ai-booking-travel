package broker

// hertzGRTerms are Autohellas' terms and conditions for Greece, per brand (OTA vendor code),
// generated from their "Autohellas Greece 2026 Terms & Conditions" workbook, wording unchanged.
// Regenerate them from the new workbook when Autohellas updates its terms.
var hertzGRTerms = map[string][]TermsAndConditionsItem{
	"ZE": {
		{
			Title:       `Driver's Qualitifications`,
			HtmlContent: `<p>A valid driving license held for at least one (1) year is required. The driver must hold a valid credit or debit card accepted by Hertz for all necessary pre-authorization amounts as well as pre-payments and settlement of final rental charges. Accepted credit/debit cards are Visa, Mastercard, American Express, Diners and Unionpay. Maestro debit cards are also accepted but only for payments and not for guarantee provision.</p>`,
		},
		{
			Title:       `Age Limits & Young Driver Surcharge`,
			HtmlContent: `<ul><li>From 19 to 20 years old – groups A, B, J, Z, AE, ZE, V2: €10 per day + VAT</li><li>From 19 to 20 years old – other categories not available</li><li>From 21 to 24 years old – groups A, B, Z, J, V2, AE: No additional charge</li><li>From 21 to 24 years old – groups C, D, F, K, M, P, T, R, BE: €10 per day + VAT</li><li>From 21 to 24 years old – groups E, EM, N, U, V, DE, NM, O, V3, V4, Y, EE, FE: €15 per day + VAT</li><li>From 21 to 24 years old – other categories not available</li><li>25 years and older – all car categories: No additional charge</li></ul>`,
		},
		{
			Title:       `Driving Restrictions`,
			HtmlContent: `<p>Vehicle may not be driven across any of the Greek borders without prior written approval by Hertz Greece. Failure to comply with these restrictions makes void all insurance policies and customer fully liable and responsible in case of a theft, accident, damages to the vehicle and any resulting consequences.</p>`,
		},
		{
			Title:       `Grace Period`,
			HtmlContent: `<p>Detailed instructions on car pick-up will be provided for each Hertz location.</p><p>At the end of the agreed rental period Hertz allows a grace period of 29 minutes free of charge after which the customer will be charged with a minimum of one full day rental.</p>`,
		},
		{
			Title:       `Early Return / Rental Extension`,
			HtmlContent: `<p>No refund for any unused voucher value/rental days. Ongoing rental extension should be arranged with Hertz at least 48 hours before the intended/reserved rental expiry, to be charged as per valid retail tariff.</p>`,
		},
		{
			Title:       `Late Cancellation / Non-Show Policy`,
			HtmlContent: `<p>Hertz has the right, in case a customer fails either to show up or to cancel at least 48 hours prior to pick-up, to charge the client a Late Cancellation Fee of €45. If the customer fails to cancel their reservation prior to the pick-up time and does not collect the vehicle on the pick-up date, or fails to comply with the Rental Terms and Conditions (see Hertz Rental Qualifications &amp; Requirements), Hertz reserves the right to charge a No-Show Fee of €95 (billed in destination currency). A Non-Show is effective if the customer hasn't collected the vehicle within one (1) hour of the booked time/date, unless there is a travel delay (given that a flight number and/or travel details have been provided prior to the car pick-up).</p>`,
		},
		{
			Title:       `Fleet Availability`,
			HtmlContent: `<p>Hertz reserves the right to provide a different car model to that originally reserved of an equivalent or higher model type.</p>`,
		},
		{
			Title:       `Accident & Theft Coverages`,
			HtmlContent: `<p><strong>Collision Damage Waiver (CDW)</strong></p><p>Obtaining Collission Damage Waiver (CDW) will reduced renter's liability for damage to the rented vehicle, as per rental's agreement terms &amp; conditions, to the amounts shown in the table below.  The following cases are not covered: Damage to the underside of the car, the wheels &amp; tyres, broken car glass or loss of the vehicle keys and documents.</p><p><strong>Reduced Excess Coverage (SCDW)</strong></p><p>Super Collission Damage Waiver (SCDW) can be obtained only after purchasing Collission Damage Waiver (CDW) and will reduce renter's liability further, to the amounts shown in the table below.</p><p>The following cases are not covered: Damage to the underside of the car, the wheels &amp; tyres, broken car glass or loss of the vehicle keys and documents.</p><p><strong>Super Cover (SC)</strong></p><p>Τhe renter can be covered (with zero excess) in case of total theft or damage to the Hertz vehicle, provided that the renter has previously accepted the relevant coverage charges for CDW, SCDW, and TP, or has a reservation inclusive of these coverages. Moreover, Super Cover covers any damages caused on the vehicle's windshield and tyres and includes Extended Road Assistance. Damages caused to the underside and wheels of the vehicle are not covered by any insurance. Super Cover is not available for rentals made by drivers aged 19 and 20 years old.</p><p><strong>Tyre &amp; Glass Protection (TGI)</strong></p><p>Tyre and Glass Coverage provides cover against consequences of damages on tires, windscreen, side windows and back window. If TGI is not accepted, the client will be held liable for the full value of damage on tires, windscreen, side windows and back window.</p><p><strong>Theft Protection (TPC)</strong></p><p>Theft Protection (TPC) will reduce renter's liability to the amounts shown in the table below.</p>`,
		},
		{
			Title:       `Coverage Charges & Excess per Group`,
			HtmlContent: `<ul><li>Groups A, B, J, Z, AE, ZE: CDW included, excess €650; SCDW €11/day (1-3 days), €8.5/day (4-16 days), excess €160; Super Cover €12.5/day (1-3 days), €9.5/day (4-16 days), zero excess; Tyre &amp; Glass €5.5/day (max 17 days); Theft Protection included, theft excess €650</li><li>Groups C, D, F, K, M, P, R, T, BE: CDW included, excess €800; SCDW €13.5/day (1-3 days), €10.5/day (4-16 days), excess €200; Super Cover €15/day (1-3 days), €11.5/day (4-16 days), zero excess; Theft Protection included, theft excess €800</li><li>Groups E, EM, N, NM, O, U, V, Y, YL, 1, 5, 8, DE, FE, EE: CDW included, excess €1000; SCDW €20.5/day (1-3 days), €16/day (4-16 days), excess €300; Super Cover €23/day (1-3 days), €18/day (4-16 days), zero excess; Theft Protection included, theft excess €1000</li><li>Groups 9, SE, SY, 9F,V6: CDW included, excess €1300; SCDW €40/day (1-3 days), €30.5/day (4-16 days), excess €400; Super Cover €30/day (1-3 days), €23/day (4-16 days), zero excess; Theft Protection included, theft excess €1300</li><li>Groups 3, YP, 9L: CDW included, excess €3000; SCDW €68/day (1-3 days), €52.5/day (4-16 days), excess €800; Super Cover €61.5/day (1-3 days), €47/day (4-16 days), zero excess; Theft Protection included, theft excess €3000</li></ul><p>All charges are subject to VAT.</p>`,
		},
		{
			Title:       `Miscellaneous Coverages`,
			HtmlContent: `<p><strong>Third Party Liability Insurance (TPL)</strong>: Included in the rates</p><p>Covers the rented vehicle for claims from property damage and third party personal injury or death to the amounts specified by Greek Law.</p><p><strong>Personal Accident Coverage (PAC)</strong>: 8.50€ per day</p><p>Covering driver’s medical expenses resulting from an accident with the rented vehicle up to €1,500 and for the loss of life or total/partial permanent disablement: up to €50,000.</p><p><strong>Additional Drivers (ADR)</strong>: 7.50€ per day</p><p>Only the main driver is covered to drive the vehicle; any additional driver can be added by purchasing ADR.</p><p><strong>Cross Border Fee (CBF)</strong>: N/A</p><p>N/A</p><p><strong>Island/Ferryboat Travel Cover (FBF)</strong>: 20.00€ per day with a maximum of 3 days</p><p>Allowed only with the prior written approval of Autohellas-Hertz and provided that the customer has accepted the Ferry boat fee. Otherwise, the renter is liable for all damages caused to the rented vehicle during its transfer, regardless of any coverage and liability restrictions accepted and paid at the start of the rental. The ferry boat fee is available for purchase locally at any Autohellas-Hertz station.</p>`,
		},
		{
			Title:       `Fuel Policy & Refueling Service Charge`,
			HtmlContent: `<p>For rental durations of 1-2 days, the customer should return the car with the same quantity of fuel as shown at the time of pickup. In the event that the customer returns it with less quantity, the customer will be charged the value of the missing fuel. Additionally, a Refueling Service Charge will be applied, which today is €16.</p><p>For rental durations of three (3) days or more, the customer is required to prepay the cost of the quantity of fuel which is in the tank of the rented car at the pickup time. Upon returning the car –and only at this time– the customer may ask to be refunded with the value of the fuel in the tank at that time. In this instance, if the quantity of fuel in the tank is less than the one the customer has prepaid, then the Refuelling Service Charge will be applied, which today is €16. Any refund due will be processed through customer's credit/debit card.</p><p>The price per litre of fuel used is the average fuel price for the area where the car was rented.</p><p>The Refueling Service Charge is €16, while Hertz reserves the right to change it in the future without prior notice. Any fuel refunds, disputes etc. must be resolved with Hertz staff at the time</p><p>of the vehicle’s return. Hertz will not consider any disputes once the final invoice has been issued.</p><p>A fuel policy leaflet is available on all Hertz Desks in five languages.</p><p><strong>Charge:</strong> 16€ per case, plus the cost of fuel</p>`,
		},
		{
			Title:       `Charging Service Charge`,
			HtmlContent: `<p>For electric car groups, the battery level of the vehicle upon return by the renter should be at least 15%, which if less, the  re-charge service charge will be applied. Renter is fully responsible in case of loss or damage to the vehicle´s charging cable.  Cost will vary depending on model and manufacturer. In case of immobilisation of an electric car due to low battery level the renter  will be charged a fee of €25, for the transportation of the vehicle to the nearest  Autohellas-Hertz charging point.</p><p><strong>Charge:</strong> € 10 per Case</p>`,
		},
		{
			Title:       `Delivery & Collection`,
			HtmlContent: `<p>Delivery/Collection outside Hertz offices is available on request basis.</p><p><strong>Charge:</strong> Minimum charge for Delivery or Collection close to a Hertz station is 14€ per Delivery/Collection.</p>`,
		},
		{
			Title:       `Night Delivery & Collection`,
			HtmlContent: `<p>All reservations with pickup and/or drop off between 00:00 to 06:59 will be subject to night delivery charge (NDC).</p><p>Our API will return rates inclusive of NDC, providing that the booking is within this time frame.</p><p>If the booking is outside night delivery times but customer arrives between 00:00 to 06:59 then night delivery will be charged locally.</p><p><strong>Charge:</strong> 50€ per case (free for Athens A/P)</p>`,
		},
		{
			Title:       `Rent It Here, Leave it There`,
			HtmlContent: `<p>RIHLIT charge is applied for all rentals that are returned in a different city where a Hertz office is located. Rentals between two locations (Airport/City office) in the same city are free of charge.</p><p><strong>Charge:</strong> As per xml calculation</p>`,
		},
		{
			Title:       `Baby/Infant Seat & Booster`,
			HtmlContent: `<p>Baby/infant seat &amp; booster can be provided upon request.</p><p><strong>Charge:</strong> 7€ per day per item, with a maximum of 17 days</p>`,
		},
		{
			Title:       `GPS/Navigation System`,
			HtmlContent: `<p>GPS device can be provided upon request.</p><p><strong>Charge:</strong> 8€ per day per item, with a maximum of 17 days</p>`,
		},
		{
			Title:       `Incident Administration Fee`,
			HtmlContent: `<p>In case of an accident, the renter has the obligation to pay to Hertz a non-refundable fee to cover administration expenses.</p><p><strong>Charge:</strong> 30€ per case</p>`,
		},
		{
			Title:       `Traffic Fines`,
			HtmlContent: `<p>Will be paid by the renter. If not, Hertz will pay them and will charge the customer with the full amount of traffic fine plus the Incident Administration Fee, 30€ per case, in order to cover all administration expenses.</p>`,
		},
		{
			Title:       `24h Road Assistance`,
			HtmlContent: `<p>24h Road Assistance option is available to all renters.</p><p><strong>Charge:</strong> Included in the rates</p>`,
		},
		{
			Title:       `Mileage`,
			HtmlContent: `<p>A limit of 150 kilometeres applies to one-day rentals. Any excess of this limit will result to additional charge, depending on the car group (as shown on Fleet &amp; Product tab).</p><p>For duration between 2 and 30 days there is no milage limit.</p>`,
		},
		{
			Title:       `Airport Service Charge`,
			HtmlContent: `<p>Charges concerning rentals where the pick-up and/or drop-off location is an airport.</p><p><strong>Charge:</strong> Included in the rates</p>`,
		},
	},
	"ZT": {
		{
			Title:       `Driver's Qualitifications`,
			HtmlContent: `<p>A valid driving license held for at least one (1) year is required. The driver must hold a valid credit or debit card accepted by Hertz for all necessary pre-authorization amounts as well as pre-payments and settlement of final rental charges. Accepted credit/debit cards are Visa, Mastercard, American Express, Diners and Unionpay. Maestro debit cards are also accepted but only for payments and not for guarantee provision.</p>`,
		},
		{
			Title:       `Age Limits & Young Driver Surcharge`,
			HtmlContent: `<ul><li>From 19 to 20 years old – groups A, B, J: €10 per day + VAT</li><li>From 19 to 20 years old – other categories not available</li><li>From 21 to 24 years old – groups A, B, J: No additional charge</li><li>From 21 to 24 years old – groups C, D, F, K, M, T: €10 per day + VAT</li><li>From 21 to 24 years old – groups E, N, U, NM, O: €15 per day + VAT</li><li>From 21 to 24 years old – other categories not available</li><li>25 years and older – all car categories: No additional charge</li></ul>`,
		},
		{
			Title:       `Driving Restrictions`,
			HtmlContent: `<p>Vehicle may not be driven across any of the Greek borders without prior written approval by Thrifty Greece. Failure to comply with these restrictions makes void all insurance policies and customer fully liable and responsible in case of a theft, accident, damages to the vehicle and any resulting consequences.</p>`,
		},
		{
			Title:       `Grace Period`,
			HtmlContent: `<p>Detailed instructions on car pick-up will be provided for each Thrifty location.</p><p>At the end of the agreed rental period Thrifty allows a grace period of 29 minutes free of charge after which the customer will be charged with a minimum of one full day rental.</p>`,
		},
		{
			Title:       `Early Return / Rental Extension`,
			HtmlContent: `<p>No refund for any unused voucher value/rental days. Ongoing rental extension should be arranged with Thrifty at least 48 hours before the intended/reserved rental expiry, to be charged as per valid retail tariff.</p>`,
		},
		{
			Title:       `Late Cancellation / Non-Show Policy`,
			HtmlContent: `<p>Thrifty have the right, in case a customer fails either to show up or to cancel at least 48 hours prior to pick-up, to charge the client a Late Cancellation Fee of €56. If the customer fails to cancel their reservation prior to the pick-up time and does not collect the vehicle on the pick-up date, or fails to comply with the Rental Terms and Conditions (see Hertz Rental Qualifications &amp; Requirements), Dollar and Thrifty reserve the right to charge a No-Show Fee of €119 (billed in destination currency). A Non-Show is effective if the customer hasn't collected the vehicle within one (1) hour of the booked time/date, unless there is a travel delay (given that a flight number and/or travel details have been provided prior to the car pick-up).</p>`,
		},
		{
			Title:       `Fleet Availability`,
			HtmlContent: `<p>Thrifty reserves the right to provide a different car model to that originally reserved of an equivalent or higher model type.</p>`,
		},
		{
			Title:       `Accident & Theft Coverages`,
			HtmlContent: `<p><strong>Collision Damage Waiver (CDW)</strong></p><p>Obtaining Collission Damage Waiver (CDW) will reduced renter's liability for damage to the rented vehicle, as per rental's agreement terms &amp; conditions, to the amounts shown in the table below.  The following cases are not covered: Damage to the underside of the car, the wheels &amp; tyres, broken car glass or loss of the vehicle keys and documents.</p><p><strong>Reduced Excess Coverage (SCDW)</strong></p><p>Super Collission Damage Waiver (SCDW) can be obtained only after purchasing Collission Damage Waiver (CDW) and will reduce renter's liability further, to the amounts shown in the table below.</p><p>The following cases are not covered: Damage to the underside of the car, the wheels &amp; tyres, broken car glass or loss of the vehicle keys and documents.</p><p><strong>Super Cover (SC)</strong></p><p>Τhe renter can be covered (with zero excess) in case of total theft or damage to the Thrifty vehicle, provided that the renter has previously accepted the relevant coverage charges for CDW, SCDW, and TP, or has a reservation inclusive of these coverages. Moreover, Super Cover covers any damages caused on the vehicle's windshield and tyres and includes Extended Road Assistance. Damages caused to the underside and wheels of the vehicle are not covered by any insurance. Super Cover is not available for rentals made by drivers aged 19 and 20 years old.</p><p><strong>Tyre &amp; Glass Protection (TGI)</strong></p><p>Tyre and Glass Coverage provides cover against consequences of damages on tires, windscreen, side windows and back window. If TGI is not accepted, the client will be held liable for the full value of damage on tires, windscreen, side windows and back window.</p><p><strong>Theft Protection (TPC)</strong></p><p>Theft Protection (TPC) will reduce renter's liability to the amounts shown in the table below.</p>`,
		},
		{
			Title:       `Coverage Charges & Excess per Group`,
			HtmlContent: `<ul><li>Groups A, B, J: CDW included, excess €850; SCDW €14/day (1-3 days), €10.5/day (4-16 days), excess €300; Super Cover €16.5/day (1-3 days), €13/day (4-16 days), zero excess; Tyre &amp; Glass €5.5/day (max 17 days); Theft Protection included, theft excess €850</li><li>Groups C, D, F, K, M, T: CDW included, excess €1000; SCDW €16.5/day (1-3 days), €12.5/day (4-16 days), excess €350; Super Cover €18/day (1-3 days), €14/day (4-16 days), zero excess; Theft Protection included, theft excess €1000</li><li>Groups 8, E, N, NM, O, U, YL: CDW included, excess €1300; SCDW €26/day (1-3 days), €20/day (4-16 days), excess €550; Super Cover €29/day (1-3 days), €22/day (4-16 days), zero excess; Theft Protection included, theft excess €1300</li></ul><p>All charges are subject to VAT.</p>`,
		},
		{
			Title:       `Miscellaneous Coverages`,
			HtmlContent: `<p><strong>Third Party Liability Insurance (TPL)</strong>: Included in the rates</p><p>Covers the rented vehicle for claims from property damage and third party personal injury or death to the amounts specified by Greek Law.</p><p><strong>Personal Accident Coverage (PAC)</strong>: 10€ per day</p><p>Covering driver’s medical expenses resulting from an accident with the rented vehicle up to €1,500 and for the loss of life or total/partial permanent disablement: up to €50,000.</p><p><strong>Additional Drivers (ADR)</strong>: 7.50€ per day</p><p>Only the main driver is covered to drive the vehicle; any additional driver can be added by purchasing ADR.</p><p><strong>Cross Border Fee (CBF)</strong>: N/A</p><p>N/A</p><p><strong>Island/Ferryboat Travel Cover (FBF)</strong>: 20.00€ per day with a maximum of 3 days</p><p>Allowed only with the prior written approval of Autohellas-Thrifty and provided that the customer has accepted the Ferry boat fee. Otherwise, the renter is liable for all damages caused to the rented vehicle during its transfer, regardless of any coverage and liability restrictions accepted and paid at the start of the rental. The ferry boat fee is available for purchase locally at any Autohellas-Thrifty station</p>`,
		},
		{
			Title:       `Fuel Policy & Refueling Service Charge`,
			HtmlContent: `<p>For rental durations of 1-2 days, the customer should return the car with the same quantity of fuel as shown at the time of pickup. In the event that the customer returns it with less quantity, the customer will be charged the value of the missing fuel. Additionally, a "Refueling Service" charge will be applied which today is €20.</p><p>For rental durations of three (3) days or more, the customer is required to prepay the cost of the quantity of fuel which is in the tank of the rented car at the pickup time. Upon returning the car –and only at this time– the customer may ask to be refunded with the value of the fuel in the tank at that time. In this instance, if the quantity of fuel in the tank is less than the one the customer has prepaid, then the 'Refuelling Service" charge will be applied, which today is €20. Any refund due will be processed through customer's credit/debit card.</p><p>The price per litre of fuel used is the average fuel price for the area where the car was rented, as published in the Price Observatory of the Hellenic Ministry of Development and Investments (site: www.fuelprices.gr).</p><p>The 'Refueling Service Charge'  is €20, while Dollar reserves the right to change it in the future without prior notice. Any fuel refunds, disputes etc. must be resolved with Dollar staff at the time of the vehicle’s return. Dollar will not consider any disputes once the final invoice has been issued. A fuel policy leaflet is available on all Dollar Desks in five languages.</p><p><strong>Charge:</strong> 20€ per case, plus the cost of fuel</p>`,
		},
		{
			Title:       `Delivery & Collection`,
			HtmlContent: `<p>Delivery/Collection outside Thrifty offices is available on request basis.</p><p><strong>Charge:</strong> Minimum charge for Delivery or Collection close to a Thrifty station is 14€ per Delivery/Collection.</p>`,
		},
		{
			Title:       `Night Delivery & Collection`,
			HtmlContent: `<p>All reservations with pickup and/or drop off between 00:00 to 06:59 will be subject to night delivery charge (NDC).</p><p>Our API will return rates inclusive of NDC, providing that the booking is within this time frame.</p><p>If the booking is outside night delivery times but customer arrives between 00:00 to 06:59 then night delivery will be charged locally.</p><p><strong>Charge:</strong> 50€ per case (free for Athens A/P)</p>`,
		},
		{
			Title:       `Rent It Here, Leave it There`,
			HtmlContent: `<p>RIHLIT charge is applied for all rentals that are returned in a different city where a Thrifty office is located. Rentals between two locations (Airport/City office) in the same city are free of charge.</p><p><strong>Charge:</strong> As per xml calculation</p>`,
		},
		{
			Title:       `Baby/Infant Seat & Booster`,
			HtmlContent: `<p>Baby/infant seat &amp; booster can be provided upon request.</p><p><strong>Charge:</strong> 7€ per day per item, with a maximum of 17 days</p>`,
		},
		{
			Title:       `GPS/Navigation System`,
			HtmlContent: `<p>GPS device can be provided upon request.</p><p><strong>Charge:</strong> 8€ per day per item, with a maximum of 14 days</p>`,
		},
		{
			Title:       `Incident Administration Fee`,
			HtmlContent: `<p>In case of an accident, the renter has the obligation to pay to Thrifty a non-refundable fee to cover administration expenses.</p><p><strong>Charge:</strong> 30€ per case</p>`,
		},
		{
			Title:       `Traffic Fines`,
			HtmlContent: `<p>Will be paid by the renter. If not, Thrifty will pay them and will charge the customer with the full amount of traffic fine plus the Incident Administration Fee, 30€ per case, in order to cover all administration expenses.</p>`,
		},
		{
			Title:       `24h Road Assistance`,
			HtmlContent: `<p>24h Road Assistance option is available to all renters.</p><p><strong>Charge:</strong> Included in the rates</p>`,
		},
		{
			Title:       `Mileage`,
			HtmlContent: `<p>A limit of 150 kilometeres applies to one-day rentals. Any excess of this limit will result to additional charge, depending on the car group (as shown on Fleet &amp; Product tab).</p><p>For duration between 2 and 30 days there is no milage limit.</p>`,
		},
		{
			Title:       `Airport Service Charge`,
			HtmlContent: `<p>Charges concerning rentals where the pick-up and/or drop-off location is an airport.</p><p><strong>Charge:</strong> Included in the rates</p>`,
		},
	},
	"ZR": {
		{
			Title:       `Driver's Qualitifications`,
			HtmlContent: `<p>A valid driving license held for at least one (1) year is required. The driver must hold a valid credit or debit card accepted by Hertz for all necessary pre-authorization amounts as well as pre-payments and settlement of final rental charges. Accepted credit/debit cards are Visa, Mastercard, American Express, Diners and Unionpay. Maestro debit cards are also accepted but only for payments and not for guarantee provision.</p>`,
		},
		{
			Title:       `Age Limits & Young Driver Surcharge`,
			HtmlContent: `<ul><li>From 19 to 20 years old – groups A, B, J: €10 per day + VAT</li><li>From 19 to 20 years old – other categories not available</li><li>From 21 to 24 years old – groups A, B, J: No additional charge</li><li>From 21 to 24 years old – groups C, D, F, K, M: €10 per day + VAT</li><li>From 21 to 24 years old – groups N, NM, O: €15 per day + VAT</li><li>From 21 to 24 years old – other categories not available</li><li>25 years and older – all car categories: No additional charge</li></ul>`,
		},
		{
			Title:       `Driving Restrictions`,
			HtmlContent: `<p>Vehicle may not be driven across any of the Greek borders without prior written approval by Dollar Greece. Failure to comply with these restrictions makes void all insurance policies and customer fully liable and responsible in case of a theft, accident, damages to the vehicle and any resulting consequences.</p>`,
		},
		{
			Title:       `Grace Period`,
			HtmlContent: `<p>Detailed instructions on car pick-up will be provided for each Dollar location.</p><p>At the end of the agreed rental period Dollar allows a grace period of 29 minutes free of charge after which the customer will be charged with a minimum of one full day rental.</p>`,
		},
		{
			Title:       `Early Return / Rental Extension`,
			HtmlContent: `<p>No refund for any unused voucher value/rental days. Ongoing rental extension should be arranged with Dollar at least 48 hours before the intended/reserved rental expiry, to be charged as per valid retail tariff.</p>`,
		},
		{
			Title:       `Late Cancellation / Non-Show Policy`,
			HtmlContent: `<p>Dollar have the right, in case a customer fails either to show up or to cancel at least 48 hours prior to pick-up, to charge the client a Late Cancellation Fee of €56. If the customer fails to cancel their reservation prior to the pick-up time and does not collect the vehicle on the pick-up date, or fails to comply with the Rental Terms and Conditions (see Hertz Rental Qualifications &amp; Requirements), Dollar and Thrifty reserve the right to charge a No-Show Fee of €119 (billed in destination currency). A Non-Show is effective if the customer hasn't collected the vehicle within one (1) hour of the booked time/date, unless there is a travel delay (given that a flight number and/or travel details have been provided prior to the car pick-up).</p>`,
		},
		{
			Title:       `Fleet Availability`,
			HtmlContent: `<p>Dollar reserves the right to provide a different car model to that originally reserved of an equivalent or higher model type.</p>`,
		},
		{
			Title:       `Accident & Theft Coverages`,
			HtmlContent: `<p><strong>Collision Damage Waiver (CDW)</strong></p><p>Obtaining Collission Damage Waiver (CDW) will reduced renter's liability for damage to the rented vehicle, as per rental's agreement terms &amp; conditions, to the amounts shown in the table below.  The following cases are not covered: Damage to the underside of the car, the wheels &amp; tyres, broken car glass or loss of the vehicle keys and documents.</p><p><strong>Reduced Excess Coverage (SCDW)</strong></p><p>Super Collission Damage Waiver (SCDW) can be obtained only after purchasing Collission Damage Waiver (CDW) and will reduce renter's liability further, to the amounts shown in the table below.</p><p>The following cases are not covered: Damage to the underside of the car, the wheels &amp; tyres, broken car glass or loss of the vehicle keys and documents.</p><p><strong>Super Cover (SC)</strong></p><p>Τhe renter can be covered (with zero excess) in case of total theft or damage to the Dollar vehicle, provided that the renter has previously accepted the relevant coverage charges for CDW, SCDW, and TP, or has a reservation inclusive of these coverages. Moreover, Super Cover covers any damages caused on the vehicle's windshield and tyres and includes Extended Road Assistance. Damages caused to the underside and wheels of the vehicle are not covered by any insurance. Super Cover is not available for rentals made by drivers aged 19 and 20 years old.</p><p><strong>Tyre &amp; Glass Protection (TGI)</strong></p><p>Tyre and Glass Coverage provides cover against consequences of damages on tires, windscreen, side windows and back window. If TGI is not accepted, the client will be held liable for the full value of damage on tires, windscreen, side windows and back window.</p><p><strong>Theft Protection (TPC)</strong></p><p>Theft Protection (TPC) will reduce renter's liability to the amounts shown in the table below.</p>`,
		},
		{
			Title:       `Coverage Charges & Excess per Group`,
			HtmlContent: `<ul><li>Groups A, B, J: CDW included, excess €950; SCDW €16.5/day (1-3 days), €12.5/day (4-16 days), excess €400; Super Cover €22/day (1-3 days), €17/day (4-16 days), zero excess; Tyre &amp; Glass €5.5/day (max 17 days); Theft Protection included, theft excess €950</li><li>Groups C, D, F, K, M: CDW included, excess €1100; SCDW €19/day (1-3 days), €15/day (4-16 days), excess €450; Super Cover €26/day (1-3 days), €20/day (4-16 days), zero excess; Theft Protection included, theft excess €1100</li><li>Groups O, N, NM: CDW included, excess €1500; SCDW €30/day (1-3 days), €23/day (4-16 days), excess €700; Super Cover €37/day (1-3 days), €28.5/day (4-16 days), zero excess; Theft Protection included, theft excess €1500</li></ul><p>All charges are subject to VAT.</p>`,
		},
		{
			Title:       `Miscellaneous Coverages`,
			HtmlContent: `<p><strong>Third Party Liability Insurance (TPL)</strong>: Included in the rates</p><p>Covers the rented vehicle for claims from property damage and third party personal injury or death to the amounts specified by Greek Law.</p><p><strong>Personal Accident Coverage (PAC)</strong>: 13.00€ per day</p><p>Covering driver’s medical expenses resulting from an accident with the rented vehicle up to €1,500 and for the loss of life or total/partial permanent disablement: up to €50,000.</p><p><strong>Additional Drivers (ADR)</strong>: 7.50€ per day</p><p>Only the main driver is covered to drive the vehicle; any additional driver can be added by purchasing ADR.</p><p><strong>Cross Border Fee (CBF)</strong>: N/A</p><p>N/A</p><p><strong>Island/Ferryboat Travel Cover (FBF)</strong>: 20.00€ per day with a maximum of 3 days</p><p>Allowed only with the prior written approval of Autohellas-Dollar and provided that the customer has accepted the Ferry boat fee. Otherwise, the renter is liable for all damages caused to the rented vehicle during its transfer, regardless of any coverage and liability restrictions accepted and paid at the start of the rental. The ferry boat fee is available for purchase locally at any Autohellas-Dollar station</p>`,
		},
		{
			Title:       `Fuel Policy & Refueling Service Charge`,
			HtmlContent: `<p>For rental durations of 1-2 days, the customer should return the car with the same quantity of fuel as shown at the time of pickup. In the event that the customer returns it with less quantity, the customer will be charged the value of the missing fuel. Additionally, a "Refueling Service" charge will be applied which today is €20.</p><p>For rental durations of three (3) days or more, the customer is required to prepay the cost of the quantity of fuel which is in the tank of the rented car at the pickup time. Upon returning the car –and only at this time– the customer may ask to be refunded with the value of the fuel in the tank at that time. In this instance, if the quantity of fuel in the tank is less than the one the customer has prepaid, then the 'Refuelling Service" charge will be applied, which today is €20. Any refund due will be processed through customer's credit/debit card.</p><p>The price per litre of fuel used is the average fuel price for the area where the car was rented, as published in the Price Observatory of the Hellenic Ministry of Development and Investments (site: www.fuelprices.gr).</p><p>The 'Refueling Service Charge'  is €20, while Dollar reserves the right to change it in the future without prior notice. Any fuel refunds, disputes etc. must be resolved with Dollar staff at the time of the vehicle’s return. Dollar will not consider any disputes once the final invoice has been issued. A fuel policy leaflet is available on all Dollar Desks in five languages.</p><p><strong>Charge:</strong> 20€ per case, plus the cost of fuel</p>`,
		},
		{
			Title:       `Night Delivery & Collection`,
			HtmlContent: `<p>All reservations with pickup and/or drop off between 00:00 to 06:59 will be subject to night delivery charge (NDC).</p><p>Our API will return rates inclusive of NDC, providing that the booking is within this time frame.</p><p>If the booking is outside night delivery times but customer arrives between 00:00 to 06:59 then night delivery will be charged locally.</p><p><strong>Charge:</strong> 50€ per case (free for Athens A/P)</p>`,
		},
		{
			Title:       `Rent It Here, Leave it There`,
			HtmlContent: `<p>RIHLIT charge is applied for all rentals that are returned in a different city where a Dollar office is located. Rentals between two locations (Airport/City office) in the same city are free of charge.</p><p><strong>Charge:</strong> As per xml calculation</p>`,
		},
		{
			Title:       `Baby/Infant Seat & Booster`,
			HtmlContent: `<p>Baby/infant seat &amp; booster can be provided upon request.</p><p><strong>Charge:</strong> 7€ per day per item, with a maximum of 17 days</p>`,
		},
		{
			Title:       `GPS/Navigation System`,
			HtmlContent: `<p>GPS device can be provided upon request.</p><p><strong>Charge:</strong> 8€ per day per item, with a maximum of 17 days</p>`,
		},
		{
			Title:       `Incident Administration Fee`,
			HtmlContent: `<p>In case of an accident, the renter has the obligation to pay to Hertz a non-refundable fee to cover administration expenses.</p><p><strong>Charge:</strong> 30€ per case</p>`,
		},
		{
			Title:       `Traffic Fines`,
			HtmlContent: `<p>Will be paid by the renter. If not, Dollar will pay them and will charge the customer with the full amount of traffic fine plus the Incident Administration Fee, 30€ per case, in order to cover all administration expenses.</p>`,
		},
		{
			Title:       `Additional Road Assistance Cover`,
			HtmlContent: `<p>Dollar offers 24h Road Assistance Service, if required, at a fixed fee of €25 per case. This fee can be waived by obtaining the additional Road Assistance cover at a cost of 8€ per rental.</p><p><strong>Charge:</strong> 8€ per rental</p>`,
		},
		{
			Title:       `Mileage`,
			HtmlContent: `<p>A limit of 150 kilometeres applies to one-day rentals. Any excess of this limit will result to additional charge, depending on the car group (as shown on Fleet &amp; Product tab).</p><p>For duration between 2 and 30 days there is no milage limit.</p>`,
		},
		{
			Title:       `Airport Service Charge`,
			HtmlContent: `<p>Charges concerning rentals where the pick-up and/or drop-off location is an airport.</p><p><strong>Charge:</strong> Included in the rates</p>`,
		},
	},
}
