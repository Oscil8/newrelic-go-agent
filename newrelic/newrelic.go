// +build linux

// Go wrapper for the New Relic Agent SDK. Requires Linux and the SDK headers
// and libraries.
package newrelic

/*
#cgo LDFLAGS: -L/usr/local/lib -lnewrelic-collector-client -lnewrelic-common -lnewrelic-transaction
#include "nr_agent_sdk/include/newrelic_collector_client.h"
#include "nr_agent_sdk/include/newrelic_common.h"
#include "nr_agent_sdk/include/newrelic_transaction.h"
#include "stdlib.h"
*/
import "C"

import (
	"errors"
	"fmt"
	"unsafe"
)

var statusMap = map[int]string{
	0: "ok",
	-0x10001: "other",
	-0x20001: "disabled",
	-0x30001: "invalid param",
	-0x30002: "invalid id",
	-0x40001: "transaction not started",
	-0x40002: "transaction in progress",
	-0x40003: "transaction not named",
}

func nrError(i C.int, name string) error {
	if int(i) < 0 {
		status, ok := statusMap[int(i)]
		if !ok {
			status = "unknown"
		}
		return errors.New(fmt.Sprintf("newrelic: %s: %s", name, status))
	}
	return nil
}

func Init(license string, appName string, lang string, langVersion string) error {
	C.newrelic_register_message_handler((*[0]byte)(C.newrelic_message_handler))
	clicense := C.CString(license)
	defer C.free(unsafe.Pointer(clicense))
	cappName := C.CString(appName)
	defer C.free(unsafe.Pointer(cappName))
	clang := C.CString(lang)
	defer C.free(unsafe.Pointer(clang))
	clangVersion := C.CString(langVersion)
	defer C.free(unsafe.Pointer(clangVersion))
	rv := C.newrelic_init(clicense, cappName, clang, clangVersion)
	return nrError(rv, "initialize")
}

func RequestShutdown(reason string) error {
	ptr := C.CString(reason)
	defer C.free(unsafe.Pointer(ptr))
	rv := C.newrelic_request_shutdown(ptr)
	return nrError(rv, "request shutdown")
}

func BeginTransaction() int64 {
	id := C.newrelic_transaction_begin()
	return int64(id)
}

func SetTransactionTypeWeb(txnID int64) error {
	rv := C.newrelic_transaction_set_type_web(C.long(txnID))
	return nrError(rv, "set transaction type web")
}

func SetTransactionTypeOther(txnID int64) error {
	rv := C.newrelic_transaction_set_type_other(C.long(txnID))
	return nrError(rv, "set transaction type other")
}

func SetTransactionName(txnID int64, name string) error {
	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cname))
	rv := C.newrelic_transaction_set_name(C.long(txnID), cname)
	return nrError(rv, "set transaction name")
}

func SetTransactionCategory(txnID int64, category string) error {
	ccategory := C.CString(category)
	defer C.free(unsafe.Pointer(ccategory))
	rv := C.newrelic_transaction_set_category(C.long(txnID), ccategory)
	return nrError(rv, "set transaction category")
}

func NoticeTransactionError(
	txnID int64,
	exception_type,
	error_message,
	stack_trace,
	delimiter string) error {
	cexception_type := C.CString(exception_type)
	defer C.free(unsafe.Pointer(cexception_type))
	cerror_message := C.CString(error_message)
	defer C.free(unsafe.Pointer(cerror_message))
	cstack_trace := C.CString(stack_trace)
	defer C.free(unsafe.Pointer(cstack_trace))
	cdelimiter := C.CString(delimiter)
	defer C.free(unsafe.Pointer(cdelimiter))
	rv := C.newrelic_transaction_notice_error(C.long(txnID),
	  cexception_type, cerror_message, cstack_trace, cdelimiter)
	return nrError(rv, "notice transaction error")
}

func AddTransactionAttribute(txnID int64, name, value string) error {
	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cname))
	cvalue := C.CString(value)
	defer C.free(unsafe.Pointer(cvalue))
	rv := C.newrelic_transaction_add_attribute(C.long(txnID), cname, cvalue)
	return nrError(rv, "add transaction attribute")
}

func BeginGenericSegment(txnID int64, parentID int64, name string) int64 {
	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cname))
	id := C.newrelic_segment_generic_begin(C.long(txnID), C.long(parentID), cname)
	return int64(id)
}

func BeginDatastoreSegment(
	txnID int64,
	parentID int64,
	table string,
	operation string,
	sql string,
	rollupName string) int64 {
	ctable := C.CString(table)
	defer C.free(unsafe.Pointer(ctable))
	coperation := C.CString(operation)
	defer C.free(unsafe.Pointer(coperation))
	csql := C.CString(sql)
	defer C.free(unsafe.Pointer(csql))
	crollupName := C.CString(rollupName)
	defer C.free(unsafe.Pointer(crollupName))
	id := C.newrelic_segment_datastore_begin(
		C.long(txnID),
		C.long(parentID),
		ctable,
		coperation,
		csql,
		crollupName,
		(*[0]byte)(C.newrelic_basic_literal_replacement_obfuscator),
	)
	return int64(id)
}

func BeginExternalSegment(txnID int64, parentID int64, host string, name string) int64 {
	chost := C.CString(host)
	defer C.free(unsafe.Pointer(chost))
	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cname))
	id := C.newrelic_segment_external_begin(C.long(txnID), C.long(parentID), chost, cname)
	return int64(id)
}

func EndSegment(txnID int64, parentID int64) error {
	rv := C.newrelic_segment_end(C.long(txnID), C.long(parentID))
	return nrError(rv, "end segment")
}

func SetTransactionRequestURL(txnID int64, url string) error {
	curl := C.CString(url)
	defer C.free(unsafe.Pointer(curl))
	rv := C.newrelic_transaction_set_request_url(C.long(txnID), curl)
	return nrError(rv, "set transaction request url")
}

func EndTransaction(txnID int64) error {
	rv := C.newrelic_transaction_end(C.long(txnID))
	return nrError(rv, "end transaction")
}

func RecordMetric(name string, val float64) error {
	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cname))
	rv := C.newrelic_record_metric(cname, C.double(val))
	return nrError(rv, "record metric")
}
