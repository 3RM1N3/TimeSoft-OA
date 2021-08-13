package lib

// 发送文件和向数据库中存储时的结构体
type DBArchiveRecord struct {
	FileID      string        // FileID      档号：由字母、数字和英文横杠 “-” 按一定规则组成
	ClientCo    string        // ClientCo    客户公司名称
	Year        int           // Year        年份
	ArchiveType byte          // ArchiveType 档案类型/类别：人事档案、文书档案、业务档案、基建档案、设备档案、 党群类、行政工作类、经营管理类
	StorageTime int           // StorageTime 保管期限（/年）：5(财务档案)、10、30、0(即为永久)
	Department  int           // Department  部门：组织部、财务部等等
	HeadOrBody  byte          // HeadOrBody  目录或正文：指该文件属于目录还是正文，包含两种值：目录、正文
	Scaner      string        // Scaner      扫描员工：员工ID
	ScanTime    int           // ScanTime    扫描上传时间
	Editor      string        // Editor      修图员工：员工ID
	EditTime    int           // EditTime    修图上传时间
	FileState   FileStateCode // FileState   文件状态
	FileName    int           // FileName    文件名：本地文件名
	SubFileNum  int           // 档号文件夹下的文件数
	Size        int64         // Size        压缩包大小
}

type FileStateCode byte

const (
	ScanOver       FileStateCode = iota // 扫描完毕待修图
	Editting                            // 修图中
	EditOver                            // 修图完毕未发生返工
	NeedRework                          // 待返工的扫描
	Reworking                           // 扫描返工中
	ReworkOver                          // 返工完毕待修图
	ReworkEditOver                      // 返工后再次修图完毕
)
