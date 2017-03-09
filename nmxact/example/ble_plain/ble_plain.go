/**
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements.  See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership.  The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License.  You may obtain a copy of the License at
 *
 *  http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */

package main

import (
	"fmt"
	"os"

	"mynewt.apache.org/newt/nmxact/bledefs"
	"mynewt.apache.org/newt/nmxact/nmble"
	"mynewt.apache.org/newt/nmxact/sesn"
	"mynewt.apache.org/newt/nmxact/xact"
)

func main() {
	// Initialize the BLE transport.
	params := nmble.XportCfg{
		SockPath:     "/tmp/blehostd-uds",
		BlehostdPath: "blehostd",
		DevPath:      "/dev/cu.usbserial-A600ANJ1",
	}
	x, err := nmble.NewBleXport(params)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error creating BLE transport: %s\n",
			err.Error())
		os.Exit(1)
	}

	// Start the BLE transport.
	if err := x.Start(); err != nil {
		fmt.Fprintf(os.Stderr, "error starting BLE transport: %s\n",
			err.Error())
		os.Exit(1)
	}
	defer x.Stop()

	// Prepare a BLE session:
	//     * We use a random address.
	//     * Peer has public address 0b:0a:0b:0a:0b:0a.
	sc := sesn.SesnCfg{
		MgmtProto: sesn.MGMT_PROTO_NMP,

		Ble: sesn.SesnCfgBle{
			OwnAddrType: bledefs.BLE_ADDR_TYPE_RANDOM,
			Peer: bledefs.BleDev{
				AddrType: bledefs.BLE_ADDR_TYPE_PUBLIC,
				Addr: bledefs.BleAddr{
					Bytes: [6]byte{0x0b, 0x0a, 0x0b, 0x0a, 0x0b, 0x0a},
				},
			},
		},
	}

	s, err := x.BuildSesn(sc)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error creating BLE session: %s\n", err.Error())
		os.Exit(1)
	}

	// Connect to the peer (open the session).
	if err := s.Open(); err != nil {
		fmt.Fprintf(os.Stderr, "error starting BLE session: %s\n", err.Error())
		os.Exit(1)
	}
	defer s.Close()

	// Send an echo command to the peer.
	c := xact.NewEchoCmd()
	c.Payload = "hello"

	res, err := c.Run(s)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error executing echo command: %s\n",
			err.Error())
		os.Exit(1)
	}

	if res.Status() != 0 {
		fmt.Printf("Peer responded negatively to echo command; status=%d\n",
			res.Status())
	}

	eres := res.(*xact.EchoResult)
	fmt.Printf("Peer echoed back: %s\n", eres.Rsp.Payload)
}
