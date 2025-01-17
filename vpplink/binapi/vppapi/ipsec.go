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
    "github.com/pkg/errors"

	types "git.fd.io/govpp.git/api/v0"

    "github.com/projectcalico/vpp-dataplane/vpplink/binapi/vppapi/interface_types"
    "github.com/projectcalico/vpp-dataplane/vpplink/binapi/vppapi/ipsec"
	"github.com/projectcalico/vpp-dataplane/vpplink/binapi/vppapi/ipsec_types"
)

func (v *Vpp) GetIPsecTunnelProtection(tunnelInterface uint32) (protections []types.IPsecTunnelProtection, err error) {
	v.lock.Lock()
	defer v.lock.Unlock()

	request := &ipsec.IpsecTunnelProtectDump{
		SwIfIndex: interface_types.InterfaceIndex(tunnelInterface),
	}
	stream := v.GetChannel().SendMultiRequest(request)
	for {
                response := &ipsec.IpsecTunnelProtectDetails{}
		stop, err := stream.ReceiveReply(response)
		if err != nil {
			return nil, errors.Wrapf(err, "error listing tunnel interface %u protections", tunnelInterface)
		}
		if stop {
			return protections, nil
		}
		protections = append(protections, types.IPsecTunnelProtection{
			SwIfIndex:   uint32(response.Tun.SwIfIndex),
			NextHop:     FromVppAddress(response.Tun.Nh),
			OutSAIndex:  response.Tun.SaOut,
			InSAIndices: response.Tun.SaIn,
		})
	}
}

func (v *Vpp) addDelIpsecSA(sa *types.IPSecSA, isAdd bool) error {
	v.lock.Lock()
	defer v.lock.Unlock()
	response := &ipsec.IpsecSadEntryAddDelV3Reply{}
	request := &ipsec.IpsecSadEntryAddDelV3{
		IsAdd: isAdd,
		Entry: ipsec_types.IpsecSadEntryV3{
			SadID:              sa.SAId,
			Spi:                sa.Spi,
			Protocol:           ipsec_types.IPSEC_API_PROTO_ESP,
			CryptoAlgorithm:    ipsec_types.IPSEC_API_CRYPTO_ALG_AES_CTR_128,
			CryptoKey:          getVPPKey(sa.CryptoKey),
			Salt:               sa.Salt,
			IntegrityKey:       getVPPKey(sa.IntegrityKey),
			IntegrityAlgorithm: ipsec_types.IPSEC_API_INTEG_ALG_SHA1_96,
			Flags:              toVppSaFlags(sa.Flags),
			UDPSrcPort:         uint16(sa.SrcPort),
			UDPDstPort:         uint16(sa.DstPort),
		},
	}
	if sa.Tunnel != nil {
		request.Entry.Tunnel = toVppTunnel(*sa.Tunnel)
	}
	return v.SendRequestAwaitReply(request, response)
}

func (v *Vpp) AddIpsecSA(sa *types.IPSecSA) error {
	return v.addDelIpsecSA(sa, true /* isAdd */)
}

func (v *Vpp) DelIpsecSA(sa *types.IPSecSA) error {
	return v.addDelIpsecSA(sa, false /* isAdd */)
}

func (v *Vpp) AddIpsecSAProtect(swIfIndex, saIn, saOut uint32) error {
	v.lock.Lock()
	defer v.lock.Unlock()

	response := &ipsec.IpsecTunnelProtectUpdateReply{}
	request := &ipsec.IpsecTunnelProtectUpdate{
		Tunnel: ipsec.IpsecTunnelProtect{
			SwIfIndex: interface_types.InterfaceIndex(swIfIndex),
			SaOut:     saOut,
			SaIn:      []uint32{saIn},
		},
	}
	return v.SendRequestAwaitReply(request, response)
}

func (v *Vpp) DelIpsecSAProtect(swIfIndex uint32) error {
	v.lock.Lock()
	defer v.lock.Unlock()

	response := &ipsec.IpsecTunnelProtectDelReply{}
	request := &ipsec.IpsecTunnelProtectDel{
		SwIfIndex: interface_types.InterfaceIndex(swIfIndex),
	}
	return v.SendRequestAwaitReply(request, response)
}

func (v *Vpp) AddIpsecInterface() (uint32, error) {
	v.lock.Lock()
	defer v.lock.Unlock()

	response := &ipsec.IpsecItfCreateReply{}
	request := &ipsec.IpsecItfCreate{
		Itf: ipsec.IpsecItf{
			UserInstance: ^uint32(0),
		},
	}
        err := v.SendRequestAwaitReply(request, response)
	if err != nil {
		return InvalidSwIfIndex, err
	}
	return uint32(response.SwIfIndex), nil
}

func (v *Vpp) DelIpsecInterface(swIfIndex uint32) error {
	v.lock.Lock()
	defer v.lock.Unlock()

	response := &ipsec.IpsecItfDeleteReply{}
	request := &ipsec.IpsecItfDelete{
		SwIfIndex: interface_types.InterfaceIndex(swIfIndex),
	}
	return v.SendRequestAwaitReply(request, response)
}
