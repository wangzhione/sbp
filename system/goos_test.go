package system

import "testing"

func TestGOOS(t *testing.T) {
	// 平台 bool 常量同一目标平台只能有一个为 true。
	flags := map[string]bool{
		"aix":       Aix,
		"android":   Android,
		"darwin":    Darwin,
		"dragonfly": Dragonfly,
		"freebsd":   Freebsd,
		"hurd":      Hurd,
		"illumos":   Illumos,
		"ios":       Ios,
		"js":        Js,
		"linux":     Linux,
		"nacl":      Nacl,
		"netbsd":    Netbsd,
		"openbsd":   Openbsd,
		"plan9":     Plan9,
		"solaris":   Solaris,
		"wasip1":    Wasip1,
		"windows":   Windows,
		"zos":       Zos,
	}

	count := 0
	for _, value := range flags {
		if value {
			count++
		}
	}
	if count != 1 {
		t.Fatalf("enabled GOOS flags = %d, want 1", count)
	}
}
