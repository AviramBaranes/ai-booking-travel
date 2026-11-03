package broker

// Autohellas availability responses for the tests. We have no UAT access yet, so these are the
// samples in Autohellas' implementation guide (version 2.2.0, section 4.1.1), trimmed to a few cars
// and wrapped in the SOAP envelope the WSDL describes. Changes from the guide's same-station
// sample, so the tests cover what it leaves out: the first car carries passenger and baggage counts
// (as in the "Sample Requests" document) and an "or similar" model name, the second car is
// OnRequest, and an included "Airport Service" charge comes from the guide's reservation sample.
// Replace them with captured UAT responses once we have access.

// hertzGRSameStationFixture is the guide's 7-day rental at Athens Airport: two available cars, one
// automatic, and one OnRequest car.
const hertzGRSameStationFixture = `<?xml version="1.0" encoding="utf-8"?>
<s:Envelope xmlns:s="http://schemas.xmlsoap.org/soap/envelope/">
<s:Body>
<VehicleAvailRateResult xmlns="http://www.hertz.gr/xml">
<OTA_VehAvailRateRS TimeStamp="2014-04-16T00:08:00" Target="Test" Version="1.0" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance" xmlns="http://www.opentravel.org/OTA/2003/05">
	<Success/>
	<Warnings>
		<Warning Type="3" Code="221">Not valid car</Warning>
	</Warnings>
	<VehAvailRSCore>
		<VehRentalCore PickUpDateTime="2014-05-01T10:00:00" ReturnDateTime="2014-05-08T10:00:00">
			<PickUpLocation LocationCode="030">Athens - International Airport "El. Venizelos"</PickUpLocation>
			<ReturnLocation LocationCode="030">Athens - International Airport "El. Venizelos"</ReturnLocation>
		</VehRentalCore>
		<VehVendorAvails>
			<VehVendorAvail>
				<Vendor xsi:type="CompanyNameType" Code="ZE"/>
				<VehAvails>
					<VehAvail>
						<VehAvailCore Status="Available">
							<Vehicle AirConditionInd="true" TransmissionType="Manual" PassengerQuantity="4" BaggageQuantity="2" VendorCarType="A" Code="MDMR" CodeContext="SIPP">
								<VehType VehicleCategory="1" DoorCount="4"/>
								<VehClass Size="1"/>
								<VehMakeModel Name="CHEVROLET MATIZ or similar" Code="MDMR"/>
								<PictureURL>http://images.autohellas.gr/CarImages/i10.gif</PictureURL>
							</Vehicle>
							<RentalRate>
								<RateDistance Unlimited="true" DistUnitName="Km" VehiclePeriodUnitName="RentalPeriod"/>
								<VehicleCharges>
									<VehicleCharge TaxInclusive="true" Description="Rental Rate" IncludedInRate="true" IncludedInEstTotalInd="true" Purpose="1" CurrencyCode="EUR" Amount="123.00">
										<TaxAmounts>
											<TaxAmount Total="28.29" CurrencyCode="EUR" Percentage=" 23" Description="VAT"/>
										</TaxAmounts>
									</VehicleCharge>
									<VehicleCharge TaxInclusive="true" Description="Airport Service" IncludedInRate="true" IncludedInEstTotalInd="true" Purpose="5" CurrencyCode="EUR" Amount="0.00"/>
								</VehicleCharges>
								<RateQualifier RateQualifier="ITRH296"/>
							</RentalRate>
							<TotalCharge RateTotalAmount="123.00" EstimatedTotalAmount="123.00" CurrencyCode="EUR"/>
						</VehAvailCore>
						<VehAvailInfo/>
					</VehAvail>
					<VehAvail>
						<VehAvailCore Status="OnRequest">
							<Vehicle AirConditionInd="true" TransmissionType="Manual" VendorCarType="B" Code="ECMR" CodeContext="SIPP">
								<VehType VehicleCategory="1"/>
								<VehClass Size="3"/>
								<VehMakeModel Name="NISSAN MICRA" Code="ECMR"/>
								<PictureURL>http://images.autohellas.gr/CarImages/Thrifty/Micra.jpg</PictureURL>
							</Vehicle>
							<RentalRate>
								<RateDistance Unlimited="true" DistUnitName="Km" VehiclePeriodUnitName="RentalPeriod"/>
								<VehicleCharges>
									<VehicleCharge TaxInclusive="true" Description="Rental Rate" IncludedInRate="true" IncludedInEstTotalInd="true" Purpose="1" CurrencyCode="EUR" Amount="135.00"/>
								</VehicleCharges>
								<RateQualifier RateQualifier="ITRH296"/>
							</RentalRate>
							<TotalCharge RateTotalAmount="135.00" EstimatedTotalAmount="135.00" CurrencyCode="EUR"/>
						</VehAvailCore>
						<VehAvailInfo/>
					</VehAvail>
					<VehAvail>
						<VehAvailCore Status="Available">
							<Vehicle AirConditionInd="true" TransmissionType="Automatic" VendorCarType="K" Code="CDAR" CodeContext="SIPP">
								<VehType VehicleCategory="1" DoorCount="4"/>
								<VehClass Size="4"/>
								<VehMakeModel Name="HYUNDAI ACCENT AUTO" Code="CDAR"/>
								<PictureURL>http://images.autohellas.gr/CarImages/2013FleetUpdate/leon2013_P.jpg</PictureURL>
							</Vehicle>
							<RentalRate>
								<RateDistance Unlimited="true" DistUnitName="Km" VehiclePeriodUnitName="RentalPeriod"/>
								<VehicleCharges>
									<VehicleCharge TaxInclusive="true" Description="Rental Rate" IncludedInRate="true" IncludedInEstTotalInd="true" Purpose="1" CurrencyCode="EUR" Amount="210.00"/>
								</VehicleCharges>
								<RateQualifier RateQualifier="ITRH296"/>
							</RentalRate>
							<TotalCharge RateTotalAmount="210.00" EstimatedTotalAmount="210.00" CurrencyCode="EUR"/>
						</VehAvailCore>
						<VehAvailInfo/>
					</VehAvail>
				</VehAvails>
				<Info>
					<LocationDetails Name="International Airport “El. Venizelos”">
						<Address FormattedInd="false">
							<AddressLine>Athens - International Airport “El. Venizelos”</AddressLine>
							<CityName>Athens</CityName>
							<CountryName Code="GR"/>
						</Address>
						<Telephone PhoneTechType="1" PhoneNumber="210/3534.900"/>
						<AdditionalInfo>
							<ParkLocation Location="1"/>
							<CounterLocation Location="1"/>
						</AdditionalInfo>
					</LocationDetails>
				</Info>
			</VehVendorAvail>
		</VehVendorAvails>
	</VehAvailRSCore>
</OTA_VehAvailRateRS>
</VehicleAvailRateResult>
</s:Body>
</s:Envelope>`

