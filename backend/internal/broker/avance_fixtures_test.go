package broker

// Live Wheelsys responses for Athens Airport (022), 15/11/2026 to 22/11/2026, captured through
// the egress proxy. Each is trimmed to the groups, options and stations the tests exercise;
// every attribute kept is verbatim. FAD is on request, and the FDW0 quote holds the two deposit
// edge cases: E-GQ's theft waiver reports dep="0", and GM3's FDW waiver has no exc attribute.

const avanceFixtureQuoteVCH = `<pricequote id="e916bf31-135c-4123-ab94-3e654de7edc4" validto="29/09/2026 05:27" duration="7" countrycode="GR" taxinclusive="true" currency="EUR" timetaken="46" graceperiod="240" owprepaid="false" oohprepaid="false" oooprepaid="false" fuelpolicy="SL">
	<rates>
		<category availability="AVAILABLE" cat="Passenger Cars" code="A" acriss="MBMR" imageurl="https://wheels-assets.s3.eu-central-1.amazonaws.com/10268/fleet/groupPhoto-A.jpg" pax="4" bags="0" doors="4" suitcases="1" model="Skoda Citigo or similar" totalrate="7766" baserate="7766" onewaycharge="0" ownet="0" outofhours="0" oohnet="0" outofoffice="0" ooonet="0" excess="70000" excessapplies="true" IncKlm="99999" unlimited="true" addklmrate="0" prepaidamount="7766" rentalamount="7766" rentalnet="6263" hourscharged="0" deliveryratenet="0" collectionratenet="0">
			<raterules />
			<options>
				<option code="CDW" rate="0" firstfree="false" inclusive="true" chargetype="I" prepaid="true" mandatory="true" calc="D" exc="70000" dep="70000" net="0" />
				<option code="FDW" rate="8400" firstfree="false" inclusive="false" chargetype="I" mandatory="false" calc="D" exc="0" dep="0" net="6774" />
				<option code="THIRD" rate="0" firstfree="false" inclusive="true" chargetype="I" prepaid="true" mandatory="true" calc="D" net="0" />
				<option code="THW" rate="0" firstfree="false" inclusive="true" chargetype="I" prepaid="true" mandatory="true" calc="D" exc="70000" dep="70000" net="0" />
				<option code="ADD" rate="2800" firstfree="true" inclusive="false" chargetype="E" mandatory="false" calc="D" net="2258" />
				<option code="BS" rate="4900" firstfree="false" inclusive="false" chargetype="E" mandatory="false" calc="D" net="3952" />
				<option code="CS" rate="4900" firstfree="false" inclusive="false" chargetype="E" mandatory="false" calc="D" net="3952" />
				<option code="SNC" rate="5000" firstfree="false" inclusive="false" chargetype="E" mandatory="false" calc="D" net="4032" />
				<option code="YDR" rate="7001" firstfree="false" inclusive="false" chargetype="E" mandatory="false" calc="D" net="5646" />
			</options>
			<taxes tax1code="VAT" tax1rate="2400" tax1amount="1503" tax2code="LCT" tax2rate="0" tax2amount="0" />
		</category>
		<category availability="AVAILABLE" cat="Passenger Cars" code="D" acriss="HDMR" imageurl="https://wheels-assets.s3.eu-central-1.amazonaws.com/10268/fleet/groupPhoto-D.jpg" pax="5" bags="0" doors="4" suitcases="1" model="Toyota Auris or similar" totalrate="11463" baserate="11463" onewaycharge="0" ownet="0" outofhours="0" oohnet="0" outofoffice="0" ooonet="0" excess="80000" excessapplies="true" IncKlm="99999" unlimited="true" addklmrate="0" prepaidamount="11463" rentalamount="11463" rentalnet="9244" hourscharged="0" deliveryratenet="0" collectionratenet="0">
			<raterules />
			<options>
				<option code="CDW" rate="0" firstfree="false" inclusive="true" chargetype="I" prepaid="true" mandatory="true" calc="D" exc="80000" dep="80000" net="0" />
				<option code="FDW" rate="9800" firstfree="false" inclusive="false" chargetype="I" mandatory="false" calc="D" exc="0" dep="0" net="7903" />
				<option code="THIRD" rate="0" firstfree="false" inclusive="true" chargetype="I" prepaid="true" mandatory="true" calc="D" net="0" />
				<option code="THW" rate="0" firstfree="false" inclusive="true" chargetype="I" prepaid="true" mandatory="true" calc="D" exc="80000" dep="80000" net="0" />
				<option code="ADD" rate="2800" firstfree="true" inclusive="false" chargetype="E" mandatory="false" calc="D" net="2258" />
				<option code="BS" rate="4900" firstfree="false" inclusive="false" chargetype="E" mandatory="false" calc="D" net="3952" />
				<option code="CS" rate="4900" firstfree="false" inclusive="false" chargetype="E" mandatory="false" calc="D" net="3952" />
				<option code="SNC" rate="5000" firstfree="false" inclusive="false" chargetype="E" mandatory="false" calc="D" net="4032" />
				<option code="YDR" rate="7001" firstfree="false" inclusive="false" chargetype="E" mandatory="false" calc="D" net="5646" />
			</options>
			<taxes tax1code="VAT" tax1rate="2400" tax1amount="2219" tax2code="LCT" tax2rate="0" tax2amount="0" />
		</category>
		<category availability="AVAILABLE" cat="SUV" code="E-GQ" acriss="DFAE" imageurl="https://wheels-assets.s3.eu-central-1.amazonaws.com/10268/fleet/groupPhoto-E-GQ.jpg" pax="5" bags="0" doors="5" suitcases="2" model="Geely EX5 Electric or similar" totalrate="21859" baserate="21859" onewaycharge="0" ownet="0" outofhours="0" oohnet="0" outofoffice="0" ooonet="0" excess="95000" excessapplies="true" IncKlm="99999" unlimited="true" addklmrate="0" prepaidamount="21859" rentalamount="21859" rentalnet="17628" hourscharged="0" deliveryratenet="0" collectionratenet="0">
			<raterules />
			<options>
				<option code="CDW" rate="0" firstfree="false" inclusive="true" chargetype="I" prepaid="true" mandatory="true" calc="D" exc="95000" dep="95000" net="0" />
				<option code="FDW" rate="11200" firstfree="false" inclusive="false" chargetype="I" mandatory="false" calc="D" exc="0" dep="0" net="9032" />
				<option code="THIRD" rate="0" firstfree="false" inclusive="true" chargetype="I" prepaid="true" mandatory="true" calc="D" net="0" />
				<option code="THW" rate="0" firstfree="false" inclusive="true" chargetype="I" prepaid="true" mandatory="true" calc="D" exc="95000" dep="95000" net="0" />
				<option code="ADD" rate="2800" firstfree="true" inclusive="false" chargetype="E" mandatory="false" calc="D" net="2258" />
				<option code="BS" rate="4900" firstfree="false" inclusive="false" chargetype="E" mandatory="false" calc="D" net="3952" />
				<option code="CS" rate="4900" firstfree="false" inclusive="false" chargetype="E" mandatory="false" calc="D" net="3952" />
				<option code="SNC" rate="5000" firstfree="false" inclusive="false" chargetype="E" mandatory="false" calc="D" net="4032" />
			</options>
			<taxes tax1code="VAT" tax1rate="2400" tax1amount="4231" tax2code="LCT" tax2rate="0" tax2amount="0" />
		</category>
		<category availability="AVAILABLE" cat="Passenger Cars" code="GM3" acriss="MDAR" imageurl="https://wheels-assets.s3.eu-central-1.amazonaws.com/10268/fleet/groupPhoto-GM3.jpg" pax="4" bags="0" doors="4" suitcases="1" model="Peugeot 108 automatic or similar" totalrate="10818" baserate="10818" onewaycharge="0" ownet="0" outofhours="0" oohnet="0" outofoffice="0" ooonet="0" excess="70000" excessapplies="true" IncKlm="99999" unlimited="true" addklmrate="0" prepaidamount="10818" rentalamount="10818" rentalnet="8724" hourscharged="0" deliveryratenet="0" collectionratenet="0">
			<raterules />
			<options>
				<option code="CDW" rate="0" firstfree="false" inclusive="true" chargetype="I" prepaid="true" mandatory="true" calc="D" exc="70000" dep="70000" net="0" />
				<option code="FDW" rate="8400" firstfree="false" inclusive="false" chargetype="I" mandatory="false" calc="D" dep="0" net="6774" />
				<option code="THIRD" rate="0" firstfree="false" inclusive="true" chargetype="I" prepaid="true" mandatory="true" calc="D" net="0" />
				<option code="THW" rate="0" firstfree="false" inclusive="true" chargetype="I" prepaid="true" mandatory="true" calc="D" exc="70000" dep="70000" net="0" />
				<option code="ADD" rate="2800" firstfree="true" inclusive="false" chargetype="E" mandatory="false" calc="D" net="2258" />
				<option code="BS" rate="4900" firstfree="false" inclusive="false" chargetype="E" mandatory="false" calc="D" net="3952" />
				<option code="CS" rate="4900" firstfree="false" inclusive="false" chargetype="E" mandatory="false" calc="D" net="3952" />
				<option code="SNC" rate="5000" firstfree="false" inclusive="false" chargetype="E" mandatory="false" calc="D" net="4032" />
				<option code="YDR" rate="7001" firstfree="false" inclusive="false" chargetype="E" mandatory="false" calc="D" net="5646" />
			</options>
			<taxes tax1code="VAT" tax1rate="2400" tax1amount="2094" tax2code="LCT" tax2rate="0" tax2amount="0" />
		</category>
		<category availability="ONREQUEST" cat="Passenger Van" code="FAD" acriss="FVAR" imageurl="https://wheels-assets.s3.eu-central-1.amazonaws.com/10268/fleet/groupPhoto-FAD.jpg" pax="7" bags="0" doors="4" suitcases="1" model="Fiat Doblo 7s automatic or similar" totalrate="26128" baserate="26128" onewaycharge="0" ownet="0" outofhours="0" oohnet="0" outofoffice="0" ooonet="0" excess="95000" excessapplies="true" IncKlm="99999" unlimited="true" addklmrate="0" prepaidamount="26128" rentalamount="26128" rentalnet="21071" hourscharged="0" deliveryratenet="0" collectionratenet="0">
			<raterules />
			<options>
				<option code="CDW" rate="0" firstfree="false" inclusive="true" chargetype="I" prepaid="true" mandatory="true" calc="D" exc="95000" dep="95000" net="0" />
				<option code="FDW" rate="11200" firstfree="false" inclusive="false" chargetype="I" mandatory="false" calc="D" exc="0" dep="0" net="9032" />
				<option code="THIRD" rate="0" firstfree="false" inclusive="true" chargetype="I" prepaid="true" mandatory="true" calc="D" net="0" />
				<option code="THW" rate="0" firstfree="false" inclusive="true" chargetype="I" prepaid="true" mandatory="true" calc="D" exc="95000" dep="95000" net="0" />
				<option code="ADD" rate="2800" firstfree="true" inclusive="false" chargetype="E" mandatory="false" calc="D" net="2258" />
				<option code="BS" rate="4900" firstfree="false" inclusive="false" chargetype="E" mandatory="false" calc="D" net="3952" />
				<option code="CS" rate="4900" firstfree="false" inclusive="false" chargetype="E" mandatory="false" calc="D" net="3952" />
				<option code="SNC" rate="5000" firstfree="false" inclusive="false" chargetype="E" mandatory="false" calc="D" net="4032" />
			</options>
			<taxes tax1code="VAT" tax1rate="2400" tax1amount="5057" tax2code="LCT" tax2rate="0" tax2amount="0" />
		</category>
	</rates>
</pricequote>`

