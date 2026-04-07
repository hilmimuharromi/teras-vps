package proxmox

import (
	"encoding/json"
	"fmt"
)

// VMConfig represents VM configuration
type VMConfig struct {
	Name        string `json:"name"`
	Cores       int    `json:"cores"`
	Memory      int    `json:"memory"`
	Disk        int    `json:"disk"`
	IPConfig0   string `json:"ipconfig0"`
	BootDisk    string `json:"bootdisk"`
	OSType      string `json:"ostype"`
	Network     string `json:"net0"`
	SCSIHW      string `json:"scsihw"`
}

// CreateVM creates a new VM in Proxmox
func (c *Client) CreateVM(vmid int, name string, cores int, memory int, diskGB int, template string) error {
	path := fmt.Sprintf("/nodes/%s/qemu", c.Node)

	config := VMConfig{
		Name:      name,
		Cores:     cores,
		Memory:    memory,
		Disk:      diskGB,
		OSType:    "l26",
		Network:   "virtio,bridge=vmbr0",
		SCSIHW:    "virtio-scsi-pci",
		BootDisk:  "scsi0",
	}

	// If template is specified, clone from it
	if template != "" {
		// Clone from template
		return c.CloneVM(vmid, name, template)
	}

	_, err := c.Request("POST", path, config)
	return err
}

// CloneVM clones an existing VM template
func (c *Client) CloneVM(vmid int, name string, templateID int) error {
	path := fmt.Sprintf("/nodes/%s/qemu/%d/clone", c.Node, templateID)

	config := map[string]interface{}{
		"newid": vmid,
		"name":  name,
		"full":  true,
	}

	_, err := c.Request("POST", path, config)
	return err
}

// DeleteVM deletes a VM from Proxmox
func (c *Client) DeleteVM(vmid int) error {
	path := fmt.Sprintf("/nodes/%s/qemu/%d", c.Node, vmid)

	_, err := c.Request("DELETE", path, nil)
	return err
}

// StartVM starts a VM
func (c *Client) StartVM(vmid int) error {
	path := fmt.Sprintf("/nodes/%s/qemu/%d/status/start", c.Node, vmid)

	_, err := c.Request("POST", path, nil)
	return err
}

// StopVM stops a VM
func (c *Client) StopVM(vmid int) error {
	path := fmt.Sprintf("/nodes/%s/qemu/%d/status/stop", c.Node, vmid)

	_, err := c.Request("POST", path, nil)
	return err
}

// RebootVM reboots a VM
func (c *Client) RebootVM(vmid int) error {
	path := fmt.Sprintf("/nodes/%s/qemu/%d/status/reboot", c.Node, vmid)

	_, err := c.Request("POST", path, nil)
	return err
}

// GetVMStatus gets the current status of a VM
func (c *Client) GetVMStatus(vmid int) (string, error) {
	path := fmt.Sprintf("/nodes/%s/qemu/%d/status/current", c.Node, vmid)

	data, err := c.Request("GET", path, nil)
	if err != nil {
		return "", err
	}

	var status struct {
		Status string `json:"status"`
	}

	if err := json.Unmarshal(data, &status); err != nil {
		return "", err
	}

	return status.Status, nil
}

// GetVMStats gets VM statistics (CPU, memory, network, disk)
func (c *Client) GetVMStats(vmid int) (map[string]interface{}, error) {
	path := fmt.Sprintf("/nodes/%s/qemu/%d/status/current", c.Node, vmid)

	data, err := c.Request("GET", path, nil)
	if err != nil {
		return nil, err
	}

	var stats map[string]interface{}
	if err := json.Unmarshal(data, &stats); err != nil {
		return nil, err
	}

	return stats, nil
}

// ListVMs lists all VMs on the node
func (c *Client) ListVMs() ([]map[string]interface{}, error) {
	path := fmt.Sprintf("/nodes/%s/qemu", c.Node)

	data, err := c.Request("GET", path, nil)
	if err != nil {
		return nil, err
	}

	var vms []map[string]interface{}
	if err := json.Unmarshal(data, &vms); err != nil {
		return nil, err
	}

	return vms, nil
}

// UpdateVMConfig updates VM configuration
func (c *Client) UpdateVMConfig(vmid int, config map[string]interface{}) error {
	path := fmt.Sprintf("/nodes/%s/qemu/%d/config", c.Node, vmid)

	_, err := c.Request("PUT", path, config)
	return err
}

// ResizeDisk resizes a VM disk
func (c *Client) ResizeDisk(vmid int, disk string, sizeGB int) error {
	path := fmt.Sprintf("/nodes/%s/qemu/%d/resize", c.Node, vmid)

	config := map[string]interface{}{
		"disk": disk,
		"size": fmt.Sprintf("%dG", sizeGB),
	}

	_, err := c.Request("PUT", path, config)
	return err
}

// GetVMConfig gets VM configuration
func (c *Client) GetVMConfig(vmid int) (map[string]interface{}, error) {
	path := fmt.Sprintf("/nodes/%s/qemu/%d/config", c.Node, vmid)

	data, err := c.Request("GET", path, nil)
	if err != nil {
		return nil, err
	}

	var config map[string]interface{}
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, err
	}

	return config, nil
}
