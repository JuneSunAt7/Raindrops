package server

import (
	"bufio"
	"encoding/json"
	"errors"
	"io"
	"net"
	"os"

	"github.com/JuneSunAt7/Raindrops/logger"
	"github.com/pterm/pterm"
)
var SESSION_PASSWD string
type Credentials struct {
    Username string `json:"username"`
    Password string `json:"password"`
}

type CredArr []Credentials

func (p *CredArr) FromJSON(r io.Reader) error {
    en := json.NewDecoder(r)
    return en.Decode(p)
}

var Uname string

func GetCred() (*CredArr, error) {
    f, err := os.Open("user_creds.db")
    if err != nil {
        return nil, err
    }
    defer f.Close()
    var creds CredArr
    err = creds.FromJSON(f)
    if err != nil {
        return nil, err
    }

    return &creds, nil
}

func AuthenticateClient(conn net.Conn) error {
    creds, err := GetCred()
    if err != nil {
        return err
    }
    if len(*creds) == 0 {
        pterm.Error.Println("Нет ни одного зарегистрированного пользователя")
        return errors.New("Нет ни одного зарегистрированного пользователя")
    }

    reader := bufio.NewScanner(conn)

    // Validate user
    reader.Scan()
    uname := reader.Text()
    Uname = uname

    if CheckUserCert(Uname) {
        logger.Println("Новое подключение", uname, "проверен")
        conn.Write([]byte("1"))
        return nil
    } else {
        conn.Write([]byte("0"))

        reader.Scan()
        passwd := reader.Text()

        for _, cred := range *creds {
            if cred.Username == uname && cred.Password == passwd {
                logger.Println("Новое подключение ", uname)
                SESSION_PASSWD = passwd
                conn.Write([]byte("1"))
                return nil
            }
        }
    }

    conn.Write([]byte("0"))
    return nil
}

