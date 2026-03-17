.PHONY: gen

gen:
	@echo "Generating sqlc files for auth service..."
	cd services/auth/db && sqlc generate
	@echo "Generating mocks for auth service..."
	mockgen -source=services/auth/db/querier.go -destination=services/auth/mocks/querier_mock.go -package=mocks
	@echo "Generating sqlc files for locations service..."
	cd services/booking/db && sqlc generate
	@echo "Generating mocks for locations service..."
	mockgen -source=services/booking/db/querier.go -destination=services/booking/mocks/querier_mock.go -package=locations_mocks
	@echo "Done!"
