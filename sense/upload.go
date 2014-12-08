package sense

import (
	proto "code.google.com/p/goprotobuf/proto"
	"github.com/hello/sense/hello"
	"log"
)

type UploadService struct {
	client *SenseProtobufClient
}

func (s *UploadService) Upload(temp int32) (*hello.SyncResponse, error) {

	data := &hello.PeriodicData{}
	data.Temperature = &temp
	device_id := "D05FB81BE1E0"
	data.DeviceId = &device_id

	buff, err := proto.Marshal(data)
	req, err := s.client.NewProtobufRequest("POST", "/in/morpheus/pb2", buff)
	if err != nil {
		log.Println("yo yo")
		return &hello.SyncResponse{}, err
	}

	resp, err := s.client.Do(req)
	if err != nil {
		log.Println("yo yo yoyo yooyoy")
		return resp, err
	}

	return resp, err
}
