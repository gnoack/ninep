package ninep

import "testing"

func TestModeString(t *testing.T) {
	for mode, want := range map[uint32]string{
		ModeDir:        "d---------",
		ModeUserRead:   "-r--------",
		ModeUserWrite:  "--w-------",
		ModeUserExec:   "---x------",
		ModeGroupRead:  "----r-----",
		ModeGroupWrite: "-----w----",
		ModeGroupExec:  "------x---",
		ModeOtherRead:  "-------r--",
		ModeOtherWrite: "--------w-",
		ModeOtherExec:  "---------x",
	} {
		got := ModeString(mode)
		if got != want {
			t.Errorf("bad mode string for mode 0x%x, got: %v, want: %v", mode, got, want)
		}
	}
}
