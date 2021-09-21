.PHONY: run
run:
	go build -o app && \
	/bin/echo -n "Warn" > log_level && \
	export LOG_LEVEL_FILE="$(shell pwd)/log_level" && \
	./app

debug:
	/bin/echo -n "Debug" > log_level

info:
	/bin/echo -n "Info" > log_level

warn:
	/bin/echo -n "Warn" > log_level

error:
	/bin/echo -n "Error" > log_level

clean:
	rm app log_level
