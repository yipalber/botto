package main

import (
	"fmt"
	"log"
	"os"

	pb "github.com/barcatfigaro/botto/vessel-service/proto/vessel"
	micro "github.com/micro/go-micro"
)

const (
	defaultHost = "localhost:27017"
)

func createDummyData(repo Repository) {
	defer repo.Close()
	vessels := []*pb.Vessel{
		{Id: "vessel001", Name: "Kane's Salty Secret", MaxWeight: 200000, Capacity: 500},
	}
	for _, v := range vessels {
		repo.Create(v)
	}
}

func main() {
	host := os.Getenv("DB_HOST")
	if host == "" {
		host = defaultHost
	}

	session, err := CreateSession(host)
	defer session.Close()

	if err != nil {
		log.Fatalf("error connecting to datastore: %v\n", err)
	}

	repo := &VesselRepository{session.Copy()}

	createDummyData(repo)

	srv := micro.NewService(
		micro.Name("go.micro.srv.vessel"),
		micro.Version("latest"),
	)

	srv.Init()

	pb.RegisterVesselServiceHandler(srv.Server(), &service{session})

	if err := srv.Run(); err != nil {
		fmt.Println(err)
	}
}
