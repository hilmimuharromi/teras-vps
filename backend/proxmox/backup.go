package proxmox

import (
	"encoding/json"
	"fmt"
)

// Snapshot represents a VM snapshot
type Snapshot struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Snaptime    int64  `json:"snaptime"`
}

// CreateSnapshot creates a VM snapshot
func (c *Client) CreateSnapshot(vmid int, name string, description string) error {
	path := fmt.Sprintf("/nodes/%s/qemu/%d/snapshot", c.Node, vmid)

	config := map[string]interface{}{
		"snapname":    name,
		"description": description,
	}

	_, err := c.Request("POST", path, config)
	return err
}

// ListSnapshots lists all snapshots for a VM
func (c *Client) ListSnapshots(vmid int) ([]Snapshot, error) {
	path := fmt.Sprintf("/nodes/%s/qemu/%d/snapshot", c.Node, vmid)

	data, err := c.Request("GET", path, nil)
	if err != nil {
		return nil, err
	}

	var snapshots []Snapshot
	if err := json.Unmarshal(data, &snapshots); err != nil {
		return nil, err
	}

	return snapshots, nil
}

// DeleteSnapshot deletes a VM snapshot
func (c *Client) DeleteSnapshot(vmid int, name string) error {
	path := fmt.Sprintf("/nodes/%s/qemu/%d/snapshot/%s", c.Node, vmid, name)

	_, err := c.Request("DELETE", path, nil)
	return err
}

// RollbackSnapshot rolls back a VM to a snapshot
func (c *Client) RollbackSnapshot(vmid int, name string) error {
	path := fmt.Sprintf("/nodes/%s/qemu/%d/snapshot/%s/rollback", c.Node, vmid, name)

	_, err := c.Request("POST", path, nil)
	return err
}
