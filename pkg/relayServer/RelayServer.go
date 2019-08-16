package relayServer

type relayServer struct {

}

type RelayServer struct{
	relayServer
}

func CreateRelayServer() RelayServer{
	return RelayServer{

	}
}

func (r relayServer) Start(){
	println("Start")
}