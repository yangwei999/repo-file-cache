package controllers

import (
	"github.com/opensourceways/repo-file-cache/models"
)

type FileController struct {
	baseController
}

func (fc *FileController) Prepare() {
	fc.apiPrepare()
}

// @Title Set
// @Description store the file
// @Param	body		body 	models.FileUpdateOption	true		"body for storing file"
// @Success 201 {string} "successfully"
// @Failure 400 error_parsing_api_body:     parse payload of request failed
// @Failure 401 error_missing_input_param:  missing some input parameters
// @Failure 402 error_has_same_file:        there are same files which are located in same path
// @Failure 403 error_not_same_file:        the name of files to be stored is not same
// @Failure 404 error_invalid_file_path:    the file path is invalid
// @Failure 500 error_system_error:         system error
// @router / [post]
func (fc *FileController) Set() {
	action := "store the file"

	input := &models.FileUpdateOption{}
	if fr := fc.fetchInputPayload(input); fr != nil {
		fc.sendFailedResultAsResp(fr, action)
		return
	}

	if merr := input.Validate(); merr != nil {
		fc.sendModelErrorAsResp(merr, action)
		return
	}

	if merr := input.Update(); merr != nil {
		fc.sendModelErrorAsResp(merr, action)
		return
	}

	fc.sendResponse(action+" successfully", 0)
}

// @Title Get
// @Description Get the stored file
// @Success 200 {object} models.FilesInfo
// @Failure 400 missing_url_path_parameter: missing url path parameter
// @Failure 500 system_error:               system error
// @router /:platform/:org/:repo/:branch/:filename [get]
func (fc *FileController) Get() {
	action := "list files"

	b := models.Branch{
		Platform: fc.GetString(":platform"),
		Org:      fc.GetString(":org"),
		Repo:     fc.GetString(":repo"),
		Branch:   fc.GetString(":branch"),
	}

	summary, _ := fc.GetBool("summary")

	r, merr := models.GetFiles(b, fc.GetString(":filename"), summary)
	if merr != nil {
		fc.sendModelErrorAsResp(merr, action)
		return
	}

	fc.sendSuccessResp(r)
}
