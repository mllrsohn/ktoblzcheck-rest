package main

/*
#cgo LDFLAGS: -lktoblzcheck
#include <ktoblzcheck.h>
*/
import "C"

import (
	"fmt"
	"strconv"
)

type AccountNumberCheck struct {
	ptr *C.struct_AccountNumberCheck
}

type AccountNumberResult int

type Record struct {
	BankID   string
	Name     string
	Location string
}

type UnknownRecord string

func (f UnknownRecord) Error() string {
	return fmt.Sprintf("unknown BLZ: %v", string(f))
}

func (check *AccountNumberCheck) FindBank(bankID string) (Record, error) {
	bank := C.AccountNumberCheck_findBank(check.ptr, C.CString(bankID))
	if bank == nil {
		return Record{}, UnknownRecord(bankID)
	}
	var record = Record{
		BankID: strconv.Itoa(int(C.AccountNumberCheck_Record_bankId(bank))),
	}
	if StringEncoding() == "UTF-8" {
		record.Name = C.GoString(C.AccountNumberCheck_Record_bankName(bank))
		record.Location = C.GoString(C.AccountNumberCheck_Record_location(bank))
	} else {
		record.Name = toUtf8([]byte(C.GoString(C.AccountNumberCheck_Record_bankName(bank))))
		record.Location = toUtf8([]byte(C.GoString(C.AccountNumberCheck_Record_location(bank))))
	}
	return record, nil
}

func (check *AccountNumberCheck) Check(bankID string, accountID string) AccountNumberResult {
	return AccountNumberResult(C.AccountNumberCheck_check(check.ptr, C.CString(bankID), C.CString(accountID)))
}

func StringEncoding() string {
	return C.GoString(C.AccountNumberCheck_stringEncoding())
}

func NewAccountNumberCheck(dataDir string) AccountNumberCheck {
	return AccountNumberCheck{
		ptr: C.AccountNumberCheck_new_file(C.CString(dataDir)),
	}
}

func NewDefaultAccountNumberCheck() AccountNumberCheck {
	return AccountNumberCheck{
		ptr: C.AccountNumberCheck_new(),
	}
}

func toUtf8(iso8859Buf []byte) string {
	buf := make([]rune, len(iso8859Buf))
	for i, b := range iso8859Buf {
		buf[i] = rune(b)
	}
	return string(buf)
}