const avanceFixtureQuoteFDW0 = `<pricequote id="fd2c67a0-2fc2-4b7d-8616-c1657f1f2c38" validto="29/09/2026 05:27" duration="7" countrycode="GR" taxinclusive="true" currency="EUR" timetaken="46" graceperiod="240" owprepaid="false" oohprepaid="false" oooprepaid="false" fuelpolicy="SL">
	<rates>
		<category availability="AVAILABLE" cat="Passenger Cars" code="A" acriss="MBMR" imageurl="https://wheels-assets.s3.eu-central-1.amazonaws.com/10268/fleet/groupPhoto-A.jpg" pax="4" bags="0" doors="4" suitcases="1" model="Skoda Citigo or similar" totalrate="20140" baserate="20140" onewaycharge="0" ownet="0" outofhours="0" oohnet="0" outofoffice="0" ooonet="0" excess="0" excessapplies="true" IncKlm="99999" unlimited="true" addklmrate="0" prepaidamount="20140" rentalamount="11740" rentalnet="9468" hourscharged="0" deliveryratenet="0" collectionratenet="0">
			<raterules />
			<options>
				<option code="CDW" rate="0" firstfree="false" inclusive="true" chargetype="I" prepaid="true" mandatory="true" calc="D" exc="70000" net="0" />
				<option code="FDW" rate="8400" firstfree="false" inclusive="true" chargetype="I" prepaid="true" mandatory="true" calc="D" exc="0" dep="15000" net="6774" />
				<option code="THIRD" rate="0" firstfree="false" inclusive="true" chargetype="I" prepaid="true" mandatory="true" calc="D" net="0" />
				<option code="THW" rate="0" firstfree="false" inclusive="true" chargetype="I" prepaid="true" mandatory="true" calc="D" exc="70000" net="0" />
				<option code="ADD" rate="2800" firstfree="true" inclusive="false" chargetype="E" mandatory="false" calc="D" net="2258" />
				<option code="BS" rate="4900" firstfree="false" inclusive="false" chargetype="E" mandatory="false" calc="D" net="3952" />
				<option code="CS" rate="4900" firstfree="false" inclusive="false" chargetype="E" mandatory="false" calc="D" net="3952" />
				<option code="SNC" rate="5000" firstfree="false" inclusive="false" chargetype="E" mandatory="false" calc="D" net="4032" />
				<option code="YDR" rate="7001" firstfree="false" inclusive="false" chargetype="E" mandatory="false" calc="D" net="5646" />
			</options>
			<taxes tax1code="VAT" tax1rate="2400" tax1amount="3898" tax2code="LCT" tax2rate="0" tax2amount="0" />
		</category>
		<category availability="AVAILABLE" cat="Passenger Cars" code="D" acriss="HDMR" imageurl="https://wheels-assets.s3.eu-central-1.amazonaws.com/10268/fleet/groupPhoto-D.jpg" pax="5" bags="0" doors="4" suitcases="1" model="Toyota Auris or similar" totalrate="24964" baserate="24964" onewaycharge="0" ownet="0" outofhours="0" oohnet="0" outofoffice="0" ooonet="0" excess="0" excessapplies="true" IncKlm="99999" unlimited="true" addklmrate="0" prepaidamount="24964" rentalamount="15164" rentalnet="12229" hourscharged="0" deliveryratenet="0" collectionratenet="0">
			<raterules />
			<options>
				<option code="CDW" rate="0" firstfree="false" inclusive="true" chargetype="I" prepaid="true" mandatory="true" calc="D" exc="80000" net="0" />
				<option code="FDW" rate="9800" firstfree="false" inclusive="true" chargetype="I" prepaid="true" mandatory="true" calc="D" exc="0" dep="15000" net="7903" />
				<option code="THIRD" rate="0" firstfree="false" inclusive="true" chargetype="I" prepaid="true" mandatory="true" calc="D" net="0" />
				<option code="THW" rate="0" firstfree="false" inclusive="true" chargetype="I" prepaid="true" mandatory="true" calc="D" exc="80000" net="0" />
				<option code="ADD" rate="2800" firstfree="true" inclusive="false" chargetype="E" mandatory="false" calc="D" net="2258" />
				<option code="BS" rate="4900" firstfree="false" inclusive="false" chargetype="E" mandatory="false" calc="D" net="3952" />
				<option code="CS" rate="4900" firstfree="false" inclusive="false" chargetype="E" mandatory="false" calc="D" net="3952" />
				<option code="SNC" rate="5000" firstfree="false" inclusive="false" chargetype="E" mandatory="false" calc="D" net="4032" />
				<option code="YDR" rate="7001" firstfree="false" inclusive="false" chargetype="E" mandatory="false" calc="D" net="5646" />
			</options>
			<taxes tax1code="VAT" tax1rate="2400" tax1amount="4832" tax2code="LCT" tax2rate="0" tax2amount="0" />
		</category>
		<category availability="AVAILABLE" cat="SUV" code="E-GQ" acriss="DFAE" imageurl="https://wheels-assets.s3.eu-central-1.amazonaws.com/10268/fleet/groupPhoto-E-GQ.jpg" pax="5" bags="0" doors="5" suitcases="2" model="Geely EX5 Electric or similar" totalrate="31440" baserate="31440" onewaycharge="0" ownet="0" outofhours="0" oohnet="0" outofoffice="0" ooonet="0" excess="0" excessapplies="true" IncKlm="99999" unlimited="true" addklmrate="0" prepaidamount="31440" rentalamount="20241" rentalnet="16323" hourscharged="0" deliveryratenet="0" collectionratenet="0">
			<raterules />
			<options>
				<option code="CDW" rate="0" firstfree="false" inclusive="true" chargetype="I" prepaid="true" mandatory="true" calc="D" exc="95000" net="0" />
				<option code="FDW" rate="11200" firstfree="false" inclusive="true" chargetype="I" prepaid="true" mandatory="true" calc="D" exc="0" dep="15000" net="9032" />
				<option code="THIRD" rate="0" firstfree="false" inclusive="true" chargetype="I" prepaid="true" mandatory="true" calc="D" net="0" />
				<option code="THW" rate="0" firstfree="false" inclusive="true" chargetype="I" prepaid="true" mandatory="true" calc="D" exc="95000" dep="0" net="0" />
				<option code="ADD" rate="2800" firstfree="true" inclusive="false" chargetype="E" mandatory="false" calc="D" net="2258" />
				<option code="BS" rate="4900" firstfree="false" inclusive="false" chargetype="E" mandatory="false" calc="D" net="3952" />
				<option code="CS" rate="4900" firstfree="false" inclusive="false" chargetype="E" mandatory="false" calc="D" net="3952" />
				<option code="SNC" rate="5000" firstfree="false" inclusive="false" chargetype="E" mandatory="false" calc="D" net="4032" />
			</options>
			<taxes tax1code="VAT" tax1rate="2400" tax1amount="6085" tax2code="LCT" tax2rate="0" tax2amount="0" />
		</category>
		<category availability="AVAILABLE" cat="Passenger Cars" code="GM3" acriss="MDAR" imageurl="https://wheels-assets.s3.eu-central-1.amazonaws.com/10268/fleet/groupPhoto-GM3.jpg" pax="4" bags="0" doors="4" suitcases="1" model="Peugeot 108 automatic or similar" totalrate="22966" baserate="22966" onewaycharge="0" ownet="0" outofhours="0" oohnet="0" outofoffice="0" ooonet="0" excess="0" excessapplies="true" IncKlm="99999" unlimited="true" addklmrate="0" prepaidamount="22966" rentalamount="14566" rentalnet="11747" hourscharged="0" deliveryratenet="0" collectionratenet="0">
			<raterules />
			<options>
				<option code="CDW" rate="0" firstfree="false" inclusive="true" chargetype="I" prepaid="true" mandatory="true" calc="D" exc="70000" net="0" />
				<option code="FDW" rate="8400" firstfree="false" inclusive="true" chargetype="I" prepaid="true" mandatory="true" calc="D" dep="15000" net="6774" />
				<option code="THIRD" rate="0" firstfree="false" inclusive="true" chargetype="I" prepaid="true" mandatory="true" calc="D" net="0" />
				<option code="THW" rate="0" firstfree="false" inclusive="true" chargetype="I" prepaid="true" mandatory="true" calc="D" exc="70000" net="0" />
				<option code="ADD" rate="2800" firstfree="true" inclusive="false" chargetype="E" mandatory="false" calc="D" net="2258" />
				<option code="BS" rate="4900" firstfree="false" inclusive="false" chargetype="E" mandatory="false" calc="D" net="3952" />
				<option code="CS" rate="4900" firstfree="false" inclusive="false" chargetype="E" mandatory="false" calc="D" net="3952" />
				<option code="SNC" rate="5000" firstfree="false" inclusive="false" chargetype="E" mandatory="false" calc="D" net="4032" />
				<option code="YDR" rate="7001" firstfree="false" inclusive="false" chargetype="E" mandatory="false" calc="D" net="5646" />
			</options>
			<taxes tax1code="VAT" tax1rate="2400" tax1amount="4445" tax2code="LCT" tax2rate="0" tax2amount="0" />
		</category>
		<category availability="ONREQUEST" cat="Passenger Van" code="FAD" acriss="FVAR" imageurl="https://wheels-assets.s3.eu-central-1.amazonaws.com/10268/fleet/groupPhoto-FAD.jpg" pax="7" bags="0" doors="4" suitcases="1" model="Fiat Doblo 7s automatic or similar" totalrate="44492" baserate="44492" onewaycharge="0" ownet="0" outofhours="0" oohnet="0" outofoffice="0" ooonet="0" excess="0" excessapplies="true" IncKlm="99999" unlimited="true" addklmrate="0" prepaidamount="44492" rentalamount="33293" rentalnet="26849" hourscharged="0" deliveryratenet="0" collectionratenet="0">
			<raterules />
			<options>
				<option code="CDW" rate="0" firstfree="false" inclusive="true" chargetype="I" prepaid="true" mandatory="true" calc="D" exc="95000" net="0" />
				<option code="FDW" rate="11200" firstfree="false" inclusive="true" chargetype="I" prepaid="true" mandatory="true" calc="D" exc="0" dep="15000" net="9032" />
				<option code="THIRD" rate="0" firstfree="false" inclusive="true" chargetype="I" prepaid="true" mandatory="true" calc="D" net="0" />
				<option code="THW" rate="0" firstfree="false" inclusive="true" chargetype="I" prepaid="true" mandatory="true" calc="D" exc="95000" net="0" />
				<option code="ADD" rate="2800" firstfree="true" inclusive="false" chargetype="E" mandatory="false" calc="D" net="2258" />
				<option code="BS" rate="4900" firstfree="false" inclusive="false" chargetype="E" mandatory="false" calc="D" net="3952" />
				<option code="CS" rate="4900" firstfree="false" inclusive="false" chargetype="E" mandatory="false" calc="D" net="3952" />
				<option code="SNC" rate="5000" firstfree="false" inclusive="false" chargetype="E" mandatory="false" calc="D" net="4032" />
			</options>
			<taxes tax1code="VAT" tax1rate="2400" tax1amount="8611" tax2code="LCT" tax2rate="0" tax2amount="0" />
		</category>
	</rates>
</pricequote>`

