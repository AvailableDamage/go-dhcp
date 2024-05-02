package main

import (
	"encoding/json"
	"log"
	"net"
	"os"
)

var confPath = "/etc/goDHCP/conf.json"

type option struct {
	id 			byte
	length 		byte
	value 		[]byte
}

type configuration struct {
	Interface  		string `json:"Interface"`
	Nameserver 		string `json:"Nameserver"`
	Gateway    		string `json:"Gateway"`
	LeaseTime  		int    `json:"LeaseTime"`
	ServerID 		string `json:"Server ID"`
	Pools []struct {
		Name    	string `json:"Name"`
		Interface 	string `json:Interface`
		StartIP 	net.IP `json:"StartIP"`
		EndIP   	net.IP `json:"EndIP"`
		Gateway 	string `json:Gateway`
		Nameserver 	string `json:Nameserver`
		Options 	map[string]interface{} `json:Options`
	} `json:"Pools"`
}

func loadconf() (configuration){

	file,err := os.Open(confPath)
	defer file.Close()

	log.Println("Loading config ", confPath)
	if err != nil {
		log.Panicln("Cant open configuration: ", err)
	}
	config := configuration{}
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&config)

	if err != nil {
		log.Panicln("Cant read configuration: ", err)
	}

	for i := range config.Pools {
		if config.Pools[i].Gateway == "" {
			config.Pools[i].Gateway = config.Gateway
		}

		if config.Pools[i].Nameserver == "" {
			config.Pools[i].Nameserver = config.Nameserver
		}
		if config.Pools[i].Interface == "" {
			config.Pools[i].Interface = config.Interface
		}

	}

	return config
}

func getOptionValue(option optionID) ([]byte){
	log.Println("Get options value")
	for key,value := range DhcpConf.Pools[0].Options {
		if optionCodeToString[option] == key {
			return parseString(value.(string))
		}
	}
	return []byte{}
}
