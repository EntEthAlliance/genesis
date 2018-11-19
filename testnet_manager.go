package main

import (
	"fmt"
	"log"
	"errors"
	db "./db"
	deploy "./deploy"
	util "./util"
	eos "./blockchains/eos"
	eth "./blockchains/ethereum"
	state "./state"
)

type DeploymentDetails struct {
	Servers    []int
	Blockchain string
	Nodes      int
	Image      string
}


func AddTestNet(details DeploymentDetails) error {
	defer state.DoneBuilding()
	servers, err := db.GetServers(details.Servers)

	if err != nil {
		log.Println(err.Error())
		state.BuildError = err
		return err
	}
	fmt.Println("Got the Servers")
	
	config := deploy.Config{Nodes: details.Nodes, Image: details.Image, Servers: details.Servers}
	fmt.Printf("Created the build configuration : %+v \n",config)

	newServerData := deploy.Build(&config, servers) //TODO: Restructure distribution of nodes over servers
	fmt.Println("Built the docker containers")

	
	switch(details.Blockchain){
		case "eos":
			eos.Eos(details.Nodes,newServerData);
		case "ethereum":
			eth.Ethereum(4000000,15468,15468,details.Nodes,newServerData)
		default:
			state.BuildError = errors.New("Unknown blockchain")
			return errors.New("Unknown blockchain")
	}

	testNetId := db.InsertTestNet(db.TestNet{Id: -1, Blockchain: details.Blockchain, Nodes: details.Nodes, Image: details.Image})

	for _, server := range newServerData {
		db.UpdateServerNodes(server.Id,0)
		for i, ip := range server.Ips {
			node := db.Node{Id: -1, TestNetId: testNetId, Server: server.Id, LocalId: i, Ip: ip} 
			_,err := db.InsertNode(node)
			if err != nil {
				log.Println(err.Error())
			}
		}
	}
	return nil
}

func GetLastTestNetId() int {
	testNets := db.GetAllTestNets()
	highestId := -1

	for _, testNet := range testNets {
		if testNet.Id > highestId {
			highestId = testNet.Id
		}
	}
	return highestId
}

func GetNextTestNetId() string {
	highestId := GetLastTestNetId()
	return fmt.Sprintf("%d",highestId+1)
}

func RebuildTestNet(id int) {
	panic("Not Implemented")
}

func RemoveTestNet(id int) error {
	nodes, err := db.GetAllNodesByTestNet(id)
	if err != nil {
		return err
	}
	for _, node := range nodes {
		server, _, err := db.GetServer(node.Server)
		if err != nil {
			return err
		}
		util.SshExec(server.Addr, fmt.Sprintf("~/local_deploy/deploy --kill=%d", node.LocalId))
	}
	return nil
}
