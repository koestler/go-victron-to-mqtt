package ftpServer

import (
	"github.com/koestler/go-ve-sensor/storage"
	"image/jpeg"
	"log"
	"github.com/disintegration/imaging"
	"bytes"
)

func onFileClose(vf *VirtualFile) {
	//log.Printf("ftpServer: onFileClose device.name=%v filePath=%v", vf.device.Name, vf.filePath)

	inpImg, err := jpeg.Decode(vf)
	if err != nil {
		log.Printf("onFileClose: cannot read file as jpg, err=%v", err)
		return;
	}

	oupImg := imaging.Resize(inpImg, 640, 0, imaging.Box)

	imageBuffer := new(bytes.Buffer)
	err = jpeg.Encode(imageBuffer, oupImg, nil)
	if err != nil {
		log.Printf("onFileClose: cannot encode as jpg, err=%v", err)
		return
	}

	picture := storage.Picture{
		Created:   vf.modified,
		JpegThumb: imageBuffer.Bytes(),
		JpegRaw:   vf.buffer,
	}
	storage.PictureDb.SetPicture(vf.device, &picture)
}
