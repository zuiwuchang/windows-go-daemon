package main

/*
#ifndef UNICODE
#define UNICODE
#endif

#include <windows.h>
#include <TlHelp32.h>
int native_InstallService(char* name,
	char* show,
	char* description,
	char* run,
	DWORD start
){
	//打開 scm
	SC_HANDLE scm = OpenSCManager(NULL,NULL,SC_MANAGER_CREATE_SERVICE);
	if(!scm){
		return 1;
	}

	int result = 0;
	int n = MultiByteToWideChar(CP_UTF8,0,name,-1,0,0);
	wchar_t* wName = (wchar_t*)malloc(sizeof(wchar_t) * n);
	if(!wName){
		result = 5;
		goto END;
	}
	MultiByteToWideChar(CP_UTF8,0,name,-1,wName,n);

	n = MultiByteToWideChar(CP_UTF8,0,show,-1,0,0);
	wchar_t* wShow = (wchar_t*)malloc(sizeof(wchar_t) * n);
	if(!wShow){
		result = 5;
		goto END;
	}
	MultiByteToWideChar(CP_UTF8,0,show,-1,wShow,n);

	n = MultiByteToWideChar(CP_UTF8,0,description,-1,0,0);
	wchar_t* wDescription = (wchar_t*)malloc(sizeof(wchar_t) * n);
	if(!wDescription){
		result = 5;
		goto END;
	}
	MultiByteToWideChar(CP_UTF8,0,description,-1,wDescription,n);

	n = MultiByteToWideChar(CP_UTF8,0,run,-1,0,0);
	wchar_t* wRun = (wchar_t*)malloc(sizeof(wchar_t) * n);
	if(!wRun){
		result = 5;
		goto END;
	}
	MultiByteToWideChar(CP_UTF8,0,run,-1,wRun,n);


	//創建服務
	SC_HANDLE s = CreateService(scm,
		wName,
		wShow,
		SERVICE_ALL_ACCESS,
		SERVICE_WIN32_OWN_PROCESS,	//exe中只提供單個服務
		start,					//手動啓動
		SERVICE_ERROR_NORMAL,
		wRun,
		NULL,
		NULL,
		NULL,
		NULL,
		NULL
	);
	if(s){
		//修改服務描述
		SERVICE_DESCRIPTION serviceDescription;
		serviceDescription.lpDescription = wDescription;
		if(!ChangeServiceConfig2(s,SERVICE_CONFIG_DESCRIPTION,&serviceDescription))
		{
			result = 4;
		}
		//關閉 服務
		CloseServiceHandle(s);
	}else{
		if(ERROR_SERVICE_EXISTS == GetLastError()){
			result = 2;
		}else{
			result = 3;
		}
	}
END:
	if(wName){
		free(wName);
	}
	if(wShow){
		free(wShow);
	}
	if(wDescription){
		free(wDescription);
	}
	if(wRun){
		free(wRun);
	}

	//關閉 scm
	CloseServiceHandle(scm);
	return result;
}
int native_UnstallService(char* name){
	//打開 scm
	SC_HANDLE scm = OpenSCManager(NULL,NULL,SC_MANAGER_CREATE_SERVICE);
	if(!scm){
		return 1;
	}

	int result = 0;
	int n = MultiByteToWideChar(CP_UTF8,0,name,-1,0,0);
	wchar_t* wName = (wchar_t*)malloc(sizeof(wchar_t) * n);
	if(!wName){
		result = 5;
		goto END;
	}
	MultiByteToWideChar(CP_UTF8,0,name,-1,wName,n);


	//打開 服務
	SC_HANDLE s = OpenService(scm,wName,SERVICE_ALL_ACCESS);
	if(s){
		if(!DeleteService(s))
		{
			result = 3;
		}
		//關閉 服務句柄
		CloseServiceHandle(s);
	}else{
		result = 2;
	}
END:
	if(wName){
		free(wName);
	}
	CloseServiceHandle(scm);
	return result;
}
//設置工作 目錄
int native_SetCurrentDirectory(char* dir){
	int n = MultiByteToWideChar(CP_UTF8,0,dir,-1,0,0);
	wchar_t* wDir = (wchar_t*)malloc(sizeof(wchar_t) * n);
	if(!wDir){
		return 1;
	}
	MultiByteToWideChar(CP_UTF8,0,dir,-1,wDir,n);

	int result = 0;
	if(!SetCurrentDirectory(wDir)){
		result = 1;
	}
	free(wDir);
	return result;
}
//創建 進程
int native_CreateProcess(char* bin,char* params,char* dir,DWORD* pid,void* hProcess,void* hThread){
	int result = 0;
	int n = MultiByteToWideChar(CP_UTF8,0,bin,-1,0,0);
	wchar_t* wBin = (wchar_t*)malloc(sizeof(wchar_t) * n);
	if(!wBin){
		result = 1;
		goto END;
	}
	MultiByteToWideChar(CP_UTF8,0,bin,-1,wBin,n);

	n = MultiByteToWideChar(CP_UTF8,0,params,-1,0,0);
	wchar_t* wParams = (wchar_t*)malloc(sizeof(wchar_t) * n);
	if(!wParams){
		result = 1;
		goto END;
	}
	MultiByteToWideChar(CP_UTF8,0,params,-1,wParams,n);

	n = MultiByteToWideChar(CP_UTF8,0,dir,-1,0,0);
	wchar_t* wDir = (wchar_t*)malloc(sizeof(wchar_t) * n);
	if(!wDir){
		result = 1;
		goto END;
	}
	MultiByteToWideChar(CP_UTF8,0,dir,-1,wDir,n);


	PROCESS_INFORMATION process;
	STARTUPINFO startupinfo;
	memset(&startupinfo,0,sizeof(STARTUPINFO));
	startupinfo.cb=sizeof(STARTUPINFO);
	if(CreateProcess(wBin,
		wParams,
		NULL,
		NULL,
		FALSE,
		CREATE_SUSPENDED,
		NULL,
		wDir,
		&startupinfo,
		&process)
	){
		*pid = process.dwProcessId;
		*((HANDLE*)hProcess) = process.hProcess;
		*((HANDLE*)hThread) = process.hThread;
	}else{
		result = 1;
	}
END:
	if(wBin){
		free(wBin);
	}
	if(wParams){
		free(wParams);
	}
	if(wDir){
		free(wDir);
	}
	return result;
}

VOID WINAPI ServiceMain(
	DWORD dwArgc,
	LPTSTR* lpszArgv
);
VOID WINAPI ServiceControlHandler(
	DWORD fdwControl
);
SERVICE_STATUS g_status;
SERVICE_STATUS_HANDLE g_hServer;
wchar_t g_Name[MAX_PATH];
typedef uintptr_t (__stdcall* control_ft)();
control_ft g_control_run;
control_ft g_control_wait;
control_ft g_control_close;
int native_service_main(char* name,
	control_ft run,
	control_ft wait,
	control_ft close
	){
	//設置回調
	g_control_run = run;
	g_control_wait = wait;
	g_control_close = close;

	//獲取服務名
	MultiByteToWideChar(CP_UTF8,0,name,-1,g_Name,MAX_PATH);

	//定義待 註冊服務 數組
	//每個服務 會被運行在一個新啓的 線程中
	SERVICE_TABLE_ENTRY serviceTable[] = {
		{
			g_Name,	//服務名
			ServiceMain	//服務 入口點函數
		},
		{NULL,NULL} // 以空數組 代表結束
	};

	//註冊服務 並等待服務停止
	StartServiceCtrlDispatcher(serviceTable);
	return 0;
}
VOID WINAPI ServiceMain(
	DWORD dwArgc,
	LPTSTR* lpszArgv
)
{
	//註冊 ControlHandler
	//返回的 句柄 不需要關閉
	g_hServer = RegisterServiceCtrlHandler(
		g_Name,//服務名 必須和 SERVICE_TABLE_ENTRY中一致
		ServiceControlHandler	//控制回調 在 StartServiceCtrlDispatcher 線程中被回調
	);
	if(!g_hServer){
		//error
		return;
	}
	//通知 scm 已經開始 繼續 初始化 服務
	g_status.dwServiceType = SERVICE_WIN32_OWN_PROCESS;
	g_status.dwCurrentState = SERVICE_START_PENDING;	//設置當前服務 狀態
	g_status.dwControlsAccepted   = SERVICE_ACCEPT_STOP | SERVICE_ACCEPT_SHUTDOWN;	//設置允許的 ControlHandler
	g_status.dwWin32ExitCode = 0;
	g_status.dwServiceSpecificExitCode = 0;
	g_status.dwCheckPoint   = 0;
	g_status.dwWaitHint   = 0;
	SetServiceStatus(g_hServer,&g_status);

	//運行服務
	if(g_control_run() != 0) //初始化 並運行
	{
		//初始化 失敗 停止服務
		g_status.dwWin32ExitCode = ERROR_SERVICE_SPECIFIC_ERROR ;
		g_status.dwServiceSpecificExitCode  = 1;
		g_status.dwCurrentState = SERVICE_STOPPED;
		SetServiceStatus(g_hServer,&g_status);
		return ;
	}


	//通知 scm 服務 初始化 完成 正在運行
	g_status.dwWaitHint = 0;
	g_status.dwCheckPoint = 0;
	g_status.dwCurrentState = SERVICE_RUNNING;
	SetServiceStatus(g_hServer,&g_status);

	//等待服務結束
	g_control_wait();
}
VOID WINAPI ServiceControlHandler(
	DWORD fdwControl
){
	switch(fdwControl)
	{
	case SERVICE_CONTROL_STOP:	//停止 服務
	case SERVICE_CONTROL_SHUTDOWN:
		//通知 scm 已經 得到 停止通知
		g_status.dwWaitHint = 1000 * 6;	//設置預計 完成 需要時間
		g_status.dwCheckPoint = 0;		//更新 執行 進度
		g_status.dwCurrentState = SERVICE_STOP_PENDING;	//設置當前服務 狀態

		SetServiceStatus(g_hServer,&g_status);


		//執行關閉操作
		g_control_close();
		//等待服務結束
		g_control_wait();

		//通知 scm 已經 停止服務
		g_status.dwWaitHint = 0;
		g_status.dwCheckPoint = 0;
		g_status.dwCurrentState = SERVICE_STOPPED;	//設置當前服務 狀態
		SetServiceStatus(g_hServer,&g_status);
		return;
	}

	SetServiceStatus(g_hServer,&g_status);
}
typedef uintptr_t (__stdcall* for_process_ft)(DWORD pid,DWORD pPid);
#include <stdio.h>
void native_for_process(for_process_ft for_process){
	HANDLE handle = CreateToolhelp32Snapshot(TH32CS_SNAPPROCESS,0);
	PROCESSENTRY32 process;
	process.dwSize = sizeof(PROCESSENTRY32);
	if(INVALID_HANDLE_VALUE != handle){
		if(Process32First(handle,&process)){
			do{
				for_process(process.th32ProcessID,process.th32ParentProcessID);
			}while(Process32Next(handle,&process));
		}

		CloseHandle(handle);
	}
}
int native_KillByPid(DWORD pid){
	int result = 0;
	HANDLE handle = OpenProcess(PROCESS_ALL_ACCESS,FALSE,pid);
	if(!handle){
		result = 1;
		return result;
	}
	if(!TerminateProcess(handle,1)){
		result = 1;
	}
	CloseHandle(handle);
	return result;
}
int native_SetToken(){
	HANDLE hToken;
	if(OpenProcessToken(GetCurrentProcess(),TOKEN_ADJUST_PRIVILEGES | TOKEN_QUERY,&hToken)){
		TOKEN_PRIVILEGES tkp;

		LookupPrivilegeValue( NULL,SE_DEBUG_NAME,&tkp.Privileges[0].Luid );//修改进程权限
		tkp.PrivilegeCount=1;
		tkp.Privileges[0].Attributes=SE_PRIVILEGE_ENABLED;
		AdjustTokenPrivileges( hToken,FALSE,&tkp,sizeof tkp,NULL,NULL );//通知系统修改进程权限

		if(GetLastError() != ERROR_SUCCESS){
			return 1;
		}
		return 0;
	}
	return 1;
}
*/
import "C"

