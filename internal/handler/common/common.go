package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	res "github.com/lin-snow/ech0/internal/handler/response"
	commonModel "github.com/lin-snow/ech0/internal/model/common"
	service "github.com/lin-snow/ech0/internal/service/common"
	errorUtil "github.com/lin-snow/ech0/internal/util/err"
)

type CommonHandler struct {
	commonService service.CommonServiceInterface
}

func NewCommonHandler(commonService service.CommonServiceInterface) *CommonHandler {
	return &CommonHandler{
		commonService: commonService,
	}
}

// UploadImage 上传图片
func (commonHandler *CommonHandler) UploadImage() gin.HandlerFunc {
	return res.Execute(func(ctx *gin.Context) res.Response {
		// 提取上传的 File数据
		file, err := ctx.FormFile("file")
		if err != nil {
			return res.Response{
				Msg: commonModel.INVALID_REQUEST_BODY,
				Err: err,
			}
		}

		// 提取userid
		userId := ctx.MustGet("userid").(uint)

		// 调用 CommonService 上传文件
		imageUrl, err := commonHandler.commonService.UploadImage(userId, file)
		if err != nil {
			return res.Response{
				Msg: commonModel.INVALID_REQUEST_BODY,
				Err: err,
			}
		}

		return res.Response{
			Data: imageUrl,
			Msg:  commonModel.UPLOAD_SUCCESS,
		}
	})
}

// DeleteImage 删除图片
func (commonHandler *CommonHandler) DeleteImage() gin.HandlerFunc {
	return res.Execute(func(ctx *gin.Context) res.Response {
		userId := ctx.MustGet("userid").(uint)

		var imageDto commonModel.ImageDto
		if err := ctx.ShouldBindJSON(&imageDto); err != nil {
			return res.Response{
				Msg: commonModel.INVALID_REQUEST_BODY,
				Err: err,
			}
		}

		if err := commonHandler.commonService.DeleteImage(userId, imageDto.URL, imageDto.SOURCE); err != nil {
			ctx.JSON(http.StatusOK, commonModel.Fail[string](errorUtil.HandleError(&commonModel.ServerError{
				Msg: "",
				Err: err,
			})))
			return res.Response{
				Msg: "",
				Err: err,
			}
		}

		return res.Response{
			Msg: commonModel.DELETE_SUCCESS,
		}
	})
}

// GetStatus 获取Echo状态
func (commonHandler *CommonHandler) GetStatus() gin.HandlerFunc {
	return res.Execute(func(ctx *gin.Context) res.Response {
		_, err := commonHandler.commonService.GetSysAdmin()
		if err != nil {
			return res.Response{
				Code: commonModel.InitInstallCode,
				Msg:  commonModel.SIGNUP_FIRST,
			}
		}

		status, err := commonHandler.commonService.GetStatus()
		if err != nil {
			return res.Response{
				Msg: "",
				Err: err,
			}
		}

		return res.Response{
			Data: status,
			Msg:  commonModel.GET_STATUS_SUCCESS,
		}
	})

}

// GetHeatMap 获取热力图数据
func (commonHandler *CommonHandler) GetHeatMap() gin.HandlerFunc {
	return res.Execute(func(ctx *gin.Context) res.Response {
		// 调用 Service 层获取热力图数据
		heatMap, err := commonHandler.commonService.GetHeatMap()
		if err != nil {
			return res.Response{
				Msg: "",
				Err: err,
			}
		}

		return res.Response{
			Data: heatMap,
			Msg:  commonModel.GET_HEATMAP_SUCCESS,
		}
	})

}

// GetRss 获取RSS
func (commonHandler *CommonHandler) GetRss(ctx *gin.Context) {
	atom, err := commonHandler.commonService.GenerateRSS(ctx)
	if err != nil {
		ctx.JSON(http.StatusOK, commonModel.Fail[string](errorUtil.HandleError(&commonModel.ServerError{
			Msg: "",
			Err: err,
		})))
		return
	}

	ctx.Data(http.StatusOK, "application/rss+xml; charset=utf-8", []byte(atom))
}

// UploadAudio 上传音频
func (commonHandler *CommonHandler) UploadAudio() gin.HandlerFunc {
	return res.Execute(func(ctx *gin.Context) res.Response {
		// 提取userid
		userId := ctx.MustGet("userid").(uint)

		// 提取上传的 File数据
		file, err := ctx.FormFile("file")
		if err != nil {
			return res.Response{
				Msg: commonModel.INVALID_REQUEST_BODY,
				Err: err,
			}
		}

		audioUrl, err := commonHandler.commonService.UploadMusic(userId, file)
		if err != nil {
			return res.Response{
				Msg: "",
				Err: err,
			}
		}

		return res.Response{
			Data: audioUrl,
			Msg:  commonModel.UPLOAD_SUCCESS,
		}
	})
}

// DeleteAudio 删除音频
func (commonHandler *CommonHandler) DeleteAudio() gin.HandlerFunc {
	return res.Execute(func(ctx *gin.Context) res.Response {
		// 提取userid
		userId := ctx.MustGet("userid").(uint)

		if err := commonHandler.commonService.DeleteMusic(userId); err != nil {
			return res.Response{
				Msg: "",
				Err: err,
			}
		}

		return res.Response{
			Msg: commonModel.DELETE_SUCCESS,
		}
	})

}

// GetPlayMusic 获取可播放的音乐
func (commonHandler *CommonHandler) GetPlayMusic() gin.HandlerFunc {
	return res.Execute(func(ctx *gin.Context) res.Response {
		musicUrl := commonHandler.commonService.GetPlayMusicUrl()

		return res.Response{
			Data: musicUrl,
			Msg:  commonModel.GET_MUSIC_URL_SUCCESS,
		}
	})
}

// PlayMusic 播放音乐
func (commonHandler *CommonHandler) PlayMusic(ctx *gin.Context) {
	commonHandler.commonService.PlayMusic(ctx)
}
