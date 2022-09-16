package api

import (
	"github.com/gin-gonic/gin"
	"path"
	"taogin/config/atom"
	"taogin/config/global"
	"taogin/core/response"
)

type Upload struct {
}

func NewUpload() *Upload {
	return &Upload{}
}

func (this *Upload) Upload(ctx *gin.Context) {
	file, err := ctx.FormFile("file")
	if err != nil {
		global.LOG.Error(err)
		response.Error(ctx, atom.ERROR_CODE_UPLOAD_FAILE, atom.GetMsgByCode(atom.ERROR_CODE_UPLOAD_FAILE))
	}
	path := path.Join("storage/upload", file.Filename)
	//上传本地
	err = ctx.SaveUploadedFile(file, path)
	if err != nil {
		global.LOG.Error(err)
		response.Error(ctx, atom.ERROR_CODE_UPLOAD_FAILE, atom.GetMsgByCode(atom.ERROR_CODE_UPLOAD_FAILE))
	}
}
