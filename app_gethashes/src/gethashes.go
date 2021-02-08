package main

import (
	"errors"
	"io"
	"os"

        "crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/hex"
        "net/http"

        "github.com/gin-gonic/gin"
)

/*
 * Examples for using this with curl.
 * curl -XGET http://localhost:7001/file/md5 -F file=@main.go
 * curl -XGET http://localhost:7001/file/all -F file=@main.go
 * curl -XGET http://localhost:7001/string/md5?ct=<STRING TO GET HASH OF>
 * curl -XGET http://localhost:7001/string/all?ct=<STRING TO GET HASH OF>
 *
 * This will give you a nice output that is easy to read.
 * curl -XGET http://localhost:7001/file/all -F file=@main.go | python -m json.tool
 */

type FileResponse struct {
	Errorcode int
	Errormsg error
	Filename string
	Hashes []HashResponse
}

type StringResponse struct {
	Errorcode int
	Errormsg error
	Plaintext string
	Hashes []HashResponse
}

type HashResponse struct {
	Hash string
	Type string
}

func GetFileHash(c *gin.Context) FileResponse {
	var response FileResponse

	// Get file and save it to disk
	file, err := c.FormFile("file")
	if err != nil {
		response.Errorcode = 1
		response.Errormsg = err
		return response
	}
	response.Filename = file.Filename

	err = c.SaveUploadedFile(file, "/tmp/"+file.Filename)
	if err != nil {
		response.Errorcode = 1
		response.Errormsg = err
		return response
	}

	// Open file and defer closing the file pointer
	f, err := os.Open("/tmp/" + file.Filename)
	if err != nil {
		response.Errorcode = 1
		response.Errormsg = err
		return response
	}
	defer f.Close()

	// Get sum for whatever hash type has been asked for
	switch c.PostForm("hash") {
	case "md5":
		h := md5.New()
		if _, err := io.Copy(h, f); err != nil {
			response.Errorcode = 1
			response.Errormsg = err
			return response
		}
		response.Hashes = append(response.Hashes,
			HashResponse{Hash: hex.EncodeToString(h.Sum(nil)), Type: "md5"})
	default:
		response.Errorcode = 1
		response.Errormsg = errors.New("Hash type " + c.PostForm("hash") + " is not supported")
		return response
	}
	return response
}

func GetStringHash(c *gin.Context) StringResponse {
	response := StringResponse{Plaintext: c.PostForm("ct")}
	switch c.PostForm("hash") {
	case "all":
		sum0 := md5.Sum([]byte(response.Plaintext))
		response.Hashes = append(response.Hashes,
			HashResponse{Hash: hex.EncodeToString(sum0[:]), Type: "md5"})
		sum1 := sha1.Sum([]byte(response.Plaintext))
		response.Hashes = append(response.Hashes,
			HashResponse{Hash: hex.EncodeToString(sum1[:]), Type: "sha1"})
		sum2 := sha256.Sum224([]byte(response.Plaintext))
		response.Hashes = append(response.Hashes,
			HashResponse{Hash: hex.EncodeToString(sum2[:]), Type: "sha224"})
		sum3 := sha256.Sum256([]byte(response.Plaintext))
		response.Hashes = append(response.Hashes,
			HashResponse{Hash: hex.EncodeToString(sum3[:]), Type: "sha256"})
		sum4 := sha512.Sum384([]byte(response.Plaintext))
		response.Hashes = append(response.Hashes,
			HashResponse{Hash: hex.EncodeToString(sum4[:]), Type: "sha384"})
		sum5 := sha512.Sum512([]byte(response.Plaintext))
		response.Hashes = append(response.Hashes,
			HashResponse{Hash: hex.EncodeToString(sum5[:]), Type: "sha512"})
		sum6 := sha512.Sum512_224([]byte(response.Plaintext))
		response.Hashes = append(response.Hashes,
			HashResponse{Hash: hex.EncodeToString(sum6[:]), Type: "sha512_224"})
		sum7 := sha512.Sum512_256([]byte(response.Plaintext))
		response.Hashes = append(response.Hashes,
			HashResponse{Hash: hex.EncodeToString(sum7[:]), Type: "sha512_256"})
	case "md5":
		sum := md5.Sum([]byte(response.Plaintext))
		response.Hashes = append(response.Hashes,
			HashResponse{Hash: hex.EncodeToString(sum[:]), Type: "md5"})
	case "sha1":
		sum := sha1.Sum([]byte(response.Plaintext))
		response.Hashes = append(response.Hashes,
			HashResponse{Hash: hex.EncodeToString(sum[:]), Type: "sha1"})
	case "sha224":
		sum := sha256.Sum224([]byte(response.Plaintext))
		response.Hashes = append(response.Hashes,
			HashResponse{Hash: hex.EncodeToString(sum[:]), Type: "sha224"})
	case "sha256":
		sum := sha256.Sum256([]byte(response.Plaintext))
		response.Hashes = append(response.Hashes,
			HashResponse{Hash: hex.EncodeToString(sum[:]), Type: "sha256"})
	case "sha384":
		sum := sha512.Sum384([]byte(response.Plaintext))
		response.Hashes = append(response.Hashes,
			HashResponse{Hash: hex.EncodeToString(sum[:]), Type: "sha384"})
	case "sha512":
		sum := sha512.Sum512([]byte(response.Plaintext))
		response.Hashes = append(response.Hashes,
			HashResponse{Hash: hex.EncodeToString(sum[:]), Type: "sha512"})
	case "sha512_224":
		sum := sha512.Sum512_224([]byte(response.Plaintext))
		response.Hashes = append(response.Hashes,
			HashResponse{Hash: hex.EncodeToString(sum[:]), Type: "sha512_224"})
	case "sha512_256":
		sum := sha512.Sum512_256([]byte(response.Plaintext))
		response.Hashes = append(response.Hashes,
			HashResponse{Hash: hex.EncodeToString(sum[:]), Type: "sha512_256"})
	default:
		response.Errorcode = 1
		response.Errormsg = errors.New("Hash type " + c.PostForm("hash") + " is not supported")
	}
	return response
}

func main() {
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()
	router.LoadHTMLGlob("templates/*.html")
	//
	// Static files and resources
	//
	router.Static("/img", "./img")

	//
	// Main page
	//
	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", nil)
	})

	//
	// IS THE SERVICE ALIVE?
	//
	router.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"Message": "pong",})
	})

	//
	// Hash stuff
	//
	router.POST("/string", func(c *gin.Context) {
		c.JSON(http.StatusOK, GetStringHash(c))
	})

	router.POST("/file", func(c *gin.Context) {
		c.JSON(http.StatusOK, GetFileHash(c))
	})

	//
	// Start the server
	//
	router.Run(":7001")
}
