package notifications

import "testing"

func TestVoucherAttachments(t *testing.T) {
	for _, tc := range []struct {
		includeERP bool
		want       []string
	}{
		{includeERP: false, want: []string{TermsFileName}},
		{includeERP: true, want: []string{TermsFileName, ERPFileName}},
	} {
		attachments, err := voucherAttachments(tc.includeERP)
		if err != nil {
			t.Fatalf("voucherAttachments(%v): %v", tc.includeERP, err)
		}

		var got []string
		for _, a := range attachments {
			got = append(got, a.Filename)
		}
		if len(got) != len(tc.want) || got[0] != tc.want[0] || (len(got) == 2 && got[1] != tc.want[1]) {
			t.Errorf("voucherAttachments(%v) = %v, want %v", tc.includeERP, got, tc.want)
		}
	}
}
