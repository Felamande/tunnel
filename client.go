package main

import (
	"io"
	"net"

	"encoding/binary"

	"github.com/codegangsta/cli"
	"github.com/qiniu/log"
)

func clientCmd() cli.Command {
	return cli.Command{
		Name:   "client",
		Action: ClientAction,
		Flags: []cli.Flag{
			cli.StringFlag{
				Name: "connect, c",
			},
			cli.StringFlag{
				Name: "local-port",
			},
			cli.IntFlag{
				Name: "remote-port",
			},
		},
	}
}

func ClientAction(c *cli.Context) error {
	cAddr, err := net.ResolveTCPAddr("tcp", c.String("connect"))
	if err != nil {
		log.Error(err)
		return nil
	}
	lAddr, err := net.ResolveTCPAddr("tcp", "localhost:"+c.String("local-port"))
	if err != nil {
		log.Error(err)
		return nil
	}

	pconn, err := net.DialTCP("tcp", nil, cAddr)
	if err != nil {
		log.Error(err)
		return nil
	}
	log.Info("start dialing server", cAddr)
	portPacket := make([]byte, 5)
	portPacket[0] = SwitchIP

	binary.BigEndian.PutUint32(portPacket[1:], uint32(c.Int("remote-port")))
	n, err := pconn.Write(portPacket)
	if err != nil {
		log.Error(err)
		return nil
	}
	log.Info("send port packet ok,length", n)
	b := make([]byte, 1)
	_, err = pconn.Read(b)
	if err != nil {
		log.Error(err)
		return nil
	}
	log.Info("receive answer", b)
	if len(b) == 0 {
		log.Info("switch port failed", "no response from server")
		return nil
	}
	if b[0] != SwitchIPOK {
		log.Info("switch port failed", "invalid flag", uint8(b[0]))
		return nil
	}

	for {
		lconn, err := net.DialTCP("tcp", nil, lAddr)
		if err != nil {
			log.Info(err)
			continue
		}
		go io.Copy(lconn, pconn)
		io.Copy(pconn, lconn)
	}

}
