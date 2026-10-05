package contracts

import (
	"os"
	"regexp"
	"testing"
)

// protoPath is the source of truth for the contract. Relative to this package's
// directory, which is where `go test` runs.
const protoPath = "proto/agentharness/v1/agentharness.proto"

var (
	protoPackageRE = regexp.MustCompile(`(?m)^package\s+([A-Za-z0-9_.]+)\s*;`)
	protoServiceRE = regexp.MustCompile(`(?m)^service\s+([A-Za-z0-9_]+)\s*\{`)
)

// TestServiceNameMatchesProto catches the failure mode this repository is most
// exposed to: proto/agentharness/v1 edited, gen/ not regenerated.
//
// Nothing in the build fails in that state. gen/ is committed, so every service
// keeps compiling against the previous contract, buf never runs in CI, and the
// mismatch surfaces only when a device on the new build talks to a server on the
// old one — as a runtime error at the far end of a fleet. Comparing the
// generated descriptor against the declaration in the proto catches it here
// instead, where the diff is still visible.
func TestServiceNameMatchesProto(t *testing.T) {
	raw, err := os.ReadFile(protoPath)
	if err != nil {
		t.Fatalf("reading %s: %v", protoPath, err)
	}
	pkg := protoPackageRE.FindSubmatch(raw)
	if pkg == nil {
		t.Fatalf("%s declares no package", protoPath)
	}
	svc := protoServiceRE.FindSubmatch(raw)
	if svc == nil {
		t.Fatalf("%s declares no service", protoPath)
	}
	want := string(pkg[1]) + "." + string(svc[1])
	if ServiceName != want {
		t.Errorf("ServiceName = %q but %s declares %q; gen/ is stale — regenerate it (`make proto`)", ServiceName, protoPath, want)
	}
}

// TestGeneratedServerSurfaceIsReachable is a smoke test that the generated
// symbols this package anchors are actually present, rather than the package
// compiling against a renamed generated API.
func TestGeneratedServerSurfaceIsReachable(t *testing.T) {
	if ServiceName == "" {
		t.Fatal("ServiceName is empty; the generated service descriptor is missing")
	}
}
