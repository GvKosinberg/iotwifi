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
func (c *Command) RemoveApInterface() {
	cmd := exec.Command("iw", "dev", "wlan1", "del")
	cmd.Start()
	cmd.Wait()
}

// ConfigureApInterface configured the AP interface.
func (c *Command) ConfigureApInterface() {
	cmd := exec.Command("ifconfig", "wlan1", c.SetupCfg.HostApdCfg.Ip)
	cmd.Start()
	cmd.Wait()
}

// UpApInterface ups the AP Interface.
func (c *Command) UpApInterface() {
	cmd := exec.Command("ifconfig", "wlan1", "up")
	cmd.Start()
	cmd.Wait()
}

// AddApInterface adds the AP interface.
func (c *Command) AddApInterface() {
	cmd := exec.Command("iw", "phy", "phy1", "interface", "add", "wlan1", "type", "__ap")
	cmd.Start()
	cmd.Wait()
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
		"--no-hosts", // Don't read the hostnames in /etc/hosts.
		"--interface=wlan1",
		"--keep-in-foreground",
		"--log-queries",
		"--no-resolv",
		"--address=" + c.SetupCfg.DnsmasqCfg.Address,
		"--dhcp-range=" + c.SetupCfg.DnsmasqCfg.DhcpRange,
		"--dhcp-vendorclass=" + c.SetupCfg.DnsmasqCfg.VendorClass,
		"--dhcp-authoritative",
		"--log-facility=-",
		"--log-dhcp",
	}

	cmd := exec.Command("dnsmasq", args...)
	go c.Runner.ProcessCmd("dnsmasq", cmd)
}

// Add bridge (br0)
// BrideAPtoEth bridges the connection from eth0 to uap0
func (c *Command) BridgeAPtoEth() {

	cmd_sed := exec.Command("sed", "-i", "s/#?net.ipv4.ip_forward.*/net.ipv4.ip_forward = 1/", "/etc/sysctl.conf")
	cmd_sed.Run()
	cmd_sysctl := exec.Command("sysctl", "-p")
	cmd_sysctl.Run()

	iptables0_args := []string{
        "-t",
        "nat",
        "-A",
        "POSTROUTING",
        "-o",
        "eth0",
        "-j",
        "MASQUERADE",
    }
	cmd_iptables0 := exec.Command("iptables", iptables0_args...)
  cmd_iptables0.Run()

  iptables1_args := []string{
	      "-A",
	      "FORWARD",
	      "-i",
	      "eth0",
				"-o",
				"wlan1",
				"-m",
				"state",
				"--state",
				"RELATED, ESTABLISHED",
	      "-j",
	      "ACCEPT",
  }
	cmd_iptables1 := exec.Command("iptables", iptables1_args...)
  cmd_iptables1.Run()

	iptables2_args := []string{
				"-A",
				"FORWARD",
				"-i",
				"wlan1",
				"-o",
				"eth0",
				"-j",
				"ACCEPT",
	}
	cmd_iptables2 := exec.Command("iptables", iptables2_args...)
	cmd_iptables2.Run()
}
