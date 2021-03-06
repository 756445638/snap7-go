package snap7go

//#cgo CFLAGS: -I .
//#include "snap7.h"
//#include <stdlib.h>
import "C"
import (
	"fmt"
	"unsafe"
)

func Srv_Create() (server S7Object) {
	server = C.Srv_Create()
	return
}

func NewS7Server() *S7Server {
	server := &S7Server{}
	server.server = Srv_Create()
	return server
}

type S7Server struct {
	server S7Object
}

func (s *S7Server) SetEventsCallback(handle func(*TSrvEvent)) error {
	return Srv_SetEventsCallback(s.server, func(usrPtr uintptr, event *TSrvEvent) {
		handle(event)
	}, uintptr(s.server))
}

func (s *S7Server) SetReadEventsCallback(handle func(*TSrvEvent)) error {
	return Srv_SetReadEventsCallback(s.server, func(usrPtr uintptr, event *TSrvEvent) {
		handle(event)
	}, uintptr(s.server))
}

func (s *S7Server) SetRWAreaCallback(handle func(sender int32, operation Operation, tag *PS7Tag, userData uintptr) SrvErrCode) error {
	return Srv_SetRWAreaCallback(s.server,
		func(_ uintptr, sender int32, operation Operation, tag *PS7Tag, userData uintptr) SrvErrCode {
			return handle(sender, Operation(operation), tag, userData)
		}, uintptr(s.server))
}

func (s *S7Server) SetRWAreaCallbackInterface(handle RWAreaCallbackInterface) error {
	return Srv_SetRWAreaCallback(s.server,
		func(_ uintptr, sender int32, operation Operation, tag *PS7Tag, userData uintptr) SrvErrCode {
			if operation == OperationRead {
				data := make([]byte, DataLength(S7WL(tag.WordLen), tag.Size))
				errCode := handle.Read(sender, tag, data)
				if errCode != 0 {
					return errCode
				}
				CopyToC(data, userData)
				return errCode
			} else {
				// write
				data := GetBytesFromC(userData, int(DataLength(S7WL(tag.WordLen), tag.Size)))
				return handle.Write(sender, tag, data)
			}
		}, uintptr(s.server))
}

type RWAreaCallbackInterface interface {
	Read(sender int32, tag *PS7Tag, ret []byte) (errCode SrvErrCode)
	Write(sender int32, tag *PS7Tag, data []byte) (errCode SrvErrCode)
}

func (s *S7Server) Destroy() {
	C.Srv_Destroy((*C.S7Object)(unsafe.Pointer(&s.server)))
	return
}

/*
	ParamNumber 为P_u16_LocalPort的时候 value的数据是uint16 其他情况类似的
*/
//int S7API Srv_GetParam(S7Object Server, int ParamNumber, void *pValue);
func (s *S7Server) GetParam(paraNumber ParamNumber) (value interface{}, err error) {
	var pValue unsafe.Pointer
	switch paraNumber {
	case P_u16_LocalPort:
		pValue = unsafe.Pointer(new(uint16))
	case P_i32_WorkInterval:
		pValue = unsafe.Pointer(new(int32))
	case P_i32_PDURequest:
		pValue = unsafe.Pointer(new(int32))
	case P_i32_MaxClients:
		pValue = unsafe.Pointer(new(int32))
	default:
		err = fmt.Errorf("paraNumber err")
		return
	}
	var code C.int = C.Srv_GetParam(s.server, C.int(paraNumber), pValue)
	err = Srv_ErrorText(code)
	if err != nil {
		return
	}
	switch paraNumber {
	case P_u16_LocalPort:
		value = *(*uint16)(pValue)
	case P_i32_WorkInterval:
		value = *(*int32)(pValue)
	case P_i32_PDURequest:
		value = *(*int32)(pValue)
	case P_i32_MaxClients:
		value = *(*int32)(pValue)
	}
	return
}

func (s *S7Server) SetParam(paraNumber ParamNumber, value interface{}) (err error) {
	var pValue unsafe.Pointer
	pValue = Value_Pvalue(paraNumber, value)
	var code C.int = C.Srv_SetParam(s.server, C.int(paraNumber), pValue)
	err = Srv_ErrorText(code)
	return
}

// int S7API Srv_StartTo(S7Object Server, const char *Address);
func (s *S7Server) StartTo(Address string) (err error) {
	address := C.CString(Address)
	defer func() {
		C.free(unsafe.Pointer(address))
	}()
	var code C.int = C.Srv_StartTo(s.server, address)
	err = Srv_ErrorText(code)
	return
}

// func Srv_Start(S7Object Server);
func (s *S7Server) Start() (err error) {
	var code C.int = C.Srv_Start(s.server)
	err = Srv_ErrorText(code)
	return
}