const avanceFixtureQuoteYoung = `<pricequote id="78672682-d65c-4909-9833-ca2e71e15076" validto="29/09/2026 05:27" duration="7" countrycode="GR" taxinclusive="true" currency="EUR" timetaken="46" graceperiod="240" owprepaid="false" oohprepaid="false" oooprepaid="false" fuelpolicy="SL">
	<rates>
		<category availability="AVAILABLE" cat="Passenger Cars" code="A" acriss="MBMR" imageurl="https://wheels-assets.s3.eu-central-1.amazonaws.com/10268/fleet/groupPhoto-A.jpg" pax="4" bags="0" doors="4" suitcases="1" model="Skoda Citigo or similar" totalrate="14767" baserate="7766" onewaycharge="0" ownet="0" outofhours="0" oohnet="0" outofoffice="0" ooonet="0" excess="70000" excessapplies="true" IncKlm="99999" unlimited="true" addklmrate="0" prepaidamount="7766" rentalamount="7766" rentalnet="6263" hourscharged="0" deliveryratenet="0" collectionratenet="0">
			<raterules />
			<options>
				<option code="CDW" rate="0" firstfree="false" inclusive="true" chargetype="I" prepaid="true" mandatory="true" calc="D" exc="70000" dep="70000" net="0" />
				<option code="FDW" rate="8400" firstfree="false" inclusive="false" chargetype="I" mandatory="false" calc="D" exc="0" dep="0" net="6774" />
				<option code="THIRD" rate="0" firstfree="false" inclusive="true" chargetype="I" prepaid="true" mandatory="true" calc="D" net="0" />
				<option code="THW" rate="0" firstfree="false" inclusive="true" chargetype="I" prepaid="true" mandatory="true" calc="D" exc="70000" dep="70000" net="0" />
				<option code="ADD" rate="2800" firstfree="true" inclusive="false" chargetype="E" mandatory="false" calc="D" net="2258" />
				<option code="BS" rate="4900" firstfree="false" inclusive="false" chargetype="E" mandatory="false" calc="D" net="3952" />
				<option code="CS" rate="4900" firstfree="false" inclusive="false" chargetype="E" mandatory="false" calc="D" net="3952" />
				<option code="SNC" rate="5000" firstfree="false" inclusive="false" chargetype="E" mandatory="false" calc="D" net="4032" />
				<option code="YDR" rate="7001" firstfree="false" inclusive="false" chargetype="E" mandatory="true" calc="D" net="5646" />
			</options>
			<taxes tax1code="VAT" tax1rate="2400" tax1amount="2858" tax2code="LCT" tax2rate="0" tax2amount="0" />
		</category>
	</rates>
</pricequote>`

