package main

import (
	"io"
	"net"
	"strconv"

	"encoding/binary"

	"github.com/codegangsta/cli"
	"github.com/qiniu/log"
)

func serverCmd() cli.Command {
	return cli.Command{
		Name: "server",
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:  "listen, l",
				Value: "30000",
				Usage: "server port of persistent connection",
			},
		},
		Action: ServerAction,
	}
}

func ServerAction(c *cli.Context) error {
	serverPort := c.String("listen")
	if serverPort == "" {
		return &Error{"ServerAction", "Get flag -listen", "invalid port"}
	}
	serverAddr, err := net.ResolveTCPAddr("tcp", ":"+serverPort)
	if err != nil {
		return &Error{"ServerAction", "ResolveTCPAddr", err}
	}

	persistentListener, err := net.ListenTCP("tcp", serverAddr)
	if err != nil {
		return &Error{"ServerAction", "ListenTCP", err}
	}
	log.Info("start persistent listening at", serverAddr)
	for {
		pconn, err := persistentListener.Accept()
		if err != nil {
			return &Error{"ServerAction", "Accept", err}
		}
		log.Info("start persistent connection at", serverAddr)
		log.Info("start accepting port packet")

		ipSwitchPacket := make([]byte, 5)
		_, err = pconn.Read(ipSwitchPacket)
		if err != nil {
			log.Error(err)
			pconn.Close()
			continue
		}
		log.Info("receive port packet", ipSwitchPacket)

		if ipSwitchPacket[0] != SwitchIP {
			log.Error("switch ip error: invalid flag", ipSwitchPacket[0])
			pconn.Close()
			continue
		}
		intAddr := int(binary.BigEndian.Uint32(ipSwitchPacket[1:]))
		rport := strconv.Itoa(intAddr)

		rAddr, err := net.ResolveTCPAddr("tcp", ":"+rport)
		if err != nil {
			log.Error(err)
			pconn.Close()
			continue
		}
		_, err = pconn.Write([]byte{SwitchIPOK})
		if err != nil {
			log.Info("Get port for forwarding", rAddr)
		}
		l, err := net.ListenTCP("tcp", rAddr)
		if err != nil {
			log.Error(err)
			pconn.Close()
			continue
		}
		log.Info("Start listening at", rAddr)
		for {
			rconn, err := l.Accept()
			if err != nil {
				log.Error(err)
				continue
			}
			log.Info("forwarding connection from", rconn.RemoteAddr())
			go io.Copy(pconn, rconn)
			io.Copy(rconn, pconn)
		}
	}

	// return nil
}
