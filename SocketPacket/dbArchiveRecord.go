package SocketPacket

// ID            序列码：每条记录唯一，用于识别文件
//
// ArchiveType   档案类型/类别：人事档案、文书档案、业务档案、基建档案、设备档案、 党群类、行政工作类、经营管理类
//
// StorageTime   保管期限（/年）：5(财务档案)、10、30、0(即为永久)
//
// Department    部门：组织部、财务部等等
//
// HeadOrBody    目录或正文：指该文件属于目录还是正文，包含两种值：目录、正文
//
// FileID        档号：由字母、数字和英文横杠 “-” 按一定规则组成
//
// Scaner/Editor 扫描和修图员工：员工ID
//
// FileState     文件状态：扫描完毕待修图、修图中、修图完毕未发生返工、待返工的扫描、扫描返工中、返工完毕待修图、返工后再次修图完毕
//
// FileName      文件名：本地文件名
type DBArchiveRecord struct {
	ID          string
	Year        string
	ArchiveType int
	StorageTime int
	Department  int
	HeadOrBody  int
	FileID      int
	Scaner      string
	ScanTime    int
	Editor      string
	EditTime    int
	FileState   int
	FileName    int
}
