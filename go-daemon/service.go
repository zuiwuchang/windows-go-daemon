package main

import (
	"go-daemon/configure"
	"log"
	"path/filepath"
	"time"
)

const (
	commandExit = 1
)

type ICommand interface {
	Command() uint16
	Param() interface{}
}
type commandImpl struct {
	command uint16
	param   interface{}
}

func (c commandImpl) Command() uint16 {
	return c.command
}
func (c commandImpl) Param() interface{} {
	return c.param
}
func NewCommand(cmd uint16, param interface{}) ICommand {
	return commandImpl{
		command: cmd,
		param:   param,
	}
}

type Service struct {
	chExit chan (bool)

	chCmd chan (ICommand)

	chCreate      chan (bool)
	chCreateOk    chan (Process)
	chProcessExit chan (Process)

	flagRun bool

	process Process
}

func newService() *Service {
	return &Service{
		chExit:        make(chan (bool), 1),
		chCreate:      make(chan (bool), 2),
		chCreateOk:    make(chan (Process)),
		chProcessExit: make(chan (Process)),

		chCmd: make(chan (ICommand), 1),
	}
}
func (s *Service) Run() error {
	cnf := configure.GetConfigure()
	var dir string
	if cnf.Directory == "" {
		abs, e := filepath.Abs(cnf.Bin)
		if e != nil {
			return e
		}
		dir = filepath.Dir(abs)
	} else if filepath.IsAbs(cnf.Directory) {
		dir = cnf.Directory
	} else {
		abs, e := filepath.Abs(cnf.Directory)
		if e != nil {
			return e
		}
		dir = abs
	}

	s.flagRun = true
	go s.run(cnf.Bin, cnf.Params, dir)
	return nil
}
func (s *Service) run(bin, params, dir string) {
	go s.createProcess(bin, params, dir)
	for s.flagRun {
		select {
		case <-s.chCreate:
			go s.createProcess(bin, params, dir)
		case p := <-s.chCreateOk:
			s.resumeProcess(p)
			//go s.Close()
		case p := <-s.chProcessExit:
			//殺死子進程
			p.KillChilds()

			//創建新進程
			go s.createProcess(bin, params, dir)

		case cmd := <-s.chCmd:
			if cmd.Command() == commandExit {
				//殺死 進程
				p := s.process
				p.Kill()

				//殺死子進程
				p.KillChilds()

				close(s.chExit)
				return
			}
		}
	}
	close(s.chExit)
}

func (s *Service) createProcess(bin, params, dir string) {
	process, e := CreateProcess(bin, params, dir)
	log.Println(bin, params, dir)
	if e == nil { //創建成功
		log.Println("success", process)
		s.chCreateOk <- process
	} else { //創建 任務 失敗
		log.Println(e)

		//等待 1秒 後 重新創建
		time.Sleep(time.Second)
		s.chCreate <- true
		return
	}
}
func (s *Service) resumeProcess(process Process) {
	s.process = process
	//恢復進程 運行
	process.Resume()

	//等待 線程結束
	go func() {
		process.Wait()
		process.Close()

		//通知 進程 已結束
		s.chProcessExit <- process
	}()
}

func (s *Service) Close() {
	s.chCmd <- NewCommand(commandExit, nil)
}
func (s *Service) Wait() {
	<-s.chExit
	return
}
