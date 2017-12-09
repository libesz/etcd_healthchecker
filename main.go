package main

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/coreos/etcd/clientv3"
)

func main() {
	target := os.Getenv("ETCD_ENDPOINT")
	if target == "" {
		log.Fatalf("ETCD_ENDPOINT is not set!")
	}
	log.Println("Etcd healthchecker started up! ETCD_ENDPOINT:", target)
	lastMemberAmount := -1
	for {
		cli, err := clientv3.New(clientv3.Config{
			Endpoints:   []string{target},
			DialTimeout: 1 * time.Second,
		})
		if err != nil {
			log.Println("Connection failed to", target)
		} else {
			resp, err := cli.MemberList(context.Background())
			if err != nil {
				log.Println("Failed to get memberlist")
			}
			if lastMemberAmount != len(resp.Members) {
				log.Println("Member amount changed to:", len(resp.Members))
				lastMemberAmount = len(resp.Members)
			}
			cli.Close()
		}
		time.Sleep(time.Second)
	}
}