// hertzGROneWayFixture is the guide's one-way rental from Athens to Thessaloniki Airport, whose
// drop-off fee is paid at the station. The guide's cars are all OnRequest, as one-way rentals
// mostly are; the tests flip them to Available.
const hertzGROneWayFixture = `<?xml version="1.0" encoding="utf-8"?>
<s:Envelope xmlns:s="http://schemas.xmlsoap.org/soap/envelope/">
<s:Body>
<VehicleAvailRateResult xmlns="http://www.hertz.gr/xml">
<OTA_VehAvailRateRS TimeStamp="2014-04-16T00:08:00" Target="Test" Version="1.0" xmlns="http://www.opentravel.org/OTA/2003/05">
	<Success/>
	<VehAvailRSCore>
		<VehRentalCore PickUpDateTime="2014-04-18T10:00:00" ReturnDateTime="2014-04-20T10:00:00">
			<PickUpLocation LocationCode="030">Athens - International Airport "El. Venizelos"</PickUpLocation>
			<ReturnLocation LocationCode="042">Thessaloniki - Airport</ReturnLocation>
		</VehRentalCore>
		<VehVendorAvails>
			<VehVendorAvail>
				<VehAvails>
					<VehAvail>
						<VehAvailCore Status="OnRequest">
							<Vehicle AirConditionInd="true" TransmissionType="Manual" VendorCarType="A" Code="MDMR" CodeContext="SIPP">
								<VehType VehicleCategory="1" DoorCount="4"/>
								<VehClass Size="1"/>
								<VehMakeModel Name="CHEVROLET MATIZ" Code="MDMR"/>
								<PictureURL>http://images.autohellas.gr/CarImages/i10.gif</PictureURL>
							</Vehicle>
							<RentalRate>
								<RateDistance Unlimited="true" DistUnitName="Km" VehiclePeriodUnitName="RentalPeriod"/>
								<VehicleCharges>
									<VehicleCharge TaxInclusive="true" Description="Drop Off" IncludedInRate="false" IncludedInEstTotalInd="true" Purpose="2" CurrencyCode="EUR" Amount="150.06"/>
									<VehicleCharge TaxInclusive="true" Description="Rental Rate" IncludedInRate="true" IncludedInEstTotalInd="true" Purpose="1" CurrencyCode="EUR" Amount="84.00"/>
								</VehicleCharges>
								<RateQualifier RateQualifier="ITRH296"/>
							</RentalRate>
							<TotalCharge RateTotalAmount="84.00" EstimatedTotalAmount="234.06" CurrencyCode="EUR"/>
							<Fees>
								<Fee TaxInclusive="true" Description="Pay on arrival amount" IncludedInRate="false" IncludedInEstTotalInd="true" Purpose="23" CurrencyCode="EUR" Amount="150.06"/>
							</Fees>
						</VehAvailCore>
						<VehAvailInfo/>
					</VehAvail>
				</VehAvails>
			</VehVendorAvail>
		</VehVendorAvails>
	</VehAvailRSCore>
</OTA_VehAvailRateRS>
</VehicleAvailRateResult>
</s:Body>
</s:Envelope>`

