package iotwifi

import (
	"os/exec"

	"github.com/bhoriuchi/go-bunyan/bunyan"
)

// Command for device network commands.
type Command struct {
	Log      bunyan.Logger
	Runner   CmdRunner
	SetupCfg *SetupCfg
}

// RemoveApInterface removes the AP interface.
func (c *Command) SetClient() {
	cmd := exec.Command("iw", "dev", "wlan0", "del")
	cmd.Run()
	cmd := exec.Command("iw", "phy", "phy0", "interface", "add", "wlan0", "type", "managed")
	cmd.Run()
}

// RemoveApInterface removes the AP interface.
func (c *Command) RemoveApInterface() {
	cmd := exec.Command("iw", "dev", "wlan1", "del")
	cmd.Run()
}

// ConfigureApInterface configured the AP interface.
func (c *Command) ConfigureApInterface() {
	cmd := exec.Command("ifconfig", "wlan1", c.SetupCfg.HostApdCfg.Ip)
	cmd.Run()
}

// UpApInterface ups the AP Interface.
func (c *Command) UpApInterface() {
	cmd := exec.Command("ifconfig", "wlan1", "up")
	cmd.Run()
}

// AddApInterface adds the AP interface.
func (c *Command) AddApInterface() {
	cmd := exec.Command("iw", "phy", "phy1", "interface", "add", "wlan1", "type", "__ap")
	cmd.Run()
}

// CheckInterface checks the AP interface.
func (c *Command) CheckApInterface() {
	cmd := exec.Command("ifconfig", "wlan1")
	go c.Runner.ProcessCmd("ifconfig_wlan1", cmd)
}

// StartWpaSupplicant starts wpa_supplicant.
func (c *Command) StartWpaSupplicant() {

	args := []string{
		"-Dnl80211",
		"-iwlan0",
		"-c/etc/wpa_supplicant/wpa_supplicant.conf",
	}

	cmd := exec.Command("wpa_supplicant", args...)
	go c.Runner.ProcessCmd("wpa_supplicant", cmd)
}

// StartDnsmasq starts dnsmasq.
func (c *Command) StartDnsmasq() {
	// hostapd is enabled, fire up dnsmasq
	args := []string{
		"--interface=wlan1",
		"--listen-address=192.168.27.1",
		"--dhcp-range=" + c.SetupCfg.DnsmasqCfg.DhcpRange,
		"--bind-interfaces",
		"--bogus-priv",
		"--log-dhcp",
		"--keep-in-foreground",
		"--no-hosts", // Don't read the hostnames in /etc/hosts.
		"--log-queries",
		"--no-resolv",
		"--address=" + c.SetupCfg.DnsmasqCfg.Address,
		"--dhcp-vendorclass=" + c.SetupCfg.DnsmasqCfg.VendorClass,
		"--dhcp-authoritative",
		"--log-facility=-",
		"--no-dhcp-interface=wlan0",
		"--cache-size=650",
	}

	cmd := exec.Command("dnsmasq", args...)
	go c.Runner.ProcessCmd("dnsmasq", cmd)
}
