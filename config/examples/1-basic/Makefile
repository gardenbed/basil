.PHONY: run
run:
	go build -o app && \
	/bin/echo -n "service-1:8080,service-2:8080,service-3:8080" > endpoints && \
	export ENDPOINTS_FILE=$(shell pwd)/endpoints && \
	export TIMEOUT=4m && \
	export LOG_LEVEL=info && \
	./app --port 3000 && \
	rm app endpoints
