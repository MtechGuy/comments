## run: run the cmd/api application
.PHONY: run
run: 
	@echo 'Running application...'
	@go run ./cmd/api -port=4000 -env=production