const avanceFixtureQuoteError = `<pricequote id="d28acf02-0928-41bd-82c9-1940227b02c3" validto="29/09/2026 05:27" duration="0" taxinclusive="true" currency="EUR" timetaken="0" graceperiod="0" owprepaid="false" oohprepaid="false" oooprepaid="false">
	<errors>
		<error code="ERR/102">Pickup Station not found</error>
	</errors>
	<rates />
</pricequote>`

const avanceFixtureStations = `<response>
	<station code="022" name="Athens International Airport" lat="37.936764" long="23.946270" icon="fa-plane" country="GR" delradius="0">
		<StationInformation StationType="Airport" StationTypeOTA="1" StationTypeDesc="Both counter and car in terminal" Active="true" TimeZone="GTB Standard Time" NonShowGraceHours="0" Tax1Percent="24.00" Tax2Percent="0" Tax1Name="VAT" Tax2Name="LCT" Address="Desk in Arrivals Hall 30" ZipCode="19019" City="Spata1" Phone="+30 2103533088 ">
			<PickupInstructions>
				<PickupInfo>Our office is located in the Arrivals Hall of Athens International Airport</PickupInfo>
			</PickupInstructions>
			<OperationHours>
				<WorkHour Day="Monday" Closed="false" OpensAt="00:00" ClosesAt="23:59" />
				<WorkHour Day="Tuesday" Closed="false" OpensAt="00:00" ClosesAt="23:59" />
				<WorkHour Day="Wednesday" Closed="false" OpensAt="00:00" ClosesAt="23:59" />
				<WorkHour Day="Thursday" Closed="false" OpensAt="00:00" ClosesAt="23:59" />
				<WorkHour Day="Friday" Closed="false" OpensAt="00:00" ClosesAt="23:59" />
				<WorkHour Day="Saturday" Closed="false" OpensAt="00:00" ClosesAt="23:59" />
				<WorkHour Day="Sunday" Closed="false" OpensAt="00:00" ClosesAt="23:59" />
			</OperationHours>
		</StationInformation>
	</station>
</response>`

