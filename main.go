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
			defer cli.Close()
			resp, err := cli.MemberList(context.Background())
			if err != nil {
				log.Println("Failed to get memberlist")
			}
			if lastMemberAmount != len(resp.Members) {
				lastMemberAmount = len(resp.Members)
				memberList := ""
				for _, element := range resp.Members {
					memberList += element.GetName()
					memberList += ", "
				}
				log.Println("Memberlist changed:", memberList)
			}
		}
		time.Sleep(time.Second)
	}
}
