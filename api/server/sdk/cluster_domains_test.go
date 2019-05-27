/*
Package sdk is the gRPC implementation of the SDK gRPC server
Copyright 2019 Portworx

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package sdk

import (
	"context"
	"testing"

	"github.com/libopenstorage/gossip/types"
	"github.com/libopenstorage/openstorage/api"
	"github.com/libopenstorage/openstorage/pkg/clusterdomain"
	"github.com/stretchr/testify/assert"
)

func TestSdkClusterDomainsEnumerate(t *testing.T) {

	// Create server and client connection
	s := newTestServer(t)
	defer s.Stop()

	// Create response
	infos := []*clusterdomain.ClusterDomainInfo{
		{
			Name:  "zone1",
			State: types.CLUSTER_DOMAIN_STATE_ACTIVE,
		},
		{
			Name:  "zone2",
			State: types.CLUSTER_DOMAIN_STATE_INACTIVE,
		},
	}
	s.MockCluster().EXPECT().EnumerateDomains().Return(infos, nil).Times(1)

	// Setup client
	c := api.NewOpenStorageClusterDomainsClient(s.Conn())

	// Get info
	r, err := c.Enumerate(context.Background(), &api.SdkClusterDomainsEnumerateRequest{})
	assert.NoError(t, err)
	assert.NotNil(t, r.GetClusterDomainNames())
	assert.Equal(t, 2, len(r.GetClusterDomainNames()), "Unexpected no. of cluster domains")
}

func TestSdkClusterDomainsInspect(t *testing.T) {

	// Create server and client connection
	s := newTestServer(t)
	defer s.Stop()

	// Create response
	info := &clusterdomain.ClusterDomainInfo{
		Name:  "zone1",
		State: types.CLUSTER_DOMAIN_STATE_ACTIVE,
	}
	s.MockCluster().EXPECT().InspectDomain("zone1").Return(info, nil).Times(1)

	// Setup client
	c := api.NewOpenStorageClusterDomainsClient(s.Conn())

	// Get info
	r, err := c.Inspect(context.Background(), &api.SdkClusterDomainInspectRequest{
		ClusterDomainName: "zone1",
	})
	assert.NoError(t, err)
	assert.True(t, r.GetIsActive(), "Unexpected cluster domain status")

	info = &clusterdomain.ClusterDomainInfo{
		Name:  "zone2",
		State: types.CLUSTER_DOMAIN_STATE_INACTIVE,
	}
	s.MockCluster().EXPECT().InspectDomain("zone2").Return(info, nil).Times(1)

	// Get info
	r, err = c.Inspect(context.Background(), &api.SdkClusterDomainInspectRequest{
		ClusterDomainName: "zone2",
	})
	assert.NoError(t, err)
	assert.False(t, r.GetIsActive(), "Unexpected cluster domain status")

}

func TestSdkClusterDomainsActivate(t *testing.T) {

	// Create server and client connection
	s := newTestServer(t)
	defer s.Stop()

	// Setup client
	c := api.NewOpenStorageClusterDomainsClient(s.Conn())

	_, err := c.Activate(context.Background(), &api.SdkClusterDomainActivateRequest{})
	assert.Error(t, err, "Expected an error on empty activate request")

	s.MockCluster().EXPECT().UpdateDomainState("zone2", types.CLUSTER_DOMAIN_STATE_ACTIVE).Return(nil).Times(1)

	// Get info
	_, err = c.Activate(context.Background(), &api.SdkClusterDomainActivateRequest{ClusterDomainName: "zone2"})
	assert.NoError(t, err)

}

func TestSdkClusterDomainsDeactivate(t *testing.T) {

	// Create server and client connection
	s := newTestServer(t)
	defer s.Stop()

	// Setup client
	c := api.NewOpenStorageClusterDomainsClient(s.Conn())

	_, err := c.Deactivate(context.Background(), &api.SdkClusterDomainDeactivateRequest{})
	assert.Error(t, err, "Expected an error on empty deactivate request")

	s.MockCluster().EXPECT().UpdateDomainState("zone2", types.CLUSTER_DOMAIN_STATE_INACTIVE).Return(nil).Times(1)

	// Get info
	_, err = c.Deactivate(context.Background(), &api.SdkClusterDomainDeactivateRequest{ClusterDomainName: "zone2"})
	assert.NoError(t, err)

}
