package config


type GRPCConfig struct{
	Host string
	Port string
}


func loadGRPC()GRPCConfig{
	return GRPCConfig{
		Host:getEnv("GRPC_HOST","0.0.0.0"),
		Port:getEnv("GRPC_PORT","50052"),
	}
}

func (c GRPCConfig) Addr()string{
	return c.Host + ":" + c.Port	
}



