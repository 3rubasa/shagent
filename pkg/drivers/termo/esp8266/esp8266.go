package esp8266

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"sync"
	"time"
)

const connType = "tcp4"
const bufSize = 1024
const smallBufSize = 256

type Sensor struct {
	port     int
	t        float64
	mux      *sync.Mutex
	tempTS   time.Time
	cacheTTL time.Duration
	l        net.Listener
}

func New(port int, cacheTTL time.Duration) *Sensor {
	return &Sensor{
		port:     port,
		t:        -275.0,
		mux:      &sync.Mutex{},
		tempTS:   time.Date(1800, 1, 1, 0, 0, 0, 0, time.Local),
		cacheTTL: cacheTTL,
	}
}

// Handles incoming requests.
func (s *Sensor) handleRequest(conn net.Conn) {
	defer conn.Close()

	var err error

	now := time.Now()

	buf := make([]byte, 0, bufSize)
	b := make([]byte, smallBufSize)

	// Read the incoming connection into the buffer.
	bytesRead := 0

	for {
		n, err := conn.Read(b)
		if err != nil && err != io.EOF && !os.IsTimeout(err) {
			log.Println("ERROR: ESP8266 termo sensor: failed to read from socket: ", err)
			return
		}

		bytesRead += n
		if bytesRead > bufSize {
			log.Println("ERROR: ESP8266 termo sensor: read from socket more data than can fit into the buffer: ", bytesRead)
			return
		}

		buf = append(buf, b[0:n]...)

		if err == io.EOF || os.IsTimeout(err) {
			break
		}
	}

	var data SensorData

	err = json.Unmarshal(buf, &data)
	if err != nil {
		log.Println("ERROR: ESP8266 termo sensor: Failed to unmarshal sensor data: ", err, "; the data that came from sensor: ", buf)
		return
	}

	if data.Temp > 900 { // invalid value is 998, but because it's float type, we don't want to do direct comparison
		log.Println("ERROR: ESP8266 termo sensor: temperature value is invalid: ", data.Temp)
		return
	}

	s.mux.Lock()
	s.t = data.Temp
	s.tempTS = now
	s.mux.Unlock()
}

func (s *Sensor) Start() error {
	var err error
	// Listen for incoming connections.
	s.l, err = net.Listen(connType, fmt.Sprintf(":%d", s.port))
	if err != nil {
		log.Println("ERROR: ESP8266 termo sensor: Failed to start listening for incoming connections", err)
		return err
	}

	log.Println("Dbg: ESP8266 termo driver: Listening on ", s.l.Addr())
	go func() {
		for {
			// Listen for an incoming connection.
			conn, err := s.l.Accept()
			if err != nil {
				log.Println("NOTICE: ESP8266 termo driver: Accept() exited with an error ", err)
				break
			}

			// Handle connections in a new goroutine.
			go s.handleRequest(conn)
		}
	}()

	return nil
}

func (s *Sensor) Stop() error {
	if s.l != nil {
		s.l.Close()
	}

	return nil
}

func (s *Sensor) Get() (float64, error) {
	s.mux.Lock()
	t := s.t
	tempTS := s.tempTS
	s.mux.Unlock()

	if time.Since(tempTS) > s.cacheTTL {
		return -100, ErrNoFreshData
	}

	return t, nil
}
