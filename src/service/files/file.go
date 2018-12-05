package files

import (
	"bytes"
	"constants"
	"fmt"
	"gopkg.in/goftp"
	"mime/multipart"
	"os"
	"oss"
	"strings"
	"time"
)

import (
	"appcfg"
	"common"
	"errors"
	"io"
	"io/ioutil"
	. "logger"
	"net/http"
)

var UploadFailErr = errors.New("upload file failed")

type OSSInfo struct {
	Endpoint string
	ID       string
	Key      string
	Bucket   string
	Dir      string
}

type FTPInfo struct {
	Host     string
	Port     string
	Username string
	Password string
	Dir      string
}

func GetPlatformOSSDirName() string {
	if FST == constants.FILE_SERVER_OSS {
		return OssDirName
	} else if FST == constants.FILE_SERVER_FTP {
		return appcfg.GetString("ftp_dir_name", "")
	}
	return ""
}

func GetPlatformFTPInfo() *FTPInfo {
	return ftpInfo
}

func getFTP() (f *goftp.FTP, err error) {
	if f, err = goftp.Connect(appcfg.GetString("ftp_address", "")); err != nil {
		Err(err)
		return
	}

	if err = f.Login(appcfg.GetString("ftp_username", ""), appcfg.GetString("ftp_password", "")); err != nil {
		f.Close()
		Err(err)
		return
	}

	return
}

func UploadFileByBuffer(buffer []byte, objectName string, expireDays int) (fileID string, err error) {
	if FST == constants.FILE_SERVER_OSS {
		options := make([]oss.Option, 0)
		if expireDays > 0 {
			options = append(options, oss.Expires(time.Unix(time.Now().Unix()+int64(expireDays*24*3600), 0)))
		}
		if err = getBucket().PutObject(objectName, bytes.NewReader(buffer), options...); err != nil {
			Err(err)
			return
		}
	} else if FST == constants.FILE_SERVER_FTP {
		var f *goftp.FTP
		if f, err = getFTP(); err != nil {
			return
		} else {
			defer f.Close()
			tmpFile, _ := ioutil.TempFile(os.TempDir(), "tmp_")
			defer os.Remove(tmpFile.Name())

			if err = ioutil.WriteFile(tmpFile.Name(), buffer, 0644); err != nil {
				Err(err)
				return
			}

			if err = f.Stor(objectName, tmpFile); err != nil {
				Err(err)
			}
		}

	}

	fileID = objectName
	return
}

func DecryptDownloadFile(fileID string, filename string) (err error) {
	//check if exists
	_, err = os.Stat(filename)
	if err == nil || os.IsExist(err) {
		//exist do nothing
	} else {
		common.Execute("touch", filename)
		Info("get file:", fileID, "dest", filename)

		tmpFile, _ := ioutil.TempFile(os.TempDir(), "tmp_")
		defer os.Remove(tmpFile.Name())

		if FST == constants.FILE_SERVER_OSS {
			if err = downloadFile(fileID, tmpFile.Name(), getBucket()); err != nil {
				return
			}
		} else if FST == constants.FILE_SERVER_FTP {
			var f *goftp.FTP
			if f, err = getFTP(); err != nil {
				return
			} else {
				defer f.Close()

				if _, err = f.Retr(fileID, func(r io.Reader) error {
					if _, err = io.Copy(tmpFile, r); err != nil {
						return err
					}

					return err
				}); err != nil {
					Err(err)
					return
				}
			}
		}

		stat, _ := tmpFile.Stat()
		Info("file size:", tmpFile.Name(), stat.Size())

		if _, err = common.Execute(deFile, tmpFile.Name()); err != nil {
			return
		}
		if _, err = common.Execute("mv", tmpFile.Name(), filename); err != nil {
			return
		}
	}

	return
}

