package webssh

import (
	"anylinker/common/log"
	"bufio"
	"bytes"
	"fmt"
	"github.com/mitchellh/go-homedir"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"golang.org/x/crypto/ssh"
	"io/ioutil"
	"os"
	"strings"
	"time"
)

// 仅支持用户密码登录
func NewSshClient(user, password, ip string) (*ssh.Client, error) {
	config := &ssh.ClientConfig{
		Timeout:         time.Second * 5,
		User:            user,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(), //这个可以， 但是不够安全
		//HostKeyCallback: hostKeyCallBackFunc(h.Host),
	}
	//if h.Type == "password" {
	config.Auth = []ssh.AuthMethod{ssh.Password(password)}
	//} else {
	//	config.Auth = []ssh.AuthMethod{publicKeyAuthFunc(h.Key)}
	//}
	fmt.Println(user,password,ip)
	addr := fmt.Sprintf("%s:%d", ip, 22)
	c, err := ssh.Dial("tcp", addr, config)
	if err != nil {
		return nil, err
	}
	return c, nil
}
func hostKeyCallBackFunc(host string) ssh.HostKeyCallback {
	hostPath, err := homedir.Expand("~/.ssh/known_hosts")
	if err != nil {
		log.Fatal("find known_hosts's home dir failed", zap.Error(err))
	}
	file, err := os.Open(hostPath)
	if err != nil {
		log.Fatal("can't find known_host file:", zap.Error(err))
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var hostKey ssh.PublicKey
	for scanner.Scan() {
		fields := strings.Split(scanner.Text(), " ")
		if len(fields) != 3 {
			continue
		}
		if strings.Contains(fields[0], host) {
			var err error
			hostKey, _, _, _, err = ssh.ParseAuthorizedKey(scanner.Bytes())
			if err != nil {
				log.Error("error parsing %q: %v", zap.Error(errors.New(fields[2])), zap.Error(err))
			}
			break
		}
	}
	if hostKey == nil {
		log.Error("no hostkey for %s,%v", zap.Error(errors.New(host)), zap.Error(err))
	}
	return ssh.FixedHostKey(hostKey)
}

func publicKeyAuthFunc(kPath string) ssh.AuthMethod {
	keyPath, err := homedir.Expand(kPath)
	if err != nil {
		log.Error("find key's home dir failed", zap.Error(err))
	}
	key, err := ioutil.ReadFile(keyPath)
	if err != nil {
		log.Error("ssh key file read failed", zap.Error(err))
	}
	// Create the Signer for this private key.
	signer, err := ssh.ParsePrivateKey(key)
	if err != nil {
		log.Error("ssh key signer failed", zap.Error(err))
	}
	return ssh.PublicKeys(signer)
}
func runCommand(client *ssh.Client, command string) (stdout string, err error) {
	session, err := client.NewSession()
	if err != nil {
		//log.Print(err)
		return
	}
	defer session.Close()

	var buf bytes.Buffer
	session.Stdout = &buf
	err = session.Run(command)
	if err != nil {
		//log.Print(err)
		return
	}
	stdout = string(buf.Bytes())

	return
}

