package sense

import (
	proto "code.google.com/p/goprotobuf/proto"
	"encoding/hex"
	"github.com/hello/sense/hello"
	"log"
)

type UploadService struct {
	client *SenseProtobufClient
}

func (s *UploadService) Upload(temp int32, deviceId, aesKey string) (*hello.SyncResponse, error) {
	keybytes, _ := hex.DecodeString(aesKey)
	log.Printf("key bytes: %v\n", keybytes)
	data := &hello.BatchedPeriodicData{}

	for i := 0; i < 2; i++ {
		periodic := &hello.PeriodicData{}
		periodic.Temperature = &temp
		data.Data = append(data.Data, periodic)
	}

	data.DeviceId = &deviceId
	fw := int32(888)
	data.FirmwareVersion = &fw

	buff, err := proto.Marshal(data)
	req, err := s.client.NewProtobufRequest("POST", "/in/sense/batch", buff, aesKey)
	if err != nil {
		return &hello.SyncResponse{}, err
	}

	resp, err := s.client.Do(req, aesKey)
	if err != nil {
		return resp, err
	}

	return resp, err
}
