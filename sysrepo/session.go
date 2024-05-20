package sysrepo

/*
#include "lib.h"
*/
import "C"
import (
	"fmt"
	"unsafe"
)

type Session struct {
	conn    *Connection
	session *C.sr_session_ctx_t
}

func NewSession(conn *Connection) (*Session, error) {
	core := Core{}

	var session *C.sr_session_ctx_t
	if err := core.SessionStart(conn.conn, C.SR_DS_RUNNING, &session); err != nil {
		return nil, fmt.Errorf("failed to start session: %w", err)
	}

	return &Session{
		conn:    conn,
		session: session,
	}, nil
}

func (s *Session) Stop() error {
	core := Core{}

	return core.SessionStop(s.session)
}

func (s *Session) GetDataByXpath(xpath string) (*Data, error) {
	core := Core{}

	cxpath := C.CString(xpath)
	defer C.free(unsafe.Pointer(cxpath))

	data, err := core.GetData(
		s.session,
		cxpath,
		C.uint32_t(0), // default depth
		C.uint32_t(0), // default timeout
		C.uint32_t(0), // default options
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get data by xpath: %w", err)
	}

	return data, nil
}

func (s *Session) SetItemStr(xpath string, value string) error {
	core := Core{}

	cxpath := C.CString(xpath)
	defer C.free(unsafe.Pointer(cxpath))

	cvalue := C.CString(value)
	defer C.free(unsafe.Pointer(cvalue))

	return core.SetItemStr(s.session, cxpath, cvalue, nil, 0)
}

func (s *Session) DeleteItem(xpath string) error {
	core := Core{}

	cxpath := C.CString(xpath)
	defer C.free(unsafe.Pointer(cxpath))

	return core.DeleteItem(s.session, cxpath)
}

func (s *Session) Commit() error {
	core := Core{}

	return core.Commit(s.session, 0)
}

func (s *Session) NotificationsSubscribe(
	module string,
	xpath string,
	callback EventNotificationCallback,
) (*NotificationSubscription, error) {
	core := Core{}

	if err := EventNotificationCallbackRegister(xpath, s, callback); err != nil {
		return nil, err
	}

	cmodule := C.CString(module)
	defer C.free(unsafe.Pointer(cmodule))

	cxpath := C.CString(xpath)
	defer C.free(unsafe.Pointer(cxpath))

	var subscription *C.sr_subscription_ctx_t
	err := core.NotificationsSubscribe(
		s.session,
		cmodule,
		cxpath,
		nil,
		nil,
		C.sr_event_notif_cb(C._sr_event_notif_cb),
		nil,
		C.uint32_t(0),
		&subscription,
	)
	if err != nil {
		return nil, err
	}

	return &NotificationSubscription{
		session: s,
		context: subscription,
		xpath:   xpath,
	}, nil
}

func (s *Session) NotificationSend(node *LibYangNode) error {
	core := Core{}

	return core.NotificationSend(s.session, node.Node, 0, 0)
}

func (s *Session) NewYangContext() (*LibYangContext, error) {
	core := Core{}

	ctx := core.NewYangContext(s.conn.conn)
	if ctx == nil {
		return nil, fmt.Errorf("Failed to NewYangContext")
	}

	return &LibYangContext{
		context: ctx.context,
	}, nil
}

func (s *Session) CloneWithContext(session *C.sr_session_ctx_t) *Session {
	return &Session{
		conn:    s.conn,
		session: session,
	}
}
