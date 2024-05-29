package main

import (
	"bytes"
	"fmt"
	"github.com/axgle/mahonia"
	"github.com/veqryn/go-email/email"
	"golang.org/x/text/encoding/simplifiedchinese"
	"io"
	"io/ioutil"
	"log"
	"mime"
	"strings"
	"unicode/utf8"
)

func main() {
	// MTIzysfFtrChyrG54rW5wffI/XNsyse1xLncwO274bzGa2dmZGoNCg0KDQoNCnNoaWNzQHNvbWUu
	//Y29tDQo=

	// yP0NCg0KDQoNCnNoaWNzQHNvbWUuY29tDQo=
	gbkAndutf := `Date: Tue, 28 May 2024 14:00:00 +0800
From: "shics@some.com" <shics@some.com>
To: wu <wu@some.com>
Subject: neirong
X-Priority: 3
X-Has-Attach: yes
X-Mailer: Foxmail 7.2.25.254[cn]
Mime-Version: 1.0
Message-ID: <2024052813394238477732@some.com>
Content-Type: multipart/mixed;
	boundary="----=_001_NextPart575676200741_=----"

This is a multi-part message in MIME format.

------=_001_NextPart575676200741_=----
Content-Type: multipart/alternative;
	boundary="----=_002_NextPart477020825382_=----"


------=_002_NextPart477020825382_=----
Content-Type: text/plain;
	charset="GB2312"
Content-Transfer-Encoding: base64

MzI0DQoNCg0KDQpzaGljc0Bzb21lLmNvbQ0K

------=_002_NextPart477020825382_=----
Content-Type: text/html;
	charset="GB2312"
Content-Transfer-Encoding: quoted-printable

<html></html>
------=_002_NextPart477020825382_=------

------=_001_NextPart575676200741_=----
Content-Type: application/octet-stream;
	name="=?GB2312?B?yq+088uoLnR4dA==?="
Content-Transfer-Encoding: base64
Content-Disposition: attachment;
	filename="=?gb2312?B?yscudHh0?="

5piv5b635Zu95qCTZmtsamfmlLbliLA=

------=_001_NextPart575676200741_=------

.`

	utfAndGbk := `Date: Wed, 29 May 2024 13:26:25 +0800
From: "shics@some.com" <shics@some.com>
To: wu <wu@some.com>
Subject: neirong
X-Priority: 3
X-Has-Attach: yes
X-Mailer: Foxmail 7.2.25.254[cn]
Mime-Version: 1.0
Message-ID: <2024052816534267913365@some.com>
Content-Type: multipart/mixed;
	boundary="----=_001_NextPart440335318122_=----"

This is a multi-part message in MIME format.

------=_001_NextPart440335318122_=----
Content-Type: multipart/alternative;
	boundary="----=_002_NextPart704270113574_=----"


------=_002_NextPart704270113574_=----
Content-Type: text/plain;
	charset="utf-8"
Content-Transfer-Encoding: base64

MzI0DQoNCg0KDQpzaGljc0Bzb21lLmNvbQ0K

------=_002_NextPart704270113574_=----
Content-Type: text/html;
	charset="utf-8"
Content-Transfer-Encoding: quoted-printable

<html></html>
------=_002_NextPart704270113574_=------

------=_001_NextPart440335318122_=----
Content-Type: application/octet-stream;
	name="=?utf-8?B?55+z5aSn5qCTLnR4dA==?="
Content-Transfer-Encoding: base64
Content-Disposition: attachment;
	filename="=?utf-8?B?55+z5aSn5qCTLnR4dA==?="

ytW1vbSry6jKx7XCufo=

------=_001_NextPart440335318122_=------

.`

	var rawEmail *strings.Reader
	rawEmail = strings.NewReader(gbkAndutf)
	rawEmail = strings.NewReader(utfAndGbk)
	//rawEmail = strings.NewReader(gbkAndutf)

	//mail.ReadMessage()

	// 解析邮件
	parsedEmail, err := email.ParseMessage(rawEmail)
	if err != nil {
		log.Fatal(err)
	}

	// 探究如何实现的
	decode(parsedEmail.Header.Subject())

	dec := simplifiedchinese.GBK.NewDecoder()
	for _, part := range parsedEmail.Parts {
		if _, key := part.Header["Content-Disposition"]; key {
			var body []byte
			body = part.Body

			res := utf8.Valid(body)
			if res {
				fmt.Println(string(part.Body))
			} else {
				data, _ := dec.Bytes(part.Body)
				fmt.Printf("Body:%s\n", data)
			}
		}
	}
}

func decode(data string) string {
	dec := new(mime.WordDecoder)

	dec.CharsetReader = func(charset string, input io.Reader) (io.Reader, error) {
		fmt.Printf("charset%v\n", charset)

		switch charset {
		case "gb2312":
			content, err := ioutil.ReadAll(input)
			if err != nil {
				return nil, err
			}

			utf8str := ConvertToString(string(content), "gbk", "utf-8")
			t := bytes.NewReader([]byte(utf8str))
			return t, nil
		case "gb18030":
			fmt.Println()

		default:
			return nil, fmt.Errorf("unhandle charset")
		}

		return nil, nil
	}

	ret, err := dec.Decode(data)
	if err != nil {
		ret, _ = dec.DecodeHeader(data)
	}

	return ret
}

func ConvertToString(src string, srcCode string, tagCode string) string {
	srcCoder := mahonia.NewDecoder(srcCode)
	srcResult := srcCoder.ConvertString(src)
	tagCoder := mahonia.NewDecoder(tagCode)
	_, cdata, _ := tagCoder.Translate([]byte(srcResult), true)
	result := string(cdata)
	return result
}