func EncryptUploadFileByName(name string, objectName string, needTmp bool) (location string, err error) {
	rn := name
	in, _ := os.OpenFile(name, os.O_RDWR, 0666)
	stat, _ := in.Stat()
	if needTmp {
		tmpFile, _ := ioutil.TempFile(os.TempDir(), "tmp_")
		defer os.Remove(tmpFile.Name())

		io.Copy(tmpFile, in)

		stat, _ = tmpFile.Stat()
		rn = tmpFile.Name()
	}

	common.Execute(enFile, rn)
	md5, _ := common.GenMD5F(rn)
	ron := objectName + "_" + md5

	if stat.Size() > 0 {
		Info("file size:", name, ron, rn, stat.Size())
		if FST == constants.FILE_SERVER_OSS {
			return uploadFile(rn, ron, getBucket())
		} else if FST == constants.FILE_SERVER_FTP {
			var f *goftp.FTP
			if f, err = getFTP(); err != nil {
				return
			} else {
				defer f.Close()

				var file *os.File
				if file, err = os.Open(rn); err != nil {
					Err(err)
					return
				}
				if err = f.Stor(ron, file); err != nil {
					Err(err)
					return
				}
				return ron, nil
			}
		}
	} else {
		Err("file size:", name, ron, rn, stat.Size())
		err = fmt.Errorf("upload file empty")
	}

	return
}

func EncryptUploadFileByNameNoMd5(name, objectName string, needTmp bool) (location string, err error) {
	rn := name
	in, _ := os.OpenFile(name, os.O_RDWR, 0666)
	stat, _ := in.Stat()
	if needTmp {
		tmpFile, _ := ioutil.TempFile(os.TempDir(), "tmp_")
		defer os.Remove(tmpFile.Name())

		io.Copy(tmpFile, in)

		stat, _ = tmpFile.Stat()
		rn = tmpFile.Name()
	}

	common.Execute(enFile, rn)
	// md5, _ := common.GenMD5F(rn)
	// ron := objectName + "_" + md5

	if stat.Size() > 0 {
		Info("file size:", name, objectName, rn, stat.Size())

		if FST == constants.FILE_SERVER_OSS {
			return uploadFile(rn, objectName, getBucket())
		} else if FST == constants.FILE_SERVER_FTP {
			var f *goftp.FTP
			if f, err = getFTP(); err != nil {
				return
			} else {
				defer f.Close()

				var file *os.File
				if file, err = os.Open(rn); err != nil {
					Err(err)
					return
				}
				if err = f.Stor(objectName, file); err != nil {
					Err(err)
					return
				}
				return objectName, nil
			}
		}
	} else {
		Err("file size:", name, objectName, rn, stat.Size())
		err = fmt.Errorf("upload file empty")
	}

	return
}

func UploadFileByName(name string, objectName string, expireDays int) (location string, err error) {
	rn := name
	in, _ := os.OpenFile(name, os.O_RDWR, 0666)
	if stat, e := in.Stat(); err != nil {
		Err(e)
		err = e
		return
	} else if stat.Size() > 0 {
		Info("file size:", name, objectName, rn, stat.Size())

		if FST == constants.FILE_SERVER_OSS {
			options := make([]oss.Option, 0)
			if expireDays > 0 {
				options = append(options, oss.Expires(time.Unix(time.Now().Unix()+int64(expireDays*24*3600), 0)))
			}
			if err = getBucket().PutObjectFromFile(objectName, rn, options...); err != nil {
				Err(err)
				return
			}
			location = objectName
			return
		} else if FST == constants.FILE_SERVER_FTP {
			var f *goftp.FTP
			if f, err = getFTP(); err != nil {
				return
			} else {
				defer f.Close()

				var file *os.File
				if file, err = os.Open(rn); err != nil {
					Err(err)
					return
				}
				if err = f.Stor(objectName, file); err != nil {
					Err(err)
					return
				}
				return objectName, nil
			}
		}
	} else {
		Err("file size:", name, objectName, rn, stat.Size())
		err = fmt.Errorf("upload file empty")
	}

	return
}

func Encrypt(name string) error {
	if f, err := os.OpenFile(name, os.O_RDWR, 0666); err != nil {
		if strings.Contains(err.Error(), constants.FILE_NOT_EXIST_MSG) {
			Info("file not exist", name)
		} else {
			Err(err)
		}
		return err
	} else if stat, err := f.Stat(); err != nil {
		Err(err)
		return err
	} else if stat.Size() > 0 {
		if _, err := common.Execute(enFile, name); err != nil {
			return err
		}
	} else {
		Err("encrypt file size 0", name)
		return constants.FileSizeZero
	}

	return nil
}

