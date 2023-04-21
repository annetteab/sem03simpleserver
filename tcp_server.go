package main

import (
	"io"
	"log"
	"net"
	"sync"
	"github.com/annetteab/is105sem03/mycrypt"
)

func main() {

	var wg sync.WaitGroup

	server, err := net.Listen("tcp", "172.17.0.3:8080")
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("bundet til %s", server.Addr().String())
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			log.Println("før server.Accept() kallet")
			conn, err := server.Accept()
			if err != nil {
				return
			}
			go func(c net.Conn) {
				defer c.Close()
				for {
					buf := make([]byte, 1024)
					n, err := c.Read(buf)
					if err != nil {
						if err != io.EOF {
							log.Println(err)
						}
						return // fra for løkke
					}

					log.Println("Melding fra klient mottatt: ", string(buf[:n]))
                                        decryptedMsg := mycrypt.Krypter([]rune(string(buf[:n])), mycrypt.ALF_SEM03, len(mycrypt.ALF_SEM03)-4)
                                        log.Println("Melding fra klient dekryptert: ", string(decryptedMsg))

                                        switch string(decryptedMsg) {
                                        case "ping":
                                            resp := mycrypt.Krypter([]rune("pong"), mycrypt.ALF_SEM03, 4)
                                            log.Println("Svar: ", string(resp))
                                            _, err = c.Write([]byte(string(resp)))

					default:
                                                resp := mycrypt.Krypter([]rune(string(decryptedMsg)), mycrypt.ALF_SEM03, 4)
                                                log.Println("Svar: ", string(resp))
                                                _, err = c.Write([]byte(string(resp)))

					}
					if err != nil {
						if err != io.EOF {
							log.Println(err)
						}
						return // fra for løkke
					}
				}
			}(conn)
		}
	}()
	wg.Wait()
}
