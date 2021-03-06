package com.tibco.tgdb.pdu.impl;

import com.tibco.tgdb.exception.TGException;
import com.tibco.tgdb.pdu.TGInputStream;
import com.tibco.tgdb.pdu.TGOutputStream;
import com.tibco.tgdb.pdu.VerbId;

import java.io.IOException;

/**
 * Copyright 2016 TIBCO Software Inc. All rights reserved.
 * 
 * Licensed under the Apache License, Version 2.0 (the "License"); You may not use this file except 
 * in compliance with the License.
 * A copy of the License is included in the distribution package with this file.
 * You also may obtain a copy of the License at http://www.apache.org/licenses/LICENSE-2.0
 * 
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 * <p/>
 * File name :MetadataRequest
 * Created on: 12/24/14
 * Created by: chung
 * <p/>
 * SVN Id: $Id: MetadataRequest.java 583 2016-03-15 02:02:39Z vchung $
 */
public class MetadataRequest extends AbstractProtocolMessage {

    @Override
    public VerbId getVerbId() {
        return VerbId.MetadataRequest;
    }

    @Override
    public boolean isUpdateable() {
        return true;
    }

	@Override
	protected void writePayload(TGOutputStream os) throws TGException,
			IOException {
		// TODO Auto-generated method stub
		
	}

	@Override
	protected void readPayload(TGInputStream is) throws TGException,
			IOException {
		// TODO Auto-generated method stub
		
	}
}