// func Srv_Stop(S7Object Server)
func (s *S7Server) Stop() (err error) {
	var code C.int = C.Srv_Stop(s.server)
	err = Srv_ErrorText(code)
	return
}

//typedef uint16_t   word;
// func Srv_RegisterArea(S7Object Server, int AreaCode, word Index, void *pUsrData, int Size)
func (s *S7Server) RegisterArea(AreaCode SrvAreaType, Index uint16, pUsrData []byte) (err error) {
	var code C.int = C.Srv_RegisterArea(
		s.server, C.int(AreaCode), C.uint16_t(Index), unsafe.Pointer(&pUsrData[0]), C.int(len(pUsrData)))
	err = Srv_ErrorText(code)
	return
}

// func Srv_UnregisterArea(S7Object Server, int AreaCode, word Index);
func (s *S7Server) UnregisterArea(AreaCode SrvAreaType, Index uint16) (err error) {
	var code C.int = C.Srv_UnregisterArea(s.server, C.int(AreaCode), C.uint16_t(Index))
	err = Srv_ErrorText(code)
	return
}

// func Srv_LockArea(S7Object Server, int AreaCode, word Index);
func (s *S7Server) LockArea(AreaCode SrvAreaType, Index uint16) (err error) {
	var code C.int = C.Srv_LockArea(s.server, C.int(AreaCode), C.uint16_t(Index))
	err = Srv_ErrorText(code)
	return
}

// func Srv_UnlockArea(S7Object Server, int AreaCode, word Index);
func (s *S7Server) UnlockArea(AreaCode SrvAreaType, Index uint16) (err error) {
	var code C.int = C.Srv_UnlockArea(s.server, C.int(AreaCode), C.uint16_t(Index))
	err = Srv_ErrorText(code)
	return
}

// func Srv_GetStatus(S7Object Server, int *ServerStatus, int *CpuStatus, int *ClientsCount);
func (s *S7Server) GetStatus() (ServerStatus S7ServerStatus, CpuStatus S7CpuStatus, ClientsCount int32, err error) {
	var code C.int = C.Srv_GetStatus(s.server, (*C.int)(unsafe.Pointer(&ServerStatus)), (*C.int)(unsafe.Pointer(&CpuStatus)), (*C.int)(unsafe.Pointer(&ClientsCount)))
	err = Srv_ErrorText(code)
	return
}

// func Srv_SetCpuStatus(S7Object Server, int CpuStatus);
func (s *S7Server) SetCpuStatus(CpuStatus S7CpuStatus) (err error) {
	var code C.int = C.Srv_SetCpuStatus(s.server, C.int(CpuStatus))
	err = Srv_ErrorText(code)
	return
}

// func Srv_ClearEvents(S7Object Server);
func (s *S7Server) ClearEvents() (err error) {
	var code C.int = C.Srv_ClearEvents(s.server)
	err = Srv_ErrorText(code)
	return
}

// func Srv_PickEvent(S7Object Server, TSrvEvent *pEvent, int *EvtReady);
/*
	pEvent == nil && err == nil 表示没有事件发生
*/
func (s *S7Server) PickEvent() (pEvent *TSrvEvent, err error) {
	var EvtReady int
	pEvent = new(TSrvEvent)
	var code C.int = C.Srv_PickEvent(s.server, (*C.TSrvEvent)(unsafe.Pointer(pEvent)), (*C.int)(unsafe.Pointer(&EvtReady)))
	err = Srv_ErrorText(code)
	if err != nil {
		pEvent = nil
		return
	}
	if EvtReady == 0 {
		pEvent = nil
	}
	return
}

// func Srv_GetMask(S7Object Server, int MaskKind, longword *Mask);  uint32_t
func (s *S7Server) GetMask(MaskKind MaskKind) (Mask uint32, err error) {
	var code C.int = C.Srv_GetMask(s.server, C.int(MaskKind), (*C.longword)(unsafe.Pointer(&Mask)))
	err = Srv_ErrorText(code)
	return
}

// func Srv_SetMask(S7Object Server, int MaskKind, longword Mask);
func (s *S7Server) SetMask(MaskKind MaskKind, Mask uint32) (err error) {
	var code C.int = C.Srv_SetMask(s.server, C.int(MaskKind), C.longword(Mask))
	err = Srv_ErrorText(code)
	return
}

// func Srv_EventText(TSrvEvent *Event, char *Text, int TextLen)
func Srv_EventText(Event *TSrvEvent) (text string, err error) {
	const length = 1024
	var buff [length]byte
	var code C.int = C.Srv_EventText(
		(*C.TSrvEvent)(unsafe.Pointer(Event)),
		(*C.char)(unsafe.Pointer(&buff[0])),
		length)
	err = Srv_ErrorText(code)
	if err != nil {
		return
	}
	text = string(buff[:])
	return
}