import (
	"errors"
	"log"
	"os"
	"path/filepath"
	"syscall"
	"unsafe"
)

func SetToken() (e error) {
	if 0 != C.native_SetToken() {
		e = errors.New("SetToken error")
	}
	return
}

//安裝服務
var ErrorChangeServiceConfig2 error = errors.New("ChangeServiceConfig2 error")

func InstallService(name,
	show,
	description,
	bin string,
	auto bool,
) (e error) {
	cName := C.CString(name)
	cShow := C.CString(show)
	cDescription := C.CString(description)
	cBin := C.CString(bin)

	var start C.DWORD
	if auto {
		start = C.SERVICE_AUTO_START
	} else {
		start = C.SERVICE_DEMAND_START
	}
	switch C.native_InstallService(
		cName,
		cShow,
		cDescription,
		cBin,

		start,
	) {
	case 0:
	case 1:
		e = errors.New("OpenSCManager error")
	case 2:
		e = errors.New("Service is already exists")
	case 3:
		e = errors.New("CreateService error")
	case 4:
		e = ErrorChangeServiceConfig2
	case 5:
		e = errors.New("bad malloc")
	default:
		e = errors.New("Unknow")
	}

	C.free(unsafe.Pointer(cName))
	C.free(unsafe.Pointer(cShow))
	C.free(unsafe.Pointer(cDescription))
	C.free(unsafe.Pointer(cBin))
	return
}

