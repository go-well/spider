package spy

import (
	"github.com/zgwit/spider/silk"
	"net"
	"os"
	"strings"
)

func Open() {
	conn, err := net.Dial("tcp", "127.0.0.1:1206")
	if err != nil {
		return
	}

	parser := silk.NewParser(conn, func(p *silk.Package) {
		if p.Type >= 0x80 {
			//string(p.Data)
			return
		}
		switch p.Type {
		case silk.Close:
			_ = conn.Close()
		case silk.Heartbeat:
		case silk.ConnectAck:
		case silk.SubscribeAck:
		case silk.UnsubscribeAck:
		case silk.PublishAck:
		case silk.Message: //topic, message
			str := string(p.Data)
			index := strings.Index(str, ",")
			if index > -1 {
				topic := str[:index]
				message := str[index+1:]
			} else {
				topic := str
			}

		case silk.TunnelCreate: //net,addr 例 tcp,127.0.0.1:8080
			str := strings.Split(string(p.Data), ",")
			conn, err := net.Dial(str[0], str[1])

		case silk.TunnelCreateAck: //id(uint16)
		case silk.TunnelClose: //id
		case silk.TunnelCloseAck: //id
		case silk.TunnelTransferData: //id,data
		//TunnelTransferEnd

		case silk.SystemShell: // /bin/sh

		case silk.SystemShellAck: //tunnel id(uint16)
		case silk.SystemExecute: //command string
		case silk.SystemExecuteAck: //stdout
		case silk.SystemConfig:
		case silk.SystemConfigAck: //yaml
		case silk.SystemDbQuery: //sql
		case silk.SystemDbQueryAck: //json
		case silk.SystemDbExec: //sql
		case silk.SystemDbExecAck: //text

		case silk.StatsHost:
		case silk.StatsHostAck: //json
		case silk.StatsCpu:
		case silk.StatsCpuAck: //json
		case silk.StatsMem:
		case silk.StatsMemAck: //json
		case silk.StatsDisk:
		case silk.StatsDiskAck: //json
		case silk.StatsNet:
		case silk.StatsNetAck: //json
		case silk.StatsUser:
		case silk.StatsUserAck: //json

		case silk.FsList: //path
		case silk.FsListAck: //json
		case silk.FsMkDir: //path
		case silk.FsMkDirAck: //
		case silk.FsRemove: //path
		case silk.FsRemoveAck: //
		case silk.FsRename: //path,path
			str := strings.Split(string(p.Data), ",")
			_ = os.Rename(str[0], str[1])
		case silk.FsRenameAck: //
		case silk.FsStats: //path
		case silk.FsStatsOk: //json
		case silk.FsDownload: //path
		case silk.FsDownloadAck: //id
		case silk.FsUpload: //path
		case silk.FsUploadAck: //id
		case silk.FsTransferData: //id,data
		case silk.FsTransferEnd: //id

		}
	})

}
