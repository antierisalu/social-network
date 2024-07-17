package app

import (
	"encoding/base64"
	"fmt"
	"os"
	"path/filepath"
)

type Avatars struct {
	NewAvatar AvatarImgData
	Dir       string
	Path      string
}

type AvatarImgData struct {
	UserID           int
	Base64String     string
	Buffer           []byte
	FilePath         string
	ShortPath        string
	FileName         string
	FormatedFileName string
	FileType         string
}

func (a *Avatars) InitPath() {
	dir, err := os.Getwd()
	if err != nil {
		fmt.Printf("failed to get current directory: %v", err)
		return
	}

	avatarsDir := filepath.Join(dir, a.Dir)
	if _, err := os.Stat(avatarsDir); os.IsNotExist(err) {
		if err := os.MkdirAll(avatarsDir, 0755); err != nil {
			fmt.Printf("failed to create directory: %v", err)
			return
		}
		fmt.Println("Avatars directory created:", avatarsDir)
	} else if err != nil {
		fmt.Printf("error checking directory: %v", err)
		return
	}
	a.Path = avatarsDir
}

func (a *Avatars) SaveNewAvatar() error {
	a.InitPath()

	err := a.NewAvatar.DecodeBase64()
	if err != nil {
		return fmt.Errorf("error: type Avatars method(SaveNewAvatar):\n%s", err)
	}

	err = a.NewAvatar.GetFileType()
	if err != nil {
		return fmt.Errorf("error: type Avatars method(SaveNewAvatar):\n%s", err)
	}

	err = a.NewAvatar.FormatedName()
	if err != nil {
		return fmt.Errorf("error: type Avatars method(SaveNewAvatar):\n%s", err)
	}

	a.NewAvatar.FilePath = filepath.Join(a.Path, fmt.Sprintf("%s%s", a.NewAvatar.FormatedFileName, a.NewAvatar.FileType))
	a.NewAvatar.ShortPath = fmt.Sprintf("/avatars/%s%s", a.NewAvatar.FormatedFileName, a.NewAvatar.FileType)

	err = a.SaveImage()
	if err != nil {
		return fmt.Errorf("error: type Avatars method(SaveNewAvatar):\n%s", err)
	}
	return nil
}

func (a *Avatars) SaveImage() error {
	file, err := os.Create(a.NewAvatar.FilePath)
	if err != nil {
		return fmt.Errorf("error: type Avatars method(SaveImage):\n%s", err)

	}
	defer file.Close()

	_, err = file.Write(a.NewAvatar.Buffer)
	if err != nil {
		return fmt.Errorf("error: type Avatars method(SaveImage):\n%s", err)
	} else { // remove this as needed
		fmt.Println("Avatar saved as:", a.NewAvatar.ShortPath)
	}
	return nil
}

func (img *AvatarImgData) FormatedName() error {
	if img.UserID < 0 {
		return fmt.Errorf("error: type AvatarImgData method(FromatedName) generating a formated file name from id: %v", img.UserID)
	}
	img.FormatedFileName = fmt.Sprint(img.UserID)
	return nil
}

func (img *AvatarImgData) GetFileType() error {
	img.FileType = filepath.Ext(img.FileName)
	if img.FileType == "" {
		return fmt.Errorf("error: type AvatarImgData method(GetFileType) getting fileType from file: %s", img.FileName)
	}
	return nil
}

func (img *AvatarImgData) DecodeBase64() error {
	decodedData, err := base64.StdEncoding.DecodeString(img.Base64String)
	if err != nil {
		return fmt.Errorf("error: type AvatarImgData method(DecodeBase64) failed to decode base64 string: %v", err)
	}
	img.Buffer = decodedData
	return nil
}
