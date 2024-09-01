package main

import (
	"context"
	"fmt"
	"github.com/khicago/irr"
	"github.com/khicago/irr/irc"
	lark "github.com/larksuite/oapi-sdk-go/v3"

	larkcard "github.com/larksuite/oapi-sdk-go/v3/card"
	larkim "github.com/larksuite/oapi-sdk-go/v3/service/im/v1"
)

const (
	TextMsgType        = "text"
	ShareCardMsgType   = "share_chat"
	InteractiveMsgType = "interactive"

	BizError irc.Code = 400
)

var (
	ErrOAPIBizError = BizError.Error("oapi biz error")
)

type (
	LarkOAPIClient struct {
		client *lark.Client
	}
)

func (cli *LarkOAPIClient) IMSendCard(ctx context.Context, openID string, jsonContent string) (*larkim.CreateMessageResp, error) {
	// 发送卡片
	resp, err := cli.client.Im.Message.Create(context.Background(), larkim.NewCreateMessageReqBuilder().
		ReceiveIdType(larkim.ReceiveIdTypeOpenId).
		Body(larkim.NewCreateMessageReqBodyBuilder().
			MsgType(InteractiveMsgType).
			ReceiveId(openID).
			Content(jsonContent).
			Build()).
		Build())
	if err != nil {
		return nil, irr.Wrap(err, "send card error")
	}
	if resp.Code != 0 {
		return nil, irr.Wrap(ErrOAPIBizError, "send card failed, code= %d, msg= %s", resp.Code, resp.Msg).SetCode(int64(resp.Code))
	}
	return resp, nil
}

func main() {
	// 创建 API Client
	cli := LarkOAPIClient{
		client: lark.NewClient("cli_a64b25c10a38500d", "ArrK2p9QY5chdqrXStdQKNF2ZFGpO4Qs"),
	}

	openID := "ou_f58ef2ee869f60e752db75e88b1c3526"

	// 定义卡片内容
	messageCard := larkcard.NewMessageCard().
		Config(larkcard.NewMessageCardConfig().
			WideScreenMode(true).
			Build()).
		//Header(larkcard.NewMessageCardHeader().
		//	Title(larkcard.NewMessageCardPlainText().Content("卡片标题").Build()).
		//	Build()).
		Elements([]larkcard.MessageCardElement{
			larkcard.NewMessageCardDiv().
				Text(larkcard.NewMessageCardLarkMd().
					Content("这是卡片内容").
					Build()).
				Build(),
		}).CardLink(larkcard.NewMessageCardURL().Url("https://memnexus.kenv.tech/api/v1/aaa")).
		Build()

	content, err := messageCard.JSON()
	if err != nil {
		fmt.Println("生成卡片内容失败:", err)
		return
	}
	fmt.Println("卡片内容:", content)

	resp, err := cli.IMSendCard(context.Background(), openID, content)
	if err != nil {
		fmt.Println("发送卡片失败:", err)
		return
	}
	fmt.Println("发送卡片成功，消息ID:", resp.Data.MessageId)
}
