# Go-Client-Server
Go Server &amp; client application using TCP and UDP api

  To run servers:
                                  --Ubuntu/Debian Terminal--
  cd .../Go-Client-Server/server/
  go run serverTCP.go //for TCP server
  go run serverUDP.go //for UDP server

  To run Clients:
                                  --Ubuntu/Debian Terminal--
  cd .../Go-Client-Server/client/Shell\ Scripts/
  chmod +x run_clientsTCP.sh && chmod +x run_clientsUDP.sh //Maybe necessary
  ./run_clientsTCP.sh //for TCP clients
  ./run_clientsUDP.sh //for UDP clients

  To storage logs:
    Each client(TCP or UDP) adds their average execution time to a txt file(log-meanTime-TCPClients.txt and log-meanTime-UDPClients.txt)
