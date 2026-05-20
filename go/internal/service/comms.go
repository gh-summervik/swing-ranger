package service

import (
	"fmt"

	"github.com/summervik/swing-ranger/internal/config"
)

type CommsService struct{
	verbose bool
}

func NewCommsService(cfg config.Config) *CommsService{
	return &CommsService{
		verbose: cfg.Verbose,
	}
}

func (s *CommsService) Communicate(msg string){
	if s.verbose {
		fmt.Println(msg)
	}
}

func Communicate(msg string, err error){
	if err != nil{
		fmt.Println(err)

	}
}