//卸載服務
func UnstallService(name string) (e error) {
	cName := C.CString(name)
	switch C.native_UnstallService(cName) {
	case 0:
	case 1:
		e = errors.New("OpenSCManager error")
	case 2:
		e = errors.New("OpenService error")
	case 3:
		e = errors.New("DeleteService error")
	case 5:
		e = errors.New("bad malloc")
	default:
		e = errors.New("Unknow")
	}
	C.free(unsafe.Pointer(cName))
	return
}

//初始化 工作目錄
func InitWorkDirectory() (e error) {
	var dir string
	dir, e = filepath.Abs(os.Args[0])
	if e != nil {
		return
	}
	dir = filepath.Dir(dir)

	cDir := C.CString(dir)
	if C.native_SetCurrentDirectory(cDir) != 0 {
		e = errors.New("SetCurrentDirectory error")
	}
	C.free(unsafe.Pointer(cDir))
	return
}

//運行一個 程式 並等待 其結束
func CreateProcess(bin, param, dir string) (p Process, e error) {
	cBin := C.CString(bin)
	cParam := C.CString(param)
	cDir := C.CString(dir)
	var pid C.DWORD
	switch C.native_CreateProcess(cBin, cParam, cDir,
		&pid,
		unsafe.Pointer(&p.Process),
		unsafe.Pointer(&p.Thread),
	) {
	case 0:
		p.Pid = (uint)(pid)
	case 1:
		e = errors.New("CreateProcess error")
	default:
		e = errors.New("unknow")
	}

	C.free(unsafe.Pointer(cBin))
	C.free(unsafe.Pointer(cParam))
	C.free(unsafe.Pointer(cDir))
	return
}

