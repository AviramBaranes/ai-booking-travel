.PHONY: gen

gen:
	@echo "Generating sqlc files for accounts service..."
	cd services/accounts/db && sqlc generate
	@echo "Generating mocks for accounts service..."
	mockgen -source=services/accounts/db/querier.go -destination=services/accounts/mocks/querier_mock.go -package=mocks
	@echo "Generating sqlc files for locations service..."
	cd services/booking/db && sqlc generate
	@echo "Generating mocks for locations service..."
	mockgen -source=services/booking/db/querier.go -destination=services/booking/mocks/querier_mock.go -package=locations_mocks
	@echo "Generating sqlc files for reservation service..."
	cd services/reservation/db && sqlc generate
	@echo "Generating mocks for reservation service..."
	mockgen -source=services/reservation/db/querier.go -destination=services/reservation/mocks/querier_mock.go -package=mocks
	@echo "Done!"
