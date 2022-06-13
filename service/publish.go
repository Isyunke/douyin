package service

import (
	"bytes"
	"errors"
	"io"
	"mime/multipart"
	"os"
	"os/exec"
	"path"
	"runtime"
	"strconv"

	"github.com/warthecatalyst/douyin/api"
	"github.com/warthecatalyst/douyin/config"
	"github.com/warthecatalyst/douyin/dao"
	"github.com/warthecatalyst/douyin/filex"
	"github.com/warthecatalyst/douyin/idgenerator"
	"github.com/warthecatalyst/douyin/logx"
	"github.com/warthecatalyst/douyin/model"
	"github.com/warthecatalyst/douyin/oss"
)

func getUploadPath(userId int64, fileName string) string {
	return config.OssConf.BucketDirectory + "/" + strconv.FormatInt(userId, 10) + "/" + fileName
}

// getUploadURL 得到一名用户对应的云端存储路径
func getUploadURL(userId int64, fileName string) string {
	return "https://" + config.OssConf.Bucket + "." + config.OssConf.Url + "/" + getUploadPath(userId, fileName)
}

// checkFileExt 检查文件的拓展名
func checkFileExt(fileName string) bool {
	//检查文件的扩展名
	ext := path.Ext(fileName)
	ext = string(bytes.ToLower([]byte(ext)))
	for _, legalExt := range config.VideoConf.AllowedExts {
		if ext == legalExt {
			return true
		}
	}
	return false
}

// checkFileMaxSize 检查文件的大小，单位为B
func checkFileMaxSize(videoSize int64) bool {
	return videoSize <= config.VideoConf.UploadMaxSize*filex.MB
}

// extractCoverFromVideo 从视频中截取图像的第一帧
func extractCoverFromVideo(pathVideo, pathImg string) error {
	binPath := "./bin/" // 最后带斜杠
	if runtime.GOOS == "windows" {
		binPath += "windows/"
	} else if runtime.GOOS == "darwin" {
		binPath += "darwin/"
	} else {
		binPath += "Linux/"
	}
	frameExtractionTime := "0"
	image_mode := "image2"
	vtime := "0.001"

	// create the command
	cmd := exec.Command(binPath+"ffmpeg",
		"-i", pathVideo,
		"-y",
		"-f", image_mode,
		"-ss", frameExtractionTime,
		"-t", vtime,
		"-y", pathImg)

	// run the command and don't wait for it to finish. waiting exec is run
	// fmt.Println(cmd.String())
	err := cmd.Start()
	if err != nil {
		return err
	}
	err = cmd.Wait()
	if err != nil {
		return err
	}
	return nil
}

func PublishVideoInfo(data *multipart.FileHeader, userId int64, title string) error {
	return newPublishVideoInfoFlow(data, userId, title).Do()
}

type publishVideoInfoFlow struct {
	data   *multipart.FileHeader
	userId int64
	title  string
}

func newPublishVideoInfoFlow(data *multipart.FileHeader, userId int64, title string) *publishVideoInfoFlow {
	return &publishVideoInfoFlow{
		data:   data,
		userId: userId,
		title:  title,
	}
}

func (p *publishVideoInfoFlow) Do() error {
	fileName := p.data.Filename
	//首先检查video扩展名和大小
	if !checkFileExt(fileName) {
		logx.DyLogger.Error("wrong input video type")
		return errors.New("视频格式不符合要求")
	}
	if !checkFileMaxSize(p.data.Size) {
		logx.DyLogger.Error("file extends the maximum value")
		return errors.New("视频过大")
	}

	//然后把文件保存至本地
	err := p.saveFile(config.VideoConf.SavePath)
	if err != nil {
		logx.DyLogger.Error("Saving goes wrong")
		return errors.New("save Error")
	}
	//截取视频的第一帧作为cover
	saveDir := path.Join(config.VideoConf.SavePath, strconv.FormatInt(p.userId, 10))
	saveVideo := saveDir + "/" + p.data.Filename
	coverName := filex.GetFileNameWithOutExt(fileName) + "_cover" + ".jpeg"
	saveCover := saveDir + "/" + coverName
	err = extractCoverFromVideo(saveVideo, saveCover)
	if err != nil {
		logx.DyLogger.Error("Saving goes wrong")
		return errors.New("save Error")
	}

	//上传视频和封面
	logx.DyLogger.Info("Saving Complete, in Upload")
	err = p.uploadServer()
	if err != nil {
		logx.DyLogger.Error("An error occurs in upload the file to Aliyun")
		return errors.New("upload Error")
	}
	err = p.uploadCoverToServer(saveCover, coverName)
	if err != nil {
		logx.DyLogger.Error("An error occurs in upload the file to Aliyun")
		return errors.New("upload Error")
	}

	//调用dao层函数操作数据库添加对应的Video字段
	video := model.Video{
		VideoID:       idgenerator.GenerateVid(), //随机生成VideoId,之后可以进行调整
		VideoName:     p.title,
		UserID:        p.userId,
		FavoriteCount: 0,
		CommentCount:  0,
		PlayURL:       getUploadURL(p.userId, p.data.Filename),
		CoverURL:      getUploadURL(p.userId, coverName),
	}

	err = dao.NewPublishDaoInstance().AddVideo(&video)
	if err != nil {
		logx.DyLogger.Error("An error occurs in ")
		return errors.New("database Error")
	}
	return nil
}

