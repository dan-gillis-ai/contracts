// Package contracts is the Go entry point for this repository's wire contract.
//
// The contract itself is defined in proto/agentharness/v1/agentharness.proto
// and generated into the nested gen/ module, which is what every service in the
// fleet imports (github.com/dan-gillis-ai/contracts/gen/agentharness/v1). This
// package exists for two reasons, and the second is the important one.
//
// First, it is the single place that names the service. Consumers currently
// hard-code "agentharness.v1.DeviceTunnel" wherever they need the fully-qualified
// name; naming it once here means a rename has somewhere to land.
//
// Second, and more concretely: it gives the ROOT module a Go package.
//
// That is not cosmetic. The root module previously contained no Go packages at
// all — only proto/, buf.yaml and the nested gen/ module — so the root-level
// `govulncheck ./...` in .github/workflows/security.yml matched nothing and
// exited 2 ("no packages matched the provided patterns"), failing the
// Vulnerability audit job before it scanned a single line of the generated
// code. Because the repo's Go code lives in a nested module, a root-level
// package pattern cannot reach it; the audit only sees the contract if the root
// module itself has a package, and this is that package.
//
// It deliberately calls BOTH halves of the generated API. govulncheck reports
// at symbol granularity, following calls it can actually reach, so anchoring
// only NewDeviceTunnelClient would leave the server half of the generated
// surface unreached and quietly narrow the audit. Referencing the registration
// helpers too means the scan walks both.
//
// Do not delete this file on the grounds that it adds no behaviour. Its value is
// that the security audit has something real to scan; a root module with no
// packages fails that job.
package contracts

import (
	agentharnessv1 "github.com/dan-gillis-ai/contracts/gen/agentharness/v1"
	"google.golang.org/grpc"
)

// ServiceName is the fully-qualified gRPC service name the proto declares.
//
// It is read from the generated descriptor rather than spelled out as a literal
// here, so it cannot drift from gen/ on its own; contracts_test.go checks it
// against the proto declaration, which is what catches gen/ itself being stale.
//
// Callers that need it — for a gateway route, a service-map entry, a
// reflection probe — should read it from here rather than repeating the
// literal.
var ServiceName = agentharnessv1.DeviceTunnel_ServiceDesc.ServiceName

// RegisterServer binds srv to s under the generated service descriptor.
//
// It is a pass-through with one purpose: to make the server half of the
// generated API a symbol govulncheck can trace from the root module.
func RegisterServer(s grpc.ServiceRegistrar, srv agentharnessv1.DeviceTunnelServer) {
	agentharnessv1.RegisterDeviceTunnelServer(s, srv)
}

// NewClient returns a client for the generated service over cc.
//
// As with RegisterServer, this exists so the client half of the generated API
// is reachable from the root module.
func NewClient(cc grpc.ClientConnInterface) agentharnessv1.DeviceTunnelClient {
	return agentharnessv1.NewDeviceTunnelClient(cc)
}
