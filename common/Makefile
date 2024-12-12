gen:
	@protoc --go_out=. --go_opt=paths=source_relative \
		--go-grpc_out=. --go-grpc_opt=paths=source_relative \
		api/protobuf/oms.proto
	@protoc --go_out=. --go_opt=paths=source_relative \
		--go-grpc_out=. --go-grpc_opt=paths=source_relative \
		api/protobuf/coordinator.proto
	@avrogo -s "_event" -p avro -d ./api/avro ./api/avro/*.avsc