package ws

import (
	"anylinker/core/webssh"
	"bytes"
	"strconv"

	"github.com/gofiber/websocket/v2"
	"github.com/sirupsen/logrus"
)


// handle webSocket connection.
// first,we establish a ssh connection to ssh server when a webSocket comes;
// then we deliver ssh data via ssh connection between browser and ssh server.
// That is, read webSocket data from browser (e.g. 'ls' command) and send data to ssh server via ssh connection;
// the other hand, read returned ssh data from ssh server and write back to browser via webSocket API.
func WsSsh(c *websocket.Conn) {

	//cIp := c.ClientIP()
	//userM, err := getAuthUser(c)
	//if handleError(c, err) {
	//	return
	//}
	cols, err := strconv.Atoi(c.Query("cols", "120"))
	if webssh.WshandleError(c, err) {
		return
	}
	rows, err := strconv.Atoi(c.Query("rows", "32"))
	if webssh.WshandleError(c, err) {
		return
	}
	username := c.Query("username", "root")
	password := c.Query("password", "")
	addr := c.Query("addr", "")
	//idx, err := parseParamID(c)
	//if wshandleError(wsConn, err) {
	//	return
	//}
	//mc, err := models.MachineFind(idx)
	//if wshandleError(wsConn, err) {
	//	return
	//}
	client, err := webssh.NewSshClient(username,password,addr)
	if webssh.WshandleError(c, err) {
		return
	}
	defer client.Close()
	//startTime := time.Now()

	ssConn, err := webssh.NewSshConn(cols, rows, client)

	if webssh.WshandleError(c, err) {
		return
	}
	defer ssConn.Close()

	quitChan := make(chan bool, 3)

	var logBuff = new(bytes.Buffer)
	// most messages are ssh output, not webSocket input
	go ssConn.ReceiveWsMsg(c, logBuff, quitChan)
	go ssConn.SendComboOutput(c, quitChan)
	go ssConn.SessionWait(quitChan)

	<-quitChan
	//write logs
	//xtermLog := models.TermLog{
	//	EndTime:     time.Now(),
	//	StartTime:   startTime,
	//	UserId:      userM.ID,
	//	Log:         logBuff.String(),
	//	MachineId:   idx,
	//	MachineName: mc.Name,
	//	MachineIp:   mc.Ip,
	//	MachineHost: mc.Host,
	//	UserName:    userM.Username,
	//	Ip:          cIp,
	//}
	//
	//err = xtermLog.Create()
	//if wshandleError(wsConn, err) {
	//	return
	//}
	logrus.Info("websocket finished")
}
