package main

/*
#cgo LDFLAGS: -lktoblzcheck
#include <iban.h>
*/
import "C"

// Iban Check
type IbanCheck struct {
	ptr *C.struct_IbanCheck
}

type Iban struct {
	ptr *C.struct_Iban
}

func NewIbanCheck(dataDir string) IbanCheck {
	return IbanCheck{
		ptr: C.IbanCheck_new(C.CString(dataDir)),
	}
}
func (check *IbanCheck) Free() {
	C.IbanCheck_free(check.ptr)
}

func (check *IbanCheck) Check(Iban string, country string) int {
	return int(C.IbanCheck_check_str(check.ptr, C.CString(Iban), C.CString(country)))
}

func (check *IbanCheck) BicPosition(Iban string) (start int, end int) {
	var cstart C.int
	var cend C.int

	C.IbanCheck_bic_position(check.ptr, C.CString(Iban), &cstart, &cend)
	return int(cstart), int(cend)
}

func NewIban(iban string) Iban {
	return Iban{
		ptr: C.Iban_new(C.CString(iban), 1),
	}
}

func (check *Iban) Free() {
	C.Iban_free(check.ptr)
}

func (check *Iban) TransmissionForm() string {
	return C.GoString(C.Iban_transmissionForm(check.ptr))
}

func (check *Iban) Printable() string {
	return C.GoString(C.Iban_printableForm(check.ptr))
}
