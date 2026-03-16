module github.com/zackyyhw/latihan-datacontract/domain-order

go 1.25.0

require (
	github.com/lib/pq v1.11.2
	github.com/zackyyhw/latihan-datacontract/contracts/order/v1 v0.0.0
	golang.org/x/net v0.48.0 // indirect
	golang.org/x/sys v0.39.0 // indirect
	golang.org/x/text v0.32.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20251202230838-ff82c1b0f217 // indirect
	google.golang.org/grpc v1.79.2
	google.golang.org/protobuf v1.36.10 // indirect
)

replace github.com/zackyyhw/latihan-datacontract/contracts/order/v1 => ../contracts/order/v1