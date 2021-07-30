module TimeSoft-OA/Server

go 1.16

require (
	TimeSoft-OA/SocketPacket v0.0.0
	github.com/mattn/go-sqlite3 v1.14.8 // indirect
	golang.org/x/text v0.3.6 // indirect
)

replace TimeSoft-OA/SocketPacket v0.0.0 => ../SocketPacket
