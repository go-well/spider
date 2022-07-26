package silk

type Type uint8

const (
	Close Type = iota
	Ping
	Pong
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
	Spawn
	SpawnAck
)

const (
	TunnelCreate    Type = iota + 20 //net,addr 例 tcp,127.0.0.1:8080
	TunnelCreateAck                  //id(uint16)
	TunnelClose                      //id
	TunnelCloseAck                   //id
	TunnelData                       //id,data
	TunnelDataAck                    //id
	TunnelEnd                        //id
	TunnelError                      //id,error

)

const (
	SystemShell      Type = iota + 40 // /bin/sh
	SystemShellAck                    //tunnel id(uint16)
	SystemExecute                     //command
	SystemExecuteAck                  //stdout
	SystemStart                       //command
	SystemStartAck                    //stdout
	SystemKill
	SystemKillAck
	SystemEnvironment
	SystemEnvironmentAck //json
	SystemConfig
	SystemConfigAck //yaml、json

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
	DatabaseQuery     Type = iota + 80 //sql
	DatabaseQueryAck                   //json
	DatabaseExec                       //sql
	DatabaseExecAck                    //json
	DatabaseMeta                       //
	DatabaseMetaAck                    //json
	DatabaseDriver                     //
	DatabaseDriverAck                  //text
	DatabaseSource                     //
	DatabaseSourceAck                  //text

)

const (
	FsList      Type = iota + 100 //path
	FsListAck                     //json
	FsMkDir                       //path
	FsMkDirAck                    //
	FsRemove                      //path
	FsRemoveAck                   //
	FsRename                      //path,path
	FsRenameAck                   //
	FsStats                       //path
	FsStatsAck                    //json

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

func no() {

}