// hertzGRErrorFixture wraps an OTA error, which Autohellas returns inside an HTTP 200.
const hertzGRErrorFixture = `<?xml version="1.0" encoding="utf-8"?>
<s:Envelope xmlns:s="http://schemas.xmlsoap.org/soap/envelope/">
<s:Body>
<VehicleAvailRateResult xmlns="http://www.hertz.gr/xml">
<OTA_VehAvailRateRS Version="1.0" xmlns="http://www.opentravel.org/OTA/2003/05">
	<Errors>
		<Error Type="3" Code="CODE" ShortText="TEXT"/>
	</Errors>
</OTA_VehAvailRateRS>
</VehicleAvailRateResult>
</s:Body>
</s:Envelope>`

// hertzGRFaultFixture is a standard SOAP 1.1 fault, as a WCF service returns for a call it cannot
// process at all.
const hertzGRFaultFixture = `<?xml version="1.0" encoding="utf-8"?>
<s:Envelope xmlns:s="http://schemas.xmlsoap.org/soap/envelope/">
<s:Body>
<s:Fault>
	<faultcode>s:Client</faultcode>
	<faultstring xml:lang="en-US">Authentication failed</faultstring>
</s:Fault>
</s:Body>
</s:Envelope>`

// hertzGRReservationFixture is the guide's confirmed reservation (section 6.1), trimmed and
// wrapped in the SOAP envelope. Its grand total is 164.01 EUR.
const hertzGRReservationFixture = `<?xml version="1.0" encoding="utf-8"?>
<s:Envelope xmlns:s="http://schemas.xmlsoap.org/soap/envelope/">
<s:Body>
<VehicleReservationResult xmlns="http://www.hertz.gr/xml">
<OTA_VehResRS Target="Test" TimeStamp="2007-12-11T00:45:32" Version="1.0" xmlns="http://www.opentravel.org/OTA/2003/05">
	<Success/>
	<VehResRSCore ReservationStatus="Committed">
		<VehReservation>
			<Customer>
				<Primary>
					<PersonName>
						<GivenName>John</GivenName>
						<Surname>Last</Surname>
					</PersonName>
				</Primary>
			</Customer>
			<VehSegmentCore>
				<ConfID ID="GR400454846" Type="14"/>
				<VehRentalCore PickUpDateTime="2008-06-01T10:00:00" ReturnDateTime="2008-06-07T10:00:00">
					<PickUpLocation LocationCode="030">Athens - International Airport "El. Venizelos"</PickUpLocation>
					<ReturnLocation LocationCode="030">Athens - International Airport "El. Venizelos"</ReturnLocation>
				</VehRentalCore>
				<Vehicle AirConditionInd="true" TransmissionType="Manual">
					<VehType VehicleCategory="1"/>
					<VehClass Size="3"/>
					<VehMakeModel Code="ECMR" Name="HYUNDAI ATOS A/C"/>
				</Vehicle>
				<RentalRate>
					<RateDistance DistUnitName="Km" Unlimited="true" VehiclePeriodUnitName="RentalPeriod"/>
					<RateQualifier RateQualifier="ITRH000"/>
				</RentalRate>
				<TotalCharge CurrencyCode="EUR" EstimatedTotalAmount="164.01" RateTotalAmount="137.82"/>
				<TPA_Extensions>
					<VehGroup Code=" A"/>
				</TPA_Extensions>
			</VehSegmentCore>
		</VehReservation>
	</VehResRSCore>
</OTA_VehResRS>
</VehicleReservationResult>
</s:Body>
</s:Envelope>`

// hertzGRCancelFixture is the guide's cancellation response (section 12.0).
const hertzGRCancelFixture = `<?xml version="1.0" encoding="utf-8"?>
<s:Envelope xmlns:s="http://schemas.xmlsoap.org/soap/envelope/">
<s:Body>
<VehicleCancelReservationResult xmlns="http://www.hertz.gr/xml">
<OTA_VehCancelRS Target="Test" TimeStamp="2007-12-11T01:32:27" Version="1.0" xmlns="http://www.opentravel.org/OTA/2003/05">
	<Success/>
	<VehCancelRSCore CancelStatus="Cancelled">
		<UniqueID ID="GR400454846" Type="14"/>
	</VehCancelRSCore>
	<VehCancelRSInfo/>
</OTA_VehCancelRS>
</VehicleCancelReservationResult>
</s:Body>
</s:Envelope>`

// hertzGRCancelErrorFixture is a cancellation refused with an OTA error.
const hertzGRCancelErrorFixture = `<?xml version="1.0" encoding="utf-8"?>
<s:Envelope xmlns:s="http://schemas.xmlsoap.org/soap/envelope/">
<s:Body>
<VehicleCancelReservationResult xmlns="http://www.hertz.gr/xml">
<OTA_VehCancelRS Version="1.0" xmlns="http://www.opentravel.org/OTA/2003/05">
	<Errors>
		<Error Type="3" Code="CODE" ShortText="TEXT"/>
	</Errors>
</OTA_VehCancelRS>
</VehicleCancelReservationResult>
</s:Body>
</s:Envelope>`
