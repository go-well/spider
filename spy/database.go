package spy

import (
	"encoding/json"
	"github.com/go-well/spider/silk"
	"os"
	"xorm.io/xorm"
)

func RegisterXORM(engine *xorm.Engine) {
	RegisterHandler(silk.DatabaseQuery, func(c *Client, p *silk.Package) {
		p.Type = silk.DatabaseQueryAck
		sql := string(p.Data)
		res, err := engine.Query(sql)
		if err != nil {
			p.SetError(err.Error())
		} else {
			p.Data, _ = json.Marshal(res)
		}
		_ = c.Send(p)
	})

	RegisterHandler(silk.DatabaseExec, func(c *Client, p *silk.Package) {
		p.Type = silk.DatabaseExecAck
		sql := string(p.Data)
		res, err := engine.Exec(sql)
		if err != nil {
			p.SetError(err.Error())
		} else {
			p.Data, _ = json.Marshal(res)
		}
		_ = c.Send(p)
	})

	RegisterHandler(silk.DatabaseMeta, func(c *Client, p *silk.Package) {
		p.Type = silk.DatabaseMetaAck
		res, err := engine.DBMetas()
		if err != nil {
			p.SetError(err.Error())
		} else {
			p.Data, _ = json.Marshal(res)
		}
		_ = c.Send(p)
	})

	RegisterHandler(silk.DatabaseDriver, func(c *Client, p *silk.Package) {
		p.Type = silk.DatabaseDriverAck
		name := engine.DriverName()
		p.Data = []byte(name)
		_ = c.Send(p)
	})

	RegisterHandler(silk.DatabaseSource, func(c *Client, p *silk.Package) {
		p.Type = silk.DatabaseSourceAck
		name := engine.DataSourceName()
		p.Data = []byte(name)
		_ = c.Send(p)
	})

	RegisterHandler(silk.DatabaseDump, func(c *Client, p *silk.Package) {
		p.Type = silk.DatabaseDumpAck
		file, err := os.CreateTemp("", "database.*.sql")
		if err != nil {
			p.SetError(err.Error())
			_ = c.Send(p)
			return
		}
		defer file.Close()
		err = engine.DumpAll(file)
		if err != nil {
			p.SetError(err.Error())
		} else {
			p.Data = []byte(file.Name())
		}
		_ = c.Send(p)
	})

	RegisterHandler(silk.DatabaseImport, func(c *Client, p *silk.Package) {
		p.Type = silk.DatabaseImportAck
		filename := string(p.Data)
		res, err := engine.ImportFile(filename)
		if err != nil {
			p.SetError(err.Error())
		} else {
			p.Data, _ = json.Marshal(res)
		}
		_ = c.Send(p)
	})

}
