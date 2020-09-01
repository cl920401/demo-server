package tools

import qrcode "github.com/skip2/go-qrcode"

func QRFile(url, file string, margin int) error {
	return qrcode.WriteFile(url, qrcode.Medium, margin, file)
}
