package main

func main() {

	confs := []ConfigFunc{withInitWorkers(3), withMinWorkers(2), withPollPeriod(10), withMaxWorkers(5)}

	pool := GetPool(confs...)

	server := GetNewServer(&ServerConfig{
		Addr: ":3000",
	})

	server.Pool = pool

	if err := server.Serve(); err != nil {
		panic(err)
	}

}
