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
 * File name: VerbQueryResponse.go
 * Created on: Sep 30, 2018
 * Created by: achavan
 * SVN id: $id: $
 *
 */

type QueryResponseMessage struct {
	*AbstractProtocolMessage
	entityStream types.TGInputStream
	hasResult    bool
	totalCount   int
	resultCount  int
	result       int
	queryHashId  int64
}

func DefaultQueryResponseMessage() *QueryResponseMessage {
	// We must register the concrete type for the encoder and decoder (which would
	// normally be on a separate machine from the encoder). On each end, this tells the
	// engine which concrete type is being sent that implements the interface.
	gob.Register(QueryResponseMessage{})

	newMsg := QueryResponseMessage{
		AbstractProtocolMessage: DefaultAbstractProtocolMessage(),
	}
	newMsg.hasResult = false
	newMsg.totalCount = 0
	newMsg.resultCount = 0
	newMsg.result = 0
	newMsg.queryHashId = 0
	newMsg.verbId = VerbQueryResponse
	newMsg.BufLength = int(reflect.TypeOf(newMsg).Size())
	return &newMsg
}

// Create New Message Instance
func NewQueryResponseMessage(authToken, sessionId int64) *QueryResponseMessage {
	newMsg := DefaultQueryResponseMessage()
	newMsg.authToken = authToken
	newMsg.sessionId = sessionId
	newMsg.BufLength = int(reflect.TypeOf(*newMsg).Size())
	return newMsg
}

/////////////////////////////////////////////////////////////////
// Helper functions for QueryResponseMessage
/////////////////////////////////////////////////////////////////

func (msg *QueryResponseMessage) GetResultCount() int {
	return msg.resultCount
}

func (msg *QueryResponseMessage) GetResult() int {
	return msg.result
}

func (msg *QueryResponseMessage) GetQueryHashId() int64 {
	return msg.queryHashId
}

func (msg *QueryResponseMessage) GetEntityStream() types.TGInputStream {
	return msg.entityStream
}

func (msg *QueryResponseMessage) GetHasResult() bool {
	return msg.hasResult
}

func (msg *QueryResponseMessage) GetTotalCount() int {
	return msg.totalCount
}

func (msg *QueryResponseMessage) SetResultCount(size int) {
	msg.resultCount = size
}

func (msg *QueryResponseMessage) SetResult(size int) {
	msg.result = size
}

func (msg *QueryResponseMessage) SetQueryHashId(id int64) {
	msg.queryHashId = id
}

func (msg *QueryResponseMessage) SetEntityStream(eStream types.TGInputStream) {
	msg.entityStream = eStream
}

func (msg *QueryResponseMessage) SetHasResult(resultFlag bool) {
	msg.hasResult = resultFlag
}

func (msg *QueryResponseMessage) SetTotalCount(count int) {
	msg.totalCount = count
}

/////////////////////////////////////////////////////////////////
// Implement functions from Interface ==> TGMessage
/////////////////////////////////////////////////////////////////

// FromBytes constructs a message object from the input buffer in the byte format
func (msg *QueryResponseMessage) FromBytes(buffer []byte) (types.TGMessage, types.TGError) {
	logger.Log(fmt.Sprint("Entering QueryResponseMessage:FromBytes"))
	if len(buffer) < 0 {
		logger.Error(fmt.Sprint("ERROR: Returning QueryResponseMessage:FromBytes w/ Error: Invalid Message Buffer"))
		return nil, exception.CreateExceptionByType(types.TGErrorInvalidMessageLength)
	}

	is := iostream.NewProtocolDataInputStream(buffer)

	// First member attribute / element of message header is BufLength
	bufLen, err := is.ReadInt()
	if err != nil {
		logger.Error(fmt.Sprint("ERROR: Returning QueryResponseMessage:FromBytes w/ Error in reading buffer length from message buffer"))
		return nil, err
	}
	logger.Log(fmt.Sprintf("Inside QueryResponseMessage:FromBytes read bufLen as '%+v'", bufLen))
	if bufLen != len(buffer) {
		errMsg := fmt.Sprint("Buffer length mismatch")
		return nil, exception.GetErrorByType(types.TGErrorInvalidMessageLength, types.INTERNAL_SERVER_ERROR, errMsg, "")
	}

	logger.Log(fmt.Sprint("Inside QueryResponseMessage:FromBytes - about to APMReadHeader"))
	err = APMReadHeader(msg, is)
	if err != nil {
		errMsg := fmt.Sprintf("Unable to recreate message from '%+v' in byte format", buffer)
		return nil, exception.GetErrorByType(types.TGErrorIOException, types.INTERNAL_SERVER_ERROR, errMsg, "")
	}

	logger.Log(fmt.Sprint("Inside QueryResponseMessage:FromBytes - about to ReadPayload"))
	err = msg.ReadPayload(is)
	if err != nil {
		errMsg := fmt.Sprintf("Unable to recreate message from '%+v' in byte format", buffer)
		return nil, exception.GetErrorByType(types.TGErrorIOException, types.INTERNAL_SERVER_ERROR, errMsg, "")
	}

	logger.Log(fmt.Sprintf("QueryResponseMessage::FromBytes resulted in '%+v'", msg))
	return msg, nil
}

