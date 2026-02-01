//go:build linux || android

package tun

// The linux platform specific tun parts

import (
	"fmt"

	"github.com/vishvananda/netlink"
	wgtun "golang.zx2c4.com/wireguard/tun"
)

// Configures the TUN adapter with the correct IPv6 address and MTU.
func (tun *TunAdapter) setup(ifname string, addr string, mtu uint64,reuseexist bool) error {

	if ifname == "auto" {
		ifname = "\000"
	}
	tun.log.Warnf("Warning: pre setup params  %+v",mtu)
	// wgtun.CreateTUN()

	iface, err := wgtun.CreateTUN(ifname, int(mtu))
	if err != nil {
		return fmt.Errorf("failed to create TUN: %w", err)
	}
	tun.iface = iface
	if mtu, err := iface.MTU(); err == nil {
		tun.mtu = getSupportedMTU(uint64(mtu))
	} else {
		tun.mtu = 0
	}
	if addr != "" {
		tun.log.Infof("setup addt %s reuse %s",addr,reuseexist)
		return tun.setupAddress(addr,reuseexist)
	}
	return nil
}

// Configures the "utun" adapter from an existing file descriptor.
func (tun *TunAdapter) setupFD(fd int32, addr string, mtu uint64) error {
	return fmt.Errorf("setup via FD not supported on this platform")
}

// Configures the TUN adapter with the correct IPv6 address and MTU. Netlink
// is used to do this, so there is not a hard requirement on "ip" or "ifconfig"
// to exist on the system, but this will fail if Netlink is not present in the
// kernel (it nearly always is).
func (tun *TunAdapter) setupAddress(addr string,reuseexist bool) error {
	nladdr, err := netlink.ParseAddr(addr)
	if err != nil {
		return fmt.Errorf("couldn't parse address %q: %w", addr, err)
	}
	nlintf, err := netlink.LinkByName(tun.Name())
	if err != nil {
		return fmt.Errorf("failed to find link by name: %w", err)
	}
	if reuseexist {

tun.log.Infof("start reusage logic %t mtu check",reuseexist)
		currentmtu:=nlintf.Attrs().MTU
		if tun.mtu!= uint64(currentmtu) { fmt.Println(" mtu ok ")}
		

tun.log.Infof("start reusage logic %t addr check",reuseexist)
addrrs,getaddrs:= netlink.AddrList(nlintf,netlink.FAMILY_V6)
if getaddrs != nil {
	return fmt.Errorf("failed get addrs %s",getaddrs )
}

	for _,x := range addrrs {

tun.log.Infof("start addr logic current addr %+v req %+v",x,nladdr)
if x.Equal(*nladdr) {
	break
}
tun.log.Infof("stop addr logic current addr %+v req %+v",x,nladdr)

	} 
return nil


	} // end reuse exist
	if err := netlink.AddrAdd(nlintf, nladdr); err != nil {
		return fmt.Errorf("failed to add address to link: %w", err)
	}
	if err := netlink.LinkSetMTU(nlintf, int(tun.mtu)); err != nil {
		return fmt.Errorf("failed to set link MTU: %w", err)
	}
	if err := netlink.LinkSetUp(nlintf); err != nil {
		return fmt.Errorf("failed to bring link up: %w", err)
	}
	// Friendly output
	tun.log.Infof("Interface name: %s", tun.Name())
	tun.log.Infof("Interface IPv6: %s", addr)
	tun.log.Infof("Interface MTU: %d", tun.mtu)
	return nil
}
