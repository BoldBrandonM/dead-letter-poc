.PHONY: emit-all
emit-all:
	cd ./src/main && go run emit_all.go

.PHONY: receive-text
receive-text:
	cd ./src/main && go run receive_text.go $(arg)

.PHONY: receive-bytes
receive-bytes:
	cd ./src/main && go run receive_bytes.go $(arg)

.PHONY: receive-dead-letters
receive-dead-letters:
	cd ./src/main && go run receive_dead_letters.go
