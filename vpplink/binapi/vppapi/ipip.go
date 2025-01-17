// Copyright (C) 2019 Cisco Systems Inc.
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
	types "git.fd.io/govpp.git/api/v0"
	"github.com/pkg/errors"
	"github.com/projectcalico/vpp-dataplane/vpplink/binapi/vppapi/interface_types"
	"github.com/projectcalico/vpp-dataplane/vpplink/binapi/vppapi/ip_types"
	"github.com/projectcalico/vpp-dataplane/vpplink/binapi/vppapi/ipip"
)

func (v *Vpp) ListIPIPTunnels() ([]*types.IPIPTunnel, error) {
	v.Lock()
	defer v.Unlock()

	tunnels := make([]*types.IPIPTunnel, 0)
	request := &ipip.IpipTunnelDump{
		SwIfIndex: interface_types.InterfaceIndex(InvalidInterface),
	}
	stream := v.GetChannel().SendMultiRequest(request)
	for {
		response := &ipip.IpipTunnelDetails{}
		stop, err := stream.ReceiveReply(response)
		if err != nil {
			return nil, errors.Wrapf(err, "error listing IPIP tunnels")
		}
		if stop {
			break
		}
		tunnels = append(tunnels, &types.IPIPTunnel{
			Src:       response.Tunnel.Src.ToIP(),
			Dst:       response.Tunnel.Dst.ToIP(),
			TableID:   response.Tunnel.TableID,
			SwIfIndex: uint32(response.Tunnel.SwIfIndex),
		})
	}
	return tunnels, nil
}

func (v *Vpp) AddIPIPTunnel(tunnel *types.IPIPTunnel) (uint32, error) {
	response := &ipip.IpipAddTunnelReply{}
	request := &ipip.IpipAddTunnel{
		Tunnel: ipip.IpipTunnel{
			Instance: ^uint32(0),
			Src:      ip_types.AddressFromIP(tunnel.Src),
			Dst:      ip_types.AddressFromIP(tunnel.Dst),
			TableID:  tunnel.TableID,
		},
	}
	err := v.SendRequestAwaitReply(request, response)
	if err != nil {
		return InvalidSwIfIndex, err
	}
	tunnel.SwIfIndex = uint32(response.SwIfIndex)
	return uint32(response.SwIfIndex), nil
}

func (v *Vpp) DelIPIPTunnel(tunnel *types.IPIPTunnel) (err error) {
	response := &ipip.IpipDelTunnelReply{}
	request := &ipip.IpipDelTunnel{
		SwIfIndex: interface_types.InterfaceIndex(tunnel.SwIfIndex),
	}
	err = v.SendRequestAwaitReply(request, response)
	if err != nil {
		return err
	}
	return nil
}
