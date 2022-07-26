package silk

type Type uint8

const (
	Close Type = iota
	Heartbeat
	Connect //json{..info...}
	ConnectAck
	Subscribe //topic
	SubscribeAck
	Unsubscribe //topic
	UnsubscribeAck
	Publish //topic, message
	PublishAck
	Message //topic, message

)

const (
	TunnelCreate       Type = iota + 20 //net,addr 例 tcp,127.0.0.1:8080
	TunnelCreateAck                     //id(uint16)
	TunnelClose                         //id
	TunnelCloseAck                      //id
	TunnelTransferData                  //id,data
	//TunnelTransferEnd

)

const (
	SystemShell      Type = iota + 40 // /bin/sh
	SystemShellAck                    //tunnel id(uint16)
	SystemExecute                     //command string
	SystemExecuteAck                  //stdout
	SystemKill
	SystemKillAck
	SystemConfig
	SystemConfigAck  //yaml
	SystemDbQuery    //sql
	SystemDbQueryAck //json
	SystemDbExec     //sql
	SystemDbExecAck  //text

)

const (
	StatsHost    Type = iota + 60
	StatsHostAck      //json
	StatsCpu
	StatsCpuAck //json
	StatsCpuTimes
	StatsCpuTimesAck //json
	StatsMem
	StatsMemAck //json
	StatsDisk
	StatsDiskAck      //json
	StatsDiskUsage    //path
	StatsDiskUsageAck //json
	StatsNet
	StatsNetAck //json
	StatsConnection
	StatsConnectionAck //json
	StatsUser
	StatsUserAck //json
)

const (
	FsList      Type = iota + 80 //path
	FsListAck                    //json
	FsMkDir                      //path
	FsMkDirAck                   //
	FsRemove                     //path
	FsRemoveAck                  //
	FsRename                     //path,path
	FsRenameAck                  //
	FsStats                      //path
	FsStatsAck                   //json

	FsDownload           //path
	FsDownloadContent    //id,data
	FsDownloadContentAck //id
	FsDownloadEnd        //id

	FsUpload           //path
	FsUploadAck        //id
	FsUploadContent    //id,data
	FsUploadContentAck //id
	FsUploadEnd        //id
	FsUploadEndAck     //id

)

const (
	ext Type = iota + 100
)

func no() {

}
