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
	TunnelCreate     Type = iota + 20 //net,addr 例 tcp,127.0.0.1:8080
	TunnelCreateAck                   //id(uint16)
	TunnelClose                       //id
	TunnelCloseAck                    //id
	TunnelData                        //id,data
	TunnelDataAck                     //id
	TunnelDataEnd                     //id
	TunnelDataEndAck                  //id
	TunnelError                       //id,error

)

const (
	TaskCreate    Type = iota + 40 //command
	TaskCreateAck                  //id(uint16)
	TaskData                       //id,data
	TaskDataAck                    //id
	TaskDataEnd                    //id
	TaskKill                       //id
	TaskKillAck                    //id
	TaskRun                        //commmand
	TaskRunAck                     //output
	TaskStart                      //commmand
	TaskStartAck                   //

)

const (
	StatsHost    Type = iota + 60
	StatsHostAck      //json
	StatsCpu
	StatsCpuAck //json
	StatsCpuUsage
	StatsCpuUsageAck //json
	StatsMem
	StatsMemAck //json
	StatsDisk
	StatsDiskAck      //json
	StatsDiskUsage    //path
	StatsDiskUsageAck //json
	StatsDiskIO       //path
	StatsDiskIOAck    //json
	StatsNet
	StatsNetAck //json
	StatsNetIO
	StatsNetIOAck //json
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
	DatabaseDump                       //
	DatabaseDumpAck                    //filename
	DatabaseImport                     //filename
	DatabaseImportAck                  //

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

	FsUpload      //path
	FsUploadAck   //id
	FsDownload    //path
	FsDownloadAck //id
	FsData        //id,data
	FsDataAck     //id
	FsDataEnd     //id
	FsDataEndAck  //id
)

func no() {

}
