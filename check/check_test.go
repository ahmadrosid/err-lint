package check_test

import (
	"err-lint/check"
	"testing"
)

func TestValidateContains(t *testing.T) {
	if check.ContainsCorrectErrHandler("if err != nil {") != true {
		t.Errorf("expected true got false")
	}

	if check.ContainsCorrectErrHandler("return err") != true {
		t.Errorf("expected true got false")
	}

	if check.ContainsCorrectErrHandler("}); err != nil {") != true {
		t.Errorf("expected true got false")
	}

	if check.ContainsCorrectErrHandler("\tif err != nil && err != redis.ErrNil {") != true {
		t.Errorf("expected true got false")
	}

	if check.ContainsCorrectErrHandler("\tif err != nil && strings.Contains(") != true {
		t.Errorf("expected true got false")
	}

	if check.ContainsCorrectErrHandler("return fmt, err") != true {
		t.Errorf("expected true got false")
	}

	if check.ContainsCorrectErrHandler("return (rgreviewInappropriateContentReportDetail)(*resDetail).ToEntity(), nil") != false {
		t.Errorf("expected false got true")
	}
}
