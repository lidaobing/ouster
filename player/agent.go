package player

import (
	"github.com/tiancaiamao/ouster/packet"
	"log"
)

func (this *Player) loop() {
	var msg interface{}
	for {
		select {
		case msg, ok := <-this.client:
			if !ok {
				// kick the player off...
				return
			} else {
				this.handleClientMessage(msg)
			}
		case msg = <-this.Scene2player:
			this.handleSceneMessage(msg)
		case <-this.aoi:
			// 来自aoi的消息
		case <-this.heartbeat:
			this.heartBeat()
		}
	}
}

func (player *Player) Go() {
	read := make(chan interface{}, 1)
	write := make(chan packet.Packet, 1)
	player.send = write
	player.client = read

	// open a goroutine to read from conn
	go func() {
		for {
			data, err := packet.Read(player.conn)
			if err != nil {
				log.Println(err)
				player.conn.Close()
				close(read)
				return
			}
			// log.Println("packet before send to chan", data)
			read <- data
			// log.Println("packet after send to chan", data)
		}
	}()

	// open a goroutine to write to conn
	go func() {
		for {
			pkt := <-write
			err := packet.Write(player.conn, pkt.Id, pkt.Obj)
			if err != nil {
				log.Println(err)
				continue
			}
		}
	}()

	player.loop()
}
