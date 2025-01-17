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

	"github.com/projectcalico/vpp-dataplane/vpplink/binapi/vppapi/ipsec_types"
    "github.com/projectcalico/vpp-dataplane/vpplink/binapi/vppapi/tunnel_types"
)

func fromVppSaFlags(vppFlags ipsec_types.IpsecSadFlags) types.SaFlags {
	return types.SaFlags(vppFlags)
}

func toVppSaFlags(flags types.SaFlags) ipsec_types.IpsecSadFlags {
	return ipsec_types.IpsecSadFlags(flags)
}

func getVPPKey(in []byte) ipsec_types.Key {
	return ipsec_types.Key{
		Length: uint8(len(in)),
		Data:   in,
	}
}

func toVppTunnel(tunnel types.Tunnel) tunnel_types.Tunnel {
	return tunnel_types.Tunnel{
		Src:     toVppAddress(tunnel.Src),
		Dst:     toVppAddress(tunnel.Dst),
		TableID: tunnel.TableID,
	}
}

func fromVppTunnel(tunnel tunnel_types.Tunnel) types.Tunnel {
	return types.Tunnel{
		Src:     FromVppAddress(tunnel.Src),
		Dst:     FromVppAddress(tunnel.Dst),
		TableID: tunnel.TableID,
	}
}

func GetSaFlagNone() types.SaFlags {
	return types.SaFlags(ipsec_types.IPSEC_API_SAD_FLAG_NONE)
}
func GetSaFlagUseEsn() types.SaFlags {
	return types.SaFlags(ipsec_types.IPSEC_API_SAD_FLAG_USE_ESN)
}
func GetSaFlagAntiReplay() types.SaFlags {
	return types.SaFlags(ipsec_types.IPSEC_API_SAD_FLAG_USE_ANTI_REPLAY)
}
func GetSaFlagIsTunnel() types.SaFlags {
	return types.SaFlags(ipsec_types.IPSEC_API_SAD_FLAG_IS_TUNNEL)
}
func GetSaFlagIsTunnelV6() types.SaFlags {
	return types.SaFlags(ipsec_types.IPSEC_API_SAD_FLAG_IS_TUNNEL_V6)
}
func GetSaFlagUdpEncap() types.SaFlags {
	return types.SaFlags(ipsec_types.IPSEC_API_SAD_FLAG_UDP_ENCAP)
}
func GetSaFlagIsInbound() types.SaFlags {
	return types.SaFlags(ipsec_types.IPSEC_API_SAD_FLAG_IS_INBOUND)
}
func GetSaFlagAsync() types.SaFlags {
	return types.SaFlags(ipsec_types.IPSEC_API_SAD_FLAG_ASYNC)
}