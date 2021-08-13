## 文件：Client\GUI\code\FormLogin.go

#### func NewFormLogin(owner vcl.IComponent) (root *TFormLogin)  

由res2go IDE插件自动生成，不要编辑。

vcl.Application.CreateForm(&FormLogin)

## 文件：Client\GUI\code\FormLoginImpl.go

#### func (f *TFormLogin) OnButtonLoginClick(sender vcl.IObject) 

::private::

#### func (f *TFormLogin) OnFormCreate(sender vcl.IObject) 



#### func (f *TFormLogin) OnLabel1Click(sender vcl.IObject) 



## 文件：Client\GUI\code\MainForm.go

#### func NewMainForm(owner vcl.IComponent) (root *TMainForm)  

由res2go IDE插件自动生成，不要编辑。

vcl.Application.CreateForm(&MainForm)

## 文件：Client\GUI\code\MainFormImpl.go

#### func (f *TMainForm) OnFormCreate(sender vcl.IObject) 

::private::

#### func (f *TMainForm) OnPageControl1Change(sender vcl.IObject) 



#### func (f *TMainForm) RefreshPageScan() 



## 文件：Client\GUI\code\main.go

#### func main() 



## 文件：Client\lib.go

#### func SignUpAccount(address string, signupJson lib.SignUpJson) error 

注册账号

#### func Login(address string, loginJson lib.LoginJson) error 

用户登录

#### func sendUDPMsg(address string, packType lib.PacketType, jsonStruct interface{}) error 

发送udp消息，jsonStruct为要发送的结构体

#### func GetMacAddrs() ([]string, error) 

获取本机MAC地址

#### func ClientReceiveFile(fileList []string, conn *net.TCPConn) (string, *os.File, error) 

客户端从远端接收文件

## 文件：Client\lib_test.go

#### func TestGetMacAddrs(t *testing.T) 



## 文件：Client\main.go

#### func main() 



## 文件：Client\scanTask.go

#### func SetProjectDir(dirPath string) error 

设置项目文件夹

#### func CheckVerf(dirPath string) (bool, error) 

校验.verf文件

#### func GenVerf(dirPath string, fileNum int) (string, error) 

生成.verf校验字符串

fileNum设置为-1，则读取文件夹生成字符串，否则使用fileNum值

#### func DirWatcher(dirPath string, increaseCh chan int, errCh chan error, scanOver *bool) 

监测项目文件夹

#### func GetClientCo() []string 

获取当前已有客户

## 文件：Client\scanTask_test.go

#### func TestGenVerf(t *testing.T) 



#### func TestCheckVerf(t *testing.T) 



#### func TestSetProjectDir(t *testing.T) 



#### func TestSetProjectDir1(t *testing.T) 



#### func TestSetProjectDir2(t *testing.T) 



## 文件：Server\lib.go

#### func GetUnreviewUser(unreviewedUserList *[][2]string) error 

获取未审核员工信息

#### func AllPass() error 

令未审核的员工全部通过

#### func PartPass(unreviewedList *[][2]string, passedIndex *[]int) error 

部分通过

#### func intListContains(l []int, i int) bool 



#### func ResetPwd(phone, newPwd string) error 

重置密码

## 文件：Server\lib_test.go

#### func TestGetUnreviewUser(t *testing.T) 

测试获取未审核用户

#### func TestAllPass(t *testing.T) 

测试全部员工通过注册

#### func TestPartPass(t *testing.T) 

测试部分员工通过注册

#### func TestResetPwd(t *testing.T) 

测试重置密码

## 文件：Server\main.go

#### func init() 



#### func main() 



## 文件：Server\socketServers.go

#### func TCPServer(address string) error 

收发文件的TCP服务器

#### func SignupAndLogin() 

用于注册和登录的udp服务器

#### func processTCPConn(conn net.Conn) 

处理TCP连接

#### func processUDPMsg(udpMessage []byte, addr *net.UDPAddr, listen *net.UDPConn) 

处理UDP消息 待优化

#### func processSignin(udpMessage []byte, addr *net.UDPAddr, listen *net.UDPConn) 

处理用户注册

#### func processLogin(udpMessage []byte, addr *net.UDPAddr, listen *net.UDPConn) 

处理用户登录

#### func ServerSendFile(conn net.Conn) error 

TCP服务器发送文件给远端

## 文件：createDoc.go

#### func main() 



#### func walkFunc(path string, info fs.FileInfo, err error) error 



#### func genDoc(path string) error 



#### func printAndWrite(s string) 



## 文件：lib\JsonStructs.go

## 文件：lib\aboutFile.go

#### func (h *FileSendHead) MakeHead() ([]byte, error) 

上传文件的请求头

下载文件的请求头

#### func (h FileReceiveHead) MakeHead() ([]byte, error) 



#### func MakeHead(SRType PacketType, some interface{}) ([]byte, error) 

创建请求头字节切片

#### func SendFile(fileName string, conn net.Conn) error 

发送文件至远端

#### func ReceiveFile(conn net.Conn) (string, *os.File, error) 

从远端接收文件，返回本地文件名，文件指针和错误

## 文件：lib\dbArchiveRecord.go

## 文件：lib\functions.go

#### func Uint16ToByte(i uint16) ([]byte, error) 



#### func ByteToUint16(b []byte) (uint16, error) 



#### func Uint32ToByte(i uint32) ([]byte, error) 



#### func ByteToUint32(b []byte) (uint32, error) 



#### func Zip(srcFileOrDir, destZip string) error 

Zip 将文件或目录压缩为.zip文件

#### func Unzip(srcZip, destDir string) error 

Unzip 将.zip文件解压至目录，如果目录不存在则自动创建

#### func GbkToUtf8(s []byte) (string, error) 



#### func MD5(s string) string 

计算字符串的md5值

## 文件：lib\lib.go

#### func (c ReportCode) Pack() []byte 



#### func (c ReportCode) ToError() error 



## 文件：lib\lib_test.go

#### func TestByteToUint16(t *testing.T) 



#### func TestUint16ToByte(t *testing.T) 



## 文件：srcCodeLines.go

#### func main() 



#### func callback(path string, info fs.FileInfo, err error) error 



#### func countByte(byteSlice []byte, b byte) int 



## 文件：test\main.go

#### func main() 



#### func SignupAndLogin() 

用于注册和登录的udp服务器

#### func b() 

测试解压文件

