package exception

import (
	"fmt"
	"github.com/TIBCOSoftware/tgdb-client/client/goAPI/types"
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
 * File name: TGErrorIOException.go
 * Created on: Oct 20, 2018
 * Created by: achavan
 * SVN id: $id: $
 *
 */

type IOException struct {
	*types.TGDBError
}

// Create New IOException Instance
func DefaultTGIOException() *IOException {
	newException := IOException{
		TGDBError: types.DefaultTGDBError(),
	}
	newException.ErrorType = types.TGErrorIOException
	return &newException
}

func NewTGIOException(eCode string, eType int, eMsg, eDetails string) *IOException {
	newException := DefaultTGIOException()
	newException.ErrorCode = eCode
	newException.ErrorType = eType
	newException.ErrorMsg = eMsg
	newException.ErrorDetails = eDetails
	return newException
}

func NewTGIOExceptionWithMsg(msg string) *IOException {
	newException := DefaultTGIOException()
	newException.ErrorMsg = msg
	return newException
}

/////////////////////////////////////////////////////////////////
// Implement functions from Interface ==> TGError
/////////////////////////////////////////////////////////////////

func (e *IOException) GetErrorCode() string {
	return e.ErrorCode
}

func (e *IOException) GetErrorType() int {
	return e.ErrorType
}

func (e *IOException) GetErrorMsg() string {
	return e.ErrorMsg
}

func (e *IOException) GetErrorDetails() string {
	return e.ErrorDetails
}

/////////////////////////////////////////////////////////////////
// Implement functions from Interface ==> error
/////////////////////////////////////////////////////////////////

func (e *IOException) Error() string {
	errMsg := fmt.Sprintf("ErrorCode: %s, ErrorType: %d, ErrorMessage: %s, ErrorDetails: %s", e.ErrorCode, e.ErrorType, e.ErrorMsg, e.ErrorDetails)
	return errMsg
}
