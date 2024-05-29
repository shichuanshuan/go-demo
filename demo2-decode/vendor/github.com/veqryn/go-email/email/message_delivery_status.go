// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package email

import (
	"bufio"
	"bytes"
	"errors"
	"io"
	"net/textproto"
)

// HasDeliveryStatusMessage returns true if this Message has a
// content type of "message/delivery-status" and has a non-nil SubMessage
// containing the delivery status information.
func (m *Message) HasDeliveryStatusMessage() bool {
	contentType, _, err := m.Header.ContentType()
	if err != nil {
		return false
	}
	return contentType == "message/delivery-status" && m.SubMessage != nil
}

// DeliveryStatusMessageDNS returns the message DNS information,
// or an error if HasDeliveryStatusMessage would return false.
func (m *Message) DeliveryStatusMessageDNS() (Header, error) {
	if !m.HasDeliveryStatusMessage() {
		return Header{}, errors.New("Message does not have media content of type message/delivery-status")
	}
	return m.SubMessage.Header, nil
}

// DeliveryStatusRecipientDNS returns the message recipients' DNS information,
// or an error if HasDeliveryStatusMessage would return false.
func (m *Message) DeliveryStatusRecipientDNS() ([]Header, error) {
	recipientDNS := make([]Header, 0, 1)
	if !m.HasDeliveryStatusMessage() {
		return recipientDNS, errors.New("Message does not have media content of type message/delivery-status")
	}
	var err error
	var recipientHeaders textproto.MIMEHeader
	tp := textproto.NewReader(bufio.NewReader(bytes.NewReader(m.SubMessage.Body)))
	for err != io.EOF {
		recipientHeaders, err = tp.ReadMIMEHeader()
		if err != nil && err != io.EOF {
			return nil, err
		}
		recipientDNS = append(recipientDNS, Header(recipientHeaders))
	}
	return recipientDNS, nil
}
