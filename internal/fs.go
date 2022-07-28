package internal

import (
	"encoding/binary"
	"github.com/go-well/spider/silk"
	"io"
	"os"
	"time"
)

type dir struct {
	Name  string `json:"name"`
	IsDir bool   `json:"isDir"`
}

type stat struct {
	Name  string    `json:"name"`
	Size  int64     `json:"size"`
	IsDir bool      `json:"isDir"`
	Time  time.Time `json:"time"`
}

func (c *Client) Download(filename string, localFile string) error {
	p, err := c.Ask(&silk.Package{Type: silk.FsDownload, Data: []byte(filename)})
	if err != nil {
		return err
	}

	file, err := os.Create(localFile)
	if err != nil {
		return err
	}

	id := binary.BigEndian.Uint16(p.Data)
	c.files.Store(id, file)

	return nil
}

func (c *Client) Upload(filename string, localFile string) error {
	p, err := c.Ask(&silk.Package{Type: silk.FsUpload, Data: []byte(filename)})
	if err != nil {
		return err
	}

	file, err := os.Open(localFile)
	if err != nil {
		return err
	}

	id := binary.BigEndian.Uint16(p.Data)
	c.files.Store(id, file)

	//主动发送第一个数据包
	p.Type = silk.FsData
	p.Data = make([]byte, 512)
	binary.BigEndian.PutUint16(p.Data, id)
	n, e := file.Read(p.Data[2:])
	if e != nil {
		if e == io.EOF {
			p.Type = silk.FsDataEnd
		} else {
			p.SetError(e.Error())
		}
	} else {
		p.Data = p.Data[:2+n]
	}
	_ = c.Send(p)

	return nil
}

func init() {

	//处理下载
	RegisterHandler(silk.FsData, func(c *Client, p *silk.Package) {
		p.Type = silk.FsDataAck
		id := binary.BigEndian.Uint16(p.Data)
		f, ok := c.files.Load(id)
		if !ok {
			p.SetError("file not exists")
			_ = c.Send(p)
			return
		}

		file := f.(*os.File)
		_, err := file.Write(p.Data[2:])
		if err != nil {
			p.SetError(err.Error())
		} else {
			p.Data = p.Data[:2]
		}
		_ = c.Send(p)
	})

	//处理上传
	RegisterHandler(silk.FsDataAck, func(c *Client, p *silk.Package) {
		p.Type = silk.FsData
		id := binary.BigEndian.Uint16(p.Data)
		f, ok := c.files.Load(id)
		if !ok {
			p.SetError("file not exists")
			_ = c.Send(p)
			return
		}

		file := f.(*os.File)
		p.Data = make([]byte, 512)
		binary.BigEndian.PutUint16(p.Data, id)
		n, err := file.Read(p.Data[2:])
		if err != nil {
			if err == io.EOF {
				p.Type = silk.FsDataEnd
			} else {
				p.SetError(err.Error())
			}
		} else {
			p.Data = p.Data[:2+n]
		}
		_ = c.Send(p)
	})

	//处理结束
	RegisterHandler(silk.FsDataEnd, func(c *Client, p *silk.Package) {
		p.Type = silk.FsDataEndAck
		id := binary.BigEndian.Uint16(p.Data)
		f, ok := c.files.LoadAndDelete(id)
		if !ok {
			p.SetError("file not exists")
			_ = c.Send(p)
			return
		}

		file := f.(*os.File)
		//file.Sync()
		err := file.Close()
		if err != nil {
			p.SetError(err.Error())
		}
		_ = c.Send(p)
	})

}
