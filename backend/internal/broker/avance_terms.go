package broker

// avanceTerms are Avance's terms and conditions, generated from the Terms sheet of their
// "Avance Fleet & Terms 2026" workbook, wording unchanged. Regenerate them from the new workbook
// when Avance updates its terms.
var avanceTerms = []TermsAndConditionsItem{
	{
		Title:       "After hours deliveries – collections",
		HtmlContent: `<p>Cost paid by client : 28€ per service</p>`,
	},
	{
		Title:       "First Additional Driver",
		HtmlContent: `<p>Free of Charge , Second or more drivers with additional charge. The maximum number of drivers is set at 3 persons.</p>`,
	},
	{
		Title:       "Maximum Charge of Extras (Snow Chains, CS/BS )",
		HtmlContent: `<p>Over 7 days the maximum charge is set at 60€/rental for Child and Booster Seats/ Maximum Charge of Snow chains is set to 50€/rental</p>`,
	},
	{
		Title:       "Delivery Service - Collections",
		HtmlContent: `<p>Upon Request - Minimum Charge 20€ per service</p>`,
	},
	{
		Title:       "Prices Inlude",
		HtmlContent: `<p>Collision Damage Waiver with Excess based on Category*</p><p>Free Mileages</p>`,
	},
	{
		Title:       "Guarantee",
		HtmlContent: `<p>Credit/Debit card under the name of driver required for guarantee upon delivery of the car.</p><p>Upon arrival, station will preauthorize/charge an amount from the driver's credit card as a deposit for guarantee which will be fully released/refunded at the end of the rental, provided no other charges arise. Please note that the refund can take up to 14 business days, depending on the issuing bank of the credit card. The amount of the deposit varies according to the car group.</p>`,
	},
	{
		Title:       "Reservations",
		HtmlContent: `<p>All reservations made are based on the vehicle category, not for a specific model.</p><p>In case of non-availability of a vehicle in the category pre-selected by the customer, always with reference to a confirmed reservation, Avance reserves the right to provide to the customer a higher category car, at no extra charge for the customer.</p>`,
	},
	{
		Title:       "Grace Period",
		HtmlContent: `<p>During the peak season, from June 15th to September 15th, and the Easter period, from April 6th to April 13th this year, the grace period will be 3h or until the end of working hours for reservations with a prior notice of delay (flight number etc.).</p><p>In addition, during the low season, the grace period will be until the end of working hours for reservations with or without a prior notice of delay.</p>`,
	},
	{
		Title:       "Means of Payment",
		HtmlContent: `<p>The Renter is required to produce a credit card/ debid card as guarantee for the rental, even if payment is made in cash.</p><p>Rental is prepaid. Accepted Credit / Debit Cards are Visa and Mastercard.</p>`,
	},
	{
		Title:       "Drivers Age",
		HtmlContent: `<p>Minimum age limit of 19 years for the use of cars of category A,A1,A3,B,C,C3,D,GM,GM3,GMH with the additional charge for a young driver at €10 per day * available at some stations. For the rest of the categories, the 23 years apply, without additional charge, and for all the stations of the Avance Network.</p><p>Maximum age limit is 75 years.</p>`,
	},
	{
		Title:       "Driving License",
		HtmlContent: `<p>Production of a valid license for driving in the Greek territory and the E.U., issued at least 1 year prior to the rental’s starting date, is required in hard copy.</p><p>Citizens of countries outside the E.U. are required to be in possession of a driver’s license valid both in Greece and in the European Union, or otherwise holders of an international driver’s license which they are required to present together with the valid driver’s license issued in their country. *</p>`,
	},
	{
		Title:       "Road Assistance",
		HtmlContent: `<p>Avance Vehicles are covered by 24hours road assistance. Road assistance’s contact details are stated on the rental contract.</p>`,
	},
	{
		Title:       "Fuel",
		HtmlContent: `<p>The renter has to return the car with the same quantity of fuel it had when collected. In case the car is returned with a quantity of fuel less than the original (namely, less than the quantity of fuel the vehicle had when collected by the customer), the customer shall be charged with the fuel difference.</p><p>No customer refund is provided in case of the car returned with a quantity of fuel greater than that collected.</p><p>Vehicles are delivered with a minimum quantity of fuel, corresponding to ¼ of the reservoir.</p>`,
	},
	{
		Title:       "Accident",
		HtmlContent: `<p>In case of an accident, the customer is charged with 30€ + Vat for the accident file administration costs.</p>`,
	},
	{
		Title:       "Replacement",
		HtmlContent: `<p>Where replacement of the vehicle deemed as necessary, up to 24 hours are required for areas in the Mainland, while for any other area within the Greek territory, a period of up to 72 hours is required.Driving outside road Network is prohibited.</p>`,
	},
	{
		Title:       "Traffic Tickets-Fines",
		HtmlContent: `<p>Traffic tickets and fines due to violations of the Highway Code are borne by the Renter.</p><p>In addition, the Renter undertakes the payment of administration costs, where required.</p><p>Please note that the administration costs may vary depending location with max. charge of 30€+Vat.</p>`,
	},
	{
		Title:       "Ferry Restrictions",
		HtmlContent: `<p>Transporting the car by ship is only permitted with the written approval of the company.</p><p>In the event the company accepts the request, the lessee can choose the additional coverage of the transportation by ship (Ferry Fee) which is 20€ per day with a maximum charge of 60€. In any other case, the renter is fully responsible for any damages that may be caused during its transportation, regardless of any additional insurance coverage accepted and paid for at the start of the rental.</p><p>In case the Vehicle breaks down on an island, Avance will not be able to provide a replacement locally. The rental vehicle will have to, if possible,  be repaired locally in case of mechanical damage and any expenses will be refunded upon client's return once all necessary invoices issued to our corporate details is provided. In any case, the renter must contact Avance before proceeding to any actions and advise the cost of the damage.</p><p>If the rental vehicle is immobilized, the renter will be responsible to arrange and pay for the transfer of the car to the nearest land destination, to which will be provided with a replacement vehicle.</p>`,
	},
	{
		Title:       "Car Transport Out of the Country",
		HtmlContent: `<p>Transport and driving of the vehicle outside the borders of the Greek Territory is prohibited.</p>`,
	},
	{
		Title:       "Claims/Complaints",
		HtmlContent: `<p>Are accepted by email only at customer-service@avance.gr maximum 30 days after the completion of the rental.</p><p>Avance reserves the right to response within 10 business days.</p>`,
	},
	{
		Title:       "Drivers License",
		HtmlContent: `<table><thead><tr><th>Drivers must present</th><th>Country</th></tr></thead><tbody><tr><td>Valid National Driving Licence</td><td>EU Nationals (including Norway, Iceland, United Kingdom of Great Britain and N. Ireland, Switzerland &amp; Liechtenstein nationals), USA, Australia, Canada.</td></tr><tr><td>Only card-type Valid National Driving License</td><td>Albania, Armenia, Azerbaijan, Bahamas, Bahrain, Belarus, Bosnia and Herzegovina, Brazil, Central African Republic, Chile, Costa Rica, Cuba, Democratic Republic of the Congo, Ecuador, Georgia, Ghana, Guyana, Holy See, Indonesia, Iran, Iraq, Israel, Ivory Coast, Kazakhstan, Kenya, Kuwait, Kyrgyzstan, Liberia, Mexico, Mongolia, Montenegro, Morocco, Niger, North Macedonia, Pakistan, Peru, Philippines, Qatar, Republic of Korea, Republic of Moldova, San Marino, Saudi Arabia, Senegal, Serbia, Seychelles, South Africa, Tajikistan, Thailand, Tunisia, Turkey, Turkmenistan, Ukraine, United Arab Emirates, Uruguay, Uzbekistan, Venezuela, Vietnam, Zimbabwe</td></tr><tr><td>National driving license and official English translation</td><td>Republic of China, Japan</td></tr><tr><td>Both Valid National AND International Driving License</td><td>All other countries</td></tr></tbody></table>`,
	},
	{
		Title:       "Cancellation Fee",
		HtmlContent: `<p>Low Season: Cancellation of a reservation must be done at least same day  prior to the scheduled date and hour of commencement of the vehicle’s rental (always according to the exact hour of reservation), as this has been made in the company’s reservations system with no extra charge.</p><p>High Season (Pick up July - August): In case of a cancellation notice received by Avance less than fourty eight (48) hours prior to the scheduled commencement of the vehicle’s, there shall be 50% of total amount cancellation fee applied.</p>`,
	},
	{
		Title:       "Insurance Policy",
		HtmlContent: `<p><strong>Extra Incurances:</strong></p><p>The renter is fully released of the liability for any car damage:</p><p><strong>SCDW – Super Collision Damage Waiver:</strong></p><ul><li>over 350€+VAT provided he/she accepts and signs the relevant terms of the company’s contract and pays the amount of 6€ on top of the CDW for each rental day in respect of the categories A,A1,A3,B,C,C3,GM,GM3,GMH</li><li>over 400€+VAT provided he/she accepts and signs the relevant terms of the company’s contract and pays the amount of 7€ on top of the CDW for each rental day in respect of the categories C1C,C2C,C2CM,D,DX,DX3,G2,G2X,G3X,K</li><li>over 475€+VAT provided he/she accepts and signs the relevant terms of the company’s contract and pays the amount of 8€ on top of the CDW for each rental day in respect of the categories F7,FAD,FFD,FR1,G4,GD,E-GQ,GQX,GQ,GQ1,GQ3,H1,J1,J2,K2,K2C,K2F,K2MAN</li><li>over 625€+VAT provided he/she accepts and signs the relevant terms of the company’s contract and pays the amount of 9€ on top of the CDW for each rental day in respect of the categories F1,F2D,F3D,FD,FDL,FF1,FFL,FFO,FFOA,FR2,J4,J4H,J5,J5A,J7</li><li>over 750€+VAT provided he/she accepts and signs the relevant terms of the company’s contract and pays the amount of 15€ on top of the CDW for each rental day in respect of the categories L1</li><li>over 1000€+VAT provided he/she accepts and signs the relevant terms of the company’s contract and pays the amount of 25€ on top of the CDW for each rental day in respect of the categories JL,J3,L3</li></ul><p>As a necessary condition for release from liability, the damage must be proven not to have been caused by a Highway Code violation.</p><p>The renter is fully released of the liability for any car damage provided he/she accepts and signs the relevant terms of the company’s contract and pays the amount of:</p><p><strong>FDW – Full Damage Waiver:</strong></p><ul><li>12€ on top of the CDW for each rental day in respect of the categories A,A1,A3,B,C,C3,GM,GM3,GMH</li><li>14€ on top of the CDW for each rental day in respect of the categories C1C,C2C,C2CM,D,DX,DX3,G2,G2X,G3X,K</li><li>16€ on top of the CDW for each rental day in respect of the categories F7,FAD,FFD,FR1,G4,GD,E-GQ,GQX,GQ,GQ1,GQ3,H1,J1,J2,K2,K2C,K2F,K2MAN</li><li>18€ on top of the CDW for each rental day in respect of the categories F1,F2D,F3D,FD,FDL,FF1,FFL,FFO,FFOA,FR2,J3,J4,J4H,J5,J5A,J7</li><li>30€ on top of the CDW for each rental day in respect of the categories L1</li><li>50€ on top of the CDW for each rental day in respect of the categories JL,J3,L3</li></ul><p>As a necessary condition for release from liability, the damage must be proven not to have been caused by a Highway Code violation.</p><p>No Insurance Covers Damages to Interior, Glass, Underneath and tyres Damage to Convertible Soft Top is not covered and will be paid in full by the customer.</p><p><strong>Insurance Limits:</strong></p><p>Vehicles are insured against third parties (excluding the driver) for death and bodily harm up to the amount of one million two hundred and twenty thousand Euros (1,220,000.00€) per accident and for material damages by third parties (excluding the AVANCE vehicle) up to the amount of one million two hundred and twenty thousand Euros (1,220,000.00€).</p><p>In case of change of the above limits, the minimum limits, as defined by the applicable Greek legislation, shall always apply. Damages exceeding the insurance coverage limit, or excluded by it, shall be fully paid by the customer.</p><p><strong>WUG</strong></p><p>The renter is fully released of the liability for any damage to the car’s tires or glass only provided he/she accepts and signs the relevant terms of the company’s contract and pays the amount of:</p><ul><li>10€/ rental day for the categories A,A1,A3,B,C,C3,GM,GM3,GMH</li><li>12€/ rental day for the categories C1C,C2C,C2CM,D,DX,DX3,G2,G2X,G3X,K</li><li>14€/ rental day for the categories F7,FAD,FFD,FR1,G4,GD,E-GQ,GQX,GQ,GQ1,GQ3,H1,J1,J2,K2,K2C,K2F,K2MAN</li><li>15€/ rental day for the categories F1,F2D,F3D,FD,FDL,FF1,FFL,FFO,FFOA,FR2,J4,J4H,J5,J5A,J7</li><li>25€/ rental day for the categories L1</li><li>45€/ rental day for the categories JL,J3,L3</li></ul><p><strong>TP</strong></p><p>The renter is fully released of the liability for any damage to the car in the event of theft, provided he/she accepts and signs the relevant terms of the company’s contract and pays the amount of:</p><ul><li>10€/ rental day for the categories A,A3,B,C,C3,GM,GM3,GMH</li><li>12€/ rental day for the categories C1C,C2C,C2CM,D,DX,DX3,G2,G2X,G3X,K</li><li>14€/ rental day for the categories F7,FAD,FFD,FR1,G4,GD,E-GQ,GQX,GQ,GQ1,GQ3,H1,J1,J2,K2,K2C,K2F,K2MAN</li><li>15€/ rental day for the categories F1,F2D,F3D,FD,FDL,FF1,FFL,FFO,FFOA,FR2,J4,J4H,J5,J5A,J7</li><li>25€/ rental day for the categories L1</li><li>45€/ rental day for the categories JL,J3,L3</li></ul><p>Main condition for the above coverage is the absence of negligence, even slight, or of unlawful conduct during the vehicle’s rental.</p>`,
	},
}
