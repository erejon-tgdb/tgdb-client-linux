package pdu

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"github.com/TIBCOSoftware/tgdb-client/client/goAPI/exception"
	"github.com/TIBCOSoftware/tgdb-client/client/goAPI/iostream"
	"github.com/TIBCOSoftware/tgdb-client/client/goAPI/types"
	"reflect"
	"strings"
)

/**
 * Copyright 2018-19 TIBCO Software Inc. All rights reserved.
 *
 * Licensed under the Apache License, Version 2.0 (the "License"); You may not use this file except
 * in compliance with the License.
 * A copy of the License is included in the distribution package with this file.
 * You also may obtain a copy of the License at http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF DirectionAny KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 * File name: VerbHandShakeRequest.go
 * Created on: Sep 30, 2018
 * Created by: achavan
 * SVN id: $id: $
 *
 */

//type requestType int

const (
	InvalidRequest = iota
	InitiateRequest
	ChallengeAccepted
)

type HandShakeRequestMessage struct {
	*AbstractProtocolMessage
	sslMode       bool
	challenge     int
	handshakeType int
}

func DefaultHandShakeRequestMessage() *HandShakeRequestMessage {
	// We must register the concrete type for the encoder and decoder (which would
	// normally be on a separate machine from the encoder). On each end, this tells the
	// engine which concrete type is being sent that implements the interface.
	gob.Register(HandShakeRequestMessage{})

	newMsg := HandShakeRequestMessage{
		AbstractProtocolMessage: DefaultAbstractProtocolMessage(),
	}
	newMsg.IsUpdatable = true
	newMsg.VerbId = VerbHandShakeRequest
	newMsg.sslMode = false
	newMsg.challenge = 0
	newMsg.handshakeType = InvalidRequest
	newMsg.BufLength = int(reflect.TypeOf(newMsg).Size())
	return &newMsg
}

// Create New Message Instance
func NewHandShakeRequestMessage(authToken, sessionId int64) *HandShakeRequestMessage {
	newMsg := DefaultHandShakeRequestMessage()
	newMsg.AuthToken = authToken
	newMsg.SessionId = sessionId
	newMsg.BufLength = int(reflect.TypeOf(*newMsg).Size())
	return newMsg
}

/////////////////////////////////////////////////////////////////
// Helper functions for HandShakeRequestMessage
/////////////////////////////////////////////////////////////////

func (msg *HandShakeRequestMessage) GetSslMode() bool {
	return msg.sslMode
}

func (msg *HandShakeRequestMessage) GetChallenge() int {
	return msg.challenge
}

func (msg *HandShakeRequestMessage) GetRequestType() int {
	return msg.handshakeType
}

func (msg *HandShakeRequestMessage) SetSslMode(mode bool) {
	msg.sslMode = mode
}

func (msg *HandShakeRequestMessage) SetChallenge(challenge int) {
	msg.challenge = challenge
}

func (msg *HandShakeRequestMessage) SetRequestType(rType int) {
	msg.handshakeType = rType
}

/////////////////////////////////////////////////////////////////
// Implement functions from Interface ==> TGMessage
/////////////////////////////////////////////////////////////////

// FromBytes constructs a message object from the input buffer in the byte format
func (msg *HandShakeRequestMessage) FromBytes(buffer []byte) (types.TGMessage, types.TGError) {
	logger.Log(fmt.Sprint("Entering HandShakeRequestMessage:FromBytes"))
	if len(buffer) < 0 {
		logger.Error(fmt.Sprint("ERROR: Returning HandShakeRequestMessage:FromBytes w/ Error: Invalid Message Buffer"))
		return nil, exception.CreateExceptionByType(types.TGErrorInvalidMessageLength)
	}

	is := iostream.NewProtocolDataInputStream(buffer)

	// First member attribute / element of message header is BufLength
	bufLen, err := is.ReadInt()
	if err != nil {
		logger.Error(fmt.Sprint("ERROR: Returning HandShakeRequestMessage:FromBytes w/ Error in reading buffer length from message buffer"))
		return nil, err
	}
	logger.Log(fmt.Sprintf("Inside HandShakeRequestMessage:FromBytes read bufLen as '%+v'", bufLen))
	if bufLen != len(buffer) {
		errMsg := fmt.Sprint("Buffer length mismatch")
		return nil, exception.GetErrorByType(types.TGErrorInvalidMessageLength, types.INTERNAL_SERVER_ERROR, errMsg, "")
	}

	logger.Log(fmt.Sprint("Inside HandShakeRequestMessage:FromBytes - about to readHeader"))
	err = msg.readHeader(is)
	if err != nil {
		errMsg := fmt.Sprintf("Unable to recreate message from '%+v' in byte format", buffer)
		return nil, exception.GetErrorByType(types.TGErrorIOException, types.INTERNAL_SERVER_ERROR, errMsg, "")
	}

	logger.Log(fmt.Sprint("Inside HandShakeRequestMessage:FromBytes - about to ReadPayload"))
	err = msg.ReadPayload(is)
	if err != nil {
		errMsg := fmt.Sprintf("Unable to recreate message from '%+v' in byte format", buffer)
		return nil, exception.GetErrorByType(types.TGErrorIOException, types.INTERNAL_SERVER_ERROR, errMsg, "")
	}

	logger.Log(fmt.Sprintf("HandShakeRequestMessage::FromBytes resulted in '%+v'", msg))
	return msg, nil
}

