package model

import (
"reflect"


apiv3 "github.com/projectcalico/libcalico-go/lib/apis/v3"
"github.com/projectcalico/libcalico-go/lib/namespace"
"regexp"
)




func init() {
	registerResourceInfo(
		apiv3.KindBGPPeer,
		"bgppeers",
		reflect.TypeOf(apiv3.BGPPeer{}),
	)

}
type IpamResourceKey struct {
	// The name of the resource.
	Key string

}

func (key IpamResourceKey) defaultPath() (string, error) {
	return key.defaultDeletePath()
}

func (key IpamResourceKey) defaultDeletePath() (string, error) {

	return key.Key, nil
}

func (key IpamResourceKey) defaultDeleteParentPaths() ([]string, error) {
	return nil, nil
}

func (key IpamResourceKey) valueType() (reflect.Type, error) {

	return reflect.TypeOf(rawString("")), nil
}

func (key IpamResourceKey) String() string {

	return key.Key
}

type IpamResourceListOptions struct {
	// The name of the resource.
	Name string
	// The namespace of the resource.  Not required if the resource is not namespaced.
	Namespace string
	// The resource kind.
	Kind string
	// Whether the name is prefix rather than the full name.
	Prefix bool
}

// If the Kind, Namespace and Name are specified, but the Name is a prefix then the
// last segment of this path is a prefix.
func (options IpamResourceListOptions) IsLastSegmentIsPrefix() bool {
	return len(options.Kind) != 0 &&
		(len(options.Namespace) != 0 || !namespace.IsNamespaced(options.Kind)) &&
		len(options.Name) != 0 &&
		options.Prefix
}

func (options IpamResourceListOptions) KeyFromDefaultPath(path string) Key {
	if path == "/calico/ipam/v2/host" {
		return nil
	}
	reg := regexp.MustCompile(`^/calico/ipam/v2/host/(.*)/ipv4/block/(.*)$`)
	params := reg.FindStringSubmatch(path)

	return ResourceKey{
		params[1]+":"+params[2],
		"",
		"",
	}
}

func (options IpamResourceListOptions) defaultPathRoot() string {
	return "/calico/ipam/v2/host"
	/*ri, ok := resourceInfoByKind[strings.ToLower(options.Kind)]
	if !ok {
		log.Panic("Unexpected resource kind: " + options.Kind)
	}

	k := "/calico/resources/v3/projectcalico.org/" + ri.plural
	if namespace.IsNamespaced(options.Kind) {
		if options.Namespace == "" {
			return k
		}
		k = k + "/" + options.Namespace
	}
	if options.Name == "" {
		return k
	}
	return k + "/" + options.Name*/
}

func (options IpamResourceListOptions) String() string {
	return options.Kind
}
