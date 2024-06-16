/* 
Package main пакет для запуска клиентского приложения
*/
package main

import (
	"github.com/KartoonYoko/gophkeeper/internal/app/cliclient"
)

var buildVersion string = "N/A"
var buildDate string = "N/A"

func main() {
	vi := &cliclient.VersionInfo{
		BuildDate: buildDate,
		Version:   buildVersion,
	}
	cliclient.Run(vi)
}
