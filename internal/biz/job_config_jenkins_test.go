package biz

import (
	"context"
	"cron/internal/basic/tracing"
	"cron/internal/models"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"log"
	"strings"
	"testing"
)

func init() {
	//err := os.Chdir("E:/WorkApps/go/src/incron/")
	//if err != nil {
	//	log.Fatal(err)
	//}
	// 日志写入
	go tracing.MysqlCollectorListen()
}

func TestJenKins(t *testing.T) {
	conf := &models.CronConfig{
		Id:           0,
		Env:          "",
		EntryId:      0,
		Type:         0,
		Name:         "",
		Spec:         "",
		Protocol:     0,
		Command:      []byte(`{"cmd": "", "rpc": {"addr": "", "body": "", "proto": "", "action": "", "header": [], "method": "GRPC", "actions": []}, "sql": {"driver": "mysql", "source": {"id": 0, "port": "", "title": "", "database": "", "hostname": "", "password": "", "username": ""}, "interval": 0, "statement": [], "err_action": 1, "err_action_name": ""}, "http": {"url": "", "body": "", "header": [{"key": "", "value": ""}], "method": "GET"}, "jenkins": {"name": "card", "params": [{"key": "", "value": ""}], "source": {"id": 10}}}`),
		Remark:       "",
		Status:       0,
		StatusRemark: "",
		StatusDt:     "",
		UpdateDt:     "",
		CreateDt:     "",
		MsgSet:       nil,
	}

	ctx := context.Background()
	job := NewJobConfig(conf)

	job.jenkins(ctx, job.commandParse.Jenkins)
}

func TestJenkins2(t *testing.T) {
	html := `<html>
    <head>
        <meta http-equiv="Content-Type" content="text/html;charset=ISO-8859-1"/>
        <title>Error 404 Not Found</title>
    </head>
    <body>
        <h2>HTTP ERROR 404 Not Found</h2>
        <table>
            <tr>
                <th>URI:</th>
                <td>/queue/item/34129/</td>
            </tr>
            <tr>
                <th>STATUS:</th>
                <td>404</td>
            </tr>
            <tr>
                <th>MESSAGE:</th>
                <td>Not Found</td>
            </tr>
            <tr>
                <th>SERVLET:</th>
                <td>Stapler</td>
            </tr>
        </table>
        <hr />
        <a href="https://eclipse.org/jetty">Powered by Jetty:// 10.0.12</a>
        <hr />
    </body>
</html>
`

	dom, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		log.Fatalln(err)
	}

	title := ""
	dom.Find("head title").Each(func(i int, selection *goquery.Selection) {
		title = selection.Text()
		fmt.Println(selection.Text())
	})
	table := map[string]string{}
	dom.Find("body table tr").Each(func(i int, selection *goquery.Selection) {
		th := selection.Find("th").Text()
		table[th[:len(th)-1]] = selection.Find("td").Text()
		fmt.Println(selection.Text())
	})
	fmt.Println(title, table)

	queue := strings.Split(strings.Trim(table["URI"], "/"), "/")
	fmt.Println(queue, queue[len(queue)-1])

}
