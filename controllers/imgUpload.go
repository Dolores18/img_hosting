package controllers

import (
	"fmt"
	"img_hosting/config"
	"img_hosting/dao"
	"img_hosting/middleware"
	"img_hosting/models"
	"img_hosting/pkg/logger"
	"img_hosting/services"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func Uploads(c *gin.Context) {
	db := models.GetDB()
	log := logger.GetLogger() //必须实例化
	config.LoadConfig()

	//获取用户信息
	claims, err := middleware.ParseAndValidateToken(c)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	userid := claims.UserID
	// Multipart form 。获取一张或者多张图片
	form, err := c.MultipartForm()
	if err != nil {
		c.JSON(400, gin.H{"error": "图片上传失败，请重新上传"})
		return
	}
	if form == nil {
		c.JSON(400, gin.H{"error": "图片上传失败，请重新上传"})
		return
	}
	files := form.File["upload[]"]
	println(files)

	for _, file := range files {
		println(file.Filename)
		//这里加一些对文件格式以及文件大小的检查
		valid, message, img_extension, imgsize := services.CheckImg(file.Filename, file.Size)
		println("文件名为：%s", file.Filename)
		if !valid {
			c.JSON(400, gin.H{"error": message})

			return
		}
		// 确保文件名安全
		name, extension, err := services.SanitizeFileName(file.Filename)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		image_name := name
		img_extension = extension
		// 读取文件内容
		fileContent, err := file.Open()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		defer fileContent.Close()

		// 将文件内容读取为字节切片
		fileBytes, err := io.ReadAll(fileContent)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		hash_img := services.HashFileName(fileBytes)
		hashedFilename := hash_img + img_extension

		if c.SaveUploadedFile(file, "./statics/uploads/"+hashedFilename) != nil {
			c.JSON(400, gin.H{"error": "文件保存失败"})
			return
		}

		imageUrl := config.AppConfigInstance.Url.Imgurl + hash_img + img_extension
		fmt.Println(imageUrl)
		//对文件进行哈希处理，文件名改成哈希文件名

		//保存文件到数据库
		dao.CreateImg(db, userid, imageUrl, image_name, img_extension, hash_img, imgsize, img_extension)
		c.JSON(200, gin.H{"msg": "图片成功上传", "data": imageUrl})

	}

	log.WithFields(logrus.Fields{"Name": userid}).Info("用户上传图片成功")

}

/*
curl -X POST ^
http://127.0.0.1:8080/upload ^
-F "file=@E:\down2\ginStudy\uploads\p1833285.jpg" ^
-H "Content-Type: multipart/form-data"
curl -X POST ^
http://127.0.0.1:8080/sigin ^
-H "Content-Type: application/json" ^
-d "{\"name\": \"mim23\", \"age\": 20, \"psd\": \"5d65466678@\"}"


curl -X POST ^
http://127.0.0.1:8080/imgupload ^
-H "Content-Type: application/json" ^
-d "{\"name\": \"mim23\", \"age\": 20, \"psd\": \"5d65466678@\"}"

"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VybmFtZSI6Im1pbWkyMyIsImV4cCI6MTcxOTM1NTU1NX0.oesiJ5wkaSKQtjAuP3vJzK-EYUMdfbKEFL0hWK6HOSg"


curl -X POST ^
http://127.0.0.1:8080/imgupload ^
-H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VybmFtZSI6Im1pbTIzIiwiZXhwIjoxNzE5NDUwMDQzfQ.6m6t_gBb0nx4JkzOvNAUtrNY_-Dt0fzA7E3UL9WOvTY" ^
-H "Content-Type: multipart/form-data" ^
-F "upload[]=@E:\\down2\\ginStudy\\uploads\\12.jpg"
-F "upload[]=@E:\\down2\\ginStudy\\uploads\\test.jpg"

curl上传的表名要一致




*/
