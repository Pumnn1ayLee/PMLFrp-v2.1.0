package main

import (
	"bufio"
	"database/sql"
	"fmt"
	"image/color"
	"log"
	"math/rand"
	"net"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"syscall"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/go-toast/toast"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

var IP = [...]string{"1.12.245.232"}
var ID = [...]string{"Guangzhou:"}
var REMOTE_PORT string
var KEY bool = false
var sign_in bool = false
var ms time.Duration
var a = app.New()
var w = a.NewWindow("SplashNirvana2K")
var tabs *container.AppTabs
var tabs_suc *container.AppTabs
var tabs_sign_up *container.AppTabs
var Content_All *fyne.Container
var Content_Login_All *fyne.Container
var Content_SignUp_All *fyne.Container

func Db_Insert(username string, userpassword string) bool {
	db, err := sql.Open("mysql", "root:159632@tcp(1.12.245.232:3306)/PMLUSERS")
	if err != nil {
		fmt.Println(err)
	}
	defer db.Close()

	UserName := fmt.Sprintf("'%s'", username)
	UserPassWord := fmt.Sprintf("'%s'", userpassword)

	dbQuery := fmt.Sprintf("INSERT INTO pml_users VALUES(user_id,%s,%s,NOW())", UserName, UserPassWord)

	_, err2 := db.Query(dbQuery)
	if err2 != nil {
		log.Fatal(err2)
		return false
	}
	return true
}

func Db_Select(username string, userpassword string) bool {
	db, err := sql.Open("mysql", "root:159632@tcp(1.12.245.232:3306)/PMLUSERS")
	if err != nil {
		fmt.Println(err)
	}
	defer db.Close()

	UserName := fmt.Sprintf("'%s'", username)
	UserPassWord := fmt.Sprintf("'%s'", userpassword)

	dbQuery := fmt.Sprintf("SELECT *from pml_users WHERE BINARY user_name=%s AND user_password=%s", UserName, UserPassWord)

	rows, err2 := db.Query(dbQuery)
	if err2 != nil {
		return false
	}
	defer rows.Close()

	for rows.Next() {
		var col1 int
		var col2 string
		var col3 string
		var col4 string

		err3 := rows.Scan(&col1, &col2, &col3, &col4)
		if err3 != nil {
			log.Fatal(err3)
		}
		if col2 == username && col3 == userpassword {
			return true
		} else {
			return false
		}
	}
	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("查询完成")
	return false
}

func create_remo_port() string {
	rand.Seed(time.Now().UnixNano())

	nums := rand.Intn(9000) + 1000

	port := strconv.Itoa(nums)

	return port
}


func Change_Text(port string) string {
	filePath := "bin/frp_0.44.0_windows_amd64/frpc.ini"
	file, err := os.OpenFile(filePath, os.O_RDWR, 0644)
	if err != nil {
		log.Fatalf("Failed to open file:%s", err)
	}
	defer file.Close()

	// 读取文件内容
	scanner := bufio.NewScanner(file)
	var lines []string
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		fmt.Println("Failed to read file:", err)
	}

	//修改指定内容的行
	var newLines []string
	var nums string = create_remo_port()
	for _, line := range lines {
		if strings.HasPrefix(line, "local_port = ") {
			line = "local_port = " + port
		}
		if strings.HasPrefix(line, "remote_port = ") {
			line = "remote_port = " + nums
			REMOTE_PORT = nums
		}
		newLines = append(newLines, line)
	}

	// 将修改后的内容写回文件
	file, err = os.Create("bin/frp_0.44.0_windows_amd64/frpc.ini")
	if err != nil {
		fmt.Println("Failed to create file:", err)
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	for _, line := range newLines {
		fmt.Fprintln(writer, line)
	}
	writer.Flush()

	// 打开修改后的文件
	file1, err := os.OpenFile(filePath, os.O_RDWR, 0644)
	if err != nil {
		log.Fatalf("Failed to open file:%s", err)
	}
	defer file1.Close()
	fmt.Println("file ok")
	fmt.Printf("File content after modification: %sand%s\n", port, nums)
	return nums
}

func Open_cmd(value bool, key bool) {
	path := "bin/frp_0.44.0_windows_amd64/"
	cmd := exec.Command("cmd", "/c", "frpc.exe")
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	cmd1 := exec.Command("taskkill", "/F", "/IM", "frpc.exe")

	cmd.Dir = path
	if value == true && key == true {
		cmd.Start()
	}
	if value == false && key == true {
		cmd1.Run()
	}

}

func Noti_Win_True() {
	ip := IP[0] + ":" + REMOTE_PORT
	notification := toast.Notification{
		AppID:   "Microsoft.Windows.Shell.RunDialog",
		Title:   "SplashNirvana2K",
		Message: "Your IP and Port is:" + ip,
		Actions: []toast.Action{
			{"protocol", "My GitHub", "https://github.com/Pumnn1ayLee/PMLfrp-v1.0"},
		},
	}
	err := notification.Push()
	if err != nil {
		log.Fatalln(err)
	}
}

func Noti_Win_False() {
	notification := toast.Notification{
		AppID:   "Microsoft.Windows.Shell.RunDialog",
		Title:   "SplashNirvana2K",
		Message: "Your IP and Port are exactly wrong!Please modify your IP and PORT! ",
		Actions: []toast.Action{
			{"protocol", "My GitHub", "https://github.com/Pumnn1ayLee/PMLfrp-v1.0"},
		},
	}
	err := notification.Push()
	if err != nil {
		log.Fatalln(err)
	}
}

func Noti_Set_False() {
	notification := toast.Notification{
		AppID:   "Microsoft.Windows.Shell.RunDialog",
		Title:   "SplashNirvana2K",
		Message: "Please Set Guangzhou in true and press Connect Button!",
		Actions: []toast.Action{
			{"protocol", "My GitHub", "https://github.com/Pumnn1ayLee/PMLfrp-v1.0"},
		},
	}
	err := notification.Push()
	if err != nil {
		log.Fatalln(err)
	}
}

// 服务器延迟方法
func Ping_Server() time.Duration {

	serverPort := "22"

	// 创建TCP连接
	conn, err := net.DialTimeout("tcp", IP[0]+":"+serverPort, 3*time.Second)
	if err != nil {
		fmt.Println("无法连接到服务器:", err)
		return 0
	}
	defer conn.Close()

	// 获取连接的远程地址
	remoteAddr := conn.RemoteAddr().(*net.TCPAddr)
	fmt.Printf("已连接到服务器 %s\n", remoteAddr.IP.String())

	// 发送Ping消息并测量延迟
	startTime := time.Now()
	_, err = conn.Write([]byte("Ping"))
	if err != nil {
		fmt.Println("发送Ping消息失败:", err)
		return 0
	}

	// 接收Pong消息
	response := make([]byte, 4)
	_, err = conn.Read(response)
	if err != nil {
		fmt.Println("接收Pong消息失败:", err)
		return 0
	}
	endTime := time.Now()

	// 计算延迟
	latency := endTime.Sub(startTime)
	fmt.Printf("延迟: %v\n", latency)
	return latency
}

func Init_Console(w fyne.Window) {
	//背景图片
	img := canvas.NewImageFromFile("bin/picture/Windows.jpg")
	//延迟参数返回
	ms = Ping_Server()


	//连接按钮
	connect_button := widget.NewButton("Connect", func() {
		if KEY == true {
			Noti_Win_True()
			Open_cmd(true, KEY)
		}
		if KEY == false{
			Noti_Set_False()
		}
	})
	//刷新按钮
	refresh_button := widget.NewButton("Refresh", func() {
		ms = Ping_Server()
		Init_Console(w)
		KEY = false
		if sign_in{
		w.SetContent(Content_Login_All)
		}
		if sign_in == false{
		w.SetContent(Content_All)
		}
	})

	button1 := container.NewVBox(connect_button, refresh_button)

	MS := ms.String()
	check := widget.NewCheck(ID[0]+MS, func(value2 bool) {
		log.Println("Check set to", value2)
		if value2 == true {
			KEY = true
		}
		if value2 == false {
			KEY = false
		}
	})

	text1 := widget.NewLabel("The list of available Channels is as follows:")

	content_channel := container.New(layout.NewBorderLayout(text1, button1, check, nil), layout.NewSpacer(), text1, check, button1)

	//Home项
	entry_user_id := widget.NewEntry()
	entry_user_id.SetPlaceHolder("Your User_Id:")

	entry_user_pas := widget.NewPasswordEntry()
	entry_user_pas.SetPlaceHolder("Your PassWord:")

	//成功登录画面
	suceText := canvas.NewText("Login Successful Welcome to SplashNirvana2K,Please wait a few times!", color.Black)
	content_suce := container.New(layout.NewCenterLayout(), suceText)

	//失败连接画面
	faiText := canvas.NewText("Your ID or Password maybe wrong,Please input in right or sign up", color.Black)
	content_fai := container.New(layout.NewCenterLayout(), faiText)

	//成功注册画面
	suce_up_Text := canvas.NewText("Sign up Successful Welcome to SplashNirvana2K,Please wait a few times!", color.Black)
	content_suce_up := container.New(layout.NewCenterLayout(), suce_up_Text)

	//失败注册画面
	fai_up_Text := canvas.NewText("Your Sign up maybe wrong!", color.Black)
	content_fai_up := container.New(layout.NewCenterLayout(), fai_up_Text)

	//成功登录连接后本地端口
	entry_loc_port := widget.NewEntry()
	entry_loc_port.SetPlaceHolder("Your Local Port:")

	//注册用户密码条
	entry_sign_up_id := widget.NewEntry()
	entry_sign_up_id.SetPlaceHolder("Your New Name:")

	entry_sign_up_pas := widget.NewPasswordEntry()
	entry_sign_up_pas.SetPlaceHolder("Your New Password")

	//本地端口按钮
	Local_Port_Button := widget.NewButton("OK", func() {
		Change_Text(entry_loc_port.Text)
	})

	//去往注册的按钮
	To_Sign_Up_Button := widget.NewButton("Sign Up", func() {
		w.SetContent(Content_SignUp_All)
	})

	//注册按钮
	Sign_Up_Button := widget.NewButton("OK", func() {
		fmt.Println("Your New User is:", entry_sign_up_id)
		fmt.Println("Your New User Password is:", entry_sign_up_pas)
		var signup bool = Db_Insert(entry_sign_up_id.Text, entry_sign_up_pas.Text)
		if signup {
			w.SetContent(content_suce_up)
			time.Sleep(2 * time.Second)
			w.SetContent(Content_All)
		}
		if signup == false {
			w.SetContent(content_fai_up)
			time.Sleep(2 * time.Second)
			w.SetContent(Content_All)
		}
	})
	//总窗格
	tabs = container.NewAppTabs(
		container.NewTabItem("Channel Lists", content_channel),
		container.NewTabItem("ReadMe Important!", widget.NewLabel("Please first open your MC and then open the LAN port,\n fill in the port and IP address to enter HOME, and then\n check the button in Channels!\n It's best not to repeatedly check!")),
		container.NewTabItem("About us", widget.NewLabel("Welcome to SplashNirvana2K and En_nuyeux\n\n\n Github:https://github.com/Pumnn1ayLee\n\nGithub:https://github.com/Ennuyeux233")),
	)

	//连接成功后的窗格
	tabs_suc = container.NewAppTabs(
		container.NewTabItem("Channel Lists", content_channel),
		container.NewTabItem("ReadMe Important!", widget.NewLabel("Please first open your MC and then open the LAN port,\n fill in the port and IP address to enter HOME, and then\n check the button in Channels!\n It's best not to repeatedly check!")),
		container.NewTabItem("About us", widget.NewLabel("Welcome to SplashNirvana2K and En_nuyeux\n\n\n Github:https://github.com/Pumnn1ayLee\n\nGithub:https://github.com/Ennuyeux233")),
	)

	content_suce_login := container.NewVBox(entry_loc_port, Local_Port_Button)

	tabs_suc.Append(container.NewTabItemWithIcon("Home", theme.HomeIcon(), content_suce_login))

	tabs_suc.SetTabLocation(container.TabLocationLeading)

	//注册窗格
	tabs_sign_up = container.NewAppTabs(
		container.NewTabItem("Channel Lists", content_channel),
		container.NewTabItem("ReadMe Important!", widget.NewLabel("Please first open your MC and then open the LAN port,\n fill in the port and IP address to enter HOME, and then\n check the button in Channels!\n It's best not to repeatedly check!")),
		container.NewTabItem("About us", widget.NewLabel("Welcome to SplashNirvana2K and En_nuyeux\n\n\n Github:https://github.com/Pumnn1ayLee\n\nGithub:https://github.com/Ennuyeux233")),
	)

	content_sign_up := container.NewVBox(entry_sign_up_id, entry_sign_up_pas, Sign_Up_Button)

	tabs_sign_up.Append(container.NewTabItemWithIcon("Home", theme.HomeIcon(), content_sign_up))

	tabs_sign_up.SetTabLocation(container.TabLocationLeading)

	Sign_In_Button := widget.NewButton("Sign In", func() {
		log.Println("Your UserID is:", entry_user_id.Text)
		log.Println("Input Password is:", entry_user_pas.Text)

		sign_in = Db_Select(entry_user_id.Text, entry_user_pas.Text)
		if sign_in {
			w.SetContent(content_suce)
			time.Sleep(2 * time.Second)
			w.SetContent(Content_Login_All)
		}
		if sign_in == false {
			w.SetContent(content_fai)
			time.Sleep(2 * time.Second)
			w.SetContent(Content_All)
		}
	})
	content := container.NewVBox(entry_user_id, entry_user_pas, Sign_In_Button, To_Sign_Up_Button)

	tabs.Append(container.NewTabItemWithIcon("Home", theme.HomeIcon(), content))

	tabs.SetTabLocation(container.TabLocationLeading)

	Content_Login_All = container.New(layout.NewMaxLayout(), img, tabs_suc)
	Content_SignUp_All = container.New(layout.NewMaxLayout(), img, tabs_sign_up)
	Content_All = container.New(layout.NewMaxLayout(), img, tabs)

	w.SetContent(Content_All)
	w.Resize(fyne.Size{Width: 500, Height: 300})
}

func main() {
	Init_Console(w)
	w.ShowAndRun()
}