// ToBytes converts a message object into byte format to be sent over the network to TGDB server
func (msg *QueryResponseMessage) ToBytes() ([]byte, int, types.TGError) {
	logger.Log(fmt.Sprint("Entering QueryResponseMessage:ToBytes"))
	os := iostream.DefaultProtocolDataOutputStream()

	logger.Log(fmt.Sprint("Inside QueryResponseMessage:ToBytes - about to APMWriteHeader"))
	err := APMWriteHeader(msg, os)
	if err != nil {
		errMsg := fmt.Sprintf("Unable to export message '%+v' in byte format", msg)
		return nil, -1, exception.GetErrorByType(types.TGErrorIOException, types.INTERNAL_SERVER_ERROR, errMsg, "")
	}

	logger.Log(fmt.Sprint("Inside QueryResponseMessage:ToBytes - about to WritePayload"))
	err = msg.WritePayload(os)
	if err != nil {
		errMsg := fmt.Sprintf("Unable to export message '%+v' in byte format", msg)
		return nil, -1, exception.GetErrorByType(types.TGErrorIOException, types.INTERNAL_SERVER_ERROR, errMsg, "")
	}

	_, err = os.WriteIntAt(0, os.GetLength())
	if err != nil {
		logger.Error(fmt.Sprint("ERROR: Returning QueryResponseMessage:ToBytes w/ Error in writing buffer length"))
		return nil, -1, err
	}
	logger.Log(fmt.Sprintf("QueryResponseMessage::ToBytes results bytes-on-the-wire in '%+v'", os.GetBuffer()))
	return os.GetBuffer(), os.GetLength(), nil
}

// GetAuthToken gets the authToken
func (msg *QueryResponseMessage) GetAuthToken() int64 {
	return msg.HeaderGetAuthToken()
}

// GetIsUpdatable checks whether this message updatable or not
func (msg *QueryResponseMessage) GetIsUpdatable() bool {
	return msg.GetUpdatableFlag()
}

// GetMessageByteBufLength gets the MessageByteBufLength. This method is called after the toBytes() is executed.
func (msg *QueryResponseMessage) GetMessageByteBufLength() int {
	return msg.HeaderGetMessageByteBufLength()
}

// GetRequestId gets the requestId for the message. This will be used as the CorrelationId
func (msg *QueryResponseMessage) GetRequestId() int64 {
	return msg.HeaderGetRequestId()
}

// GetSequenceNo gets the sequenceNo of the message
func (msg *QueryResponseMessage) GetSequenceNo() int64 {
	return msg.HeaderGetSequenceNo()
}

// GetSessionId gets the session id
func (msg *QueryResponseMessage) GetSessionId() int64 {
	return msg.HeaderGetSessionId()
}

// GetTimestamp gets the Timestamp
func (msg *QueryResponseMessage) GetTimestamp() int64 {
	return msg.HeaderGetTimestamp()
}

// GetVerbId gets verbId of the message
func (msg *QueryResponseMessage) GetVerbId() int {
	return msg.HeaderGetVerbId()
}

// SetAuthToken sets the authToken
func (msg *QueryResponseMessage) SetAuthToken(authToken int64) {
	msg.HeaderSetAuthToken(authToken)
}

// SetDataOffset sets the offset at which data starts in the payload
func (msg *QueryResponseMessage) SetDataOffset(dataOffset int16) {
	msg.HeaderSetDataOffset(dataOffset)
}

// SetIsUpdatable sets the updatable flag
func (msg *QueryResponseMessage) SetIsUpdatable(updateFlag bool) {
	msg.SetUpdatableFlag(updateFlag)
}

// SetMessageByteBufLength sets the message buffer length
func (msg *QueryResponseMessage) SetMessageByteBufLength(bufLength int) {
	msg.HeaderSetMessageByteBufLength(bufLength)
}

// SetRequestId sets the request id
func (msg *QueryResponseMessage) SetRequestId(requestId int64) {
	msg.HeaderSetRequestId(requestId)
}