func EncryptUploadFile1(f *os.File, objectName string) (location string, err error) {
	rn := f.Name()
	in := f
	stat, _ := in.Stat()

	common.Execute(enFile, rn)
	md5, _ := common.GenMD5F(rn)
	ron := objectName + "_" + md5

	if stat.Size() > 0 {
		Info("file size:", rn, ron, rn, stat.Size())

		if FST == constants.FILE_SERVER_OSS {
			return uploadFile(rn, ron, getBucket())
		} else if FST == constants.FILE_SERVER_FTP {
			var f *goftp.FTP
			if f, err = getFTP(); err != nil {
				return
			} else {
				defer f.Close()

				var file *os.File
				if file, err = os.Open(rn); err != nil {
					Err(err)
					return
				}
				if err = f.Stor(ron, file); err != nil {
					Err(err)
					return
				}

				return ron, nil
			}
		}
	} else {
		Err("file size:", rn, ron, rn, stat.Size())
		err = fmt.Errorf("upload file empty")
	}

	return
}

func DeleteFile(fileID string, platform string) (err error) {
	if fileID == "" {
		return fmt.Errorf("fdfs delete empty fileid")
	}

	if FST == constants.FILE_SERVER_OSS {
		if err = getBucket().DeleteObject(fileID); err != nil {
			Err(err)
		} else {
			Info("deleted file:", fileID)
		}
	} else if FST == constants.FILE_SERVER_FTP {
		var f *goftp.FTP
		if f, err = getFTP(); err != nil {
			return
		} else {
			defer f.Close()

			if err = f.Dele(fileID); err != nil {
				Err(err)
				return
			}
		}
	}

	return
}

func downloadFile(fileID string, saveName string, b *oss.Bucket) (err error) {
	if FST == constants.FILE_SERVER_OSS {
		if err = b.GetObjectToFile(fileID, saveName); err != nil {
			Err(err)
		}
	} else if FST == constants.FILE_SERVER_FTP {
		var f *goftp.FTP
		if f, err = getFTP(); err != nil {
			return
		} else {
			defer f.Close()

			var file *os.File
			if file, err = os.Open(saveName); err != nil {
				Err(err)
				return
			}
			if _, err = f.Retr(fileID, func(r io.Reader) error {
				if _, err = io.Copy(file, r); err != nil {
					return err
				}

				return err
			}); err != nil {
				Err(err)
				return
			}
		}
	}

	return
}

func DeleteFiles(files []string) (err error) {
	if len(files) <= 0 {
		return
	}
	if FST == constants.FILE_SERVER_OSS {
		if _, err = bucket.DeleteObjects(files); err != nil {
			Err(err)
		}
	} else if FST == constants.FILE_SERVER_FTP {
		var f *goftp.FTP
		if f, err = getFTP(); err != nil {
			return
		} else {
			defer f.Close()

			for _, fe := range files {
				if err = f.Dele(fe); err != nil {
					Err(err)
					return
				}
			}
		}
	}

	return
}

func uploadFile(filename string, objectName string, b *oss.Bucket) (fileID string, err error) {
	Info("uploadFile:", filename, objectName, b)
	if err = b.PutObjectFromFile(objectName, filename); err != nil {
		Err(err)
		return
	}
	fileID = objectName
	return
}

func getBucket() *oss.Bucket {
	if FST == constants.FILE_SERVER_OSS {
		return bucket
	}

	return nil
}

func initBucket(oss_endpoint string, oss_accessid string, oss_accesskey string, bucketName string) (b *oss.Bucket) {
	if client, err := oss.New(oss_endpoint, oss_accessid, oss_accesskey, appcfg.GetInt32("oss_concurrency", 2000)); err != nil {
		panic(err)
	} else {
		if b, err = client.Bucket(bucketName); err != nil {
			panic(err)
		}
	}
	return
}

func postFile(filename string, targetUrl string) error {
	bodyBuf := &bytes.Buffer{}
	bodyWriter := multipart.NewWriter(bodyBuf)

	//关键的一步操作
	bodyWriter.WriteField("file_name", filename)
	fileWriter, err := bodyWriter.CreateFormFile("file_content", filename)
	if err != nil {
		Err(err)
		return err
	}

	//打开文件句柄操作
	fh, err := os.Open(filename)
	if err != nil {
		Err(err)
		return err
	}

	//iocopy
	_, err = io.Copy(fileWriter, fh)
	if err != nil {
		Err(err)
		return err
	}

	contentType := bodyWriter.FormDataContentType()
	bodyWriter.Close()

	resp, err := http.Post(targetUrl, contentType, bodyBuf)
	if err != nil {
		Err(err)
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		resp_body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			Err(err)
			return err
		}
		Err("upload resp status:", resp.Status, "content:", string(resp_body))
		return UploadFailErr
	}
	return nil
}