func (p *publishVideoInfoFlow) saveFile(savePath string) error {
	userSavePath := path.Join(savePath, strconv.FormatInt(p.userId, 10))
	if flag, _ := filex.PathExists(userSavePath); !flag {
		err := os.Mkdir(userSavePath, os.ModePerm)
		if err != nil {
			logx.DyLogger.Errorf("Error in making directory: %s", err)
			return err
		}
	}
	src, err := p.data.Open()
	if err != nil {
		logx.DyLogger.Error("Error in Saving", err)
		return err
	}
	defer src.Close()

	out, err := os.Create(userSavePath + "/" + p.data.Filename)
	if err != nil {
		logx.DyLogger.Error("Error in Saving", err)
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, src)
	return err
}

func (p *publishVideoInfoFlow) uploadServer() error {
	src, err := p.data.Open()
	if err != nil {
		logx.DyLogger.Error("Error in upload3:", err)
		return err
	}
	defer src.Close()

	// 先将文件流上传至BucketDirectory目录下
	err = oss.UploadFromReader(getUploadPath(p.userId, p.data.Filename), src)
	if err != nil {
		logx.DyLogger.Error("Error in upload4:", err)
		return err
	}

	logx.DyLogger.Info("file upload success")
	return nil
}

// uploadCoverToServer 把封面上传到云端
func (p *publishVideoInfoFlow) uploadCoverToServer(filePath, fileName string) error {
	if err := oss.UploadFromFile(getUploadPath(p.userId, fileName), filePath); err != nil {
		logx.DyLogger.Error(err)
		return err
	}

	return nil
}

func PublishListInfo(userId int64) (*VideoList, error) {
	return newPublishListInfoFlow(userId).Do()
}

type publishListInfoFlow struct {
	userId int64
}

func newPublishListInfoFlow(userId int64) *publishListInfoFlow {
	return &publishListInfoFlow{userId: userId}
}

func (p *publishListInfoFlow) Do() (*VideoList, error) {
	videoIds, err := dao.NewPublishDaoInstance().GetVideoPublistList(p.userId)
	if err != nil {
		return nil, err
	}
	var videolist VideoList
	for _, videoId := range videoIds {
		user, err := videoService.getUserFromVideoId(videoId)
		if err != nil {
			return nil, err
		}
		video, err := dao.NewVideoDaoInstance().GetVideoFromId(videoId)

		videolist = append(videolist, api.Video{
			Id:            videoId,
			Author:        *user,
			PlayUrl:       video.PlayURL,
			CoverUrl:      video.CoverURL,
			FavoriteCount: int64(video.FavoriteCount),
			CommentCount:  int64(video.CommentCount),
			IsFavorite:    true, //被videolist返回的肯定点赞过了
		})
	}
	return &videolist, nil
}

type PublishService struct{}

var publishService = &PublishService{}

// 构造 Video 切片
func newVideoList(videos []model.Video) []api.Video {
	var v []api.Video
	for _, i := range videos {
		v = append(v, api.Video{
			Id: i.VideoID,
			Author: api.User{
				Id: i.UserID,
			},
			PlayUrl:       i.PlayURL,
			CoverUrl:      i.CoverURL,
			CommentCount:  int64(i.CommentCount),
			FavoriteCount: int64(i.FavoriteCount),
		})
	}

	return v
}
