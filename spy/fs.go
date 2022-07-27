package spy

import (
	"encoding/binary"
	"encoding/json"
	"github.com/go-well/spider/silk"
	"io"
	"os"
	"strings"
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

func init() {
	RegisterHandler(silk.FsList, func(c *Client, p *silk.Package) {
		p.Type = silk.FsListAck
		path := string(p.Data)
		dirs, err := os.ReadDir(path)
		if err != nil {
			p.SetError(err.Error())
			_ = c.Send(p)
			return
		}
		ds := make([]dir, 0)
		for _, d := range dirs {
			ds = append(ds, dir{Name: d.Name(), IsDir: d.IsDir()})
		}
		p.Data, err = json.Marshal(ds)
		_ = c.Send(p)
	})

	RegisterHandler(silk.FsMkDir, func(c *Client, p *silk.Package) {
		p.Type = silk.FsMkDirAck
		path := string(p.Data)
		err := os.Mkdir(path, os.ModePerm)
		if err != nil {
			p.SetError(err.Error())
		}
		_ = c.Send(p)
	})

	RegisterHandler(silk.FsRemove, func(c *Client, p *silk.Package) {
		p.Type = silk.FsRemoveAck
		path := string(p.Data)
		err := os.RemoveAll(path)
		if err != nil {
			p.SetError(err.Error())
		}
		_ = c.Send(p)
	})

	RegisterHandler(silk.FsRename, func(c *Client, p *silk.Package) {
		p.Type = silk.FsRenameAck
		path := string(p.Data)
		str := strings.Split(path, ",")
		if len(str) < 2 {
			p.SetError("old,new")
			_ = c.Send(p)
			return
		}
		err := os.Rename(str[0], str[1])
		if err != nil {
			p.SetError(err.Error())
		}
		_ = c.Send(p)
	})

	RegisterHandler(silk.FsStats, func(c *Client, p *silk.Package) {
		p.Type = silk.FsStatsAck
		path := string(p.Data)
		st, err := os.Stat(path)
		if err != nil {
			p.SetError(err.Error())
			_ = c.Send(p)
			return
		}
		s := &stat{Name: st.Name(), IsDir: st.IsDir(), Size: st.Size(), Time: st.ModTime()}
		p.Data, err = json.Marshal(s)
		_ = c.Send(p)
	})

	RegisterHandler(silk.FsDownload, func(c *Client, p *silk.Package) {
		p.Type = silk.FsDownloadAck
		path := string(p.Data)
		file, err := os.Open(path)
		if err != nil {
			p.SetError(err.Error())
			_ = c.Send(p)
			return
		}

		//缓存
		p.Data = make([]byte, 2)
		id := c.newFile(file)
		binary.BigEndian.PutUint16(p.Data, id)
		_ = c.Send(p)

		//主动发送第一个数据包
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
	})

	RegisterHandler(silk.FsUpload, func(c *Client, p *silk.Package) {
		p.Type = silk.FsUploadAck
		path := string(p.Data)
		file, err := os.OpenFile(path, os.O_CREATE, os.ModePerm)
		if err != nil {
			p.SetError(err.Error())
			_ = c.Send(p)
			return
		}

		//缓存
		id := c.newFile(file)
		p.Data = make([]byte, 2)
		binary.BigEndian.PutUint16(p.Data, id)

		_ = c.Send(p)
	})

	//处理上传
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

	//处理下载
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