// ToBytes converts a message object into byte format to be sent over the network to TGDB server
func (msg *HandShakeRequestMessage) ToBytes() ([]byte, int, types.TGError) {
	logger.Log(fmt.Sprint("Entering HandShakeRequestMessage:ToBytes"))
	os := iostream.DefaultProtocolDataOutputStream()

	logger.Log(fmt.Sprint("Inside HandShakeRequestMessage:ToBytes - about to writeHeader"))
	err := msg.writeHeader(os)
	if err != nil {
		errMsg := fmt.Sprintf("Unable to export message '%+v' in byte format", msg)
		return nil, -1, exception.GetErrorByType(types.TGErrorIOException, types.INTERNAL_SERVER_ERROR, errMsg, "")
	}

	logger.Log(fmt.Sprint("Inside HandShakeRequestMessage:ToBytes - about to WritePayload"))
	err = msg.WritePayload(os)
	if err != nil {
		errMsg := fmt.Sprintf("Unable to export message '%+v' in byte format", msg)
		return nil, -1, exception.GetErrorByType(types.TGErrorIOException, types.INTERNAL_SERVER_ERROR, errMsg, "")
	}

	_, err = os.WriteIntAt(0, os.GetLength())
	if err != nil {
		logger.Error(fmt.Sprint("ERROR: Returning HandShakeRequestMessage:ToBytes w/ Error in writing buffer length"))
		return nil, -1, err
	}
	logger.Log(fmt.Sprintf("HandShakeRequestMessage::ToBytes results bytes-on-the-wire in '%+v'", os.GetBuffer()))
	return os.GetBuffer(), os.GetLength(), nil
}

// GetAuthToken gets the authToken
func (msg *HandShakeRequestMessage) GetAuthToken() int64 {
	return msg.getAuthToken()
}

// GetIsUpdatable checks whether this message updatable or not
func (msg *HandShakeRequestMessage) GetIsUpdatable() bool {
	return msg.getIsUpdatable()
}

// GetMessageByteBufLength gets the MessageByteBufLength. This method is called after the toBytes() is executed.
func (msg *HandShakeRequestMessage) GetMessageByteBufLength() int {
	return msg.getMessageByteBufLength()
}

// GetRequestId gets the requestId for the message. This will be used as the CorrelationId
func (msg *HandShakeRequestMessage) GetRequestId() int64 {
	return msg.getRequestId()
}

// GetSequenceNo gets the sequenceNo of the message
func (msg *HandShakeRequestMessage) GetSequenceNo() int64 {
	return msg.getSequenceNo()
}

// GetSessionId gets the session id
func (msg *HandShakeRequestMessage) GetSessionId() int64 {
	return msg.getSessionId()
}

// GetTimestamp gets the Timestamp
func (msg *HandShakeRequestMessage) GetTimestamp() int64 {
	return msg.getTimestamp()
}

// GetVerbId gets verbId of the message
func (msg *HandShakeRequestMessage) GetVerbId() int {
	return msg.getVerbId()
}

// SetAuthToken sets the authToken
func (msg *HandShakeRequestMessage) SetAuthToken(authToken int64) {
	msg.setAuthToken(authToken)
}

// SetRequestId sets the request id
func (msg *HandShakeRequestMessage) SetRequestId(requestId int64) {
	msg.setRequestId(requestId)
}

// SetSessionId sets the session id
func (msg *HandShakeRequestMessage) SetSessionId(sessionId int64) {
	msg.setSessionId(sessionId)
}

// SetTimestamp sets the timestamp
func (msg *HandShakeRequestMessage) SetTimestamp(timestamp int64) types.TGError {
	return msg.setTimestamp(timestamp)
}

func (msg *HandShakeRequestMessage) String() string {
	var buffer bytes.Buffer
	buffer.WriteString("HandShakeRequestMessage:{")
	buffer.WriteString(fmt.Sprintf("SslMode: %+v", msg.sslMode))
	buffer.WriteString(fmt.Sprintf(", Challenge: %d", msg.challenge))
	buffer.WriteString(fmt.Sprintf(", HandshakeType: %d", msg.handshakeType))
	buffer.WriteString(fmt.Sprintf(", BufLength: %d", msg.BufLength))
	strArray := []string{buffer.String(), msg.messageToString()+"}"}
	msgStr := strings.Join(strArray, ", ")
	return  msgStr
}

// UpdateSequenceAndTimeStamp updates the SequenceAndTimeStamp, if message is mutable
// @param timestamp
// @return TGMessage on success, error on failure
func (msg *HandShakeRequestMessage) UpdateSequenceAndTimeStamp(timestamp int64) types.TGError {
	return msg.updateSequenceAndTimeStamp(timestamp)
}