type Process struct {
	Pid     uint
	Process uintptr
	Thread  uintptr
}

func (p *Process) Resume() {
	C.ResumeThread((C.HANDLE)(p.Thread))
}
func (p *Process) Wait() {
	C.WaitForSingleObject((C.HANDLE)(p.Process), C.INFINITE)
}
func (p *Process) Close() {
	if p.Process == 0 {
		return
	}
	C.CloseHandle((C.HANDLE)(p.Process))
	C.CloseHandle((C.HANDLE)(p.Thread))

	p.Process = 0
	p.Thread = 0
}
func (p *Process) Kill() {
	C.TerminateProcess((C.HANDLE)(p.Process), 1)

	p.Close()
}
func (p *Process) KillChilds() {
	ps := make(map[uint]uint)
	C.native_for_process(
		(*[0]byte)(unsafe.Pointer(syscall.NewCallback(func(pid, pPid uint) uintptr {
			if pPid != 0 && pid != p.Pid {
				ps[pid] = pPid
			}
			return 0
		}))),
	)
	log.Println(ps)
	keys := make(map[uint]bool)
	kills := make(map[uint]bool)
	var ok bool
	for pid, _ := range ps {
		if _, ok = keys[pid]; ok {
			kills[pid] = true
			continue
		}

		if p.isChild(pid, ps, keys) {
			kills[pid] = true
		}
	}
	log.Println("kills", kills)
	for pid, _ := range kills {
		log.Println(
			"kill %v %v",
			pid,
			C.native_KillByPid((C.DWORD)(pid)),
		)
	}
}
func (p *Process) isChild(pid uint, ps map[uint]uint, keys map[uint]bool) bool {
	tmp := make(map[uint]bool)
	var ok bool
	for pid != 0 {
		if _, ok = tmp[pid]; ok {
			//loop busy
			return false
		} else if _, ok = keys[pid]; ok {
			//find
			for k, _ := range tmp {
				keys[k] = true
			}
			return true
		} else if pid == p.Pid {
			//find
			for k, _ := range tmp {
				keys[k] = true
			}
			return true
		}

		tmp[pid] = true
		//get
		if pid, ok = ps[pid]; !ok {
			return false
		}

	}
	return false
}
func service_main(s *Service, name string) {
	cName := C.CString(name)
	C.native_service_main(
		cName,
		(*[0]byte)(
			unsafe.Pointer(
				syscall.NewCallback(func() uintptr {
					if e := s.Run(); e != nil {
						return 1
					}
					return 0
				}),
			),
		),
		(*[0]byte)(
			unsafe.Pointer(
				syscall.NewCallback(func() uintptr {
					s.Wait()
					return 0
				}),
			),
		),
		(*[0]byte)(
			unsafe.Pointer(
				syscall.NewCallback(func() uintptr {
					s.Close()
					return 0
				}),
			),
		),
	)
	C.free(unsafe.Pointer(cName))
}