// SetSequenceNo sets the sequenceNo
func (msg *QueryResponseMessage) SetSequenceNo(sequenceNo int64) {
	msg.HeaderSetSequenceNo(sequenceNo)
}

// SetSessionId sets the session id
func (msg *QueryResponseMessage) SetSessionId(sessionId int64) {
	msg.HeaderSetSessionId(sessionId)
}

// SetTimestamp sets the timestamp
func (msg *QueryResponseMessage) SetTimestamp(timestamp int64) types.TGError {
	if !(msg.isUpdatable || timestamp != -1) {
		logger.Error(fmt.Sprint("ERROR: Returning APMReadHeader:setTimestamp as !msg.IsUpdatable && timestamp != -1"))
		errMsg := fmt.Sprintf("Mutating a readonly message '%s'", GetVerb(msg.verbId).name)
		return exception.GetErrorByType(types.TGErrorGeneralException, types.INTERNAL_SERVER_ERROR, errMsg, "")
	}
	msg.HeaderSetTimestamp(timestamp)
	return nil
}

// SetVerbId sets verbId of the message
func (msg *QueryResponseMessage) SetVerbId(verbId int) {
	msg.HeaderSetVerbId(verbId)
}

func (msg *QueryResponseMessage) String() string {
	var buffer bytes.Buffer
	buffer.WriteString("QueryResponseMessage:{")
	buffer.WriteString(fmt.Sprintf("HasResult: %+v", msg.hasResult))
	buffer.WriteString(fmt.Sprintf(", TotalCount: %d", msg.totalCount))
	buffer.WriteString(fmt.Sprintf(", ResultCount: %d", msg.resultCount))
	buffer.WriteString(fmt.Sprintf(", Result: %d", msg.result))
	buffer.WriteString(fmt.Sprintf(", QueryHashId: %d", msg.queryHashId))
	//buffer.WriteString(fmt.Sprintf(", EntityStream: %+v", msg.entityStream))
	buffer.WriteString(fmt.Sprintf(", BufLength: %d", msg.BufLength))
	strArray := []string{buffer.String(), msg.APMMessageToString()+"}"}
	msgStr := strings.Join(strArray, ", ")
	return  msgStr
}

// UpdateSequenceAndTimeStamp updates the SequenceAndTimeStamp, if message is mutable
// @param timestamp
// @return TGMessage on success, error on failure
func (msg *QueryResponseMessage) UpdateSequenceAndTimeStamp(timestamp int64) types.TGError {
	return msg.SetSequenceAndTimeStamp(timestamp)
}

// ReadHeader reads the bytes from input stream and constructs a common header of network packet
func (msg *QueryResponseMessage) ReadHeader(is types.TGInputStream) types.TGError {
	return APMReadHeader(msg, is)
}

// WriteHeader exports the values of the common message header attributes to output stream
func (msg *QueryResponseMessage) WriteHeader(os types.TGOutputStream) types.TGError {
	return APMWriteHeader(msg, os)
}

