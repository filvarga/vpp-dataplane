// Copyright (C) 2020 Cisco Systems Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or
// implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package vppapi

import (
	"fmt"
	types "git.fd.io/govpp.git/api/v0"
	"github.com/projectcalico/vpp-dataplane/vpplink/binapi/vppapi/interface_types"
	"github.com/projectcalico/vpp-dataplane/vpplink/binapi/vppapi/ip_types"
	"github.com/pkg/errors"
	"github.com/projectcalico/vpp-dataplane/vpplink/binapi/vppapi/vxlan"
)

func (v *Vpp) ListVXLanTunnels() ([]types.VXLanTunnel, error) {
	v.Lock()
	defer v.Unlock()

	tunnels := make([]types.VXLanTunnel, 0)
	request := &vxlan.VxlanTunnelV2Dump{
		SwIfIndex: interface_types.InterfaceIndex(InvalidInterface),
	}
	stream := v.GetChannel().SendMultiRequest(request)
	for {
		response := &vxlan.VxlanTunnelV2Details{}
		stop, err := stream.ReceiveReply(response)
		if err != nil {
			return nil, errors.Wrapf(err, "error listing VXLan tunnels")
		}
		if stop {
			break
		}
		tunnels = append(tunnels, types.VXLanTunnel{
			SrcAddress:     response.SrcAddress.ToIP(),
			DstAddress:     response.DstAddress.ToIP(),
			SrcPort:        response.SrcPort,
			DstPort:        response.DstPort,
			Vni:            response.Vni,
			DecapNextIndex: response.DecapNextIndex,
			SwIfIndex:      uint32(response.SwIfIndex),
		})
	}
	return tunnels, nil
}
func (v *Vpp) addDelVXLanTunnel(tunnel *types.VXLanTunnel, isAdd bool) (swIfIndex uint32, err error) {
	v.Lock()
	defer v.Unlock()

	response := &vxlan.VxlanAddDelTunnelV3Reply{}
	request := &vxlan.VxlanAddDelTunnelV3{
		IsAdd:          isAdd,
		Instance:       ^uint32(0),
		SrcAddress:     ip_types.AddressFromIP(tunnel.SrcAddress),
		DstAddress:     ip_types.AddressFromIP(tunnel.DstAddress),
		SrcPort:        tunnel.SrcPort,
		DstPort:        tunnel.DstPort,
		Vni:            tunnel.Vni,
		DecapNextIndex: tunnel.DecapNextIndex,
		IsL3:           true,
	}
	err = v.GetChannel().SendRequest(request).ReceiveReply(response)
	opStr := "Del"
	if isAdd {
		opStr = "Add"
	}
	if err != nil {
	// TODO: return invalid interface here from types
		return InvalidSwIfIndex, errors.Wrapf(err, "%s vxlan Tunnel failed", opStr)
	} else if response.Retval != 0 {
		return InvalidSwIfIndex, fmt.Errorf("%s vxlan Tunnel failed with retval %d", opStr, response.Retval)
	}
	tunnel.SwIfIndex = uint32(response.SwIfIndex)
	return uint32(response.SwIfIndex), nil
}

func (v *Vpp) AddVXLanTunnel(tunnel *types.VXLanTunnel) (swIfIndex uint32, err error) {
	return v.addDelVXLanTunnel(tunnel, true)
}

func (v *Vpp) DelVXLanTunnel(tunnel *types.VXLanTunnel) (err error) {
	_, err = v.addDelVXLanTunnel(tunnel, false)
	return err
}
