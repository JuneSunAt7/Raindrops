package client

import (
	"net"

	"time"

	"github.com/pterm/pterm"
)

func AllCertificates(conn net.Conn) {
	pterm.DefaultBasicText.WithStyle(pterm.NewStyle(pterm.FgBlue)).Println("Все сертификаты")

	conn.Write([]byte("certs\n")) // Bind with server
	buffer := make([]byte, 4096)
	n, _ := conn.Read(buffer)
	pterm.FgGreen.Println(string(buffer[:n])) // Read answer

}
func CreateCert(conn net.Conn) {
	conn.Write([]byte("create\n")) // Bind with server
	buffer := make([]byte, 4096)
	n, _ := conn.Read(buffer)
	if string(buffer[:n]) == "0" {
		pterm.DefaultBasicText.WithStyle(pterm.NewStyle(pterm.FgRed)).Println("Сертификат уже создан!")
	} else {
		p, _ := pterm.DefaultProgressbar.WithTotal(10).WithTitle("...Создание сертификата...").Start()

		for i := 0; i < p.Total; i++ {
			if i == 6 {
				time.Sleep(time.Second * 3) // ProgressBar - uploader
			}
			p.UpdateTitle("...Обработка данных...")
			pterm.Success.Println("Успешное создание сертификата")
			p.Increment()
			time.Sleep(time.Millisecond * 350)

		}
	}
}