// ReadPayload reads the bytes from input stream and constructs message specific payload attributes
func (msg *QueryResponseMessage) ReadPayload(is types.TGInputStream) types.TGError {
	logger.Log(fmt.Sprint("Entering QueryResponseMessage:ReadPayload"))
	avail, err := is.(*iostream.ProtocolDataInputStream).Available()
	if err != nil {
		logger.Error(fmt.Sprint("ERROR: Returning QueryResponseMessage:ReadPayload w/ Error in reading available bytes from message buffer"))
		return err
	}
	logger.Log(fmt.Sprintf("QueryResponseMessage:ReadPayload read avail as '%+v'", avail))
	if avail == 0 {
		logger.Log(fmt.Sprintf("Query response has no data"))
		errMsg := fmt.Sprint("Query response has no data")
		return exception.GetErrorByType(types.TGErrorProtocolNotSupported, types.INTERNAL_SERVER_ERROR, errMsg, "")
	}
	msg.entityStream = is

	_, err = is.(*iostream.ProtocolDataInputStream).ReadInt() // buf length
	if err != nil {
		logger.Error(fmt.Sprint("ERROR: Returning QueryResponseMessage:ReadPayload w/ Error in reading buffer length from message buffer"))
		return err
	}

	_, err = is.(*iostream.ProtocolDataInputStream).ReadInt() // checksum
	if err != nil {
		logger.Error(fmt.Sprint("ERROR: Returning QueryResponseMessage:ReadPayload w/ Error in reading checksum from message buffer"))
		return err
	}

	result, err := is.(*iostream.ProtocolDataInputStream).ReadInt() // query result
	if err != nil {
		logger.Error(fmt.Sprint("ERROR: Returning QueryResponseMessage:ReadPayload w/ Error in reading query result from message buffer"))
		return err
	}
	logger.Log(fmt.Sprintf("QueryResponseMessage:ReadPayload read result as '%+v'", result))

	hashId, err := is.(*iostream.ProtocolDataInputStream).ReadLong() // query hash id
	if err != nil {
		logger.Error(fmt.Sprint("ERROR: Returning QueryResponseMessage:ReadPayload w/ Error in reading query hashId from message buffer"))
		return err
	}
	logger.Log(fmt.Sprintf("QueryResponseMessage:ReadPayload read hashId as '%+v'", hashId))

	syntax, err := is.(*iostream.ProtocolDataInputStream).ReadByte()
	if err != nil {
		logger.Error(fmt.Sprint("ERROR: Returning QueryResponseMessage:ReadPayload w/ Error in reading syntax from message buffer"))
		return err
	}
	logger.Log(fmt.Sprintf("QueryResponseMessage:ReadPayload read syntax as '%+v'", syntax))

	resultCount, err := is.(*iostream.ProtocolDataInputStream).ReadInt()
	if err != nil {
		logger.Error(fmt.Sprint("ERROR: Returning QueryResponseMessage:ReadPayload w/ Error in reading resultCount from message buffer"))
		return err
	}
	logger.Log(fmt.Sprintf("QueryResponseMessage:ReadPayload read resultCount as '%+v'", resultCount))
	if resultCount > 0 {
		msg.SetHasResult(true)
	}

	if syntax == 1 {
		totalCount, err := is.(*iostream.ProtocolDataInputStream).ReadInt()
		if err != nil {
			logger.Error(fmt.Sprint("ERROR: Returning QueryResponseMessage:ReadPayload w/ Error in reading totalCount from message buffer"))
			return err
		}
		logger.Log(fmt.Sprintf("QueryResponseMessage:ReadPayload read totalCount as '%+v'", totalCount))
		msg.SetTotalCount(totalCount)
		logger.Log(fmt.Sprintf("Query has '%d' result entities and %d total entities", resultCount, totalCount))
	} else {
		logger.Log(fmt.Sprintf("Query has '%d' result count", resultCount))
	}

	msg.SetResult(result)
	msg.SetQueryHashId(hashId)
	msg.SetResultCount(resultCount)
	logger.Log(fmt.Sprint("Returning QueryResponseMessage:ReadPayload"))
	return nil
}

// WritePayload exports the values of the message specific payload attributes to output stream
func (msg *QueryResponseMessage) WritePayload(os types.TGOutputStream) types.TGError {
	// No-Op for Now
	return nil
}

/////////////////////////////////////////////////////////////////
// Implement functions from Interface ==> encoding/BinaryMarshaller
/////////////////////////////////////////////////////////////////

func (msg *QueryResponseMessage) MarshalBinary() ([]byte, error) {
	// A simple encoding: plain text.
	var b bytes.Buffer
	_, err := fmt.Fprintln(&b, msg.BufLength, msg.verbId, msg.sequenceNo, msg.timestamp,
		msg.requestId, msg.dataOffset, msg.authToken, msg.sessionId, msg.isUpdatable, msg.hasResult, msg.totalCount,
		msg.resultCount, msg.result, msg.queryHashId)
	if err != nil {
		logger.Error(fmt.Sprintf("ERROR: Returning QueryResponseMessage:MarshalBinary w/ Error: '%+v'", err.Error()))
		return nil, err
	}
	return b.Bytes(), nil
}

/////////////////////////////////////////////////////////////////
// Implement functions from Interface ==> encoding/BinaryUnmarshaller
/////////////////////////////////////////////////////////////////

// UnmarshalBinary modifies the receiver so it must take a pointer receiver.
func (msg *QueryResponseMessage) UnmarshalBinary(data []byte) error {
	// A simple encoding: plain text.
	b := bytes.NewBuffer(data)
	_, err := fmt.Fscanln(b, &msg.BufLength, &msg.verbId, &msg.sequenceNo,
		&msg.timestamp, &msg.requestId, &msg.dataOffset, &msg.authToken, &msg.sessionId, &msg.isUpdatable,
		&msg.hasResult, &msg.totalCount, &msg.resultCount, &msg.result, &msg.queryHashId)
	if err != nil {
		logger.Error(fmt.Sprintf("ERROR: Returning QueryResponseMessage:UnmarshalBinary w/ Error: '%+v'", err.Error()))
		return err
	}
	return nil
}
