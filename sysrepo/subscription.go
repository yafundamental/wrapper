package sysrepo

/*
#include "lib.h"
*/
import "C"

type NotificationSubscription struct {
	session *Session
	context *C.sr_subscription_ctx_t
	xpath   string
}

func (s *NotificationSubscription) Unsubscribe() error {
	core := Core{}

	return core.NotificationUnsubscribe(s.context)
}
