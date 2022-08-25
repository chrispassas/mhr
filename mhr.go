package mhr

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"strconv"
	"time"
)

type Result struct {
	Hash      string
	Timestamp time.Time
	HitRate   int
	NoData    bool
}

var (
	ErrMaxHashes = fmt.Errorf("Exceeded max hashes per-request of 1,000")
)

func Search(ctx context.Context, hashes []string) (results []Result, err error) {

	if len(hashes) > 1000 {
		err = ErrMaxHashes
		return results, err
	}

	var (
		d        = &net.Dialer{}
		conn     net.Conn
		message  []byte
		response []byte
	)

	if conn, err = d.DialContext(ctx, "tcp", "hash.cymru.com:43"); err != nil {
		err = fmt.Errorf("d.DialContext() error:%w", err)
		return results, err
	}
	defer conn.Close()

	message = createMessage(hashes)
	var n int
	if n, err = conn.Write(message); err != nil {
		err = fmt.Errorf("conn.Write() n:%d error:%w", n, err)
		return results, err
	}

	if response, err = io.ReadAll(conn); err != nil {
		err = fmt.Errorf("io.ReadAll() error:%w", err)
		return results, err
	}

	log.Printf("%s", string(response))

	results, err = parseResponse(response)

	return results, err
}

func parseResponse(response []byte) (results []Result, err error) {
	var pound = []byte("#")
	var space = []byte(" ")
	for _, line := range bytes.Split(response, []byte("\n")) {
		if bytes.HasPrefix(line, pound) {
			continue
		}

		pieces := bytes.Split(line, space)
		if len(pieces) != 3 {
			continue
		}
		/*
			START HERE
			finish parsing results line by line

		*/
		var result Result

		result.Hash = string(pieces[0])
		unix, _ := strconv.Atoi(string(pieces[1]))

		result.Timestamp = time.Unix(int64(unix), 0)

		if string(pieces[2]) == "NO_DATA" {
			result.NoData = true
		} else {
			result.HitRate, _ = strconv.Atoi(string(pieces[2]))
		}

		results = append(results, result)
	}

	return results, err
}

func createMessage(hashes []string) (message []byte) {
	var msg bytes.Buffer
	msg.WriteString("begin\n")

	for _, hash := range hashes {
		msg.WriteString(hash + "\n")
	}

	msg.WriteString("end\n")
	message = msg.Bytes()
	return message
}
