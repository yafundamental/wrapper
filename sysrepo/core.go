package sysrepo

/*
#cgo CFLAGS: -I$/opt/dev/sysrepo/src -I$/sandbox/sysrepo -I$/usr/include/libyang/libyang.h
#cgo LDFLAGS: -L$/opt/dev/sysrepo/src -L$/sandbox/sysrepo -L$/usr/include/libyang
#cgo LDFLAGS: -lsysrepo -lyang -lpcre2-8 -lpcre2-16 -lpcre2-32 -lpcre2-posix -ldl
#include "lib.h"
*/
import "C"

import (
	"fmt"
	"unsafe"
)

// Ядро враппера. Слой, ближе всего лежащий к Sysrepo.
type Core struct{}

func (c Core) Connect(
	opts C.sr_conn_options_t,
	conn **C.sr_conn_ctx_t,
) error {
	rc := C.sr_connect(opts, conn)

	return ParseError(rc)
}

func (c Core) Disconnect(
	conn *C.sr_conn_ctx_t,
) error {
	rc := C.sr_disconnect(conn)

	return ParseError(rc)
}

func (c Core) SessionStart(
	conn *C.sr_conn_ctx_t,
	datastore C.sr_datastore_t,
	session **C.sr_session_ctx_t,
) error {
	rc := C.sr_session_start(conn, datastore, session)

	return ParseError(rc)
}

func (c Core) SessionStop(
	session *C.sr_session_ctx_t,
) error {
	rc := C.sr_session_stop(session)

	return ParseError(rc)
}

func (c Core) GetData(
	session *C.sr_session_ctx_t,
	xpath *C.char,
	maxDepth C.uint32_t,
	timeoutMs C.uint32_t,
	opts C.sr_get_options_t,
) (*Data, error) {
	var data *C.sr_data_t
	rc := C.sr_get_data(session, xpath, maxDepth, timeoutMs, opts, &data)

	if err := ParseError(rc); err != nil {
		return nil, err
	}

	if data == nil {
		return nil, ErrNodeNotFound
	}

	return (*Data)(data), nil

}

func (c Core) SetItemStr(
	session *C.sr_session_ctx_t,
	xpath *C.char,
	value *C.char,
	origin *C.char,
	opts C.uint32_t,
) error {
	rc := C.sr_set_item_str(session, xpath, value, origin, opts)

	return ParseError(rc)
}

func (c Core) DeleteItem(
	session *C.sr_session_ctx_t,
	xpath *C.char,
) error {
	rc := C.sr_delete_item(session, xpath, 0)

	return ParseError(rc)
}

func (c Core) Commit(
	session *C.sr_session_ctx_t,
	timeoutMs C.uint32_t,
) error {
	rc := C.sr_apply_changes(session, timeoutMs)

	return ParseError(rc)
}

func (c Core) LydNewPath(
	parent *C.lyd_node_s,
	ctx *C.ly_ctx_s,
	path *C.char,
	value *C.char,
	opts C.uint32_t,
) (*LibYangNode, error) {
	var node *C.lyd_node_s
	rc := C.lyd_new_path(parent, ctx, path, value, opts, &node)
	if err := ParseYangError(rc); err != nil {
		return nil, err
	}

	return &LibYangNode{
		Node: node,
	}, nil
}

func (c Core) NewYangContext(conn *C.sr_conn_ctx_t) *LibYangContext {
	return &LibYangContext{
		context: C.sr_acquire_context(conn),
	}
}

func (c Core) NotificationsSubscribe(
	session *C.sr_session_ctx_t,
	moduleName *C.char,
	xpath *C.char,
	startTime *C.timespec_s,
	stopTime *C.timespec_s,
	callback C.sr_event_notif_cb,
	privateData unsafe.Pointer,
	opts C.uint32_t,
	subscription **C.sr_subscription_ctx_t,
) error {
	rc := C.sr_notif_subscribe(
		session,
		moduleName,
		xpath,
		startTime,
		stopTime,
		callback,
		privateData,
		opts,
		subscription,
	)

	return ParseError(rc)
}

func (c Core) NotificationSend(
	session *C.sr_session_ctx_t,
	notification *C.lyd_node_s,
	timeoutMs C.uint32_t,
	wait C.int,
) error {
	fmt.Println("gg")
	rc := C.sr_notif_send_tree(session, notification, timeoutMs, wait)

	return ParseError(rc)
}

func (c Core) NotificationUnsubscribe(
	subscription *C.sr_subscription_ctx_t,
) error {
	rc := C.sr_unsubscribe(subscription)

	return ParseError(rc)
}
