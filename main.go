package main

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/coreos/etcd/clientv3"
)

func bootstrap(target string) {
	for {
		cli, err := clientv3.New(clientv3.Config{
			Endpoints:   []string{target},
			DialTimeout: 1 * time.Second,
		})
		if err != nil {
			log.Println("Connection failed to", target)
			time.Sleep(time.Second)
			continue
		} else {
			log.Println("Connected! Bootstrapping data...")
			time.Sleep(10 * time.Second)
			//cli.Delete(context.Background(), "magic")
			cli.Put(context.Background(), "magic", "42")
			cli.Close()
			break
		}
	}
}

func work(target string) {
	for {
		cli, err := clientv3.New(clientv3.Config{
			Endpoints:   []string{target},
			DialTimeout: 1 * time.Second,
		})
		if err != nil {
			log.Println("Connection failed to", target)
			time.Sleep(time.Second)
			continue
		} else {
			rch := cli.Watch(context.Background(), "magic")
			for wresp := range rch {
				for _, ev := range wresp.Events {
					log.Printf("%s %q : %q\n", ev.Type, ev.Kv.Key, ev.Kv.Value)
				}
			}
			cli.Close()
		}
	}
}

func watchMembers(target string) {
	lastMemberAmount := -1
	for {
		cli, err := clientv3.New(clientv3.Config{
			Endpoints:   []string{target},
			DialTimeout: 1 * time.Second,
		})
		if err != nil {
			log.Println("Connection failed to", target)
			lastMemberAmount = -1
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
			time.Sleep(time.Second)
		}
	}
}

func main() {
	target := os.Getenv("ETCD_ENDPOINT")
	if target == "" {
		log.Fatalf("ETCD_ENDPOINT is not set!")
	}
	if len(os.Args) == 2 && os.Args[1] == "bootstrap" {
		bootstrap(target)
		os.Exit(0)
	} else {
		log.Println("Etcd healthchecker started up! ETCD_ENDPOINT:", target)
		go watchMembers(target)
		work(target)
	}
}
