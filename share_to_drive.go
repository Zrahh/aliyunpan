// Copyright (c) 2020 tickstep.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package aliyunpan

import (
	"fmt"
	"github.com/tickstep/aliyunpan-api/aliyunpan/apierror"
	"github.com/tickstep/aliyunpan-api/aliyunpan/apiutil"
	"github.com/tickstep/library-go/logger"
	"strings"
)

type(
	ShareToDriveParam struct {
	    
	    AutoRename      bool `json:"auto_rename"`
		// 源文件ID
		FileId        string `json:"file_id"`
		// share id
		ShareId        string `json:"share_id"`

		// 目标网盘ID
		ToDriveId     string `json:"to_drive_id"`
		// 目标文件夹ID
		ToParentFileId string `json:"to_parent_file_id"`
	}

	ShareToDriveResult struct {
		// 文件ID
		FileId string
		// 是否成功
		Success bool
	}
)

// ShareToDrive 转存文件
func (p *PanClient) ShareToDrive(param []*ShareToDriveParam) ([]*ShareToDriveResult, *apierror.ApiError) {
	// url
	fullUrl := &strings.Builder{}
	fmt.Fprintf(fullUrl, "%s/v3/batch", API_URL)
	logger.Verboseln("do request url: " + fullUrl.String())

	// data
	requests,e := p.getShareToDriveBatchRequestList(param)
	if e != nil {
		return nil, e
	}
	batchParam := BatchRequestParam{
		Requests: requests,
		Resource: "file",
	}

	// request
	result,err := p.BatchTask(fullUrl.String(), &batchParam)
	if err != nil {
		logger.Verboseln("file move error ", err)
		return nil, apierror.NewFailedApiError(err.Error())
	}

	// parse result
	r := []*ShareToDriveResult{}
	for _,item := range result.Responses{
		r = append(r, &ShareToDriveResult{
			FileId: item.Id,
			Success:     item.Status == 200,
		})
	}
	return r, nil
}

func (p *PanClient) getShareToDriveBatchRequestList(param []*ShareToDriveParam) (BatchRequestList, *apierror.ApiError) {
	if param == nil {
		return nil, apierror.NewFailedApiError("参数不能为空")
	}

	r := BatchRequestList{}
	for _,item := range param {
		r = append(r, &BatchRequest{
			Id:      item.FileId,
			Method:  "POST",
			Url:     "/file/copy",
			Headers: map[string]string{
				"Content-Type": "application/json",
			},
			Body:    apiutil.GetMapSet(item),
		})
	}
	return r, nil
}
