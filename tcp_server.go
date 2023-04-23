package main

import (
	"io"
"strconv"
"fmt"
	"log"
	"net"
	"sync"
        "strings"
"github.com/Omrrrr89/is105sem03/mycrypt"
"funtemps/conv"
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
dekryptertMelding := mycrypt.Krypter([]rune(string(buf[:n])), mycrypt.ALF_SEM03, len(mycrypt.ALF_SEM03)-4)
log.Println("Dekrypter melding: ", string(dekryptertMelding))
msg := string(dekryptertMelding)  
				switch msg  {
  				        case "ping":
                                         _, err = c.Write([]byte("pong"))
                                        case" kjevik":
parts := strings.Split(msg, ";")
                                    if len(parts) < 4 {
                                  log.Println("Invalid input message")
                                       return     }
                     t, err := strconv.ParseFloat(strings.TrimSpace(parts[3]), 64)
                     if err != nil {
                     log.Println(err) }
                     f := conv.CelsiusToFahrenheit(t)

                      response := fmt.Sprintf("%.2f Celsius er %.2f Fahrenheit", t, f)
                           _, err = c.Write([]byte(response))

					default:
						_, err = c.Write(buf[:n])
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