// ReadHeader reads the bytes from input stream and constructs a common header of network packet
func (msg *HandShakeRequestMessage) ReadHeader(is types.TGInputStream) types.TGError {
	return msg.readHeader(is)
}

// WriteHeader exports the values of the common message header attributes to output stream
func (msg *HandShakeRequestMessage) WriteHeader(os types.TGOutputStream) types.TGError {
	return msg.writeHeader(os)
}

// ReadPayload reads the bytes from input stream and constructs message specific payload attributes
func (msg *HandShakeRequestMessage) ReadPayload(is types.TGInputStream) types.TGError {
	logger.Log(fmt.Sprint("Entering HandShakeRequestMessage:ReadPayload"))
	//For Testing purpose only.
	rType, err := is.(*iostream.ProtocolDataInputStream).ReadByte()
	if err != nil {
		logger.Error(fmt.Sprint("ERROR: Returning HandShakeRequestMessage:ReadPayload w/ Error in reading rType from message buffer"))
		return err
	}
	logger.Log(fmt.Sprintf("Inside HandShakeRequestMessage:ReadPayload read rType as '%+v'", rType))

	mode, err := is.(*iostream.ProtocolDataInputStream).ReadBoolean()
	if err != nil {
		logger.Error(fmt.Sprint("ERROR: Returning HandShakeRequestMessage:ReadPayload w/ Error in reading mode from message buffer"))
		return err
	}
	logger.Log(fmt.Sprintf("Inside HandShakeRequestMessage:ReadPayload read mode as '%+v'", mode))

	challenge, err := is.(*iostream.ProtocolDataInputStream).ReadInt()
	if err != nil {
		logger.Error(fmt.Sprint("ERROR: Returning HandShakeRequestMessage:ReadPayload w/ Error in reading challenge from message buffer"))
		return err
	}
	logger.Log(fmt.Sprintf("HandShakeRequestMessage:ReadPayload read challenge as '%+v'", challenge))

	msg.SetRequestType(int(rType))
	msg.SetSslMode(mode)
	msg.SetChallenge(challenge)
	logger.Log(fmt.Sprint("Returning HandShakeRequestMessage:ReadPayload"))
	return nil
}

// WritePayload exports the values of the message specific payload attributes to output stream
func (msg *HandShakeRequestMessage) WritePayload(os types.TGOutputStream) types.TGError {
	startPos := os.GetPosition()
	logger.Log(fmt.Sprintf("Entering HandShakeRequestMessage:WritePayload at output buffer position: '%d'", startPos))
	os.(*iostream.ProtocolDataOutputStream).WriteByte(msg.GetRequestType())
	os.(*iostream.ProtocolDataOutputStream).WriteBoolean(msg.GetSslMode())
	os.(*iostream.ProtocolDataOutputStream).WriteInt(msg.GetChallenge())
	currPos := os.GetPosition()
	length := currPos - startPos
	logger.Log(fmt.Sprintf("Returning HandShakeRequestMessage::WritePayload at output buffer position at: %d after writing %d payload bytes", currPos, length))
	return nil
}

/////////////////////////////////////////////////////////////////
// Implement functions from Interface ==> encoding/BinaryMarshaler
/////////////////////////////////////////////////////////////////

func (msg *HandShakeRequestMessage) MarshalBinary() ([]byte, error) {
	// A simple encoding: plain text.
	var b bytes.Buffer
	_, err := fmt.Fprintln(&b, msg.BufLength, msg.VerbId, msg.SequenceNo, msg.Timestamp,
		msg.RequestId, msg.DataOffset, msg.AuthToken, msg.SessionId, msg.IsUpdatable, msg.sslMode, msg.challenge, msg.handshakeType)
	if err != nil {
		logger.Error(fmt.Sprintf("ERROR: Returning HandShakeRequestMessage:MarshalBinary w/ Error: '%+v'", err.Error()))
		return nil, err
	}
	return b.Bytes(), nil
}

/////////////////////////////////////////////////////////////////
// Implement functions from Interface ==> encoding/BinaryUnmarshaler
/////////////////////////////////////////////////////////////////

// UnmarshalBinary modifies the receiver so it must take a pointer receiver.
func (msg *HandShakeRequestMessage) UnmarshalBinary(data []byte) error {
	// A simple encoding: plain text.
	b := bytes.NewBuffer(data)
	_, err := fmt.Fscanln(b, &msg.BufLength, &msg.VerbId, &msg.SequenceNo,
		&msg.Timestamp, &msg.RequestId, &msg.DataOffset, &msg.AuthToken, &msg.SessionId, &msg.IsUpdatable,
		&msg.sslMode, &msg.challenge, &msg.handshakeType)
	if err != nil {
		logger.Error(fmt.Sprintf("ERROR: Returning HandShakeRequestMessage:UnmarshalBinary w/ Error: '%+v'", err.Error()))
		return err
	}
	return nil
}