const avanceFixtureOptions = `<response>
	<option code="ADD" name="Additional Driver" quant="true" />
	<option code="BS" name="Booster Seat" quant="true" />
	<option code="CDW" name="Collision Damage Waiver" quant="false" />
	<option code="CS" name="Baby Seat" quant="true" />
	<option code="FDW" name="Full Damage Waiver" quant="false" />
	<option code="SNC" name="Snow Chains" quant="false" />
	<option code="THIRD" name="Third Party Insurance" quant="false" />
	<option code="THW" name="Theft Waiver with Excess" quant="false" />
	<option code="YDR" name="Young Driver Fee 19-22" quant="false" />
</response>`

// Live new-res responses from the September test bookings: a confirmed booking, and one that came
// back on request.
const avanceFixtureNewResConfirmed = `<?xml version="1.0" encoding="utf-8"?><response xmlns:xsd="http://www.w3.org/2001/XMLSchema" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance"><reservation irn="5RLDDC" status="OK" refno="AIBT-TEST-005" res-status="RES" checkouturl="https://checkout.wheelsys.ms/10268/5RLDDC/AIBOOKINGTRAVEL/20261015/" /></response>`

const avanceFixtureNewResOnRequest = `<?xml version="1.0" encoding="utf-8"?><response xmlns:xsd="http://www.w3.org/2001/XMLSchema" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance"><reservation irn="5NS061" status="OK" refno="AIBT-TEST-001" res-status="REQ" /></response>`

// The cancel-res response as documented by Wheelsys.
const avanceFixtureCancelRes = `<?xml version="1.0" encoding="utf-8" ?>
<response>
	<reservation irn="3QW187" status="OK" res-status="CNC" ></reservation>
</response>`

// Error shapes: the documentation puts the ERR code in status, while the quote reports errors in an
// <errors> block, so both are covered.
const avanceFixtureResStatusError = `<?xml version="1.0" encoding="utf-8"?><response><reservation status="ERR/111" /></response>`

const avanceFixtureResErrorsBlock = `<?xml version="1.0" encoding="utf-8"?><response><errors><error code="ERR/110">Price Quote Not Found or has Expired</error></errors></response>`

const avanceFixtureCancelAlreadyCancelled = `<?xml version="1.0" encoding="utf-8"?><response><reservation irn="5RLDDC" status="ERR/104" /></response>`
