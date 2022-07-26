package spy

import (
	"encoding/json"
	"github.com/go-well/spider/silk"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/host"
	"github.com/shirou/gopsutil/mem"
	"github.com/shirou/gopsutil/net"
)

func init() {
	RegisterHandler(silk.StatsHost, func(c *Client, p *silk.Package) {
		p.Type = silk.StatsHostAck
		info, err := host.Info()
		if err != nil {
			p.SetError(err.Error())
		} else {
			p.Data, _ = json.Marshal(info)
		}
		_ = c.Send(p)
	})

	RegisterHandler(silk.StatsCpu, func(c *Client, p *silk.Package) {
		p.Type = silk.StatsCpuAck
		info, err := cpu.Info()
		if err != nil {
			p.SetError(err.Error())
		} else {
			p.Data, _ = json.Marshal(info)
		}
		_ = c.Send(p)
	})

	RegisterHandler(silk.StatsCpuTimes, func(c *Client, p *silk.Package) {
		p.Type = silk.StatsCpuTimesAck
		info, err := cpu.Times(true)
		if err != nil {
			p.SetError(err.Error())
		} else {
			p.Data, _ = json.Marshal(info)
		}
		_ = c.Send(p)
	})

	RegisterHandler(silk.StatsMem, func(c *Client, p *silk.Package) {
		p.Type = silk.StatsMemAck
		info, err := mem.VirtualMemory()
		if err != nil {
			p.SetError(err.Error())
		} else {
			p.Data, _ = json.Marshal(info)
		}
		_ = c.Send(p)
	})

	RegisterHandler(silk.StatsDisk, func(c *Client, p *silk.Package) {
		p.Type = silk.StatsDiskAck
		info, err := disk.Partitions(true)
		if err != nil {
			p.SetError(err.Error())
		} else {
			p.Data, _ = json.Marshal(info)
		}
		_ = c.Send(p)
	})

	RegisterHandler(silk.StatsDiskUsage, func(c *Client, p *silk.Package) {
		p.Type = silk.StatsDiskUsageAck
		path := string(p.Data)
		info, err := disk.Usage(path)
		if err != nil {
			p.SetError(err.Error())
		} else {
			p.Data, _ = json.Marshal(info)
		}
		_ = c.Send(p)
	})

	RegisterHandler(silk.StatsNet, func(c *Client, p *silk.Package) {
		p.Type = silk.StatsNetAck
		info, err := net.Interfaces()
		if err != nil {
			p.SetError(err.Error())
		} else {
			p.Data, _ = json.Marshal(info)
		}
		_ = c.Send(p)
	})

	RegisterHandler(silk.StatsConnection, func(c *Client, p *silk.Package) {
		p.Type = silk.StatsConnectionAck
		kind := string(p.Data)
		info, err := net.Connections(kind)
		if err != nil {
			p.SetError(err.Error())
		} else {
			p.Data, _ = json.Marshal(info)
		}
		_ = c.Send(p)
	})

	RegisterHandler(silk.StatsUser, func(c *Client, p *silk.Package) {
		p.Type = silk.StatsUserAck
		info, err := host.Users()
		if err != nil {
			p.SetError(err.Error())
		} else {
			p.Data, _ = json.Marshal(info)
		}
		_ = c.Send(p)
	})
}
